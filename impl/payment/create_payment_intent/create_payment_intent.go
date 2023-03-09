package create_payment_intent

import (
	"context"
	"errors"
	"github.com/aaronangxz/AffiliateManager/impl/tracking/track_checkout"
	"github.com/aaronangxz/AffiliateManager/logger"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentintent"
	"google.golang.org/protobuf/proto"
	"os"
)

type CreatePaymentIntent struct {
	c   echo.Context
	ctx context.Context
	req *pb.CreatePaymentIntentRequest
}

func New(c echo.Context) *CreatePaymentIntent {
	p := new(CreatePaymentIntent)
	p.c = c
	p.ctx = logger.NewCtx(p.c)
	logger.Info(p.ctx, "CreatePaymentIntent Initialized")
	return p
}

func (p *CreatePaymentIntent) CreatePaymentIntentImpl() (*string, *resp.Error) {
	if err := p.verifyCreatePaymentIntent(); err != nil {
		logger.Error(p.ctx, err)
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(p.calculateOrderAmount()),
		Currency: stripe.String(string(stripe.CurrencyMYR)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		logger.Error(p.ctx, err)
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_FAIL)
	}

	return proto.String(pi.ClientSecret), nil
}

func (p *CreatePaymentIntent) calculateOrderAmount() int64 {
	var total int64
	if p.req.GetTickets().CitizenTicketCount != nil {
		total += p.req.GetTickets().GetCitizenTicketCount() * track_checkout.CitizenTicket
	}
	if p.req.GetTickets().TouristTicketCount != nil {
		total += p.req.GetTickets().GetTouristTicketCount() * track_checkout.TouristTicket
	}
	return total / 100
}

func (p *CreatePaymentIntent) verifyCreatePaymentIntent() error {
	p.req = new(pb.CreatePaymentIntentRequest)
	if err := p.c.Bind(p.req); err != nil {
		return err
	}
	if p.req == nil {
		return errors.New("request cannot be empty")
	}
	if p.req.Tickets == nil {
		return errors.New("no ticket found")
	}
	return nil
}
