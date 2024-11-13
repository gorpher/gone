package httputil

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/samber/lo"
)

// IsMobile 识别手机号码
func IsMobile(mobile string) bool {
	if mobile != "" {
		result, err := regexp.MatchString(`^1\d{10}$`, mobile)
		if err != nil {
			return false
		}
		return result
	}
	return true
}

// IsEmail 验证邮箱
func IsEmail(email string) bool {
	if email != "" {
		pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
		reg := regexp.MustCompile(pattern)
		return reg.MatchString(email)
	}
	return true
}

const (
	MaxUsernameLength  = 20
	MinUsernameLength  = 3
	MaxNicknameLength  = 20
	MinNicknameLength  = 2
	MaxGroupnameLength = 20
	MinGroupnameLength = 2
	MaxPasswordLength  = 20
	MinPasswordLength  = 6
)

// IsUserName 用户名检查
func IsUserName(username string) bool {
	//用户名长度3-20个字符、 仅能使用半角的小写字母、数字和_ - 符号和首字符只能是数字或字母,(Windows最大用户面长度20）
	if len(username) < MinUsernameLength || len(username) > MaxUsernameLength {
		return false
	}
	firstChar := username[0]
	// 首字符只能是数字或字母
	if !(unicode.IsDigit(rune(firstChar)) || unicode.IsLower(rune(firstChar))) {
		return false
	}

	for _, r := range username {
		if !(unicode.IsDigit(r) || unicode.IsLower(r) || strings.Contains("_-", string(r))) {
			return false
		}
	}

	return true
}

// IsNickname 姓名
// 长度2-20个UTF-8字符
// 可使用除了禁止输入字符之外的任意字符
func IsNickname(nickname string) bool {
	if IsIncludeForbiddenChar(nickname) {
		return false
	}
	length := lo.RuneLength(nickname)
	if length > MaxNicknameLength || length < MinNicknameLength {
		return false
	}
	return true
}

// IsGroupName 部门名称
// 长度2-20个UTF-8字符
// 可使用除了禁止输入字符之外的任意字符
func IsGroupName(nickname string) bool {
	if IsIncludeForbiddenChar(nickname) {
		return false
	}
	length := lo.RuneLength(nickname)
	if length > MaxGroupnameLength || length < MinGroupnameLength {
		return false
	}
	return true
}

// IsIncludeForbiddenChar 是否包含禁用字符,空格 ; 分号 -- 连续的减号 " " 引号 { } 大括号 [ ] 中括号 ( ) 小括号
func IsIncludeForbiddenChar(value string) bool {
	forbiddenChars := "{}[](); \""
	for _, char := range value {
		if strings.Contains(forbiddenChars, string(char)) {
			return true
		}
	}
	return strings.Contains(value, "--")
}

// IsPassword
// 1. 设置密码规则有，最小长度默认最小长度6个字符，字符不能小于6个，默认最大密码长度（20）
// 2. 设置包含：大写字母、小写字母、数字和字符（默认都不开启）
// 3. 开启限制后，用户注册必须按照规范注册用户
// 4. 已经注册过的用户，登录时不按照密码规则校验
// 5. 密码特殊字符包括： !@#$%^&.
// 6. 默认密码不开启验证时，密码必须属于写字母、小写字母、数字和字符中一种或多种组合
func IsPassword(password string, chains ...func(string) bool) bool {
	if len(password) > MaxPasswordLength || len(password) < MinPasswordLength {
		return false
	}
	for _, r := range password {
		if !(unicode.IsDigit(r) || unicode.IsUpper(r) || unicode.IsLower(r) || IsPasswordSpecialLetter(r)) {
			return false
		}

	}
	for _, fn := range chains {
		if !fn(password) {
			return false
		}
	}
	return true
}

const PasswordSpecialLetter = "!@#$%^&."

// 1大写，2小写，4数字，8符号
const (
	PasswordUpper int = 1 << iota
	PasswordLower
	PasswordDigit
	PasswordSpecial
)

func IsPasswordSpecialLetter(r rune) bool {
	for _, c := range PasswordSpecialLetter {
		if c == r {
			return true
		}
	}
	return false
}

// IsHash 是否是有效的hash
func IsHash(s string) bool {
	for _, r := range s {
		if !(unicode.IsDigit(r) || unicode.IsLower(r)) {
			return false
		}
	}
	return true
}
