package httputil

import (
	"encoding/json"
	"fmt"
	"github.com/gorpher/gone/httputil/binding"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/mileusna/useragent"
	"github.com/mitchellh/mapstructure"
)

func GetUserAgent(r *http.Request) useragent.UserAgent {
	uaStr := r.Header.Get("User-Agent")
	if uaStr != "" {
		ua := useragent.Parse(uaStr)
		return ua
	}
	return useragent.UserAgent{
		OS: useragent.Windows,
	}
}

// GetOrigin 获取客户端origin
func GetOrigin(r *http.Request) string {
	scheme := "http"
	host := r.Host
	forwardedHost := r.Header.Get("X-Forwarded-Host")
	if forwardedHost != "" {
		host = forwardedHost
	}
	forwardedProto := r.Header.Get("X-Forwarded-Proto")
	if forwardedProto == "https" {
		scheme = forwardedProto
	}

	return fmt.Sprintf("%s://%s", scheme, host)
}

func GetClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		ip = strings.Split(ip, ", ")[0]
		return ip
	}
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		ip = strings.Split(ip, ", ")[0]
		return ip
	}

	ip = r.Header.Get("Origin")
	if ip != "" && !strings.HasPrefix(ip, "http") {
		ip = net.ParseIP(ip).String()
		if ip != "" {
			return ip
		}
	}

	ip = r.RemoteAddr
	ip = strings.Split(ip, ",")[0]
	if strings.Contains(ip, ":") {
		ip, _, _ = net.SplitHostPort(ip) //nolint
	}
	return ip
}

func ShouldJson(r *http.Request, obj any) error {
	contentType := r.Header.Get("Content-Type")
	splitN := strings.SplitN(contentType, ";", 2)
	if len(splitN) > 0 {
		contentType = splitN[0]
	}
	bindType := binding.GetBinding(r.Method, contentType)
	switch bindType { // nolint
	case binding.Form:
		if err := r.ParseForm(); err != nil {
			return err
		}
		if len(r.Form) == 0 {
			break
		}
		m := make(map[string]any, len(r.Form))
		for k := range r.Form {
			m[k] = r.Form.Get(k)
		}
		if err := MapStructDecode(m, obj); err != nil {
			return err
		}
	case binding.JSON:
		if err := json.NewDecoder(r.Body).Decode(obj); err != nil {
			return err
		}
	}
	return ValidateStruct(obj)
}

// MapStructDecode takes an input structure and uses reflection to translate it to
// the output structure. output must be a pointer to a map or struct.
func MapStructDecode(input interface{}, output interface{}) error {
	config := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           output,
		TagName:          "json",
		WeaklyTypedInput: true,
		Squash:           true,
		DecodeHook:       mapstructure.ComposeDecodeHookFunc(mapstructure.StringToTimeHookFunc(time.RFC3339)),
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}
