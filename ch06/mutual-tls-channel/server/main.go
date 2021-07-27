package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"

	"github.com/gofrs/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"

	pb "productinfo/ecommerce"
)

const (
	port    = ":50051"
	crtFile = "cert/server.crt"
	keyFile = "cert/server.key"
	caFile  = "cert/ca.crt" // CA의 공개키(인증서)
)

// server is used to implement ecommerce/product_info.
type server struct {
	productMap map[string]*pb.Product
	pb.UnimplementedProductInfoServer
}

// func newServer() *routeGuideServer {
// 	s := &routeGuideServer{routeNotes: make(map[string][]*pb.RouteNote)}
// 	s.loadFeatures(*jsonDBFile)
// 	return s
// }

// AddProduct implements ecommerce.AddProduct
func (s *server) AddProduct(ctx context.Context,
	in *pb.Product) (*pb.ProductID, error) {
	out, err := uuid.NewV4()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error while generating Product ID", err)
	}
	in.Id = out.String()
	if s.productMap == nil {
		s.productMap = make(map[string]*pb.Product)
	}
	s.productMap[in.Id] = in
	log.Printf("Product %v : %v - Added.", in.Id, in.Name)
	return &pb.ProductID{Value: in.Id}, status.New(codes.OK, "").Err()
}

// GetProduct implements ecommerce.GetProduct
func (s *server) GetProduct(ctx context.Context, in *pb.ProductID) (*pb.Product, error) {
	product, exists := s.productMap[in.Value]
	if exists && product != nil {
		log.Printf("Product %v : %v - Retrieved.", product.Id, product.Name)
		return product, status.New(codes.OK, "").Err()
	}
	return nil, status.Errorf(codes.NotFound, "Product does not exist.", in.Value)
}

func main() {
	certificate, err := tls.LoadX509KeyPair(crtFile, keyFile) // 1. 인증서 생성
	if err != nil {
		log.Fatalf("failed to load Key pair: %s", err)
	}

	certPool := x509.NewCertPool() // 2. 인증서 풀 생성
	ca, err := ioutil.ReadFile(caFile)
	if err != nil {
		log.Fatalf("could not read ca certification: %s", err)
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok { // 3. ca 공개키(인증서)를 인증서 풀에 등록
		log.Fatalf("failed to append ca certification")
	}

	opts := []grpc.ServerOption{
		grpc.Creds(
			credentials.NewTLS(&tls.Config{ // 4. TLS 활성화
				ClientAuth:   tls.RequireAndVerifyClientCert,
				Certificates: []tls.Certificate{certificate},
				ClientCAs:    certPool,
			},
			)),
	}

	s := grpc.NewServer(opts...)

	pb.RegisterProductInfoServer(s, &server{})

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
