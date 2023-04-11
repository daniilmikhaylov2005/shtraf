package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	pb "github.com/daniilmikhaylov2005/shtraf/api"
	"google.golang.org/grpc"
)

func main() {
	inn := flag.String("inn", "", "ИНН для поиска компании")
	flag.Parse()

	conn, _ := grpc.Dial("localhost:8800", grpc.WithInsecure(), grpc.WithBlock())
	defer conn.Close()

	client := pb.NewRusProfileClient(conn)

	in, err := client.GetInfo(context.Background(), &pb.GetInfoRequest{
		Inn: *inn,
	})
	if err != nil {
		log.Fatalf("[err] %v", err)
	}

	fmt.Printf(`
		ИНН: %s,
		КПП: %s,
		Название компании: %s,
		ФИО Директора: %s.
	`, in.GetInn(), in.GetKpp(), in.GetTitle(), in.GetFullnameDirector())

}
