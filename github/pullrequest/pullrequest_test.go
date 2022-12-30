package pullrequest

import (
	"io"
	"net/http"
	"testing"

	"github.com/cli/go-gh/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
	fakegraphql "github.com/discentem/arborlocker/graphql"
	"github.com/stretchr/testify/assert"
)

func TestLinesFromTextDescription(t *testing.T) {
	tests := []struct {
		text         string
		ignoreLineFn func(string) bool
		want         []string
	}{
		{
			text: `abort if no args provided to sl pr pull
	---
	Stack created with [Sapling](https://sapling-scm.com). Best reviewed with [ReviewStack](https://reviewstack.dev/facebook/sapling/pull/357).
	* __->__ #357`,
			// nil causes LinesFromTextDescription to use default ignoreLine function
			ignoreLineFn: nil,
			want: []string{
				"abort if no args provided to sl pr pull",
				"\tStack created with [Sapling](https://sapling-scm.com). Best reviewed with [ReviewStack](https://reviewstack.dev/facebook/sapling/pull/357).",
				"\t* __->__ #357",
			},
		},
	}

	for _, test := range tests {
		got := LinesFromTextDescription(test.text, test.ignoreLineFn)
		assert.Equal(t, test.want, got)
	}

}

func TestLinesFromHTMLDescription(t *testing.T) {
	tests := []struct {
		content string
		want    []string
	}{
		{
			content: `<p dir=\"auto\">util: add zsh prompt</p>\n<p dir=\"auto\">Summary: Add shell prompt for Sapling.  The current prompt only supports git<br>\nand hg, but not sl.  Since the <code class=\"notranslate\">hg</code> prompt should work the same for <code class=\"notranslate\">sl</code>, use<br>\n<code class=\"notranslate\">_hg_prompt</code> for Sapling if a <code class=\"notranslate\">.sl</code> directory exists</p>\n<p dir=\"auto\">Test Plan: TODO: copy <code class=\"notranslate\">.hg</code> testcases from <code class=\"notranslate\">eden/scm/tests/test-fb-ext-scm-prompt-hg.t</code></p>\n<hr>\n<p dir=\"auto\">Stack created with <a href=\"https://sapling-scm.com\" rel=\"nofollow\">Sapling</a>. Best reviewed with <a href=\"https://reviewstack.dev/facebook/sapling/pull/348\" rel=\"nofollow\">ReviewStack</a>.</p>\n<ul dir=\"auto\">\n<li><a class=\"issue-link js-issue-link\" data-error-text=\"Failed to load title\" data-id=\"1509816830\" data-permission-text=\"Title is private\" data-url=\"https://github.com/facebook/sapling/issues/349\" data-hovercard-type=\"pull_request\" data-hovercard-url=\"/facebook/sapling/pull/349/hovercard\" href=\"https://github.com/facebook/sapling/pull/349\">#349</a></li>\n<li><strong>-&gt;</strong> <a class=\"issue-link js-issue-link\" data-error-text=\"Failed to load title\" data-id=\"1508364579\" data-permission-text=\"Title is private\" data-url=\"https://github.com/facebook/sapling/issues/348\" data-hovercard-type=\"pull_request\" data-hovercard-url=\"/facebook/sapling/pull/348/hovercard\" href=\"https://github.com/facebook/sapling/pull/348\">#348</a></li>\n</ul>`,
			want: []string{
				"https://github.com/facebook/sapling/pull/349",
				"https://github.com/facebook/sapling/pull/348",
			},
		},
	}
	for _, test := range tests {
		lines, err := LinesFromHTMLDescription(test.content)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, test.want, lines)
	}
}

func TestQuery(t *testing.T) {
	tests := []struct {
		client api.GQLClient
		number int
		owner  string
		name   string
		want   HTMLBodyQuery
	}{
		{
			client: fakegraphql.NewFakeGQLClient("blah", func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, err := io.WriteString(w, `{
					"data": {
						"repository": {
							"nameWithOwner":"Acme/BigProject",
							"pullRequest": {
								"bodyHTML":"blah"
							}
						}
					}
				}`)
				if err != nil {
					t.Error(err)
				}
			}),
			number: 1,
			owner:  "Acme",
			name:   "BigProject",
			want: HTMLBodyQuery{
				Repository: Repository{
					NameWithOwner: *graphql.NewString("Acme/BigProject"),
					PullRequest: PullRequest{
						BodyHTML: *graphql.NewString("blah"),
					},
				},
			},
		},
	}
	for _, test := range tests {
		got, err := Query(test.client, test.owner, test.name, test.number)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, test.want, got)

	}
}
