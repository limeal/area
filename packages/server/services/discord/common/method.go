package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// It takes a JSON response from the API, and returns the time at which the rate limit will reset, and
// any errors that occurred
func HandleRateLimitError(encode []byte) (time.Time, error) {
	body := make(map[string]interface{})

	err := json.Unmarshal(encode, &body)
	if err != nil {
		return time.Now(), err
	}

	if body["retry_after"] == nil {
		return time.Now(), errors.New("retry_after not found")
	}

	retryAfter, ok := body["retry_after"].(float64)
	if !ok {
		return time.Now(), errors.New("retry_after is not a float64")
	}

	fmt.Println("Beep boop, I'm a bot. I'm going to sleep for", retryAfter, "seconds")
	return time.Now().Add(time.Duration(retryAfter) * time.Second), nil
}

// If the rate limit is in effect, return true. Otherwise, return false
func IsRateLimited(store *map[string]interface{}) bool {
	if (*store)["rate:time:wait"] != nil {
		if time.Now().Before((*store)["rate:time:wait"].(time.Time)) {
			return true
		}
		(*store)["rate:time:wait"] = nil
	}
	return false
}

// If the response is a rate limit error, then set the rate limit time to wait in the store
func SetRateLimit(
	httpResp *http.Response,
	store *map[string]interface{},
	encode []byte,
) (bool, error) {
	var err error
	if httpResp != nil && httpResp.StatusCode == 429 {
		(*store)["rate:time:wait"], err = HandleRateLimitError(encode)
		if err != nil {
			return false, err
		}
		return false, nil
	}
	return true, nil
}
