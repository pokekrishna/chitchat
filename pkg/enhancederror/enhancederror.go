package enhancederror

import "reflect"

func IsEqual(source error, target error) bool {
	if source == nil && target == nil {
		return true
	} else if reflect.TypeOf(source) == reflect.TypeOf(target) &&
		source.Error() == target.Error() {
		return true
	} else {
		return false
	}
}
