package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// https://docs.github.com/en/rest/users/users?apiVersion=2022-11-28#get-a-user
type GitHubUser struct {
}

func githubSuccessHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	fmt.Println("callback from github", code)
	token, err := tokenForCode(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("access token", token)
	user, err := getGithubUser(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("user", user)
	sessionManager.Put(r.Context(), "pie", "pizza")
	// TODO: Find or create user according to GitHub user_id
	// TODO: Set (our) user_id, user on session

	// DONE
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

}

func tokenForCode(code string) (string, error) {
	// TODO: Validate the code
	url := "https://github.com/login/oauth/access_token"

	payload := map[string]string{
		"client_id":     os.Getenv("GH_CLIENT_ID"),
		"client_secret": os.Getenv("GH_CLIENT_SECRET"),
		"code":          code,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var tokenResp map[string]any
	if err := json.Unmarshal(respBody, &tokenResp); err != nil {
		return "", err
	}

	token, ok := tokenResp["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("no access toke in response: %v", tokenResp)
	}
	return token, nil

}

func getGithubUser(token string) (map[string]any, error) {
	url := "https://api.github.com/user"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API responded with %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)

	var user map[string]any
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}

	return user, nil

}
