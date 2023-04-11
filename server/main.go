package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"

	pb "github.com/daniilmikhaylov2005/shtraf/api"
	"github.com/daniilmikhaylov2005/shtraf/server/models"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opesun/goquery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Server struct {
	log *zap.Logger
	pb.UnimplementedRusProfileServer
}

func main() {
	port := flag.Int("port", 8800, "port for a service")
	host := flag.String("host", "localhost", "host for a service")
	flag.Parse()

	hostPort := fmt.Sprintf("%s:%d", *host, *port)

	logger := zap.NewExample()
	defer logger.Sync()

	listener, err := net.Listen("tcp", hostPort)
	if err != nil {
		logger.Error("failed to up the listener",
			zap.String("host", *host),
			zap.Int("port", *port),
		)
		panic(err)
	}

	s := grpc.NewServer()
	pb.RegisterRusProfileServer(s, &Server{log: logger})

	logger.Info("Server started.",
		zap.String("url", hostPort),
	)

	go func() {
		if err := s.Serve(listener); err != nil {
			logger.Error("failed to up the grpc server",
				zap.String("host", *host),
				zap.Int("port", *port),
			)
			panic(err)
		}
	}()

	// rest handlers
	conn, err := grpc.DialContext(
		context.Background(),
		hostPort,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Info("failed to start connect to grpc server")
		panic(err)
	}

	gwMux := runtime.NewServeMux()

	err = pb.RegisterRusProfileHandler(context.Background(), gwMux, conn)
	if err != nil {
		logger.Info("failed to register handler")
		panic(err)
	}

	gPort := *port + 1
	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", gPort),
		Handler: gwMux,
	}

	if err := gwServer.ListenAndServe(); err != nil {
		panic(err)
	}

}

func (s *Server) GetInfo(ctx context.Context, in *pb.GetInfoRequest) (*pb.Company, error) {
	urlToData := fmt.Sprint("https://www.rusprofile.ru/search?query=", in.GetInn())
	result, err := goquery.ParseUrl(urlToData)
	if err != nil {
		s.log.Error("failed to get data", zap.String("error string", err.Error()))
		return nil, status.Error(codes.Unknown, "failed to get data")
	}
	divDirector := result.Find("div.company-row.hidden-parent")
	director := divDirector.Find("span.company-info__text").First().Text()
	info := models.Company{
		Inn:      result.Find("#clip_inn").Text(),
		Kpp:      result.Find("#clip_kpp").Text(),
		Title:    result.Find(".company-name").Text(),
		Director: director,
	}
	return &pb.Company{
		Inn:              info.Inn,
		Kpp:              info.Kpp,
		Title:            info.Title,
		FullnameDirector: director,
	}, nil
}
