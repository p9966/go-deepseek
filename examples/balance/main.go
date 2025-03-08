package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/p9966/go-deepseek"
)

func main() {
	client := deepseek.NewClient(os.Getenv("DEEPSEEK_API_KEY"))
	if resp, err := client.GetBalance(context.TODO()); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(resp.IsAvailable)
	}
}
