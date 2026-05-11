package spvalidator

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"mime"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func Dir(path string) error {
	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		return nil
	}
	return fail("dir", path, nil, "path must be an existing directory")
}

func DirPath(path string) error {
	if cleanPath(path) {
		return nil
	}
	return fail("dirpath", path, nil, "value must be a directory path")
}

func File(path string) error {
	info, err := os.Stat(path)
	if err == nil && !info.IsDir() {
		return nil
	}
	return fail("file", path, nil, "path must be an existing file")
}

func FilePath(path string) error {
	if cleanPath(path) {
		return nil
	}
	return fail("filepath", path, nil, "value must be a file path")
}

func Image(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fail("image", path, nil, "path must be a readable image")
	}
	defer f.Close()
	if _, _, err := image.DecodeConfig(f); err == nil {
		return nil
	}
	return fail("image", path, nil, "path must contain a supported image")
}

func MIMEType(value string) error {
	if mediaType, _, err := mime.ParseMediaType(value); err == nil && strings.Contains(mediaType, "/") {
		return nil
	}
	return fail("mimetype", value, nil, "value must be a MIME media type")
}

func IsDefault(value any) error {
	if isZero(value) {
		return nil
	}
	return fail("isdefault", value, nil, "value must be the default zero value")
}

func Len(value any, length int) error {
	if got, ok := lengthOf(value); ok {
		if got == length {
			return nil
		}
		return fail("len", value, length, "value length must equal required length")
	}
	if cmp, ok := compareOrder(value, length); ok && cmp == 0 {
		return nil
	}
	return fail("len", value, length, "value must equal required length")
}

func Max(value any, max any) error {
	if got, ok := lengthOf(value); ok {
		limit, ok := numericParam(max)
		if ok && float64(got) <= limit {
			return nil
		}
		return fail("max", value, max, "value length must be less than or equal to maximum")
	}
	if cmp, ok := compareOrder(value, max); ok && cmp <= 0 {
		return nil
	}
	return fail("max", value, max, "value must be less than or equal to maximum")
}

func Min(value any, min any) error {
	if got, ok := lengthOf(value); ok {
		limit, ok := numericParam(min)
		if ok && float64(got) >= limit {
			return nil
		}
		return fail("min", value, min, "value length must be greater than or equal to minimum")
	}
	if cmp, ok := compareOrder(value, min); ok && cmp >= 0 {
		return nil
	}
	return fail("min", value, min, "value must be greater than or equal to minimum")
}

func OneOf(value any, choices ...any) error {
	for _, choice := range choices {
		if equalAny(value, choice) {
			return nil
		}
	}
	return fail("oneof", value, choices, "value must be one of the allowed choices")
}

func NoneOf(value any, choices ...any) error {
	for _, choice := range choices {
		if equalAny(value, choice) {
			return fail("noneof", value, choices, "value must not be one of the disallowed choices")
		}
	}
	return nil
}

func Unique(value any) error {
	v := reflect.ValueOf(value)
	if !v.IsValid() {
		return fail("unique", value, nil, "value must be a collection or string")
	}
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return fail("unique", value, nil, "value must be a collection or string")
		}
		v = v.Elem()
	}
	seen := map[string]struct{}{}
	switch v.Kind() {
	case reflect.String:
		for _, r := range v.String() {
			key := string(r)
			if _, exists := seen[key]; exists {
				return fail("unique", value, nil, "value must contain unique elements")
			}
			seen[key] = struct{}{}
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			key := fmt.Sprintf("%#v", v.Index(i).Interface())
			if _, exists := seen[key]; exists {
				return fail("unique", value, nil, "value must contain unique elements")
			}
			seen[key] = struct{}{}
		}
	case reflect.Map:
		return nil
	default:
		return fail("unique", value, nil, "value must be a collection or string")
	}
	return nil
}

func ValidateFn(value any, methodName ...string) (err error) {
	name := "Validate"
	if len(methodName) > 0 && methodName[0] != "" {
		name = methodName[0]
	}
	v := reflect.ValueOf(value)
	if !v.IsValid() {
		return fail("validateFn", value, name, "value is nil")
	}
	method := v.MethodByName(name)
	if !method.IsValid() && v.Kind() != reflect.Pointer && v.CanAddr() {
		method = v.Addr().MethodByName(name)
	}
	if !method.IsValid() {
		return fail("validateFn", value, name, "validation method does not exist")
	}
	if method.Type().NumIn() != 0 || method.Type().NumOut() > 1 {
		return fail("validateFn", value, name, "validation method must accept no arguments and return zero or one value")
	}
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fail("validateFn", value, name, fmt.Sprintf("validation method panicked: %v", recovered))
		}
	}()
	results := method.Call(nil)
	if len(results) == 0 {
		return nil
	}
	result := results[0]
	if result.Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		if result.IsNil() {
			return nil
		}
		return fail("validateFn", value, name, result.Interface().(error).Error())
	}
	if result.Kind() == reflect.Bool {
		if result.Bool() {
			return nil
		}
		return fail("validateFn", value, name, "validation method returned false")
	}
	return fail("validateFn", value, name, "validation method must return error or bool")
}

func cleanPath(path string) bool {
	return strings.TrimSpace(path) != "" && !strings.ContainsRune(path, 0) && filepath.Clean(path) != "."
}
