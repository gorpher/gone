package httputil

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorpher/gone/core"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type JsonRawBody map[string]any

// JsonListBody 分页响应体
type JsonListBody struct {
	List  any    `json:"list"`
	Total int64  `json:"total"`
	Code  int    `json:"code,omitempty"`
	Msg   string `json:"msg,omitempty"`
}

// JsonDataBody json响应体
type JsonDataBody struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

// Ok 返回成功信息, params作为动态参数，默认没有参数则返回204
func Ok(w http.ResponseWriter, params ...any) {
	if len(params) == 0 || params[0] == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	data := params[0]
	str, ok := data.(string)
	if ok {
		w.Header().Set("Content-Length", strconv.FormatInt(int64(len(str)), 10))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(str)) // nolint
		return
	}
	by, ok := data.([]byte)
	if ok {
		w.Header().Set("Content-Length", strconv.FormatInt(int64(len(by)), 10))
		w.WriteHeader(http.StatusOK)
		w.Write(by) // nolint
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data) // nolint
}

// OkList 返回成功列表
func OkList(w http.ResponseWriter, list any, total int64) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(JsonRawBody{ // nolint
		"list":  list,
		"total": total,
	})
}

// Bad 错误的请求
func Bad(w http.ResponseWriter, params ...any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if len(params) == 0 || params[0] == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(JsonRawBody{ // nolint
			"error": "invalid request params",
			"msg":   "invalid request params",
			"code":  http.StatusBadRequest,
		})
		return
	}
	data := params[0]
	BadError(w, http.StatusBadRequest, data, params[1:]...)
}

// BadError 返回错误信息
func BadError(w http.ResponseWriter, status int, data any, params ...any) {
	if data == nil {
		w.WriteHeader(status)
		return
	}
	lang := w.Header().Get("Accept-Language")
	if lang == "" {
		lang = "en"
	}
	switch v := data.(type) {
	case string:
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(JsonRawBody{"msg": v, "error": v, "code": status}) // nolint
	case validator.ValidationErrors:
		w.WriteHeader(status)
		var msg string
		for _, fieldError := range v {
			msg = validatorMsg(fieldError)
			break
		}
		json.NewEncoder(w).Encode(JsonRawBody{"msg": msg, "error": msg, "code": status}) // nolint
		return
	case core.LocalMessageInterface:
		msg := v.Local(lang)
		w.Header().Del("Accept-Language")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(JsonRawBody{"msg": msg, "error": msg, "code": status}) // nolint
		return
	case error:
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(JsonRawBody{"msg": v.Error(), "error": v.Error(), "code": status}) // nolint
	default:
		w.WriteHeader(status)
	}
}

// Forbidden Forbidden
func Forbidden(w http.ResponseWriter, err any) {
	BadError(w, http.StatusForbidden, 0, err)
}

// BadW 返回错误信息
func BadW(w http.ResponseWriter, msg string) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(fmt.Sprintf("{\"msg\":\"%s\"}", msg))) //nolint
	return errors.New(msg)
}

var fieldMap = map[string]string{}
var tagMap = map[string]string{}

func SetTagMap(t map[string]string) {
	tagMap = t
}

func SetFieldMap(t map[string]string) {
	fieldMap = t
}

// validatorMsg
func validatorMsg(e validator.FieldError) string {
	field := e.Field()
	fieldVal, findField := fieldMap[field]
	tagVal, findTag := tagMap[e.Tag()]
	if findField && findTag {
		return fieldVal + tagVal
	}
	if !findField && findTag {
		return "param:" + field + tagVal
	}
	return "invalid request params"
}
