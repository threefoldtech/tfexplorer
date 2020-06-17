package workloads

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/threefoldtech/tfexplorer/mw"
	"github.com/threefoldtech/tfexplorer/pkg/capacity"
	"github.com/threefoldtech/tfexplorer/pkg/escrow"
	"github.com/threefoldtech/tfexplorer/pkg/workloads/types"
	"go.mongodb.org/mongo-driver/mongo"
)

// Setup injects and initializes directory package
func Setup(parent *mux.Router, db *mongo.Database, escrow escrow.Escrow, planner capacity.Planner) error {
	if err := types.Setup(context.TODO(), db); err != nil {
		return err
	}

	var api API
	api.escrow = escrow
	api.capacityPlanner = planner
	reservations := parent.PathPrefix("/reservations").Subrouter()

	reservations.HandleFunc("", mw.AsHandlerFunc(api.create)).Methods(http.MethodPost).Name("reservation-create")
	reservations.HandleFunc("", mw.AsHandlerFunc(api.list)).Methods(http.MethodGet).Name("reservation-list")
	reservations.HandleFunc("/{res_id:\\d+}", mw.AsHandlerFunc(api.get)).Methods(http.MethodGet).Name("reservation-get")
	reservations.HandleFunc("/{res_id:\\d+}/sign/provision", mw.AsHandlerFunc(api.signProvision)).Methods(http.MethodPost).Name("reservation-sign-provision")
	reservations.HandleFunc("/{res_id:\\d+}/sign/delete", mw.AsHandlerFunc(api.signDelete)).Methods(http.MethodPost).Name("reservation-sign-delete")

	reservations.HandleFunc("/workloads/{node_id}", mw.AsHandlerFunc(api.workloads)).Queries("from", "{from:\\d+}").Methods(http.MethodGet).Name("workloads-poll")
	reservations.HandleFunc("/workloads/{gwid:\\d+-\\d+}", mw.AsHandlerFunc(api.workloadGet)).Methods(http.MethodGet).Name("workload-get")
	reservations.HandleFunc("/workloads/{gwid:\\d+-\\d+}/{node_id}", mw.AsHandlerFunc(api.workloadPutResult)).Methods(http.MethodPut).Name("workloads-results")
	reservations.HandleFunc("/workloads/{gwid:\\d+-\\d+}/{node_id}", mw.AsHandlerFunc(api.workloadPutDeleted)).Methods(http.MethodDelete).Name("workloads-deleted")

	reservations.HandleFunc("/pools", mw.AsHandlerFunc(api.setupPool)).Methods(http.MethodPost).Name("pool-create")
	reservations.HandleFunc("/pools/{id:\\d+}", mw.AsHandlerFunc(api.getPool)).Methods(http.MethodGet).Name("pool-get")
	reservations.HandleFunc("/pools/owner/{owner:\\d+}", mw.AsHandlerFunc(api.listPools)).Methods(http.MethodGet).Name("pool-get")

	return nil
}
