package main

import (
	goFlag "flag"
	"fmt"
	"net"
	"time"

	"github.com/TRON-US/soter-order-service/config"
	"github.com/TRON-US/soter-order-service/logger"
	"github.com/TRON-US/soter-order-service/model"
	orderPb "github.com/TRON-US/soter-order-service/proto"
	"github.com/TRON-US/soter-order-service/service"

	"github.com/prometheus/common/log"
	flag "github.com/spf13/pflag"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var (
	cfgDir  = flag.StringP("dir", "d", "./", "server config file dir")
	cfgName = flag.StringP("name", "n", "config", "server config file filename")
)

func main() {
	defer recoverPanic()
	// Load config.
	flag.CommandLine.AddGoFlagSet(goFlag.CommandLine)
	flag.Parse()

	// Set the time zone to UTC+8.
	time.Local = time.FixedZone("CST", 3600*8)

	// Get config from config.yaml.
	conf, err := config.NewConfiguration(*cfgName, *cfgDir)
	if err != nil {
		panic(err)
	}

	// Init logger config.
	err = logger.InitLogger(conf.Logger.Output, zapcore.Level(conf.Logger.Level))
	if err != nil {
		panic(err)
	}

	// New tcp listener.
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", conf.Server.Host, conf.Server.Port))
	if err != nil {
		panic(err)
	}

	// New gRPC server.
	s := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
		}))

	// Init database connection.
	db, err := model.NewDatabase(conf)
	if err != nil {
		panic(err)
	}

	// Register gRPC server.
	orderPb.RegisterOrderServiceServer(s, &service.Server{DbConn: db})

	logger.Logger.Info(fmt.Sprintf("Server started, listening on port %v...", conf.Server.Port))
	if err = s.Serve(lis); err != nil {
		panic(err)
	}
}

// Global panic recover.
func recoverPanic() {
	if rec := recover(); rec != nil {
		log.Error(fmt.Sprintf("Failed to generate server, reasons: [%v]", rec))
	}
}
