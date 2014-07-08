package querymaker

import (
	"sort"
	"testing"
)

func TestOrderedHash(t *testing.T) {
	var hash IHash
	func() {
		hash = NewOrderedHash(Hash{
			"foo": 1,
			"bar": 2,
			"baz": 3,
		}, "foo", "bar", "baz")
		expect(t, hash.Keys(), []string{"foo", "bar", "baz"})
	}()
	func() {
		hash = NewOrderedHash(Hash{
			"foo": 1,
			"bar": 2,
			"baz": 3,
		})
		keys := hash.Keys()
		exp := []string{"foo", "bar", "baz"}
		sort.Strings(keys)
		sort.Strings(exp)
		expect(t, keys, exp)
	}()
}
