package compiler

import "testing"

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
			name:         "heaviest files are scheduled first",
			fileWeights:  []int{1, 100, 1, 1},
			checkerCount: 2,
			want:         []int{1, 0, 1, 1},
		},
		{
			name:         "descending weights balance across least-loaded checker",
			fileWeights:  []int{8, 7, 6, 5, 4, 3, 2, 1},
			checkerCount: 3,
			want:         []int{0, 1, 2, 2, 1, 0, 0, 1},
		},
		{
			name:         "ties are assigned in file order",
			fileWeights:  []int{1, 1, 1, 1, 1, 1, 1, 1},
			checkerCount: 4,
			want:         []int{0, 1, 2, 3, 0, 1, 2, 3},
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
