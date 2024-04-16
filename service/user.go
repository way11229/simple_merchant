package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/way11229/simple_merchant/domain"
	mysql_sqlc "github.com/way11229/simple_merchant/repo/mysql/sqlc"
	"github.com/way11229/simple_merchant/utils"
)

type UserService struct {
	mysqlStore domain.MysqlStore

	mailerClient domain.MailerClient

	loginTokenExpireSeconds                    time.Duration
	userEmailVerificationCodeLen               uint
	userEmailVerificationCodeMaxTry            uint
	userEmailVerificationCodeExpiredSeconds    time.Duration
	userEmailVerificationCodeIssueLimitSeconds time.Duration

	verificationEmailSubject string
	verificationEmailContent string
}

func NewUserService(
	mysqlStore domain.MysqlStore,

	mailerClient domain.MailerClient,

	loginTokenExpireSeconds time.Duration,
	userEmailVerificationCodeLen uint,
	userEmailVerificationCodeMaxTry uint,
	userEmailVerificationCodeExpiredSeconds time.Duration,
	userEmailVerificationCodeIssueLimitSeconds time.Duration,

	verificationEmailSubject string,
	verificationEmailContent string,
) domain.UserService {
	return &UserService{
		mysqlStore: mysqlStore,

		mailerClient: mailerClient,

		loginTokenExpireSeconds:                    loginTokenExpireSeconds,
		userEmailVerificationCodeLen:               userEmailVerificationCodeLen,
		userEmailVerificationCodeMaxTry:            userEmailVerificationCodeMaxTry,
		userEmailVerificationCodeExpiredSeconds:    userEmailVerificationCodeExpiredSeconds,
		userEmailVerificationCodeIssueLimitSeconds: userEmailVerificationCodeIssueLimitSeconds,

		verificationEmailSubject: verificationEmailSubject,
		verificationEmailContent: verificationEmailContent,
	}
}

func (u *UserService) CreateUser(ctx context.Context, input *domain.CreateUserParams) (*domain.CreateUserResult, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	if err := u.checkUserEmailExist(ctx, input.Email); err != nil {
		return nil, err
	}

	encodePwd, err := u.encodeUserPassword(input.Password)
	if err != nil {
		return nil, err
	}

	txCreateUserParams := &domain.MysqlTxCreateUserParams{
		CreateUserParams: mysql_sqlc.CreateUserParams{
			Name:     input.Name,
			Email:    input.Email,
			Password: encodePwd,
		},
	}
	newUserId, err := u.mysqlStore.TxCreateUser(ctx, txCreateUserParams)
	if err != nil {
		log.Printf("CreateUser error = %v, params = %v", err, txCreateUserParams)
		return nil, domain.ErrUnknown
	}

	return &domain.CreateUserResult{
		UserId: newUserId,
	}, nil
}

func (u *UserService) DeleteUserById(ctx context.Context, userId uint32) error {
	user, err := u.mysqlStore.GetUserById(ctx, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrRecordNotFound
		}

		log.Printf("GetUserById error = %v, params = %d", err, userId)
		return domain.ErrUnknown
	}

	if err := u.mysqlStore.TxDeleteUser(ctx, &domain.MysqlTxDeleteUserParams{
		UserId: user.ID,
	}); err != nil {
		return err
	}

	return nil
}

/********************
 ********************
 ** private method **
 ********************
 ********************/

func (u *UserService) checkUserEmailExist(ctx context.Context, email string) error {
	if _, err := u.mysqlStore.GetUserByEmail(ctx, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		log.Printf("GetUserByEmail error = %v, params = %s", err, email)
		return domain.ErrUnknown
	}

	return domain.ErrUserEmailDuplicated
}

func (u *UserService) encodeUserPassword(pwd string) (string, error) {
	encodePwd, err := utils.BcryptEncryptToHex(pwd)
	if err != nil {
		log.Println(err)
		return "", domain.ErrUnknown
	}

	return encodePwd, nil
}

func (u *UserService) checkUserPassword(encodePwd, inputPwd string) bool {
	if err := utils.BcryptCompareWithHex(encodePwd, inputPwd); err != nil {
		log.Println(err)
		return false
	}

	return true
}
