package booking

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	u "github.com/aaronangxz/AffiliateManager/test_suite/test_builder/test_utils"
	"github.com/aaronangxz/AffiliateManager/utils"
	"github.com/labstack/gommon/log"
	"google.golang.org/protobuf/proto"
	"time"
)

const (
	CitizenTicket = 8800
	TouristTicket = 9800
)

type Booking struct {
	BookingDetails *pb.BookingDetails
}

func New() *Booking {
	b := new(Booking)
	b.BookingDetails = new(pb.BookingDetails)
	orm.ENV = "TEST"
	return b
}

func (b *Booking) SetBookingStatus(status int64) *Booking {
	b.BookingDetails.BookingStatus = proto.Int64(status)
	return b
}

func (b *Booking) SetBookingDay(day string) *Booking {
	b.BookingDetails.BookingDay = proto.String(day)
	return b
}

func (b *Booking) SetBookingSlot(slot int64) *Booking {
	b.BookingDetails.BookingSlot = proto.Int64(slot)
	return b
}

func (b *Booking) SetTransactionTime(time int64) *Booking {
	b.BookingDetails.TransactionTime = proto.Int64(time)
	return b
}

func (b *Booking) SetPaymentStatus(status int64) *Booking {
	b.BookingDetails.PaymentStatus = proto.Int64(status)
	return b
}

func (b *Booking) SetCitizenTicketCount(count int64) *Booking {
	b.BookingDetails.CitizenTicketCount = proto.Int64(count)
	return b
}

func (b *Booking) SetCitizenTicketTotal(total int64) *Booking {
	b.BookingDetails.CitizenTicketTotal = proto.Int64(total)
	return b
}

func (b *Booking) SetTouristTicketCount(count int64) *Booking {
	b.BookingDetails.TouristTicketCount = proto.Int64(count)
	return b
}

func (b *Booking) SetTouristTicketTotal(total int64) *Booking {
	b.BookingDetails.TouristTicketTotal = proto.Int64(total)
	return b
}

func (b *Booking) checkCustomerInfo() *Booking {
	if b.BookingDetails.CustomerInfo == nil {
		b.BookingDetails.CustomerInfo = make([]*pb.CustomerInfo, 0)
	}
	return b
}

func (b *Booking) AddCustomer() *Booking {
	b.checkCustomerInfo()
	b.BookingDetails.CustomerInfo = append(b.BookingDetails.CustomerInfo, new(pb.CustomerInfo))
	return b
}

func (b *Booking) SetCustomerInfo(info *pb.CustomerInfo) *Booking {
	b.checkCustomerInfo()
	b.BookingDetails.CustomerInfo = append(b.BookingDetails.CustomerInfo, info)
	return b
}

func (b *Booking) SetCustomerName(name string) *Booking {
	b.checkCustomerInfo()
	b.BookingDetails.CustomerInfo[len(b.BookingDetails.CustomerInfo)-1].CustomerName = proto.String(name)
	return b
}

func (b *Booking) SetCustomerMobile(mobile string) *Booking {
	b.checkCustomerInfo()
	b.BookingDetails.CustomerInfo[len(b.BookingDetails.CustomerInfo)-1].CustomerMobile = proto.String(mobile)
	return b
}

func (b *Booking) SetCustomerEmail(email string) *Booking {
	b.checkCustomerInfo()
	b.BookingDetails.CustomerInfo[len(b.BookingDetails.CustomerInfo)-1].CustomerEmail = proto.String(email)
	return b
}

func (b *Booking) SetTicketType(ticket int64) *Booking {
	b.checkCustomerInfo()
	b.BookingDetails.CustomerInfo[len(b.BookingDetails.CustomerInfo)-1].TicketType = proto.Int64(ticket)
	return b
}

func (b *Booking) SetTicketPrice(price int64) *Booking {
	b.checkCustomerInfo()
	b.BookingDetails.CustomerInfo[len(b.BookingDetails.CustomerInfo)-1].TicketPrice = proto.Int64(price)
	return b
}

