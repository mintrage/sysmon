package notifier

import (
	"fmt"
	"net/http"
	"net/url"
)

func SendAlert(token, chatID, message string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	data := url.Values{
		"chat_id": {chatID},
		"text":    {message},
	}

	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		return fmt.Errorf("ошибка сети при отправке в TG: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram вернул ошибку, код статуса: %d", resp.StatusCode)
	}

	return nil
}
