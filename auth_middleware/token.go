package auth_middleware

import (
	"context"
	"fmt"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/twinj/uuid"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type AccessDetails struct {
	AccessUuid string
	UserId     int64
	Role       int64
}

// CreateToken creates a token upon user login
func CreateToken(ctx context.Context, userId, userRole int64) (*TokenDetails, error) {
	td := &TokenDetails{}

	td.AtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userId
	atClaims["exp"] = td.AtExpires
	atClaims["role"] = userRole
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userId
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}
	return td, nil
}

// CreateAuth creates an auth upon user login using token
func CreateAuth(ctx context.Context, userId int64, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := orm.RedisInstance().Set(context.Background(), td.AccessUuid, strconv.Itoa(int(userId)), at.Sub(now)).Err()
	if errAccess != nil {
		logger.Error(ctx, errAccess)
		return errAccess
	}
	logger.Info(ctx, "Added to Redis: %v", td.AccessUuid)
	errRefresh := orm.RedisInstance().Set(context.Background(), td.RefreshUuid, strconv.Itoa(int(userId)), rt.Sub(now)).Err()
	if errRefresh != nil {
		logger.Error(ctx, errRefresh)
		return errRefresh
	}
	logger.Info(ctx, "Added to Redis: %v", td.RefreshUuid)
	return nil
}

// ExtractToken extracts token from the request Authorization header
func extractToken(ctx context.Context, r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		logger.Info(ctx, "Extracted token:%v", strArr[1])
		return strArr[1]
	}
	logger.Info(ctx, "Extracted token: empty")
	return ""
}

// verifyToken verifies if a token is valid
func verifyToken(ctx context.Context, r *http.Request) (*jwt.Token, error) {
	tokenString := extractToken(ctx, r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		logger.ErrorMsg(ctx, "Error during VerifyToken: %v", err)
		return nil, err
	}
	return token, nil
}

// TokenValid checks if a token is indeed valid and not randomly generated
func TokenValid(ctx context.Context, r *http.Request) error {
	token, err := verifyToken(ctx, r)
	if err != nil {
		logger.ErrorMsg(ctx, "Error during TokenValid: %v", err)
		return err
	}
	if !token.Valid {
		logger.ErrorMsg(ctx, "Error during TokenValid: token is invalid")
		return err
	}
	_, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		logger.ErrorMsg(ctx, "Error during TokenValid: token claim is invalid")
		return err
	}
	return nil
}

// ExtractTokenMetadata extracts AccessDetails information from the token
func ExtractTokenMetadata(ctx context.Context, r *http.Request) (*AccessDetails, error) {
	token, err := verifyToken(ctx, r)
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			logger.Error(ctx, err)
			return nil, err
		}
		userId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			logger.Error(ctx, err)
			return nil, err
		}
		role, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["role"]), 10, 64)
		if err != nil {
			logger.Error(ctx, err)
			return nil, err
		}
		return &AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
			Role:       role,
		}, nil
	}
	return nil, err
}

// FetchAuth fetches the corresponding userid of an auth from cache
func FetchAuth(ctx context.Context, authD *AccessDetails) (int64, error) {
	userid, err := orm.RedisInstance().Get(context.Background(), authD.AccessUuid).Result()
	if err != nil {
		logger.ErrorMsg(ctx, "Error during FetchAuth: %v", err)
		return 0, err
	}
	logger.Info(ctx, "Fetched from Redis: user_id: %v", userid)
	userID, _ := strconv.ParseInt(userid, 10, 64)
	return userID, nil
}

func DeleteAuth(ctx context.Context, givenUuid string) (int64, error) {
	deleted, err := orm.RedisInstance().Del(context.Background(), givenUuid).Result()
	if err != nil {
		logger.Error(ctx, err)
		return 0, err
	}
	if deleted != 0 {
		logger.Info(ctx, "Successfully deleted auth token")
	}
	return deleted, nil
}

func DeleteRefresh(ctx context.Context, refreshUuid string) (int64, error) {
	deleted, err := orm.RedisInstance().Del(context.Background(), refreshUuid).Result()
	if err != nil {
		logger.Error(ctx, err)
		return 0, err
	}
	if deleted != 0 {
		logger.Info(ctx, "Successfully deleted refresh token")
	}
	return deleted, nil
}

func Refresh(c echo.Context) error {
	mapToken := map[string]string{}
	if err := c.Bind(&mapToken); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	refreshToken := mapToken["refresh_token"]

	//verify the token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "Refresh token expired")
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return c.JSON(http.StatusUnauthorized, err)
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			return c.JSON(http.StatusUnprocessableEntity, err)
		}
		userId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, "Error occurred")
		}
		userRole, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["role"]), 10, 64)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, "Error occurred")
		}
		//Delete the previous Refresh Token
		deleted, delErr := DeleteAuth(context.Background(), refreshUuid)
		if delErr != nil || deleted == 0 { //if any goes wrong
			return c.JSON(http.StatusUnauthorized, "unauthorized")
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := CreateToken(context.Background(), userId, userRole)
		if createErr != nil {
			return c.JSON(http.StatusForbidden, createErr.Error())
		}
		//save the tokens metadata to redis
		saveErr := CreateAuth(context.Background(), userId, ts)
		if saveErr != nil {
			return c.JSON(http.StatusForbidden, saveErr.Error())
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		return c.JSON(http.StatusCreated, tokens)
	} else {
		return c.JSON(http.StatusUnauthorized, "refresh expired")
	}
}
