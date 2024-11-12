package authed

import (
	"github.com/gorpher/gone/codec"
	"github.com/gorpher/gone/core"
	"time"
)

type SubjectType string

const (
	SubjectTypeAuthToken    string = "auth_token"
	SubjectTypeRefreshToken string = "refresh_token"
)

// UserSession 用户会话
type UserSession struct {
	OsName     string         `json:"os_name"`
	OsArch     string         `json:"os_arch"`
	ClientName string         `json:"client_name"`
	Username   string         `json:"username"`
	Nickname   string         `json:"nickname"`
	ID         string         `json:"id"`
	ExpiredAt  int64          `json:"expired_at"` // 单位秒，时间戳
	Roles      []string       `json:"roles,omitempty"`
	Scopes     []string       `json:"scopes,omitempty"`
	Extends    map[string]any `json:"extends,omitempty"`
	token      string
}

func (u *UserSession) GetToken() string {
	return u.token
}

type Payload struct {
	*UserSession
	codec.Payload
}

func (p *Payload) SetExpired(t time.Time) {
	if p.UserSession != nil {
		p.UserSession.ExpiredAt = t.Unix()
		p.ExpiredAt = t.Unix()
	}
	p.ExpirationTime = core.NewTime(t)
}
