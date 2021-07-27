// Go to ${grpc-up-and-running}/samples/ch02/productinfo
// Optional: Execute protoc -I proto proto/product_info.proto --go_out=plugins=grpc:go/product_info
// Execute go get -v github.com/grpc-up-and-running/samples/ch02/productinfo/golang/product_info
// Execute go run go/client/main.go

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"time"

	pb "productinfo/ecommerce"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	address  = "localhost:50051"
	hostname = "localhost.modusign.co.kr"
	crtFile  = "cert/client.crt"
	keyFile  = "cert/client.key"
	caFile   = "cert/ca.crt" // CA의 공개키(인증서)
)

func main() {
	certificate, err := tls.LoadX509KeyPair(crtFile, keyFile) // 1. 인증서 생성
	if err != nil {
		log.Fatalf("could not load client key pair: %s", err)
	}

	certPool := x509.NewCertPool() // 2. 인증서 풀 생성
	ca, err := ioutil.ReadFile(caFile)
	if err != nil {
		log.Fatalf("could not read ca certificate: %s", err)
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok { // 3. ca 공개키(인증서)를 인증서 풀에 등록
		log.Fatalf("failed to append ca certs")
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{ // 4. TLS 활성화
			ServerName:   hostname, // 서버 인증서의 CN(Common Name)
			Certificates: []tls.Certificate{certificate},
			RootCAs:      certPool,
		})),
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
