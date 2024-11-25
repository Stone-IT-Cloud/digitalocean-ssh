package droplets

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/digitalocean/godo"
	"github.com/jedib0t/go-pretty/v6/table"
)

var authenticate = func () (*godo.Client, error) {
	token := os.Getenv("DIGITALOCEAN_ACCESS_TOKEN")
	if token == "" {
		return nil, errors.New("the environment variable DIGITALOCEAN_ACCESS_TOKEN is not set")
	}
	client := godo.NewFromToken(token)
	return client, nil
}

// ListDroplets lists all droplets in the DigitalOcean account and prints them in a table format.
func ListDroplets() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleColoredDark)

	client, err := authenticate()
	if err != nil {
		log.Fatalf("Error authenticating: %v", err)
	}

	d, _, err := client.Droplets.List(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to retrieve droplets from DigitalOcean API: %v", err)
	}

	t.AppendHeader(table.Row{"id", "name", "region", "image_name", "image_id", "size", "price $/mo", "public_addr", "private_addr"})
	


	for _, droplet := range d {
		dropletID := droplet.ID
		dropletName := droplet.Name
		dropletRegion := droplet.Region.Slug
		imageName := droplet.Image.Name
		imageID := droplet.Image.ID
		dropletSize := droplet.Size.Slug
		priceMonthly := droplet.Size.PriceMonthly
		privateAddr := ""
		publicAddr := ""
		if len(droplet.Networks.V4) > 0 {
			privateAddr = droplet.Networks.V4[0].IPAddress
		}
		if len(droplet.Networks.V4) >= 2 {
			publicAddr = droplet.Networks.V4[1].IPAddress
		}
				
		t.AppendRow([]interface{}{dropletID, dropletName, dropletRegion, imageName, imageID, dropletSize, fmt.Sprintf("$%.2f", priceMonthly), publicAddr, privateAddr})
		t.Render()
	}
}