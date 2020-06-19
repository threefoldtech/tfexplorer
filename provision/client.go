package provision

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/threefoldtech/tfexplorer"
	"github.com/threefoldtech/tfexplorer/client"
	"github.com/threefoldtech/tfexplorer/models/generated/workloads"
	"github.com/threefoldtech/tfexplorer/pkg/capacity"
	wrklds "github.com/threefoldtech/tfexplorer/pkg/workloads"
	"github.com/threefoldtech/tfexplorer/schema"
)

// ReservationClient is a client to deploy and delete reservations
type ReservationClient struct {
	explorer *client.Client
	userID   *tfexplorer.UserIdentity
}

// NewReservationClient creates a new reservation client
func NewReservationClient(explorer *client.Client, userID *tfexplorer.UserIdentity) *ReservationClient {
	return &ReservationClient{
		explorer: explorer,
		userID:   userID,
	}
}

// Deploy deploys the reservation
func (r *ReservationClient) Deploy(reservation workloads.Reservation, currencies []string) (wrklds.ReservationCreateResponse, error) {
	reservationToCreate, err := r.DryRun(reservation, currencies)
	if err != nil {
		return wrklds.ReservationCreateResponse{}, nil
	}

	fmt.Printf("%+v", reservationToCreate)

	response, err := r.explorer.Workloads.Create(reservationToCreate)
	if err != nil {
		return wrklds.ReservationCreateResponse{}, errors.Wrap(err, "failed to send reservation")
	}

	return response, nil
}

// DryRun will return the reservation to deploy and marshals the data of the reservation
func (r *ReservationClient) DryRun(reservation workloads.Reservation, currencies []string) (workloads.Reservation, error) {
	userID := int64(r.userID.ThreebotID)
	signer, err := client.NewSigner(r.userID.Key().PrivateKey.Seed())
	if err != nil {
		return workloads.Reservation{}, errors.Wrap(err, "could not load signer")
	}

	reservation.CustomerTid = userID

	// set allowed the currencies as provided by the user
	reservation.DataReservation.Currencies = currencies

	bytes, err := json.Marshal(reservation.DataReservation)
	if err != nil {
		return workloads.Reservation{}, err
	}

	reservation.Json = string(bytes)
	_, signature, err := signer.SignHex(reservation.Json)
	if err != nil {
		return workloads.Reservation{}, errors.Wrap(err, "failed to sign the reservation")
	}

	reservation.CustomerSignature = signature

	return reservation, nil
}

// DeployCapacityPool deploys the reservation
func (r *ReservationClient) DeployCapacityPool(reservation capacity.Reservation, currencies []string) (wrklds.CapacityPoolCreateResponse, error) {
	reservationToCreate, err := r.DryRunCapacity(reservation, currencies)
	if err != nil {
		return wrklds.CapacityPoolCreateResponse{}, nil
	}

	fmt.Printf("%+v", reservationToCreate)

	response, err := r.explorer.Workloads.PoolCreate(reservationToCreate)
	if err != nil {
		return wrklds.CapacityPoolCreateResponse{}, errors.Wrap(err, "failed to send reservation")
	}

	return response, nil
}

// DryRunCapacity will return the reservation to deploy and marshals the data of the reservation
func (r *ReservationClient) DryRunCapacity(reservation capacity.Reservation, currencies []string) (capacity.Reservation, error) {
	userID := int64(r.userID.ThreebotID)
	signer, err := client.NewSigner(r.userID.Key().PrivateKey.Seed())
	if err != nil {
		return capacity.Reservation{}, errors.Wrap(err, "could not load signer")
	}

	reservation.CustomerTid = userID

	// set allowed the currencies as provided by the user
	reservation.DataReservation.Currencies = currencies

	bytes, err := json.Marshal(reservation.DataReservation)
	if err != nil {
		return capacity.Reservation{}, err
	}

	reservation.JSON = string(bytes)
	_, signature, err := signer.SignHex(reservation.JSON)
	if err != nil {
		return capacity.Reservation{}, errors.Wrap(err, "failed to sign the reservation")
	}

	reservation.CustomerSignature = signature

	return reservation, nil
}

// DeleteReservation deletes a reservation by id
func (r *ReservationClient) DeleteReservation(resID int64) error {
	userID := int64(r.userID.ThreebotID)

	reservation, err := r.explorer.Workloads.Get(schema.ID(resID))
	if err != nil {
		return errors.Wrap(err, "failed to get reservation info")
	}

	signer, err := client.NewSigner(r.userID.Key().PrivateKey.Seed())
	if err != nil {
		return errors.Wrapf(err, "failed to load signer")
	}

	_, signature, err := signer.SignHex(resID, reservation.Json)
	if err != nil {
		return errors.Wrap(err, "failed to sign the reservation")
	}

	if err := r.explorer.Workloads.SignDelete(schema.ID(resID), schema.ID(userID), signature); err != nil {
		return errors.Wrapf(err, "failed to sign deletion of reservation: %d", resID)
	}

	fmt.Printf("Reservation %v marked as to be deleted\n", resID)
	return nil
}
