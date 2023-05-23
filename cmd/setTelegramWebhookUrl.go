package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Response struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	tgBotToken := os.Getenv("TG_BOT_TOKEN")
	if tgBotToken == "" {
		fmt.Println("Please set TG_BOT_TOKEN environment variable")
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fmt.Println("Please provide a URL as the first argument")
		os.Exit(1)
	}

	url := os.Args[1]
	setWebhookUrl := "https://api.telegram.org/bot" + tgBotToken + "/setWebhook?url=" + url + "/api/v1/bot"

	resp, err := http.Get(setWebhookUrl)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	if !response.Ok {
		fmt.Printf("Error: %s\n", response.Description)
		os.Exit(1)
	}

	fmt.Println(string(body))
	fmt.Println("Webhook URL is: " + setWebhookUrl)
}
