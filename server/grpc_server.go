package server

import (
	"context"
	_ "embed"
	"time"

	"github.com/cetnfurkan/core/config"
	"github.com/labstack/gommon/log"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	rkentry "github.com/rookie-ninja/rk-entry/v2/entry"
	rkgrpc "github.com/rookie-ninja/rk-grpc/v2/boot"
	"google.golang.org/grpc"
)

var (
	_ grpc.UnaryServerInterceptor = UnaryLogHandler
)

type (
	GrpcServer struct {
		cfg   *config.Server
		boot  []byte
		entry *rkgrpc.GrpcEntry
	}

	grpcServerOption func(*GrpcServer) error
)

func WithRegisterGrpcFunc(f ...rkgrpc.GrpcRegFunc) grpcServerOption {
	return func(server *GrpcServer) error {
		server.entry.AddRegFuncGrpc(f...)
		return nil
	}
}

func WithRegisterGrpcGatewayFunc(f ...rkgrpc.GwRegFunc) grpcServerOption {
	return func(server *GrpcServer) error {
		server.entry.AddRegFuncGw(f...)
		return nil
	}
}

// NewGRPCServer creates a new gRPC server instance.
//
// It takes a config instance, a database instance and a cache instance
// and returns a new server interface instance.
//
// It will panic if it fails to create a new gRPC server instance.
func NewGRPCServer(cfg *config.Server, name string, boot []byte, opts ...grpcServerOption) Server {
	server := &GrpcServer{
		boot: boot,
		cfg:  cfg,
	}

	// Bootstrap basic entries from boot config.
	rkentry.BootstrapBuiltInEntryFromYAML(server.boot)
	rkentry.BootstrapPluginEntryFromYAML(server.boot)

	// Bootstrap grpc entry from boot config
	res := rkgrpc.RegisterGrpcEntryYAML(server.boot)

	if res[name] == nil {
		log.Fatal("Grpc entry not found in boot config.")
	}

	// Get GrpcEntry
	server.entry = res[name].(*rkgrpc.GrpcEntry)

	for _, opt := range opts {
		if err := opt(server); err != nil {
			log.Fatal("Unable to apply grpc server option: ", err)
		}
	}

	server.entry.ServerOpts = append(
		server.entry.ServerOpts,
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			UnaryLogHandler,
		)),
	)

	return server
}

func (server *GrpcServer) Start() {
	go server.start()
}

func (server *GrpcServer) start() {
	// Bootstrap
	server.entry.Bootstrap(context.Background())

	// Wait for shutdown signal
	rkentry.GlobalAppCtx.WaitForShutdownSig()

	// Wait for shutdown sig
	server.entry.Interrupt(context.Background())
}

func (server *GrpcServer) Stop() error {
	return nil
}

func UnaryLogHandler(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

	start := time.Now()

	resp, err = handler(ctx, req)

	end := time.Now()
	latency := end.Sub(start)

	tags := grpc_ctxtags.Extract(ctx).Values()
	ip, _ := tags["ip"].(string)

	if err != nil {
		log.Errorf("[ERROR] %s took: %s ip: %s err: %s req: %s", info.FullMethod, latency, ip, err, req)
	} else {
		log.Infof("[INFO] %s took: %s ip: %s req: %s", info.FullMethod, latency, ip, req)
	}

	return
}
