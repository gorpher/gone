package cookie

import (
	"net/http"
	"net/url"
	"time"
)

var optional = &Optionals{
	Path:     "/",
	HTTPOnly: true,
	SameSite: http.SameSiteLaxMode,
}

type Optionals struct {
	Path       string    // optional
	Domain     string    // optional
	Expires    time.Time // optional
	RawExpires string    // for reading cookies only

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   int
	Secure   bool
	HTTPOnly bool
	SameSite http.SameSite
}

func Set(w http.ResponseWriter, name, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:       name,
		Value:      url.QueryEscape(value),
		Path:       optional.Path,
		Domain:     optional.Domain,
		Expires:    optional.Expires,
		RawExpires: optional.RawExpires,
		MaxAge:     optional.MaxAge,
		Secure:     optional.Secure,
		HttpOnly:   optional.HTTPOnly,
		SameSite:   optional.SameSite,
	})
}

func SetWithExpires(w http.ResponseWriter, name, value string, expires int64) {
	http.SetCookie(w, &http.Cookie{
		Name:       name,
		Value:      url.QueryEscape(value),
		Path:       optional.Path,
		Domain:     optional.Domain,
		Expires:    time.Unix(expires, 0),
		RawExpires: optional.RawExpires,
		MaxAge:     optional.MaxAge,
		Secure:     optional.Secure,
		HttpOnly:   optional.HTTPOnly,
		SameSite:   optional.SameSite,
	})
}

func Get(r *http.Request, key string) *http.Cookie {
	cookies := r.Cookies()
	for i := range cookies {
		if cookies[i].Name == key {
			return cookies[i]
		}
	}
	return nil
}

func GetValue(r *http.Request, key string) string {
	cookies := r.Cookies()
	for i := range cookies {
		if cookies[i].Name == key {
			if cookies[i].Value != "" {
				v, _ := url.QueryUnescape(cookies[i].Value)
				return v
			}
		}
	}
	return ""
}

func Delete(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:       name,
		Value:      "",
		Path:       optional.Path,
		Domain:     optional.Domain,
		Expires:    optional.Expires,
		RawExpires: optional.RawExpires,
		MaxAge:     0,
		Secure:     optional.Secure,
		HttpOnly:   optional.HTTPOnly,
		SameSite:   optional.SameSite,
	})
}
