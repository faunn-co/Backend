package get_referral_recent_list

import (
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/aaronangxz/AffiliateManager/utils"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
	"strconv"
	"time"
)

type GetReferralRecentList struct {
	c   echo.Context
	req *pb.GetReferralRecentListRequest
}

func New(c echo.Context) *GetReferralRecentList {
	g := new(GetReferralRecentList)
	g.c = c
	return g
}

func (g *GetReferralRecentList) GetReferralRecentListImpl() (*pb.ReferralRecent, *resp.Error) {
	if err := g.verifyGetReferralRecentList(); err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}

	var (
		c []*pb.ReferralClicks
		e []*pb.ReferralEarnings
	)

	var (
		start int64
		end   int64
	)
	p := g.c.QueryParam("period")

	switch p {
	case strconv.FormatInt(int64(pb.TimeSelectorPeriod_PERIOD_WEEK), 10):
		start, end = utils.WeekStartEndDate(time.Now().Unix())
		break
	case strconv.FormatInt(int64(pb.TimeSelectorPeriod_PERIOD_MONTH), 10):
		fallthrough
	default:
		start, end = utils.MonthStartEndDate(time.Now().Unix())
		break
	}

	end = utils.Min(end, time.Now().Unix())

	if err := orm.DbInstance(g.c).Raw(orm.GetReferralRecentClicksQuery(), g.req.GetAffiliateId(), start, end).Scan(&c).Error; err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}
	if err := orm.DbInstance(g.c).Raw(orm.GetReferralRecentEarningsQuery(), g.req.GetAffiliateId(), start, end).Scan(&e).Error; err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}

	return &pb.ReferralRecent{
		RecentClicks:   c,
		RecentEarnings: e,
		StartTime:      proto.Int64(start),
		EndTime:        proto.Int64(end),
	}, nil
}

func (g *GetReferralRecentList) verifyGetReferralRecentList() error {
	g.req = new(pb.GetReferralRecentListRequest)
	if err := g.c.Bind(g.req); err != nil {
		return err
	}
	return nil
}
