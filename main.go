package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Response struct {
	OK     bool        `json:"ok"`
	Result interface{} `json:"result"` // Accepts any type (string or object)
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	chatID := r.URL.Query().Get("chat")
	text := r.URL.Query().Get("text")

	if chatID == "" || text == "" {
		http.Error(w, "Missing chat or text parameters", http.StatusBadRequest)
		return
	}

	sendTelegramMessage(chatID, text, w)
}

func sendTelegramMessage(chatID, text string, w http.ResponseWriter) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		http.Error(w, "Bot token not set", http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", token, chatID, text)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error sending request to Telegram: %v", err)
		if w != nil {
			http.Error(w, "Failed to send request", http.StatusInternalServerError)
		}
		return
	}
	defer resp.Body.Close()

	var telegramResp Response
	if err := json.NewDecoder(resp.Body).Decode(&telegramResp); err != nil {
		log.Printf("Error decoding Telegram response: %v", err)
		if w != nil {
			http.Error(w, "Failed to parse response", http.StatusInternalServerError)
		}
		return
	}

	if w != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(telegramResp)
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	chatID := os.Getenv("TELEGRAM_CHAT_ID")
	if chatID != "" {
		sendTelegramMessage(chatID, "HTTP Bridge has started", nil)
	}

	http.HandleFunc("/send", sendMessage)
	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
