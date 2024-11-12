package cookie

import (
	"github.com/gorpher/gone/codec"
	"net/http"
	"net/url"
	"time"
)

var (
	//16, 24, or 32 bytes to select
	hashKey = []byte{
		0xf3, 0x90, 0x19, 0x8e, 0xb8, 0x12, 0x1c, 0x56,
		0xf4, 0xde, 0x16, 0x2b, 0x8f, 0xaa, 0xf3, 0x98,
	}
	blockKey = []byte{
		0x9c, 0x93, 0x5b, 0x28, 0x13, 0x0a, 0x55, 0x49,
		0x5b, 0xfd, 0x3c, 0x63, 0x98, 0x86, 0xa9, 0x47,
	}
	cryptoKey = []byte{
		0xf3, 0x90, 0x19, 0x8e, 0xb8, 0x12, 0x1c, 0x56,
		0xf4, 0xde, 0x16, 0x2b, 0x8f, 0xaa, 0xf3, 0x98,
	}
)
var cryptoCodec = codec.NewCookieCodec(hashKey, blockKey)

func SetCodec(ck, hashKey, blockKey []byte) {
	cryptoCodec = codec.NewCookieCodec(hashKey, blockKey)
	cryptoKey = ck
}

func SetCryptoCookie(w http.ResponseWriter, key, value string, maxAge int) (err error) {
	var cryptoBytes []byte
	cryptoBytes, err = cryptoCodec.Encode(cryptoKey, []byte(value))
	if err != nil {
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     key,
		Value:    url.QueryEscape(string(cryptoBytes)),
		Path:     "/",
		Domain:   "",
		Expires:  time.Now().Add(time.Duration(maxAge)).UTC(),
		MaxAge:   maxAge,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	return err
}

func GetCryptoCookie(req *http.Request, k string) string {
	var value []byte
	value = []byte(req.Header.Get(k))
	if len(value) == 0 {
		c, err := req.Cookie(k)
		if err != nil {
			return ""
		}
		value = []byte(c.Value)
	}
	if len(value) == 0 {
		return ""
	}
	v, err := cryptoCodec.Decode(cryptoKey, value)
	if err != nil {
		return ""
	}
	return string(v)
}
