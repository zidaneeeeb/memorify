// localfile defines secret client that stores secret values
// in local file.
package localfile

import (
	"context"
	"io/ioutil"
)

// client implements secret.Client.
type client struct{}

// New returns a new client.
func New() *client {
	return &client{}
}

// Fetch fetches secret values and return it in bytes.
func (c *client) Fetch(ctx context.Context, path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
