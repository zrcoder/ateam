package ring

// Ring is a fixed-size circular buffer
type Ring[T any] struct {
	buf   []T
	start int
	count int
}

// New creates a new Ring with the given capacity
func New[T any](capacity int) *Ring[T] {
	if capacity <= 0 {
		panic("capacity must be positive")
	}
	return &Ring[T]{
		buf:   make([]T, capacity),
		start: 0,
		count: 0,
	}
}

// PushFront inserts an element at the front, overwriting the oldest element if full
func (r *Ring[T]) PushFront(val T) {
	if len(r.buf) == 0 {
		return
	}
	r.start = (r.start - 1 + len(r.buf)) % len(r.buf)
	r.buf[r.start] = val
	if r.count < len(r.buf) {
		r.count++
	}
}

// PushBack inserts an element at the back, overwriting the oldest element if full
func (r *Ring[T]) PushBack(val T) {
	if len(r.buf) == 0 {
		return
	}
	idx := (r.start + r.count) % len(r.buf)
	r.buf[idx] = val
	if r.count < len(r.buf) {
		r.count++
	} else {
		r.start = (r.start + 1) % len(r.buf)
	}
}

// Clear removes all elements
func (r *Ring[T]) Clear() {
	r.start = 0
	r.count = 0
}

// Len is an alias for Count, returns the number of elements in the ring
func (r *Ring[T]) Len() int {
	return r.count
}

// Range returns an iterator function. Each call returns (value, ok) where ok is false when iteration is done.
func (r *Ring[T]) Range() func() (T, bool) {
	i := 0
	return func() (T, bool) {
		if i >= r.count {
			var zero T
			return zero, false
		}
		idx := (r.start + i) % len(r.buf)
		v := r.buf[idx]
		i++
		return v, true
	}
}

// Slice returns all elements as a slice (oldest to newest)
func (r *Ring[T]) Slice() []T {
	result := make([]T, 0, r.count)
	for i := 0; i < r.count; i++ {
		idx := (r.start + i) % len(r.buf)
		result = append(result, r.buf[idx])
	}
	return result
}
