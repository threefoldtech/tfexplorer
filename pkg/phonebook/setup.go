package phonebook

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/threefoldtech/tfexplorer/mw"
	phonebook "github.com/threefoldtech/tfexplorer/pkg/phonebook/types"
	"github.com/zaibon/httpsig"
	"go.mongodb.org/mongo-driver/mongo"
)

// Setup injects and initializes directory package
func Setup(parent *mux.Router, db *mongo.Database, threebotConnectURL string) error {
	if err := phonebook.Setup(context.TODO(), db); err != nil {
		return err
	}

	userVerifier := httpsig.NewVerifier(mw.NewUserKeyGetter(db))

	var userAPI = UserAPI{
		verifier:              userVerifier,
		threebotConnectAPIURL: threebotConnectURL,
	}

	// versionned endpoints
	api := parent.PathPrefix("/api/v1").Subrouter()
	users := api.PathPrefix("/users").Subrouter()

	users.HandleFunc("", mw.AsHandlerFunc(userAPI.create)).Methods(http.MethodPost).Name("user-create-v1")
	users.HandleFunc("", mw.AsHandlerFunc(userAPI.list)).Methods(http.MethodGet).Name(("user-list-v1"))
	users.HandleFunc("/{user_id}", mw.AsHandlerFunc(userAPI.register)).Methods(http.MethodPut).Name("user-register-v1")
	users.HandleFunc("/{user_id}", mw.AsHandlerFunc(userAPI.get)).Methods(http.MethodGet).Name("user-get-v1")
	users.HandleFunc("/{user_id}/validate", mw.AsHandlerFunc(userAPI.validate)).Methods(http.MethodPost).Name("user-validate-v1")

	// legacy endpoints
	legacyUsers := parent.PathPrefix("/explorer/users").Subrouter()

	legacyUsers.HandleFunc("", mw.AsHandlerFunc(userAPI.create)).Methods(http.MethodPost).Name("user-create")
	legacyUsers.HandleFunc("", mw.AsHandlerFunc(userAPI.list)).Methods(http.MethodGet).Name(("user-list"))
	legacyUsers.HandleFunc("/{user_id}", mw.AsHandlerFunc(userAPI.register)).Methods(http.MethodPut).Name("user-register")
	legacyUsers.HandleFunc("/{user_id}", mw.AsHandlerFunc(userAPI.get)).Methods(http.MethodGet).Name("user-get")
	legacyUsers.HandleFunc("/{user_id}/validate", mw.AsHandlerFunc(userAPI.validate)).Methods(http.MethodPost).Name("user-validate")

	return nil
}
