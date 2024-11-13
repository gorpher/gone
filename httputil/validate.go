package httputil

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
)

var (
	onceValidate sync.Once
	validate     *validator.Validate
)

type SliceValidationError []error

// Error concatenates all error elements in SliceValidationError into a single string separated by \n.
func (err SliceValidationError) Error() string {
	n := len(err)
	switch n {
	case 0:
		return ""
	default:
		var b strings.Builder
		if err[0] != nil {
			fmt.Fprintf(&b, "[%d]: %s", 0, err[0].Error())
		}
		if n > 1 {
			for i := 1; i < n; i++ {
				if err[i] != nil {
					b.WriteString("\n")
					fmt.Fprintf(&b, "[%d]: %s", i, err[i].Error())
				}
			}
		}
		return b.String()
	}
}

// ValidateStruct receives any kind of type, but only performed struct or pointer to struct type.
func ValidateStruct(obj any) error {
	if obj == nil {
		return nil
	}

	value := reflect.ValueOf(obj)
	switch value.Kind() { // nolint
	case reflect.Ptr:
		return ValidateStruct(value.Elem().Interface())
	case reflect.Struct:
		return validateStruct(obj)
	case reflect.Slice, reflect.Array:
		count := value.Len()
		validateRet := make(SliceValidationError, 0)
		for i := 0; i < count; i++ {
			if err := ValidateStruct(value.Index(i).Interface()); err != nil {
				validateRet = append(validateRet, err)
			}
		}
		if len(validateRet) == 0 {
			return nil
		}
		return validateRet
	default:
		return nil
	}
}

// validateStruct receives struct type
func validateStruct(obj any) error {
	lazyinit()
	return validate.Struct(obj)
}

func lazyinit() {
	onceValidate.Do(func() {
		validate = validator.New(validator.WithRequiredStructEnabled())
		validate.RegisterValidation("isMobile", wrapperValidate(IsMobile))       // nolint
		validate.RegisterValidation("isUserName", wrapperValidate(IsUserName))   // nolint
		validate.RegisterValidation("isNickName", wrapperValidate(IsNickname))   // nolint
		validate.RegisterValidation("isGroupName", wrapperValidate(IsGroupName)) // nolint
		validate.RegisterValidation("isLessThanNow", IsLessThanNow)              // nolint
		validate.RegisterValidation("isMoreThanNow", IsMoreThanNow)              // nolint
		validate.RegisterValidation("isLessEndTime", IsLessEndTime)              // nolint
	})
}

func wrapperValidate(isOk func(string) bool) func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		if value, ok := fl.Field().Interface().(string); ok {
			return isOk(value)
		}
		return false
	}
}

func IsLessEndTime(fl validator.FieldLevel) bool {
	endTimeName := fl.Param()
	if endTime, ok := fl.Parent().FieldByName(endTimeName).Interface().(int64); ok {
		if endTime == 0 {
			return true
		}
		if startTime, ok := fl.Field().Interface().(int64); ok {
			return startTime <= endTime
		}
	}
	return true
}

// IsLessThanNow 检查时间小于当前时间戳
func IsLessThanNow(fl validator.FieldLevel) bool {
	now := time.Now().UTC().Unix()
	if value, ok := fl.Field().Interface().(int64); ok {
		return value <= now
	}
	return false
}

// IsMoreThanNow 检查时间大于当前时间戳
func IsMoreThanNow(fl validator.FieldLevel) bool {
	now := time.Now().UTC().Unix()
	if value, ok := fl.Field().Interface().(int64); ok {
		return value > now
	}
	return false
}
