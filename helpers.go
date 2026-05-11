package spvalidator

import (
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func isZero(value any) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return true
		}
		v = v.Elem()
	}
	if t, ok := v.Interface().(time.Time); ok {
		return t.IsZero()
	}
	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		return v.Len() == 0
	default:
		return v.IsZero()
	}
}

func stringValue(value any) (string, bool) {
	switch v := value.(type) {
	case string:
		return v, true
	case []byte:
		return string(v), true
	case fmt.Stringer:
		return v.String(), true
	default:
		return "", false
	}
}

func numberValue(value any) (float64, bool) {
	switch v := value.(type) {
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	case uintptr:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	case string:
		if strings.TrimSpace(v) == "" {
			return 0, false
		}
		n, err := strconv.ParseFloat(v, 64)
		return n, err == nil
	default:
		return 0, false
	}
}

func compareOrder(left any, right any) (int, bool) {
	if ln, ok := exactNumber(left); ok {
		rn, ok := exactNumber(right)
		if !ok {
			return 0, false
		}
		return ln.Cmp(rn), true
	}

	if lt, ok := left.(time.Time); ok {
		rt, ok := right.(time.Time)
		if !ok {
			return 0, false
		}
		switch {
		case lt.Before(rt):
			return -1, true
		case lt.After(rt):
			return 1, true
		default:
			return 0, true
		}
	}

	if ln, ok := numberValue(left); ok {
		rn, ok := numberValue(right)
		if !ok {
			return 0, false
		}
		switch {
		case ln < rn:
			return -1, true
		case ln > rn:
			return 1, true
		default:
			return 0, true
		}
	}

	ls, lok := left.(string)
	rs, rok := right.(string)
	if lok && rok {
		return strings.Compare(ls, rs), true
	}

	lb, lok := left.(bool)
	rb, rok := right.(bool)
	if lok && rok {
		switch {
		case lb == rb:
			return 0, true
		case !lb && rb:
			return -1, true
		default:
			return 1, true
		}
	}

	return 0, false
}

func equalAny(left any, right any) bool {
	if ln, ok := exactNumber(left); ok {
		rn, ok := exactNumber(right)
		return ok && ln.Cmp(rn) == 0
	}
	return reflect.DeepEqual(left, right)
}

func lengthOf(value any) (int, bool) {
	if value == nil {
		return 0, false
	}
	v := reflect.ValueOf(value)
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return 0, false
		}
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.String:
		return utf8.RuneCountInString(v.String()), true
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		return v.Len(), true
	default:
		return 0, false
	}
}

func numericParam(value any) (float64, bool) {
	if n, ok := numberValue(value); ok {
		return n, true
	}
	return 0, false
}

func exactNumber(value any) (*big.Rat, bool) {
	switch v := value.(type) {
	case int:
		return new(big.Rat).SetInt(big.NewInt(int64(v))), true
	case int8:
		return new(big.Rat).SetInt(big.NewInt(int64(v))), true
	case int16:
		return new(big.Rat).SetInt(big.NewInt(int64(v))), true
	case int32:
		return new(big.Rat).SetInt(big.NewInt(int64(v))), true
	case int64:
		return new(big.Rat).SetInt(big.NewInt(v)), true
	case uint:
		return new(big.Rat).SetInt(new(big.Int).SetUint64(uint64(v))), true
	case uint8:
		return new(big.Rat).SetInt(new(big.Int).SetUint64(uint64(v))), true
	case uint16:
		return new(big.Rat).SetInt(new(big.Int).SetUint64(uint64(v))), true
	case uint32:
		return new(big.Rat).SetInt(new(big.Int).SetUint64(uint64(v))), true
	case uint64:
		return new(big.Rat).SetInt(new(big.Int).SetUint64(v)), true
	case uintptr:
		return new(big.Rat).SetInt(new(big.Int).SetUint64(uint64(v))), true
	case float32:
		return exactNumber(float64(v))
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return nil, false
		}
		return new(big.Rat).SetFloat64(v), true
	case string:
		s := strings.TrimSpace(v)
		if s == "" {
			return nil, false
		}
		r, ok := new(big.Rat).SetString(s)
		return r, ok
	default:
		return nil, false
	}
}
