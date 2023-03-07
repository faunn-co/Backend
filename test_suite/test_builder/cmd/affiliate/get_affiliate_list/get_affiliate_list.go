package get_affiliate_list

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

type GetAffiliateList struct {
	req *pb.GetAffiliateListRequest
}

func New() *GetAffiliateList {
	g := new(GetAffiliateList)
	g.req = new(pb.GetAffiliateListRequest)
	return g
}

func (g *GetAffiliateList) SetTimeSelector(t *pb.TimeSelector) *GetAffiliateList {
	g.req.TimeSelector = t
	return g
}

func (g *GetAffiliateList) checkTimeSelector() *GetAffiliateList {
	if g.req.TimeSelector == nil {
		g.req.TimeSelector = new(pb.TimeSelector)
	}
	return g
}

func (g *GetAffiliateList) SetBaseTimeStamp(baseTs int64) *GetAffiliateList {
	g.checkTimeSelector()
	g.req.TimeSelector.BaseTs = proto.Int64(baseTs)
	return g
}

func (g *GetAffiliateList) SetStartTimeStamp(startTs int64) *GetAffiliateList {
	g.checkTimeSelector()
	g.req.TimeSelector.StartTs = proto.Int64(startTs)
	return g
}

func (g *GetAffiliateList) SetEndTimeStamp(endTs int64) *GetAffiliateList {
	g.checkTimeSelector()
	g.req.TimeSelector.EndTs = proto.Int64(endTs)
	return g
}

func (g *GetAffiliateList) SetPeriod(period pb.TimeSelectorPeriod) *GetAffiliateList {
	g.checkTimeSelector()
	g.req.TimeSelector.Period = proto.Int64(int64(period))
	return g
}

func (g *GetAffiliateList) fillDefaults() *GetAffiliateList {
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

func (g *GetAffiliateList) Request() *pb.GetAffiliateListRequest {
	g.fillDefaults()
	return g.req
}
