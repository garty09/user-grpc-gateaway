package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	test "user_test/gen/go/proto"
	"user_test/internal"
	"user_test/internal/config"
	"user_test/internal/model"
	"user_test/internal/pub"
)

var flagConfig = flag.String("config", "./config/config.yml", "path to the config file")

func run() error {
	// load config
	cfg, err := config.Load(*flagConfig)
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	db, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		log.Fatalf("failed connect to postgres: %s", err)
	}

	cache := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
	defer func() {
		if err = db.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	p := pub.NewProducer(cfg.Brokers)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s := grpc.NewServer()
	test.RegisterUserServiceServer(s, &internal.RPCServer{
		Repo: model.NewRepository(db, cache, p),
	})
	lis, _ := net.Listen("tcp", "localhost:"+strconv.Itoa(cfg.GRPCPort))
	go s.Serve(lis)

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err = test.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "localhost:"+strconv.Itoa(cfg.GRPCPort), opts)
	if err != nil {
		return err
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(":"+strconv.Itoa(cfg.HTTPPort), mux)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
