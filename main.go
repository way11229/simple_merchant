package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/way11229/simple_merchant/config"
	"github.com/way11229/simple_merchant/domain"
	"github.com/way11229/simple_merchant/initial_process"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	handler_grpc "github.com/way11229/simple_merchant/handler/grpc"

	pb "github.com/way11229/simple_merchant/pb"
)

const (
	GRPC_PORT = 9000
	HTTP_PORT = 8080
)

func main() {
	envConfig := config.NewConfig()

	initial_process.RunMysqlMigration(envConfig)
	mysqlConn, err := sql.Open(
		envConfig.MysqlSqlDriverName,
		envConfig.MysqlSqlDataSourceName,
	)
	if err != nil {
		log.Fatalf("mysql connection error: %v", err)
	}

	defer mysqlConn.Close()

	serviceManager := initial_process.GetServiceManager(
		envConfig,
		mysqlConn,
	)

	newServer(serviceManager)
}

func newServer(
	serviceManager *domain.ServiceManager,
) {
	grpcHandler := handler_grpc.NewGrpcHandler(
		serviceManager.UserService,
		serviceManager.AuthService,
		serviceManager.ProductService,
	)

	go newGrpcGatewayServer(grpcHandler)

	newGrpcServer(grpcHandler)
}

func newGrpcGatewayServer(
	grpcHandler *handler_grpc.GrpcHandler,
) {
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames:   true,
			EmitUnpopulated: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := pb.RegisterSimpleMerchantHandlerServer(ctx, grpcMux, grpcHandler); err != nil {
		log.Fatalf("RegisterSimpleMerchantHandler error = %v", err)
	}

	handler := grpcHandler.PanicHandler(grpcMux)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", HTTP_PORT))
	if err != nil {
		log.Fatalf("net.Listen error = %v", err)
	}

	if err := http.Serve(lis, handler); err != nil {
		log.Fatalf("http.Serve error = %v", err)
	}
}

func newGrpcServer(
	grpcHandler *handler_grpc.GrpcHandler,
) {
	sigCh := make(chan os.Signal, 1)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", GRPC_PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_recovery.UnaryServerInterceptor(),
		),
	)
	pb.RegisterSimpleMerchantServer(
		grpcServer,
		grpcHandler,
	)

	reflection.Register(grpcServer)

	log.Println("grpc server start")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	grpcServer.GracefulStop()

	log.Println("grpc server graceful shutdown End")
}
