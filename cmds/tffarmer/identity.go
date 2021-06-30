package main

import (
	"encoding/hex"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/mitchellh/go-homedir"
	"github.com/threefoldtech/tfexplorer"
	"github.com/threefoldtech/tfexplorer/client"
	"github.com/threefoldtech/tfexplorer/models/generated/phonebook"
	"github.com/threefoldtech/zos/pkg/identity"
	"github.com/urfave/cli"
)

func getSeedPath() (location string, err error) {
	// Get home directory for current user
	dir, err := homedir.Dir()
	if err != nil {
		return "", errors.Wrap(err, "Cannot get current user home directory")
	}
	if dir == "" {
		return "", errors.Wrap(err, "Cannot get current user home directory")
	}
	expandedDir, err := homedir.Expand(dir)
	if err != nil {
		return "", err
	}

	path := filepath.Join(expandedDir, ".config", "tffarmer.seed")
	return path, nil

}

func generateNewUser(c *cli.Context, url string, seedPath string) (userIdentity *tfexplorer.UserIdentity, err error) {
	var name, email string
	fmt.Print("Enter a name for the identity: ")
	fmt.Scanln(&name)
	fmt.Print("Enter an email: ")
	fmt.Scanln(&email)

	_, userIdentity, err = generateID(c, url, name, email, seedPath)
	if err != nil {
		return userIdentity, err
	}
	return userIdentity, err

}

func generateID(c *cli.Context, url, name, email, seedPath string) (user phonebook.User, ui *tfexplorer.UserIdentity, err error) {
	ui = &tfexplorer.UserIdentity{}

	k, err := identity.GenerateKeyPair()
	if err != nil {
		return phonebook.User{}, ui, err
	}

	ui = tfexplorer.NewUserIdentity(k, 0)

	user = phonebook.User{
		Name:        name,
		Email:       email,
		Pubkey:      hex.EncodeToString(ui.Key().PublicKey),
		Description: "",
	}

	log.Debug().Msg("initializing client with created key")
	httpClient, err := client.NewClient(url, ui)
	if err != nil {
		return user, ui, err
	}

	log.Debug().Msg("registering user")
	id, err := httpClient.Phonebook.Create(user)
	if err != nil {
		return user, ui, errors.Wrap(err, "failed to register user")
	}

	// Update UserData with created id
	ui.ThreebotID = uint64(id)

	// Saving new seed struct

	if err := ui.Save(seedPath); err != nil {
		return user, ui, errors.Wrap(err, "failed to save seed")
	}

	fmt.Println("Your ID is: ", id)
	fmt.Println("Seed saved in: ", seedPath, " Please make sure you have it backed up.")
	return user, ui, nil
}
