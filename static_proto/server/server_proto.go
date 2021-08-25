package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	lw "static_db_wf/lib_worker"
	pb "static_db_wf/static_proto/pb"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct {
	pb.SearchServer
}

// server is used to implement helloworld.GreeterServer.

// SayHello implements helloworld.GreeterServer
func (s server) SearchInDB(ctx context.Context, in *pb.SearchRequest) (*pb.SearchReply, error) {

	var reply pb.SearchReply
	var req lw.Request
	var vectors = in.GetVectors()
	for i := 0; i < len(vectors.Vectors); i++ {

		fmt.Print(555)
		if len(vectors.Vectors[i].Vector) != 4 {
			fmt.Print(12343)
			return nil, errors.New("传入长度不对")
		} //vectors.Vectors[i].Vector

		// vectors.Vectors = append(vectors.Vectors, vectors.Vectors[i])

		req.Vectors.Vector = append(req.Vectors.Vector, vectors.Vectors[i].Vector)
	}
	req.Topk = int(in.GetTopk())
	req.Request_type = int(in.GetSearchType())
	result := lw.Search(req)
	var vectors_dis pb.Vectors
	// fmt.Print(vectors_dis)
	for i := 0; i < len(result.TopkDistance); i++ {
		var vector_dis pb.Vector
		for j := 0; j < len(result.TopkDistance[i]); j++ {
			vector_dis.Vector = append(vector_dis.Vector, result.TopkDistance[i][j])
		}
		vectors_dis.Vectors = append(vectors_dis.Vectors, &vector_dis)
	}
	reply.Distancetopk = &vectors_dis

	// for i := 0; i < len(result.VectorGroup); i++ {
	// 	reply.Vectorsgroup.Vectorgroup = append(reply.Vectorsgroup.Vectorgroup, &pb.Vectors{})
	// 	for j := 0; j < len(result.VectorGroup[i]); j++ {

	// 		reply.Vectorsgroup.Vectorgroup[i].Vectors = append(reply.Vectorsgroup.Vectorgroup[i].Vectors, &pb.Vector{})

	// 		reply.Vectorsgroup.Vectorgroup[i].Vectors[j].Vector = result.VectorGroup[i][j]
	// 	}
	// }

	var group pb.VectorsGroup
	for i := 0; i < len(result.VectorGroup); i++ {
		var vectors_gr pb.Vectors
		for j := 0; j < len(result.VectorGroup[i]); j++ {
			var vector_gr pb.Vector
			for k := 0; k < len(result.VectorGroup[i][j]); k++ {
				vector_gr.Vector = append(vector_gr.Vector, result.VectorGroup[i][j][k])
			}
			vectors_gr.Vectors = append(vectors_gr.Vectors, &vector_gr)
		}
		group.Vectorgroup = append(group.Vectorgroup, &vectors_gr)
	}
	reply.Vectorsgroup = &group

	return &reply, nil
}
func Server_grpc() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSearchServer(s, &server{})
	fmt.Print(666)
	s.Serve(lis)

}
