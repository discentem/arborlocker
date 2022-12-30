// lint:file-ignore U1000 repository

package pullrequest

import (
	"strings"

	"github.com/discentem/arborlocker/htmlhelpers"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
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

type PRQuery struct {
	Repository struct {
		PullRequest struct {
			BodyHTML graphql.String `graphql:"bodyHTML"`
		} `graphql:"pullRequest(number: $number)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func Query(c *api.GQLClient, prNumber int) (PRQuery, error) {
	var err error
	var client api.GQLClient
	if c == nil {
		client, err = gh.GQLClient(nil)
		if err != nil {
			return PRQuery{}, err
		}
	} else {
		client = *c
	}
	var query PRQuery
	variables := map[string]interface{}{
		"number": graphql.Int(prNumber),
		"owner":  graphql.String("facebook"),
		"name":   graphql.String("sapling"),
	}
	err = client.Query("pullRequest", &query, variables)
	if err != nil {
		return PRQuery{}, err
	}
	return query, nil
}
