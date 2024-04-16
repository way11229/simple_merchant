package initial_process

import (
	"database/sql"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/way11229/simple_merchant/config"
	"github.com/way11229/simple_merchant/domain"
	"github.com/way11229/simple_merchant/service"

	mailer "github.com/way11229/simple_merchant/repo/mailer"
	mysql_store "github.com/way11229/simple_merchant/repo/mysql/store"
)

type RepositoryClientGroup struct {
	MysqlStore domain.MysqlStore
	Mailer     domain.MailerClient
}

func RunDbMigration(config *config.Config) {
	migration, err := migrate.New(config.MigrationSourceURL, config.MigrationDatabaseURL)
	if err != nil {
		log.Fatalf("cannot create new migrate instance, error: %v", err)
	}

	defer migration.Close()

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to run migrate up, error: %v", err)
	}

	log.Println("db migrated successfully")
}

func GetServiceManager(
	config *config.Config,
	mysqlConn *sql.DB,
) *domain.ServiceManager {
	return GetServiceManagerWithRepositoryClientGroup(
		config,
		GetRepositoryClientGroup(config, mysqlConn),
	)
}

func GetRepositoryClientGroup(
	config *config.Config,
	mysqlConn *sql.DB,
) *RepositoryClientGroup {
	return &RepositoryClientGroup{
		MysqlStore: mysql_store.NewStore(mysqlConn),
		Mailer:     mailer.NewMailer(),
	}
}

func GetServiceManagerWithRepositoryClientGroup(
	config *config.Config,
	repositoryClientGroup *RepositoryClientGroup,
) *domain.ServiceManager {
	return &domain.ServiceManager{
		UserService: service.NewUserService(
			repositoryClientGroup.MysqlStore,
			repositoryClientGroup.Mailer,
			time.Duration(config.LoginTokenExpireSeconds)*time.Second,
			config.UserEmailVerificationCodeLen,
			config.UserEmailVerificationCodeMaxTry,
			time.Duration(config.UserEmailVerificationCodeExpiredSeconds)*time.Second,
			time.Duration(config.UserEmailVerificationCodeIssueLimitSeconds)*time.Second,
			config.VerificationEmailSubject,
			config.VerificationEmailContent,
		),
	}
}
