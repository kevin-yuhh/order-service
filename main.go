package main

import (
	goFlag "flag"
	"fmt"
	"net"
	"time"

	"github.com/TRON-US/soter-order-service/charge"
	"github.com/TRON-US/soter-order-service/config"
	"github.com/TRON-US/soter-order-service/logger"
	"github.com/TRON-US/soter-order-service/model"
	"github.com/TRON-US/soter-order-service/service"
	orderPb "github.com/TRON-US/soter-proto/order-service"

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
	fmt.Println(`
   _____ ____  ________________     ____  ____  ____  __________ 
  / ___// __ \/_  __/ ____/ __ \   / __ \/ __ \/ __ \/ ____/ __ \
  \__ \/ / / / / / / __/ / /_/ /  / / / / /_/ / / / / __/ / /_/ /
 ___/ / /_/ / / / / /___/ _, _/  / /_/ / _, _/ /_/ / /___/ _, _/ 
/____/\____/ /_/ /_____/_/ |_|   \____/_/ |_/_____/_____/_/ |_|  
                                                                 `)
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

	// Get configure from database by env.
	strategyId, defaultTime, err := db.QueryConfig(conf.Env)
	if err != nil {
		panic(err)
	}

	// Get strategy from database by strategy id.
	strategy, err := db.QueryStrategyById(strategyId)
	if err != nil {
		panic(err)
	}

	server := &service.Server{
		DbConn: db,
		Fee:    charge.NewFeeCalculator(strategyId, strategy.Script),
		Config: conf,
		Time:   defaultTime,
	}

	server.Fee.Reload()

	go server.ClusterConsumer()

	// Register gRPC server.
	orderPb.RegisterOrderServiceServer(s, server)

	logger.Logger.Infof("[%v] server started, listening on port: [%v]", conf.Env, conf.Server.Port)
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
