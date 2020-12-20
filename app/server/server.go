package server

import (
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/joho/godotenv"
	handler "github.com/softcorp-io/hqs_department_service/handler"
	proto "github.com/softcorp-io/hqs_proto/go_hqs/hqs_department_service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Init - initialize .env variables.
func Init(zapLog *zap.Logger) {
	if err := godotenv.Load("hqs.env"); err != nil {
		zapLog.Error(fmt.Sprintf("Could not load hqs.env with err %v", err))
	}
}

// Run - runs a go microservice. Uses zap for logging and a waitGroup for async testing.
func Run(zapLog *zap.Logger, wg *sync.WaitGroup) {
	// use above to create handler
	handle := handler.NewHandler(zapLog)

	// create the service and run the service
	port, ok := os.LookupEnv("SERVICE_PORT")
	if !ok {
		zapLog.Fatal("Could not get service port")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		zapLog.Fatal(fmt.Sprintf("Failed to listen with err %v", err))
	}
	defer lis.Close()

	zapLog.Info(fmt.Sprintf("Service running on port: %s", port))

	// setup grpc

	grpcServer := grpc.NewServer()

	// register handler
	proto.RegisterDepartmentServiceServer(grpcServer, handle)

	// run the server
	wg.Done()
	if err := grpcServer.Serve(lis); err != nil {
		zapLog.Fatal(fmt.Sprintf("Failed to serve with err %v", err))
	}
}
