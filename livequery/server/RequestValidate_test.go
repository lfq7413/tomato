package server

import (
	"errors"
	"reflect"
	"testing"

	tp "github.com/lfq7413/tomato/livequery/t"
)

func Test_validateGeneral(t *testing.T) {
	type args struct {
		data tp.M
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "1",
			args:    args{data: tp.M{}},
			wantErr: errors.New("need op"),
		},
		{
			name:    "2",
			args:    args{data: tp.M{"op": 1024}},
			wantErr: errors.New("op is not string"),
		},
		{
			name:    "3",
			args:    args{data: tp.M{"op": "hello"}},
			wantErr: errors.New("op is not in [connect, subscribe, unsubscribe, update]"),
		},
		{
			name:    "4",
			args:    args{data: tp.M{"op": "connect"}},
			wantErr: nil,
		},
		{
			name:    "5",
			args:    args{data: tp.M{"op": "subscribe"}},
			wantErr: nil,
		},
		{
			name:    "6",
			args:    args{data: tp.M{"op": "unsubscribe"}},
			wantErr: nil,
		},
		{
			name:    "7",
			args:    args{data: tp.M{"op": "update"}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		if err := validateGeneral(tt.args.data); reflect.DeepEqual(err, tt.wantErr) == false {
			t.Errorf("%q. validateGeneral() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func Test_validateConnect(t *testing.T) {
	type args struct {
		data tp.M
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "1",
			args:    args{data: tp.M{}},
			wantErr: nil,
		},
		{
			name:    "2",
			args:    args{data: tp.M{"applicationId": 1024}},
			wantErr: errors.New("applicationId is not string"),
		},
		{
			name:    "3",
			args:    args{data: tp.M{"applicationId": "1024"}},
			wantErr: nil,
		},
		{
			name:    "4",
			args:    args{data: tp.M{"masterKey": 1024}},
			wantErr: errors.New("masterKey is not string"),
		},
		{
			name:    "5",
			args:    args{data: tp.M{"masterKey": "1024"}},
			wantErr: nil,
		},
		{
			name:    "6",
			args:    args{data: tp.M{"clientKey": 1024}},
			wantErr: errors.New("clientKey is not string"),
		},
		{
			name:    "7",
			args:    args{data: tp.M{"clientKey": "1024"}},
			wantErr: nil,
		},
		{
			name:    "8",
			args:    args{data: tp.M{"restAPIKey": 1024}},
			wantErr: errors.New("restAPIKey is not string"),
		},
		{
			name:    "9",
			args:    args{data: tp.M{"restAPIKey": "1024"}},
			wantErr: nil,
		},
		{
			name:    "10",
			args:    args{data: tp.M{"javascriptKey": 1024}},
			wantErr: errors.New("javascriptKey is not string"),
		},
		{
			name:    "11",
			args:    args{data: tp.M{"javascriptKey": "1024"}},
			wantErr: nil,
		},
		{
			name:    "12",
			args:    args{data: tp.M{"windowsKey": 1024}},
			wantErr: errors.New("windowsKey is not string"),
		},
		{
			name:    "13",
			args:    args{data: tp.M{"windowsKey": "1024"}},
			wantErr: nil,
		},
		{
			name:    "14",
			args:    args{data: tp.M{"sessionToken": 1024}},
			wantErr: errors.New("sessionToken is not string"),
		},
		{
			name:    "15",
			args:    args{data: tp.M{"sessionToken": "1024"}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		if err := validateConnect(tt.args.data); reflect.DeepEqual(err, tt.wantErr) == false {
			t.Errorf("%q. validateConnect() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func Test_validateSubscribe(t *testing.T) {
	type args struct {
		data tp.M
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "1",
			args:    args{data: tp.M{}},
			wantErr: errors.New("need requestId"),
		},
		{
			name:    "2",
			args:    args{data: tp.M{"requestId": "1024"}},
			wantErr: errors.New("requestId is not number"),
		},
		{
			name:    "3",
			args:    args{data: tp.M{"requestId": 1024.0}},
			wantErr: errors.New("need query"),
		},
		{
			name: "4",
			args: args{data: tp.M{
				"requestId": 1024.0,
				"query":     1024,
			}},
			wantErr: errors.New("query is not object"),
		},
		{
			name: "5",
			args: args{data: tp.M{
				"requestId": 1024.0,
				"query": map[string]interface{}{
					"className": "post",
					"where":     map[string]interface{}{},
				},
			}},
			wantErr: nil,
		},
		{
			name: "6",
			args: args{data: tp.M{
				"requestId": 1024.0,
				"query": map[string]interface{}{
					"className": "post",
					"where":     map[string]interface{}{},
				},
				"sessionToken": 1024,
			}},
			wantErr: errors.New("sessionToken is not string"),
		},
		{
			name: "7",
			args: args{data: tp.M{
				"requestId": 1024.0,
				"query": map[string]interface{}{
					"className": "post",
					"where":     map[string]interface{}{},
				},
				"sessionToken": "1024",
			}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		if err := validateSubscribe(tt.args.data); reflect.DeepEqual(err, tt.wantErr) == false {
			t.Errorf("%q. validateSubscribe() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func Test_validateUnsubscribe(t *testing.T) {
	type args struct {
		data tp.M
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "1",
			args:    args{data: tp.M{}},
			wantErr: errors.New("need requestId"),
		},
		{
			name:    "2",
			args:    args{data: tp.M{"requestId": "hello"}},
			wantErr: errors.New("requestId is not number"),
		},
		{
			name:    "3",
			args:    args{data: tp.M{"requestId": 1024.0}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		if err := validateUnsubscribe(tt.args.data); reflect.DeepEqual(err, tt.wantErr) == false {
			t.Errorf("%q. validateUnsubscribe() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func Test_validateQuery(t *testing.T) {
	type args struct {
		data tp.M
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "1",
			args:    args{data: tp.M{}},
			wantErr: errors.New("need className"),
		},
		{
			name:    "2",
			args:    args{data: tp.M{"className": 1024}},
			wantErr: errors.New("className is not string"),
		},
		{
			name:    "3",
			args:    args{data: tp.M{"className": "post"}},
			wantErr: errors.New("need where"),
		},
		{
			name: "4",
			args: args{data: tp.M{
				"className": "post",
				"where":     1024,
			}},
			wantErr: errors.New("where is not object"),
		},
		{
			name: "5",
			args: args{data: tp.M{
				"className": "post",
				"where":     map[string]interface{}{},
			}},
			wantErr: nil,
		},
		{
			name: "6",
			args: args{data: tp.M{
				"className": "post",
				"where":     map[string]interface{}{},
				"fields":    1024,
			}},
			wantErr: errors.New("fields is not []string"),
		},
		{
			name: "7",
			args: args{data: tp.M{
				"className": "post",
				"where":     map[string]interface{}{},
				"fields":    []interface{}{},
			}},
			wantErr: errors.New("minItems is not 1"),
		},
		{
			name: "8",
			args: args{data: tp.M{
				"className": "post",
				"where":     map[string]interface{}{},
				"fields":    []interface{}{"hello", 1024},
			}},
			wantErr: errors.New("fields is not []string"),
		},
		{
			name: "9",
			args: args{data: tp.M{
				"className": "post",
				"where":     map[string]interface{}{},
				"fields":    []interface{}{"hello", "world"},
			}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		if err := validateQuery(tt.args.data); reflect.DeepEqual(err, tt.wantErr) == false {
			t.Errorf("%q. validateQuery() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
