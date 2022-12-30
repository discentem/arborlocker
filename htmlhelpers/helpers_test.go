package htmlhelpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	content = `<p dir=\"auto\">util: add zsh prompt</p>\n<p dir=\"auto\">Summary: Add shell prompt for Sapling.  The current prompt only supports git<br>\nand hg, but not sl.  Since the <code class=\"notranslate\">hg</code> prompt should work the same for <code class=\"notranslate\">sl</code>, use<br>\n<code class=\"notranslate\">_hg_prompt</code> for Sapling if a <code class=\"notranslate\">.sl</code> directory exists</p>\n<p dir=\"auto\">Test Plan: TODO: copy <code class=\"notranslate\">.hg</code> testcases from <code class=\"notranslate\">eden/scm/tests/test-fb-ext-scm-prompt-hg.t</code></p>\n<hr>\n<p dir=\"auto\">Stack created with <a href=\"https://sapling-scm.com\" rel=\"nofollow\">Sapling</a>. Best reviewed with <a href=\"https://reviewstack.dev/facebook/sapling/pull/348\" rel=\"nofollow\">ReviewStack</a>.</p>\n<ul dir=\"auto\">\n<li><a class=\"issue-link js-issue-link\" data-error-text=\"Failed to load title\" data-id=\"1509816830\" data-permission-text=\"Title is private\" data-url=\"https://github.com/facebook/sapling/issues/349\" data-hovercard-type=\"pull_request\" data-hovercard-url=\"/facebook/sapling/pull/349/hovercard\" href=\"https://github.com/facebook/sapling/pull/349\">#349</a></li>\n<li><strong>-&gt;</strong> <a class=\"issue-link js-issue-link\" data-error-text=\"Failed to load title\" data-id=\"1508364579\" data-permission-text=\"Title is private\" data-url=\"https://github.com/facebook/sapling/issues/348\" data-hovercard-type=\"pull_request\" data-hovercard-url=\"/facebook/sapling/pull/348/hovercard\" href=\"https://github.com/facebook/sapling/pull/348\">#348</a></li>\n</ul>`
)

func TestGetHTMLSubset(t *testing.T) {
	tests := []struct {
		content string
		want    string
	}{
		{
			content: content,
			want:    "<ul dir=\"\\&#34;auto\\&#34;\">\\n<li><a class=\"\\&#34;issue-link\" js-issue-link\\\"=\"\" data-error-text=\"\\&#34;Failed\" to=\"\" load=\"\" title\\\"=\"\" data-id=\"\\&#34;1509816830\\&#34;\" data-permission-text=\"\\&#34;Title\" is=\"\" private\\\"=\"\" data-url=\"\\&#34;https://github.com/facebook/sapling/issues/349\\&#34;\" data-hovercard-type=\"\\&#34;pull_request\\&#34;\" data-hovercard-url=\"\\&#34;/facebook/sapling/pull/349/hovercard\\&#34;\" href=\"\\&#34;https://github.com/facebook/sapling/pull/349\\&#34;\">#349</a></li>\\n<li><strong>-&gt;</strong> <a class=\"\\&#34;issue-link\" js-issue-link\\\"=\"\" data-error-text=\"\\&#34;Failed\" to=\"\" load=\"\" title\\\"=\"\" data-id=\"\\&#34;1508364579\\&#34;\" data-permission-text=\"\\&#34;Title\" is=\"\" private\\\"=\"\" data-url=\"\\&#34;https://github.com/facebook/sapling/issues/348\\&#34;\" data-hovercard-type=\"\\&#34;pull_request\\&#34;\" data-hovercard-url=\"\\&#34;/facebook/sapling/pull/348/hovercard\\&#34;\" href=\"\\&#34;https://github.com/facebook/sapling/pull/348\\&#34;\">#348</a></li>\\n</ul>",
		},
	}
	for _, test := range tests {
		got, err := getHTMLSubset("ul", test.content)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, test.want, got)
	}
}

func TestGetLinks(t *testing.T) {
	tests := []struct {
		content string
		want    []string
	}{
		{
			content: content,
			want: []string{
				"https://github.com/facebook/sapling/pull/349",
				"https://github.com/facebook/sapling/pull/348",
			},
		},
	}
	for _, test := range tests {
		got, err := GetLinks(test.content)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, test.want, got)
	}
}
