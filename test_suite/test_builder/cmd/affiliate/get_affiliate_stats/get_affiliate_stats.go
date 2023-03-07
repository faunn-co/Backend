package get_affiliate_stats

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

type GetAffiliateStats struct {
	req *pb.GetAffiliateStatsRequest
}

func New() *GetAffiliateStats {
	g := new(GetAffiliateStats)
	g.req = new(pb.GetAffiliateStatsRequest)
	return g
}

func (g *GetAffiliateStats) SetTimeSelector(t *pb.TimeSelector) *GetAffiliateStats {
	g.req.TimeSelector = t
	return g
}

func (g *GetAffiliateStats) checkTimeSelector() *GetAffiliateStats {
	if g.req.TimeSelector == nil {
		g.req.TimeSelector = new(pb.TimeSelector)
	}
	return g
}

func (g *GetAffiliateStats) SetBaseTimeStamp(baseTs int64) *GetAffiliateStats {
	g.checkTimeSelector()
	g.req.TimeSelector.BaseTs = proto.Int64(baseTs)
	return g
}

func (g *GetAffiliateStats) SetStartTimeStamp(startTs int64) *GetAffiliateStats {
	g.checkTimeSelector()
	g.req.TimeSelector.StartTs = proto.Int64(startTs)
	return g
}

func (g *GetAffiliateStats) SetEndTimeStamp(endTs int64) *GetAffiliateStats {
	g.checkTimeSelector()
	g.req.TimeSelector.EndTs = proto.Int64(endTs)
	return g
}

func (g *GetAffiliateStats) SetPeriod(period pb.TimeSelectorPeriod) *GetAffiliateStats {
	g.checkTimeSelector()
	g.req.TimeSelector.Period = proto.Int64(int64(period))
	return g
}

func (g *GetAffiliateStats) fillDefaults() *GetAffiliateStats {
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

func (g *GetAffiliateStats) Request() *pb.GetAffiliateStatsRequest {
	g.fillDefaults()
	return g.req
}
