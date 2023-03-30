package encrypt

import (
	"context"
	"github.com/aaronangxz/AffiliateManager/logger"
	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(ctx context.Context, pwd string) string {
	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		logger.Error(ctx, err)
	}
	return string(hash)
}

func ComparePasswords(ctx context.Context, hashedPwd string, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	if err != nil {
		logger.Error(ctx, err)
		return false
	}
	return true
}
