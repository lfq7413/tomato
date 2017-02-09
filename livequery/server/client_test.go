package server

import "testing"
import tp "github.com/lfq7413/tomato/livequery/t"
import "reflect"

func Test_transformUpdateOperators(t *testing.T) {
	type args struct {
		object         tp.M
		originalObject tp.M
	}
	tests := []struct {
		name string
		args args
		want tp.M
	}{
		{
			name: "1",
			args: args{
				object:         nil,
				originalObject: nil,
			},
			want: nil,
		},
		{
			name: "2",
			args: args{
				object:         tp.M{"key": "hello"},
				originalObject: tp.M{},
			},
			want: tp.M{"key": "hello"},
		},
		{
			name: "3",
			args: args{
				object: tp.M{
					"key": map[string]interface{}{
						"__op":   "Increment",
						"amount": 10,
					},
				},
				originalObject: tp.M{},
			},
			want: tp.M{"key": 10.0},
		},
		{
			name: "4",
			args: args{
				object: tp.M{
					"key": map[string]interface{}{
						"__op":   "Increment",
						"amount": 10,
					},
				},
				originalObject: tp.M{"key": 20},
			},
			want: tp.M{"key": 30.0},
		},
		{
			name: "5",
			args: args{
				object: tp.M{
					"key": map[string]interface{}{
						"__op":    "Add",
						"objects": []interface{}{"hello", "world"},
					},
				},
				originalObject: tp.M{},
			},
			want: tp.M{"key": []interface{}{"hello", "world"}},
		},
		{
			name: "6",
			args: args{
				object: tp.M{
					"key": map[string]interface{}{
						"__op":    "Add",
						"objects": []interface{}{"world"},
					},
				},
				originalObject: tp.M{"key": []interface{}{"hello"}},
			},
			want: tp.M{"key": []interface{}{"hello", "world"}},
		},
		{
			name: "7",
			args: args{
				object: tp.M{
					"key": map[string]interface{}{
						"__op":    "AddUnique",
						"objects": []interface{}{"hello", "world"},
					},
				},
				originalObject: tp.M{},
			},
			want: tp.M{"key": []interface{}{"hello", "world"}},
		},
		{
			name: "8",
			args: args{
				object: tp.M{
					"key": map[string]interface{}{
						"__op":    "AddUnique",
						"objects": []interface{}{"world"},
					},
				},
				originalObject: tp.M{"key": []interface{}{"hello"}},
			},
			want: tp.M{"key": []interface{}{"hello", "world"}},
		},
		{
			name: "9",
			args: args{
				object: tp.M{
					"key": map[string]interface{}{
						"__op":    "AddUnique",
						"objects": []interface{}{"hello", "world"},
					},
				},
				originalObject: tp.M{"key": []interface{}{"hello"}},
			},
			want: tp.M{"key": []interface{}{"hello", "world"}},
		},
		{
			name: "10",
			args: args{
				object: tp.M{
					"key": map[string]interface{}{
						"__op":    "Remove",
						"objects": []interface{}{"hello", "world"},
					},
				},
				originalObject: tp.M{},
			},
			want: tp.M{"key": []interface{}{}},
		},
		{
			name: "11",
			args: args{
				object: tp.M{
					"key": map[string]interface{}{
						"__op":    "Remove",
						"objects": []interface{}{"world"},
					},
				},
				originalObject: tp.M{"key": []interface{}{"hello", "world"}},
			},
			want: tp.M{"key": []interface{}{"hello"}},
		},
		{
			name: "12",
			args: args{
				object: tp.M{
					"key": map[string]interface{}{
						"__op": "Delete",
					},
				},
				originalObject: tp.M{"key": "hello"},
			},
			want: tp.M{},
		},
		{
			name: "13",
			args: args{
				object: tp.M{
					"key": map[string]interface{}{
						"__op":    "AddRelation",
						"objects": []interface{}{},
					},
				},
				originalObject: tp.M{},
			},
			want: tp.M{
				"key": map[string]interface{}{"__type": "Relation"},
			},
		},
		{
			name: "14",
			args: args{
				object: tp.M{
					"key": map[string]interface{}{
						"__op": "AddRelation",
						"objects": []interface{}{
							map[string]interface{}{
								"__type":    "Pointer",
								"className": "Player",
								"objectId":  "Vx4nudeWn",
							},
						},
					},
				},
				originalObject: tp.M{},
			},
			want: tp.M{
				"key": map[string]interface{}{
					"__type":    "Relation",
					"className": "Player",
				},
			},
		},
		{
			name: "15",
			args: args{
				object: tp.M{
					"key": map[string]interface{}{
						"__op":    "RemoveRelation",
						"objects": []interface{}{},
					},
				},
				originalObject: tp.M{},
			},
			want: tp.M{
				"key": map[string]interface{}{"__type": "Relation"},
			},
		},
		{
			name: "16",
			args: args{
				object: tp.M{
					"key": map[string]interface{}{
						"__op": "RemoveRelation",
						"objects": []interface{}{
							map[string]interface{}{
								"__type":    "Pointer",
								"className": "Player",
								"objectId":  "Vx4nudeWn",
							},
						},
					},
				},
				originalObject: tp.M{},
			},
			want: tp.M{
				"key": map[string]interface{}{
					"__type":    "Relation",
					"className": "Player",
				},
			},
		},
	}
	for _, tt := range tests {
		transformUpdateOperators(tt.args.object, tt.args.originalObject)
		if reflect.DeepEqual(tt.args.object, tt.want) == false {
			t.Errorf("%q. transformUpdateOperators() = %v, want %v", tt.name, tt.args.object, tt.want)
		}
	}
}
