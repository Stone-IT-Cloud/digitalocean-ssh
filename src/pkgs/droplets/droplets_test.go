package droplets

import (
	"context"
	"os"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	t.Run("missing token", func(t *testing.T) {
		os.Unsetenv("DIGITALOCEAN_ACCESS_TOKEN")
		client, err := authenticate()
		
		assert.Nil(t, client)
		assert.EqualError(t, err, "the environment variable DIGITALOCEAN_ACCESS_TOKEN is not set")
	})

	t.Run("valid token", func(t *testing.T) {
		os.Setenv("DIGITALOCEAN_ACCESS_TOKEN", "invalid-token")
		client, err := authenticate()
		assert.NotNil(t, client)
		assert.NoError(t, err)
		assert.IsType(t, &godo.Client{}, client)
		
		// DO does not validates the token on authentication, it just configures the barer in the client
		accnt, _, err := client.Account.Get(context.Background())
		assert.NotNil(t, accnt)
		assert.NoError(t, err)
	})
}