package ginutil

import (
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"strconv"
	"strings"
)

func GetInt64Param(c *gin.Context, key string) int64 {
	i, err := strconv.ParseInt(c.Param(key), 10, 64)
	if err != nil {
		return 0
	}
	return i
}

func GetInt64Query(c *gin.Context, key string) int64 {
	i, err := strconv.ParseInt(c.Query(key), 10, 64)
	if err != nil {
		return 0
	}
	return i
}

func GetInt64(r *http.Request, key string) int64 {
	valueStr := r.FormValue(key)
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return 0
	}
	return value
}

func GetInt(r *http.Request, key string) int {
	valueStr := r.FormValue(key)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0
	}
	return value
}

func ParamInt64(c *gin.Context, key string) int64 {
	valueStr := c.Param(key)
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return 0
	}
	return value
}

func GetPathInt64(c *gin.Context, key string, defaultValue int64) int64 {
	str := c.Param(key)
	i, err := strconv.ParseInt(str, 10, 64)
	if err == nil {
		return i
	}
	return defaultValue
}

func ShouldJson(c *gin.Context, v any) error {
	return c.ShouldBindJSON(v)
}
func GetClientIP(c *gin.Context) string {
	ip := c.GetHeader("X-Forwarded-For")
	if ip != "" {
		ip = strings.Split(ip, ", ")[0]
		return ip
	}
	ip = c.GetHeader("X-Real-IP")
	if ip != "" {
		ip = strings.Split(ip, ", ")[0]
		return ip
	}
	ip = c.RemoteIP()
	ip = strings.Split(ip, ", ")[0]

	if strings.Contains(ip, ":") {
		ip, _, _ = net.SplitHostPort(ip) //nolint
	}
	return ip
}

func UidGet(c *gin.Context) int64 {
	value, exist := c.Get("uid")
	if !exist {
		return 0
	}
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
