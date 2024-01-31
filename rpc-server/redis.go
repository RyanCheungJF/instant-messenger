package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
)

// put structs in front of function definitions to signify
// is a method of the struct
// public attrs or methods signified by capital letters
type RedisClient struct {
	cli *redis.Client
}

// in golang the ctx var is the context instance from the background function
// know info about the env it is executed in
// https://www.digitalocean.com/community/tutorials/how-to-use-contexts-in-go
func (c *RedisClient) InitClient(ctx context.Context) error {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// connection test
	err := client.Ping(ctx).Err()
	if err != nil {
		// better to pass error back up to top level
		return err
	}

	c.cli = client
	return nil
}

func (c *RedisClient) SaveMessage(ctx context.Context, roomID string, message *Message) error {
	// convert to json
	text, err := json.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}

	member := &redis.Z{
		Score:  float64(message.Timestamp),
		Member: text,
	}

	_, err = c.cli.ZAdd(ctx, roomID, *member).Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *RedisClient) GetMessagesById(ctx context.Context, roomID string, start, end int64, reverse bool) ([]*Message, error) {
	var (
		rawMessages []string
		messages    []*Message
		err         error
	)

	// get the messages in either in order (first msg is earliest) else reverse
	if reverse {
		rawMessages, err = c.cli.ZRevRange(ctx, roomID, start, end).Result()
		if err != nil {
			return nil, err
		}
	} else {
		rawMessages, err = c.cli.ZRange(ctx, roomID, start, end).Result()
		if err != nil {
			return nil, err
		}
	}

	for _, msg := range rawMessages {
		temp := &Message{}
		// deserialize message into a different format 
		err := json.Unmarshal([]byte(msg), temp)
		if err != nil {
			return nil, err
		}
		messages = append(messages, temp)
	}

	return messages, nil
}
