package resp

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
	"net/http"
	"unicode"
)

type Error struct {
	err     error
	errCode *int64
}

func (e *Error) Error() error {
	return e.err
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
			ErrorMsg:  proto.String(string(append([]rune{unicode.ToUpper([]rune(err.err.Error())[0])}, []rune(err.err.Error()[1:])...))),
		},
	})
}

func NotAuthenticatedResp(c echo.Context) error {
	return c.JSON(http.StatusUnauthorized, pb.GenericResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(int64(pb.GlobalErrorCode_ERROR_NOT_AUTHORISED)),
			ErrorMsg:  proto.String("Not Authenticated"),
		},
	})
}

func NotAuthorisedResp(c echo.Context) error {
	return c.JSON(http.StatusForbidden, pb.GenericResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(int64(pb.GlobalErrorCode_ERROR_NO_ACCESS)),
			ErrorMsg:  proto.String("No Permission"),
		},
	})
}
