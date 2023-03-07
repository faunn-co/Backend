package get_referral_trend

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/utils"
	"google.golang.org/protobuf/proto"
	"time"
)

var (
	defaultBaseTs       = time.Now().Unix()
	defaultPeriod       = int64(pb.TimeSelectorPeriod_PERIOD_MONTH)
	defaultStartTs      = time.Now().Unix() - utils.DAY
	defaultEndTs        = time.Now().Unix()
	defaultTimeSelector = &pb.TimeSelector{
		BaseTs:  proto.Int64(defaultBaseTs),
		StartTs: nil,
		EndTs:   nil,
		Period:  proto.Int64(defaultPeriod),
	}
)

type GetReferralTrend struct {
	req *pb.GetReferralTrendRequest
}

func New() *GetReferralTrend {
	g := new(GetReferralTrend)
	g.req = new(pb.GetReferralTrendRequest)
	return g
}

func (g *GetReferralTrend) SetTimeSelector(t *pb.TimeSelector) *GetReferralTrend {
	g.req.TimeSelector = t
	return g
}

func (g *GetReferralTrend) checkTimeSelector() *GetReferralTrend {
	if g.req.TimeSelector == nil {
		g.req.TimeSelector = new(pb.TimeSelector)
	}
	return g
}

func (g *GetReferralTrend) SetBaseTimeStamp(baseTs int64) *GetReferralTrend {
	g.checkTimeSelector()
	g.req.TimeSelector.BaseTs = proto.Int64(baseTs)
	return g
}

func (g *GetReferralTrend) SetStartTimeStamp(startTs int64) *GetReferralTrend {
	g.checkTimeSelector()
	g.req.TimeSelector.StartTs = proto.Int64(startTs)
	return g
}

func (g *GetReferralTrend) SetEndTimeStamp(endTs int64) *GetReferralTrend {
	g.checkTimeSelector()
	g.req.TimeSelector.EndTs = proto.Int64(endTs)
	return g
}

func (g *GetReferralTrend) SetPeriod(period pb.TimeSelectorPeriod) *GetReferralTrend {
	g.checkTimeSelector()
	g.req.TimeSelector.Period = proto.Int64(int64(period))
	return g
}

func (g *GetReferralTrend) fillDefaults() *GetReferralTrend {
	if g.req.TimeSelector == nil {
		g.SetTimeSelector(defaultTimeSelector)
	} else {
		if g.req.TimeSelector.Period != nil {
			if g.req.GetTimeSelector().GetPeriod() == int64(pb.TimeSelectorPeriod_PERIOD_RANGE) {
				if g.req.GetTimeSelector().StartTs == nil {
					g.SetStartTimeStamp(defaultStartTs)
				}
				if g.req.GetTimeSelector().EndTs == nil {
					g.SetEndTimeStamp(defaultEndTs)
				}
			} else {
				if g.req.GetTimeSelector().BaseTs == nil {
					g.SetBaseTimeStamp(defaultBaseTs)
				}
			}
		}
	}
	return g
}

func (g *GetReferralTrend) Request() *pb.GetReferralTrendRequest {
	g.fillDefaults()
	return g.req
}
