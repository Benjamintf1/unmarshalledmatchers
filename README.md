Gomega matchers for more complicated json matching
==================================

This package provides [Gomega](https://github.com/onsi/gomega) matchers to match against either subsets of json/xml/yml, as well 
as json/xml/yml that has unordered lists. You can see an example of what that means below.

The functions of this library should follow the following pattern

```
(Match|Contain)(Unordered|Unordered)(JSON|XML|YAML)( aStructureToUnmarshall, optional (With(Un)orderedKeys(keys))
```

JsonMatchers()
-------------------
```go 
import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/Benjamintf1/ExpandedUnmarshalledMatchers"
)

//Match with exception keys
Expect(`{"a":[1,2,3],"b":[1,2,3],"c":[1,2,3]}`).Should(
	MatchUnorderedJSON(`{"a":[1,2,3],"b":[3,2,1],"c":[3,2,1]}`,
		WithOrderedListKeys("a"))) 
		
Expect(`{"a":[1,2,3],"b":[1,2,3],"c":[1,2,3]}`).Should(
	MatchOrderedJSON(`{"a":[3,2,1],"b":[1,2,3],"c":[1,2,3]}`,
		WithUnorderedListKeys("a")))

//Contain with exception keys
Expect(`{"a":[1,2,3],"b":[1,2,3],"c":[1,2,3]}`).Should(
	ContainUnorderedJSON(`{"a":[1,2,3],"b":[3,2,1]}`,
		WithOrderedListKeys("a")))
		
Expect(`{"a":[1,2,3],"b":[1,2,3],"c":[1,2,3]}`).Should(
	ContainOrderedJSON(`{"a":[3,2,1],"b":[1,2,3]}`, 
		WithUnorderedListKeys("a")))


```
