package unmarshalledmatchers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestUnmarshalledMatchers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "unmarshalledmatchers Suite")
}
