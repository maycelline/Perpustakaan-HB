package controllers

import (
	"Perpustakaan-HB/model"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

func GetPopularBooksFromCache() []model.Book {
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

	var books []model.Book
	_ = json.Unmarshal([]byte(value), &books)

	return books
}

func SetPopularBooksCache(books []model.Book) {
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

	//Popular book direset setiap 1 minggu
	err = client.Set(ctx, "books", converted, time.Duration(time.Monday)).Err()
	if err != nil {
		log.Println(err)
		return
	} else {
		log.Println("Cache set")
	}
}
