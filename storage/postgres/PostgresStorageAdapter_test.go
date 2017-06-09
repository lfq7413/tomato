package postgres

import (
	"database/sql"
	"log"
	"reflect"
	"testing"
	"time"

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
			want:    "char(24)",
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
				"key":  types.M{},
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

func Test_joinTablesForSchema(t *testing.T) {
	type args struct {
		schema types.M
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "1",
			args: args{schema: nil},
			want: []string{},
		},
		{
			name: "2",
			args: args{schema: types.M{"className": "post"}},
			want: []string{},
		},
		{
			name: "3",
			args: args{schema: types.M{
				"className": "post",
				"fields": types.M{
					"name": types.M{"type": "String"},
				},
			}},
			want: []string{},
		},
		{
			name: "4",
			args: args{schema: types.M{
				"className": "post",
				"fields": types.M{
					"name": types.M{"type": "String"},
					"user": types.M{"type": "Relation"},
				},
			}},
			want: []string{"_Join:user:post"},
		},
	}
	for _, tt := range tests {
		if got := joinTablesForSchema(tt.args.schema); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. joinTablesForSchema() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_buildWhereClause(t *testing.T) {
	type args struct {
		schema types.M
		query  types.M
		index  int
	}
	tests := []struct {
		name    string
		args    args
		want    *whereClause
		wantErr error
	}{
		{
			name: "1",
			args: args{
				schema: nil,
				query:  nil,
				index:  1,
			},
			want: &whereClause{
				pattern: "",
				values:  types.S{},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "2",
			args: args{
				schema: types.M{},
				query:  types.M{},
				index:  1,
			},
			want: &whereClause{
				pattern: "",
				values:  types.S{},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "3",
			args: args{
				schema: types.M{},
				query: types.M{
					"key": types.M{"$exists": false},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: "",
				values:  types.S{},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "4",
			args: args{
				schema: types.M{},
				query: types.M{
					"key.sub": "hello",
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `"key"->'sub' = '"hello"'`,
				values:  types.S{},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "4.1",
			args: args{
				schema: types.M{},
				query: types.M{
					"key.sub": true,
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `"key"->'sub' = 'true'`,
				values:  types.S{},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "4.2",
			args: args{
				schema: types.M{},
				query: types.M{
					"key.sub": 1024,
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `"key"->'sub' = '1024'`,
				values:  types.S{},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "4.3",
			args: args{
				schema: types.M{},
				query: types.M{
					"key.sub": types.M{"key": "hello"},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `"key"->'sub' = '{"key":"hello"}'`,
				values:  types.S{},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "4.4",
			args: args{
				schema: types.M{},
				query: types.M{
					"key.sub": types.S{"hello", "world"},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `"key"->'sub' = '["hello","world"]'`,
				values:  types.S{},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "5",
			args: args{
				schema: types.M{},
				query: types.M{
					"key.subkey.sub": "hello",
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `"key"->'subkey'->'sub' = '"hello"'`,
				values:  types.S{},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "6",
			args: args{
				schema: types.M{},
				query: types.M{
					"key": "hello",
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `"key" = $1`,
				values:  types.S{"hello"},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "7",
			args: args{
				schema: types.M{},
				query: types.M{
					"key": true,
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `"key" = $1`,
				values:  types.S{true},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "8",
			args: args{
				schema: types.M{},
				query: types.M{
					"key": 10.24,
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `"key" = $1`,
				values:  types.S{10.24},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "9",
			args: args{
				schema: types.M{},
				query: types.M{
					"key": 1024,
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `"key" = $1`,
				values:  types.S{1024},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "10",
			args: args{
				schema: types.M{},
				query: types.M{
					"$or": types.S{
						types.M{"key": 10},
						types.M{"key": 20},
					},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `("key" = $1 OR "key" = $2)`,
				values:  types.S{10, 20},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "11",
			args: args{
				schema: types.M{},
				query: types.M{
					"$and": types.S{
						types.M{"key": 10},
						types.M{"key": 20},
					},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `("key" = $1 AND "key" = $2)`,
				values:  types.S{10, 20},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "12",
			args: args{
				schema: types.M{
					"fields": types.M{
						"key": types.M{"type": "Array"},
					},
				},
				query: types.M{
					"key": types.M{
						"$ne": "hello",
					},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `NOT array_contains("key", $1)`,
				values:  types.S{`["hello"]`},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "12.1",
			args: args{
				schema: types.M{
					"fields": types.M{
						"key": types.M{"type": "Array"},
					},
				},
				query: types.M{
					"key": types.M{
						"$ne": nil,
					},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `NOT array_contains("key", $1)`,
				values:  types.S{`[null]`},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "13",
			args: args{
				schema: types.M{},
				query: types.M{
					"key": types.M{
						"$ne": nil,
					},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `"key" IS NOT NULL`,
				values:  types.S{},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "14",
			args: args{
				schema: types.M{},
				query: types.M{
					"key": types.M{
						"$ne": "hello",
					},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `("key" <> $1 OR "key" IS NULL)`,
				values:  types.S{"hello"},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "15",
			args: args{
				schema: types.M{},
				query: types.M{
					"key": types.M{
						"$eq": "hello",
					},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `"key" = $1`,
				values:  types.S{"hello"},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "15.1",
			args: args{
				schema: types.M{},
				query: types.M{
					"key": types.M{
						"$eq": nil,
					},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `"key" IS NULL`,
				values:  types.S{},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "16",
			args: args{
				schema: types.M{
					"fields": types.M{
						"key": types.M{
							"type": "Array",
							"contents": types.M{
								"type": "String",
							},
						},
					},
				},
				query: types.M{
					"key": types.M{
						"$in": types.S{"hello", "world"},
					},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `("key" && ARRAY[$1,$2])`,
				values:  types.S{"hello", "world"},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "17",
			args: args{
				schema: types.M{
					"fields": types.M{
						"key": types.M{
							"type": "Array",
							"contents": types.M{
								"type": "String",
							},
						},
					},
				},
				query: types.M{
					"key": types.M{
						"$in": types.S{"hello", nil, "world"},
					},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `("key" IS NULL OR "key" && ARRAY[$1,$2])`,
				values:  types.S{"hello", "world"},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "18",
			args: args{
				schema: types.M{
					"fields": types.M{
						"key": types.M{
							"type": "Array",
						},
					},
				},
				query: types.M{
					"key": types.M{
						"$in": types.S{"hello", "world"},
					},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: ` array_contains("key", $1)`,
				values:  types.S{`["hello","world"]`},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "19",
			args: args{
				schema: types.M{
					"fields": types.M{
						"key": types.M{
							"type": "String",
						},
					},
				},
				query: types.M{
					"key": types.M{
						"$in": types.S{"hello", "world"},
					},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `"key"  IN ($1,$2)`,
				values:  types.S{"hello", "world"},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "20",
			args: args{
				schema: types.M{
					"fields": types.M{
						"key": types.M{
							"type": "String",
						},
					},
				},
				query: types.M{
					"key": types.M{
						"$in": types.S{},
					},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `"key" IS NULL`,
				values:  types.S{},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "21",
			args: args{
				schema: types.M{
					"fields": types.M{
						"key": types.M{
							"type": "Array",
						},
					},
				},
				query: types.M{
					"key": types.M{
						"$nin": types.S{"hello", "world"},
					},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: ` NOT  array_contains("key", $1)`,
				values:  types.S{`["hello","world"]`},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "22",
			args: args{
				schema: types.M{
					"fields": types.M{
						"key": types.M{
							"type": "String",
						},
					},
				},
				query: types.M{
					"key": types.M{
						"$nin": types.S{"hello", "world"},
					},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `"key"  NOT  IN ($1,$2)`,
				values:  types.S{"hello", "world"},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "23",
			args: args{
				schema: types.M{
					"fields": types.M{
						"key": types.M{
							"type": "Array",
						},
					},
				},
				query: types.M{
					"key": types.M{
						"$all": types.S{"hello", "world"},
					},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `array_contains_all("key", $1::jsonb)`,
				values:  types.S{`["hello","world"]`},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "24",
			args: args{
				schema: types.M{
					"fields": types.M{},
				},
				query: types.M{
					"key": types.M{
						"$exists": true,
					},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `"key" IS NOT NULL`,
				values:  types.S{},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "25",
			args: args{
				schema: types.M{
					"fields": types.M{
						"key": types.M{
							"type": "String",
						},
					},
				},
				query: types.M{
					"key": types.M{
						"$exists": false,
					},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `"key" IS NULL`,
				values:  types.S{},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "26",
			args: args{
				schema: types.M{
					"fields": types.M{},
				},
				query: types.M{
					"key": types.M{
						"$nearSphere": types.M{
							"longitude": 10.0,
							"latitude":  10.0,
						},
						"$maxDistance": 1.0,
					},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `ST_distance_sphere("key"::geometry, POINT($1, $2)::geometry) <= $3`,
				values:  types.S{10.0, 10.0, 6371000.0},
				sorts:   []string{`ST_distance_sphere("key"::geometry, POINT($1, $2)::geometry) ASC`},
			},
			wantErr: nil,
		},
		{
			name: "27",
			args: args{
				schema: types.M{
					"fields": types.M{},
				},
				query: types.M{
					"key": types.M{
						"$within": types.M{
							"$box": types.S{
								types.M{
									"longitude": 10.0,
									"latitude":  20.0,
								},
								types.M{
									"longitude": 20.0,
									"latitude":  10.0,
								},
							},
						},
					},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `"key"::point <@ $1::box`,
				values:  types.S{"((10, 20), (20, 10))"},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "27.1",
			args: args{
				schema: types.M{
					"fields": types.M{},
				},
				query: types.M{
					"key": types.M{
						"$geoWithin": types.M{
							"$polygon": types.S{
								types.M{
									"__type":    "GeoPoint",
									"longitude": 10.0,
									"latitude":  20.0,
								},
								types.M{
									"__type":    "GeoPoint",
									"longitude": 20.0,
									"latitude":  10.0,
								},
								types.M{
									"__type":    "GeoPoint",
									"longitude": 20.0,
									"latitude":  20.0,
								},
							},
						},
					},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `"key"::point <@ $1::polygon`,
				values:  types.S{"((10, 20), (20, 10), (20, 20))"},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "28",
			args: args{
				schema: types.M{
					"fields": types.M{},
				},
				query: types.M{
					"key": types.M{
						"$regex": `abc`,
					},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `"key" ~ 'abc'`,
				values:  types.S{},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "29",
			args: args{
				schema: types.M{
					"fields": types.M{},
				},
				query: types.M{
					"key": types.M{
						"$regex":   `abc`,
						"$options": "i",
					},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `"key" ~* 'abc'`,
				values:  types.S{},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "30",
			args: args{
				schema: types.M{
					"fields": types.M{},
				},
				query: types.M{
					"key": types.M{
						"$regex":   `abc efg`,
						"$options": "x",
					},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `"key" ~ 'abcefg'`,
				values:  types.S{},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "31",
			args: args{
				schema: types.M{
					"fields": types.M{
						"key": types.M{
							"type": "Array",
						},
					},
				},
				query: types.M{
					"key": types.M{
						"__type":   "Pointer",
						"objectId": "1024",
					},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `array_contains("key", $1)`,
				values:  types.S{`[{"__type":"Pointer","objectId":"1024"}]`},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "32",
			args: args{
				schema: types.M{
					"fields": types.M{},
				},
				query: types.M{
					"key": types.M{
						"__type":   "Pointer",
						"objectId": "hello",
					},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `"key" = $1`,
				values:  types.S{"hello"},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "33",
			args: args{
				schema: types.M{
					"fields": types.M{},
				},
				query: types.M{
					"key": types.M{
						"__type": "Date",
						"iso":    "2017-01-02T15:04:05.000Z",
					},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `"key" = $1`,
				values:  types.S{"2017-01-02T15:04:05.000Z"},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "34",
			args: args{
				schema: types.M{
					"fields": types.M{},
				},
				query: types.M{
					"key": types.M{
						"$gt": 1024,
					},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `"key" > $1`,
				values:  types.S{1024},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "35",
			args: args{
				schema: types.M{
					"fields": types.M{},
				},
				query: types.M{
					"key": types.M{
						"$lt": 1024,
					},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `"key" < $1`,
				values:  types.S{1024},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "36",
			args: args{
				schema: types.M{
					"fields": types.M{},
				},
				query: types.M{
					"key": types.M{
						"$gte": 1024,
					},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `"key" >= $1`,
				values:  types.S{1024},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "37",
			args: args{
				schema: types.M{
					"fields": types.M{},
				},
				query: types.M{
					"key": types.M{
						"$lte": 1024,
					},
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `"key" <= $1`,
				values:  types.S{1024},
				sorts:   []string{},
			},
			wantErr: nil,
		},
		{
			name: "38",
			args: args{
				schema: types.M{
					"fields": types.M{},
				},
				query: types.M{
					"key": types.S{"hello"},
				},
				index: 1,
			},
			want:    nil,
			wantErr: errs.E(errs.OperationForbidden, `Postgres doesn't support this query type yet ["hello"]`),
		},
		{
			name: "39",
			args: args{
				schema: types.M{
					"fields": types.M{},
				},
				query: types.M{
					"key": nil,
				},
				index: 1,
			},
			want: &whereClause{
				pattern: `"key" IS NULL`,
				values:  types.S{},
				sorts:   []string{},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		got, err := buildWhereClause(tt.args.schema, tt.args.query, tt.args.index)
		if !reflect.DeepEqual(err, tt.wantErr) {
			t.Errorf("%q. buildWhereClause() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. buildWhereClause() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_removeWhiteSpace(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{s: ""},
			want: "",
		},
		{
			name: "2",
			args: args{s: "abc#abc\ndef#def\n"},
			want: "abcdef",
		},
		{
			name: "3",
			args: args{s: "#abc\n#def\n"},
			want: "",
		},
		{
			name: "4",
			args: args{s: "abc   def"},
			want: "abcdef",
		},
		{
			name: "5",
			args: args{s: "   abc"},
			want: "abc",
		},
		{
			name: "6",
			args: args{s: "abc#def\n#ghi\n  jkl  mno   "},
			want: "abcjklmno",
		},
	}
	for _, tt := range tests {
		if got := removeWhiteSpace(tt.args.s); got != tt.want {
			t.Errorf("%q. removeWhiteSpace() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_processRegexPattern(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{s: ""},
			want: "",
		},
		{
			name: "2",
			args: args{s: `^abc\Qedf$^\E`},
			want: `^abcedf\$\^`,
		},
		{
			name: "3",
			args: args{s: `abc\Qedf$^\E$`},
			want: `abcedf\$\^$`,
		},
		{
			name: "4",
			args: args{s: `abc\Qedf$^\E`},
			want: `abcedf\$\^`,
		},
	}
	for _, tt := range tests {
		if got := processRegexPattern(tt.args.s); got != tt.want {
			t.Errorf("%q. processRegexPattern() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_createLiteralRegex(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{s: ""},
			want: "",
		},
		{
			name: "2",
			args: args{s: "a"},
			want: "a",
		},
		{
			name: "3",
			args: args{s: "abcDEF123"},
			want: "abcDEF123",
		},
		{
			name: "4",
			args: args{s: "abc'edf'"},
			want: "abc''edf''",
		},
		{
			name: "5",
			args: args{s: "abc^$"},
			want: `abc\^\$`,
		},
	}
	for _, tt := range tests {
		if got := createLiteralRegex(tt.args.s); got != tt.want {
			t.Errorf("%q. createLiteralRegex() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_literalizeRegexPart(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{s: ""},
			want: "",
		},
		{
			name: "2",
			args: args{s: `\Q\E`},
			want: "",
		},
		{
			name: "3",
			args: args{s: `\Qabc\E`},
			want: "abc",
		},
		{
			name: "4",
			args: args{s: `\Q\abc\E`},
			want: `\\abc`,
		},
		{
			name: "5",
			args: args{s: `abc\Qabc\E`},
			want: "abcabc",
		},
		{
			name: "6",
			args: args{s: `\Q`},
			want: "",
		},
		{
			name: "7",
			args: args{s: `\Qabc`},
			want: "abc",
		},
		{
			name: "8",
			args: args{s: `\Q\abc`},
			want: `\\abc`,
		},
		{
			name: "9",
			args: args{s: `abc\Qabc`},
			want: "abcabc",
		},
		{
			name: "10",
			args: args{s: `abc\Q\Eabc\E`},
			want: "abcabc",
		},
		{
			name: "11",
			args: args{s: `abc\Q\Eabc\E`},
			want: "abcabc",
		},
		{
			name: "12",
			args: args{s: `\Eabc\Q\Eabc\E`},
			want: "abcabc",
		},
		{
			name: "13",
			args: args{s: `\Q\Eabc\E`},
			want: "abc",
		},
		{
			name: "14",
			args: args{s: `'abc'`},
			want: "''abc''",
		},
		{
			name: "15",
			args: args{s: `abc`},
			want: "abc",
		},
		{
			name: "16",
			args: args{s: `\Q'abc'\E`},
			want: "''abc''",
		},
		{
			name: "17",
			args: args{s: `\Q$^*\E`},
			want: `\$\^\*`,
		},
	}
	for _, tt := range tests {
		if got := literalizeRegexPart(tt.args.s); got != tt.want {
			t.Errorf("%q. literalizeRegexPart() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func openDB() *sql.DB {
	db, err := sql.Open("postgres", "postgres://postgres:123456@192.168.99.100:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func TestPostgresAdapter_ensureSchemaCollectionExists(t *testing.T) {
	db := openDB()
	p := NewPostgresAdapter("", db)
	tests := []struct {
		name       string
		wantErr    error
		initialize func()
		clean      func()
	}{
		{
			name:       "1",
			wantErr:    nil,
			initialize: func() {},
			clean: func() {
				db.Exec(`DROP TABLE "_SCHEMA"`)
			},
		},
	}
	for _, tt := range tests {
		tt.initialize()
		if err := p.ensureSchemaCollectionExists(); reflect.DeepEqual(err, tt.wantErr) == false {
			t.Errorf("%q. PostgresAdapter.ensureSchemaCollectionExists() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
		if p.ClassExists("_SCHEMA") == false {
			t.Errorf("%q. PostgresAdapter.ensureSchemaCollectionExists() _SCHEMA is not Exists", tt.name)
		}
		tt.clean()
	}
	db.Close()
}

func TestPostgresAdapter_ClassExists(t *testing.T) {
	db := openDB()
	p := NewPostgresAdapter("", db)
	type args struct {
		name string
	}
	tests := []struct {
		name       string
		args       args
		want       bool
		initialize func(string)
		clean      func(string)
	}{
		{
			name:       "1",
			args:       args{name: "post"},
			want:       false,
			initialize: func(name string) {},
			clean:      func(name string) {},
		},
		{
			name: "2",
			args: args{name: "post"},
			want: true,
			initialize: func(name string) {
				db.Exec(`CREATE TABLE IF NOT EXISTS "` + name + `" ( "title" varChar(120) )`)
			},
			clean: func(name string) {
				db.Exec(`DROP TABLE "` + name + `"`)
			},
		},
	}
	for _, tt := range tests {
		tt.initialize(tt.args.name)
		if got := p.ClassExists(tt.args.name); got != tt.want {
			t.Errorf("%q. PostgresAdapter.ClassExists() = %v, want %v", tt.name, got, tt.want)
		}
		tt.clean(tt.args.name)
	}
	db.Close()
}

func TestPostgresAdapter_createTable(t *testing.T) {
	db := openDB()
	p := NewPostgresAdapter("", db)
	clean := func(name string) {
		db.Exec(`DROP TABLE "` + name + `"`)
		db.Exec(`DROP TABLE "_SCHEMA"`)
	}
	type args struct {
		className string
		schema    types.M
	}
	tests := []struct {
		name       string
		args       args
		wantErr    error
		initialize func(string)
		clean      func(string)
	}{
		{
			name: "1",
			args: args{
				className: "post",
				schema:    nil,
			},
			wantErr:    nil,
			initialize: func(name string) {},
			clean:      clean,
		},
		{
			name: "1.1",
			args: args{
				className: "_User",
				schema:    nil,
			},
			wantErr:    nil,
			initialize: func(name string) {},
			clean:      clean,
		},
		{
			name: "2",
			args: args{
				className: "post",
				schema: types.M{
					"fields": types.M{
						"title": types.M{"type": "String"},
					},
				},
			},
			wantErr:    nil,
			initialize: func(name string) {},
			clean:      clean,
		},
		{
			name: "3",
			args: args{
				className: "post",
				schema: types.M{
					"fields": types.M{
						"title":    types.M{"type": "String"},
						"objectId": types.M{"type": "String"},
					},
				},
			},
			wantErr:    nil,
			initialize: func(name string) {},
			clean:      clean,
		},
		{
			name: "4",
			args: args{
				className: "post",
				schema: types.M{
					"fields": types.M{
						"title":    types.M{"type": "String"},
						"objectId": types.M{"type": "String"},
						"user":     types.M{"type": "Relation"},
					},
				},
			},
			wantErr:    nil,
			initialize: func(name string) {},
			clean: func(name string) {
				db.Exec(`DROP TABLE "` + name + `"`)
				db.Exec(`DROP TABLE "_SCHEMA"`)
				db.Exec(`DROP TABLE "_Join:user:post"`)
			},
		},
		{
			name: "5",
			args: args{
				className: "_User",
				schema: types.M{
					"fields": types.M{
						"objectId": types.M{"type": "String"},
						"name":     types.M{"type": "String"},
					},
				},
			},
			wantErr:    nil,
			initialize: func(name string) {},
			clean:      clean,
		},
		{
			name: "6",
			args: args{
				className: "TypeClass",
				schema: types.M{
					"fields": types.M{
						"StringKey":      types.M{"type": "String"},
						"DateKey":        types.M{"type": "Date"},
						"ObjectKey":      types.M{"type": "Object"},
						"FileKey":        types.M{"type": "File"},
						"BooleanKey":     types.M{"type": "Boolean"},
						"PointerKey":     types.M{"type": "Pointer"},
						"NumberKey":      types.M{"type": "Number"},
						"GeoPointKey":    types.M{"type": "GeoPoint"},
						"ArrayKey":       types.M{"type": "Array"},
						"StringArrayKey": types.M{"type": "Array", "contents": types.M{"type": "String"}},
					},
				},
			},
			wantErr:    nil,
			initialize: func(name string) {},
			clean:      clean,
		},
	}
	for _, tt := range tests {
		tt.initialize(tt.args.className)
		if err := p.createTable(tt.args.className, tt.args.schema, nil); reflect.DeepEqual(err, tt.wantErr) == false {
			t.Errorf("%q. PostgresAdapter.createTable() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
		tt.clean(tt.args.className)
	}
}

func TestPostgresAdapter_CreateClass(t *testing.T) {
	db := openDB()
	p := NewPostgresAdapter("", db)
	clean := func(name string) {
		db.Exec(`DROP TABLE "` + name + `"`)
		db.Exec(`DROP TABLE "_SCHEMA"`)
	}
	type args struct {
		className string
		schema    types.M
	}
	tests := []struct {
		name       string
		args       args
		want       types.M
		wantErr    error
		initialize func(string, types.M)
		clean      func(string)
	}{
		{
			name: "1",
			args: args{
				className: "post",
				schema:    types.M{"className": "post"},
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
			wantErr: nil,
			clean:   clean,
		},
		{
			name: "2",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"title": types.M{"type": "String"},
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
			wantErr: nil,
			clean:   clean,
		},
		{
			name: "3",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"title": types.M{"type": "String"},
					},
				},
			},
			want:    nil,
			wantErr: errs.E(errs.DuplicateValue, "Class post already exists."),
			initialize: func(name string, schema types.M) {
				p.CreateClass(name, schema)
			},
			clean: clean,
		},
	}
	for _, tt := range tests {
		if tt.initialize != nil {
			tt.initialize(tt.args.className, tt.args.schema)
		}
		got, err := p.CreateClass(tt.args.className, tt.args.schema)
		if reflect.DeepEqual(err, tt.wantErr) == false {
			t.Errorf("%q. PostgresAdapter.CreateClass() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			tt.clean(tt.args.className)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. PostgresAdapter.CreateClass() = %v, want %v", tt.name, got, tt.want)
		}
		tt.clean(tt.args.className)
	}
}

func TestPostgresAdapter_PerformInitialization(t *testing.T) {
	db := openDB()
	p := NewPostgresAdapter("", db)
	type args struct {
		options types.M
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		clean   func()
	}{
		{
			name:    "1",
			args:    args{options: nil},
			wantErr: nil,
			clean:   func() {},
		},
		{
			name: "2",
			args: args{
				options: types.M{
					"VolatileClassesSchemas": []types.M{
						types.M{
							"className": "_Hooks",
							"fields": types.M{
								"functionName": types.M{"type": "String"},
								"className":    types.M{"type": "String"},
								"triggerName":  types.M{"type": "String"},
								"url":          types.M{"type": "String"},
							},
							"classLevelPermissions": types.M{},
						},
					},
				},
			},
			wantErr: nil,
			clean: func() {
				db.Exec(`DROP TABLE "_SCHEMA"`)
				db.Exec(`DROP TABLE "_Hooks"`)
			},
		},
	}
	for _, tt := range tests {
		if err := p.PerformInitialization(tt.args.options); reflect.DeepEqual(err, tt.wantErr) == false {
			t.Errorf("%q. PostgresAdapter.PerformInitialization() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
		tt.clean()
	}
}

func TestPostgresAdapter_SetClassLevelPermissions(t *testing.T) {
	db := openDB()
	p := NewPostgresAdapter("", db)
	type args struct {
		className string
		CLPs      types.M
	}
	tests := []struct {
		name       string
		args       args
		wantErr    error
		initialize func()
		clean      func()
	}{
		{
			name: "1",
			args: args{
				className: "post",
				CLPs:      types.M{"find": types.M{"*": true}},
			},
			wantErr: nil,
			initialize: func() {
				p.CreateClass("post", types.M{"className": "post", "fields": types.M{"title": types.M{"type": "String"}}})
			},
			clean: func() {
				db.Exec(`DROP TABLE "_SCHEMA"`)
				db.Exec(`DROP TABLE "post"`)
			},
		},
		{
			name: "2",
			args: args{
				className: "user",
				CLPs:      types.M{"find": types.M{"*": true}},
			},
			wantErr:    nil,
			initialize: func() {},
			clean: func() {
				db.Exec(`DROP TABLE "_SCHEMA"`)
			},
		},
	}
	p.PerformInitialization(nil)
	for _, tt := range tests {
		tt.initialize()
		if err := p.SetClassLevelPermissions(tt.args.className, tt.args.CLPs); reflect.DeepEqual(err, tt.wantErr) == false {
			t.Errorf("%q. PostgresAdapter.SetClassLevelPermissions() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
		tt.clean()
	}
}

func TestPostgresAdapter_AddFieldIfNotExists(t *testing.T) {
	db := openDB()
	p := NewPostgresAdapter("", db)
	type args struct {
		className string
		fieldName string
		fieldType types.M
	}
	tests := []struct {
		name       string
		args       args
		wantErr    error
		initialize func()
		clean      func()
	}{
		{
			name: "1",
			args: args{
				className: "post",
				fieldName: "title",
				fieldType: types.M{"type": "String"},
			},
			wantErr:    nil,
			initialize: func() {},
			clean: func() {
				db.Exec(`DROP TABLE "post"`)
				db.Exec(`DROP TABLE "_SCHEMA"`)
			},
		},
		{
			name: "2",
			args: args{
				className: "post",
				fieldName: "title",
				fieldType: types.M{"type": "String"},
			},
			wantErr: nil,
			initialize: func() {
				p.CreateClass("post", types.M{"fields": types.M{"title": types.M{"type": "String"}}})
			},
			clean: func() {
				db.Exec(`DROP TABLE "post"`)
				db.Exec(`DROP TABLE "_SCHEMA"`)
			},
		},
		{
			name: "3",
			args: args{
				className: "post",
				fieldName: "title",
				fieldType: types.M{"type": "String"},
			},
			wantErr: nil,
			initialize: func() {
				p.CreateClass("post", types.M{"fields": types.M{"id": types.M{"type": "String"}}})
			},
			clean: func() {
				db.Exec(`DROP TABLE "post"`)
				db.Exec(`DROP TABLE "_SCHEMA"`)
			},
		},
		{
			name: "4",
			args: args{
				className: "post",
				fieldName: "user",
				fieldType: types.M{"type": "Relation"},
			},
			wantErr: nil,
			initialize: func() {
				p.CreateClass("post", types.M{"fields": types.M{"id": types.M{"type": "String"}}})
			},
			clean: func() {
				db.Exec(`DROP TABLE "post"`)
				db.Exec(`DROP TABLE "_SCHEMA"`)
				db.Exec(`DROP TABLE "_Join:user:post"`)
			},
		},
	}
	for _, tt := range tests {
		tt.initialize()
		if err := p.AddFieldIfNotExists(tt.args.className, tt.args.fieldName, tt.args.fieldType); reflect.DeepEqual(err, tt.wantErr) == false {
			t.Errorf("%q. PostgresAdapter.AddFieldIfNotExists() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
		tt.clean()
	}
}

func TestPostgresAdapter_DeleteClass(t *testing.T) {
	db := openDB()
	p := NewPostgresAdapter("", db)
	type args struct {
		className string
	}
	tests := []struct {
		name       string
		args       args
		want       types.M
		wantErr    error
		initialize func()
		clean      func()
	}{
		{
			name:    "1",
			args:    args{className: "post"},
			want:    types.M{},
			wantErr: nil,
			initialize: func() {
				p.ensureSchemaCollectionExists()
			},
			clean: func() {
				db.Exec(`DROP TABLE "_SCHEMA"`)
			},
		},
		{
			name:    "2",
			args:    args{className: "post"},
			want:    types.M{},
			wantErr: nil,
			initialize: func() {
				p.CreateClass("post", types.M{})
			},
			clean: func() {
				db.Exec(`DROP TABLE "_SCHEMA"`)
			},
		},
	}
	for _, tt := range tests {
		tt.initialize()
		got, err := p.DeleteClass(tt.args.className)
		if reflect.DeepEqual(err, tt.wantErr) == false {
			t.Errorf("%q. PostgresAdapter.DeleteClass() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. PostgresAdapter.DeleteClass() = %v, want %v", tt.name, got, tt.want)
		}
		tt.clean()
	}
}

func TestPostgresAdapter_DeleteAllClasses(t *testing.T) {
	db := openDB()
	p := NewPostgresAdapter("", db)
	tests := []struct {
		name       string
		wantErr    error
		initialize func()
		clean      func()
	}{
		{
			name:       "1",
			wantErr:    nil,
			initialize: func() {},
			clean:      func() {},
		},
		{
			name:    "2",
			wantErr: nil,
			initialize: func() {
				p.ensureSchemaCollectionExists()
			},
			clean: func() {},
		},
		{
			name:    "3",
			wantErr: nil,
			initialize: func() {
				p.CreateClass("post", types.M{})
			},
			clean: func() {},
		},
		{
			name:    "4",
			wantErr: nil,
			initialize: func() {
				p.CreateClass("post", types.M{"className": "post", "fields": types.M{"title": types.M{"type": "String"}}})
				p.CreateClass("user", types.M{"className": "user", "fields": types.M{"role": types.M{"type": "Relation"}}})
			},
			clean: func() {},
		},
	}
	for _, tt := range tests {
		tt.initialize()
		if err := p.DeleteAllClasses(); reflect.DeepEqual(err, tt.wantErr) == false {
			t.Errorf("%q. PostgresAdapter.DeleteAllClasses() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
		tt.clean()
	}
}

func TestPostgresAdapter_DeleteFields(t *testing.T) {
	db := openDB()
	p := NewPostgresAdapter("", db)
	initialize := func(className string, schema types.M) {
		p.CreateClass(className, schema)
	}
	clean := func() {
		db.Exec(`DROP TABLE "post"`)
		db.Exec(`DROP TABLE "_SCHEMA"`)
	}
	type args struct {
		className  string
		schema     types.M
		fieldNames []string
	}
	tests := []struct {
		name       string
		args       args
		wantErr    error
		initialize func(className string, schema types.M)
		clean      func()
	}{
		{
			name: "1",
			args: args{
				className:  "post",
				schema:     types.M{},
				fieldNames: []string{},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "2",
			args: args{
				className: "post",
				schema: types.M{
					"fields": types.M{
						"title":   types.M{"type": "String"},
						"content": types.M{"type": "String"},
					},
				},
				fieldNames: []string{"content"},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "3",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"title":   types.M{"type": "String"},
						"content": types.M{"type": "String"},
						"user":    types.M{"type": "Relation"},
					},
				},
				fieldNames: []string{"content", "user"},
			},
			wantErr:    nil,
			initialize: initialize,
			clean: func() {
				db.Exec(`DROP TABLE "post"`)
				db.Exec(`DROP TABLE "_SCHEMA"`)
				db.Exec(`DROP TABLE "_Join:user:post"`)
			},
		},
		{
			name: "4",
			args: args{
				className: "post",
				schema: types.M{
					"fields": types.M{
						"title":   types.M{"type": "String"},
						"content": types.M{"type": "String"},
						"name":    types.M{"type": "String"},
					},
				},
				fieldNames: []string{"content", "name"},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
	}
	for _, tt := range tests {
		tt.initialize(tt.args.className, tt.args.schema)
		if err := p.DeleteFields(tt.args.className, tt.args.schema, tt.args.fieldNames); reflect.DeepEqual(err, tt.wantErr) == false {
			t.Errorf("%q. PostgresAdapter.DeleteFields() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
		tt.clean()
	}
}

func TestPostgresAdapter_GetAllClasses(t *testing.T) {
	db := openDB()
	p := NewPostgresAdapter("", db)
	tests := []struct {
		name       string
		want       []types.M
		wantErr    error
		initialize func()
		clean      func()
	}{
		{
			name:       "1",
			want:       []types.M{},
			wantErr:    nil,
			initialize: func() {},
			clean: func() {
				db.Exec(`DROP TABLE "_SCHEMA"`)
			},
		},
		{
			name: "2",
			want: []types.M{
				types.M{
					"className": "post",
					"fields":    types.M{"title": map[string]interface{}{"type": "String"}},
					"classLevelPermissions": types.M{
						"find":     types.M{"*": true},
						"get":      types.M{"*": true},
						"create":   types.M{"*": true},
						"update":   types.M{"*": true},
						"delete":   types.M{"*": true},
						"addField": types.M{"*": true},
					},
				},
				types.M{
					"className": "user",
					"fields":    types.M{"name": map[string]interface{}{"type": "String"}},
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
			wantErr: nil,
			initialize: func() {
				p.CreateClass("post", types.M{"className": "post", "fields": types.M{"title": types.M{"type": "String"}}})
				p.CreateClass("user", types.M{"className": "user", "fields": types.M{"name": types.M{"type": "String"}}})
			},
			clean: func() {
				db.Exec(`DROP TABLE "_SCHEMA"`)
				db.Exec(`DROP TABLE "post"`)
				db.Exec(`DROP TABLE "user"`)
			},
		},
	}
	for _, tt := range tests {
		tt.initialize()
		got, err := p.GetAllClasses()
		if reflect.DeepEqual(err, tt.wantErr) == false {
			t.Errorf("%q. PostgresAdapter.GetAllClasses() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. PostgresAdapter.GetAllClasses() = %v, want %v", tt.name, got, tt.want)
		}
		tt.clean()
	}
}

func TestPostgresAdapter_GetClass(t *testing.T) {
	db := openDB()
	p := NewPostgresAdapter("", db)
	type args struct {
		className string
	}
	tests := []struct {
		name       string
		args       args
		want       types.M
		wantErr    error
		initialize func()
		clean      func()
	}{
		{
			name:    "1",
			args:    args{className: "post"},
			want:    types.M{},
			wantErr: nil,
			initialize: func() {
				p.ensureSchemaCollectionExists()
			},
			clean: func() {
				db.Exec(`DROP TABLE "_SCHEMA"`)
			},
		},
		{
			name: "2",
			args: args{className: "post"},
			want: types.M{
				"className": "post",
				"fields":    types.M{"title": map[string]interface{}{"type": "String"}},
				"classLevelPermissions": types.M{
					"find":     types.M{"*": true},
					"get":      types.M{"*": true},
					"create":   types.M{"*": true},
					"update":   types.M{"*": true},
					"delete":   types.M{"*": true},
					"addField": types.M{"*": true},
				},
			},
			wantErr: nil,
			initialize: func() {
				p.CreateClass("post", types.M{"className": "post", "fields": types.M{"title": types.M{"type": "String"}}})
			},
			clean: func() {
				db.Exec(`DROP TABLE "_SCHEMA"`)
				db.Exec(`DROP TABLE "post"`)
			},
		},
	}
	for _, tt := range tests {
		tt.initialize()
		got, err := p.GetClass(tt.args.className)
		if reflect.DeepEqual(err, tt.wantErr) == false {
			t.Errorf("%q. PostgresAdapter.GetClass() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. PostgresAdapter.GetClass() = %v, want %v", tt.name, got, tt.want)
		}
		tt.clean()
	}
}

func TestPostgresAdapter_CreateObject(t *testing.T) {
	db := openDB()
	p := NewPostgresAdapter("", db)
	clean := func(className string) {
		db.Exec(`DROP TABLE "` + className + `"`)
		db.Exec(`DROP TABLE "_SCHEMA"`)
	}
	type args struct {
		className string
		schema    types.M
		object    types.M
	}
	tests := []struct {
		name       string
		args       args
		wantErr    error
		initialize func(className string, schema types.M)
		clean      func(className string)
	}{
		{
			name: "1",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key1": types.M{"type": "String"},
					},
				},
				object: types.M{
					"key1": "hello",
				},
			},
			wantErr: nil,
			initialize: func(className string, schema types.M) {
				p.CreateClass(className, schema)
			},
			clean: clean,
		},
		{
			name: "2",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"authData": types.M{"type": "Object"},
					},
				},
				object: types.M{
					"_auth_data_facebook": types.M{"id": "1024"},
				},
			},
			wantErr: nil,
			initialize: func(className string, schema types.M) {
				p.CreateClass(className, schema)
			},
			clean: clean,
		},
		{
			name: "3",
			args: args{
				className: "_User",
				schema: types.M{
					"className": "_User",
					"fields":    types.M{},
				},
				object: types.M{
					"_email_verify_token": "abc",
					"_failed_login_count": 10,
					"_perishable_token":   "abc",
				},
			},
			wantErr: nil,
			initialize: func(className string, schema types.M) {
				p.CreateClass(className, schema)
			},
			clean: clean,
		},
		{
			name: "4",
			args: args{
				className: "_User",
				schema: types.M{
					"className": "_User",
					"fields":    types.M{},
				},
				object: types.M{
					"_password_history": types.S{"hello", "world"},
				},
			},
			wantErr: nil,
			initialize: func(className string, schema types.M) {
				p.CreateClass(className, schema)
			},
			clean: clean,
		},
		{
			name: "5",
			args: args{
				className: "_User",
				schema: types.M{
					"className": "_User",
					"fields":    types.M{},
				},
				object: types.M{
					"_email_verify_token_expires_at": types.M{"type": "Date", "iso": "2017-03-01T09:10:10.000Z"},
					"_account_lockout_expires_at":    types.M{"type": "Date", "iso": "2017-03-01T09:10:10.000Z"},
					"_perishable_token_expires_at":   types.M{"type": "Date", "iso": "2017-03-01T09:10:10.000Z"},
					"_password_changed_at":           types.M{"type": "Date", "iso": "2017-03-01T09:10:10.000Z"},
				},
			},
			wantErr: nil,
			initialize: func(className string, schema types.M) {
				p.CreateClass(className, schema)
			},
			clean: clean,
		},
		{
			name: "6",
			args: args{
				className: "_User",
				schema: types.M{
					"className": "_User",
					"fields":    types.M{},
				},
				object: types.M{
					"_email_verify_token_expires_at": types.M{"type": "Date", "iso": ""},
					"_account_lockout_expires_at":    types.M{"type": "Date", "iso": ""},
					"_perishable_token_expires_at":   types.M{"type": "Date", "iso": ""},
					"_password_changed_at":           types.M{"type": "Date", "iso": ""},
				},
			},
			wantErr: nil,
			initialize: func(className string, schema types.M) {
				p.CreateClass(className, schema)
			},
			clean: clean,
		},
		{
			name: "7",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key1": types.M{"type": "Date"},
						"key2": types.M{"type": "Pointer", "targetClass": "user"},
					},
				},
				object: types.M{
					"key1": types.M{"type": "Date", "iso": "2017-03-01T09:10:10.000Z"},
					"key2": types.M{"type": "Pointer", "className": "user", "objectId": "1024"},
				},
			},
			wantErr: nil,
			initialize: func(className string, schema types.M) {
				p.CreateClass(className, schema)
			},
			clean: clean,
		},
		{
			name: "8",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key1": types.M{"type": "Array"},
						"key2": types.M{"type": "Object"},
					},
				},
				object: types.M{
					"key1": types.S{"hello", "world"},
					"key2": types.M{"key": "hello"},
				},
			},
			wantErr: nil,
			initialize: func(className string, schema types.M) {
				p.CreateClass(className, schema)
			},
			clean: clean,
		},
		{
			name: "9",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key1": types.M{"type": "String"},
						"key2": types.M{"type": "Number"},
						"key3": types.M{"type": "Boolean"},
					},
				},
				object: types.M{
					"key1": "hello",
					"key2": 1024,
					"key3": true,
				},
			},
			wantErr: nil,
			initialize: func(className string, schema types.M) {
				p.CreateClass(className, schema)
			},
			clean: clean,
		},
		{
			name: "10",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key1": types.M{"type": "File"},
					},
				},
				object: types.M{
					"key1": types.M{"type": "File", "name": "icon.png"},
				},
			},
			wantErr: nil,
			initialize: func(className string, schema types.M) {
				p.CreateClass(className, schema)
			},
			clean: clean,
		},
		{
			name: "11",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key1": types.M{"type": "GeoPoint"},
					},
				},
				object: types.M{
					"key1": types.M{"longitude": 10.0, "latitude": 10.0},
				},
			},
			wantErr: nil,
			initialize: func(className string, schema types.M) {
				p.CreateClass(className, schema)
			},
			clean: clean,
		},
		{
			name: "12",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"_rperm": types.M{"type": "Array"},
						"_wperm": types.M{"type": "Array"},
					},
				},
				object: types.M{
					"_rperm": types.S{"*", "role:1024", "1024"},
					"_wperm": types.S{"role:1024", "1024"},
				},
			},
			wantErr: nil,
			initialize: func(className string, schema types.M) {
				p.CreateClass(className, schema)
			},
			clean: clean,
		},
		{
			name: "13",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"objectId": types.M{"type": "String"},
					},
				},
				object: types.M{
					"objectId": "1024",
				},
			},
			wantErr: errs.E(errs.DuplicateValue, "A duplicate value for a field with unique values was provided"),
			initialize: func(className string, schema types.M) {
				p.CreateClass(className, schema)
				p.CreateObject(className, schema, types.M{"objectId": "1024"})
			},
			clean: clean,
		},
		{
			name: "14",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key1": types.M{"type": "String"},
					},
				},
				object: nil,
			},
			wantErr: nil,
			initialize: func(className string, schema types.M) {
				p.CreateClass(className, schema)
			},
			clean: clean,
		},
	}
	for _, tt := range tests {
		tt.initialize(tt.args.className, tt.args.schema)
		if err := p.CreateObject(tt.args.className, tt.args.schema, tt.args.object); reflect.DeepEqual(err, tt.wantErr) == false {
			t.Errorf("%q. PostgresAdapter.CreateObject() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
		tt.clean(tt.args.className)
	}
}

func TestPostgresAdapter_Find(t *testing.T) {
	db := openDB()
	p := NewPostgresAdapter("", db)
	initialize := func(className string, schema types.M, objects []types.M) {
		p.CreateClass(className, schema)
		for _, object := range objects {
			p.CreateObject(className, schema, object)
		}
	}
	clean := func(className string) {
		db.Exec(`DROP TABLE "` + className + `"`)
		db.Exec(`DROP TABLE "_SCHEMA"`)
	}
	type args struct {
		className   string
		schema      types.M
		query       types.M
		options     types.M
		dataObjects []types.M
	}
	tests := []struct {
		name       string
		args       args
		want       []types.M
		wantErr    error
		initialize func(className string, schema types.M, objects []types.M)
		clean      func(className string)
	}{
		{
			name: "1",
			args: args{
				className: "post",
				schema:    types.M{},
				query:     types.M{"key": "hello"},
				options:   types.M{},
			},
			want:    []types.M{},
			wantErr: nil,
			initialize: func(className string, schema types.M, objects []types.M) {
				p.ensureSchemaCollectionExists()
			},
			clean: func(className string) {
				db.Exec(`DROP TABLE "_SCHEMA"`)
			},
		},
		{
			name: "2",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields":    types.M{"key": types.M{"type": "String"}},
				},
				query:       types.M{"key": "hello"},
				options:     types.M{},
				dataObjects: []types.M{},
			},
			want:       []types.M{},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "3",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields":    types.M{"key": types.M{"type": "String"}},
				},
				query:   types.M{"key": "hello"},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": "hello"},
					types.M{"key": "world"},
				},
			},
			want: []types.M{
				types.M{"key": "hello"},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "4",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields":    types.M{"key": types.M{"type": "String"}},
				},
				query:   types.M{},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": "hello"},
					types.M{"key": "world"},
				},
			},
			want: []types.M{
				types.M{"key": "hello"},
				types.M{"key": "world"},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "5-limit",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields":    types.M{"key": types.M{"type": "String"}},
				},
				query:   types.M{},
				options: types.M{"limit": 1},
				dataObjects: []types.M{
					types.M{"key": "hello"},
					types.M{"key": "world"},
				},
			},
			want: []types.M{
				types.M{"key": "hello"},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "6-skip",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields":    types.M{"key": types.M{"type": "String"}},
				},
				query:   types.M{},
				options: types.M{"skip": 1},
				dataObjects: []types.M{
					types.M{"key": "hello"},
					types.M{"key": "world"},
				},
			},
			want: []types.M{
				types.M{"key": "world"},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "7-sort",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields":    types.M{"key": types.M{"type": "String"}},
				},
				query:   types.M{},
				options: types.M{"sort": []string{"key"}},
				dataObjects: []types.M{
					types.M{"key": "hello"},
					types.M{"key": "world"},
				},
			},
			want: []types.M{
				types.M{"key": "hello"},
				types.M{"key": "world"},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "8-sort",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields":    types.M{"key": types.M{"type": "String"}},
				},
				query:   types.M{},
				options: types.M{"sort": []string{"-key"}},
				dataObjects: []types.M{
					types.M{"key": "hello"},
					types.M{"key": "world"},
				},
			},
			want: []types.M{
				types.M{"key": "world"},
				types.M{"key": "hello"},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "9-keys",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "String"},
					},
				},
				query:   types.M{},
				options: types.M{"keys": []string{"key"}},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": "world"},
					types.M{"key": "hi", "key2": "friend"},
				},
			},
			want: []types.M{
				types.M{"key": "hello"},
				types.M{"key": "hi"},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "10-keys",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "String"},
					},
				},
				query:   types.M{},
				options: types.M{"keys": []string{"key2"}},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": "world"},
					types.M{"key": "hi", "key2": "friend"},
				},
			},
			want: []types.M{
				types.M{"key2": "world"},
				types.M{"key2": "friend"},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "11-Pointer",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{
							"type":        "Pointer",
							"targetClass": "user",
						},
					},
				},
				query:   types.M{},
				options: types.M{},
				dataObjects: []types.M{
					types.M{
						"key": types.M{
							"className": "user",
							"__type":    "Pointer",
							"objectId":  "123456789012345678901111",
						},
					},
				},
			},
			want: []types.M{
				types.M{
					"key": types.M{
						"className": "user",
						"__type":    "Pointer",
						"objectId":  "123456789012345678901111",
					},
				},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "12-Relation",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{
							"type":        "Relation",
							"targetClass": "user",
						},
						"key2": types.M{"type": "String"},
					},
				},
				query:       types.M{},
				options:     types.M{},
				dataObjects: []types.M{types.M{"key2": "hello"}},
			},
			want: []types.M{
				types.M{
					"key": types.M{
						"className": "user",
						"__type":    "Relation",
					},
					"key2": "hello",
				},
			},
			wantErr:    nil,
			initialize: initialize,
			clean: func(className string) {
				db.Exec(`DROP TABLE "` + className + `"`)
				db.Exec(`DROP TABLE "_SCHEMA"`)
				db.Exec(`DROP TABLE "_Join:key:post"`)
			},
		},
		{
			name: "13-GeoPoint",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{
							"type": "GeoPoint",
						},
					},
				},
				query:   types.M{},
				options: types.M{},
				dataObjects: []types.M{
					types.M{
						"key": types.M{
							"__type":    "GeoPoint",
							"longitude": -30.5,
							"latitude":  40.5,
						},
					},
				},
			},
			want: []types.M{
				types.M{
					"key": types.M{
						"__type":    "GeoPoint",
						"longitude": -30.5,
						"latitude":  40.5,
					},
				},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "14-File",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{
							"type": "File",
						},
					},
				},
				query:   types.M{},
				options: types.M{},
				dataObjects: []types.M{
					types.M{
						"key": types.M{
							"__type": "File",
							"name":   "icon.png",
						},
					},
				},
			},
			want: []types.M{
				types.M{
					"key": types.M{
						"__type": "File",
						"name":   "icon.png",
					},
				},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "15",
			args: args{
				className: "post-Object",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{
							"type": "Object",
						},
					},
				},
				query:   types.M{},
				options: types.M{},
				dataObjects: []types.M{
					types.M{
						"key": types.M{
							"sub": "hello",
						},
					},
				},
			},
			want: []types.M{
				types.M{
					"key": types.M{
						"sub": "hello",
					},
				},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "16-Array",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{
							"type": "Array",
						},
					},
				},
				query:   types.M{},
				options: types.M{},
				dataObjects: []types.M{
					types.M{
						"key": types.S{"hello", "world"},
					},
				},
			},
			want: []types.M{
				types.M{
					"key": types.S{"hello", "world"},
				},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "17-Boolean Number",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "Boolean"},
						"key2": types.M{"type": "Number"},
					},
				},
				query:   types.M{},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": true, "key2": 10.24},
				},
			},
			want: []types.M{
				types.M{"key": true, "key2": 10.24},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "18-_rperm",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"_rperm": types.M{
							"type": "Array",
						},
					},
				},
				query:   types.M{},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"_rperm": types.S{"hello", "world"}},
				},
			},
			want: []types.M{
				types.M{"_rperm": types.S{"hello", "world"}},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "19-_wperm",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"_wperm": types.M{
							"type": "Array",
						},
					},
				},
				query:   types.M{},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"_wperm": types.S{"hello", "world"}},
				},
			},
			want: []types.M{
				types.M{"_wperm": types.S{"hello", "world"}},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "20-objectId...",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"objectId":  types.M{"type": "String"},
						"createdAt": types.M{"type": "Date"},
						"updatedAt": types.M{"type": "Date"},
					},
				},
				query:   types.M{},
				options: types.M{},
				dataObjects: []types.M{
					types.M{
						"objectId": "123456789012345678901111",
						"createdAt": types.M{
							"__type": "Date",
							"iso":    "2006-01-02T15:04:05.000Z",
						},
						"updatedAt": types.M{
							"__type": "Date",
							"iso":    "2006-01-02T15:04:05.000Z",
						},
					},
				},
			},
			want: []types.M{
				types.M{
					"objectId":  "123456789012345678901111",
					"createdAt": "2006-01-02T15:04:05.000Z",
					"updatedAt": "2006-01-02T15:04:05.000Z",
				},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "21-expiresAt...",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"expiresAt":                      types.M{"type": "Date"},
						"_email_verify_token_expires_at": types.M{"type": "Date"},
						"_account_lockout_expires_at":    types.M{"type": "Date"},
						"_perishable_token_expires_at":   types.M{"type": "Date"},
						"_password_changed_at":           types.M{"type": "Date"},
					},
				},
				query:   types.M{},
				options: types.M{},
				dataObjects: []types.M{
					types.M{
						"expiresAt": types.M{
							"__type": "Date",
							"iso":    "2006-01-02T15:04:05.000Z",
						},
						"_email_verify_token_expires_at": types.M{
							"__type": "Date",
							"iso":    "2006-01-02T15:04:05.000Z",
						},
						"_account_lockout_expires_at": types.M{
							"__type": "Date",
							"iso":    "2006-01-02T15:04:05.000Z",
						},
						"_perishable_token_expires_at": types.M{
							"__type": "Date",
							"iso":    "2006-01-02T15:04:05.000Z",
						},
						"_password_changed_at": types.M{
							"__type": "Date",
							"iso":    "2006-01-02T15:04:05.000Z",
						},
					},
				},
			},
			want: []types.M{
				types.M{
					"expiresAt": types.M{
						"__type": "Date",
						"iso":    "2006-01-02T15:04:05.000Z",
					},
					"_email_verify_token_expires_at": types.M{
						"__type": "Date",
						"iso":    "2006-01-02T15:04:05.000Z",
					},
					"_account_lockout_expires_at": types.M{
						"__type": "Date",
						"iso":    "2006-01-02T15:04:05.000Z",
					},
					"_perishable_token_expires_at": types.M{
						"__type": "Date",
						"iso":    "2006-01-02T15:04:05.000Z",
					},
					"_password_changed_at": types.M{
						"__type": "Date",
						"iso":    "2006-01-02T15:04:05.000Z",
					},
				},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "22-Date",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "Date"},
					},
				},
				query:   types.M{},
				options: types.M{},
				dataObjects: []types.M{
					types.M{
						"key": types.M{
							"__type": "Date",
							"iso":    "2006-01-02T15:04:05.000Z",
						},
					},
				},
			},
			want: []types.M{
				types.M{
					"key": types.M{
						"__type": "Date",
						"iso":    "2006-01-02T15:04:05.000Z",
					},
				},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "23-where-.",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "Object"},
					},
				},
				query:   types.M{"key.subKey": "hello"},
				options: types.M{},
				dataObjects: []types.M{
					types.M{
						"key": types.M{
							"subKey": "hello",
							"sub":    types.M{"key": "world"},
						},
					},
				},
			},
			want: []types.M{
				types.M{
					"key": types.M{
						"subKey": "hello",
						"sub":    map[string]interface{}{"key": "world"},
					},
				},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "24-where-.",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "Object"},
					},
				},
				query:   types.M{"key.sub": types.M{"key": "world"}},
				options: types.M{},
				dataObjects: []types.M{
					types.M{
						"key": types.M{
							"subKey": "hello",
							"sub":    types.M{"key": "world"},
						},
					},
				},
			},
			want: []types.M{
				types.M{
					"key": types.M{
						"subKey": "hello",
						"sub":    map[string]interface{}{"key": "world"},
					},
				},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "25-where-.",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "Object"},
					},
				},
				query:   types.M{"key.sub.key": "world"},
				options: types.M{},
				dataObjects: []types.M{
					types.M{
						"key": types.M{
							"subKey": "hello",
							"sub":    types.M{"key": "world"},
						},
					},
				},
			},
			want: []types.M{
				types.M{
					"key": types.M{
						"subKey": "hello",
						"sub":    map[string]interface{}{"key": "world"},
					},
				},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "26-where-bool",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "Boolean"},
					},
				},
				query:   types.M{"key": true},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": true},
					types.M{"key": false},
				},
			},
			want: []types.M{
				types.M{"key": true},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "27-where-float64",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "Number"},
					},
				},
				query:   types.M{"key": 10.24},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": 10.24},
					types.M{"key": 20.48},
				},
			},
			want: []types.M{
				types.M{"key": 10.24},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "28-where-int",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "Number"},
					},
				},
				query:   types.M{"key": 1024},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": 1024},
					types.M{"key": 2048},
				},
			},
			want: []types.M{
				types.M{"key": 1024.0},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "29-where-$or",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "Number"},
						"key2": types.M{"type": "String"},
					},
				},
				query: types.M{
					"$or": types.S{
						types.M{"key": 10.24},
						types.M{"key": 20.48},
					},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": 10.24, "key2": "hello"},
					types.M{"key": 20.48, "key2": "world"},
				},
			},
			want: []types.M{
				types.M{"key": 10.24, "key2": "hello"},
				types.M{"key": 20.48, "key2": "world"},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "30-where-$and",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "Number"},
						"key2": types.M{"type": "String"},
					},
				},
				query: types.M{
					"$and": types.S{
						types.M{"key": 10.24},
						types.M{"key2": "hello"},
					},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": 10.24, "key2": "hello"},
					types.M{"key": 20.48, "key2": "world"},
				},
			},
			want: []types.M{
				types.M{"key": 10.24, "key2": "hello"},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "31-where-$ne",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "String"},
					},
				},
				query: types.M{
					"key": types.M{"$ne": nil},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": "hello"},
					types.M{"key": nil},
				},
			},
			want: []types.M{
				types.M{"key": "hello"},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "32-where-$ne",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "Number"},
					},
				},
				query: types.M{
					"key": types.M{"$ne": "world"},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": 10.24},
					types.M{"key": nil, "key2": 10.24},
				},
			},
			want: []types.M{
				types.M{"key": "hello", "key2": 10.24},
				types.M{"key2": 10.24},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "33-where-$ne",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "Array"},
					},
				},
				query: types.M{
					"key": types.M{"$ne": "world"},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": types.S{"hello", "world"}},
					types.M{"key": types.S{"hello", nil}},
				},
			},
			want: []types.M{
				types.M{"key": types.S{"hello", nil}},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "33-where-$ne",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "Array"},
					},
				},
				query: types.M{
					"key": types.M{"$ne": nil},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": types.S{"hello", "world"}},
					types.M{"key": types.S{"hello", nil}},
				},
			},
			want: []types.M{
				types.M{"key": types.S{"hello", "world"}},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "34-where-$eq",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "String"},
					},
				},
				query: types.M{
					"key2": types.M{"$eq": "world"},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": "world"},
					types.M{"key": "hello", "key2": nil},
				},
			},
			want: []types.M{
				types.M{"key": "hello", "key2": "world"},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "35-where-$eq",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "String"},
					},
				},
				query: types.M{
					"key2": types.M{"$eq": nil},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": "world"},
					types.M{"key": "hello", "key2": nil},
				},
			},
			want: []types.M{
				types.M{"key": "hello"},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "36-where-text[]",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"_rperm": types.M{"type": "Array"},
					},
				},
				query: types.M{
					"_rperm": types.M{"$in": types.S{"hello"}},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"_rperm": types.S{"hello", "world"}},
				},
			},
			want: []types.M{
				types.M{"_rperm": types.S{"hello", "world"}},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "37-where-text[]",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"_rperm": types.M{"type": "Array"},
					},
				},
				query: types.M{
					"_rperm": types.M{"$in": types.S{"hello", nil}},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"_rperm": types.S{"hello", "world"}},
				},
			},
			want: []types.M{
				types.M{"_rperm": types.S{"hello", "world"}},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "37.1-where-text[]",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":    types.M{"type": "String"},
						"_rperm": types.M{"type": "Array"},
					},
				},
				query: types.M{
					"_rperm": types.M{"$in": types.S{"hello", nil}},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": "hi", "_rperm": types.S{"hello", "world"}},
					types.M{"key": "hi"},
				},
			},
			want: []types.M{
				types.M{"key": "hi", "_rperm": types.S{"hello", "world"}},
				types.M{"key": "hi"},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "38-where-$in",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "Array"},
					},
				},
				query: types.M{
					"key": types.M{"$in": types.S{"hello"}},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": types.S{"hello", "world"}},
				},
			},
			want: []types.M{
				types.M{"key": types.S{"hello", "world"}},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "39-where-$in",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "String"},
					},
				},
				query: types.M{
					"key": types.M{"$in": types.S{"hello"}},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": "hello"},
					types.M{"key": "world"},
				},
			},
			want: []types.M{
				types.M{"key": "hello"},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "40-where-$nin",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "Array"},
					},
				},
				query: types.M{
					"key": types.M{"$nin": types.S{"hello"}},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": types.S{"hello", "world"}},
					types.M{"key": types.S{"hi", "world"}},
				},
			},
			want: []types.M{
				types.M{"key": types.S{"hi", "world"}},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "41-where-$nin",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "String"},
					},
				},
				query: types.M{
					"key": types.M{"$nin": types.S{"hello"}},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": "hello"},
					types.M{"key": "world"},
				},
			},
			want: []types.M{
				types.M{"key": "world"},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "42-where-$all",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "Array"},
					},
				},
				query: types.M{
					"key": types.M{"$all": types.S{"hello", "world"}},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": types.S{"hello", "world"}},
					types.M{"key": types.S{"hi", "world"}},
				},
			},
			want: []types.M{
				types.M{"key": types.S{"hello", "world"}},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "43-where-$exists",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "String"},
					},
				},
				query: types.M{
					"key": types.M{"$exists": true},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": "world"},
					types.M{"key": nil, "key2": "world"},
				},
			},
			want: []types.M{
				types.M{"key": "hello", "key2": "world"},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "44-where-$exists",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "String"},
					},
				},
				query: types.M{
					"key": types.M{"$exists": false},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": "world"},
					types.M{"key": nil, "key2": "world"},
				},
			},
			want: []types.M{
				types.M{"key2": "world"},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "45-where-$nearSphere",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "GeoPoint"},
					},
				},
				query: types.M{
					"key": types.M{
						"$nearSphere": types.M{
							"__type":    "GeoPoint",
							"longitude": 10.00,
							"latitude":  10.00,
						},
						"$maxDistance": 0.1,
					},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{
						"key": types.M{
							"__type":    "GeoPoint",
							"longitude": 10.20,
							"latitude":  10.20,
						},
					},
					types.M{
						"key": types.M{
							"__type":    "GeoPoint",
							"longitude": 10.10,
							"latitude":  10.10,
						},
					},
					types.M{
						"key": types.M{
							"__type":    "GeoPoint",
							"longitude": 90.10,
							"latitude":  90.10,
						},
					},
				},
			},
			want: []types.M{
				types.M{
					"key": types.M{
						"__type":    "GeoPoint",
						"longitude": 10.10,
						"latitude":  10.10,
					},
				},
				types.M{
					"key": types.M{
						"__type":    "GeoPoint",
						"longitude": 10.20,
						"latitude":  10.20,
					},
				},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "46-where-$within",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "GeoPoint"},
					},
				},
				query: types.M{
					"key": types.M{
						"$within": types.M{
							"$box": types.S{
								types.M{
									"__type":    "GeoPoint",
									"latitude":  5.0,
									"longitude": 5.0,
								},
								types.M{
									"__type":    "GeoPoint",
									"latitude":  15.0,
									"longitude": 15.0,
								},
							},
						},
					},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{
						"key": types.M{
							"__type":    "GeoPoint",
							"longitude": 10.10,
							"latitude":  10.10,
						},
					},
					types.M{
						"key": types.M{
							"__type":    "GeoPoint",
							"longitude": 90.10,
							"latitude":  90.10,
						},
					},
				},
			},
			want: []types.M{
				types.M{
					"key": types.M{
						"__type":    "GeoPoint",
						"longitude": 10.10,
						"latitude":  10.10,
					},
				},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "47-where-$regex",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "String"},
					},
				},
				query: types.M{
					"key": types.M{
						"$regex": "^h",
					},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": "hello"},
					types.M{"key": "world"},
				},
			},
			want: []types.M{
				types.M{"key": "hello"},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "48-where-$regex",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "String"},
					},
				},
				query: types.M{
					"key": types.M{
						"$regex":   "^h",
						"$options": "i",
					},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": "hello"},
					types.M{"key": "Hello"},
				},
			},
			want: []types.M{
				types.M{"key": "hello"},
				types.M{"key": "Hello"},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "49-where-$regex",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "String"},
					},
				},
				query: types.M{
					"key": types.M{
						"$regex":   "^h e",
						"$options": "x",
					},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": "hello"},
					types.M{"key": "Hello"},
				},
			},
			want: []types.M{
				types.M{"key": "hello"},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "50-where-$Pointer",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "Array"},
					},
				},
				query: types.M{
					"key": types.M{
						"__type":    "Pointer",
						"objectId":  "123456789012345678901111",
						"className": "user",
					},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{
						"key": types.S{
							types.M{
								"__type":    "Pointer",
								"objectId":  "123456789012345678901111",
								"className": "user",
							},
						},
					},
					types.M{
						"key": types.S{
							types.M{
								"__type":    "Pointer",
								"objectId":  "123456789012345678902222",
								"className": "user",
							},
						},
					},
				},
			},
			want: []types.M{
				types.M{
					"key": types.S{
						map[string]interface{}{
							"__type":    "Pointer",
							"objectId":  "123456789012345678901111",
							"className": "user",
						},
					},
				},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "51-where-$Pointer",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "Pointer", "targetClass": "user"},
					},
				},
				query: types.M{
					"key": types.M{
						"__type":    "Pointer",
						"objectId":  "123456789012345678901111",
						"className": "user",
					},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{
						"key": types.M{
							"__type":    "Pointer",
							"objectId":  "123456789012345678901111",
							"className": "user",
						},
					},
					types.M{
						"key": types.M{
							"__type":    "Pointer",
							"objectId":  "123456789012345678902222",
							"className": "user",
						},
					},
				},
			},
			want: []types.M{
				types.M{
					"key": types.M{
						"__type":    "Pointer",
						"objectId":  "123456789012345678901111",
						"className": "user",
					},
				},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "52-where-$Date",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "Date"},
					},
				},
				query: types.M{
					"key": types.M{
						"__type": "Date",
						"iso":    "2006-01-02T15:04:05.000Z",
					},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{
						"key": types.M{
							"__type": "Date",
							"iso":    "2006-01-02T15:04:05.000Z",
						},
					},
					types.M{
						"key": types.M{
							"__type": "Date",
							"iso":    "2007-01-02T15:04:05.000Z",
						},
					},
				},
			},
			want: []types.M{
				types.M{
					"key": types.M{
						"__type": "Date",
						"iso":    "2006-01-02T15:04:05.000Z",
					},
				},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "53-where-$gt",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "Number"},
					},
				},
				query: types.M{
					"key": types.M{"$gt": 15.0},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": 10.0},
					types.M{"key": 15.0},
					types.M{"key": 20.0},
				},
			},
			want: []types.M{
				types.M{"key": 20.0},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "54-where-$lt",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "Number"},
					},
				},
				query: types.M{
					"key": types.M{"$lt": 15.0},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": 10.0},
					types.M{"key": 15.0},
					types.M{"key": 20.0},
				},
			},
			want: []types.M{
				types.M{"key": 10.0},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "55-where-$gte",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "Number"},
					},
				},
				query: types.M{
					"key": types.M{"$gte": 15.0},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": 10.0},
					types.M{"key": 15.0},
					types.M{"key": 20.0},
				},
			},
			want: []types.M{
				types.M{"key": 15.0},
				types.M{"key": 20.0},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "56-where-$lte",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "Number"},
					},
				},
				query: types.M{
					"key": types.M{"$lte": 15.0},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": 10.0},
					types.M{"key": 15.0},
					types.M{"key": 20.0},
				},
			},
			want: []types.M{
				types.M{"key": 10.0},
				types.M{"key": 15.0},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "57-where-$gt",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "String"},
					},
				},
				query: types.M{
					"key": types.M{"$gt": "def"},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": "abc"},
					types.M{"key": "def"},
					types.M{"key": "hij"},
				},
			},
			want: []types.M{
				types.M{"key": "hij"},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "58-where-$lt",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "String"},
					},
				},
				query: types.M{
					"key": types.M{"$lt": "def"},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": "abc"},
					types.M{"key": "def"},
					types.M{"key": "hij"},
				},
			},
			want: []types.M{
				types.M{"key": "abc"},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "59-where-$gt",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "Date"},
					},
				},
				query: types.M{
					"key": types.M{
						"$gt": types.M{
							"__type": "Date",
							"iso":    "2006-01-02T15:04:05.000Z",
						},
					},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{
						"key": types.M{
							"__type": "Date",
							"iso":    "2005-01-02T15:04:05.000Z",
						},
					},
					types.M{
						"key": types.M{
							"__type": "Date",
							"iso":    "2007-01-02T15:04:05.000Z",
						},
					},
				},
			},
			want: []types.M{
				types.M{
					"key": types.M{
						"__type": "Date",
						"iso":    "2007-01-02T15:04:05.000Z",
					},
				},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "60-where-$lt",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "Date"},
					},
				},
				query: types.M{
					"key": types.M{
						"$lt": types.M{
							"__type": "Date",
							"iso":    "2006-01-02T15:04:05.000Z",
						},
					},
				},
				options: types.M{},
				dataObjects: []types.M{
					types.M{
						"key": types.M{
							"__type": "Date",
							"iso":    "2005-01-02T15:04:05.000Z",
						},
					},
					types.M{
						"key": types.M{
							"__type": "Date",
							"iso":    "2007-01-02T15:04:05.000Z",
						},
					},
				},
			},
			want: []types.M{
				types.M{
					"key": types.M{
						"__type": "Date",
						"iso":    "2005-01-02T15:04:05.000Z",
					},
				},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "61-where-null",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "String"},
					},
				},
				query:   types.M{"key2": nil},
				options: types.M{},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": "world"},
					types.M{"key": "hello", "key2": nil},
				},
			},
			want: []types.M{
				types.M{"key": "hello"},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
	}
	for _, tt := range tests {
		tt.initialize(tt.args.className, tt.args.schema, tt.args.dataObjects)
		got, err := p.Find(tt.args.className, tt.args.schema, tt.args.query, tt.args.options)
		tt.clean(tt.args.className)
		if reflect.DeepEqual(err, tt.wantErr) == false {
			t.Errorf("%q. PostgresAdapter.Find() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. PostgresAdapter.Find() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestPostgresAdapter_DeleteObjectsByQuery(t *testing.T) {
	db := openDB()
	p := NewPostgresAdapter("", db)
	initialize := func(className string, schema types.M, objects []types.M) {
		p.CreateClass(className, schema)
		for _, object := range objects {
			p.CreateObject(className, schema, object)
		}
	}
	clean := func(className string) {
		db.Exec(`DROP TABLE "` + className + `"`)
		db.Exec(`DROP TABLE "_SCHEMA"`)
	}
	type args struct {
		className   string
		schema      types.M
		query       types.M
		dataObjects []types.M
	}
	tests := []struct {
		name       string
		args       args
		wantErr    error
		initialize func(className string, schema types.M, objects []types.M)
		clean      func(className string)
	}{
		{
			name: "1",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "String"},
					},
				},
				query: types.M{},
				dataObjects: []types.M{
					types.M{"key": "hello"},
					types.M{"key": "hi"},
				},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "2",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "String"},
					},
				},
				query: types.M{"key": "hi"},
				dataObjects: []types.M{
					types.M{"key": "hello"},
					types.M{"key": "hi"},
				},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "3",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "String"},
					},
				},
				query: types.M{"key": "world"},
				dataObjects: []types.M{
					types.M{"key": "hello"},
					types.M{"key": "hi"},
				},
			},
			wantErr:    errs.E(errs.ObjectNotFound, "Object not found."),
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "4",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "String"},
					},
				},
				query: types.M{"key": "world"},
				dataObjects: []types.M{
					types.M{"key": "hello"},
					types.M{"key": "hi"},
				},
			},
			wantErr: errs.E(errs.ObjectNotFound, "Object not found."),
			initialize: func(className string, schema types.M, objects []types.M) {
				p.ensureSchemaCollectionExists()
			},
			clean: clean,
		},
	}
	for _, tt := range tests {
		tt.initialize(tt.args.className, tt.args.schema, tt.args.dataObjects)
		err := p.DeleteObjectsByQuery(tt.args.className, tt.args.schema, tt.args.query)
		tt.clean(tt.args.className)
		if reflect.DeepEqual(err, tt.wantErr) == false {
			t.Errorf("%q. PostgresAdapter.DeleteObjectsByQuery() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestPostgresAdapter_Count(t *testing.T) {
	db := openDB()
	p := NewPostgresAdapter("", db)
	initialize := func(className string, schema types.M, objects []types.M) {
		p.CreateClass(className, schema)
		for _, object := range objects {
			p.CreateObject(className, schema, object)
		}
	}
	clean := func(className string) {
		db.Exec(`DROP TABLE "` + className + `"`)
		db.Exec(`DROP TABLE "_SCHEMA"`)
	}
	type args struct {
		className   string
		schema      types.M
		query       types.M
		dataObjects []types.M
	}
	tests := []struct {
		name       string
		args       args
		want       int
		wantErr    error
		initialize func(className string, schema types.M, objects []types.M)
		clean      func(className string)
	}{
		{
			name: "1",
			args: args{
				className:   "post",
				schema:      types.M{},
				query:       types.M{},
				dataObjects: []types.M{},
			},
			want:    0,
			wantErr: nil,
			initialize: func(className string, schema types.M, objects []types.M) {
				p.ensureSchemaCollectionExists()
			},
			clean: func(className string) {
				db.Exec(`DROP TABLE "_SCHEMA"`)
			},
		},
		{
			name: "2",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "String"},
					},
				},
				query: types.M{"key": "world"},
				dataObjects: []types.M{
					types.M{"key": "hello"},
					types.M{"key": "hi"},
				},
			},
			want:       0,
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "3",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "String"},
					},
				},
				query: types.M{"key": "hello"},
				dataObjects: []types.M{
					types.M{"key": "hello"},
					types.M{"key": "hi"},
				},
			},
			want:       1,
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "4",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key": types.M{"type": "String"},
					},
				},
				query: types.M{},
				dataObjects: []types.M{
					types.M{"key": "hello"},
					types.M{"key": "hi"},
				},
			},
			want:       2,
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
	}
	for _, tt := range tests {
		tt.initialize(tt.args.className, tt.args.schema, tt.args.dataObjects)
		got, err := p.Count(tt.args.className, tt.args.schema, tt.args.query)
		tt.clean(tt.args.className)
		if reflect.DeepEqual(err, tt.wantErr) == false {
			t.Errorf("%q. PostgresAdapter.Count() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. PostgresAdapter.Count() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestPostgresAdapter_FindOneAndUpdate(t *testing.T) {
	db := openDB()
	p := NewPostgresAdapter("", db)
	initialize := func(className string, schema types.M, objects []types.M) {
		p.CreateClass(className, schema)
		for _, object := range objects {
			p.CreateObject(className, schema, object)
		}
	}
	clean := func(className string) {
		db.Exec(`DROP TABLE "` + className + `"`)
		db.Exec(`DROP TABLE "_SCHEMA"`)
	}
	type args struct {
		className   string
		schema      types.M
		query       types.M
		update      types.M
		dataObjects []types.M
	}
	tests := []struct {
		name       string
		args       args
		want       types.M
		wantErr    error
		initialize func(className string, schema types.M, objects []types.M)
		clean      func(className string)
	}{
		{
			name: "1-NULL",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "String"},
					},
				},
				query:  types.M{"key2": nil},
				update: types.M{"key2": "world"},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": "world"},
					types.M{"key": "hi", "key2": nil},
				},
			},
			want:       types.M{"key": "hi", "key2": "world"},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "2-authData",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"authData": types.M{"type": "Object"},
					},
				},
				query:  types.M{"authData.facebook": types.M{"id": "2048"}},
				update: types.M{"_auth_data_facebook": types.M{"id": "512"}},
				dataObjects: []types.M{
					types.M{"_auth_data_facebook": types.M{"id": "1024"}},
					types.M{"_auth_data_facebook": types.M{"id": "2048"}},
				},
			},
			want:       types.M{"authData": types.M{"facebook": map[string]interface{}{"id": "512"}}},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "3-authData",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"authData": types.M{"type": "Object"},
					},
				},
				query:  types.M{"authData.facebook": types.M{"id": "2048"}},
				update: types.M{"_auth_data_facebook": types.M{"__op": "Delete"}},
				dataObjects: []types.M{
					types.M{"_auth_data_facebook": types.M{"id": "1024"}},
					types.M{"_auth_data_facebook": types.M{"id": "2048"}},
				},
			},
			want:       types.M{},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "4-updatedAt",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"updatedAt": types.M{"type": "Date"},
					},
				},
				query:  types.M{"updatedAt": "2006-01-02T15:04:05.000Z"},
				update: types.M{"updatedAt": "2007-01-02T15:04:05.000Z"},
				dataObjects: []types.M{
					types.M{"updatedAt": types.M{"__type": "Date", "iso": "2006-01-02T15:04:05.000Z"}},
				},
			},
			want:       types.M{"updatedAt": "2007-01-02T15:04:05.000Z"},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "5-string",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "String"},
					},
				},
				query:  types.M{"key": "hi", "key2": "go"},
				update: types.M{"key": "hello", "key2": "golang"},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": "world"},
					types.M{"key": "hi", "key2": "go"},
				},
			},
			want:       types.M{"key": "hello", "key2": "golang"},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "6-bool",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "Boolean"},
					},
				},
				query:  types.M{"key2": true},
				update: types.M{"key2": false},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": false},
					types.M{"key": "hi", "key2": true},
				},
			},
			want:       types.M{"key": "hi", "key2": false},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "7-float64",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "Number"},
					},
				},
				query:  types.M{"key2": 20.5},
				update: types.M{"key2": 10.5},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": 30.5},
					types.M{"key": "hi", "key2": 20.5},
				},
			},
			want:       types.M{"key": "hi", "key2": 10.5},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "8-int",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "Number"},
					},
				},
				query:  types.M{"key2": 205},
				update: types.M{"key2": 105},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": 305},
					types.M{"key": "hi", "key2": 205},
				},
			},
			want:       types.M{"key": "hi", "key2": 105.0},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "9-time",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "Date"},
					},
				},
				query:  types.M{"key2": "2007-01-02T15:04:05.000Z"},
				update: types.M{"key2": time.Date(2008, time.January, 2, 15, 4, 5, 0, time.UTC)},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": types.M{"__type": "Date", "iso": "2006-01-02T15:04:05.000Z"}},
					types.M{"key": "hi", "key2": types.M{"__type": "Date", "iso": "2007-01-02T15:04:05.000Z"}},
				},
			},
			want:       types.M{"key": "hi", "key2": types.M{"__type": "Date", "iso": "2008-01-02T15:04:05.000Z"}},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "10-Increment",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "Number"},
					},
				},
				query:  types.M{"key2": 20},
				update: types.M{"key2": types.M{"__op": "Increment", "amount": 10}},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": 10},
					types.M{"key": "hi", "key2": 20},
				},
			},
			want:       types.M{"key": "hi", "key2": 30.0},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "11-Add",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "Array"},
					},
				},
				query:  types.M{"key": "hi"},
				update: types.M{"key2": types.M{"__op": "Add", "objects": types.S{"k3"}}},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": types.S{"j1", "j2"}},
					types.M{"key": "hi", "key2": types.S{"k1", "k2"}},
				},
			},
			want:       types.M{"key": "hi", "key2": types.S{"k1", "k2", "k3"}},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "12-Delete",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "String"},
					},
				},
				query:  types.M{"key": "hi"},
				update: types.M{"key2": types.M{"__op": "Delete"}},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": "world"},
					types.M{"key": "hi", "key2": "world"},
				},
			},
			want:       types.M{"key": "hi"},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "13-Remove",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "Array"},
					},
				},
				query:  types.M{"key": "hi"},
				update: types.M{"key2": types.M{"__op": "Remove", "objects": types.S{"k3"}}},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": types.S{"j1", "j2"}},
					types.M{"key": "hi", "key2": types.S{"k1", "k2", "k3"}},
				},
			},
			want:       types.M{"key": "hi", "key2": types.S{"k1", "k2"}},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "14-AddUnique",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "Array"},
					},
				},
				query:  types.M{"key": "hi"},
				update: types.M{"key2": types.M{"__op": "AddUnique", "objects": types.S{"k3", "k4"}}},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": types.S{"j1", "j2"}},
					types.M{"key": "hi", "key2": types.S{"k1", "k2", "k3"}},
				},
			},
			want:       types.M{"key": "hi", "key2": types.S{"k1", "k2", "k3", "k4"}},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "15-Pointer",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "Pointer", "targetClass": "user"},
					},
				},
				query:  types.M{"key": "hi"},
				update: types.M{"key2": types.M{"__type": "Pointer", "className": "user", "objectId": "123456789012345678903333"}},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": types.M{"__type": "Pointer", "className": "user", "objectId": "123456789012345678901111"}},
					types.M{"key": "hi", "key2": types.M{"__type": "Pointer", "className": "user", "objectId": "123456789012345678902222"}},
				},
			},
			want:       types.M{"key": "hi", "key2": types.M{"__type": "Pointer", "className": "user", "objectId": "123456789012345678903333"}},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "16-Date",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "Date"},
					},
				},
				query:  types.M{"key": "hi"},
				update: types.M{"key2": types.M{"__type": "Date", "iso": "2008-01-02T15:04:05.000Z"}},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": types.M{"__type": "Date", "iso": "2006-01-02T15:04:05.000Z"}},
					types.M{"key": "hi", "key2": types.M{"__type": "Date", "iso": "2007-01-02T15:04:05.000Z"}},
				},
			},
			want:       types.M{"key": "hi", "key2": types.M{"__type": "Date", "iso": "2008-01-02T15:04:05.000Z"}},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "17-File",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "File"},
					},
				},
				query:  types.M{"key": "hi"},
				update: types.M{"key2": types.M{"__type": "File", "name": "hi.png"}},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": types.M{"__type": "File", "name": "hello.png"}},
					types.M{"key": "hi", "key2": types.M{"__type": "File", "name": "world.png"}},
				},
			},
			want:       types.M{"key": "hi", "key2": types.M{"__type": "File", "name": "hi.png"}},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "18-GeoPoint",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "GeoPoint"},
					},
				},
				query:  types.M{"key": "hi"},
				update: types.M{"key2": types.M{"__type": "GeoPoint", "longitude": 30, "latitude": 30}},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": types.M{"__type": "GeoPoint", "longitude": 10, "latitude": 10}},
					types.M{"key": "hi", "key2": types.M{"__type": "GeoPoint", "longitude": 20, "latitude": 20}},
				},
			},
			want:       types.M{"key": "hi", "key2": types.M{"__type": "GeoPoint", "longitude": 30.0, "latitude": 30.0}},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "19-Object",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "Object"},
					},
				},
				query:  types.M{"key": "hi"},
				update: types.M{"key2": types.M{"key2": "go"}},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": types.M{"key1": "hello", "key2": "world"}},
					types.M{"key": "hi", "key2": types.M{"key1": "hi", "key2": "world"}},
				},
			},
			want:       types.M{"key": "hi", "key2": types.M{"key1": "hi", "key2": "go"}},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "20-Object-Delete",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "Object"},
					},
				},
				query:  types.M{"key": "hi"},
				update: types.M{"key2.key2": types.M{"__op": "Delete"}},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": types.M{"key1": "hello", "key2": "world"}},
					types.M{"key": "hi", "key2": types.M{"key1": "hi", "key2": "world"}},
				},
			},
			want:       types.M{"key": "hi", "key2": types.M{"key1": "hi"}},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "21-Object-Delete",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "Object"},
					},
				},
				query:  types.M{"key": "hi"},
				update: types.M{"key2.key2": types.M{"__op": "Delete"}, "key2.key1": "hello"},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": types.M{"key1": "hello", "key2": "world"}},
					types.M{"key": "hi", "key2": types.M{"key1": "hi", "key2": "world"}},
				},
			},
			want:       types.M{"key": "hi", "key2": types.M{"key1": "hello"}},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "21.1-Object-Increment",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "Object"},
					},
				},
				query:  types.M{"key": "hi"},
				update: types.M{"key2.key2": types.M{"__op": "Increment", "amount": 10.0}},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": types.M{"key1": "hello", "key2": 10}},
					types.M{"key": "hi", "key2": types.M{"key1": "hi", "key2": 10}},
				},
			},
			want:       types.M{"key": "hi", "key2": types.M{"key1": "hi", "key2": 20.0}},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "21.2-Object-Increment",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "Object"},
					},
				},
				query:  types.M{"key": "hi"},
				update: types.M{"key2.key2": types.M{"__op": "Increment", "amount": -10}},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": types.M{"key1": "hello", "key2": 10}},
					types.M{"key": "hi", "key2": types.M{"key1": "hi", "key2": 10.5}},
				},
			},
			want:       types.M{"key": "hi", "key2": types.M{"key1": "hi", "key2": 0.5}},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "22-text[]",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":    types.M{"type": "String"},
						"_rperm": types.M{"type": "Array"},
					},
				},
				query:  types.M{"key": "hi"},
				update: types.M{"_rperm": types.S{"hello"}},
				dataObjects: []types.M{
					types.M{"key": "hello", "_rperm": types.S{"hello", "world"}},
					types.M{"key": "hi", "_rperm": types.S{"hi", "world"}},
				},
			},
			want:       types.M{"key": "hi", "_rperm": types.S{"hello"}},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "23-Array",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "Array"},
					},
				},
				query:  types.M{"key": "hi"},
				update: types.M{"key2": types.S{"hello"}},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": types.S{"hello", "world"}},
					types.M{"key": "hi", "key2": types.S{"hi", "world"}},
				},
			},
			want:       types.M{"key": "hi", "key2": types.S{"hello"}},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "24-Unsupport",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "String"},
					},
				},
				query:  types.M{"key": "hi"},
				update: types.M{"key2": types.M{"key": "hello"}},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": "world"},
					types.M{"key": "hi", "key2": "world"},
				},
			},
			want:       nil,
			wantErr:    errs.E(errs.OperationForbidden, `Postgres doesn't support update {"key":"hello"} yet`),
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "25",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "String"},
					},
				},
				query:  types.M{"key": "haha"},
				update: types.M{"key2": "go"},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": "world"},
					types.M{"key": "hi", "key2": "world"},
				},
			},
			want:       types.M{},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
	}
	for _, tt := range tests {
		tt.initialize(tt.args.className, tt.args.schema, tt.args.dataObjects)
		got, err := p.FindOneAndUpdate(tt.args.className, tt.args.schema, tt.args.query, tt.args.update)
		tt.clean(tt.args.className)
		if reflect.DeepEqual(err, tt.wantErr) == false {
			t.Errorf("%q. PostgresAdapter.FindOneAndUpdate() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. PostgresAdapter.FindOneAndUpdate() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestPostgresAdapter_UpsertOneObject(t *testing.T) {
	db := openDB()
	p := NewPostgresAdapter("", db)
	initialize := func(className string, schema types.M, objects []types.M) {
		p.CreateClass(className, schema)
		for _, object := range objects {
			p.CreateObject(className, schema, object)
		}
	}
	clean := func(className string) {
		db.Exec(`DROP TABLE "` + className + `"`)
		db.Exec(`DROP TABLE "_SCHEMA"`)
	}
	type args struct {
		className   string
		schema      types.M
		query       types.M
		update      types.M
		dataObjects []types.M
	}
	tests := []struct {
		name       string
		args       args
		want       types.M
		wantErr    error
		initialize func(className string, schema types.M, objects []types.M)
		clean      func(className string)
	}{
		{
			name: "1",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "String"},
					},
				},
				query:  types.M{"key": "hi"},
				update: types.M{"key2": "go"},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": "world"},
					types.M{"key": "hi", "key2": "world"},
				},
			},
			want:       types.M{"key": "hi", "key2": "go"},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "2",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "String"},
					},
				},
				query:  types.M{"key": "haha"},
				update: types.M{"key2": "go"},
				dataObjects: []types.M{
					types.M{"key": "hello", "key2": "world"},
					types.M{"key": "hi", "key2": "world"},
				},
			},
			want:       types.M{"key": "haha", "key2": "go"},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
	}
	for _, tt := range tests {
		tt.initialize(tt.args.className, tt.args.schema, tt.args.dataObjects)
		err := p.UpsertOneObject(tt.args.className, tt.args.schema, tt.args.query, tt.args.update)
		got, err := p.Find(tt.args.className, tt.args.schema, tt.args.query, types.M{})
		tt.clean(tt.args.className)

		if reflect.DeepEqual(err, tt.wantErr) == false {
			t.Errorf("%q. PostgresAdapter.UpsertOneObject() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}

		if err != nil || len(got) != 1 || reflect.DeepEqual(got[0], tt.want) == false {
			t.Errorf("%q. PostgresAdapter.UpsertOneObject() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestPostgresAdapter_EnsureUniqueness(t *testing.T) {
	db := openDB()
	p := NewPostgresAdapter("", db)
	initialize := func(className string, schema types.M) {
		p.CreateClass(className, schema)
	}
	clean := func(className string) {
		db.Exec(`DROP TABLE "` + className + `"`)
		db.Exec(`DROP TABLE "_SCHEMA"`)
	}
	type args struct {
		className  string
		schema     types.M
		fieldNames []string
	}
	tests := []struct {
		name       string
		args       args
		wantErr    error
		initialize func(className string, schema types.M)
		clean      func(className string)
	}{
		{
			name: "1",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "String"},
					},
				},
				fieldNames: []string{"key"},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "2",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "String"},
					},
				},
				fieldNames: []string{"key", "key2"},
			},
			wantErr:    nil,
			initialize: initialize,
			clean:      clean,
		},
		{
			name: "3",
			args: args{
				className: "post",
				schema: types.M{
					"className": "post",
					"fields": types.M{
						"key":  types.M{"type": "String"},
						"key2": types.M{"type": "String"},
					},
				},
				fieldNames: []string{"key"},
			},
			wantErr: nil,
			initialize: func(className string, schema types.M) {
				p.CreateClass(className, schema)
				p.EnsureUniqueness(className, schema, []string{"key"})
			},
			clean: clean,
		},
	}
	for _, tt := range tests {
		tt.initialize(tt.args.className, tt.args.schema)
		err := p.EnsureUniqueness(tt.args.className, tt.args.schema, tt.args.fieldNames)
		tt.clean(tt.args.className)
		if reflect.DeepEqual(err, tt.wantErr) == false {
			t.Errorf("%q. PostgresAdapter.EnsureUniqueness() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
