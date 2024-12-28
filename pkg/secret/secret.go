package secret

import (
	"context"
)

// Client is a client that stores secret values and provides
// a way to fetch those values.
//
// All clients that implement this interface should handle
// authentication and authorization on their own. For example
// in the initialization function.
//
// For user, it is recommended to cache the fetched secret
// values to reduce fetching cost.
type Client interface {
	// Fetch fetches secret values and return it in bytes.
	Fetch(ctx context.Context, path string) ([]byte, error)
}
