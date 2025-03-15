package querybuilder

import (
	"reflect"
	"testing"
)

func Test_appendSliceWhereIN(t *testing.T) {
	type args struct {
		value  interface{}
		values []interface{}
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			name: "Test 1",
			args: args{
				value: 1,
			},
			want: []interface{}{1},
		},
		{
			name: "Test 2",
			args: args{
				value:  1,
				values: []interface{}{2, 3},
			},
			want: []interface{}{1, 2, 3},
		},
		{
			name: "Test 3",
			args: args{
				value:  nil,
				values: []interface{}{1, 2, 3},
			},
			want: []interface{}{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := appendSliceWhereIN(tt.args.value, tt.args.values...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendSliceWhereIN() = %v, want %v", got, tt.want)
			}
		})
	}
}
