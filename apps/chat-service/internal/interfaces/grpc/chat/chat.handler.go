package chat

import (
	"context"
	"time"

	bootstrap "golang-social-media/apps/chat-service/internal/infrastructure/bootstrap"
	commandcontracts "golang-social-media/apps/chat-service/internal/application/command/contracts"
	"golang-social-media/apps/chat-service/internal/interfaces/grpc/mappers"
	"golang-social-media/pkg/logger"
	chatv1 "golang-social-media/pkg/gen/chat/v1"
)

type Handler struct {
	createMessageCmd commandcontracts.CreateMessageCommand
	dtoMapper        mappers.MessageDTOMapper
	chatv1.UnimplementedChatServiceServer
}

func NewHandler(deps *bootstrap.Dependencies, dtoMapper mappers.MessageDTOMapper) *Handler {
	return &Handler{
		createMessageCmd: deps.CreateMessageCmd,
		dtoMapper:        dtoMapper,
	}
}

func (h *Handler) CreateMessage(ctx context.Context, req *chatv1.CreateMessageRequest) (*chatv1.CreateMessageResponse, error) {
	startTime := time.Now()

	// Prepare request using mapper
	requestStart := time.Now()
	cmdReq := commandcontracts.CreateMessageCommandRequest{
		SenderID:   req.GetSenderId(),
		ReceiverID: req.GetReceiverId(),
		Content:    req.GetContent(),
	}
	requestDuration := time.Since(requestStart)

	// Execute command
	commandStart := time.Now()
	msg, err := h.createMessageCmd.Execute(ctx, cmdReq)
	commandDuration := time.Since(commandStart)

	if err != nil {
		totalDuration := time.Since(startTime)
		logger.Component("chat.grpc.create_message").
			Error().
			Err(err).
			Str("sender_id", req.GetSenderId()).
			Str("receiver_id", req.GetReceiverId()).
			Dur("request_prep_ms", requestDuration).
			Dur("command_exec_ms", commandDuration).
			Dur("total_ms", totalDuration).
			Msg("failed to create message")
		return nil, err
	}

	// Build response using mapper
	responseStart := time.Now()
	resp := h.dtoMapper.ToCreateMessageResponse(msg)
	responseDuration := time.Since(responseStart)

	totalDuration := time.Since(startTime)

	logger.Component("chat.grpc.create_message").
		Info().
		Str("message_id", msg.ID).
		Dur("request_prep_ms", requestDuration).
		Dur("command_exec_ms", commandDuration).
		Dur("response_build_ms", responseDuration).
		Dur("total_ms", totalDuration).
		Msg("gRPC request completed")

	return resp, nil
}
