package ginutil

import (
	"github.com/gin-gonic/gin"
	"github.com/gorpher/gone/core"
	"github.com/rs/zerolog/log"
	"net/http"
)

// OkList 返回成功列表
func OkList(c *gin.Context, list interface{}, total int64) {
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

// Ok 返回成功信息, params作为动态参数，默认没有参数则返回204
func Ok(c *gin.Context, params ...interface{}) {
	if len(params) == 0 || params[0] == nil {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	data := params[0]
	str, ok := data.(string)
	if ok {
		c.Status(http.StatusOK)
		c.Abort()
		_, err := c.Writer.WriteString(str)
		if err != nil {
			return
		}
		return
	}
	var bys []byte
	bys, ok = data.([]byte)
	if ok {
		c.Status(http.StatusOK)
		c.Abort()
		_, err := c.Writer.Write(bys)
		if err != nil {
			return
		}
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, data)
}

// Bad 错误的请求
func Bad(c *gin.Context, params ...interface{}) {
	if len(params) == 0 || params[0] == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "invalid request params", "code": http.StatusBadRequest})
		return
	}
	data := params[0]
	BadError(c, http.StatusBadRequest, data, params[1:]...)
}

func BadRequest(ctx *gin.Context, v any) error {
	Bad(ctx, v)
	return nil
}
func StatusOK(c *gin.Context, data interface{}) error {
	Ok(c, data)
	return nil
}

func Status(ctx *gin.Context, status int) error {
	ctx.Status(status)
	return nil
}

// Cookie 设置Cooke
func Cookie(c *gin.Context, name, value string, maxAge int) {
	c.SetCookie(name, value, maxAge, "/", "", false, false)
}

// BadError 返回错误信息
func BadError(c *gin.Context, status int, data interface{}, params ...interface{}) {
	xUserAgent := c.GetHeader("x-user-agent")
	log.Debug().Str("url", c.Request.URL.String()).
		Str("method", c.Request.Method).
		Int("status", status).
		Str("user-agent", xUserAgent).
		Interface("params", params).
		Str("ip", GetClientIP(c)).Msg("Bad Request")
	if data == nil {
		c.AbortWithStatus(status)
		return
	}
	lang := c.GetHeader("Accept-Language")
	if lang == "" {
		lang = "en"
	}
	switch v := data.(type) {
	case string:
		c.AbortWithStatusJSON(status, gin.H{"msg": v, "code": status})
	case core.LocalMessageInterface:
		msg := v.Local(lang)
		c.Header("Accept-Language", "")
		c.AbortWithStatusJSON(status, gin.H{"msg": msg, "error": msg, "code": status})
		return
	case error:
		c.AbortWithStatusJSON(status, gin.H{"msg": v.Error(), "code": status})
		c.Error(v) // nolint
	default:
		c.AbortWithStatus(status)
	}
}
