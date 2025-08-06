package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type TelegramBot struct {
	Token  string
	ChatID string
}

type TelegramResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description"`
	ErrorCode   int    `json:"error_code"`
}

func NewTelegramBot(token, chatID string) *TelegramBot {
	return &TelegramBot{
		Token:  token,
		ChatID: chatID,
	}
}

func (t *TelegramBot) SendDocument(filePath, caption string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument", t.Token)

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add chat_id field
	if err := writer.WriteField("chat_id", t.ChatID); err != nil {
		return fmt.Errorf("failed to write chat_id field: %v", err)
	}

	// Add caption field if provided
	if caption != "" {
		if err := writer.WriteField("caption", caption); err != nil {
			return fmt.Errorf("failed to write caption field: %v", err)
		}
	}

	// Add document field
	part, err := writer.CreateFormFile("document", filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf("failed to create form file: %v", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("failed to copy file content: %v", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}

	req, err := http.NewRequest("POST", url, &body)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	var telegramResp TelegramResponse
	if err := json.Unmarshal(respBody, &telegramResp); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	if !telegramResp.OK {
		return fmt.Errorf("telegram API error: %s (code: %d)", telegramResp.Description, telegramResp.ErrorCode)
	}

	return nil
}

func (t *TelegramBot) SendMessage(message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.Token)

	payload := map[string]interface{}{
		"chat_id": t.ChatID,
		"text":    message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	var telegramResp TelegramResponse
	if err := json.Unmarshal(respBody, &telegramResp); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	if !telegramResp.OK {
		return fmt.Errorf("telegram API error: %s (code: %d)", telegramResp.Description, telegramResp.ErrorCode)
	}

	return nil
}