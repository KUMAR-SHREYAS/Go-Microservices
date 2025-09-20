package data

import "testing"

func TestValidationChecks(t *testing.T) { // testing.T here is a testing stats
	p := &Product{
		Name:  "nics",
		Price: 1.00,
		SKU:   "abs-bec-def",
		// SKU:   "abs",
	}
	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
