// lint:file-ignore U1000 repository

package pullrequest

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
	"github.com/google/go-github/github"
)

// LinesFromTextDescription parses t, assumed to be a plain text from a Github Pull Request description,
func LinesFromTextDescription(text string, ignoreLineFn func(string) bool) []string {
	lines := strings.Split(text, "\n")
	flines := []string{}

	if ignoreLineFn == nil {
		ignoreLineFn = func(t string) bool {
			return t == "\t---" || t == "\t" || t == "\t\t\t" || t == "---" || t == "\n"
		}
	}

	for _, l := range lines {
		if ignoreLineFn(l) {
			continue
		}
		flines = append(flines, l)
	}
	return flines
}

var ErrMatchFuncNil = errors.New("match function must not be nil")
var ErrStackLinksNotFound = errors.New("stack links not found")

func StackLines(header string, match func(header, line string) bool, text string) ([]string, error) {
	lines := LinesFromTextDescription(text, nil)
	if match == nil {
		return []string{}, ErrMatchFuncNil
	}
	fmt.Println(lines)
	for i, l := range lines {
		// Once we find the standard text that Sapling adds, assumed that the rest of the lines are PR numbers
		if match(header, l) {
			// * #349
			// * __->__ #348
			return lines[i+1:], nil
		}
	}
	return []string{}, ErrStackLinksNotFound

}

func PRNumFromLine(line string) (*int, error) {
	n, err := strconv.Atoi(strings.SplitN(line, "#", 1)[1])
	if err != nil {
		return nil, err
	}
	return &n, nil
}

type PullRequest struct {
	Body     graphql.String `graphql:"body"`
	BodyHTML graphql.String `graphql:"bodyHTML"`
}

type Repository struct {
	PullRequest   PullRequest    `graphql:"pullRequest(number: $number)"`
	NameWithOwner graphql.String `graphql:"nameWithOwner"`
}

type RepoQuery struct {
	Repository Repository `graphql:"repository(owner: $owner, name: $name)"`
}

func Query(c api.GQLClient, owner, project string, prNumber int) (RepoQuery, error) {
	var err error
	var client api.GQLClient
	if c == nil {
		client, err = gh.GQLClient(nil)
		if err != nil {
			return RepoQuery{}, err
		}
	} else {
		client = c
	}
	var query RepoQuery
	variables := map[string]interface{}{
		"number": graphql.Int(prNumber),
		"owner":  graphql.String(owner),
		"name":   graphql.String(project),
	}
	err = client.Query("pullRequest", &query, variables)
	if err != nil {
		return RepoQuery{}, err
	}
	return query, nil
}

// Copied with <3 from https://groob.io/tutorial/go-github-webhook/
func RunWebhook(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
	}

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Printf("could not parse webhook: err=%s\n", err)
		return
	}

	switch e := event.(type) {
	case *github.PullRequestEvent:
		log.Print(e.PullRequest.GetURL())
	case *github.PingEvent:
		log.Print(e)
	default:
		log.Printf("unknown event type %s\n", github.WebHookType(r))
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Print(err)
	}
}
