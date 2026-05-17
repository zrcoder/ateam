package ring

import "testing"

func TestPushBack(t *testing.T) {
	r := New[int](3)
	r.PushBack(1)
	r.PushBack(2)
	r.PushBack(3)

	var result []int
	iter := r.Range()
	for v, ok := iter(); ok; v, ok = iter() {
		result = append(result, v)
	}
	if len(result) != 3 {
		t.Errorf("len(result) = %d, want 3", len(result))
	}
}

func TestPushFront(t *testing.T) {
	r := New[int](3)
	r.PushFront(1)
	r.PushFront(2)
	r.PushFront(3)

	var result []int
	iter := r.Range()
	for v, ok := iter(); ok; v, ok = iter() {
		result = append(result, v)
	}
	if len(result) != 3 {
		t.Errorf("len(result) = %d, want 3", len(result))
	}
	if result[0] != 3 || result[2] != 1 {
		t.Errorf("Range() = %v, want [3 2 1]", result)
	}
}

func TestPushBackOverwrite(t *testing.T) {
	r := New[int](3)
	r.PushBack(1)
	r.PushBack(2)
	r.PushBack(3)
	r.PushBack(4)

	var result []int
	iter := r.Range()
	for v, ok := iter(); ok; v, ok = iter() {
		result = append(result, v)
	}
	if result[0] != 2 || result[1] != 3 || result[2] != 4 {
		t.Errorf("Range() = %v, want [2 3 4]", result)
	}
}

func TestPushFrontOverwrite(t *testing.T) {
	r := New[int](3)
	r.PushFront(1)
	r.PushFront(2)
	r.PushFront(3)
	r.PushFront(4)

	var result []int
	iter := r.Range()
	for v, ok := iter(); ok; v, ok = iter() {
		result = append(result, v)
	}
	if result[0] != 4 || result[1] != 3 || result[2] != 2 {
		t.Errorf("Range() = %v, want [4 3 2]", result)
	}
}

func TestClear(t *testing.T) {
	r := New[int](3)
	r.PushBack(1)
	r.PushBack(2)
	r.Clear()

	var count int
	iter := r.Range()
	for v, ok := iter(); ok; v, ok = iter() {
		_ = v
		count++
	}
	if count != 0 {
		t.Errorf("Range after Clear() = %d elements, want 0", count)
	}
}

func TestRangeEmpty(t *testing.T) {
	r := New[int](3)

	var count int
	iter := r.Range()
	for v, ok := iter(); ok; v, ok = iter() {
		_ = v
		count++
	}
	if count != 0 {
		t.Errorf("Range on empty ring = %d iterations, want 0", count)
	}
}

func TestRangeOrder(t *testing.T) {
	r := New[int](4)
	for i := 1; i <= 4; i++ {
		r.PushBack(i)
	}

	i := 1
	iter := r.Range()
	for v, ok := iter(); ok; v, ok = iter() {
		if v != i {
			t.Errorf("element at position = %d, want %d", v, i)
		}
		i++
	}
}

func TestRangeConsumeAll(t *testing.T) {
	r := New[int](3)
	r.PushBack(1)
	r.PushBack(2)

	iter := r.Range()
	for v, ok := iter(); ok; v, ok = iter() {
		_ = v
	}

	var count int
	for v, ok := iter(); ok; v, ok = iter() {
		_ = v
		count++
	}
	if count != 0 {
		t.Errorf("After consume, second iteration = %d, want 0", count)
	}
}

func TestLen(t *testing.T) {
	r := New[int](5)
	if r.Len() != 0 {
		t.Errorf("empty ring Len() = %d, want 0", r.Len())
	}

	r.PushBack(1)
	if r.Len() != 1 {
		t.Errorf("after 1 PushBack Len() = %d, want 1", r.Len())
	}

	r.PushBack(2)
	r.PushBack(3)
	if r.Len() != 3 {
		t.Errorf("after 3 PushBack Len() = %d, want 3", r.Len())
	}

	r.PushBack(4)
	r.PushBack(5)
	if r.Len() != 5 {
		t.Errorf("after 5 PushBack Len() = %d, want 5", r.Len())
	}

	r.PushBack(6) // overwrite oldest
	if r.Len() != 5 {
		t.Errorf("after 6th PushBack (overflow) Len() = %d, want 5", r.Len())
	}
}

func TestSlice(t *testing.T) {
	r := New[int](5)
	empty := r.Slice()
	if len(empty) != 0 {
		t.Errorf("empty ring Slice() len = %d, want 0", len(empty))
	}

	r.PushBack(10)
	r.PushBack(20)
	r.PushBack(30)

	s := r.Slice()
	if len(s) != 3 {
		t.Errorf("Slice() len = %d, want 3", len(s))
	}
	if s[0] != 10 || s[1] != 20 || s[2] != 30 {
		t.Errorf("Slice() = %v, want [10 20 30]", s)
	}
}

func TestSliceOrder(t *testing.T) {
	r := New[int](3)
	r.PushBack(1)
	r.PushBack(2)
	r.PushBack(3)

	s := r.Slice()
	if len(s) != 3 {
		t.Errorf("Slice() len = %d, want 3", len(s))
	}
	// PushBack adds at back, so oldest is at start (index 0)
	if s[0] != 1 || s[1] != 2 || s[2] != 3 {
		t.Errorf("Slice() = %v, want [1 2 3]", s)
	}
}

func TestSliceAfterOverflow(t *testing.T) {
	r := New[int](3)
	r.PushBack(1)
	r.PushBack(2)
	r.PushBack(3)
	r.PushBack(4) // overwrites 1

	s := r.Slice()
	if len(s) != 3 {
		t.Errorf("Slice() len = %d, want 3", len(s))
	}
	// After overflow: 1 is gone, [2, 3, 4] remain (oldest to newest)
	if s[0] != 2 || s[1] != 3 || s[2] != 4 {
		t.Errorf("Slice() = %v, want [2 3 4]", s)
	}
}

func TestSliceAfterPushFront(t *testing.T) {
	r := New[string](3)
	r.PushFront("a")
	r.PushFront("b")
	r.PushFront("c")

	s := r.Slice()
	// PushFront: newest at front (start), so order is c, b, a (newest to oldest)
	// Actually for PushFront, the newest goes to front, oldest at back
	if s[0] != "c" || s[1] != "b" || s[2] != "a" {
		t.Errorf("Slice() = %v, want [c b a]", s)
	}
}

func TestCapacity(t *testing.T) {
	r := New[int](3)
	if len(r.buf) != 3 {
		t.Errorf("len(r.buf) = %d, want 3", len(r.buf))
	}
}
