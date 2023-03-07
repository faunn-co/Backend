package get_available_slot_test

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/test_suite"
)

func run(reqBody *pb.GetAvailableSlotRequest, tokens *pb.Tokens, date string) (int, *pb.GetAvailableSlotResponse) {
	err, response := test_suite.NewMockTest(test_suite.GetAvailableSlot).QueryParam("date", date).Req(reqBody, tokens).Response()
	resp := response.(*pb.GetAvailableSlotResponse)
	return err, resp
}
