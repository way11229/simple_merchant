package config

import "github.com/way11229/simple_merchant/utils"

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

	VerifyEmailSubject string `mapstructure:"VERIFY_EMAIL_SUBJECT"`
	VerifyEmailContent string `mapstructure:"VERIFY_EMAIL_CONTENT"`
}

func NewConfig() *Config {
	return &Config{
		SqlDriverName:        utils.GetEnv("SQL_DRIVER_NAME"),
		SqlDataSourceName:    utils.GetEnv("SQL_DATA_SOURCE_NAME"),
		MigrationSourceURL:   utils.GetEnv("MIGRATION_SOURCE_URL"),
		MigrationDatabaseURL: utils.GetEnv("MIGRATION_DATABASE_URL"),

		LoginTokenExpireSeconds:                    utils.ConvertStringToUint(utils.GetEnv("LOGIN_TOKEN_EXPIRE_SECONDS")),
		UserEmailVerificationCodeLen:               utils.ConvertStringToUint(utils.GetEnv("USER_EMAIL_VERIFICATION_CODE_LEN")),
		UserEmailVerificationCodeMaxTry:            utils.ConvertStringToUint(utils.GetEnv("USER_EMAIL_VERIFICATION_CODE_MAX_TRY")),
		UserEmailVerificationCodeExpiredSeconds:    utils.ConvertStringToUint(utils.GetEnv("USER_EMAIL_VERIFICATION_CODE_EXPIRED_SECONDS")),
		UserEmailVerificationCodeIssueLimitSeconds: utils.ConvertStringToUint(utils.GetEnv("USER_EMAIL_VERIFICATION_CODE_ISSUE_LIMIT_SECONDS")),

		VerifyEmailSubject: utils.GetEnv("VERIFY_EMAIL_SUBJECT"),
		VerifyEmailContent: utils.GetEnv("VERIFY_EMAIL_CONTENT"),
	}
}
