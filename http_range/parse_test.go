package http_range

import (
	"reflect"
	"testing"
)

func Test_parseRange(t *testing.T) {
	type args struct {
		s    string
		size int64
	}
	tests := []struct {
		name    string
		args    args
		want    []httpRange
		wantErr bool
	}{
		{
			name: "first 10 bytes",
			args: args{
				s:    "bytes=0-9",
				size: 10,
			},
			want: []httpRange{
				{
					start:  0,
					length: 10,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseRange(tt.args.s, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseRange() = %v, want %v", got, tt.want)
			}
		})
	}
}
