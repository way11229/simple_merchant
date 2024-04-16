package config

import (
	"fmt"

	"github.com/way11229/simple_merchant/utils"
)

type Config struct {
	MysqlSqlDriverName        string `mapstructure:"MYSQL_SQL_DRIVER_NAME"`
	MysqlSqlDataSourceName    string `mapstructure:"MYSQL_SQL_DATA_SOURCE_NAME"`
	MysqlMigrationSourceURL   string `mapstructure:"MYSQL_MIGRATION_SOURCE_URL"`
	MysqlMigrationDatabaseURL string `mapstructure:"MYSQL_MIGRATION_DATABASE_URL"`

	LoginTokenExpireSeconds                    uint `mapstructure:"LOGIN_TOKEN_EXPIRE_SECONDS"`
	UserEmailVerificationCodeLen               uint `mapstructure:"USER_EMAIL_VERIFICATION_CODE_LEN"`
	UserEmailVerificationCodeMaxTry            uint `mapstructure:"USER_EMAIL_VERIFICATION_CODE_MAX_TRY"`
	UserEmailVerificationCodeExpiredSeconds    uint `mapstructure:"USER_EMAIL_VERIFICATION_CODE_EXPIRED_SECONDS"`
	UserEmailVerificationCodeIssueLimitSeconds uint `mapstructure:"USER_EMAIL_VERIFICATION_CODE_ISSUE_LIMIT_SECONDS"`

	VerificationEmailSubject string `mapstructure:"VERIFICATION_EMAIL_SUBJECT"`
	VerificationEmailContent string `mapstructure:"VERIFICATION_EMAIL_CONTENT"`

	SymmetricKey string `mapstructure:"SYMMETRIC_KEY"`
}

func NewConfig() *Config {
	return &Config{
		MysqlSqlDriverName:        utils.GetEnv("MYSQL_SQL_DRIVER_NAME"),
		MysqlSqlDataSourceName:    utils.GetEnv("MYSQL_SQL_DATA_SOURCE_NAME"),
		MysqlMigrationSourceURL:   utils.GetEnv("MYSQL_MIGRATION_SOURCE_URL"),
		MysqlMigrationDatabaseURL: utils.GetEnv("MYSQL_MIGRATION_DATABASE_URL"),

		LoginTokenExpireSeconds:                    convertStringToUintAndPanicIfError(utils.GetEnv("LOGIN_TOKEN_EXPIRE_SECONDS")),
		UserEmailVerificationCodeLen:               convertStringToUintAndPanicIfError(utils.GetEnv("USER_EMAIL_VERIFICATION_CODE_LEN")),
		UserEmailVerificationCodeMaxTry:            convertStringToUintAndPanicIfError(utils.GetEnv("USER_EMAIL_VERIFICATION_CODE_MAX_TRY")),
		UserEmailVerificationCodeExpiredSeconds:    convertStringToUintAndPanicIfError(utils.GetEnv("USER_EMAIL_VERIFICATION_CODE_EXPIRED_SECONDS")),
		UserEmailVerificationCodeIssueLimitSeconds: convertStringToUintAndPanicIfError(utils.GetEnv("USER_EMAIL_VERIFICATION_CODE_ISSUE_LIMIT_SECONDS")),

		VerificationEmailSubject: utils.GetEnv("VERIFICATION_EMAIL_SUBJECT"),
		VerificationEmailContent: utils.GetEnv("VERIFICATION_EMAIL_CONTENT"),

		SymmetricKey: utils.GetEnv("SYMMETRIC_KEY"),
	}
}

func convertStringToUintAndPanicIfError(str string) uint {
	rtn, err := utils.ConvertStringToUint(str)
	if err != nil {
		panic(fmt.Sprintf("ConvertStringToUint error = %v", err))
	}

	return rtn
}
