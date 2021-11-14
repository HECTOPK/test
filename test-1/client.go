package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"os"
	"test-1/api"
	"test-1/auth"

	"time"
)

const (
	clientCert = "cert/client-cert.pem"
	clientKey  = "cert/client-key.pem"
	rootCert   = "cert/ca-cert.pem"
)

func main() {
	if len(os.Args) != 4 {
		panic("bad number of args (should be 4)")
	}
	mode, host, phone := os.Args[1], os.Args[2], os.Args[3]
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	conn, err := grpc.DialContext(ctx, host, grpc.WithTransportCredentials(getClientCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := api.NewApiServiceClient(conn)
	req := &api.GetCodeRequest{Phone: phone}
	var resp *api.GetCodeResponse
	ctx = metadata.AppendToOutgoingContext(ctx, "auth-token", auth.GetToken())
	switch mode {
	case "unary":
		resp, err = client.GetCode(ctx, req)
	case "stream":
		var stream api.ApiService_GetcodeStreamClient
		stream, err = client.GetcodeStream(ctx)
		if err != nil {
			break
		}
		err = stream.Send(req)
		if err != nil {
			break
		}
		resp, err = stream.Recv()
	default:
		panic("bad request type")
	}
	if err != nil {
		panic(err)
	}
	if resp.Code == "" {
		fmt.Println("not found")
		return
	}
	fmt.Println(resp.Code)
}

func getClientCredentials() credentials.TransportCredentials {
	ca, err := os.ReadFile(rootCert)
	if err != nil {
		panic(err)
	}
	pool := x509.NewCertPool()
	if pool.AppendCertsFromPEM(ca) != true {
		panic("error read root certificate")
	}
	cert, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		panic(err)
	}
	conf := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}
	return credentials.NewTLS(conf)
}
