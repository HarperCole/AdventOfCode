package until

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
