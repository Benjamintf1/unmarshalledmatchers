package unmarshalledmatchers

import "reflect"

type UnmarshalledDeepMatcher struct {
	Ordered            bool
	InvertOrderingKeys map[interface{}]bool
	Subset             bool
}

type yamlMap map[interface{}]interface{}
type jsonMap map[string]interface{}

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
			return matcher.deepEqualOrderedList(a.([]interface{}), b.([]interface{}), errorPath)
		} else {
			return matcher.deepEqualUnorderedList(a.([]interface{}), b.([]interface{}), errorPath)
		}
	case jsonMap:
		return matcher.deepEqualMap(toYamlMap(a.(jsonMap)), toYamlMap(b.(jsonMap)), errorPath)
	case yamlMap:
		return matcher.deepEqualMap(a.(yamlMap), b.(yamlMap), errorPath)
	default:
		return a == b, errorPath
	}
}

func (matcher *UnmarshalledDeepMatcher) deepEqualMap(a yamlMap, b yamlMap, errorPath []interface{}) (bool, []interface{}) {
	if matcher.Subset {
		if len(a) > len(b) {
			return false, errorPath
		}
	} else {
		if len(a) != len(b) {
			return false, errorPath
		}
	}

	for k, v1 := range a {
		v2, ok := b[k]
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

func (matcher *UnmarshalledDeepMatcher) deepEqualUnorderedList(a []interface{}, b []interface{}, errorPath []interface{}) (bool, []interface{}) {
	matched := make([]bool, len(b))

	if matcher.Subset {
		if len(a) > len(b) {
			return false, errorPath
		}
	} else {
		if len(a) != len(b) {
			return false, errorPath
		}
	}

	for _, v1 := range a {
		foundMatch := false
		for j, v2 := range b {
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

func (matcher *UnmarshalledDeepMatcher) deepEqualOrderedList(a []interface{}, b []interface{}, errorPath []interface{}) (bool, []interface{}) {
	if matcher.Subset {
		if len(a) > len(b) {
			return false, errorPath
		}
	} else {
		if len(a) != len(b) {
			return false, errorPath
		}
	}

	for i, v := range a {
		elementEqual, keyPath := matcher.deepEqualRecursive(v, b[i], false)
		if !elementEqual {
			return false, append(keyPath, i)
		}
	}
	return true, errorPath
}

func toYamlMap(map1 jsonMap) (yamlMap){
	convert := make(yamlMap)
	for key, value := range map1 {
		convert[key] = value
	}
	return convert
}

