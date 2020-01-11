package data

import "fmt"

type (
	// Spot - a simple spot type
	Spot struct {
		// ID - The spot identifier
		ID string

		// OpenDate - The available date of the spot
		OpenDate string

		// RegDate - The reg date of the spot
		RegDate string

		// RegisteredBy - The user who registered the spot
		RegisteredBy string
	}

)

// Key - the key for this spot
func (s Spot) Key() string {
	return fmt.Sprintf("%v-%s", s.ID, s.OpenDate)
}

// IsZeroValue - returns true if all elements of the struct are their zero-value. This is primarily used to make testing easier.
func (s Spot) IsZeroValue() bool {
	return s.ID == "" && s.OpenDate == "" && s.RegDate == "" && s.RegisteredBy == ""
}
