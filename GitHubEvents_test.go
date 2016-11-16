package main

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAsJira(t *testing.T) {
	list := makeEventListTestData()
	buf := new(bytes.Buffer)
	list.AsJira(buf)
	assert.Contains(t, buf.String(), "|Hello|")
	assert.Contains(t, buf.String(), "|tuser|Hello|http://somewhere.com|")
}

func TestAsCSV(t *testing.T) {
	list := makeEventListTestData()
	buf := new(bytes.Buffer)
	list.AsCSV(buf)
	assert.Contains(t, buf.String(), "Hello,")
	assert.Contains(t, buf.String(), ",tuser,Hello,http://somewhere.com")
}

func TestAsHTML(t *testing.T) {
	list := makeEventListTestData()
	buf := new(bytes.Buffer)
	list.AsHTML(buf)
	assert.Contains(t, buf.String(), "<td>tuser</td><td>Hello</td><td><a href=\"http://somewhere.com\">http://somewhere.com</a></td></tr>")
}

func TestAsJSON(t *testing.T) {
	list := makeEventListTestData()
	buf := new(bytes.Buffer)
	list.AsJSON(buf)
	assert.Contains(t, buf.String(), "\"Login\":\"tuser\",\"Title\":\"Hello\",\"URL\":\"http://somewhere.com\"}]")
}

func TestAsText(t *testing.T) {
	list := makeEventListTestData()
	buf := new(bytes.Buffer)
	list.AsText(buf)
	assert.Contains(t, buf.String(), "| Hello | \x1b[32m")
	assert.Contains(t, buf.String(), "\x1b[0m | tuser  | Hello | http://somewhere.com |")
}

func makeEventListTestData() EventList {
	repo := "Hello"
	created := time.Now()
	login := "tuser"
	url := "http://somewhere.com"
	list := EventList{
		summaries: &[]EventSummary{
			EventSummary{
				Repository: &repo,
				LastUsed:   &created,
				Login:      &login,
				Title:      &repo,
				URL:        &url,
			},
		},
	}
	return list
}