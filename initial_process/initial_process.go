package initial_process

import (
	"database/sql"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/way11229/simple_merchant/config"
	"github.com/way11229/simple_merchant/domain"
	"github.com/way11229/simple_merchant/pkg/auth_token_maker"
	"github.com/way11229/simple_merchant/service"

	mailer "github.com/way11229/simple_merchant/repo/mailer"
	mysql_store "github.com/way11229/simple_merchant/repo/mysql/store"
	redis "github.com/way11229/simple_merchant/repo/redis"
)

type RepositoryClientGroup struct {
	MysqlStore     domain.MysqlStore
	Mailer         domain.MailerClient
	AuthTokenMaker auth_token_maker.AuthTokenMaker
	RedisClient    domain.RedisClient
}

func RunMysqlMigration(config *config.Config) {
	migration, err := migrate.New(config.MysqlMigrationSourceURL, config.MysqlMigrationDatabaseURL)
	if err != nil {
		log.Fatalf("cannot create new mysql migrate instance, error: %v", err)
	}

	defer migration.Close()

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to run mysql migrate up, error: %v", err)
	}

	log.Println("mysql migrated successfully")
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
	pasetoMaker, err := auth_token_maker.NewPasetoMaker(config.SymmetricKey)
	if err != nil {
		log.Fatalf("NewPasetoMaker error = %v", err)
	}

	return &RepositoryClientGroup{
		MysqlStore:     mysql_store.NewStore(mysqlConn),
		Mailer:         mailer.NewMailer(),
		AuthTokenMaker: pasetoMaker,
		RedisClient:    redis.NewRedisClient(config.RedisAddr, config.RedisPwd),
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
		AuthService: service.NewAuthService(
			repositoryClientGroup.MysqlStore,
			repositoryClientGroup.AuthTokenMaker,
			repositoryClientGroup.RedisClient,
			time.Duration(config.LoginTokenExpireSeconds)*time.Second,
			time.Duration(config.LoginTokenCacheExpireSeconds)*time.Second,
		),
		ProductService: service.NewProductService(
			repositoryClientGroup.MysqlStore,
			repositoryClientGroup.RedisClient,
			time.Duration(config.RecommendedProductCacheExpiredSeconds)*time.Second,
		),
	}
}
