package pagination

import (
	"math"
	"net/http"
	"strconv"
)

const (
	defaultPerPage  = 50
	defaultShoulder = 2
)

type Pagination struct {
	request   *http.Request
	shoulder  int
	totalRows int
}

func NewPagination(r *http.Request) (p *Pagination) {
	p = &Pagination{
		request:  r,
		shoulder: defaultShoulder,
	}
	return
}

func (p Pagination) Current() int {
	page, _ := strconv.Atoi(p.request.FormValue("p"))
	if page < 1 {
		page = 1
	}
	return page
}

func (p Pagination) Links() (links []Link) {
	lastPage := p.pageCount()
	if lastPage == 1 {
		return
	}

	currentPage := p.Current()
	query := p.request.URL.Query()

	// Previous, 1, gap, windowPages, gap, lastPage, Next
	links = make([]Link, 0, lastPage+4)

	// Previous
	prevPage := currentPage - 1
	if prevPage < 1 {
		prevPage = 1
	}
	links = append(links, Link{
		Page:   prevPage,
		IsPrev: true,
		Query:  query,
	})

	// Page 1
	links = append(links, Link{
		Page:  1,
		Query: query,
	})

	windowPages := p.windowPages()

	// Gap, if necessary
	if windowPages[0] > 2 {
		links = append(links, Link{
			IsGap: true,
		})
	}

	// Window pages
	for _, page := range windowPages {
		links = append(links, Link{
			Page:  page,
			Query: query,
		})
	}

	// Gap, if necessary
	if windowPages[len(windowPages)-1]+1 < lastPage {
		links = append(links, Link{
			IsGap: true,
			Query: query,
		})
	}

	// Last page
	links = append(links, Link{
		Page:  lastPage,
		Query: query,
	})

	// Next
	nextPage := currentPage + 1
	if nextPage > lastPage {
		nextPage = lastPage
	}
	links = append(links, Link{
		Page:   nextPage,
		IsNext: true,
		Query:  query,
	})

	return
}

func (p Pagination) Offset() int {
	return (p.Current() - 1) * p.PerPage()
}

func (p Pagination) PerPage() int {
	perPage, _ := strconv.Atoi(p.request.FormValue("pp"))
	if perPage < 1 {
		perPage = defaultPerPage
	}
	return perPage
}

func (p *Pagination) SetTotal(rows int) {
	p.totalRows = rows
}

func (p Pagination) pageCount() int {
	return int(math.Ceil(float64(p.totalRows) / float64(p.PerPage())))
}

// Window Pages are the pages between 1 and the last page number, centered
// on the current page when the current page is between 1 and the last page.
//
// This is how you get pagination that looks like:
// 1 ... 3 4 5 6 7 ... 10
//
// The shoulder defines how many pages on each side of the current page to
// display.
func (p Pagination) windowPages() (pages []int) {
	pages = make([]int, 0, p.shoulder*2+1)

	lastPage := p.pageCount()
	if lastPage <= 1 {
		return
	}

	// The easy one: all pages fit inside the window
	if p.shoulder*2+3 >= lastPage {
		for page := 2; page < lastPage; page++ {
			pages = append(pages, page)
		}
		return
	}

	// The tough one: create a window centered around the current page
	currentPage := p.Current()
	minPage := 2
	maxPage := lastPage - 1

	// Extend the lower bound to the left by the shoulder width
	if lower := currentPage - p.shoulder; lower > minPage {
		minPage = lower
	}

	// Extend lower bound further if the upper bound is up against the last page
	if lower := maxPage - p.shoulder*2; lower < minPage {
		minPage = lower
	}

	// Extend the upper bound to the right by the shoulder width
	if upper := currentPage + p.shoulder; upper < maxPage {
		maxPage = upper
	}

	// Extend upper bound further if the lower bound is up against the page 1.
	if upper := minPage + p.shoulder*2; upper > maxPage {
		maxPage = upper
	}

	for page := minPage; page <= maxPage; page++ {
		pages = append(pages, page)
	}

	return pages
}
