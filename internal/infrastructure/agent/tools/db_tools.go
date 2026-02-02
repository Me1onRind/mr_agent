package tools

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/Me1onRind/mr_agent/internal/config"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/db"
	"github.com/Me1onRind/mr_agent/internal/protocol/agent"
)

func excuteSQL() (*Tool, error) {
	return generateTool(
		"excute_sql", "excute sql", ExcuteSQLHander,
	)
}

func ExcuteSQLHander(ctx context.Context, input *agent.ExcuteSQLInput) (*agent.ExcuteSQLOutput, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}
	sqlText := strings.TrimSpace(input.SQL)
	if sqlText == "" {
		return nil, errors.New("sql is empty")
	}

	label := config.DBLabel(input.DBLabel)
	if isQuerySQL(sqlText) {
		rows, err := queryRows(ctx, label, sqlText)
		if err != nil {
			return &agent.ExcuteSQLOutput{Error: err.Error()}, nil
		}
		return &agent.ExcuteSQLOutput{Rows: rows}, nil
	}

	result := db.GetMasterDB(ctx, label).Exec(sqlText)
	output := &agent.ExcuteSQLOutput{
		RowsAffected: result.RowsAffected,
	}
	if result.Error != nil {
		output.Error = result.Error.Error()
		return output, nil
	}
	return output, nil
}

func isQuerySQL(sqlText string) bool {
	trimmed := strings.TrimSpace(sqlText)
	if trimmed == "" {
		return false
	}
	first := strings.ToLower(strings.Fields(trimmed)[0])
	switch first {
	case "select", "show", "describe", "desc", "explain", "with":
		return true
	default:
		return false
	}
}

func queryRows(ctx context.Context, label config.DBLabel, sqlText string) ([]map[string]any, error) {
	query := sqlText
	if !hasLimitClause(query) {
		query = strings.TrimRight(query, ";")
		query = query + " LIMIT 100"
	}
	rows, err := db.GetSlaveDB(ctx, label).Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	results := make([]map[string]any, 0)
	for rows.Next() {
		values := make([]any, len(cols))
		valuePtrs := make([]any, len(cols))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}
		rowMap := make(map[string]any, len(cols))
		for i, col := range cols {
			rowMap[col] = normalizeSQLValue(values[i], colTypes[i])
		}
		results = append(results, rowMap)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func hasLimitClause(sqlText string) bool {
	lower := strings.ToLower(sqlText)
	return strings.Contains(lower, " limit ")
}

func normalizeSQLValue(value any, colType *sql.ColumnType) any {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case []byte:
		if colType != nil {
			name := strings.ToLower(colType.DatabaseTypeName())
			if strings.Contains(name, "blob") || strings.Contains(name, "binary") {
				return v
			}
		}
		return string(v)
	default:
		return v
	}
}
