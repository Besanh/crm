package response

import (
	"fmt"
	"net/http"
)

func Pagination(data, total, limit, offset any) (int, any) {
	return http.StatusOK, map[string]any{
		"data":   data,
		"limit":  limit,
		"offset": offset,
		"total":  total,
	}
}
func Data(code int, data any) (int, any) {
	return code, map[string]any{
		"data": data,
	}
}

func NewResponse(code int, data any) (int, any) {
	return code, data
}

func NewOKResponse(data any) (int, any) {
	return http.StatusOK, map[string]any{
		"data":    data,
		"code":    http.StatusOK,
		"content": "successfully",
	}
}

func OK(data any) (int, any) {
	return http.StatusOK, data
}

func Created(data map[string]any) (int, any) {
	result := map[string]any{
		"code":    http.StatusCreated,
		"content": "successfully",
	}
	for key, value := range data {
		result[key] = value
	}
	return http.StatusCreated, result
}

func NewErrorResponse(code int, msg any) (int, any) {
	return code, map[string]any{
		"error":   http.StatusText(code),
		"code":    code,
		"content": msg,
	}
}

func ServiceUnavailable() (int, any) {
	return http.StatusServiceUnavailable, map[string]any{
		"error":   http.StatusText(http.StatusServiceUnavailable),
		"code":    http.StatusServiceUnavailable,
		"content": http.StatusText(http.StatusServiceUnavailable),
	}
}

func ServiceUnavailableMsg(msg any) (int, any) {
	return http.StatusServiceUnavailable, map[string]any{
		"error":   http.StatusText(http.StatusServiceUnavailable),
		"code":    http.StatusServiceUnavailable,
		"content": msg,
	}
}

func BadRequest() (int, any) {
	return http.StatusBadRequest, map[string]any{
		"error":   http.StatusText(http.StatusBadRequest),
		"code":    http.StatusBadRequest,
		"content": http.StatusText(http.StatusBadRequest),
	}
}

func BadRequestMsg(msg any) (int, any) {
	return http.StatusBadRequest, map[string]any{
		"error":   http.StatusText(http.StatusBadRequest),
		"code":    http.StatusBadRequest,
		"content": msg,
	}
}

func NotFound() (int, any) {
	return http.StatusNotFound, map[string]any{
		"error":   http.StatusText(http.StatusNotFound),
		"code":    http.StatusNotFound,
		"content": http.StatusText(http.StatusNotFound),
	}
}

func NotFoundMsg(msg any) (int, any) {
	return http.StatusNotFound, map[string]any{
		"error":   http.StatusText(http.StatusNotFound),
		"code":    http.StatusNotFound,
		"content": msg,
	}
}

func Forbidden() (int, any) {
	return http.StatusForbidden, map[string]any{
		"error":   "Do not have permission for the request.",
		"code":    http.StatusForbidden,
		"content": http.StatusText(http.StatusForbidden),
	}
}

func ForbiddenLevel(level string) (int, any) {
	return http.StatusForbidden, map[string]any{
		"error":   fmt.Sprintf("level %s not have permission for the request.", level),
		"code":    http.StatusForbidden,
		"content": http.StatusText(http.StatusForbidden),
	}
}

func Unauthorized() (int, any) {
	return http.StatusUnauthorized, map[string]any{
		"error":   http.StatusText(http.StatusUnauthorized),
		"code":    http.StatusUnauthorized,
		"content": http.StatusText(http.StatusUnauthorized),
	}
}

func NewCreatedResponse(data map[string]any) (int, any) {
	result := map[string]any{
		"code":    http.StatusCreated,
		"content": "successfully",
	}
	for key, value := range data {
		result[key] = value
	}
	return http.StatusCreated, result
}

func EmptyData() []any {
	return []any{}
}

func Empty() any {
	return any(nil)
}
