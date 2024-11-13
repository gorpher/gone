package httputil

import (
	"context"
	"net/http"
	"strconv"
)

func UidSet(r *http.Request, uid int64) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), "uid", uid)) // nolint
}

func SetContext(r *http.Request, key string, value any) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), key, value)) // nolint
}

func UidGet(r *http.Request) int64 {
	value := r.Context().Value("uid")
	if value == nil {
		return 0
	}
	switch s := value.(type) {
	case int64:
		return s
	case int32:
		return int64(s)
	case int8:
		return int64(s)
	case int16:
		return int64(s)
	case string:
		v, _ := strconv.ParseInt(s, 10, 64) //nolint
		return v
	default:
		return 0
	}
}
