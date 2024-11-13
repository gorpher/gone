package authed

import (
	"errors"
	"fmt"
	"github.com/gorpher/gone/cache"
	"github.com/gorpher/gone/codec"
	"github.com/gorpher/gone/cookie"
	"github.com/gorpher/gone/core"
	"github.com/gorpher/gone/osutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Authed struct {
	Issuer               string
	Audience             []string // example: appname
	TokenDuration        time.Duration
	RefreshTokenDuration time.Duration
	MultiSession         bool
	// ===============================
	cookieName  string // example: appname
	cryptoKey   []byte
	cryptoCodec codec.CryptoCodec
	objectCodec codec.ObjectCodec
	store       cache.Cache
}

type OptFunc func(session *Authed) *Authed

func WithCookieCode(hashKey, key []byte) OptFunc {
	return func(s *Authed) *Authed {
		s.cryptoCodec = codec.NewCookieCodec(hashKey, key)
		return s
	}
}
func WithJwtCode(alg string) OptFunc {
	return func(s *Authed) *Authed {
		s.cryptoCodec = codec.NewJwtCodec(alg)
		return s
	}
}

// WithCryptoKey 初始化加密jwt的hs512的密钥key
func WithCryptoKey(key []byte) OptFunc {
	return func(s *Authed) *Authed {
		s.cryptoKey = key
		return s
	}
}

func WithCookieName(cookieName string) OptFunc {
	return func(s *Authed) *Authed {
		s.cookieName = cookieName
		return s
	}
}
func WithCache(c cache.Cache) OptFunc {
	return func(s *Authed) *Authed {
		s.store = c
		return s
	}
}
func WithMultiSession() OptFunc {
	return func(s *Authed) *Authed {
		s.MultiSession = true
		return s
	}
}

var cryptoKey = []byte{
	0x9c, 0x93, 0x5b, 0x28, 0x13, 0x0a, 0x55, 0x49,
	0x5b, 0xfd, 0x3c, 0x63, 0x98, 0x86, 0xa9, 0x47,
	0xf3, 0x90, 0x19, 0x8e, 0xb8, 0x12, 0x1c, 0x56,
	0xf4, 0xde, 0x16, 0x2b, 0x8f, 0xaa, 0xf3, 0x98,
}
var ErrorInvalidSession = errors.New("invalid session")
var ErrorInvalidPayload = errors.New("invalid payload")
var ErrorInvalidRefreshToken = errors.New("invalid refresh token")
var ErrorInvalidToken = errors.New("invalid token")

