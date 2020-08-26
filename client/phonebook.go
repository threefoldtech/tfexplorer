package client

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/threefoldtech/tfexplorer/models/phonebook"
	"github.com/threefoldtech/tfexplorer/schema"
)

type httpPhonebook struct {
	*httpClient
}

func (p *httpPhonebook) Create(user phonebook.User) (schema.ID, error) {
	var out phonebook.User
	if _, err := p.post(p.url("users"), user, &out); err != nil {
		return 0, err
	}

	return out.ID, nil
}

func (p *httpPhonebook) List(name, email string, page *Pager) (output []phonebook.User, err error) {
	query := url.Values{}
	page.apply(query)
	if len(name) != 0 {
		query.Set("name", name)
	}
	if len(email) != 0 {
		query.Set("email", email)
	}

	_, err = p.get(p.url("users"), query, &output, http.StatusOK)

	return
}

func (p *httpPhonebook) Get(id schema.ID) (user phonebook.User, err error) {
	_, err = p.get(p.url("users", fmt.Sprint(id)), nil, &user, http.StatusOK)
	return
}

// Validate the signature of this message for the user, signature and message are hex encoded
func (p *httpPhonebook) Validate(id schema.ID, message, signature string) (bool, error) {
	var input struct {
		S string `json:"signature"`
		M string `json:"payload"`
	}
	input.S = signature
	input.M = message

	var output struct {
		V bool `json:"is_valid"`
	}

	_, err := p.post(p.url("users", fmt.Sprint(id), "validate"), input, &output, http.StatusOK)
	if err != nil {
		return false, err
	}

	return output.V, nil
}
