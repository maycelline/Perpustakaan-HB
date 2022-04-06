package controllers

import (
	"Perpustakaan-HB/model"
	"context"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"
)

func GetPopularBooksFromCache() []model.PopularBook {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	var ctx = context.Background()

	value, err := client.Get(ctx, "books").Result()
	if err != nil {
		return nil
	}

	var books []model.PopularBook
	_ = json.Unmarshal([]byte(value), &books)

	return books
}

func SetPopularBooksCache(books []model.PopularBook) {
	converted, err := json.Marshal(books)
	if err != nil {
		log.Println(err)
		return
	}

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	var ctx = context.Background()

	err = client.Set(ctx, "books", converted, 0).Err()
	if err != nil {
		log.Println(err)
		return
	} else {
		log.Println("Cache set")
	}
}
