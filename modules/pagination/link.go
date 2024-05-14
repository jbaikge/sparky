package pagination

import (
	"fmt"
	"net/url"
)

type Link struct {
	Page   int
	IsGap  bool
	IsPrev bool
	IsNext bool
	Query  url.Values
}

// Previous, Next, and gap links can never be active
func (link Link) IsActive(page int) bool {
	if link.IsPrev || link.IsNext || link.IsGap {
		return false
	}
	return link.Page == page
}

// Disable Previous when on page 1
// Disable Next when on the last page
// Gaps are always disabled
func (link Link) IsDisabled(page int) bool {
	if link.IsPrev && link.Page == page {
		return true
	}
	if link.IsNext && link.Page == page {
		return true
	}
	if link.IsGap {
		return true
	}
	return false
}

func (link Link) IsNumber() bool {
	return !link.IsGap && !link.IsPrev && !link.IsNext
}

func (link Link) URL() string {
	if link.Query == nil {
		return fmt.Sprintf("?p=%d", link.Page)
	}

	link.Query.Set("p", fmt.Sprint(link.Page))
	return "?" + link.Query.Encode()
}
