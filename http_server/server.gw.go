package http_server

import (
	"context"
	"fmt"
	"net/http"

	gw "static_db_wf/static_proto/pb"

	server "static_db_wf/static_proto/server"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

func Server_gw() {
	server.Server_grpc()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	//lw.Display_allfirst()
	//lw.Display_DBbydistence()
	httpEndpoint := "127.0.0.1:9090"
	grpcendpoint := "127.0.0.1:50052"
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	// HTTP转grpc
	err := gw.RegisterSearchHandlerFromEndpoint(ctx, mux, grpcendpoint, opts)
	if err != nil {
		fmt.Printf("Register handler err:%v\n", err)
	}
	//fmt.Print(1233)
	http.Handle("/", mux)
	httpServer := http.Server{Addr: httpEndpoint}

	fmt.Println("HTTP Listen on 9090")
	httpServer.ListenAndServe()

}
