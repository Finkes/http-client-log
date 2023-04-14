package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Finkes/http-client-log"
	"github.com/Finkes/http-client-log/pkg/receiver"
	"github.com/google/go-github/v41/github"
	"log"
	"net/http"
)

func main() {
	http_client_log.Init(
		http_client_log.WithFormat(receiver.FormatSummary),
		http_client_log.WithLogReceiver(nil),
		http_client_log.WithFileReceiver(".", nil),
	)
	defer http_client_log.Cleanup()

	demoSift()
	demoGet()
	demoPost()
	demoGithub()
	demoGithubRaw()
	demoRateLimit429()
	demoGithubRateLimit()

	fmt.Println("demo requests done")
}

func demoRateLimit429() {
	for i := 0; i < 13; i++ {
		http.Get("https://www.cloudflare.com/rate-limit-test/")
	}
}

func demoGet() {
	// make some dummy API calls
	http.Get("https://jsonplaceholder.typicode.com/todos/1")
}

func demoPost() {
	post := map[string]interface{}{
		"title":  "title",
		"userId": 2,
		"body":   "body",
	}
	postJSON, err := json.Marshal(post)
	if err != nil {
		log.Fatal(err)
	}

	http.Post("https://jsonplaceholder.typicode.com/posts", "application/json", bytes.NewBuffer(postJSON))
	http.Post("http://jsonplaceholder.typicode.com/posts", "application/json", bytes.NewBuffer(postJSON))
}

func demoGithub() {
	client := github.NewClient(nil)

	// list public repositories for org "github"
	opt := &github.RepositoryListByOrgOptions{Type: "public"}
	client.Repositories.ListByOrg(context.Background(), "github", opt)

	// list public repositories for org "getkin"
	opt = &github.RepositoryListByOrgOptions{Type: "public"}
	client.Repositories.ListByOrg(context.Background(), "getkin", opt)

	demoGithubRepositories(client)
}

func demoGithubRepositories(client *github.Client) {
	owner := "github"
	repo := "github"
	branch := "main"
	ctx := context.Background()
	client.Repositories.Get(context.Background(), owner, repo)
	client.Repositories.GetBranch(context.Background(), owner, repo, branch, true)
	client.Repositories.GetBranchProtection(context.Background(), owner, repo, branch)
	client.Repositories.GetCodeOfConduct(context.Background(), owner, repo)
	client.Repositories.GetPagesInfo(ctx, owner, repo)
}

func demoGithubRateLimit() {
	client := github.NewClient(nil)
	opt := &github.RepositoryListByOrgOptions{Type: "public"}

	for i := 0; i < 62; i++ {
		// list public repositories for org "getkin"
		client.Repositories.ListByOrg(context.Background(), "getkin", opt)
	}
}

func demoGithubRaw() {
	http.Get("https://api.github.com/users/octocat/orgs")
}

func demoSift() {
	// can help to debug http2 issue
	values := map[string]string{
		"$api_key":           "12345",
		"$chargeback_reason": "$fraud",
		"$chargeback_state":  "$lost",
		"$order_id":          "123456",
		"$transaction_id":    "12345",
		"$type":              "$chargeback",
		"$user_id":           "123456",
	}
	payload, err := json.Marshal(values)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post("https://api.sift.com/v205/events", "application/json",
		bytes.NewBuffer(payload))

	if err != nil {
		log.Fatal(err)
	}

	var res map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		fmt.Printf("error decoding sift response: %v", err)
	}
}
