package table

import (
	"fmt"
	"strings"
	"testing"

	"github.com/makkes/gitlab-cli/api"
)

func checkColumn(t *testing.T, cols map[string]int, col string, width int) {
	if cols[col] != width {
		t.Errorf("'%s' column has unexpected width %d", col, cols[col])
	}
}

func TestPipelineColumnWidths(t *testing.T) {
	var pipelineColumnWidthTests = []struct {
		name string
		in   []api.PipelineDetails
		out  map[string]int
	}{
		{
			"empty input",
			[]api.PipelineDetails{},
			map[string]int{
				"id":       20,
				"status":   20,
				"duration": 10,
				"url":      50,
			},
		},
		{
			"nil input",
			nil,
			map[string]int{
				"id":       20,
				"status":   20,
				"duration": 10,
				"url":      50,
			},
		},
		{
			"happy path",
			[]api.PipelineDetails{
				{
					ID:       99,
					Status:   "this is a status with more than 20 characters",
					URL:      "This is a uniform resource locator with more than 50 characters",
					Duration: 50,
				},
			},
			map[string]int{
				"id":       20,
				"status":   45,
				"url":      63,
				"duration": 10,
			},
		},
	}

	for _, tt := range pipelineColumnWidthTests {
		t.Run(tt.name, func(t *testing.T) {
			res := calcPipelineColumnWidths(tt.in)
			for k, v := range tt.out {
				checkColumn(t, res, k, v)
			}
		})
	}
}

func TestProjectColumnWidths(t *testing.T) {
	var projectColumnWidthTests = []struct {
		name string
		in   []api.Project
		out  map[string]int
	}{
		{
			"empty input",
			[]api.Project{},
			map[string]int{
				"id":   15,
				"name": 40,
				"url":  50,
			},
		},
		{
			"nil input",
			nil,
			map[string]int{
				"id":   15,
				"name": 40,
				"url":  50,
			},
		},
		{
			"happy path",
			[]api.Project{
				{
					ID:   99,
					Name: "this is a name with more than 40 characters",
					URL:  "This is a uniform resource locator with more than 50 characters",
				},
			},
			map[string]int{
				"id":   15,
				"name": 43,
				"url":  63,
			},
		},
	}

	for _, tt := range projectColumnWidthTests {
		t.Run(tt.name, func(t *testing.T) {
			res := calcProjectColumnWidths(tt.in)
			for k, v := range tt.out {
				checkColumn(t, res, k, v)
			}
		})
	}
}

func TestIssuesTable(t *testing.T) {
	var issuesTableTests = []struct {
		name   string
		writer *strings.Builder
		issues []api.Issue
		out    string
	}{
		{
			"empty input",
			&strings.Builder{},
			[]api.Issue{},
			"ID                   TITLE                          STATE      URL                                               \n",
		},
		{
			"nil input",
			&strings.Builder{},
			nil,
			"ID                   TITLE                          STATE      URL                                               \n",
		},
		{
			"happy path",
			&strings.Builder{},
			[]api.Issue{
				{
					ID:    99,
					Title: "this is a name with more than 40 characters",
					State: "this is a loong state",
					URL:   "This is a uniform resource locator with more than 50 characters",
				},
			},
			`ID                   TITLE                          STATE                 URL                                                            
0:99                 this is a name with more than… this is a loong state This is a uniform resource locator with more than 50 characters
`,
		},
	}

	for _, tt := range issuesTableTests {
		t.Run(tt.name, func(t *testing.T) {
			PrintIssues(tt.writer, tt.issues)
			if tt.writer.String() != tt.out {
				t.Errorf("'%s'", strings.Split(tt.writer.String(), "\n")[1])
				t.Errorf("Unexpected result: (%d) '%s'", len(tt.writer.String()), tt.writer.String())
			}
		})
	}
}

func TestPad(t *testing.T) {
	var padTable = []struct {
		s   string
		w   int
		out string
	}{
		{"", 0, ""},
		{"", 10, "          "},
		{"don't shorten me", 0, "don't shorten me"},
		{"i'm too short", 20, "i'm too short       "},
		{"not padded when negative", -100, "not padded when negative"},
	}

	for _, tt := range padTable {
		t.Run(fmt.Sprintf("'%s':%d", tt.s, tt.w), func(t *testing.T) {
			res := pad(tt.s, tt.w)
			if res != tt.out {
				t.Errorf("Expected '%s' to be '%s'", res, tt.out)
			}
		})
	}
}
