package unmarshalledmatchers_test

import (
	. "github.com/benjamintf1/unmarshalledmatchers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MatchUnorderedJSONMatcher", func() {
	Context("When passed stringifiables", func() {
		It("should succeed if the JSON matches", func() {
			Ω("{}").Should(MatchUnorderedJSON("{}"))
			Ω(`{"a":1}`).Should(MatchUnorderedJSON(`{"a":1}`))
			Ω(`{
			             "a":1
			         }`).Should(MatchUnorderedJSON(`{"a":1}`))
			Ω(`{"a":1, "b":2}`).Should(MatchUnorderedJSON(`{"b":2, "a":1}`))
			Ω(`{"a":1}`).ShouldNot(MatchUnorderedJSON(`{"b":2, "a":1}`))

			Ω(`{"a":"a", "b":"b"}`).ShouldNot(MatchUnorderedJSON(`{"a":"a", "b":"b", "c":"c"}`))
			Ω(`{"a":"a", "b":"b", "c":"c"}`).ShouldNot(MatchUnorderedJSON(`{"a":"a", "b":"b"}`))

			Ω(`{"a":null, "b":null}`).ShouldNot(MatchUnorderedJSON(`{"c":"c", "d":"d"}`))
			Ω(`{"a":null, "b":null, "c":null}`).ShouldNot(MatchUnorderedJSON(`{"a":null, "b":null, "d":null}`))
		})

		It("should work with byte arrays", func() {
			Ω([]byte("{}")).Should(MatchUnorderedJSON([]byte("{}")))
			Ω("{}").Should(MatchUnorderedJSON([]byte("{}")))
			Ω([]byte("{}")).Should(MatchUnorderedJSON("{}"))
		})
	})

	Context("When there are arrays that are unordered", func() {
		It("should succeed if the JSON matches", func() {
			Ω(`{"a":[1,2,3]}`).Should(MatchUnorderedJSON(`{"a":[3,2,1]}`))
			Ω(`{"a":[1,2,3],"b":[1,2,3],"c":[1,2,3]}`).Should(MatchUnorderedJSON(`{"a":[3,2,1],"b":[3,2,1],"c":[3,2,1]}`))
			Ω(`[[1,2,3],[4,5,6],[7,8,9]]`).Should(MatchUnorderedJSON(`[[9,8,7],[6,5,4],[3,2,1]]`))
			Ω(`[[1,2,3],[1,2,3],[1,2,3]]`).Should(MatchUnorderedJSON(`[[3,2,1],[3,2,1],[3,2,1]]`))
		})
	})

	Context("When some of the keys are ordered", func() {
		It("should succeed if the JSON matches", func() {
			Ω(`{"a":[1,2,3],"b":[1,2,3],"c":[1,2,3]}`).Should(MatchUnorderedJSON(`{"a":[1,2,3],"b":[3,2,1],"c":[3,2,1]}`, WithOrderedListKeys("a")))
		})

		It("should not succeed if a ordered key doesn't match", func() {
			Ω(`{"a":[1,2,3],"b":[1,2,3],"c":[1,2,3]}`).ShouldNot(MatchUnorderedJSON(`{"a":[3,2,1],"b":[3,2,1],"c":[3,2,1]}`, WithOrderedListKeys("a")))
		})
	})

	Context("When some of the keys are Unordered", func() {
		It("should succeed if the JSON matches", func() {
			Ω(`{"a":[1,2,3],"b":[1,2,3],"c":[1,2,3]}`).Should(MatchOrderedJSON(`{"a":[3,2,1],"b":[1,2,3],"c":[1,2,3]}`, WithUnorderedListKeys("a")))
		})
	})

	Context("SubsetMatching", func() {
		It("should succeed if the JSON is contained", func() {
			Ω(`{"a":[1,2,3],"b":[1,2,3],"c":[1,2,3]}`).Should(ContainOrderedJSON(`{"b":[1,2,3]}`))
			Ω(`{"a":[1,2,3],"b":[1,2,3],"c":[1,2,3]}`).Should(ContainUnorderedJSON(`{"b":[3,2,1]}`))

			Ω(`{"a":[1,2,3],"b":[1,2,3],"c":[1,2,3]}`).Should(ContainUnorderedJSON(`{"a":[1,2,3],"b":[3,2,1]}`, WithOrderedListKeys("a")))
			Ω(`{"a":[1,2,3],"b":[1,2,3],"c":[1,2,3]}`).Should(ContainOrderedJSON(`{"a":[3,2,1],"b":[1,2,3]}`, WithUnorderedListKeys("a")))

		})

		It("should Not succeed if the JSON isnt contained", func() {
			Ω(`{"a":[1,2,3],"b":[1,2,3],"c":[1,2,3]}`).ShouldNot(ContainOrderedJSON(`{"b":[3,2,1]}`))
			Ω(`{"b":[3,2,1]}`).ShouldNot(ContainOrderedJSON(`{"a":[1,2,3],"b":[1,2,3],"c":[1,2,3]}`))
			Ω(`{"a":[1,2,3],"b":[1,2,3],"c":[1,2,3]}`).ShouldNot(ContainUnorderedJSON(`{"a":[3,2,1],"b":[3,2,1]}`, WithOrderedListKeys("a")))
		})
	})

	Context("when a key mismatch is found", func() {
		It("reports the first found mismatch", func() {
			subject := ExpandedJsonMatcher{JSONToMatch: `5`}
			actual := `7`
			subject.Match(actual)

			failureMessage := subject.FailureMessage(`7`)
			Ω(failureMessage).ToNot(ContainSubstring("first mismatched key"))

			subject = ExpandedJsonMatcher{JSONToMatch: `{"a": 1, "b.g": {"c": 2, "1": ["hello", "see ya"]}}`}
			actual = `{"a": 1, "b.g": {"c": 2, "1": ["hello", "goodbye"]}}`
			subject.Match(actual)

			failureMessage = subject.FailureMessage(actual)
			Ω(failureMessage).To(ContainSubstring(`first mismatched key: "b.g"."1"`)) //Only report on the array, not the index
		})

		It("reports the first found mismatch as Json does when the key is ordered", func() {
			subject := ExpandedJsonMatcher{JSONToMatch: `5`}
			actual := `7`
			subject.Match(actual)

			failureMessage := subject.FailureMessage(`7`)
			Ω(failureMessage).ToNot(ContainSubstring("first mismatched key"))

			subject = ExpandedJsonMatcher{
				JSONToMatch: `{"a": 1, "b.g": {"c": 2, "1": ["hello", "see ya"]}}`,
				DeepMatcher: UnmarshalledDeepMatcher{
					Ordered:            false,
					Subset:             false,
					InvertOrderingKeys: map[interface{}]bool{"1": true},
				},
			}
			actual = `{"a": 1, "b.g": {"c": 2, "1": ["hello", "goodbye"]}}`
			subject.Match(actual)

			failureMessage = subject.FailureMessage(actual)
			Ω(failureMessage).To(ContainSubstring(`first mismatched key: "b.g"."1"[1]`))
		})
	})

	Context("when the expected is not valid JSON", func() {
		It("should error and explain why", func() {
			success, err := (&ExpandedJsonMatcher{JSONToMatch: `{}`}).Match(`oops`)
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("Actual 'oops' should be valid JSON"))
		})
	})

	Context("when the actual is not valid JSON", func() {
		It("should error and explain why", func() {
			success, err := (&ExpandedJsonMatcher{JSONToMatch: `oops`}).Match(`{}`)
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("Expected 'oops' should be valid JSON"))
		})
	})

	Context("when the expected is neither a string nor a stringer nor a byte array", func() {
		It("should error", func() {
			success, err := (&ExpandedJsonMatcher{JSONToMatch: 2}).Match("{}")
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("ExpandedJsonMatcher matcher requires a string, stringer, or []byte.  Got expected:\n    <int>: 2"))

			success, err = (&ExpandedJsonMatcher{JSONToMatch: nil}).Match("{}")
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("ExpandedJsonMatcher matcher requires a string, stringer, or []byte.  Got expected:\n    <nil>: nil"))
		})
	})

	Context("when the actual is neither a string nor a stringer nor a byte array", func() {
		It("should error", func() {
			success, err := (&ExpandedJsonMatcher{JSONToMatch: "{}"}).Match(2)
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("ExpandedJsonMatcher matcher requires a string, stringer, or []byte.  Got actual:\n    <int>: 2"))

			success, err = (&ExpandedJsonMatcher{JSONToMatch: "{}"}).Match(nil)
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring("ExpandedJsonMatcher matcher requires a string, stringer, or []byte.  Got actual:\n    <nil>: nil"))
		})
	})
})
