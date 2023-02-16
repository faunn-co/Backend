package utils

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"google.golang.org/protobuf/proto"
	"testing"
	"time"
)

func TestVerifyTimeSelectorFields(t *testing.T) {
	type args struct {
		t *pb.TimeSelector
	}
	tests := []struct {
		name          string
		args          args
		wantErr       bool
		wantErrString string
	}{
		{
			"nil",
			args{t: nil},
			true,
			"time_selector is required",
		},
		{
			"period nil",
			args{t: &pb.TimeSelector{
				BaseTs:  nil,
				StartTs: nil,
				EndTs:   nil,
				Period:  nil,
			}},
			true,
			"period is required",
		},
		{
			"invalid period",
			args{t: &pb.TimeSelector{
				BaseTs:  nil,
				StartTs: nil,
				EndTs:   nil,
				Period:  proto.Int64(0),
			}},
			true,
			"invalid period",
		},
		{
			"all ts nil",
			args{t: &pb.TimeSelector{
				BaseTs:  nil,
				StartTs: nil,
				EndTs:   nil,
				Period:  proto.Int64(int64(pb.TimeSelectorPeriod_PERIOD_MONTH)),
			}},
			true,
			"at least base_ts or either of start_ts / end_ts is required",
		},
		{
			"only end ts",
			args{t: &pb.TimeSelector{
				BaseTs:  nil,
				StartTs: nil,
				EndTs:   proto.Int64(time.Now().Unix()),
				Period:  proto.Int64(int64(pb.TimeSelectorPeriod_PERIOD_MONTH)),
			}},
			true,
			"start_ts is required",
		},
		{
			"only start ts",
			args{t: &pb.TimeSelector{
				BaseTs:  nil,
				StartTs: proto.Int64(time.Now().Unix()),
				EndTs:   nil,
				Period:  proto.Int64(int64(pb.TimeSelectorPeriod_PERIOD_MONTH)),
			}},
			true,
			"end_ts is required",
		},
		{
			"happy",
			args{t: &pb.TimeSelector{
				BaseTs:  nil,
				StartTs: proto.Int64(time.Now().Unix()),
				EndTs:   proto.Int64(time.Now().Unix()),
				Period:  proto.Int64(int64(pb.TimeSelectorPeriod_PERIOD_MONTH)),
			}},
			false,
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifyTimeSelectorFields(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyTimeSelectorFields() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err.Error() != tt.wantErrString {
				t.Errorf("VerifyTimeSelectorFields() Error = %v, errString %v", err.Error(), tt.wantErrString)
			}
		})
	}
}

func TestGetStartEndTimeFromTimeSelector(t *testing.T) {
	type args struct {
		t *pb.TimeSelector
	}
	tests := []struct {
		name  string
		args  args
		want  int64
		want1 int64
		want2 int64
		want3 int64
	}{
		{
			"happy day",
			args{t: &pb.TimeSelector{
				BaseTs:  proto.Int64(1676433600),
				StartTs: nil,
				EndTs:   nil,
				Period:  proto.Int64(int64(pb.TimeSelectorPeriod_PERIOD_DAY)),
			}},
			1676390400, 1676476799, 1676304000, 1676390399,
		},
		{
			"happy week",
			args{t: &pb.TimeSelector{
				BaseTs:  proto.Int64(1675526400),
				StartTs: nil,
				EndTs:   nil,
				Period:  proto.Int64(int64(pb.TimeSelectorPeriod_PERIOD_WEEK)),
			}},
			1675526400, 1676131199, 1674921600, 1675526399,
		},
		{
			"happy month",
			args{t: &pb.TimeSelector{
				BaseTs:  proto.Int64(1672848000),
				StartTs: nil,
				EndTs:   nil,
				Period:  proto.Int64(int64(pb.TimeSelectorPeriod_PERIOD_MONTH)),
			}},
			1672502400, 1675180799, 1670083200, 1672502399,
		},
		{
			"happy last 7d",
			args{t: &pb.TimeSelector{
				BaseTs:  proto.Int64(1675526400),
				StartTs: nil,
				EndTs:   nil,
				Period:  proto.Int64(int64(pb.TimeSelectorPeriod_PERIOD_LAST_7_DAYS)),
			}},
			1674921600, 1675526400, 1674316800, 1674921599,
		},
		{
			"happy last 28d",
			args{t: &pb.TimeSelector{
				BaseTs:  proto.Int64(1675526400),
				StartTs: nil,
				EndTs:   nil,
				Period:  proto.Int64(int64(pb.TimeSelectorPeriod_PERIOD_LAST_28_DAYS)),
			}},
			1673107200, 1675526400, 1670688000, 1673107199,
		},
		{
			"happy range",
			args{t: &pb.TimeSelector{
				BaseTs:  nil,
				StartTs: proto.Int64(1672848000),
				EndTs:   proto.Int64(1673452800),
				Period:  proto.Int64(int64(pb.TimeSelectorPeriod_PERIOD_RANGE)),
			}},
			1672848000, 1673539199, 0, 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, got3 := GetStartEndTimeFromTimeSelector(tt.args.t)
			if got != tt.want {
				t.Errorf("GetStartEndTimeFromTimeSelector() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetStartEndTimeFromTimeSelector() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("GetStartEndTimeFromTimeSelector() got2 = %v, want %v", got2, tt.want2)
			}
			if got3 != tt.want3 {
				t.Errorf("GetStartEndTimeFromTimeSelector() got3 = %v, want %v", got3, tt.want3)
			}
		})
	}
}
