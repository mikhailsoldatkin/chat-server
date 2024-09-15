package app

import (
	"context"
	"log"
	"net"
	"os"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/mikhailsoldatkin/chat-server/internal/interceptor"
	"github.com/mikhailsoldatkin/chat-server/internal/logger"
	"github.com/mikhailsoldatkin/chat-server/internal/tracing"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	"github.com/mikhailsoldatkin/chat-server/internal/config"
	pb "github.com/mikhailsoldatkin/chat-server/pkg/chat_v1"
	"github.com/mikhailsoldatkin/platform_common/pkg/closer"
)

// App represents the application with its dependencies and GRPC server.
type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

// NewApp initializes a new App instance with the given context and sets up the necessary dependencies.
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Run starts the GRPC server and handles graceful shutdown.
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.runGRPCServer()
}

// initDeps initializes the dependencies required by the App.
func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
		a.initLogger,
		a.initTracing,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	_, err := config.Load()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	creds, err := credentials.NewServerTLSFromFile("cert/service.pem", "cert/service.key")
	if err != nil {
		log.Fatalf("failed to load TLS credentials from files: %v", err)
	}

	a.grpcServer = grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				interceptor.ServerTracingInterceptor,
				interceptor.AuthInterceptor(a.serviceProvider.AuthClient()),
			),
		),
	)

	reflection.Register(a.grpcServer)

	pb.RegisterChatV1Server(a.grpcServer, a.serviceProvider.ChatImplementation(ctx))

	return nil
}

func (a *App) runGRPCServer() error {
	lis, err := net.Listen("tcp", a.serviceProvider.config.GRPC.Address)
	if err != nil {
		return err
	}

	log.Printf("gRPC server is running on %d", a.serviceProvider.config.GRPC.Port)

	err = a.grpcServer.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}

// initLogger initializes the app logger.
func (a *App) initLogger(_ context.Context) error {
	var level zapcore.Level
	if err := level.Set(a.serviceProvider.config.Logger.Level); err != nil {
		return err
	}

	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   a.serviceProvider.config.Logger.Filename,
		MaxSize:    a.serviceProvider.config.Logger.MaxSizeMB,
		MaxBackups: a.serviceProvider.config.Logger.MaxBackups,
		MaxAge:     a.serviceProvider.config.Logger.MaxAgeDays,
	})

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006/01/02 15:04:05")

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	developmentCfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006/01/02 15:04:05")

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	logger.Init(zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, zap.NewAtomicLevelAt(level)),
		zapcore.NewCore(fileEncoder, file, zap.NewAtomicLevelAt(level)),
	))

	return nil
}

// initTracing initializes the tracing for the application by setting up the Jaeger tracer.
func (a *App) initTracing(_ context.Context) error {
	tracing.Init(logger.Logger(), "chat-server", a.serviceProvider.config.Jaeger.Address)
	return nil
}
