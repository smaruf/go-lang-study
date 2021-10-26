package proto

import (
	"context"
	. "core-service/config"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type WalletAccountClient struct {
	Client WalletServiceClient
	Conn   *grpc.ClientConn
}
type WalletTransferClient struct {
	Client WalletServiceClient
	Conn   *grpc.ClientConn
}

type WalletSearchClient struct {
	Client WalletServiceClient
	Conn   *grpc.ClientConn
}

func NewWalletSearchClient() *WalletSearchClient {
	fmt.Println("Established connection")
	conn, err := grpc.DialContext(context.Background(),
		HostConfig.RemoteSearchHost,
		grpc.WithReturnConnectionError(),
		grpc.WithBlock(),
		grpc.WithNoProxy(),
		grpc.WithInsecure())
	if err != nil || conn.GetState() == connectivity.Idle {
		fmt.Println("Failed to connect ", err)
		panic(fmt.Sprintf("Failed to connect %v", err))
	}
	fmt.Println("Get new Client")
	client := NewWalletServiceClient(conn)
	fmt.Println(client, err)
	return &WalletSearchClient{Client: client, Conn: conn}
}

func NewWalletTransferClient() *WalletTransferClient {
	fmt.Println("Established connection")
	conn, err := grpc.DialContext(context.Background(),
		HostConfig.RemoteTransferHost,
		grpc.WithReturnConnectionError(),
		grpc.WithBlock(),
		grpc.WithNoProxy(),
		grpc.WithInsecure())
	if err != nil || conn.GetState() == connectivity.Idle {
		fmt.Println("Failed to connect ", err)
		panic(fmt.Sprintf("Failed to connect %v", err))
	}
	fmt.Println("Get new Client")
	client := NewWalletServiceClient(conn)
	fmt.Println(client, err)
	return &WalletTransferClient{Client: client, Conn: conn}
}

func NewWalletAccountClient() *WalletAccountClient {
	fmt.Println("Established connection")
	conn, err := grpc.DialContext(context.Background(),
		HostConfig.RemoteAccountHost,
		grpc.WithReturnConnectionError(),
		grpc.WithBlock(),
		grpc.WithNoProxy(),
		grpc.WithInsecure())
	if err != nil || conn.GetState() == connectivity.Idle {
		fmt.Println("Failed to connect ", err)
		panic(fmt.Sprintf("Failed to connect %v", err))
	}
	fmt.Println("Get new Client")
	client := NewWalletServiceClient(conn)
	fmt.Println(client, err)
	return &WalletAccountClient{Client: client, Conn: conn}
}
