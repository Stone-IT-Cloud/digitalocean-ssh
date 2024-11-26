package droplets

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"sort"

	"github.com/digitalocean/godo"
	"github.com/jedib0t/go-pretty/v6/table"
)

func isPrivateIp(ip string) bool {
//checks if ip is a private IP address or loopback and return true in that case
//https://en.wikipedia.org/wiki/Private_network
		// Check if the IP is a loopback address
		if ip == "127.0.0.1" {
			return true
		}

		// Check if the IP is in the private IP ranges
		privateIPBlocks := []string{
			"10.0.0.0/8",
			"172.16.0.0/12",
			"192.168.0.0/16",
		}

		for _, block := range privateIPBlocks {
			_, ipNet, err := net.ParseCIDR(block)
			if err != nil {
				log.Fatalf("Failed to parse CIDR block: %v", err)
			}
			if ipNet.Contains(net.ParseIP(ip)) {
				return true
			}
		}

		return false
}	

var authenticate = func () (*godo.Client, error) {
	token := os.Getenv("DIGITALOCEAN_ACCESS_TOKEN")
	if token == "" {
		return nil, errors.New("the environment variable DIGITALOCEAN_ACCESS_TOKEN is not set")
	}
	client := godo.NewFromToken(token)
	return client, nil
}
func doListAllDroplets(page int, pageSize int) ([]godo.Droplet, bool, error) {
	lastPage := false
	client, err := authenticate()
	options := &godo.ListOptions{
		Page: page,
		PerPage: pageSize,
		WithProjects: true,
	}
	if err != nil {
		log.Fatalf("Error authenticating: %v", err)
	}
	
	d, resp, err := client.Droplets.List(context.Background(), options)
	if err != nil {
		return nil, false, err
	}
	if resp.Links == nil || resp.Links.IsLastPage() {
		lastPage = true
	}
	return d, lastPage, err
}

func getDropletBasicInfo(page int, pageSize int) ([]DropletBasicInfo, bool, error) {
	var droplets = []DropletBasicInfo{}
	d, lastPage, err := doListAllDroplets(page, pageSize)
	if err != nil {
		return nil, false, err
	}
	for _, droplet := range d {
		privateAddr := ""
		publicAddr := ""
		if len(droplet.Networks.V4) > 0 {
			ip := droplet.Networks.V4[0].IPAddress
			if isPrivateIp(ip) {
				privateAddr = ip
			} else {
				publicAddr = ip
			}
		}

		if len(droplet.Networks.V4) >= 2 {
			ip := droplet.Networks.V4[1].IPAddress
			if isPrivateIp(ip) {
				privateAddr = ip
			} else {
				publicAddr = ip
			}
		}

		droplets = append(droplets, DropletBasicInfo{ID: droplet.ID, Name: droplet.Name, Region: droplet.Region.Slug, PrivateAddr: privateAddr, PublicAddr: publicAddr})
	}
	return droplets, lastPage, nil
}

// ListDroplets lists all droplets in the DigitalOcean account and prints them in a table format.
func ListDroplets() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleColoredDark)
	t.AppendHeader(table.Row{"id", "name", "region", "image_name", "image_id", "size", "price $/mo", "public_addr", "private_addr"})
	t.SortBy([]table.SortBy{
	    {Name: "region", Mode: table.Asc},
	    {Name: "name", Mode: table.Asc},
    })

	page:=1
	pageSize:=20
	
	for {
		d, lastPage, err := doListAllDroplets(page, pageSize)
		if err != nil {
			log.Fatalf("Failed to retrieve droplets from DigitalOcean API: %v", err)
		}
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
		}
	
		t.Render()
		page++
		if lastPage {
			break
		}
	}
}

type DropletBasicInfo struct {
	ID int
	Name string
	Region string
	PrivateAddr string
	PublicAddr string
}

func SshDroplet() {
	var droplets []DropletBasicInfo

	fmt.Println("SSH into a droplet")
	page:=1
	pageSize:=4
	d, lastPage, err := doListAllDroplets(page, pageSize)
	if err != nil {
		log.Fatalf("Failed to retrieve droplets from DigitalOcean API: %v", err)
	}
	fmt.Println(lastPage)


	for _, droplet := range d {
		privateAddr := ""
		publicAddr := ""
		if len(droplet.Networks.V4) > 0 {
			privateAddr = droplet.Networks.V4[0].IPAddress
		}
		if len(droplet.Networks.V4) >= 2 {
			publicAddr = droplet.Networks.V4[1].IPAddress
		}
		droplets = append(droplets, DropletBasicInfo{ID: droplet.ID, Name: droplet.Name, Region: droplet.Region.Slug, PrivateAddr: privateAddr, PublicAddr: publicAddr})
	}
	sort.Slice(droplets, func(i, j int) bool {
		return droplets[i].Name < droplets[j].Name
	})
	fmt.Println(droplets)
}