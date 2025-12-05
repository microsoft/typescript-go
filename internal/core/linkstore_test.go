package core_test

import (
	"sync"
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"gotest.tools/v3/assert"
)

func TestLinkStoreGet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		setup    func(*core.LinkStore[string, int])
		key      string
		wantZero bool // whether we expect zero value (first access)
	}{
		{
			name:     "creates new value on first access",
			setup:    nil,
			key:      "key1",
			wantZero: true,
		},
		{
			name: "returns existing value",
			setup: func(s *core.LinkStore[string, int]) {
				v := s.Get("key1")
				*v = 42
			},
			key:      "key1",
			wantZero: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var store core.LinkStore[string, int]
			if tt.setup != nil {
				tt.setup(&store)
			}

			value := store.Get(tt.key)
			assert.Assert(t, value != nil, "Get should return a non-nil pointer")

			if tt.wantZero {
				assert.Equal(t, *value, 0, "New value should be zero-initialized")
			} else {
				assert.Equal(t, *value, 42, "Should return existing value")
			}
		})
	}
}

func TestLinkStoreGetPointerIdentity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		key1        string
		key2        string
		wantSamePtr bool
	}{
		{
			name:        "same key returns same pointer",
			key1:        "key1",
			key2:        "key1",
			wantSamePtr: true,
		},
		{
			name:        "different keys return different pointers",
			key1:        "key1",
			key2:        "key2",
			wantSamePtr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var store core.LinkStore[string, int]

			ptr1 := store.Get(tt.key1)
			ptr2 := store.Get(tt.key2)

			if tt.wantSamePtr {
				assert.Assert(t, ptr1 == ptr2, "Expected same pointer")
			} else {
				assert.Assert(t, ptr1 != ptr2, "Expected different pointers")
			}
		})
	}
}

func TestLinkStoreHas(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		setup    func(*core.LinkStore[string, int])
		key      string
		expected bool
	}{
		{
			name:     "returns false for missing key",
			setup:    nil,
			key:      "missing",
			expected: false,
		},
		{
			name: "returns true after Get",
			setup: func(s *core.LinkStore[string, int]) {
				s.Get("key1")
			},
			key:      "key1",
			expected: true,
		},
		{
			name: "returns false for different key after Get",
			setup: func(s *core.LinkStore[string, int]) {
				s.Get("key1")
			},
			key:      "key2",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var store core.LinkStore[string, int]
			if tt.setup != nil {
				tt.setup(&store)
			}

			result := store.Has(tt.key)
			assert.Equal(t, result, tt.expected)
		})
	}
}

func TestLinkStoreTryGet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(*core.LinkStore[string, int])
		key       string
		wantNil   bool
		wantValue int
	}{
		{
			name:    "returns nil for missing key",
			setup:   nil,
			key:     "missing",
			wantNil: true,
		},
		{
			name: "returns value after Get",
			setup: func(s *core.LinkStore[string, int]) {
				v := s.Get("key1")
				*v = 42
			},
			key:       "key1",
			wantNil:   false,
			wantValue: 42,
		},
		{
			name: "returns zero value for unmodified key",
			setup: func(s *core.LinkStore[string, int]) {
				s.Get("key1") // Get but don't modify
			},
			key:       "key1",
			wantNil:   false,
			wantValue: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var store core.LinkStore[string, int]
			if tt.setup != nil {
				tt.setup(&store)
			}

			value := store.TryGet(tt.key)

			if tt.wantNil {
				assert.Assert(t, value == nil, "Expected nil")
			} else {
				assert.Assert(t, value != nil, "Expected non-nil")
				assert.Equal(t, *value, tt.wantValue)
			}
		})
	}
}

func TestLinkStoreConcurrent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "concurrent Get with same key returns same pointer",
			test: func(t *testing.T) {
				var store core.LinkStore[string, int]
				const goroutines = 100

				var wg sync.WaitGroup
				results := make([]*int, goroutines)

				wg.Add(goroutines)
				for i := range goroutines {
					go func(idx int) {
						defer wg.Done()
						results[idx] = store.Get("shared-key")
					}(i)
				}
				wg.Wait()

				first := results[0]
				for i := 1; i < goroutines; i++ {
					assert.Assert(t, results[i] == first,
						"All concurrent Get calls for the same key should return the same pointer")
				}
			},
		},
		{
			name: "concurrent Get with different keys returns different pointers",
			test: func(t *testing.T) {
				var store core.LinkStore[int, int]
				const goroutines = 100

				var wg sync.WaitGroup
				results := make([]*int, goroutines)

				wg.Add(goroutines)
				for i := range goroutines {
					go func(idx int) {
						defer wg.Done()
						results[idx] = store.Get(idx)
					}(i)
				}
				wg.Wait()

				seen := make(map[*int]bool)
				for i := range goroutines {
					assert.Assert(t, !seen[results[i]],
						"Concurrent Get calls for different keys should return different pointers")
					seen[results[i]] = true
				}
			},
		},
		{
			name: "concurrent mixed operations",
			test: func(t *testing.T) {
				var store core.LinkStore[int, int]
				const goroutines = 100
				const iterations = 100

				// Pre-populate some keys
				for i := range 10 {
					store.Get(i)
				}

				var wg sync.WaitGroup
				wg.Add(goroutines)

				for i := range goroutines {
					go func(idx int) {
						defer wg.Done()
						for j := range iterations {
							key := (idx + j) % 20

							switch j % 3 {
							case 0:
								store.Get(key)
							case 1:
								store.Has(key)
							case 2:
								store.TryGet(key)
							}
						}
					}(i)
				}
				wg.Wait()
				// If we got here without panic, the test passed
			},
		},
		{
			name: "concurrent writes to different keys",
			test: func(t *testing.T) {
				var store core.LinkStore[int, int]
				const goroutines = 100
				const keysPerGoroutine = 10

				var wg sync.WaitGroup
				wg.Add(goroutines)

				for i := range goroutines {
					go func(idx int) {
						defer wg.Done()
						for j := range keysPerGoroutine {
							key := idx*keysPerGoroutine + j
							value := store.Get(key)
							*value = key
						}
					}(i)
				}
				wg.Wait()

				// Verify all values
				for i := range goroutines * keysPerGoroutine {
					value := store.Get(i)
					assert.Equal(t, *value, i, "Value should match the key")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.test(t)
		})
	}
}

func TestLinkStoreWithStruct(t *testing.T) {
	t.Parallel()

	type TestStruct struct {
		Name  string
		Value int
	}

	tests := []struct {
		name      string
		setup     func(*core.LinkStore[string, TestStruct])
		key       string
		wantName  string
		wantValue int
	}{
		{
			name:      "zero-initialized struct",
			setup:     nil,
			key:       "key1",
			wantName:  "",
			wantValue: 0,
		},
		{
			name: "modified struct persists",
			setup: func(s *core.LinkStore[string, TestStruct]) {
				v := s.Get("key1")
				v.Name = "test"
				v.Value = 42
			},
			key:       "key1",
			wantName:  "test",
			wantValue: 42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var store core.LinkStore[string, TestStruct]
			if tt.setup != nil {
				tt.setup(&store)
			}

			value := store.Get(tt.key)
			assert.Equal(t, value.Name, tt.wantName)
			assert.Equal(t, value.Value, tt.wantValue)
		})
	}
}
