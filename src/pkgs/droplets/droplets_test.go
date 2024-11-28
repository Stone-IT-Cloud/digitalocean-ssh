package droplets

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	t.Run("missing token", func(t *testing.T) {
		token := os.Getenv("DIGITALOCEAN_ACCESS_TOKEN")
		os.Unsetenv("DIGITALOCEAN_ACCESS_TOKEN")
		client, err := authenticate()

		assert.Nil(t, client)
		assert.EqualError(t, err, "the environment variable DIGITALOCEAN_ACCESS_TOKEN is not set")
		os.Setenv("DIGITALOCEAN_ACCESS_TOKEN", token)
	})

	t.Run("valid token", func(t *testing.T) {
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

func (m *MockDropletsService) Actions(ctx context.Context, dropletID int, opt *godo.ListOptions) ([]godo.Action, *godo.Response, error) {
	return nil, nil, nil
}

type MockDropletsService struct {
	ListFn                        func(ctx context.Context, opt *godo.ListOptions) ([]godo.Droplet, *godo.Response, error)
	ListWithGPUsFn                func(context.Context, *godo.ListOptions) ([]godo.Droplet, *godo.Response, error)
	ListByNameFn                  func(context.Context, string, *godo.ListOptions) ([]godo.Droplet, *godo.Response, error)
	ListByTagFn                   func(context.Context, string, *godo.ListOptions) ([]godo.Droplet, *godo.Response, error)
	GetFn                         func(context.Context, int) (*godo.Droplet, *godo.Response, error)
	CreateFn                      func(context.Context, *godo.DropletCreateRequest) (*godo.Droplet, *godo.Response, error)
	CreateMultipleFn              func(context.Context, *godo.DropletMultiCreateRequest) ([]godo.Droplet, *godo.Response, error)
	DeleteFn                      func(context.Context, int) (*godo.Response, error)
	DeleteByTagFn                 func(context.Context, string) (*godo.Response, error)
	KernelsFn                     func(context.Context, int, *godo.ListOptions) ([]godo.Kernel, *godo.Response, error)
	SnapshotsFn                   func(context.Context, int, *godo.ListOptions) ([]godo.Image, *godo.Response, error)
	BackupsFn                     func(context.Context, int, *godo.ListOptions) ([]godo.Image, *godo.Response, error)
	ActionsFn                     func(context.Context, int, *godo.ListOptions) ([]godo.Action, *godo.Response, error)
	NeighborsFn                   func(context.Context, int) ([]godo.Droplet, *godo.Response, error)
	GetBackupPolicyFn             func(context.Context, int) (*godo.DropletBackupPolicy, *godo.Response, error)
	ListBackupPoliciesFn          func(context.Context, *godo.ListOptions) (map[int]*godo.DropletBackupPolicy, *godo.Response, error)
	ListSupportedBackupPoliciesFn func(context.Context) ([]*godo.SupportedBackupPolicy, *godo.Response, error)
}

func (m *MockDropletsService) List(ctx context.Context, opt *godo.ListOptions) ([]godo.Droplet, *godo.Response, error) {
	return m.ListFn(ctx, opt)
}
func (m *MockDropletsService) ListWithGPUs(context.Context, *godo.ListOptions) ([]godo.Droplet, *godo.Response, error) {
	panic("not implemented")
}
func (m *MockDropletsService) ListByName(context.Context, string, *godo.ListOptions) ([]godo.Droplet, *godo.Response, error) {
	panic("not implemented")
}
func (m *MockDropletsService) ListByTag(context.Context, string, *godo.ListOptions) ([]godo.Droplet, *godo.Response, error) {
	panic("not implemented")
}
func (m *MockDropletsService) Get(context.Context, int) (*godo.Droplet, *godo.Response, error) {
	panic("not implemented")
}
func (m *MockDropletsService) Create(context.Context, *godo.DropletCreateRequest) (*godo.Droplet, *godo.Response, error) {
	panic("not implemented")
}
func (m *MockDropletsService) CreateMultiple(context.Context, *godo.DropletMultiCreateRequest) ([]godo.Droplet, *godo.Response, error) {
	panic("not implemented")
}
func (m *MockDropletsService) Delete(context.Context, int) (*godo.Response, error) {
	panic("not implemented")
}
func (m *MockDropletsService) DeleteByTag(context.Context, string) (*godo.Response, error) {
	panic("not implemented")
}
func (m *MockDropletsService) Kernels(context.Context, int, *godo.ListOptions) ([]godo.Kernel, *godo.Response, error) {
	panic("not implemented")
}
func (m *MockDropletsService) Snapshots(context.Context, int, *godo.ListOptions) ([]godo.Image, *godo.Response, error) {
	panic("not implemented")
}
func (m *MockDropletsService) Backups(context.Context, int, *godo.ListOptions) ([]godo.Image, *godo.Response, error) {
	panic("not implemented")
}

/*
	 func (m *MockDropletsService) Actions(context.Context, int, *godo.ListOptions) ([]godo.Action, *godo.Response, error) {
		panic("not implemented")
	}
*/
func (m *MockDropletsService) Neighbors(context.Context, int) ([]godo.Droplet, *godo.Response, error) {
	panic("not implemented")
}
func (m *MockDropletsService) GetBackupPolicy(context.Context, int) (*godo.DropletBackupPolicy, *godo.Response, error) {
	panic("not implemented")
}
func (m *MockDropletsService) ListBackupPolicies(context.Context, *godo.ListOptions) (map[int]*godo.DropletBackupPolicy, *godo.Response, error) {
	panic("not implemented")
}
func (m *MockDropletsService) ListSupportedBackupPolicies(context.Context) ([]*godo.SupportedBackupPolicy, *godo.Response, error) {
	panic("not implemented")
}
func TestListDroplets(t *testing.T) {
	// Mock the authenticate function
	originalAuthenticate := authenticate
	defer func() { authenticate = originalAuthenticate }()
	authenticate = func() (*godo.Client, error) {
		return &godo.Client{}, nil
	}

	// Mock the Droplets.List function
	mockDroplets := []godo.Droplet{
		{
			ID:     123,
			Name:   "test-droplet",
			Region: &godo.Region{Slug: "nyc3"},
			Image:  &godo.Image{ID: 456, Name: "ubuntu-20-04-x64"},
			Size:   &godo.Size{Slug: "s-1vcpu-1gb", PriceMonthly: 5.00},
			Networks: &godo.Networks{
				V4: []godo.NetworkV4{
					{IPAddress: "192.168.1.1"},
					{IPAddress: "203.0.113.1"},
				},
			},
		},
	}

	client := &godo.Client{
		Droplets: &MockDropletsService{
			ListFn: func(ctx context.Context, opt *godo.ListOptions) ([]godo.Droplet, *godo.Response, error) {
				return mockDroplets, nil, nil
			},
		},
	}

	authenticate = func() (*godo.Client, error) {
		return client, nil
	}

	// Capture the output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	ListDroplets()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()

	// Check if the output contains the expected table headers and data
	assert.Contains(t, output, "ID")
	assert.Contains(t, output, "NAME")
	assert.Contains(t, output, "REGION")
	assert.Contains(t, output, "IMAGE_NAME")
	assert.Contains(t, output, "IMAGE_ID")
	assert.Contains(t, output, "SIZE")
	assert.Contains(t, output, "PRICE $/MO")
	assert.Contains(t, output, "PUBLIC_ADDR")
	assert.Contains(t, output, "PRIVATE_ADDR")

	assert.Contains(t, output, "123")
	assert.Contains(t, output, "test-droplet")
	assert.Contains(t, output, "nyc3")
	assert.Contains(t, output, "ubuntu-20-04-x64")
	assert.Contains(t, output, "456")
	assert.Contains(t, output, "s-1vcpu-1gb")
	assert.Contains(t, output, "$5.00")
	assert.Contains(t, output, "203.0.113.1")
	assert.Contains(t, output, "192.168.1.1")
}
