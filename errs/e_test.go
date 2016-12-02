package errs

import (
	"errors"
	"reflect"
	"testing"

	"github.com/lfq7413/tomato/types"
)

func TestE(t *testing.T) {
	type args struct {
		code int
		msg  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
	}{
		{
			name:    "E 1",
			args:    args{code: 20, msg: "hello"},
			wantErr: `{"code": 20,"error": "hello"}`,
		},
	}
	for _, tt := range tests {
		if err := E(tt.args.code, tt.args.msg); err.Error() != tt.wantErr {
			t.Errorf("%q. E() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestErrorToMap(t *testing.T) {
	type args struct {
		e error
	}
	tests := []struct {
		name string
		args args
		want types.M
	}{
		{
			name: "ErrorToMap 1",
			args: args{E(20, "hello")},
			want: types.M{"code": 20, "error": "hello"},
		},
		{
			name: "ErrorToMap 2",
			args: args{errors.New("hello")},
			want: types.M{"code": -1, "error": "hello"},
		},
	}
	for _, tt := range tests {
		if got := ErrorToMap(tt.args.e); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. ErrorToMap() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestGetErrorCode(t *testing.T) {
	type args struct {
		e error
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "GetErrorCode 1",
			args: args{E(20, "hello")},
			want: 20,
		},
		{
			name: "GetErrorCode 1",
			args: args{errors.New("hello")},
			want: 0,
		},
	}
	for _, tt := range tests {
		if got := GetErrorCode(tt.args.e); got != tt.want {
			t.Errorf("%q. GetErrorCode() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestGetErrorMessage(t *testing.T) {
	type args struct {
		e error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "GetErrorCode 1",
			args: args{E(20, "hello")},
			want: "hello",
		},
		{
			name: "GetErrorCode 1",
			args: args{errors.New("hello")},
			want: "hello",
		},
	}
	for _, tt := range tests {
		if got := GetErrorMessage(tt.args.e); got != tt.want {
			t.Errorf("%q. GetErrorMessage() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
