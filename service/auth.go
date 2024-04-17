package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	mysql_sqlc "github.com/way11229/simple_merchant/repo/mysql/sqlc"

	"github.com/way11229/simple_merchant/domain"
	"github.com/way11229/simple_merchant/pkg/auth_token_maker"
	"github.com/way11229/simple_merchant/utils"
)

type AuthService struct {
	mysqlStore domain.MysqlStore

	authTokenMaker auth_token_maker.AuthTokenMaker
	redisClient    domain.RedisClient

	loginTokenExpired      time.Duration
	loginTokenCacheExpired time.Duration
}

func NewAuthService(
	mysqlStore domain.MysqlStore,

	authTokenMaker auth_token_maker.AuthTokenMaker,
	redisClient domain.RedisClient,

	loginTokenExpired time.Duration,
	loginTokenCacheExpired time.Duration,
) domain.AuthService {
	return &AuthService{
		mysqlStore: mysqlStore,

		authTokenMaker: authTokenMaker,
		redisClient:    redisClient,

		loginTokenExpired:      loginTokenExpired,
		loginTokenCacheExpired: loginTokenCacheExpired,
	}
}

func (a *AuthService) LoginUser(ctx context.Context, input *domain.LoginUserParams) (*domain.LoginUserResult, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	user, err := a.mysqlStore.GetUserByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrEmailNotFound
		}

		log.Printf("GetUserByEmail error = %v, params = %s", err, input.Email)
		return nil, domain.ErrUnknown
	}

	if !a.checkUserPassword(user.Password, input.Password) {
		return nil, domain.ErrLoginAborted
	}

	token, err := a.getAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	if err := a.createOrUpdateUserAuth(ctx, &createOrUpdateUserAuthParams{
		UserId:         user.ID,
		Token:          token,
		GetNowTimeFunc: time.Now,
	}); err != nil {
		return nil, err
	}

	go a.setAccessTokenCache(context.Background(), &setAccessTokenCacheParams{
		UserId: user.ID,
		Token:  token,
	})

	return &domain.LoginUserResult{
		Token:            token,
		EmailHasVerified: user.EmailVerifiedAt.Valid,
	}, nil
}

func (a *AuthService) CheckAccessToken(ctx context.Context, input *domain.CheckAccessTokenParams) (*domain.CheckAccessTokenResult, error) {
	payload, err := a.authTokenMaker.VerifyToken(input.AccessToken)
	if err != nil {
		return nil, domain.ErrPermissionDeny
	}

	userIdFromUniqueCode, err := utils.ConvertStringToUint(payload.UniqueCode)
	if err != nil {
		log.Printf("ConvertStringToUint error = %v, params = %s", err, payload.UniqueCode)
		return nil, domain.ErrUnknown
	}

	userId := uint32(userIdFromUniqueCode)
	if err := a.checkAccessTokenWithCacheAndDB(ctx, &checkAccessTokenWithCacheAndDBParams{
		UserId:         userId,
		Token:          input.AccessToken,
		GetNowTimeFunc: time.Now,
	}); err != nil {
		return nil, err
	}

	return &domain.CheckAccessTokenResult{
		UserId: userId,
	}, nil
}

func (a *AuthService) LogoutUser(ctx context.Context, input *domain.LogoutUserParams) error {
	if err := a.mysqlStore.DeleteUserAuthByUserId(ctx, input.UserId); err != nil {
		log.Printf("DeleteUserAuthByUserId error = %v, params = %d", err, input.UserId)
		return domain.ErrUnknown
	}

	go a.removeAccessTokenCache(context.Background(), input.UserId)

	return nil
}

/********************
 ********************
 ** private method **
 ********************
 ********************/

func (a *AuthService) checkUserPassword(encodePwd, inputPwd string) bool {
	if err := utils.BcryptCompareWithHex(encodePwd, inputPwd); err != nil {
		log.Printf("BcryptCompareWithHex error = %v", err)
		return false
	}

	return true
}

func (a *AuthService) getAccessToken(userId uint32) (string, error) {
	token, err := a.authTokenMaker.CreateToken(fmt.Sprintf("%d", userId), a.loginTokenExpired)
	if err != nil {
		log.Printf("authTokenMaker.CreateToken error = %v", err)
		return "", domain.ErrUnknown
	}

	return token, nil
}

