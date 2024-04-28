package response

import (
	"context"
	"net/http"

	i18n_err "github.com/Risuii/frs-lib/i18n/errors"
	"github.com/Risuii/invoice/src/middleware/request"
)

func JSONSuccessResponse(ctx context.Context, w http.ResponseWriter, data interface{}) {
	JSONResponse(ctx, w, createSuccessResponse(data, request.GetRequestID(ctx)), http.StatusOK)
}

func JSONUnauthorizedResponse(ctx context.Context, w http.ResponseWriter) {
	JSONResponse(ctx, w, createErrorResponse(i18n_err.ErrUnauthorized, request.GetRequestID(ctx), request.GetLanguage(ctx)),
		http.StatusUnauthorized)
}

func JSONInternalErrorResponse(ctx context.Context, w http.ResponseWriter) {
	JSONResponse(ctx, w, createErrorResponse(i18n_err.ErrInternalServer, request.GetRequestID(ctx), request.GetLanguage(ctx)),
		http.StatusInternalServerError)
}

func JSONBadRequestResponse(ctx context.Context, w http.ResponseWriter) {
	JSONResponse(ctx, w, createErrorResponse(i18n_err.ErrBadRequest, request.GetRequestID(ctx), request.GetLanguage(ctx)),
		http.StatusBadRequest)
}

func JSONUnprocessableEntity(ctx context.Context, w http.ResponseWriter, err i18n_err.I18nError) {
	JSONResponse(ctx, w, createErrorResponse(err, request.GetRequestID(ctx), request.GetLanguage(ctx)),
		http.StatusUnprocessableEntity)
}
