package resp

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
	"net/http"
)

type Error struct {
	err     error
	errCode *int64
}

func BuildError(err error, code pb.GlobalErrorCode) *Error {
	return &Error{
		err:     err,
		errCode: proto.Int64(int64(code)),
	}
}

func JSONResp(c echo.Context, err *Error) error {
	return c.JSON(http.StatusOK, pb.GenericResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: err.errCode,
			ErrorMsg:  proto.String(err.err.Error()),
		},
	})
}
