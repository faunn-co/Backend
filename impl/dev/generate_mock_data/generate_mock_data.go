package generate_mock_data

import (
	"context"
	"github.com/aaronangxz/AffiliateManager/logger"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

type GenerateMockData struct {
	c   echo.Context
	ctx context.Context
	req *pb.GenerateMockDataRequest
}

func New(c echo.Context) *GenerateMockData {
	m := new(GenerateMockData)
	m.c = c
	m.ctx = logger.NewCtx(m.c)
	logger.Info(m.ctx, "GenerateMockData Initialized")
	return m
}

func (t *GenerateMockData) GenerateMockDataImpl() (*int64, *resp.Error) {
	return nil, nil
}
