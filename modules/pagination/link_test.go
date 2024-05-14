package pagination

import "testing"

func TestActive(t *testing.T) {
	link := Link{
		Page: 1,
	}
	if !link.IsActive(1) {
		t.Fatal("expected link to be active")
	}
	if link.IsActive(2) {
		t.Fatal("expected link to be inactive")
	}
}

func TestActiveLabel(t *testing.T) {
	link := Link{
		Page:   1,
		IsPrev: true,
	}
	if link.IsActive(1) {
		t.Fatal("expected link to be inactive")
	}
	if link.IsActive(2) {
		t.Fatal("expected link to be inactive")
	}
}

func TestDisabled(t *testing.T) {
	link := Link{
		Page: 10,
	}
	if link.IsDisabled(10) {
		t.Fatal("expected link to be enabled")
	}
	if link.IsDisabled(2) {
		t.Fatal("expected link to be enabled")
	}
}

// Gap set to true will always cause links to be disabled
func TestGap(t *testing.T) {
	link := Link{
		IsGap: true,
	}
	if !link.IsDisabled(10) {
		t.Fatal("expected link to be enabled")
	}
	if !link.IsDisabled(2) {
		t.Fatal("expected link to be enabled")
	}
}
