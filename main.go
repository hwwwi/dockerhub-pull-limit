package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func getToken(url string) string {
	var body struct {
		Token string `json:"token"`
	}

	resp, _ := http.Get(url)
	defer resp.Body.Close()
	err := json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		panic(err)
	}
	return body.Token
}

func getCurrentLimit(token string) (int, error) {
	url := "https://registry-1.docker.io/v2/ratelimitpreview/test/manifests/latest"
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return -1, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}

	if resp.StatusCode != http.StatusOK {
		return -1, fmt.Errorf("wron response: %v", resp.StatusCode)
	}

	rateLimit := resp.Header.Get("RateLimit-Remaining")
	if rateLimit == "" {
		return -1, fmt.Errorf("failed to get RateLimit-Remaining")
	}

	rateLimitRemain, err := strconv.Atoi(strings.Split(rateLimit, ";")[0])
	if err != nil {
		return -1, fmt.Errorf("failed to parse RateLimit-Remaining")
	}
	return rateLimitRemain, nil
}

func main() {
	url := "https://auth.docker.io/token?service=registry.docker.io&scope=repository:ratelimitpreview/test:pull"
	remain, err := getCurrentLimit(getToken(url))
	if err != nil {
		panic(err)
	}
	fmt.Println(remain)
}
