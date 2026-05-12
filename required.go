package spvalidator

// FieldCondition describes a struct field value condition used by required_* and excluded_* validators.
type FieldCondition struct {
	Field string
	Value any
}

// Condition creates a FieldCondition.
func Condition(field string, value any) FieldCondition {
	return FieldCondition{Field: field, Value: value}
}

func Required(value any) error {
	if !isZero(value) {
		return nil
	}
	return fail("required", value, nil, "required")
}

func RequiredIf(target any, fieldPath string, conditions ...FieldCondition) error {
	if conditionsMatch(target, conditions, true) {
		return requireField("required_if", target, fieldPath)
	}
	return nil
}

func RequiredUnless(target any, fieldPath string, conditions ...FieldCondition) error {
	if !conditionsMatch(target, conditions, true) {
		return requireField("required_unless", target, fieldPath)
	}
	return nil
}

func RequiredWith(target any, fieldPath string, otherPaths ...string) error {
	if anyFieldPresent(target, otherPaths) {
		return requireField("required_with", target, fieldPath)
	}
	return nil
}

func RequiredWithAll(target any, fieldPath string, otherPaths ...string) error {
	if allFieldsPresent(target, otherPaths) {
		return requireField("required_with_all", target, fieldPath)
	}
	return nil
}

func RequiredWithout(target any, fieldPath string, otherPaths ...string) error {
	if anyFieldMissing(target, otherPaths) {
		return requireField("required_without", target, fieldPath)
	}
	return nil
}

func RequiredWithoutAll(target any, fieldPath string, otherPaths ...string) error {
	if allFieldsMissing(target, otherPaths) {
		return requireField("required_without_all", target, fieldPath)
	}
	return nil
}

func ExcludedIf(target any, fieldPath string, conditions ...FieldCondition) error {
	if conditionsMatch(target, conditions, true) {
		return excludeField("excluded_if", target, fieldPath)
	}
	return nil
}

func ExcludedUnless(target any, fieldPath string, conditions ...FieldCondition) error {
	if !conditionsMatch(target, conditions, true) {
		return excludeField("excluded_unless", target, fieldPath)
	}
	return nil
}

func ExcludedWith(target any, fieldPath string, otherPaths ...string) error {
	if anyFieldPresent(target, otherPaths) {
		return excludeField("excluded_with", target, fieldPath)
	}
	return nil
}

func ExcludedWithAll(target any, fieldPath string, otherPaths ...string) error {
	if allFieldsPresent(target, otherPaths) {
		return excludeField("excluded_with_all", target, fieldPath)
	}
	return nil
}

func ExcludedWithout(target any, fieldPath string, otherPaths ...string) error {
	if anyFieldMissing(target, otherPaths) {
		return excludeField("excluded_without", target, fieldPath)
	}
	return nil
}

func ExcludedWithoutAll(target any, fieldPath string, otherPaths ...string) error {
	if allFieldsMissing(target, otherPaths) {
		return excludeField("excluded_without_all", target, fieldPath)
	}
	return nil
}

func requireField(tag string, target any, fieldPath string) error {
	value, err := fieldInterface(target, fieldPath)
	if err != nil {
		return fail(tag, target, fieldPath, err.Error())
	}
	if !isZero(value) {
		return nil
	}
	return failf(tag, value, fieldPath, "%s is required", fieldPath)
}

func excludeField(tag string, target any, fieldPath string) error {
	value, err := fieldInterface(target, fieldPath)
	if err != nil {
		return fail(tag, target, fieldPath, err.Error())
	}
	if isZero(value) {
		return nil
	}
	return failf(tag, value, fieldPath, "%s must be excluded", fieldPath)
}

func conditionsMatch(target any, conditions []FieldCondition, emptyMatches bool) bool {
	if len(conditions) == 0 {
		return emptyMatches
	}
	for _, condition := range conditions {
		value, err := fieldInterface(target, condition.Field)
		if err != nil || !equalAny(value, condition.Value) {
			return false
		}
	}
	return true
}

func anyFieldPresent(target any, paths []string) bool {
	for _, path := range paths {
		value, err := fieldInterface(target, path)
		if err == nil && !isZero(value) {
			return true
		}
	}
	return false
}

func allFieldsPresent(target any, paths []string) bool {
	if len(paths) == 0 {
		return false
	}
	for _, path := range paths {
		value, err := fieldInterface(target, path)
		if err != nil || isZero(value) {
			return false
		}
	}
	return true
}

func anyFieldMissing(target any, paths []string) bool {
	for _, path := range paths {
		value, err := fieldInterface(target, path)
		if err != nil || isZero(value) {
			return true
		}
	}
	return false
}

func allFieldsMissing(target any, paths []string) bool {
	if len(paths) == 0 {
		return false
	}
	for _, path := range paths {
		value, err := fieldInterface(target, path)
		if err == nil && !isZero(value) {
			return false
		}
	}
	return true
}