type createOrUpdateUserAuthParams struct {
	UserId         uint32
	Token          string
	GetNowTimeFunc domain.FuncTimeType
}

func (a *AuthService) createOrUpdateUserAuth(ctx context.Context, input *createOrUpdateUserAuthParams) error {
	userAuth, err := a.mysqlStore.GetUserAuthByUserId(ctx, input.UserId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("GetUserAuthByUserId error = %v, params = %d", err, input.UserId)
		return domain.ErrUnknown
	}

	now := input.GetNowTimeFunc()
	expiredAt := now.Add(a.loginTokenExpired)
	if errors.Is(err, sql.ErrNoRows) {
		createUserAuthParams := mysql_sqlc.CreateUserAuthParams{
			UserID:    input.UserId,
			Token:     input.Token,
			ExpiredAt: expiredAt,
		}
		if err := a.mysqlStore.CreateUserAuth(ctx, createUserAuthParams); err != nil {
			log.Printf("CreateUserAuth error = %v, params = %v", err, createUserAuthParams)
			return domain.ErrUnknown
		}
	} else {
		updateUserAuthParams := mysql_sqlc.UpdateUserAuthByIdParams{
			ID:        userAuth.ID,
			Token:     input.Token,
			ExpiredAt: expiredAt,
		}
		if err := a.mysqlStore.UpdateUserAuthById(ctx, updateUserAuthParams); err != nil {
			log.Printf("UpdateUserAuthById error = %v, params = %v", err, updateUserAuthParams)
			return domain.ErrUnknown
		}
	}

	return nil
}

type setAccessTokenCacheParams struct {
	UserId uint32
	Token  string
}

func (a *AuthService) setAccessTokenCache(ctx context.Context, input *setAccessTokenCacheParams) error {
	return a.redisClient.SetEx(ctx, &domain.SetExParams{
		Key:        a.getAccessTokenCacheKey(input.UserId),
		Value:      input.Token,
		Expiration: a.loginTokenCacheExpired,
	})
}

func (a *AuthService) getAccessTokenCacheKey(userId uint32) string {
	return fmt.Sprintf("%s%d", domain.ACCESS_TOKEN_CACHE_KEY_PREFIX, userId)
}

func (a *AuthService) removeAccessTokenCache(ctx context.Context, userId uint32) error {
	return a.redisClient.Del(ctx, &domain.DelParams{
		Keys: []string{
			a.getAccessTokenCacheKey(userId),
		},
	})
}

type checkAccessTokenWithCacheAndDBParams struct {
	UserId         uint32
	Token          string
	GetNowTimeFunc domain.FuncTimeType
}

func (a *AuthService) checkAccessTokenWithCacheAndDB(ctx context.Context, input *checkAccessTokenWithCacheAndDBParams) error {
	err := a.checkAccessTokenFromCache(ctx, input)
	if err != nil && !errors.Is(err, domain.ErrRecordNotFound) {
		return err
	}

	if err == nil {
		return nil
	}

	userAuth, err := a.mysqlStore.GetUserAuthByUserId(ctx, input.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrPermissionDeny
		}

		log.Printf("GetUserAuthByUserId error = %v, params = %d", err, input.UserId)
		return domain.ErrUnknown
	}

	if userAuth.ExpiredAt.Before(input.GetNowTimeFunc()) {
		return domain.ErrPermissionDeny
	}

	go a.setAccessTokenCache(context.Background(), &setAccessTokenCacheParams{
		UserId: input.UserId,
		Token:  userAuth.Token,
	})

	return nil
}

func (a *AuthService) checkAccessTokenFromCache(ctx context.Context, input *checkAccessTokenWithCacheAndDBParams) error {
	accessTokenCache, err := a.redisClient.Get(ctx, &domain.GetParams{
		Key: a.getAccessTokenCacheKey(input.UserId),
	})
	if err != nil {
		return err
	}

	if accessTokenCache != input.Token {
		return domain.ErrPermissionDeny
	}

	return nil
}
