package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/pkg/errors"
	"github.com/threefoldtech/tfexplorer/client"
	"github.com/threefoldtech/tfexplorer/models/generated/directory"
	"github.com/threefoldtech/tfexplorer/schema"
	"github.com/urfave/cli"
)

func registerFarm(c *cli.Context) (err error) {
	name := c.Args().First()
	if name == "" {
		return fmt.Errorf("farm name needs to be specified")
	}

	addrs := c.StringSlice("addresses")
	email := c.String("email")
	iyo := c.String("iyo_organization")

	addresses := make([]directory.WalletAddress, len(addrs))
	for i := range addrs {
		addresses[i].Asset, addresses[i].Address, err = splitAddressCode(addrs[i])
		if err != nil {
			return err
		}
	}

	farm := directory.Farm{
		Name:            name,
		ThreebotID:      int64(userid.ThreebotID),
		Email:           schema.Email(email),
		IyoOrganization: iyo,
		WalletAddresses: addresses,
	}

	farm.ID, err = db.FarmRegister(farm)
	if err != nil {
		return err
	}

	fmt.Println("Farm registered successfully")
	fmt.Println(formatFarm(farm))
	return nil
}

func updateFarm(c *cli.Context) error {
	id := c.Int64("id")
	farm, err := db.FarmGet(schema.ID(id))
	if err != nil {
		return err
	}

	addrs := c.StringSlice("addresses")
	email := c.String("email")
	iyo := c.String("iyo_organization")

	if len(addrs) > 0 {
		addresses := make([]directory.WalletAddress, len(addrs))
		for i := range addrs {
			addresses[i].Asset, addresses[i].Address, err = splitAddressCode(addrs[i])
			if err != nil {
				return err
			}
		}
		farm.WalletAddresses = addresses
	}

	if email != "" {
		farm.Email = schema.Email(email)
	}

	if iyo != "" {
		farm.IyoOrganization = iyo
	}

	if err := db.FarmUpdate(farm); err != nil {
		return errors.Wrap(err, "Failed to update farm")
	}

	fmt.Println("Farm updated successfully")
	fmt.Println(formatFarm(farm))
	return nil
}

func addIP(c *cli.Context) error {
	id := c.Int64("id")
	farm, err := db.FarmGet(schema.ID(id))
	if err != nil {
		return err
	}

	addrs := c.String("address")
	gateway := c.String("gateway")
	if addrs == "" || gateway == "" {
		return fmt.Errorf("address and gateway need to be specified")
	}

	ip, err := schema.ParseIPCidr(addrs)
	if err != nil {
		return err
	}

	gatewayIP := net.ParseIP(gateway)
	if gatewayIP == nil {
		return fmt.Errorf("gateway ip not valid")
	}

	address := directory.PublicIP{
		Address: ip,
		Gateway: schema.IP{gatewayIP},
	}

	if err := db.FarmAddIP(farm.ID, address); err != nil {
		return err
	}
	fmt.Println("IP address added successfully")
	return nil
}

func deleteIP(c *cli.Context) error {
	id := c.Int64("id")
	farm, err := db.FarmGet(schema.ID(id))
	if err != nil {
		return err
	}

	addrs := c.String("address")
	if addrs == "" {
		return fmt.Errorf("IP address needs to be specified")
	}
	ip, err := schema.ParseIPCidr(addrs)
	if err != nil {
		return err
	}

	for i := range farm.IPAddresses {
		if farm.IPAddresses[i].Address.IP.Equal(ip.IP) {
			if err := db.FarmDeleteIP(farm.ID, ip); err != nil {
				return err
			}
			fmt.Println("IP address deleted successfully")
			return nil
		}
	}
	return fmt.Errorf("IP address not found in farm")

}

func listFarms(c *cli.Context) (err error) {
	farmsRet := make([]directory.Farm, 0)
	pageNumber := 1

	for {
		pager := client.Page(pageNumber, 20)
		farms, err := db.FarmList(schema.ID(int64(userid.ThreebotID)), "", pager)
		farmsRet = append(farmsRet, farms...)
		if err != nil {
			break
		}
		if len(farms) == 0 {
			break
		}
		pageNumber++
	}
	if len(farmsRet) > 0 {
		for _, f := range farmsRet {
			fmt.Println(formatFarm(f))
			fmt.Println("-------------------------------------------")
		}
		// fmt.Println(farmsRet)
	} else {
		fmt.Println("No farms found")
	}

	return nil
}

func splitAddressCode(addr string) (string, string, error) {
	ss := strings.Split(addr, ":")
	if len(ss) != 2 {
		return "", "", fmt.Errorf("wrong format for wallet address %s, should be 'asset:address'", addr)
	}

	return ss[0], ss[1], nil
}
func formatFarm(farm directory.Farm) string {
	b := &strings.Builder{}
	fmt.Fprintf(b, "ID: %d\n", farm.ID)
	fmt.Fprintf(b, "Name: %s\n", farm.Name)
	fmt.Fprintf(b, "Email: %s\n", farm.Email)
	fmt.Fprintf(b, "Farmer TheebotID: %d\n", farm.ThreebotID)
	fmt.Fprintf(b, "Wallet addresses:\n")
	for _, a := range farm.WalletAddresses {
		fmt.Fprintf(b, "  %s:%s\n", a.Asset, a.Address)
	}
	fmt.Fprintf(b, "IP addresses:\n")
	for _, a := range farm.IPAddresses {
		fmt.Fprintf(b, "  IP-> %s   Gateway-> %s\n", a.Address, a.Gateway)
	}
	return b.String()
}
