package pullrequests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	requiredApprovals   = 2
	githubAPIURL        = "https://api.github.com/graphql"
	githubQueryTemplate = `query {
			viewer {
				login
				organization(login: \"%s\") {
					repositories(first: 100) {
						nodes {
							name
							pullRequests(first: 10, states: OPEN) {
								nodes {
									url
									createdAt
									updatedAt
									author {
										login
									}
									reviews(first: 50) {
										nodes {
											author{
												login
											}
											state
										}
									}
								}
							}
						}
					}
				}
			}
		}`
)

func GetPRs(orgName, apiKey string) (string, PullRequestContainer, error) {
	body, err := makeRequest(apiKey, orgName)
	if err != nil {
		return "", nil, err
	}

	user, prs, err := normaliseResponse(body)
	if err != nil {
		return "", nil, err
	}

	return user, prs, nil
}

func makeRequest(apiKey, orgName string) ([]byte, error) {
	template := strings.Replace(githubQueryTemplate, "\n", ` \n `, -1)
	template = strings.Join(strings.Fields(template), " ")
	graphQuery := fmt.Sprintf(template, orgName)
	q := fmt.Sprintf(`{"query": "%s"}`, graphQuery)

	client := &http.Client{}
	req, err := http.NewRequest("POST", githubAPIURL, bytes.NewBufferString(q))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf(`bearer %s`, apiKey))
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("Error: %s\n\n%s", resp.Status, body)
		return nil, errors.New(msg)
	}

	return body, nil
}

func normaliseResponse(resp []byte) (string, []PullRequest, error) {
	var data map[string]map[string]githubResp
	if err := json.Unmarshal(resp, &data); err != nil {
		return "", nil, err
	}

	var pullRequests []PullRequest

	for _, r := range data["data"]["viewer"].Organization.Repositories.Nodes {
		if len(r.PullRequests.Nodes) > 0 {
			for _, rawPR := range r.PullRequests.Nodes {
				pr := PullRequest{
					RepoName:         r.Name,
					URL:              rawPR.URL,
					CreatedAt:        rawPR.CreatedAt,
					UpdatedAt:        rawPR.UpdatedAt,
					Author:           rawPR.Author.Login,
					TotalReviews:     0,
					Approved:         false,
					ChangesRequested: false,
				}

				for _, r := range rawPR.Reviews.Nodes {
					pr.AddReviewer(r.Author.Login)
					pr.AddState(r.State)
				}

				if pr.approvals >= requiredApprovals {
					pr.Approved = true
				}

				pullRequests = append(pullRequests, pr)
			}
		}
	}
	return data["data"]["viewer"].Login, pullRequests, nil
}

type githubResp struct {
	Login        string `json:"login"`
	Organization struct {
		Repositories struct {
			Nodes []repoResp `json:"nodes"`
		} `json:"repositories"`
	} `json:"organization"`
}

type repoResp struct {
	Name         string `json:"name"`
	PullRequests struct {
		Nodes []pullRequestResp `json:"nodes"`
	} `json:"pullRequests"`
}

type pullRequestResp struct {
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Author    struct {
		Login string `json:"login"`
	} `json:"author"`
	Reviews struct {
		Nodes []reviewResp `json:"nodes"`
	} `json:"reviews"`
}

type reviewResp struct {
	Author struct {
		Login string `json:"login"`
	} `json:"author"`
	State string `json:"state"`
}
