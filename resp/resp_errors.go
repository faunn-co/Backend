package resp

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
	"net/http"
)

func JSONResp(c echo.Context, code int64, msg string) error {
	return c.JSON(http.StatusOK, response(code, msg))
}

func response(code int64, msg string) pb.GenericResponse {
	return pb.GenericResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(code),
			ErrorMsg:  proto.String(msg),
		},
	}
}
