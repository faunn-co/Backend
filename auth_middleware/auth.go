package auth_middleware

import (
	"context"
	"github.com/aaronangxz/AffiliateManager/logger"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

// globalAuthentication verifies if a request comes from a logged-in user
func globalAuthentication(c echo.Context) error {
	err := TokenValid(c.Request().Context(), c.Request())
	if err != nil {
		return err
	}
	return nil
}

// AdminAuthorization verifies if a request is from a logged-in admin
func AdminAuthorization(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := globalAuthentication(c); err != nil {
			logger.Error(context.Background(), err)
			return resp.NotAuthenticatedResp(c)
		}

		tokenAuth, err := ExtractTokenMetadata(c.Request().Context(), c.Request())
		if err != nil {
			logger.Error(context.Background(), err)
			return resp.NotAuthenticatedResp(c)
		}

		if tokenAuth.Role == int64(pb.UserRole_ROLE_AFFILIATE) {
			logger.ErrorMsg(context.Background(), "Role no permission")
			return resp.NotAuthorisedResp(c)
		}

		userId, err := FetchAuth(c.Request().Context(), tokenAuth)
		if err != nil {
			logger.Error(context.Background(), err)
			return resp.NotAuthenticatedResp(c)
		}

		if userId == 0 {
			logger.ErrorMsg(context.Background(), "Invalid credentials")
			return resp.NotAuthenticatedResp(c)
		}
		return next(c)
	}
}

// AffiliateAuthorization verifies if a request is from a logged-in admin or affiliate
func AffiliateAuthorization(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := globalAuthentication(c); err != nil {
			return resp.NotAuthenticatedResp(c)
		}

		tokenAuth, err := ExtractTokenMetadata(c.Request().Context(), c.Request())
		if err != nil {
			return resp.NotAuthenticatedResp(c)
		}

		userId, err := FetchAuth(c.Request().Context(), tokenAuth)
		if err != nil {
			return resp.NotAuthenticatedResp(c)
		}

		if userId == 0 {
			return resp.NotAuthenticatedResp(c)
		}
		return next(c)
	}
}

func DevAuthorization(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := globalAuthentication(c); err != nil {
			return resp.NotAuthenticatedResp(c)
		}

		tokenAuth, err := ExtractTokenMetadata(c.Request().Context(), c.Request())
		if err != nil {
			return resp.NotAuthenticatedResp(c)
		}

		if tokenAuth.Role != int64(pb.UserRole_ROLE_DEV) {
			logger.ErrorMsg(context.Background(), "Role no permission")
			return resp.NotAuthorisedResp(c)
		}

		userId, err := FetchAuth(c.Request().Context(), tokenAuth)
		if err != nil {
			return resp.NotAuthenticatedResp(c)
		}

		if userId == 0 {
			return resp.NotAuthenticatedResp(c)
		}
		return next(c)
	}
}
