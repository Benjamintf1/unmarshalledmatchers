package ExpandedUnmarshalledMatchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"io/ioutil"
	"fmt"
	"os"
)

func TestExpandedUnmarshalledMatchers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ExpandedUnmarshalledMatchers Suite")
}



func readFileContents(filePath string) []byte {
	f := openFile(filePath)
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(fmt.Errorf("failed to read file contents: %v", err))
	}
	return b
}

func openFile(filePath string) *os.File {
	f, err := os.Open(filePath)
	if err != nil {
		panic(fmt.Errorf("failed to open file: %v", err))
	}
	return f
}


type myStringer struct {
	a string
}

func (s *myStringer) String() string {
	return s.a
}

type StringAlias string

type myCustomType struct {
	s   string
	n   int
	f   float32
	arr []string
}