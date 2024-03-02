package cache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshal(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr string
	}{
		{"OK - nil", args{nil}, []byte{}, ""},
		{"OK - string", args{"test"}, []byte(`"test"`), ""},
		{"OK - int", args{1}, []byte("1"), ""},
		{"OK - float", args{1.1}, []byte("1.1"), ""},
		{"OK - bool", args{true}, []byte("true"), ""},
		{"OK - struct", args{struct{ A string }{"test"}}, []byte(`{"A":"test"}`), ""},
		{"KO - error", args{make(chan int)}, nil, "CACHE.VALUE.MARSHAL.ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			got, err := Marshal(tt.args.v)
			require.Equal(st, tt.want, got)
			if tt.wantErr != "" {
				require.ErrorContains(st, err, tt.wantErr)
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	type args struct {
		data []byte
		v    any
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
	}{
		{"OK - nil", args{[]byte{}, nil}, ""},
		{"OK - string", args{[]byte(`"test"`), new(string)}, ""},
		{"OK - int", args{[]byte("1"), new(int)}, ""},
		{"OK - float", args{[]byte("1.1"), new(float64)}, ""},
		{"OK - bool", args{[]byte("true"), new(bool)}, ""},
		{"OK - struct", args{[]byte(`{"A":"test"}`), new(struct{ A string })}, ""},
		{"KO - error", args{[]byte("test"), new(chan int)}, "CACHE.VALUE.UNMARSHAL.ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			err := Unmarshal(tt.args.data, tt.args.v)
			if tt.wantErr != "" {
				require.ErrorContains(st, err, tt.wantErr)
			}
		})
	}
}
