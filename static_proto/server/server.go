package server

import (
	"context"
	"fmt"
	"log"
	"net"
	tools "static_db_wf/api_ips/tools"
	lw "static_db_wf/lib_worker"
	pb "static_db_wf/static_proto/pb"

	"google.golang.org/grpc"
)

const (
	port = ":50052"
)

type server struct {
	pb.SearchServer
}

// server is used to implement helloworld.GreeterServer.

// SayHello implements helloworld.GreeterServer
func (s server) Add(ctx context.Context, in *pb.AddRequest) (*pb.AddReply, error) {
	features := in.GetFeatures()
	var request lw.Request
	for i := 0; i < len(features); i++ {
		b, err := tools.Base64Decode([]byte(features[i].Blob))
		if err != nil {
			fmt.Print(err)
			return &pb.AddReply{}, err
		}
		request.Vectors.Vector = append(request.Vectors.Vector, tools.DecodeFloat32(b))
	}
	//fmt.Print(len(request.Vectors.Vector))
	res := lw.Add(request)
	return &pb.AddReply{Status: res}, nil
}
func (s server) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteReply, error) {

	var request lw.Request
	request.I_delete = float32(in.GetRow())
	request.J_delete = float32(in.GetCol())
	//fmt.Print(len(request.Vectors.Vector))
	res := lw.Delete(request)
	return &pb.DeleteReply{Status: res}, nil
}
func (s server) SearchInDB(ctx context.Context, in *pb.SearchRequest) (*pb.SearchReply, error) {

	var reply pb.SearchReply
	var req lw.Request
	var a = in.GetFeatures()
	for i := 0; i < len(a); i++ {
		//fmt.Print(555)
		b, err := tools.Base64Decode([]byte(a[i].Blob))
		if err != nil {
			fmt.Print(err)
			return &pb.SearchReply{}, err
		}
		req.Vectors.Vector = append(req.Vectors.Vector, tools.DecodeFloat32(b))
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

	var featuresgroup []*pb.Features
	for i := 0; i < len(result.VectorGroup); i++ {
		var features pb.Features
		for j := 0; j < len(result.VectorGroup[i]); j++ {
			var vector pb.Vector
			for k := 0; k < len(result.VectorGroup[i][j]); k++ {
				vector.Vector = append(vector.Vector, result.VectorGroup[i][j][k])
			}
			b := tools.EncodeFloat32(vector.Vector)
			b = tools.Base64Encode(b)
			features.Features = append(features.Features, &pb.Feature{Blob: string(b)})
		}
		featuresgroup = append(featuresgroup, &features)
	}
	reply.Featuresgroup = featuresgroup

	return &reply, nil
}
func Server_grpc() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSearchServer(s, &server{})
	fmt.Println("grpc listening on:", port)
	go s.Serve(lis)
}
