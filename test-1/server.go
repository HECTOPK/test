package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"net"
	"os"
	"test-1/api"
	"test-1/auth"
	"test-1/storage"
	"time"
)

const (
	data       = "test_data.csv"
	serverCert = "cert/server-cert.pem"
	serverKey  = "cert/server-key.pem"
	rootCert   = "cert/ca-cert.pem"
)

func main() {
	st := storage.New(data)
	serv := &server{storage: st}
	s := grpc.NewServer(grpc.Creds(getServerCredentials()),
		grpc.UnaryInterceptor(serv.unaryInterceptor()),
		grpc.StreamInterceptor(serv.streamInterceptor()))
	api.RegisterApiServiceServer(s, serv)
	lis, err := net.Listen("tcp", ":6000")
	if err != nil {
		panic(err)
	}
	fmt.Println("start grpc server on port 6000")
	err = s.Serve(lis)
	if err != nil {
		panic(err)
	}
}

func getServerCredentials() credentials.TransportCredentials {
	cert, err := tls.LoadX509KeyPair(serverCert, serverKey)
	if err != nil {
		panic(err)
	}
	root, err := os.ReadFile(rootCert)
	if err != nil {
		panic(err)
	}
	pool := x509.NewCertPool()
	if pool.AppendCertsFromPEM(root) != true {
		panic("error read root certificate")
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    pool,
	}
	return credentials.NewTLS(config)
}

type server struct {
	api.UnimplementedApiServiceServer
	storage *storage.Storage
}

func (s *server) GetCode(ctx context.Context, request *api.GetCodeRequest) (*api.GetCodeResponse, error) {
	t := time.Now()
	code := s.storage.Find(request.Phone)
	fmt.Printf("code for %s is %s (%v)\n", request.Phone, code, time.Since(t))
	return &api.GetCodeResponse{
		Code: code,
	}, nil
}

func (s *server) GetcodeStream(stream api.ApiService_GetcodeStreamServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		t := time.Now()
		code := s.storage.Find(req.Phone)
		err = stream.Send(&api.GetCodeResponse{Code: code})
		if err != nil {
			return err
		}

		fmt.Printf("code for %s is %s (%v)\n", req.Phone, code, time.Since(t))
	}
}

func (s *server) unaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		err = s.auth(ctx)
		if err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

func (s *server) streamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := s.auth(ss.Context())
		if err != nil {
			return err
		}
		return handler(srv, ss)
	}
}

func (s *server) auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("no metadata in context")
	}
	strs := md.Get("auth-token")
	if len(strs) < 1 {
		return status.Error(codes.Unauthenticated, "no token")
	}
	if !auth.CheckToken(strs[0]) {
		return status.Error(codes.Unauthenticated, "bad token")
	}
	return nil
}
