package slot

import (
	"reflect"
	"testing"
)

func TestNewPositionsByPattern(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name string
		args args
		want []Position
	}{
		{
			name: "Normal",
			args: args{
				pattern: "4_1,3_2,2_1",
			},
			want: []Position{
				{Reel: 2, Row: 1},
				{Reel: 3, Row: 2},
				{Reel: 4, Row: 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPositionsByPattern(tt.args.pattern); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPositionsByPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergePositions(t *testing.T) {
	type args struct {
		arrs [][]Position
	}
	tests := []struct {
		name string
		args args
		want []Position
	}{
		{
			name: "Merge Same Reel",
			args: args{
				arrs: [][]Position{
					[]Position{
						{1, 1},
						{1, 3},
						{1, 5},
					},
					[]Position{
						{1, 1},
						{1, 2},
						{1, 3},
					},
				},
			},
			want: []Position{
				{1, 1},
				{1, 2},
				{1, 3},
				{1, 5},
			},
		},
		{
			name: "Merge Multi Reel",
			args: args{
				arrs: [][]Position{
					[]Position{
						{1, 1},
						{2, 3},
						{3, 5},
					},
					[]Position{
						{1, 1},
						{1, 2},
						{2, 3},
					},
					[]Position{
						{1, 3},
						{2, 2},
						{3, 1},
					},
					[]Position{
						{3, 1},
						{5, 5},
						{3, 5},
					},
				},
			},
			want: []Position{
				{1, 1},
				{1, 2},
				{1, 3},
				{2, 2},
				{2, 3},
				{3, 1},
				{3, 5},
				{5, 5},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergePositions(tt.args.arrs...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergePositions() = %v, want %v", got, tt.want)
			}
		})
	}
}
