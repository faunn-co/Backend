package user_registration

import (
	"context"
	"errors"
	"github.com/aaronangxz/AffiliateManager/impl/verification/affiliate_verification"
	"github.com/aaronangxz/AffiliateManager/impl/verification/user_verification"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
	"time"
)

const (
	ENTITY_NAME_MAX_LENGTH   = 50
	USER_NAME_MAX_LENGTH     = 50
	REFERRAL_CODE_MAX_LENGTH = 10
)

type UserRegistration struct {
	c   echo.Context
	ctx context.Context
	req *pb.UserRegistrationRequest
}

func New(c echo.Context) *UserRegistration {
	u := new(UserRegistration)
	u.c = c
	u.ctx = logger.NewCtx(u.c)
	logger.Info(u.ctx, "GetAffiliateTrend Initialized")
	return u
}

func (u *UserRegistration) UserRegistrationImpl() *resp.Error {
	if err := u.verifyUserRegistration(); err != nil {
		return resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}

	if err := u.verifyUserName(); err != nil {
		return resp.BuildError(err, pb.GlobalErrorCode_ERROR_USER_NAME_EXISTS)
	}

	if err := u.verifyUserEmail(); err != nil {
		return resp.BuildError(err, pb.GlobalErrorCode_ERROR_USER_EMAIL_EXISTS)
	}

	if err := u.verifyEntityName(); err != nil {
		return resp.BuildError(err, pb.GlobalErrorCode_ERROR_ENTITY_NAME_EXISTS)
	}

	if err := u.verifyReferralCode(); err != nil {
		return resp.BuildError(err, pb.GlobalErrorCode_ERROR_REFERRAL_CODE_EXISTS)
	}

	type User struct {
		UserId          *int64 `gorm:"primary_key"`
		UserName        *string
		UserEmail       *string
		UserContact     *string
		UserRole        *int64
		CreateTimestamp *int64
	}

	user := User{
		UserName:        u.req.UserName,
		UserEmail:       u.req.UserEmail,
		UserContact:     u.req.UserContact,
		UserRole:        proto.Int64(int64(pb.UserRole_ROLE_AFFILIATE)),
		CreateTimestamp: proto.Int64(time.Now().Unix()),
	}

	if err := orm.DbInstance(u.ctx).Table(orm.USER_TABLE).Create(&user).Error; err != nil {
		return resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}

	auth := &pb.UserAuth{
		UserId:       user.UserId,
		UserPassword: u.req.UserPassword,
	}

	if err := orm.DbInstance(u.ctx).Table(orm.USER_AUTH_TABLE).Create(&auth).Error; err != nil {
		return resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}

	affiliate := &pb.AffiliateDetailsDb{
		UserId:             user.UserId,
		EntityName:         u.req.EntityName,
		AffiliateType:      u.req.AffiliateType,
		UniqueReferralCode: u.req.PreferredReferralCode,
	}

	if err := orm.DbInstance(u.ctx).Table(orm.AFFILIATE_DETAILS_TABLE).Create(affiliate).Error; err != nil {
		return resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}

	return nil
}

func (u *UserRegistration) verifyUserName() error {
	if err := user_verification.New(u.c, u.ctx).VerifyUserName(u.req.GetUserName()); err == nil {
		return errors.New("user name already exist")
	}
	return nil
}

func (u *UserRegistration) verifyUserEmail() error {
	if err := user_verification.New(u.c, u.ctx).VerifyUserEmail(u.req.GetUserEmail()); err == nil {
		return errors.New("user email already exist")
	}
	return nil
}

func (u *UserRegistration) verifyEntityName() error {
	if err := affiliate_verification.New(u.c, u.ctx).VerifyEntityName(u.req.GetEntityName()); err == nil {
		return errors.New("entity name already exist")
	}
	return nil
}

func (u *UserRegistration) verifyReferralCode() error {
	if err := affiliate_verification.New(u.c, u.ctx).VerifyReferralCode(u.req.GetPreferredReferralCode()); err == nil {
		return errors.New("referral code already exist")
	}
	return nil
}

func (u *UserRegistration) verifyUserRegistration() error {
	u.req = new(pb.UserRegistrationRequest)
	if err := u.c.Bind(u.req); err != nil {
		return err
	}
	if u.req == nil {
		return errors.New("request body is empty")
	}

	if u.req.EntityName == nil {
		return errors.New("entity name is required")
	}

	if u.req.EntityIdentifier == nil {
		return errors.New("entity identifier is required")
	}

	if u.req.UserName == nil {
		return errors.New("user name is required")
	}

	if u.req.UserPassword == nil {
		return errors.New("user password is required")
	}

	if u.req.UserEmail == nil {
		return errors.New("user email is required")
	}

	if u.req.UserContact == nil {
		return errors.New("user contact is required")
	}

	if u.req.AffiliateType == nil {
		return errors.New("affiliate type is required")
	}

	if len(u.req.GetEntityName()) > ENTITY_NAME_MAX_LENGTH {
		return errors.New("entity name cannot be longer than 50 characters")
	}

	if isContainsSpecialChar(u.req.GetEntityName()) {
		return errors.New("entity name contains illegal characters")
	}

	if isContainsSpecialChar(u.req.GetEntityIdentifier()) || isContainsSpace(u.req.GetEntityIdentifier()) {
		return errors.New("entity identifier contains illegal characters")
	}

	if len(u.req.GetUserName()) > USER_NAME_MAX_LENGTH {
		return errors.New("user name cannot be longer than 50 characters")
	}

	if isContainsSpecialChar(u.req.GetUserName()) || isContainsSpace(u.req.GetUserName()) {
		return errors.New("user name contains illegal characters")
	}

	if isContainsSpace(u.req.GetUserPassword()) {
		return errors.New("user password cannot contain spaces")
	}

	if isContainsAtSign(u.req.GetUserEmail()) {
		return errors.New("user email format is incorrect")
	}

	if isContainsNonNumeric(u.req.GetUserContact()) {
		return errors.New("user contact format is incorrect")
	}

	if _, exists := pb.AffiliateType_name[int32(u.req.GetAffiliateType())]; !exists {
		return errors.New("invalid affiliate type")
	}

	if len(u.req.GetPreferredReferralCode()) > REFERRAL_CODE_MAX_LENGTH {
		return errors.New("referral code cannot be longer than 50 characters")
	}

	if isContainsSpecialChar(u.req.GetPreferredReferralCode()) || isContainsSpace(u.req.GetPreferredReferralCode()) {
		return errors.New("referral code contains illegal characters")
	}
	return nil
}
