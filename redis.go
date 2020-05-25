package main

import (
	"context"

	"github.com/childe/gohangout/codec"
	"github.com/go-redis/redis/v8"
)

type RedisInput struct {
	key     string
	decoder codec.Decoder

	client *redis.Client
}

func New(config map[interface{}]interface{}) interface{} {
	client := redis.NewClient(&redis.Options{
		Addr:     config["address"].(string),
		Password: config["password"].(string),
		DB:       config["db"].(int),
	})
	codertype := "json"
	if v, ok := config["codec"]; ok {
		codertype = v.(string)
	}
	key := config["key"].(string)
	decoder := codec.NewDecoder(codertype)
	return &RedisInput{
		key,
		decoder,
		client,
	}
}

func (p *RedisInput) ReadOneEvent() map[string]interface{} {
	cmd := p.client.BLPop(context.TODO(), 0, p.key)
	message := cmd.Val()[1]
	return p.decoder.Decode([]byte(message))
}

func (p *RedisInput) Shutdown() {}
