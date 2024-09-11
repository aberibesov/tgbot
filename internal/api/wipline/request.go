package wipline

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"time"
)

const url = "http://172.22.3.4"

func ApiRequest(xmlRequest []byte) ([]byte, error) {
	httpClient := http.Client{Timeout: 5 * time.Second}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(xmlRequest))
	if err != nil {
		log.Fatal(err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
