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
