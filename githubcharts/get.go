package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
)

func get(org string, repo string, outf string) error {
	f, err := os.Create(outf)
	if err != nil {
		return err
	}
	defer f.Close()

	var client *github.Client
	token := os.Getenv("GITHUB_AUTH_TOKEN")
	if token == "" {
		log.Print("No auth token")
		client = github.NewClient(nil)
	} else {
		log.Print("Using auth token")
		tc := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		))
		client = github.NewClient(tc)
	}
	opts := github.IssueListByRepoOptions{State: "all"}
	opts.PerPage = 100

	w := csv.NewWriter(f)
	for {
		issues, r, err := client.Issues.ListByRepo(org, repo, &opts)
		if limit, ok := err.(*github.RateLimitError); ok {
			log.Printf("Rate limit: %+v", limit)
			wait := limit.Rate.Reset.Sub(time.Now()) + 1*time.Minute
			log.Printf("Waiting for %v", wait)
			time.Sleep(wait)
			continue
		}
		if err != nil {
			return err
		}
		log.Printf("Issues.ListRepositoryEvents (page %d of %d): %+v", r.NextPage-1, r.LastPage, r)
		for _, issue := range issues {
			if issue.PullRequestLinks == nil {
				var num, state, createdAt, closedAt string
				if issue.Number != nil {
					num = strconv.Itoa(*issue.Number)
				}
				if issue.State != nil {
					state = *issue.State
				}
				if issue.CreatedAt != nil {
					createdAt = (*issue.CreatedAt).Format(time.RFC3339)
				}
				if issue.ClosedAt != nil {
					closedAt = (*issue.ClosedAt).Format(time.RFC3339)
				}
				if err := w.Write([]string{num, state, createdAt, closedAt}); err != nil {
					return err
				}
			}
		}
		if r.NextPage == 0 {
			break
		}
		opts.Page = r.NextPage
	}
	w.Flush()
	return w.Error()
}
