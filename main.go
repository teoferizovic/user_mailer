package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/go-redis/redis"
	"user_mailer/model"
	"user_mailer/processor"
)

var conf model.Config
var err error

type dbStore struct {
	client *redis.Client
}

var dbConn dbStore

func init(){

	//load config file
	if _, err := toml.DecodeFile("./config.toml", &conf); err != nil {
		fmt.Println(err)
	}

	//conection on Redis
	client := redis.NewClient(&redis.Options{
		Addr:     conf.RedisHost+":"+conf.RedisPort,
		Password: conf.RedisPassword,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	dbConn = dbStore{client:client}
}

func main() {

	err = processor.SubRedis(dbConn.client,conf)
	if err != nil{
		fmt.Println(err)
	}

}
