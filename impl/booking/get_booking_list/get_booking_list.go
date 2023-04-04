package get_booking_list

import (
	"context"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/aaronangxz/AffiliateManager/utils"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
)

type GetBookingList struct {
	c   echo.Context
	ctx context.Context
	req *pb.GetBookingListRequest
}

func New(c echo.Context) *GetBookingList {
	g := new(GetBookingList)
	g.c = c
	g.ctx = logger.NewCtx(g.c)
	logger.Info(g.ctx, "GetBookingList Initialized")
	return g
}

func (g *GetBookingList) GetBookingListImpl() ([]*pb.BookingBasic, *int64, *int64, *resp.Error) {
	if err := g.verifyGetBookingList(); err != nil {
		return nil, nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}
	var b []*pb.BookingBasic
	start, end, _, _ := utils.GetStartEndTimeFromTimeSelector(g.req.GetTimeSelector())
	if err := orm.DbInstance(g.ctx).Raw(orm.GetBookingListQuery(), start, end).Scan(&b).Error; err != nil {
		return nil, nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}
	return b, proto.Int64(start), proto.Int64(end), nil
}

func (g *GetBookingList) verifyGetBookingList() error {
	g.req = new(pb.GetBookingListRequest)
	if err := g.c.Bind(g.req); err != nil {
		return err
	}
	if err := utils.VerifyTimeSelectorFields(g.req.TimeSelector); err != nil {
		return err
	}
	return nil
}
