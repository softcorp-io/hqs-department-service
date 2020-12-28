package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	database "github.com/softcorp-io/hqs_department_service/database"
	handler "github.com/softcorp-io/hqs_department_service/handler"
	repository "github.com/softcorp-io/hqs_department_service/repository"
	proto "github.com/softcorp-io/hqs_proto/go_hqs/hqs_department_service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type collectionEnv struct {
	departmentCollection string
}

// Init - initialize .env variables.
func Init(zapLog *zap.Logger) {
	if err := godotenv.Load("hqs.env"); err != nil {
		zapLog.Error(fmt.Sprintf("Could not load hqs.env with err %v", err))
	}
}

func loadCollections() (collectionEnv, error) {
	departmentCollection, ok := os.LookupEnv("MONGO_DB_DEPARTMENT_COLLECTION")
	if !ok {
		return collectionEnv{}, errors.New("Required MONGO_DB_DEPARTMENT_COLLECTION")
	}
	return collectionEnv{departmentCollection}, nil
}

// Run - runs a go microservice. Uses zap for logging and a waitGroup for async testing.
func Run(zapLog *zap.Logger, wg *sync.WaitGroup) {
	// creates a database connection and closes it when done
	mongoenv, err := database.GetMongoEnv()
	if err != nil {
		zapLog.Fatal(fmt.Sprintf("Could not set up mongo env with err %v", err))
	}
	// build uri for mongodb
	mongouri := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority", mongoenv.User, mongoenv.Password, mongoenv.Host, mongoenv.DBname)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	mongo, err := database.NewMongoDatabase(ctx, zapLog, mongouri)
	if err != nil {
		zapLog.Fatal(fmt.Sprintf("Could not make connection to DB with err %v", err))
	}

	defer mongo.Disconnect(context.Background())

	database := mongo.Database(mongoenv.DBname)

	collections, err := loadCollections()
	if err != nil {
		zapLog.Fatal(fmt.Sprintf("Could not load collections with err: %v", err))
	}

	departmentCollection := database.Collection(collections.departmentCollection)

	// setting up repository
	repo := repository.NewRepository(departmentCollection)

	// use above to create handler
	handle := handler.NewHandler(zapLog, repo)

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
