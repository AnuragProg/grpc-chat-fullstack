package main


import (
	"fmt"
	"log"
	"net"
	"google.golang.org/grpc"
	"sync"

	pb "github.com/AnuragProg/grpc-chat-fullstack/pb"
)

const PORT = 3000

type ChatServer struct {
	pb.UnimplementedChatServiceServer
	UserCount uint32
	Users map[string]*pb.ChatService_ConverseServer
	UserMutex *sync.Mutex
}

func (c *ChatServer) FanOut(cs pb.ChatService_ConverseServer, current_user_id, msg string){
	c.UserMutex.Lock()
	defer c.UserMutex.Unlock()
	msg_obj := &pb.Message{Msg:msg}
	for user_id, socket := range c.Users{
		if user_id == current_user_id {
			continue
		}
		(*socket).Send(msg_obj)
	}
}

func (c *ChatServer) Converse(cs pb.ChatService_ConverseServer) error {
	c.UserMutex.Lock()	
	c.UserCount += 1
	current_user_id := fmt.Sprintf("User %v", c.UserCount)
	c.Users[current_user_id] = &cs
	c.UserMutex.Unlock()

	log.Println(current_user_id,"connected")
	c.FanOut(cs, current_user_id, current_user_id+" connected")
	cs.Send(&pb.Message{Msg:"Welcome to the chat "+current_user_id})
	for {
		msg, err := cs.Recv()
		if err != nil{
			log.Println(current_user_id,"disconnected")
			c.FanOut(cs, current_user_id, current_user_id+" disconnected")
			delete(c.Users, current_user_id)
			return nil
		}
		c.FanOut(cs, current_user_id, current_user_id+": "+msg.Msg)
	}
}

func main(){
	lis, _ := net.Listen("tcp", fmt.Sprintf(":%v", PORT))
	chat_service_registrar := grpc.NewServer()
	chat_server := ChatServer{
		UserCount : 0,
		Users : make(map[string]*pb.ChatService_ConverseServer),
		UserMutex : &sync.Mutex{},
	}
	pb.RegisterChatServiceServer(chat_service_registrar, &chat_server)
	if err := chat_service_registrar.Serve(lis); err!=nil{
		log.Fatalf(err.Error())
	}
}

