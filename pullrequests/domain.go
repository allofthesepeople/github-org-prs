package pullrequests

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	Columns = []string{
		"RepoName",
		"URL",
		"CreatedAt",
		"UpdatedAt",
		"Author",
		"TotalReviews",
		"Approved",
		"ChangesRequested",
		"Reviewers",
	}
)

type PullRequestContainer []PullRequest

func (c *PullRequestContainer) Filter(column string, value interface{}) PullRequestContainer {
	return *c
}

func (c *PullRequestContainer) Sort(column, direction string) PullRequestContainer {
	prc := *c
	s := column + "__" + direction

	switch s {
	case "RepoName__asc":
		sort.Slice(prc, func(i, j int) bool { return prc[i].RepoName < prc[j].RepoName })
	case "RepoName__desc":
		sort.Slice(prc, func(i, j int) bool { return prc[i].RepoName > prc[j].RepoName })
	case "CreatedAt__asc":
		sort.Slice(prc, func(i, j int) bool { return prc[i].CreatedAt.Unix() < prc[j].CreatedAt.Unix() })
	case "CreatedAt__desc":
		sort.Slice(prc, func(i, j int) bool { return prc[i].CreatedAt.Unix() > prc[j].CreatedAt.Unix() })
	case "UpdatedAt__asc":
		sort.Slice(prc, func(i, j int) bool { return prc[i].UpdatedAt.Unix() < prc[j].UpdatedAt.Unix() })
	case "UpdatedAt__desc":
		sort.Slice(prc, func(i, j int) bool { return prc[i].UpdatedAt.Unix() > prc[j].UpdatedAt.Unix() })
	case "Author__asc":
		sort.Slice(prc, func(i, j int) bool { return prc[i].Author < prc[j].Author })
	case "Author__desc":
		sort.Slice(prc, func(i, j int) bool { return prc[i].Author > prc[j].Author })
	case "TotalReviews__asc":
		sort.Slice(prc, func(i, j int) bool { return prc[i].TotalReviews < prc[j].TotalReviews })
	case "TotalReviews__desc":
		sort.Slice(prc, func(i, j int) bool { return prc[i].TotalReviews > prc[j].TotalReviews })
	}

	return prc
}

func (c *PullRequestContainer) Headers(cols []string) []string {
	return []string{
		"RepoName",
		"URL",
		"CreatedAt",
		"UpdatedAt",
		"Author",
		"TotalReviews",
		"Approved",
		"ChangesRequested",
		"Reviewers",
	}
}

type PullRequest struct {
	RepoName         string    `json:"repoName"`
	URL              string    `json:"url"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	Author           string    `json:"author"`
	TotalReviews     int       `json:"totalReviews"`
	Approved         bool      `json:"approved"`
	ChangesRequested bool      `json:"needReview"`
	Reviewers        []string  `json:"reviewers"`
	approvals        int
}

func (pr *PullRequest) AddReviewer(reviewer string) {
	pr.TotalReviews = pr.TotalReviews + 1

	reviewerExists := false
	for _, r := range pr.Reviewers {
		if r == reviewer {
			reviewerExists = true
		}
	}

	if reviewerExists == false {
		pr.Reviewers = append(pr.Reviewers, reviewer)
	}
}

func (pr *PullRequest) AddState(state string) {
	switch state {
	case "APPROVED":
		pr.approvals = pr.approvals + 1
	case "CHANGES_REQUESTED":
		pr.ChangesRequested = true
	}
}

func (pr *PullRequest) ToStrings(cols []string) []string {
	var retCols []string

	for _, c := range cols {
		switch c {
		case "RepoName":
			retCols = append(retCols, pr.RepoName)
		case "URL":
			retCols = append(retCols, pr.URL)
		case "CreatedAt":
			retCols = append(retCols, pr.CreatedAt.Format(time.RFC822))
		case "UpdatedAt":
			retCols = append(retCols, pr.UpdatedAt.Format(time.RFC822))
		case "Author":
			retCols = append(retCols, pr.Author)
		case "TotalReviews":
			retCols = append(retCols, strconv.Itoa(pr.TotalReviews))
		case "Approved":
			retCols = append(retCols, strconv.FormatBool(pr.Approved))
		case "ChangesRequested":
			retCols = append(retCols, strconv.FormatBool(pr.ChangesRequested))
		case "Reviewers":
			retCols = append(retCols, strings.Join(pr.Reviewers, ", "))
		}
	}

	return retCols
}
