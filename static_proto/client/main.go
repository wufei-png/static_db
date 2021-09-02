package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "static_db_wf/static_proto/pb"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func test() string {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewSearchClient(conn)

	// Contact the server and print out its response.

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var request pb.SearchRequest
	request.SearchType = 2
	request.Topk = 2
	var vector pb.Vector
	var vectors pb.Vectors
	for i := 0; i < 4; i++ {
		vector.Vector = append(vector.Vector, float32(i))
	}
	vectors.Vectors = append(vectors.Vectors, &vector)
	// fmt.Print("vector:", vector.Vector)
	// fmt.Print("vectors:", vectors.Vectors[0].Vector)
	// var vectors pb.Vectors
	// vectors.Vectors = append(vectors.Vectors, &vector)
	request.Vectors = &vectors
	// fmt.Print(request.Vectors.Vectors[0].Vector)
	r, err := c.SearchInDB(ctx, &request)
	if err != nil {
		log.Fatalf("could not search: %v", err)
	}

	fmt.Println("distancetopk:", r.GetDistancetopk())
	fmt.Println("vectorsgroup:", r.GetVectorsgroup())
	return "成功！" //可以转换为result的格式
}
func main() {
	fmt.Print(test())
}
