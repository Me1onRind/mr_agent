package tools

import (
	"context"
	"errors"
	"log/slog"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/Me1onRind/mr_agent/internal/infrastructure/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/openai/openai-go/v3"
)

var (
	ErrToolExists   = errors.New("tool already exists")
	ErrToolNotFound = errors.New("tool not found")
	ErrToolInvalid  = errors.New("tool definition invalid")
)

type Definition struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Parameters  map[string]any    `json:"parameters,omitempty"`
	Meta        map[string]string `json:"meta,omitempty"`
}

type Handler func(ctx context.Context, args jsoniter.RawMessage) (any, error)
type handler[A any, B any] func(ctx context.Context, input *A) (*B, error)
type generateToolFunc func() (*Tool, error)

func HandlerWrap[A any, B any](f handler[A, B]) Handler {
	return func(ctx context.Context, args jsoniter.RawMessage) (any, error) {
		var input A
		if err := jsoniter.Unmarshal(args, &input); err != nil {
			return nil, err
		}

		return f(ctx, &input)
	}
}

type Tool struct {
	Definition Definition
	Handler    Handler
}

type Registry struct {
	mu    sync.RWMutex
	tools map[string]Tool
}

func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]Tool),
	}
}

var DefaultRegistry = NewRegistry()

func InitAgentTools(ctx context.Context) error {
	log := logger.LoggerFromCtx(ctx)
	generateToolFuncs := []generateToolFunc{
		excuteSQL,
	}
	for _, f := range generateToolFuncs {
		tool, err := f()
		pc := reflect.ValueOf(f).Pointer()
		name := filepath.Base(runtime.FuncForPC(pc).Name())
		log.Info("generateTool", slog.String("func_name", name))
		if err != nil {
			log.Error("generateToolFuncs failed", slog.String("func_name", name))
			return err
		}
		log.Info("register add agent tool", slog.String("tool", tool.Definition.Name))
		if err := DefaultRegistry.Add(tool); err != nil {
			log.Error("register add agent tool failed", slog.String("tool", tool.Definition.Name))
			return err
		}
	}
	return nil
}

func (r *Registry) Add(tool *Tool) error {
	if tool.Definition.Name == "" {
		return ErrToolInvalid
	}
	if tool.Handler == nil {
		return ErrToolInvalid
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.tools[tool.Definition.Name]; exists {
		return ErrToolExists
	}
	r.tools[tool.Definition.Name] = *tool
	return nil
}

func (r *Registry) Get(name string) (Tool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tool, ok := r.tools[name]
	if !ok {
		return Tool{}, ErrToolNotFound
	}
	return tool, nil
}

func (r *Registry) Definitions() []Definition {
	r.mu.RLock()
	defer r.mu.RUnlock()
	defs := make([]Definition, 0, len(r.tools))
	for _, tool := range r.tools {
		defs = append(defs, tool.Definition)
	}
	sort.Slice(defs, func(i, j int) bool {
		return defs[i].Name < defs[j].Name
	})
	return defs
}

func (r *Registry) OpenAITools() []openai.ChatCompletionToolUnionParam {
	defs := r.Definitions()
	tools := make([]openai.ChatCompletionToolUnionParam, 0, len(defs))
	for _, def := range defs {
		tools = append(tools, def.ToOpenAITool())
	}
	return tools
}

type PropertyValue struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

func structToPropertyValues(input any) (map[string]PropertyValue, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}

	t := reflect.TypeOf(input)
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, errors.New("input is not a struct")
	}

	properties := make(map[string]PropertyValue, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue
		}

		jsonName, skip := jsonFieldName(field)
		if skip {
			continue
		}

		prop := PropertyValue{
			Type:        jsonSchemaType(field.Type),
			Description: field.Tag.Get("description"),
		}
		if prop.Description == "" {
			prop.Description = field.Tag.Get("desc")
		}

		properties[jsonName] = prop
	}

	return properties, nil
}

func jsonFieldName(field reflect.StructField) (string, bool) {
	tag := field.Tag.Get("json")
	if tag == "-" {
		return "", true
	}

	if tag == "" {
		return field.Name, false
	}

	parts := strings.Split(tag, ",")
	if len(parts) == 0 || parts[0] == "" {
		return field.Name, false
	}
	return parts[0], false
}

func jsonSchemaType(t reflect.Type) string {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "integer"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Slice, reflect.Array:
		if t.Elem().Kind() == reflect.Uint8 {
			return "string"
		}
		return "array"
	case reflect.Map, reflect.Struct, reflect.Interface:
		if t.PkgPath() == "time" && t.Name() == "Time" {
			return "string"
		}
		return "object"
	default:
		return "object"
	}
}

func newObjectSchema(properties map[string]PropertyValue, required ...string) map[string]any {
	schema := map[string]any{
		"type":       "object",
		"Prorerties": properties,
	}
	if len(required) > 0 {
		schema["required"] = required
	}
	return schema
}

func (d Definition) ToOpenAITool() openai.ChatCompletionToolUnionParam {
	fd := openai.FunctionDefinitionParam{
		Name: d.Name,
	}
	if d.Description != "" {
		fd.Description = openai.String(d.Description)
	}
	if d.Parameters != nil {
		fd.Parameters = openai.FunctionParameters(d.Parameters)
	}
	return openai.ChatCompletionFunctionTool(fd)
}

func (t Tool) ToOpenAITool() openai.ChatCompletionToolUnionParam {
	return t.Definition.ToOpenAITool()
}

func generateTool[A any, B any](name, desc string, f handler[A, B]) (*Tool, error) {
	var a A
	propertyValue, err := structToPropertyValues(&a)
	if err != nil {
		return nil, err
	}
	tool := Tool{
		Definition: Definition{
			Name:        name,
			Description: desc,
			Parameters:  newObjectSchema(propertyValue),
		},
		Handler: HandlerWrap(f),
	}
	return &tool, nil
}
