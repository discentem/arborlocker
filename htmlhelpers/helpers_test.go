package htmlhelpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	content = `<p dir="auto">util: add zsh prompt</p>
<p dir="auto">Summary: Add shell prompt for Sapling.  The current prompt only supports git<br>
and hg, but not sl.  Since the <code class="notranslate">hg</code> prompt should work the same for <code class="notranslate">sl</code>, use<br>
<code class="notranslate">_hg_prompt</code> for Sapling if a <code class="notranslate">.sl</code> directory exists</p>
<p dir="auto">Test Plan: TODO: copy <code class="notranslate">.hg</code> testcases from <code class="notranslate">eden/scm/tests/test-fb-ext-scm-prompt-hg.t</code></p>
<hr>
<p dir="auto">Stack created with <a href="https://sapling-scm.com" rel="nofollow">Sapling</a>. Best reviewed with <a href="https://reviewstack.dev/facebook/sapling/pull/348" rel="nofollow">ReviewStack</a>.</p>
<ul dir="auto">
<li><a class="issue-link js-issue-link" data-error-text="Failed to load title" data-id="1509816830" data-permission-text="Title is private" data-url="https://github.com/facebook/sapling/issues/349" data-hovercard-type="pull_request" data-hovercard-url="/facebook/sapling/pull/349/hovercard" href="https://github.com/facebook/sapling/pull/349">#349</a></li>
<li><strong>-&gt;</strong> <a class="issue-link js-issue-link" data-error-text="Failed to load title" data-id="1508364579" data-permission-text="Title is private" data-url="https://github.com/facebook/sapling/issues/348" data-hovercard-type="pull_request" data-hovercard-url="/facebook/sapling/pull/348/hovercard" href="https://github.com/facebook/sapling/pull/348">#348</a></li>
</ul>`
)

func TestGetHTMLSubset(t *testing.T) {
	tests := []struct {
		content string
		want    string
	}{
		{
			content: content,
			want: `<ul dir="auto">
<li><a class="issue-link js-issue-link" data-error-text="Failed to load title" data-id="1509816830" data-permission-text="Title is private" data-url="https://github.com/facebook/sapling/issues/349" data-hovercard-type="pull_request" data-hovercard-url="/facebook/sapling/pull/349/hovercard" href="https://github.com/facebook/sapling/pull/349">#349</a></li>
<li><strong>-&gt;</strong> <a class="issue-link js-issue-link" data-error-text="Failed to load title" data-id="1508364579" data-permission-text="Title is private" data-url="https://github.com/facebook/sapling/issues/348" data-hovercard-type="pull_request" data-hovercard-url="/facebook/sapling/pull/348/hovercard" href="https://github.com/facebook/sapling/pull/348">#348</a></li>
</ul>`,
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
