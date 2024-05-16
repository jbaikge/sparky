package pagination

import (
	"net/http"
	"slices"
	"testing"
)

func TestCurrent(t *testing.T) {
	tests := []struct {
		Name    string
		Path    string
		Current int
	}{
		{
			Name:    "landing page",
			Path:    "/test",
			Current: 1,
		},
	}

	for _, test := range tests {
		r, _ := http.NewRequest("GET", test.Path, nil)
		p := NewPagination(r)
		if got, expect := p.Current(), test.Current; got != expect {
			t.Fatalf("%s: %d != %d", test.Name, got, expect)
		}
	}
}

func TestPerPage(t *testing.T) {
	tests := []struct {
		Name  string
		Path  string
		Limit int
	}{
		{
			Name:  "default limit",
			Path:  "/test",
			Limit: defaultPerPage,
		},
		{
			Name:  "passed in limit",
			Path:  "/test?pp=42",
			Limit: 42,
		},
		{
			Name:  "negative limit",
			Path:  "/test?pp=-42",
			Limit: defaultPerPage,
		},
		{
			Name:  "invalid number",
			Path:  "/test?pp=forty-two",
			Limit: defaultPerPage,
		},
	}

	for _, test := range tests {
		r, _ := http.NewRequest("GET", test.Path, nil)
		p := NewPagination(r)
		if got, expect := p.PerPage(), test.Limit; got != expect {
			t.Fatalf("%s: %d != %d", test.Name, got, expect)
		}
	}
}

func TestOffset(t *testing.T) {
	tests := []struct {
		Name   string
		Path   string
		Offset int
	}{
		{
			Name:   "first page",
			Path:   "/test",
			Offset: 0,
		},
		{
			Name:   "page 3, default per-page",
			Path:   "/test?p=3",
			Offset: 2 * defaultPerPage,
		},
		{
			Name:   "page 3, custom per-page",
			Path:   "/test?p=3&pp=4",
			Offset: 2 * 4,
		},
	}

	for _, test := range tests {
		r, _ := http.NewRequest("GET", test.Path, nil)
		p := NewPagination(r)
		if got, expect := p.Offset(), test.Offset; got != expect {
			t.Fatalf("%s: %d != %d", test.Name, got, expect)
		}
	}
}

func TestWindow(t *testing.T) {
	tests := []struct {
		Name  string
		Path  string
		Total int
		Pages []int
	}{
		{
			Name:  "empty data",
			Path:  "/test",
			Total: 0,
			Pages: []int{},
		},
		{
			Name:  "one page",
			Path:  "/test",
			Total: defaultPerPage,
			Pages: []int{},
		},
		{
			Name:  "two pages",
			Path:  "/test",
			Total: 2 * defaultPerPage,
			Pages: []int{},
		},
		{
			Name:  "three pages",
			Path:  "/test",
			Total: 3 * defaultPerPage,
			Pages: []int{2},
		},
		{
			Name:  "landing page",
			Path:  "/test",
			Total: 20 * defaultPerPage,
			Pages: []int{2, 3, 4, 5, 6},
		},
		{
			Name:  "first page",
			Path:  "/test?p=1",
			Total: 20 * defaultPerPage,
			Pages: []int{2, 3, 4, 5, 6},
		},
		{
			Name:  "page 4",
			Path:  "/test?p=4",
			Total: 20 * defaultPerPage,
			Pages: []int{2, 3, 4, 5, 6},
		},
		{
			Name:  "page 5",
			Path:  "/test?p=5",
			Total: 20 * defaultPerPage,
			Pages: []int{3, 4, 5, 6, 7},
		},
		{
			Name:  "page 10",
			Path:  "/test?p=10",
			Total: 20 * defaultPerPage,
			Pages: []int{8, 9, 10, 11, 12},
		},
		{
			Name:  "page 15",
			Path:  "/test?p=15",
			Total: 20 * defaultPerPage,
			Pages: []int{13, 14, 15, 16, 17},
		},
		{
			Name:  "page 20",
			Path:  "/test?p=20",
			Total: 20 * defaultPerPage,
			Pages: []int{15, 16, 17, 18, 19},
		},
	}

	for _, test := range tests {
		r, _ := http.NewRequest("GET", test.Path, nil)
		p := NewPagination(r)
		p.SetTotal(test.Total)
		if got, expect := p.windowPages(), test.Pages; !slices.Equal(got, expect) {
			t.Fatalf("%s: %v != %v", test.Name, got, expect)
		}
	}
}

func TestLinksSinglePage(t *testing.T) {
	r, _ := http.NewRequest("GET", "/test", nil)
	p := NewPagination(r)
	p.SetTotal(defaultPerPage)
	if num := len(p.Links()); num != 0 {
		t.Fatalf("expected zero links, got %d", num)
	}
}

func TestLinksZeroResults(t *testing.T) {
	r, _ := http.NewRequest("GET", "/test", nil)
	p := NewPagination(r)
	p.SetTotal(0)
	if num := len(p.Links()); num != 0 {
		t.Fatalf("expected zero links, got %d", num)
	}
}

func TestLinksWindowed(t *testing.T) {
	r, _ := http.NewRequest("GET", "/test?p=6&pp=10", nil)
	p := NewPagination(r)
	p.SetTotal(200)
	links := p.Links()

	// Prev + 1 + gap + shoulder + 1 + shoulder + gap + last + next = 11
	if num := len(links); num != 11 {
		t.Fatalf("expected 11 links, got %d", num)
	}

	if !links[0].IsPrev {
		t.Fatal("links[0] is not a previous link")
	}
	if links[0].Page != 5 {
		t.Fatal("links[0] is not page 5")
	}
	if links[1].Page != 1 {
		t.Fatal("links[1] is not page 1")
	}
	if !links[2].IsGap {
		t.Fatal("links[2] is not a gap")
	}
	for i, link := range links[3:8] {
		if link.Page != i+4 {
			t.Fatalf("links[%d] is not page %d", i+3, i+4)
		}
	}
	if !links[8].IsGap {
		t.Fatal("links[8] is not a gap")
	}
	if links[9].Page != 20 {
		t.Fatal("links[0] is not page 20")
	}
	if !links[10].IsNext {
		t.Fatal("links[10] is not a next link")
	}
}