func (b *Booking) MakeDummyCustomer(count int64) *Booking {
	for i := 0; i < int(count); i++ {
		var ticketType = pb.TicketType_TICKET_TYPE_CITIZEN
		var ticketPrice = int64(CitizenTicket)
		if i%2 != 0 {
			ticketType = pb.TicketType_TICKET_TYPE_TOURIST
			ticketPrice = int64(TouristTicket)
		}
		b.SetCustomerInfo(&pb.CustomerInfo{
			CustomerName:   proto.String(fmt.Sprintf("Customer %v", i)),
			CustomerMobile: proto.String(fmt.Sprintf("123456%v", i)),
			CustomerEmail:  proto.String(fmt.Sprintf("Customer%v@mail.com", i)),
			TicketType:     proto.Int64(int64(ticketType)),
			TicketPrice:    proto.Int64(ticketPrice),
		})
	}
	return b
}

func (b *Booking) filDefaults() *Booking {
	if b.BookingDetails.BookingStatus == nil {
		b.BookingDetails.BookingStatus = proto.Int64(int64(pb.BookingStatus_BOOKING_STATUS_SUCCESS))
	}

	if b.BookingDetails.BookingDay == nil {
		b.BookingDetails.BookingDay = proto.String(utils.ConvertTimeStampYearMonthDay(time.Now().Unix()))
	}

	if b.BookingDetails.BookingSlot == nil {
		b.BookingDetails.BookingSlot = proto.Int64(u.RandomRange(0, 3))
	}

	if b.BookingDetails.TransactionTime == nil {
		b.BookingDetails.TransactionTime = proto.Int64(time.Now().Unix() - utils.MINUTE)
	}

	if b.BookingDetails.PaymentStatus == nil {
		b.BookingDetails.PaymentStatus = proto.Int64(int64(pb.PaymentStatus_PAYMENT_STATUS_SUCCESS))
	}

	if b.BookingDetails.CitizenTicketCount == nil {
		b.BookingDetails.CitizenTicketCount = proto.Int64(u.RandomRange(1, 5))
	}

	if b.BookingDetails.CitizenTicketTotal == nil {
		b.BookingDetails.CitizenTicketTotal = proto.Int64(CitizenTicket * b.BookingDetails.GetCitizenTicketCount())
	}

	if b.BookingDetails.TouristTicketCount == nil {
		b.BookingDetails.TouristTicketCount = proto.Int64(u.RandomRange(0, 5))
	}

	if b.BookingDetails.TouristTicketTotal == nil {
		b.BookingDetails.TouristTicketTotal = proto.Int64(TouristTicket * b.BookingDetails.GetTouristTicketCount())
	}

	if b.BookingDetails.CustomerInfo == nil {
		b.MakeDummyCustomer(b.BookingDetails.GetCitizenTicketCount() + b.BookingDetails.GetTouristTicketCount())
	}
	return b
}

func (b *Booking) Build() *Booking {
	b.filDefaults()

	type Booking struct {
		BookingId          *int64 `gorm:"primary_key"`
		BookingStatus      *int64
		BookingDay         *string
		BookingSlot        *int64
		TransactionTime    *int64
		PaymentStatus      *int64
		CitizenTicketCount *int64
		CitizenTicketTotal *int64
		TouristTicketCount *int64
		TouristTicketTotal *int64
		CustomerInfo       []byte
	}

	info, jsonErr := json.Marshal(b.BookingDetails.CustomerInfo)
	if jsonErr != nil {
		log.Error(jsonErr)
	}

	booking := Booking{
		BookingId:          nil,
		BookingStatus:      b.BookingDetails.BookingStatus,
		BookingDay:         b.BookingDetails.BookingDay,
		BookingSlot:        b.BookingDetails.BookingSlot,
		TransactionTime:    b.BookingDetails.TransactionTime,
		PaymentStatus:      b.BookingDetails.PaymentStatus,
		CitizenTicketCount: b.BookingDetails.CitizenTicketCount,
		CitizenTicketTotal: b.BookingDetails.CitizenTicketTotal,
		TouristTicketCount: b.BookingDetails.TouristTicketCount,
		TouristTicketTotal: b.BookingDetails.TouristTicketTotal,
		CustomerInfo:       info,
	}

	if err := orm.DbInstance(context.Background()).Table(orm.BOOKING_DETAILS_TABLE).Create(&booking).Error; err != nil {
		log.Error(err)
		return nil
	}
	b.BookingDetails.BookingId = booking.BookingId
	return b
}

func (b *Booking) TearDown() {
	if err := orm.DbInstance(context.Background()).Exec(fmt.Sprintf("DELETE FROM %v WHERE booking_id = %v", orm.BOOKING_DETAILS_TABLE, b.BookingDetails.GetBookingId())).Error; err != nil {
		log.Error(err)
	}
}
