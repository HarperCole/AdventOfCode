package util

type HashSet struct {
	data map[interface{}]struct{}
}

// NewHashSet creates and returns a new instance of HashSet
func NewHashSet() *HashSet {
	return &HashSet{data: make(map[interface{}]struct{})}
}

// Add adds an element to the set
func (hs *HashSet) Add(value interface{}) {
	hs.data[value] = struct{}{}
}

// Remove removes an element from the set
func (hs *HashSet) Remove(value interface{}) {
	delete(hs.data, value)
}

// Contains checks if an element exists in the set
func (hs *HashSet) Contains(value interface{}) bool {
	_, exists := hs.data[value]
	return exists
}

// Size returns the number of elements in the set
func (hs *HashSet) Size() int {
	return len(hs.data)
}

// Clear removes all elements from the set
func (hs *HashSet) Clear() {
	hs.data = make(map[interface{}]struct{})
}

// ToSlice returns all elements in the set as a slice
func (hs *HashSet) ToSlice() []interface{} {
	keys := make([]interface{}, 0, len(hs.data))
	for key := range hs.data {
		keys = append(keys, key)
	}
	return keys
}

// Union returns a new HashSet containing all elements from all provided sets
func (hs *HashSet) Union(others ...*HashSet) *HashSet {
	result := NewHashSet()
	for key := range hs.data {
		result.Add(key)
	}
	for _, other := range others {
		for key := range other.data {
			result.Add(key)
		}
	}
	return result
}

// Intersection returns a new HashSet containing only elements present in all provided sets
func (hs *HashSet) Intersection(others ...*HashSet) *HashSet {
	result := NewHashSet()
	for key := range hs.data {
		inAll := true
		for _, other := range others {
			if !other.Contains(key) {
				inAll = false
				break
			}
		}
		if inAll {
			result.Add(key)
		}
	}
	return result
}

// Difference returns a new HashSet containing elements present in the first set but not in any of the others
func (hs *HashSet) Difference(others ...*HashSet) *HashSet {
	result := NewHashSet()
	for key := range hs.data {
		inAny := false
		for _, other := range others {
			if other.Contains(key) {
				inAny = true
				break
			}
		}
		if !inAny {
			result.Add(key)
		}
	}
	return result
}

// Equals checks if two sets contain the same elements
func (hs *HashSet) Equals(other *HashSet) bool {
	if hs.Size() != other.Size() {
		return false
	}
	for key := range hs.data {
		if !other.Contains(key) {
			return false
		}
	}
	return true
}
