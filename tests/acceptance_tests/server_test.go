package acceptance_tests

import (
	"context"
	"database/sql"
	"log"
	"net"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	handler_grpc "github.com/way11229/simple_merchant/handler/grpc"
	"github.com/way11229/simple_merchant/initial_process"
	pb "github.com/way11229/simple_merchant/pb"
)

func server(ctx context.Context) (pb.SimpleMerchantClient, func()) {
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	testConfig := getTestConfigFromEnv()
	accessToken = testConfig.AccessToken

	initial_process.RunMysqlMigration(&testConfig.Config)
	mysqlConn, err := sql.Open(testConfig.MysqlSqlDriverName, testConfig.MysqlSqlDataSourceName)
	if err != nil {
		log.Fatalf("mysql connection error: %v", err)
	}

	baseServer := grpc.NewServer()
	pb.RegisterSimpleMerchantServer(baseServer, newGrpcHandler(
		testConfig,
		mysqlConn,
	))
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Printf("error serving server: %v", err)
		}
	}()

	conn, err := grpc.DialContext(
		ctx,
		"",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("error connecting to server: %v", err)
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}

		baseServer.Stop()
		mysqlConn.Close()
	}

	return pb.NewSimpleMerchantClient(conn), closer
}

func newGrpcHandler(
	config *testConfig,
	mysqlConn *sql.DB,
) *handler_grpc.GrpcHandler {
	repositoryClientGroup := initial_process.GetRepositoryClientGroup(&config.Config, mysqlConn)

	// just for test
	repositoryClientGroup.Mailer = newTestMailerClient()

	serviceManager := initial_process.GetServiceManagerWithRepositoryClientGroup(&config.Config, repositoryClientGroup)

	return handler_grpc.NewGrpcHandler(
		serviceManager.UserService,
		serviceManager.AuthService,
	)
}

func getTestConfigFromEnv() *testConfig {
	config := testConfig{}

	viper.AddConfigPath(".")
	viper.SetConfigName("test")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&config.Config); err != nil {
		panic(err)
	}

	return &config
}
