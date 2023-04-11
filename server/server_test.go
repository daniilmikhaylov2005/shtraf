package main

import (
	"context"
	"testing"

	pb "github.com/daniilmikhaylov2005/shtraf/api"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const hostPort = "localhost:8800"

func TestGetInfo(t *testing.T) {
	conn, err := grpc.DialContext(
		context.Background(),
		hostPort,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Errorf("failed to get conn: %v", err)
	}
	defer conn.Close()

	client := pb.NewRusProfileClient(conn)

	grpcRes, err := client.GetInfo(context.Background(), &pb.GetInfoRequest{Inn: "7705878037"})
	if err != nil {
		t.Errorf("failed to use client: %v", err)
	}

	assert.Equal(t, "7705878037", grpcRes.GetInn(), "The two inn should be the same")
}
