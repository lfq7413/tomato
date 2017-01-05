package postgres

import (
	"reflect"
	"testing"

	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
)

func Test_parseTypeToPostgresType(t *testing.T) {
	type args struct {
		t types.M
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{
			name:    "1",
			args:    args{t: nil},
			want:    "",
			wantErr: nil,
		},
		{
			name:    "2",
			args:    args{t: types.M{"type": "String"}},
			want:    "text",
			wantErr: nil,
		},
		{
			name:    "3",
			args:    args{t: types.M{"type": "Date"}},
			want:    "timestamp with time zone",
			wantErr: nil,
		},
		{
			name:    "4",
			args:    args{t: types.M{"type": "Object"}},
			want:    "jsonb",
			wantErr: nil,
		},
		{
			name:    "5",
			args:    args{t: types.M{"type": "File"}},
			want:    "text",
			wantErr: nil,
		},
		{
			name:    "6",
			args:    args{t: types.M{"type": "Boolean"}},
			want:    "boolean",
			wantErr: nil,
		},
		{
			name:    "7",
			args:    args{t: types.M{"type": "Pointer"}},
			want:    "char(10)",
			wantErr: nil,
		},
		{
			name:    "8",
			args:    args{t: types.M{"type": "Number"}},
			want:    "double precision",
			wantErr: nil,
		},
		{
			name:    "9",
			args:    args{t: types.M{"type": "GeoPoint"}},
			want:    "point",
			wantErr: nil,
		},
		{
			name: "10",
			args: args{
				t: types.M{
					"type":     "Array",
					"contents": types.M{"type": "String"},
				},
			},
			want:    "text[]",
			wantErr: nil,
		},
		{
			name:    "11",
			args:    args{t: types.M{"type": "Array"}},
			want:    "jsonb",
			wantErr: nil,
		},
		{
			name:    "12",
			args:    args{t: types.M{}},
			want:    "",
			wantErr: errs.E(errs.IncorrectType, "no type for  yet"),
		},
		{
			name:    "13",
			args:    args{t: types.M{"type": "Other"}},
			want:    "",
			wantErr: errs.E(errs.IncorrectType, "no type for Other yet"),
		},
	}
	for _, tt := range tests {
		got, err := parseTypeToPostgresType(tt.args.t)
		if reflect.DeepEqual(err, tt.wantErr) == false {
			t.Errorf("%q. parseTypeToPostgresType() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. parseTypeToPostgresType() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_toPostgresValue(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "1",
			args: args{value: nil},
			want: nil,
		},
		{
			name: "2",
			args: args{value: "hello"},
			want: "hello",
		},
		{
			name: "3",
			args: args{
				value: types.M{"key": "value"},
			},
			want: types.M{"key": "value"},
		},
		{
			name: "4",
			args: args{
				value: types.M{"__type": "Other"},
			},
			want: types.M{"__type": "Other"},
		},
		{
			name: "5",
			args: args{
				value: types.M{
					"__type": "Date",
					"iso":    "2006-01-02T15:04:05.000Z",
				},
			},
			want: "2006-01-02T15:04:05.000Z",
		},
		{
			name: "6",
			args: args{
				value: types.M{
					"__type": "File",
					"name":   "image.jpg",
				},
			},
			want: "image.jpg",
		},
	}
	for _, tt := range tests {
		if got := toPostgresValue(tt.args.value); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. toPostgresValue() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_transformValue(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "1",
			args: args{value: nil},
			want: nil,
		},
		{
			name: "2",
			args: args{value: "hello"},
			want: "hello",
		},
		{
			name: "3",
			args: args{
				value: types.M{"key": "value"},
			},
			want: types.M{"key": "value"},
		},
		{
			name: "4",
			args: args{
				value: types.M{"__type": "Other"},
			},
			want: types.M{"__type": "Other"},
		},
		{
			name: "5",
			args: args{
				value: types.M{
					"__type":   "Pointer",
					"objectId": "1024",
				},
			},
			want: "1024",
		},
	}
	for _, tt := range tests {
		if got := transformValue(tt.args.value); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. transformValue() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
