package auth_middleware

import (
	"context"
	"errors"
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
		if extErr := ExtractMetaFromToken(c); extErr != nil {
			logger.Error(context.Background(), extErr)
			return resp.NotAuthenticatedResp(c)
		}
		if GetUserRoleFromToken(c) == int64(pb.UserRole_ROLE_AFFILIATE) {
			logger.ErrorMsg(context.Background(), "Role no permission")
			return resp.NotAuthorisedResp(c)
		}
		return next(c)
	}
}

// AffiliateAuthorization verifies if a request is from a logged-in admin or affiliate
func AffiliateAuthorization(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if extErr := ExtractMetaFromToken(c); extErr != nil {
			logger.Error(context.Background(), extErr)
			return resp.NotAuthenticatedResp(c)
		}
		return next(c)
	}
}

func DevAuthorization(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if extErr := ExtractMetaFromToken(c); extErr != nil {
			logger.Error(context.Background(), extErr)
			return resp.NotAuthenticatedResp(c)
		}
		if GetUserRoleFromToken(c) != int64(pb.UserRole_ROLE_DEV) {
			logger.ErrorMsg(context.Background(), "Role no permission")
			return resp.NotAuthorisedResp(c)
		}
		return next(c)
	}
}

func ExtractMetaFromToken(c echo.Context) error {
	if err := globalAuthentication(c); err != nil {
		return err
	}

	tokenAuth, err := ExtractTokenMetadata(c.Request().Context(), c.Request())
	if err != nil {
		return err
	}

	userId, err := FetchAuth(c.Request().Context(), tokenAuth)
	if err != nil {
		return err
	}

	if userId == 0 {
		return errors.New("failed to fetch user_id")
	}

	c.Set("user_role", tokenAuth.Role)
	c.Set("user_id", userId)
	return nil
}

func GetUserRoleFromToken(c echo.Context) int64 {
	return c.Get("user_role").(int64)
}

func GetUserIdFromToken(c echo.Context) int64 {
	return c.Get("user_id").(int64)
}
