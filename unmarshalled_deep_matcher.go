package unmarshalledmatchers

import "reflect"

type UnmarshalledDeepMatcher struct {
	Ordered            bool
	InvertOrderingKeys map[string]bool
	Subset             bool
}

func (matcher *UnmarshalledDeepMatcher) deepEqual(a interface{}, b interface{}) (bool, []interface{}){
	return matcher.deepEqualRecursive(a, b, false)
}


func (matcher *UnmarshalledDeepMatcher) deepEqualRecursive(a interface{}, b interface{}, invertOrdering bool) (bool, []interface{}) {
	var errorPath []interface{}
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false, errorPath
	}

	switch a.(type) {
	case []interface{}:
		if (matcher.Ordered && !invertOrdering) || (!matcher.Ordered && invertOrdering){
			return matcher.deepEqualOrderedList(a, b, errorPath)
		} else {
			return matcher.deepEqualUnorderedList(a, b, errorPath)
		}
	case map[string]interface{}:
		return matcher.deepEqualMap(a, b, errorPath)
	default:
		return a == b, errorPath
	}
}

func (matcher *UnmarshalledDeepMatcher) deepEqualMap(a interface{}, b interface{}, errorPath []interface{}) (bool, []interface{}) {
	if matcher.Subset {
		if len(a.(map[string]interface{})) > len(b.(map[string]interface{})) {
			return false, errorPath
		}
	} else {
		if len(a.(map[string]interface{})) != len(b.(map[string]interface{})) {
			return false, errorPath
		}
	}

	for k, v1 := range a.(map[string]interface{}) {
		v2, ok := b.(map[string]interface{})[k]
		if !ok {
			return false, errorPath
		}

		elementEqual, keyPath := matcher.deepEqualRecursive(v1, v2, matcher.InvertOrderingKeys[k])
		if !elementEqual {
			return false, append(keyPath, k)
		}
	}
	return true, errorPath
}

func (matcher *UnmarshalledDeepMatcher) deepEqualUnorderedList(a interface{}, b interface{}, errorPath []interface{}) (bool, []interface{}) {
	matched := make([]bool, len(b.([]interface{})))

	if matcher.Subset {
		if len(a.([]interface{})) > len(b.([]interface{})) {
			return false, errorPath
		}
	} else {
		if len(a.([]interface{})) != len(b.([]interface{})) {
			return false, errorPath
		}
	}

	for _, v1 := range a.([]interface{}) {
		foundMatch := false
		for j, v2 := range b.([]interface{}) {
			if matched[j] {
				continue
			}
			elementEqual, _ := matcher.deepEqualRecursive(v1, v2, false)
			if elementEqual {
				foundMatch = true
				matched[j] = true
				break
			}
		}
		if !foundMatch {
			return false, errorPath
		}
	}

	return true, errorPath
}

func (matcher *UnmarshalledDeepMatcher) deepEqualOrderedList(a interface{}, b interface{}, errorPath []interface{}) (bool, []interface{}) {
	if matcher.Subset {
		if len(a.([]interface{})) > len(b.([]interface{})) {
			return false, errorPath
		}
	} else {
		if len(a.([]interface{})) != len(b.([]interface{})) {
			return false, errorPath
		}
	}

	for i, v := range a.([]interface{}) {
		elementEqual, keyPath := matcher.deepEqualRecursive(v, b.([]interface{})[i], false)
		if !elementEqual {
			return false, append(keyPath, i)
		}
	}
	return true, errorPath
}

