package util

import (
	"testing"
)

func TestHashSet_AddAndContains(t *testing.T) {
	set := NewHashSet()

	// Test adding and checking integers
	set.Add(1)
	set.Add(2)
	set.Add(3)

	if !set.Contains(1) {
		t.Error("Expected set to contain 1")
	}
	if !set.Contains(2) {
		t.Error("Expected set to contain 2")
	}
	if set.Contains(4) {
		t.Error("Did not expect set to contain 4")
	}

	// Test adding and checking arrays
	array1 := [2]int{1, 2}
	array2 := [2]int{3, 4}
	set.Add(array1)
	set.Add(array2)

	if !set.Contains(array1) {
		t.Error("Expected set to contain [1, 2]")
	}
	if !set.Contains(array2) {
		t.Error("Expected set to contain [3, 4]")
	}
}

func TestHashSet_Remove(t *testing.T) {
	set := NewHashSet()

	// Test removing integers
	set.Add(1)
	set.Add(2)
	set.Remove(1)

	if set.Contains(1) {
		t.Error("Did not expect set to contain 1 after removal")
	}
	if !set.Contains(2) {
		t.Error("Expected set to contain 2")
	}

	// Test removing arrays
	array1 := [2]int{1, 2}
	set.Add(array1)
	set.Remove(array1)

	if set.Contains(array1) {
		t.Error("Did not expect set to contain [1, 2] after removal")
	}
}

func TestHashSet_Size(t *testing.T) {
	set := NewHashSet()

	set.Add(1)
	set.Add(2)
	set.Add(3)

	if set.Size() != 3 {
		t.Errorf("Expected size to be 3, got %d", set.Size())
	}

	// Test size after adding duplicates
	set.Add(2)
	if set.Size() != 3 {
		t.Errorf("Expected size to still be 3 after adding duplicate, got %d", set.Size())
	}
}

func TestHashSet_Clear(t *testing.T) {
	set := NewHashSet()

	set.Add(1)
	set.Add(2)
	set.Clear()

	if set.Size() != 0 {
		t.Errorf("Expected size to be 0 after clear, got %d", set.Size())
	}

	if set.Contains(1) {
		t.Error("Did not expect set to contain 1 after clear")
	}
}

func TestHashSet_ToSlice(t *testing.T) {
	set := NewHashSet()

	set.Add(1)
	set.Add(2)
	set.Add(3)

	elements := set.ToSlice()
	expected := map[interface{}]bool{
		1: true,
		2: true,
		3: true,
	}

	for _, elem := range elements {
		if !expected[elem] {
			t.Errorf("Unexpected element in set: %v", elem)
		}
	}
}
