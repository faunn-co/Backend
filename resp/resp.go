package resp

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
	"net/http"
)

func GetAvailableSlotResponseJSON(c echo.Context, date *string, slots []*pb.BookingSlot) error {
	return c.JSON(http.StatusOK, getAvailableSlotResponse(date, slots))
}

func getAvailableSlotResponse(date *string, slots []*pb.BookingSlot) pb.GetAvailableSlotResponse {
	return pb.GetAvailableSlotResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(0),
			ErrorMsg:  proto.String("success"),
		},
		Date:         date,
		BookingSlots: slots,
	}
}
