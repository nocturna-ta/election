package cutresp

import (
	"errors"
	"github.com/nocturna-ta/golib/custerr"
	"github.com/nocturna-ta/golib/response"
	"github.com/nocturna-ta/golib/response/rest"
	"reflect"
	"testing"
)

func TestCustomErrorResponse(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want *rest.JSONResponse
	}{
		{
			name: "error is nil",
			args: args{
				err: nil,
			},
			want: rest.NewJSONResponse(),
		},
		{
			name: "error is standard",
			args: args{
				err: errors.New("error standard"),
			},
			want: rest.NewJSONResponse().SetError(errors.New("error standard")).SetMessage("error standard"),
		},
		{
			name: "success error custresp",
			args: args{
				err: &custerr.ErrChain{
					Message: "invalid branch ID",
					Cause:   errors.New("DB connection Timeout"),
					Type:    response.ErrBadRequest,
					Code:    40001,
				},
			},
			want: &rest.JSONResponse{
				Data:    nil,
				Message: "",
				Code:    400,
				Error: &rest.ErrorResponse{
					ErrorCode:    40001,
					ErrorMessage: "invalid branch ID",
				},
			},
		},
		{
			name: "message is empty",
			args: args{
				err: &custerr.ErrChain{
					Cause: errors.New("DB connection Timeout"),
					Type:  response.ErrTimeoutError,
				},
			},
			want: &rest.JSONResponse{
				Data:    nil,
				Code:    408,
				Message: "",
				Error: &rest.ErrorResponse{
					ErrorCode:    0,
					ErrorMessage: "DB connection Timeout",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := CustomErrorResponse(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CustomErrorResponse() = %v, want %v", got, tt.want)
			}

		})
	}
}
