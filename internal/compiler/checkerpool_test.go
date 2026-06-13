package compiler

import "testing"

func TestGetCheckerAssociationBlockSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		fileCount    int
		checkerCount int
		want         int
	}{
		{name: "single checker uses max block", fileCount: 100, checkerCount: 1, want: maxCheckerAssociationBlockSize},
		{name: "small project balances across checkers", fileCount: 16, checkerCount: 4, want: 1},
		{name: "medium project uses smaller locality blocks", fileCount: 128, checkerCount: 4, want: 8},
		{name: "large project caps block size", fileCount: 2000, checkerCount: 4, want: maxCheckerAssociationBlockSize},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := getCheckerAssociationBlockSize(test.fileCount, test.checkerCount); got != test.want {
				t.Fatalf("getCheckerAssociationBlockSize(%d, %d) = %d, want %d", test.fileCount, test.checkerCount, got, test.want)
			}
		})
	}
}

func TestGetCheckerAssociationsForFileWeights(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		fileWeights  []int
		checkerCount int
		want         []int
	}{
		{
			name:         "small project cycles each file",
			fileWeights:  []int{1, 1, 1, 1},
			checkerCount: 4,
			want:         []int{0, 1, 2, 3},
		},
		{
			name:         "large first file sends following files to lighter checker",
			fileWeights:  []int{100, 1, 1, 1, 1},
			checkerCount: 2,
			want:         []int{0, 1, 1, 1, 1},
		},
		{
			name: "medium project keeps contiguous locality blocks",
			fileWeights: []int{
				1, 1, 1, 1, 1, 1, 1, 1,
				1, 1, 1, 1, 1, 1, 1, 1,
				1, 1, 1, 1, 1, 1, 1, 1,
				1, 1, 1, 1, 1, 1, 1, 1,
			},
			checkerCount: 4,
			want: []int{
				0, 0, 1, 1, 2, 2, 3, 3,
				0, 0, 1, 1, 2, 2, 3, 3,
				0, 0, 1, 1, 2, 2, 3, 3,
				0, 0, 1, 1, 2, 2, 3, 3,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := getCheckerAssociationsForFileWeights(test.fileWeights, test.checkerCount)
			if len(got) != len(test.want) {
				t.Fatalf("len(getCheckerAssociationsForFileWeights(%v, %d)) = %d, want %d", test.fileWeights, test.checkerCount, len(got), len(test.want))
			}
			for i := range got {
				if got[i] != test.want[i] {
					t.Fatalf("getCheckerAssociationsForFileWeights(%v, %d)[%d] = %d, want %d", test.fileWeights, test.checkerCount, i, got[i], test.want[i])
				}
			}
		})
	}
}

func TestRefineCheckerAssociationsByGraph(t *testing.T) {
	t.Parallel()

	t.Run("moves file to import-neighbor checker within balance cap", func(t *testing.T) {
		t.Parallel()

		associations := []int{0, 0, 0, 1, 1}
		refineCheckerAssociationsByGraph(
			associations,
			[]int{1, 1, 1, 1, 1},
			[][]int{{3, 4}, nil, nil, {0}, {0}},
			2,
		)
		want := []int{1, 0, 0, 1, 1}
		for i := range want {
			if associations[i] != want[i] {
				t.Fatalf("associations[%d] = %d, want %d; associations = %v", i, associations[i], want[i], associations)
			}
		}
	})

	t.Run("does not move file past balance cap", func(t *testing.T) {
		t.Parallel()

		associations := []int{0, 0, 1, 1}
		refineCheckerAssociationsByGraph(
			associations,
			[]int{1, 1, 1, 1},
			[][]int{{2, 3}, nil, {0}, {0}},
			2,
		)
		want := []int{0, 0, 1, 1}
		for i := range want {
			if associations[i] != want[i] {
				t.Fatalf("associations[%d] = %d, want %d; associations = %v", i, associations[i], want[i], associations)
			}
		}
	})
}
