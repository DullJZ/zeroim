package resultx

import (
	"context"
	"net/http"

	"github.com/DullJZ/zeroim/pkg/xerr"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	zrpcErr "github.com/zeromicro/x/errors"
	"google.golang.org/grpc/status"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Success(data any) *Response {
	return &Response{
		Code: 200,
		Msg:  "success",
		Data: data,
	}
}

func Fail(code int, err string) *Response {
	return &Response{
		Code: code,
		Msg:  err,
		Data: nil,
	}
}

func OkHandler(_ context.Context, v any) any {
	return Success(v)
}

func ErrHandler(name string) func(ctx context.Context, err error) (int, any) {
	return func(ctx context.Context, err error) (int, any) {
		errcode := xerr.SERVER_COMMON_ERROR
		errmsg := xerr.ErrMsg(errcode)

		causeErr := errors.Cause(err)
		if e, ok := causeErr.(*zrpcErr.CodeMsg); ok {
			errcode = e.Code
			errmsg = e.Msg
		} else {
			if gstatus, ok := status.FromError(err); ok {
				errcode = int(gstatus.Code())
				errmsg = gstatus.Message()
			}
		}

		logx.WithContext(ctx).Errorf("【%s】 err: %v", name, err)
		return http.StatusBadRequest, Fail(errcode, errmsg)
	}
}