func NewAuthed(opts ...OptFunc) *Authed {
	s := &Authed{
		Issuer:               "gone.authed",
		Audience:             []string{"app"},
		TokenDuration:        time.Hour * 2,
		RefreshTokenDuration: time.Hour * 24 * 7,
		cryptoKey:            cryptoKey,
		cookieName:           "authed",
		cryptoCodec:          codec.NewJwtCodec("HS256"),
		objectCodec:          codec.JSONEncoder{},
		store:                cache.NewMemoryCache(),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Authed) FormatTokenStoreKey(key string) string {
	return fmt.Sprintf("%s/authed/token/%s", s.cookieName, key)
}
func (s *Authed) FormatRefreshTokenStoreKey(key string) string {
	return fmt.Sprintf("%s/authed/refreshtoken/%s", s.cookieName, key)
}

func (s *Authed) FormatLinkTokenStoreKey(key string) string {
	return fmt.Sprintf("%s/authed/linktk/%s", s.cookieName, key)
}

func (s *Authed) NewClaims(subject string, se *UserSession, duration time.Duration) *Payload {
	timeNow := core.Now()
	var expiredAt *core.Time
	if se.ExpiredAt != 0 {
		expiredAt = core.NewTime(time.Unix(se.ExpiredAt, 0))
	} else {
		expiredAt = core.NewTime(timeNow.Add(duration))
		se.ExpiredAt = expiredAt.Unix()
	}
	if se.ID == "" {
		se.ID = osutil.NumberID()
	}
	return &Payload{
		Payload: codec.Payload{
			Issuer:         s.Issuer,
			Audience:       s.Audience,
			ExpirationTime: expiredAt,
			IssuedAt:       &timeNow,
			NotBefore:      &timeNow,
			Subject:        subject,
			JWTID:          se.ID,
		},
		UserSession: se,
	}
}
func (s *Authed) GetHTTPSession(req *http.Request) (se *UserSession) {
	var (
		payload    = Payload{}
		err        error
		headerAuth = req.Header.Get("Authorization")
	)
	const basicPrefix = "Basic "
	const BearerPrefix = "Bearer "
	cookieAuth := cookie.GetValue(req, s.cookieName)
	if (cookieAuth == "" && headerAuth == "") || strings.HasPrefix(headerAuth, basicPrefix) {
		return nil
	}
	authToken := strings.TrimPrefix(headerAuth, BearerPrefix)
	if authToken == "" {
		authToken = cookieAuth
	}
	payload, err = s.verifyToken(authToken)
	if err != nil {
		return
	}
	if payload.UserSession == nil {
		return
	}
	se = payload.UserSession
	return
}
func (s *Authed) CreateToken(se *UserSession) (token, refresh string, err error) {
	if se == nil {
		err = ErrorInvalidSession
		return
	}
	token, refresh, err = s.createToken(s.NewClaims(SubjectTypeAuthToken, se, s.TokenDuration))
	return
}
func (s *Authed) createToken(payload *Payload) (token, refresh string, err error) {
	if payload == nil {
		err = ErrorInvalidPayload
		return
	}
	var plainByte []byte
	plainByte, err = s.objectCodec.Encode(payload)
	if err != nil {
		return
	}
	var ecryptoBase64 []byte
	ecryptoBase64, err = s.cryptoCodec.Encode(s.cryptoKey, plainByte)
	if err != nil {
		return
	}
	token = string(ecryptoBase64)
	refresh = osutil.UUID()

	err = s.store.SetWithTTL(s.FormatTokenStoreKey(payload.JWTID), token, time.Until(payload.ExpirationTime.Time))
	if err != nil {
		return
	}
	err = s.store.SetWithTTL(s.FormatRefreshTokenStoreKey(refresh), token, s.RefreshTokenDuration)
	if err != nil {
		return
	}
	err = s.store.SetWithTTL(s.FormatLinkTokenStoreKey(payload.JWTID), refresh, s.RefreshTokenDuration)
	return
}

func (s *Authed) DeleteToken(id string) (err error) {
	var freshTokenByte []byte
	freshTokenByte, err = s.store.Get(s.FormatLinkTokenStoreKey(id))
	if err != nil {
		return
	}
	err = s.store.Del(string(freshTokenByte))
	if err != nil {
		return
	}
	err = s.store.Del(s.FormatTokenStoreKey(id))
	if err != nil {
		return
	}
	return
}

func (s *Authed) DeleteTokenOnly(id string) (err error) {
	err = s.store.Del(s.FormatTokenStoreKey(id))
	return
}

func (s *Authed) RefreshToken(refreshToken string) (token, refresh string, err error) {
	if refreshToken == "" {
		err = ErrorInvalidRefreshToken
		return
	}
	var tokenBytes []byte
	tokenBytes, err = s.store.Get(s.FormatRefreshTokenStoreKey(refreshToken))
	if err != nil {
		return
	}
	var plainByte []byte
	plainByte, err = s.cryptoCodec.Decode(s.cryptoKey, tokenBytes)
	if err != nil {
		return
	}
	var payload Payload
	err = s.objectCodec.Decode(plainByte, &payload)
	if err != nil {
		return
	}
	payload.SetExpired(core.Now().Add(s.TokenDuration))
	token, refresh, err = s.createToken(&payload)
	return
}
func (s *Authed) VerifyToken(token string) (payload Payload, err error) {
	return s.verifyToken(token)
}
func (s *Authed) verifyToken(token string) (payload Payload, err error) {
	if token == "" {
		err = ErrorInvalidToken
		return
	}
	var plainByte []byte
	plainByte, err = s.cryptoCodec.Decode(s.cryptoKey, []byte(token))
	if err != nil {
		return
	}
	err = s.objectCodec.Decode(plainByte, &payload)
	if err != nil {
		return
	}
	if payload.UserSession != nil {
		payload.UserSession.token = token
	}
	if s.MultiSession {
		return
	}
	var tokenSavedByte []byte
	tokenSavedByte, err = s.store.Get(s.FormatTokenStoreKey(payload.JWTID))
	if err != nil {
		return
	}
	tokenSaved := string(tokenSavedByte)
	if tokenSaved != token {
		err = fmt.Errorf("invalid token")
		return
	}
	return
}
func (s *Authed) SetCookieToken(w http.ResponseWriter, value string, maxAge int) {
	s.SetCookie(w, s.cookieName, value, maxAge)
}
func (s *Authed) SetCookie(w http.ResponseWriter, key, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     key,
		Value:    url.QueryEscape(value),
		Path:     "/",
		Domain:   "",
		Expires:  time.Now().Add(time.Duration(maxAge)).UTC(),
		MaxAge:   maxAge,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
