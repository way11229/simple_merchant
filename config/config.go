package config

import (
	"fmt"

	"github.com/way11229/simple_merchant/utils"
)

type Config struct {
	SqlDriverName        string `mapstructure:"SQL_DRIVER_NAME"`
	SqlDataSourceName    string `mapstructure:"SQL_DATA_SOURCE_NAME"`
	MigrationSourceURL   string `mapstructure:"MIGRATION_SOURCE_URL"`
	MigrationDatabaseURL string `mapstructure:"MIGRATION_DATABASE_URL"`

	LoginTokenExpireSeconds                    uint `mapstructure:"LOGIN_TOKEN_EXPIRE_SECONDS"`
	UserEmailVerificationCodeLen               uint `mapstructure:"USER_EMAIL_VERIFICATION_CODE_LEN"`
	UserEmailVerificationCodeMaxTry            uint `mapstructure:"USER_EMAIL_VERIFICATION_CODE_MAX_TRY"`
	UserEmailVerificationCodeExpiredSeconds    uint `mapstructure:"USER_EMAIL_VERIFICATION_CODE_EXPIRED_SECONDS"`
	UserEmailVerificationCodeIssueLimitSeconds uint `mapstructure:"USER_EMAIL_VERIFICATION_CODE_ISSUE_LIMIT_SECONDS"`

	VerificationEmailSubject string `mapstructure:"VERIFICATION_EMAIL_SUBJECT"`
	VerificationEmailContent string `mapstructure:"VERIFICATION_EMAIL_CONTENT"`
}

func NewConfig() *Config {
	return &Config{
		SqlDriverName:        utils.GetEnv("SQL_DRIVER_NAME"),
		SqlDataSourceName:    utils.GetEnv("SQL_DATA_SOURCE_NAME"),
		MigrationSourceURL:   utils.GetEnv("MIGRATION_SOURCE_URL"),
		MigrationDatabaseURL: utils.GetEnv("MIGRATION_DATABASE_URL"),

		LoginTokenExpireSeconds:                    convertStringToUintAndPanicIfError(utils.GetEnv("LOGIN_TOKEN_EXPIRE_SECONDS")),
		UserEmailVerificationCodeLen:               convertStringToUintAndPanicIfError(utils.GetEnv("USER_EMAIL_VERIFICATION_CODE_LEN")),
		UserEmailVerificationCodeMaxTry:            convertStringToUintAndPanicIfError(utils.GetEnv("USER_EMAIL_VERIFICATION_CODE_MAX_TRY")),
		UserEmailVerificationCodeExpiredSeconds:    convertStringToUintAndPanicIfError(utils.GetEnv("USER_EMAIL_VERIFICATION_CODE_EXPIRED_SECONDS")),
		UserEmailVerificationCodeIssueLimitSeconds: convertStringToUintAndPanicIfError(utils.GetEnv("USER_EMAIL_VERIFICATION_CODE_ISSUE_LIMIT_SECONDS")),

		VerificationEmailSubject: utils.GetEnv("VERIFICATION_EMAIL_SUBJECT"),
		VerificationEmailContent: utils.GetEnv("VERIFICATION_EMAIL_CONTENT"),
	}
}

func convertStringToUintAndPanicIfError(str string) uint {
	rtn, err := utils.ConvertStringToUint(str)
	if err != nil {
		panic(fmt.Sprintf("ConvertStringToUint error = %v", err))
	}

	return rtn
}
