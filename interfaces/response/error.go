package response

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"strconv"

	"github.com/labstack/echo"
)

// Error はエラー情報のフォーマットを定義する
type Error struct {
	ErrorCode  string
	StatusCode int
	Message    string
}

var errors = map[string]Error{
	// 認証・認可エラー
	"AUTHENTICATION_ERROR": {"1001", 400, "Authentication Error"},
	"AUTHORIZATION_ERROR":  {"1002", 403, "Authorization Error"},
	"INVALID_TOKEN":        {"1003", 401, "Invalid Token"},
	"TOKEN_GENERATE_ERROR": {"1004", 500, "Token Generate Error"},
	// バリデーションエラー
	"METHOD_NOT_ALLOWED": {"1101", 405, "Method Not Allowed"},
	"INVALID_FORMAT":     {"1102", 400, "Invalid Request Format"},
	"INVALID_PARAMETER":  {"1103", 400, "Invalid Request Parameter"},
	// DB エラー
	"DB_CONNECTION_ERROR":    {"2001", 500, "DataBase Connection Failed"},
	"DB_REQUEST_ERROR":       {"2002", 500, "Invalid Query Requested For DataBase"},
	"DB_CONNECTION_TIME_OUT": {"2003", 504, "DataBase Connection Timed Out"},
	"DB_DUPLICATE_ENTRY":     {"2004", 500, "Duplicate key entry"},
	"DB_NOT_FOUND":           {"2005", 404, "Requested records not found"},
	// ファイル IO エラー
	"OPEN_FILE_ERROR":  {"2101", 500, "Open File Error"},
	"READ_FILE_ERROR":  {"2102", 500, "Read File Error"},
	"WRITE_FILE_ERROR": {"2103", 504, "Write File Error"},
	// 外部システム IF エラー
	"CONNECTION_ERROR":    {"2201", 503, "Connection Failed"},
	"CONNECTION_TIME_OUT": {"2202", 504, "Connection Timed Out"},
	"RESPONSE_ERROR":      {"2203", 502, "Unexpected Response"},
	// 内部システム IF エラー
	"CACHE_SERVER_ERROR": {"2301", 500, "Cache Server Error"},
	// その他エラー
	"302_FOUND":              {"9901", 302, "Found"},
	"API_VIRSION_MISMATCH":   {"9902", 412, "API Version Mismatched"},
	"INVALID_API_VIRSION":    {"9903", 412, "Invalid API Version"},
	"INVALID_CLIENT_VIRSION": {"9904", 412, "Unsupported Cilent Version"},
	"UNDER_MAINTENANCE":      {"9905", 503, "Under Maintenance"},
	"NOT_IMPLEMENTED":        {"9906", 501, "Not Implemented"},
	"INTERNAL_SERVER_ERROR":  {"9907", 500, "Internal Server Error"},
	"UNKNOWN_ERROR":          {"9999", 503, "Unknown Error"},
}

// ErrorResponse return error json response
func ErrorResponse(c echo.Context, ErrorID string, Detail string) error {
	errorObject := errors[ErrorID]
	errorJSON := map[string]string{
		"code":    errorObject.ErrorCode,
		"message": errorObject.Message,
		"detail":  Detail,
	}

	pc, _, line, _ := runtime.Caller(1)
	called := runtime.FuncForPC(pc).Name()
	errLog := map[string]string{
		"called":  called + ":" + strconv.Itoa(line),
		"code":    errorObject.ErrorCode,
		"message": errorObject.Message,
		"detail":  Detail,
	}
	errLogJSON, _ := json.Marshal(errLog)

	log.Println(fmt.Sprintf("%v: %v", "error", string(errLogJSON)))

	return c.JSON(errorObject.StatusCode, errorJSON)
}
