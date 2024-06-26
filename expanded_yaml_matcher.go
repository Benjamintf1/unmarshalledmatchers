package unmarshalledmatchers

import (
	"fmt"
	"strings"

	"github.com/onsi/gomega/format"
	"gopkg.in/yaml.v3"
)

type ExpandedYAMLMatcher struct {
	YAMLToMatch      interface{}
	firstFailurePath []interface{}
	DeepMatcher      UnmarshalledDeepMatcher
}

func (matcher *ExpandedYAMLMatcher) Match(actual interface{}) (success bool, err error) {
	actualString, expectedString, err := matcher.toStrings(actual)
	if err != nil {
		return false, err
	}

	var aval interface{}
	var eval interface{}

	if err := yaml.Unmarshal([]byte(actualString), &aval); err != nil {
		return false, fmt.Errorf("Actual '%s' should be valid YAML, but it is not.\nUnderlying error:%s", actualString, err)
	}
	if err := yaml.Unmarshal([]byte(expectedString), &eval); err != nil {
		return false, fmt.Errorf("Expected '%s' should be valid YAML, but it is not.\nUnderlying error:%s", expectedString, err)
	}
	var equal bool

	equal, matcher.firstFailurePath = matcher.DeepMatcher.deepEqual(eval, aval)
	return equal, nil
}

func (matcher *ExpandedYAMLMatcher) FailureMessage(actual interface{}) (message string) {
	actualString, expectedString, _ := matcher.toNormalisedStrings(actual)
	return format.Message(actualString, "to match YAML of", expectedString)
}

func (matcher *ExpandedYAMLMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	actualString, expectedString, _ := matcher.toNormalisedStrings(actual)
	return format.Message(actualString, "not to match YAML of", expectedString)
}

func (matcher *ExpandedYAMLMatcher) toNormalisedStrings(actual interface{}) (actualFormatted, expectedFormatted string, err error) {
	actualString, expectedString, err := matcher.toStrings(actual)
	return normalise(actualString), normalise(expectedString), err
}

func normalise(input string) string {
	var val interface{}
	err := yaml.Unmarshal([]byte(input), &val)
	if err != nil {
		panic(err) // guarded by Match
	}
	output, err := yaml.Marshal(val)
	if err != nil {
		panic(err) // guarded by Unmarshal
	}
	return strings.TrimSpace(string(output))
}

func (matcher *ExpandedYAMLMatcher) toStrings(actual interface{}) (actualFormatted, expectedFormatted string, err error) {
	actualString, ok := toString(actual)
	if !ok {
		return "", "", fmt.Errorf("ExpandedYAMLMatcher matcher requires a string, stringer, or []byte.  Got actual:\n%s", format.Object(actual, 1))
	}
	expectedString, ok := toString(matcher.YAMLToMatch)
	if !ok {
		return "", "", fmt.Errorf("ExpandedYAMLMatcher matcher requires a string, stringer, or []byte.  Got expected:\n%s", format.Object(matcher.YAMLToMatch, 1))
	}

	return actualString, expectedString, nil
}
