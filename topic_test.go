package sps

import (
	"reflect"
	"testing"
)

func Test_shift(t *testing.T) {
	type args struct {
		data   [][]byte
		offset int
	}
	tests := []struct {
		name string
		args args
		want [][]byte
	}{
		{
			name: "nil data",
			args: args{
				data:   nil,
				offset: 0,
			},
			want: nil,
		}, {
			name: "empty data",
			args: args{
				data:   [][]byte{},
				offset: 0,
			},
			want: [][]byte{},
		}, {
			name: "normal shift",
			args: args{
				data:   [][]byte{[]byte{0x01}, []byte{0x02}},
				offset: 1,
			},
			want: [][]byte{[]byte{0x02}},
		}, {
			name: "full shift",
			args: args{
				data:   [][]byte{[]byte{0x01}, []byte{0x02}},
				offset: 2,
			},
			want: [][]byte{},
		}, {
			name: "out of bounds offset",
			args: args{
				data:   [][]byte{[]byte{0x01}, []byte{0x02}},
				offset: 3,
			},
			want: [][]byte{},
		}, {
			name: "negative offset",
			args: args{
				data:   [][]byte{[]byte{0x01}, []byte{0x02}},
				offset: -3,
			},
			want: [][]byte{[]byte{0x01}, []byte{0x02}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shift(tt.args.data, tt.args.offset); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("shift() = %v, want %v", got, tt.want)
			}
		})
	}
}
