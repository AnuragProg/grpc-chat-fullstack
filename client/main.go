package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	pb "github.com/AnuragProg/grpc-chat-fullstack/pb"
	"google.golang.org/grpc"
)

func log_replies(socket pb.ChatService_ConverseClient) error {
	for{
		msg, err := socket.Recv()	
		if err != nil{
			log.Println("Closing connection with server")
			return nil
		}
		fmt.Println("Reply:", msg.GetMsg())
	}
}

func main(){
	log.Println("Client started ...")
	log.Println("Connecting to gprc server...")
	conn, err := grpc.Dial("localhost:3000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil{
		log.Println(err.Error())
	}
	client := pb.NewChatServiceClient(conn)
	input_reader := bufio.NewReader(os.Stdin)
	socket, err := client.Converse(context.TODO())
	if err != nil{
		log.Fatalf("Server down")
	}
	go log_replies(socket)
	fmt.Println("You can write messages and enter to send them")
	for {
		input, _, _ := input_reader.ReadLine()
		socket.Send(&pb.Message{Msg:string(input)})
	}
}
