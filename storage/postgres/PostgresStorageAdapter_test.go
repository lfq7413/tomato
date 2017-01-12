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

func Test_toParseSchema(t *testing.T) {
	type args struct {
		schema types.M
	}
	tests := []struct {
		name string
		args args
		want types.M
	}{
		{
			name: "1",
			args: args{
				schema: nil,
			},
			want: nil,
		},
		{
			name: "2",
			args: args{
				schema: types.M{
					"className": "post",
				},
			},
			want: types.M{
				"className": "post",
				"fields":    types.M{},
				"classLevelPermissions": types.M{
					"find":     types.M{"*": true},
					"get":      types.M{"*": true},
					"create":   types.M{"*": true},
					"update":   types.M{"*": true},
					"delete":   types.M{"*": true},
					"addField": types.M{"*": true},
				},
			},
		},
		{
			name: "3",
			args: args{
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"title": types.M{"type": "String"},
						"auth":  types.M{"type": "String"},
					},
				},
			},
			want: types.M{
				"className": "post",
				"fields": types.M{
					"title": types.M{"type": "String"},
					"auth":  types.M{"type": "String"},
				},
				"classLevelPermissions": types.M{
					"find":     types.M{"*": true},
					"get":      types.M{"*": true},
					"create":   types.M{"*": true},
					"update":   types.M{"*": true},
					"delete":   types.M{"*": true},
					"addField": types.M{"*": true},
				},
			},
		},
		{
			name: "4",
			args: args{
				schema: types.M{
					"className": "_User",
					"fields": types.M{
						"name":             types.M{"type": "String"},
						"_hashed_password": types.M{"type": "String"},
					},
				},
			},
			want: types.M{
				"className": "_User",
				"fields": types.M{
					"name": types.M{"type": "String"},
				},
				"classLevelPermissions": types.M{
					"find":     types.M{"*": true},
					"get":      types.M{"*": true},
					"create":   types.M{"*": true},
					"update":   types.M{"*": true},
					"delete":   types.M{"*": true},
					"addField": types.M{"*": true},
				},
			},
		},
		{
			name: "5",
			args: args{
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"title":  types.M{"type": "String"},
						"_wperm": types.M{"type": "Array"},
						"_rperm": types.M{"type": "Array"},
					},
				},
			},
			want: types.M{
				"className": "post",
				"fields": types.M{
					"title": types.M{"type": "String"},
				},
				"classLevelPermissions": types.M{
					"find":     types.M{"*": true},
					"get":      types.M{"*": true},
					"create":   types.M{"*": true},
					"update":   types.M{"*": true},
					"delete":   types.M{"*": true},
					"addField": types.M{"*": true},
				},
			},
		},
		{
			name: "6",
			args: args{
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"title": types.M{"type": "String"},
					},
					"classLevelPermissions": types.M{
						"addField": types.M{},
					},
				},
			},
			want: types.M{
				"className": "post",
				"fields": types.M{
					"title": types.M{"type": "String"},
				},
				"classLevelPermissions": types.M{
					"find":     types.M{"*": true},
					"get":      types.M{"*": true},
					"create":   types.M{"*": true},
					"update":   types.M{"*": true},
					"delete":   types.M{"*": true},
					"addField": types.M{},
				},
			},
		},
	}
	for _, tt := range tests {
		if got := toParseSchema(tt.args.schema); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. toParseSchema() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_toPostgresSchema(t *testing.T) {
	type args struct {
		schema types.M
	}
	tests := []struct {
		name string
		args args
		want types.M
	}{
		{
			name: "1",
			args: args{
				schema: nil,
			},
			want: nil,
		},
		{
			name: "2",
			args: args{
				schema: types.M{
					"className": "post",
				},
			},
			want: types.M{
				"className": "post",
				"fields": types.M{
					"_wperm": types.M{
						"type":     "Array",
						"contents": types.M{"type": "String"},
					},
					"_rperm": types.M{
						"type":     "Array",
						"contents": types.M{"type": "String"},
					},
				},
			},
		},
		{
			name: "3",
			args: args{
				schema: types.M{
					"className": "_User",
					"fields": types.M{
						"name": types.M{"type": "String"},
					},
				},
			},
			want: types.M{
				"className": "_User",
				"fields": types.M{
					"name": types.M{"type": "String"},
					"_wperm": types.M{
						"type":     "Array",
						"contents": types.M{"type": "String"},
					},
					"_rperm": types.M{
						"type":     "Array",
						"contents": types.M{"type": "String"},
					},
					"_hashed_password":  types.M{"type": "String"},
					"_password_history": types.M{"type": "Array"},
				},
			},
		},
	}
	for _, tt := range tests {
		if got := toPostgresSchema(tt.args.schema); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. toPostgresSchema() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_handleDotFields(t *testing.T) {
	type args struct {
		object types.M
	}
	tests := []struct {
		name string
		args args
		want types.M
	}{
		{
			name: "1",
			args: args{
				object: nil,
			},
			want: nil,
		},
		{
			name: "2",
			args: args{
				object: types.M{
					"key": "hello",
				},
			},
			want: types.M{
				"key": "hello",
			},
		},
		{
			name: "3",
			args: args{
				object: types.M{
					"key.sub": "hello",
					"key2":    "world",
				},
			},
			want: types.M{
				"key": types.M{
					"sub": "hello",
				},
				"key2": "world",
			},
		},
		{
			name: "4",
			args: args{
				object: types.M{
					"key.sub.sub": "hello",
					"key2":        "world",
				},
			},
			want: types.M{
				"key": types.M{
					"sub": types.M{
						"sub": "hello",
					},
				},
				"key2": "world",
			},
		},
		{
			name: "5",
			args: args{
				object: types.M{
					"key.sub": types.M{
						"__op": "Delete",
					},
					"key2": "world",
				},
			},
			want: types.M{
				"key": types.M{
					"sub": nil,
				},
				"key2": "world",
			},
		},
	}
	for _, tt := range tests {
		if got := handleDotFields(tt.args.object); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. handleDotFields() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_validateKeys(t *testing.T) {
	type args struct {
		object interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "1",
			args:    args{object: nil},
			wantErr: nil,
		},
		{
			name: "2",
			args: args{object: types.M{
				"key": "hello",
			}},
			wantErr: nil,
		},
		{
			name: "3",
			args: args{object: types.M{
				"key.sub": "hello",
			}},
			wantErr: errs.E(errs.InvalidNestedKey, "Nested keys should not contain the '$' or '.' characters"),
		},
		{
			name: "4",
			args: args{object: types.M{
				"key$sub": "hello",
			}},
			wantErr: errs.E(errs.InvalidNestedKey, "Nested keys should not contain the '$' or '.' characters"),
		},
		{
			name: "5",
			args: args{object: types.M{
				"key": types.M{
					"sub": "hello",
				},
			}},
			wantErr: nil,
		},
		{
			name: "6",
			args: args{object: types.M{
				"key": types.M{
					"sub.key": "hello",
				},
			}},
			wantErr: errs.E(errs.InvalidNestedKey, "Nested keys should not contain the '$' or '.' characters"),
		},
		{
			name: "7",
			args: args{object: types.M{
				"key": types.M{
					"sub$key": "hello",
				},
			}},
			wantErr: errs.E(errs.InvalidNestedKey, "Nested keys should not contain the '$' or '.' characters"),
		},
	}
	for _, tt := range tests {
		if err := validateKeys(tt.args.object); reflect.DeepEqual(err, tt.wantErr) == false {
			t.Errorf("%q. validateKeys() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
