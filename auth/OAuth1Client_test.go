package auth

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/lfq7413/tomato/types"
)

func TestOAuth_Get(t *testing.T) {
	type fields struct {
		ConsumerKey     string
		ConsumerSecret  string
		AuthToken       string
		AuthTokenSecret string
		Host            string
		OAuthParams     map[string]string
	}
	type args struct {
		path   string
		params map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    types.M
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		o := &OAuth{
			ConsumerKey:     tt.fields.ConsumerKey,
			ConsumerSecret:  tt.fields.ConsumerSecret,
			AuthToken:       tt.fields.AuthToken,
			AuthTokenSecret: tt.fields.AuthTokenSecret,
			Host:            tt.fields.Host,
			OAuthParams:     tt.fields.OAuthParams,
		}
		got, err := o.Get(tt.args.path, tt.args.params)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. OAuth.Get() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. OAuth.Get() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestOAuth_Post(t *testing.T) {
	type fields struct {
		ConsumerKey     string
		ConsumerSecret  string
		AuthToken       string
		AuthTokenSecret string
		Host            string
		OAuthParams     map[string]string
	}
	type args struct {
		path   string
		params map[string]string
		body   map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    types.M
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		o := &OAuth{
			ConsumerKey:     tt.fields.ConsumerKey,
			ConsumerSecret:  tt.fields.ConsumerSecret,
			AuthToken:       tt.fields.AuthToken,
			AuthTokenSecret: tt.fields.AuthTokenSecret,
			Host:            tt.fields.Host,
			OAuthParams:     tt.fields.OAuthParams,
		}
		got, err := o.Post(tt.args.path, tt.args.params, tt.args.body)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. OAuth.Post() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. OAuth.Post() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestOAuth_Send(t *testing.T) {
	type fields struct {
		ConsumerKey     string
		ConsumerSecret  string
		AuthToken       string
		AuthTokenSecret string
		Host            string
		OAuthParams     map[string]string
	}
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    types.M
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		o := &OAuth{
			ConsumerKey:     tt.fields.ConsumerKey,
			ConsumerSecret:  tt.fields.ConsumerSecret,
			AuthToken:       tt.fields.AuthToken,
			AuthTokenSecret: tt.fields.AuthTokenSecret,
			Host:            tt.fields.Host,
			OAuthParams:     tt.fields.OAuthParams,
		}
		got, err := o.Send(tt.args.req)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. OAuth.Send() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. OAuth.Send() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestOAuth_buildRequest(t *testing.T) {
	type fields struct {
		ConsumerKey     string
		ConsumerSecret  string
		AuthToken       string
		AuthTokenSecret string
		Host            string
		OAuthParams     map[string]string
	}
	type args struct {
		method string
		path   string
		params map[string]string
		body   map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Request
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		o := &OAuth{
			ConsumerKey:     tt.fields.ConsumerKey,
			ConsumerSecret:  tt.fields.ConsumerSecret,
			AuthToken:       tt.fields.AuthToken,
			AuthTokenSecret: tt.fields.AuthTokenSecret,
			Host:            tt.fields.Host,
			OAuthParams:     tt.fields.OAuthParams,
		}
		got, err := o.buildRequest(tt.args.method, tt.args.path, tt.args.params, tt.args.body)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. OAuth.buildRequest() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. OAuth.buildRequest() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_buildParameterString(t *testing.T) {
	type args struct {
		obj map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				obj: map[string]string{
					"b": "1",
					"a": "2",
				},
			},
			want: "a=2&b=1",
		},
	}
	for _, tt := range tests {
		if got := buildParameterString(tt.args.obj); got != tt.want {
			t.Errorf("%q. buildParameterString() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_signRequest(t *testing.T) {
	type args struct {
		req             *http.Request
		oauthParameters map[string]string
		consumerSecret  string
		authTokenSecret string
		url             string
		params          map[string]string
		body            map[string]string
	}
	req, _ := http.NewRequest("GET", "http://www.baidu.com", nil)
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				req: req,
				oauthParameters: map[string]string{
					"oauth_nonce":     "abc",
					"oauth_timestamp": "1477136609",
				},
				consumerSecret:  "123",
				authTokenSecret: "123",
				url:             "http://www.baidu.com",
				params:          nil,
				body:            nil,
			},
			want: `OAuth oauth_nonce="abc", oauth_signature="VaySwdCc1dAibfofm6oWKwkYwms%3D", oauth_signature_method="HMAC-SHA1", oauth_timestamp="1477136609", oauth_version="1.0"`,
		},
	}
	for _, tt := range tests {
		if got := signRequest(tt.args.req, tt.args.oauthParameters, tt.args.consumerSecret, tt.args.authTokenSecret, tt.args.url, tt.args.params, tt.args.body); !reflect.DeepEqual(got.Header.Get("Authorization"), tt.want) {
			t.Errorf("%q. signRequest() = %v, want %v", tt.name, got.Header.Get("Authorization"), tt.want)
		}
	}
}

func Test_buildSignatureString(t *testing.T) {
	type args struct {
		method     string
		url        string
		parameters string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				method:     "GET",
				url:        "https://api.twitter.com/1.1/account/verify_credentials.json",
				parameters: "aaa=123"},
			want: "GET&https%3A%2F%2Fapi.twitter.com%2F1.1%2Faccount%2Fverify_credentials.json&aaa=123",
		},
	}
	for _, tt := range tests {
		if got := buildSignatureString(tt.args.method, tt.args.url, tt.args.parameters); got != tt.want {
			t.Errorf("%q. buildSignatureString() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_signature(t *testing.T) {
	type args struct {
		text string
		key  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{text: "abc", key: "123"},
			want: "VAsMU9SSWDe9krP3Gr56nXC2dsQ%3D",
		},
	}
	for _, tt := range tests {
		if got := signature(tt.args.text, tt.args.key); got != tt.want {
			t.Errorf("%q. signature() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_encode(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{str: "abc"},
			want: "abc",
		},
		{
			name: "2",
			args: args{str: "abc=123"},
			want: "abc%3D123",
		},
		{
			name: "3",
			args: args{str: "abc=123!'()*"},
			want: "abc%3D123%21%27%28%29%2A",
		},
	}
	for _, tt := range tests {
		if got := encode(tt.args.str); got != tt.want {
			t.Errorf("%q. encode() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_nonce(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "nonce",
			want: 30,
		},
	}
	for _, tt := range tests {
		if got := len(nonce()); got != tt.want {
			t.Errorf("%q. nonce() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
