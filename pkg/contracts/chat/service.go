package chat

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const ServiceName = "chat.ChatService"

const createMessageMethod = "/" + ServiceName + "/CreateMessage"

type CreateMessageRequest struct {
	SenderID   string
	ReceiverID string
	Content    string
}

type CreateMessageResponse struct {
	Message Message
}

type Message struct {
	ID         string
	SenderID   string
	ReceiverID string
	Content    string
	CreatedAt  time.Time
}

type ChatServiceClient interface {
	CreateMessage(ctx context.Context, in *CreateMessageRequest, opts ...grpc.CallOption) (*CreateMessageResponse, error)
}

type chatServiceClient struct {
	cc *grpc.ClientConn
}

func NewChatServiceClient(cc *grpc.ClientConn) ChatServiceClient {
	return &chatServiceClient{cc: cc}
}

func (c *chatServiceClient) CreateMessage(ctx context.Context, in *CreateMessageRequest, opts ...grpc.CallOption) (*CreateMessageResponse, error) {
	out := new(CreateMessageResponse)
	if err := c.cc.Invoke(ctx, createMessageMethod, in, out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

type ChatServiceServer interface {
	CreateMessage(ctx context.Context, in *CreateMessageRequest) (*CreateMessageResponse, error)
}

type UnimplementedChatServiceServer struct{}

func (UnimplementedChatServiceServer) CreateMessage(ctx context.Context, in *CreateMessageRequest) (*CreateMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMessage not implemented")
}

func RegisterChatServiceServer(s grpc.ServiceRegistrar, srv ChatServiceServer) {
	s.RegisterService(&grpc.ServiceDesc{
		ServiceName: ServiceName,
		HandlerType: (*ChatServiceServer)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "CreateMessage",
				Handler:    _ChatService_CreateMessage_Handler,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "chat_service.proto",
	}, srv)
}

func _ChatService_CreateMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).CreateMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: createMessageMethod,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).CreateMessage(ctx, req.(*CreateMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}
