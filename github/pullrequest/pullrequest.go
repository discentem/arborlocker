// lint:file-ignore U1000 repository

package pullrequest

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/discentem/arborlocker/htmlhelpers"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
	"github.com/google/go-github/github"
)

// LinesFromHTMLDescription parses t, assumed to be HTML text from a Github Pull Request description.
func LinesFromHTMLDescription(content string) ([]string, error) {
	return htmlhelpers.GetLinks(content)
}

// LinesFromTextDescription parses t, assumed to be a plain text from a Github Pull Request description,
func LinesFromTextDescription(t string, ignoreLineFn func(string) bool) []string {
	lines := strings.Split(t, "\n")
	flines := []string{}

	if ignoreLineFn == nil {
		ignoreLineFn = func(t string) bool {
			return t == "\t---" || t == "\t" || t == "\t\t\t" || t == "---"
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

type PullRequest struct {
	BodyHTML graphql.String `graphql:"bodyHTML"`
}

type Repository struct {
	PullRequest   PullRequest    `graphql:"pullRequest(number: $number)"`
	NameWithOwner graphql.String `graphql:"nameWithOwner"`
}

type HTMLBodyQuery struct {
	Repository Repository `graphql:"repository(owner: $owner, name: $name)"`
}

// test
func Query(c api.GQLClient, owner, project string, prNumber int) (HTMLBodyQuery, error) {
	var err error
	var client api.GQLClient
	if c == nil {
		client, err = gh.GQLClient(nil)
		if err != nil {
			return HTMLBodyQuery{}, err
		}
	} else {
		client = c
	}
	var query HTMLBodyQuery
	variables := map[string]interface{}{
		"number": graphql.Int(prNumber),
		"owner":  graphql.String(owner),
		"name":   graphql.String(project),
	}
	err = client.Query("pullRequest", &query, variables)
	if err != nil {
		return HTMLBodyQuery{}, err
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
		log.Print(*e.PullRequest)
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
