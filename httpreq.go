package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

func httpReq(ctx context.Context, url string, postData []byte) (string, error) {
	var err error
	var request *http.Request
	if postData != nil {
		request, err = http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(postData))
		if err != nil {
			return "", err
		}
		request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	} else {
		request, err = http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return "", err
		}
		request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		request.Header.Set("Accept-Language", "en-GB,en;q=0.9,en-US;q=0.8,hu;q=0.7")
		request.Header.Set("Cache-Control", "max-age=0")
		request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36")
	}

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("api status code: %d", resp.StatusCode)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}
