package main

import (
	"context"
	"encoding/base64"
	"log"
	"time"

	pb "productinfo/ecommerce"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	address  = "localhost:50051"
	hostname = "localhost.modusign.co.kr"
	crtFile  = "cert/server.crt"
)

/////////////  커스텀 자격증명
// 1. 인증에 사용할 구조체 정의
type basicAuth struct {
	username string
	password string
}

// 2. 인증 정보를 요청 메타데이터로 변환
func (b basicAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	auth := b.username + ":" + b.password
	enc := base64.StdEncoding.EncodeToString([]byte(auth))
	return map[string]string{
		"authorization": "Basic " + enc,
	}, nil
}

// 3. 채널 보안이 필요한지 여부 지정
func (b basicAuth) RequireTransportSecurity() bool {
	return true
}

/////////////////////////////

func main() {
	creds, err := credentials.NewClientTLSFromFile(crtFile, hostname)
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	// 1. 사용자 인증 정보 설정
	auth := basicAuth{
		username: "admin",
		password: "admin1",
	}

	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(auth), // 인증 정보 전달
		grpc.WithTransportCredentials(creds),
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProductInfoClient(conn)

	// Contact the server and print out its response.
	name := "Apple iPhone 11"
	description := "Meet Apple iPhone 11. All-new dual-camera system with Ultra Wide and Night mode."
	price := float32(699.00)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.AddProduct(ctx, &pb.Product{Name: name, Description: description, Price: price})
	if err != nil {
		log.Fatalf("Could not add product: %v", err)
	}
	log.Printf("Product ID: %s added successfully", r.Value)

	product, err := c.GetProduct(ctx, &pb.ProductID{Value: r.Value})
	if err != nil {
		log.Fatalf("Could not get product: %v", err)
	}
	log.Printf("Product: %v", product.String())
}
