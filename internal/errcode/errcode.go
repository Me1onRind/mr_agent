package errcode

const (
	SuccessCode int = 0

	JsonEncodeFailedCode int = -10000000
	JsonDecodeFailedCode int = -10000001
	ParamInvalidCode     int = -10000002

	UnexpectCode int = -999999999
)

func IsWarning(code int) bool {
	return code > SuccessCode
}

var (
	ErrJsonEncodeFailed = NewError(JsonEncodeFailedCode, "Json Encode Failed")
	ErrJsonDecodeFailed = NewError(JsonDecodeFailedCode, "Json Decode Failed")
	ErrParamInvalid     = NewError(ParamInvalidCode, "Param Invalid")
)
