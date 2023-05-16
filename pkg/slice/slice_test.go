package slice

import (
	"reflect"
	"testing"
)

func TestDiffSlice(t *testing.T) {
	type args struct {
		a []int
		b []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{name: "test#1", args: args{a: []int{1, 2, 3}, b: []int{3, 4, 5}}, want: []int{1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DiffSlice(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DiffSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
