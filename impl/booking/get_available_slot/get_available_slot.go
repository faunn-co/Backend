package get_available_slot

import (
	"errors"
	"fmt"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"google.golang.org/protobuf/proto"
)

type GetAvailableSlot struct {
	c echo.Context
}

func New(c echo.Context) *GetAvailableSlot {
	d := new(GetAvailableSlot)
	d.c = c
	return d
}

func (g *GetAvailableSlot) GetAvailableSlotImpl() (*string, []*pb.BookingSlot, *resp.Error) {
	if err := g.verifyGetAvailableSlot(); err != nil {
		return nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}
	date := g.c.QueryParam("date")

	var slots []*pb.BookingSlot
	if err := orm.DbInstance(g.c).Raw(fmt.Sprintf("SELECT * FROM %v.%v WHERE date = '%v'", orm.AFFILIATE_MANAGER_DB, orm.BOOKING_SLOTS_TABLE, date)).Scan(&slots).Error; err != nil {
		log.Error(err.Error)
		return proto.String(date), nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}
	return proto.String(date), slots, nil
}

func (g *GetAvailableSlot) verifyGetAvailableSlot() error {
	if g.c.QueryParam("date") == "" {
		return errors.New("invalid date")
	}
	return nil
}
