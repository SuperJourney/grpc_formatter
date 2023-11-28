package grpcformatter

import (
	"context"
	"errors"
	"testing"

	demo "github.com/SuperJourney/grpc_formatter/proto_for_test"
	"github.com/stretchr/testify/assert"
)

var grpcformatter = NewGrpcFormatter("test")

func TestGrpcFormatter_GetUniqKey(t *testing.T) {
	type args struct {
		req []interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test-1",
			args: args{
				req: []interface{}{context.TODO(), &demo.DemoRequest{
					Id: 1,
				}},
			},
			want: "test:common-cache:f7c57f06a1d3ce117749fc98e2111668",
		},
		{
			name: "test-2",
			args: args{
				req: []interface{}{context.TODO(), &demo.DemoRequest{
					Id: 2,
				}},
			},
			want: "test:common-cache:09347d2b5d17a91e1e71347b0f779963",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, err := grpcformatter.GetUniqKey(tt.args.req...)
			assert.Equal(t, tt.want, string(x))
			assert.NoError(t, err)
		})
	}
}

func TestGrpcFormatter_MarshalWrapper(t *testing.T) {
	type args struct {
		respes []interface{}
	}
	tests := []struct {
		name      string
		args      args
		want      string
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "test-1",
			args: args{
				respes: []interface{}{&demo.DemoResponse{
					Id:   1,
					Age:  "18 years old",
					Name: "hong",
				}, nil},
			},
			want:      `{"message":"CAESBGhvbmcaDDE4IHllYXJzIG9sZA=="}`,
			assertion: assert.NoError,
		},
		{
			name: "test-2",
			args: args{
				respes: []interface{}{nil, errors.New("error")},
			},
			want:      `{"err":"error"}`,
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := grpcformatter.MarshalWrapper(tt.args.respes...)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestGrpcFormatter_UnMarshalWrapper(t *testing.T) {
	type args struct {
		respStr []byte
		resp    interface{}
	}
	tests := []struct {
		name      string
		args      args
		want      []interface{}
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "test-1",
			args: args{
				respStr: []byte(`{"message":"CAESBGhvbmcaDDE4IHllYXJzIG9sZA=="}`),
				resp:    &demo.DemoResponse{},
			},
			want: []interface{}{&demo.DemoResponse{
				Id:   1,
				Age:  "18 years old",
				Name: "hong",
			}, nil},
			assertion: assert.NoError,
		},
		{
			name: "test-2",
			args: args{
				respStr: []byte(`{"err":"error"}`),
				resp:    &demo.DemoResponse{},
			},
			want:      []interface{}{nil, errors.New("error")},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := grpcformatter.UnMarshalWrapper(tt.args.respStr, tt.args.resp)
			tt.assertion(t, err)
			assert.EqualValues(t, tt.want, got)
		})
	}
}
