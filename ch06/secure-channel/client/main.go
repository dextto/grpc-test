// Go to ${grpc-up-and-running}/samples/ch02/productinfo
// Optional: Execute protoc -I proto proto/product_info.proto --go_out=plugins=grpc:go/product_info
// Execute go get -v github.com/grpc-up-and-running/samples/ch02/productinfo/golang/product_info
// Execute go run go/client/main.go

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
	hostname = "localhost"
	crtFile  = "cert/server.crt"
)

func main() {
	creds, err := credentials.NewClientTLSFromFile(crtFile, hostname)
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}
	opts := []grpc.DialOption{
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
