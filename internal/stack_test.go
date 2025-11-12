package internal

import (
	"testing"
)

func TestNewStack_IsEmpty_ByInitialLen(t *testing.T) {
	type testCase struct {
		name      string
		initial   []int
		wantEmpty bool
	}
	tests := []testCase{
		{"len=0", nil, true},
		{"len=1", []int{42}, false},
		{"len=2", []int{1, 2}, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := NewStack(tc.initial)
			if got := s.IsEmpty(); got != tc.wantEmpty {
				t.Fatalf("IsEmpty()=%v, want %v (initial=%v)", got, tc.wantEmpty, tc.initial)
			}
			// sanity check: internal slice matches initializer (same-length copy)
			if len(s.items) != len(tc.initial) {
				t.Fatalf("len(items)=%d, want %d", len(s.items), len(tc.initial))
			}
		})
	}
}

func TestPush_From0_1_2(t *testing.T) {
	type testCase struct {
		name    string
		initial []int
		push    []int
		want    []int // final items slice content
	}
	tests := []testCase{
		{"start0_push1", nil, []int{7}, []int{7}},
		{"start0_push2", nil, []int{1, 2}, []int{1, 2}},
		{"start1_push1", []int{9}, []int{3}, []int{9, 3}},
		{"start2_push1", []int{4, 5}, []int{6}, []int{4, 5, 6}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := NewStack(append([]int(nil), tc.initial...)) // copy to avoid aliasing
			for _, v := range tc.push {
				s.Push(v)
			}
			if len(s.items) != len(tc.want) {
				t.Fatalf("len(items)=%d, want %d", len(s.items), len(tc.want))
			}
			for i := range tc.want {
				if s.items[i] != tc.want[i] {
					t.Fatalf("items[%d]=%v, want %v; full=%v want=%v", i, s.items[i], tc.want[i], s.items, tc.want)
				}
			}
			if (len(tc.want) == 0) != s.IsEmpty() {
				t.Fatalf("IsEmpty()=%v, want %v", s.IsEmpty(), len(tc.want) == 0)
			}
		})
	}
}

func TestTop_Empty_0(t *testing.T) {
	s := NewStack[int](nil)
	got, err := s.Top()
	if err == nil {
		t.Fatalf("Top() on empty: expected error, got nil")
	}
	if want := 0; got != want {
		t.Fatalf("Top() zero-value: got %v, want %v", got, want)
	}
	if !s.IsEmpty() {
		t.Fatalf("Top() on empty should keep stack empty")
	}
}

func TestTop_Single_1(t *testing.T) {
	s := NewStack([]int{99})
	got, err := s.Top()
	if err != nil {
		t.Fatalf("Top() error: %v", err)
	}
	if got != 99 {
		t.Fatalf("Top() value: got %v, want 99", got)
	}
	if !s.IsEmpty() {
		t.Fatalf("stack should be empty after popping the only element")
	}
	// Another Top should error
	_, err = s.Top()
	if err == nil {
		t.Fatalf("expected error after popping last element, got nil")
	}
}

func TestTop_LIFO_2(t *testing.T) {
	s := NewStack([]int{10, 20}) // 10 (bottom), 20 (top)

	// First pop -> 20
	v, err := s.Top()
	if err != nil {
		t.Fatalf("Top() #1 error: %v", err)
	}
	if v != 20 {
		t.Fatalf("Top() #1 value: got %v, want 20", v)
	}
	if s.IsEmpty() {
		t.Fatalf("stack should not be empty after first pop")
	}

	// Second pop -> 10
	v, err = s.Top()
	if err != nil {
		t.Fatalf("Top() #2 error: %v", err)
	}
	if v != 10 {
		t.Fatalf("Top() #2 value: got %v, want 10", v)
	}
	if !s.IsEmpty() {
		t.Fatalf("stack should be empty after popping both elements")
	}

	// Third pop (now empty) -> error + zero value
	v, err = s.Top()
	if err == nil {
		t.Fatalf("Top() on empty after 2 pops: expected error, got nil")
	}
	if v != 0 {
		t.Fatalf("Top() zero-value after empty: got %v, want 0", v)
	}
}

func TestEndToEnd_0_1_2_ThenPushAndPop(t *testing.T) {
	// Start with 0
	s := NewStack[int](nil)
	if !s.IsEmpty() {
		t.Fatalf("new stack should be empty")
	}

	// Push 1
	s.Push(1)
	if s.IsEmpty() {
		t.Fatalf("stack should not be empty after push")
	}

	// Push making it 2
	s.Push(2)
	if len(s.items) != 2 {
		t.Fatalf("len(items)=%d, want 2", len(s.items))
	}

	// Pop -> 2 then 1
	if v, err := s.Top(); err != nil || v != 2 {
		t.Fatalf("Top #1: got (%v,%v), want (2,nil)", v, err)
	}
	if v, err := s.Top(); err != nil || v != 1 {
		t.Fatalf("Top #2: got (%v,%v), want (1,nil)", v, err)
	}
	if !s.IsEmpty() {
		t.Fatalf("should be empty after popping all")
	}
}
