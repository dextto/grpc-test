package main

import (
	"context"
	"log"
	"time"

	pb "productinfo/ecommerce"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	address  = "localhost:50051"
	hostname = "localhost.modusign.co.kr" // 공개키가 허용하는 CN(Common Name)
	crtFile  = "cert/server.crt"          // server's self-signed public key (certificate)
)

func main() {
	creds, err := credentials.NewClientTLSFromFile(crtFile, hostname) // 1. 서버의 공개키로 인증서 생성(??)
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds), // 2. TLS 연결 설정
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
