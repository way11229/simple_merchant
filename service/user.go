package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/way11229/simple_merchant/domain"
	mysql_sqlc "github.com/way11229/simple_merchant/repo/mysql/sqlc"
	"github.com/way11229/simple_merchant/utils"
)

type UserService struct {
	mysqlStore domain.MysqlStore

	mailerClient domain.MailerClient

	loginTokenExpired                       time.Duration
	userEmailVerificationCodeLen            uint
	userEmailVerificationCodeMaxTry         uint
	userEmailVerificationCodeExpired        time.Duration
	userEmailVerificationCodeIssueLimitTime time.Duration

	verificationEmailSubject string
	verificationEmailContent string
}

func NewUserService(
	mysqlStore domain.MysqlStore,

	mailerClient domain.MailerClient,

	loginTokenExpired time.Duration,
	userEmailVerificationCodeLen uint,
	userEmailVerificationCodeMaxTry uint,
	userEmailVerificationCodeExpired time.Duration,
	userEmailVerificationCodeIssueLimitTime time.Duration,

	verificationEmailSubject string,
	verificationEmailContent string,
) domain.UserService {
	return &UserService{
		mysqlStore: mysqlStore,

		mailerClient: mailerClient,

		loginTokenExpired:                       loginTokenExpired,
		userEmailVerificationCodeLen:            userEmailVerificationCodeLen,
		userEmailVerificationCodeMaxTry:         userEmailVerificationCodeMaxTry,
		userEmailVerificationCodeExpired:        userEmailVerificationCodeExpired,
		userEmailVerificationCodeIssueLimitTime: userEmailVerificationCodeIssueLimitTime,

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

func (u *UserService) DeleteUserById(ctx context.Context, input *domain.DeleteUserByIdParams) error {
	user, err := u.mysqlStore.GetUserById(ctx, input.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrRecordNotFound
		}

		log.Printf("GetUserById error = %v, params = %d", err, input.UserId)
		return domain.ErrUnknown
	}

	if err := u.mysqlStore.TxDeleteUser(ctx, &domain.MysqlTxDeleteUserParams{
		UserId: user.ID,
	}); err != nil {
		return err
	}

	return nil
}

func (u *UserService) GetUserEmailVerificationCode(ctx context.Context, input *domain.GetUserEmailVerificationCodeParams) error {
	if err := input.Validate(); err != nil {
		return err
	}

	user, err := u.getUserByEmail(ctx, input.Email)
	if err != nil {
		return err
	}

	if user.EmailVerifiedAt.Valid {
		return domain.ErrEmailHasVerified
	}

	if err := u.checkAndDeleteExistUserEmailVerifications(ctx, user.ID); err != nil {
		return err
	}

	return u.triggerUserEmailVerification(ctx, &triggerUserEmailVerificationParams{
		User:           user,
		GetNowTimeFunc: time.Now,
	})
}

func (u *UserService) VerifyUserEmail(ctx context.Context, input *domain.VerifyUserEmailParams) error {
	if err := input.Validate(); err != nil {
		return err
	}

	user, err := u.getUserByEmail(ctx, input.Email)
	if err != nil {
		return err
	}

	verificationCode, err := u.mysqlStore.GetLastCreatedUserEmailVerificationCodeByUserId(ctx, user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrInvalidVerificationCode
		}

		log.Printf("GetLastCreatedUserEmailVerificationCodeByUserId error = %v, params = %d", err, user.ID)
		return domain.ErrUnknown
	}

	if verificationCode.VerificationCode != input.EmailVerificationCode {
		go u.mysqlStore.DecreaseUserEmailVerificationCodeMaxTryById(context.Background(), verificationCode.ID)
		return domain.ErrInvalidVerificationCode
	}

	now := input.GetNowTimeFunc()
	if verificationCode.ExpiredAt.Before(now) {
		return domain.ErrVerificationCodeExpired
	}

	if verificationCode.MaxTry == 0 {
		return domain.ErrInvalidVerificationCode
	}

	if err := u.mysqlStore.VerifyUserEmailById(ctx, user.ID); err != nil {
		log.Printf("VerifyUserEmailById error = %v, params = %d", err, user.ID)
		return domain.ErrUnknown
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

func (u *UserService) checkAndDeleteExistUserEmailVerifications(ctx context.Context, userId uint32) error {
	verification, err := u.mysqlStore.GetLastCreatedUserEmailVerificationCodeByUserId(ctx, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		log.Printf("GetLastCreatedUserEmailVerificationCodeByUserId error = %v, params = %d", err, userId)
		return domain.ErrUnknown
	}

	if time.Since(verification.CreatedAt) < u.userEmailVerificationCodeIssueLimitTime {
		return domain.ErrSendVerificationCodeInShortPeriod
	}

	u.mysqlStore.DeleteUserEmailVerificationCodeByUserId(ctx, userId)

	return nil
}

type triggerUserEmailVerificationParams struct {
	User           *mysql_sqlc.User
	GetNowTimeFunc domain.FuncTimeType
}

func (u *UserService) triggerUserEmailVerification(ctx context.Context, input *triggerUserEmailVerificationParams) error {
	verificationCode, err := utils.GenerateRandomCode(u.userEmailVerificationCodeLen)
	if err != nil {
		log.Printf("GenerateRandomCode error = %v", err)
		return domain.ErrUnknown
	}

	now := input.GetNowTimeFunc()
	if err := u.mysqlStore.TxCreateUserEmailVerificationCode(ctx, &domain.MysqlTxCreateUserEmailVerificationCodeParams{
		CreateUserEmailVerificationCodeParams: mysql_sqlc.CreateUserEmailVerificationCodeParams{
			UserID:           input.User.ID,
			Email:            input.User.Email,
			VerificationCode: verificationCode,
			MaxTry:           uint32(u.userEmailVerificationCodeMaxTry),
			ExpiredAt:        now.Add(u.userEmailVerificationCodeExpired),
		},
		AfterCreate: func() error {
			return u.mailerClient.Send(ctx, &domain.MailerClientSendParams{
				Sender: &domain.MailInfo{}, // use default
				Receiver: &domain.MailInfo{
					Name:    "",
					Address: input.User.Email,
				},
				Subject:     u.verificationEmailContent,
				HtmlContent: fmt.Sprintf(u.verificationEmailContent, verificationCode),
			})
		},
	}); err != nil {
		return err
	}

	return nil
}

func (u *UserService) getUserByEmail(ctx context.Context, email string) (*mysql_sqlc.User, error) {
	user, err := u.mysqlStore.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrEmailNotFound
		}

		log.Printf("GetUserByEmail error = %v, params = %s", err, email)
		return nil, domain.ErrUnknown
	}

	return &user, nil
}
