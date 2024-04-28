package response

import (
	"context"
	"encoding/json"
	"net/http"

	frsI18n "github.com/Risuii/frs-lib/i18n"
	frsI18nErr "github.com/Risuii/frs-lib/i18n/errors"
	"github.com/Risuii/invoice/src/middleware/request"
)

type Response struct {
	Data     interface{} `json:"data"`
	Error    *Error      `json:"error"`
	Success  bool        `json:"success"`
	Metadata Meta        `json:"metadata"`
}

type Meta struct {
	RequestId string `json:"request_id"`
}

type Error struct {
	Code     string `json:"code"`
	Title    string `json:"message_title"`
	Message  string `json:"message"`
	Severity string `json:"message_severity"`
}

func JSONResponse(ctx context.Context, w http.ResponseWriter, data Response, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func createSuccessResponse(data interface{}, reqId string) Response {
	return Response{
		Data:    data,
		Success: true,
		Metadata: Meta{
			RequestId: reqId,
		},
	}
}

func createErrorResponse(err frsI18nErr.I18nError, reqId, lang string, args ...interface{}) Response {

	return Response{
		Error: &Error{
			Code:     err.Error(),
			Title:    frsI18n.Title(lang, err.Error()),
			Message:  frsI18n.Message(lang, err.Error()),
			Severity: "error",
		},
		Metadata: Meta{
			RequestId: reqId,
		},
	}
}

func JSONSuccess(ctx context.Context, w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := Response{
		Data:    data,
		Success: true,
		Metadata: Meta{
			RequestId: request.GetRequestID(ctx),
		},
	}

	json.NewEncoder(w).Encode(resp)
}

func JSONError(ctx context.Context, w http.ResponseWriter, code int, err frsI18nErr.I18nError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	lang := request.GetLanguage(ctx)
	resp := Response{
		Error: &Error{
			Code:     err.Error(),
			Title:    frsI18n.Title(lang, err.Error()),
			Message:  frsI18n.Message(lang, err.Error()),
			Severity: "error",
		},
		Metadata: Meta{
			RequestId: request.GetRequestID(ctx),
		},
	}

	json.NewEncoder(w).Encode(resp)
}
