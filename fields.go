package spvalidator

import (
	"fmt"
	"reflect"
	"strings"
)

func fieldByPath(target any, path string) (reflect.Value, error) {
	if strings.TrimSpace(path) == "" {
		return reflect.Value{}, fmt.Errorf("field path is empty")
	}
	v := reflect.ValueOf(target)
	if !v.IsValid() {
		return reflect.Value{}, fmt.Errorf("target is nil")
	}
	for part := range strings.SplitSeq(path, ".") {
		if part == "" {
			return reflect.Value{}, fmt.Errorf("field path %q contains an empty segment", path)
		}
		for v.Kind() == reflect.Interface || v.Kind() == reflect.Pointer {
			if v.IsNil() {
				return reflect.Value{}, fmt.Errorf("field path %q reached nil before %q", path, part)
			}
			v = v.Elem()
		}
		switch v.Kind() {
		case reflect.Struct:
			field := v.FieldByName(part)
			if !field.IsValid() {
				return reflect.Value{}, fmt.Errorf("field %q does not exist", part)
			}
			if !field.CanInterface() {
				return reflect.Value{}, fmt.Errorf("field %q is not exported", part)
			}
			v = field
		case reflect.Map:
			if v.Type().Key().Kind() != reflect.String {
				return reflect.Value{}, fmt.Errorf("map in field path %q does not use string keys", path)
			}
			key := reflect.ValueOf(part).Convert(v.Type().Key())
			field := v.MapIndex(key)
			if !field.IsValid() {
				return reflect.Value{}, fmt.Errorf("map key %q does not exist", part)
			}
			v = field
		default:
			return reflect.Value{}, fmt.Errorf("cannot resolve %q through %s", part, v.Kind())
		}
	}
	return v, nil
}

func fieldInterface(target any, path string) (any, error) {
	v, err := fieldByPath(target, path)
	if err != nil {
		return nil, err
	}
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil, nil
		}
		v = v.Elem()
	}
	if !v.CanInterface() {
		return nil, fmt.Errorf("field %q is not exported", path)
	}
	return v.Interface(), nil
}

func EqField(target any, fieldPath string, otherPath string) error {
	left, right, err := fieldPair(target, fieldPath, otherPath)
	if err != nil {
		return fail("eqfield", target, fieldPath+"="+otherPath, err.Error())
	}
	if equalAny(left, right) {
		return nil
	}
	return failf("eqfield", left, right, "must equal %v", right)
}

func EqCSField(target any, fieldPath string, otherPath string) error {
	return EqField(target, fieldPath, otherPath)
}

func NeField(target any, fieldPath string, otherPath string) error {
	left, right, err := fieldPair(target, fieldPath, otherPath)
	if err != nil {
		return fail("nefield", target, fieldPath+"!="+otherPath, err.Error())
	}
	if !equalAny(left, right) {
		return nil
	}
	return failf("nefield", left, right, "must not equal %v", right)
}

func NeCSField(target any, fieldPath string, otherPath string) error {
	return NeField(target, fieldPath, otherPath)
}

func GtField(target any, fieldPath string, otherPath string) error {
	return orderedField("gtfield", target, fieldPath, otherPath, func(cmp int) bool { return cmp > 0 }, "must be greater than %v")
}

func GtCSField(target any, fieldPath string, otherPath string) error {
	return GtField(target, fieldPath, otherPath)
}

func GteField(target any, fieldPath string, otherPath string) error {
	return orderedField("gtefield", target, fieldPath, otherPath, func(cmp int) bool { return cmp >= 0 }, "must be greater than or equal to %v")
}

func GteCSField(target any, fieldPath string, otherPath string) error {
	return GteField(target, fieldPath, otherPath)
}

func LtField(target any, fieldPath string, otherPath string) error {
	return orderedField("ltfield", target, fieldPath, otherPath, func(cmp int) bool { return cmp < 0 }, "must be less than %v")
}

func LtCSField(target any, fieldPath string, otherPath string) error {
	return LtField(target, fieldPath, otherPath)
}

func LteField(target any, fieldPath string, otherPath string) error {
	return orderedField("ltefield", target, fieldPath, otherPath, func(cmp int) bool { return cmp <= 0 }, "must be less than or equal to %v")
}

func LteCSField(target any, fieldPath string, otherPath string) error {
	return LteField(target, fieldPath, otherPath)
}

func FieldContains(target any, fieldPath string, chars string) error {
	value, err := fieldInterface(target, fieldPath)
	if err != nil {
		return fail("fieldcontains", target, fieldPath, err.Error())
	}
	s, ok := stringValue(value)
	if !ok {
		return fail("fieldcontains", value, chars, "field must be string-like")
	}
	if strings.Contains(s, chars) {
		return nil
	}
	return failf("fieldcontains", s, chars, "must contain %q", chars)
}

func FieldExcludes(target any, fieldPath string, chars string) error {
	value, err := fieldInterface(target, fieldPath)
	if err != nil {
		return fail("fieldexcludes", target, fieldPath, err.Error())
	}
	s, ok := stringValue(value)
	if !ok {
		return fail("fieldexcludes", value, chars, "field must be string-like")
	}
	if !strings.Contains(s, chars) {
		return nil
	}
	return failf("fieldexcludes", s, chars, "must not contain %q", chars)
}

func fieldPair(target any, fieldPath string, otherPath string) (any, any, error) {
	left, err := fieldInterface(target, fieldPath)
	if err != nil {
		return nil, nil, err
	}
	right, err := fieldInterface(target, otherPath)
	if err != nil {
		return nil, nil, err
	}
	return left, right, nil
}

func orderedField(tag string, target any, fieldPath string, otherPath string, pass func(int) bool, messageFormat string) error {
	left, right, err := fieldPair(target, fieldPath, otherPath)
	if err != nil {
		return fail(tag, target, fieldPath+"~"+otherPath, err.Error())
	}
	cmp, ok := compareOrder(left, right)
	if ok && pass(cmp) {
		return nil
	}
	return failf(tag, left, right, messageFormat, right)
}
