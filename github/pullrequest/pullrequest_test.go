package pullrequest

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/cli/go-gh/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
	fakegraphql "github.com/discentem/arborlocker/graphql"
	"github.com/stretchr/testify/assert"
)

const descriptionText = `util: add zsh prompt` + "\n" +
	`Summary: Add shell prompt for Sapling.  The current prompt only supports git
and hg, but not sl.  Since the` + "`" + "hg" + "`" + " prompt should work the same for " + "`" + "sl" + "`" + `, use
` + "`" + "_hg_prompt" + "`" + "for Sapling if a " + "`" + ".sl" + "`" + " directory exists" + "\n" +
	`Test Plan: TODO: copy ` + "`" + `.hg` + "`" + ` testcases from ` + "`" + `eden/scm/tests/test-fb-ext-scm-prompt-hg.t` + "`" + `
Stack created with [Sapling](https://sapling-scm.com). Best reviewed with [ReviewStack](https://reviewstack.dev/facebook/sapling/pull/348).
	* #349
	* __->__ #348`

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
		{
			text: descriptionText,
			want: []string{
				"util: add zsh prompt",
				"Summary: Add shell prompt for Sapling.  The current prompt only supports git",
				"and hg, but not sl.  Since the`hg` prompt should work the same for `sl`, use",
				"`_hg_prompt`for Sapling if a `.sl` directory exists",
				"Test Plan: TODO: copy `.hg` testcases from `eden/scm/tests/test-fb-ext-scm-prompt-hg.t`",
				"Stack created with [Sapling](https://sapling-scm.com). Best reviewed with [ReviewStack](https://reviewstack.dev/facebook/sapling/pull/348).",
				"\t* #349",
				"\t* __->__ #348"},
		},
	}
	for _, test := range tests {
		got := LinesFromTextDescription(test.text, test.ignoreLineFn)
		assert.Equal(t, test.want, got)
	}

}

func TestStackLines(t *testing.T) {
	table := []struct {
		prefix  string
		text    string
		matchFn func(header, line string) bool
		want    []string
		wantErr error
	}{
		{
			prefix: "Stack created with [Sapling](https://sapling-scm.com). Best reviewed with [ReviewStack]",
			text:   descriptionText,
			matchFn: func(header, line string) bool {
				return strings.HasPrefix(line, header)
			},
			want: []string{
				"\t* #349",
				"\t* __->__ #348",
			},
		},
	}
	for _, test := range table {
		links, err := StackLines(test.prefix, test.matchFn, test.text)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, test.want, links)
	}
}

func TestPRNumFromLine(t *testing.T) {
	table := []struct {
		line string
		want int
	}{
		{
			line: "\t* #349",
			want: 349,
		},
	}
	for _, test := range table {
		num, err := PRNumFromLine(test.line)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, test.want, num)
	}
}

func TestPRListFromLines(t *testing.T) {
	table := []struct {
		lines []string
		want  []int
	}{
		{
			lines: []string{
				"\t* #349",
				"\t* #348",
			},
			want: []int{
				349,
				348,
			},
		},
	}
	for _, test := range table {
		nums, err := PRListFromLines(test.lines)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, test.want, nums)
	}
}

func TestQuery(t *testing.T) {
	tests := []struct {
		client api.GQLClient
		number int
		owner  string
		name   string
		want   RepoQuery
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
			want: RepoQuery{
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
