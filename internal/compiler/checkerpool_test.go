package compiler

import (
	"slices"
	"testing"
)

func TestGetCheckerAssociationWeights(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		baseWeights  []int
		importCounts []int
		want         []int
	}{
		{
			name:         "normalizes import work to syntax work",
			baseWeights:  []int{100, 50, 25},
			importCounts: []int{0, 1, 3},
			want:         []int{100, 93, 154},
		},
		{
			name:         "no imports preserves base weights",
			baseWeights:  []int{100, 50, 25},
			importCounts: []int{0, 0, 0},
			want:         []int{100, 50, 25},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := getCheckerAssociationWeights(test.baseWeights, test.importCounts)
			if !slices.Equal(got, test.want) {
				t.Fatalf("getCheckerAssociationWeights(%v, %v) = %v, want %v", test.baseWeights, test.importCounts, got, test.want)
			}
		})
	}
}

func TestGetCheckerAssociations(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		if got := getCheckerAssociations(nil, nil, 4); got != nil {
			t.Fatalf("getCheckerAssociations(nil, nil, 4) = %v, want nil", got)
		}
	})

	t.Run("balances disconnected files", func(t *testing.T) {
		t.Parallel()
		got := getCheckerAssociations(
			[]int{1, 1, 1, 1, 1, 1},
			make([][]int, 6),
			3,
		)
		want := []int{0, 1, 2, 0, 1, 2}
		if !slices.Equal(got, want) {
			t.Fatalf("getCheckerAssociations() = %v, want %v", got, want)
		}
	})

	t.Run("keeps dense components together", func(t *testing.T) {
		t.Parallel()
		got := getCheckerAssociations(
			[]int{1, 1, 1, 1, 1, 1},
			[][]int{
				{1, 2},
				{0, 2},
				{0, 1},
				{4, 5},
				{3, 5},
				{3, 4},
			},
			2,
		)
		want := []int{0, 0, 0, 1, 1, 1}
		if !slices.Equal(got, want) {
			t.Fatalf("getCheckerAssociations() = %v, want %v", got, want)
		}
	})

	t.Run("respects weighted balance cap", func(t *testing.T) {
		t.Parallel()
		weights := []int{8, 7, 6, 5, 4, 3, 2, 1}
		got := getCheckerAssociations(weights, make([][]int, len(weights)), 3)
		loads := make([]int, 3)
		for i, checkerIndex := range got {
			loads[checkerIndex] += weights[i]
		}
		for checkerIndex, load := range loads {
			if load > 13 {
				t.Fatalf("checker %d load = %d, want at most 13; associations = %v", checkerIndex, load, got)
			}
		}
	})
}
