package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/kitex_gen/rpc"
)

// IMServiceImpl implements the last service interface defined in the IDL.
type IMServiceImpl struct{}

func (s *IMServiceImpl) Send(ctx context.Context, req *rpc.SendRequest) (*rpc.SendResponse, error) {
	if err := validateSendRequest(req); err != nil {
		return nil, err
	}

	timestamp := time.Now().Unix()
	message := &Message{
		Message:   req.Message.GetText(),
		Sender:    req.Message.GetSender(),
		Timestamp: timestamp,
	}

	roomID, err := getRoomID(req.Message.GetChat())
	if err != nil {
		return nil, err
	}

	err = redisClient.SaveMessage(ctx, roomID, message)
	if err != nil {
		return nil, err
	}

	response := rpc.NewSendResponse()
	// Code is not status code 0
	response.Code, response.Msg = 0, "success"
	return response, nil
}

func (s *IMServiceImpl) Pull(ctx context.Context, req *rpc.PullRequest) (*rpc.PullResponse, error) {
	roomID, err := getRoomID(req.GetChat())
	if err != nil {
		return nil, err
	}

	start := req.GetCursor()
	end := start + int64(req.GetLimit())

	messages, err := redisClient.GetMessagesById(ctx, roomID, start, end, req.GetReverse())
	if err != nil {
		return nil, err
	}

	responseMessages := make([]*rpc.Message, 0)
	var counter int32 = 0
	var nextCursor int64 = 0
	hasMore := false
	for _, msg := range messages {
		if counter+1 > req.GetLimit() {
			hasMore = true
			nextCursor = end
			// don't return the last message, just an indicator for the response
			break
		}
		temp := &rpc.Message{
			Chat:     req.GetChat(),
			Text:     msg.Message,
			Sender:   msg.Sender,
			SendTime: msg.Timestamp,
		}
		responseMessages = append(responseMessages, temp)
		counter += 1
	}

	response := rpc.NewPullResponse()
	response.Messages = responseMessages
	response.Code = 0
	response.Msg = "success"
	response.HasMore = &hasMore
	response.NextCursor = &nextCursor

	return response, nil
}

func validateSendRequest(req *rpc.SendRequest) error {
	senders := strings.Split(req.Message.Chat, ":")
	if len(senders) != 2 {
		err := fmt.Errorf("invalid chat ID %s, required format of id1:id2", req.Message.Chat)
		return err
	}

	sender1, sender2 := senders[0], senders[1]
	if req.Message.GetSender() != sender1 && req.Message.GetSender() != sender2 {
		err := fmt.Errorf("sender %s is not in the chat room", req.Message.GetSender())
		return err
	}

	return nil
}

func getRoomID(chat string) (string, error) {
	var roomID string

	lowercase := strings.ToLower(chat)
	senders := strings.Split(lowercase, ":")
	if len(senders) != 2 {
		err := fmt.Errorf("invalid chat ID %s, required format of id1:id2", chat)
		return "", err
	}

	// make it such that users of s1:s2 and s2:s1 are stored in the same place
	sender1, sender2 := senders[0], senders[1]
	// compares the sender and receiver alphabetically, and sort them asc to form the room ID
	if comp := strings.Compare(sender1, sender2); comp == 1 {
		roomID = fmt.Sprintf("%s:%s", sender2, sender1)
	} else {
		roomID = fmt.Sprintf("%s:%s", sender1, sender2)
	}

	return roomID, nil
}
