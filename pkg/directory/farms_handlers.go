package directory

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/zaibon/httpsig"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/threefoldtech/tfexplorer/models"
	"github.com/threefoldtech/tfexplorer/mw"
	directory "github.com/threefoldtech/tfexplorer/pkg/directory/types"
	"github.com/threefoldtech/tfexplorer/schema"

	"github.com/gorilla/mux"
)

type farmKey struct{}

func (f FarmAPI) isAuthenticated(r *http.Request) bool {
	_, err := f.verifier.Verify(r)
	return err == nil
}

func (f *FarmAPI) registerFarm(r *http.Request) (interface{}, mw.Response) {
	log.Info().Msg("farm register request received")

	db := mw.Database(r)
	defer r.Body.Close()

	var info directory.Farm
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		return nil, mw.BadRequest(err)
	}

	if err := info.Validate(); err != nil {
		return nil, mw.BadRequest(err)
	}

	id, err := f.Add(r.Context(), db, info)
	if err != nil {
		return nil, mw.Error(err)
	}

	return struct {
		ID schema.ID `json:"id"`
	}{
		id,
	}, mw.Created()
}

func (f *FarmAPI) updateFarm(r *http.Request) (interface{}, mw.Response) {
	farmID := getFarmID(r.Context())

	var info directory.Farm
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		return nil, mw.BadRequest(err)
	}

	info.ID = farmID

	db := mw.Database(r)
	err := f.Update(r.Context(), db, farmID, info)
	if err != nil {
		return nil, mw.Error(err)
	}

	return nil, mw.Ok()
}

func (f *FarmAPI) listFarm(r *http.Request) (interface{}, mw.Response) {

	q := directory.FarmQuery{}
	if err := q.Parse(r); err != nil {
		return nil, err
	}
	var filter directory.FarmFilter
	filter = filter.WithFarmQuery(q)

	db := mw.Database(r)

	var findOpts []*options.FindOptions

	pager := models.PageFromRequest(r)
	findOpts = append(findOpts, pager)
	// hide the email of the farm for any non authenticated user
	if !f.isAuthenticated(r) {
		findOpts = append(findOpts, options.Find().SetProjection(bson.D{
			{Key: "email", Value: 0},
		}))
	}

	farms, total, err := f.List(r.Context(), db, filter, findOpts...)
	if err != nil {
		return nil, mw.Error(err)
	}

	pages := fmt.Sprintf("%d", models.Pages(pager, total))
	return farms, mw.Ok().WithHeader("Pages", pages)
}

func (f *FarmAPI) getFarm(r *http.Request) (interface{}, mw.Response) {
	sid := mux.Vars(r)["farm_id"]

	id, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return nil, mw.BadRequest(err)
	}

	db := mw.Database(r)

	farm, err := f.GetByID(r.Context(), db, id)
	if err != nil {
		return nil, mw.NotFound(err)
	}

	// hide the email of the farm for any non authenticated user
	if !f.isAuthenticated(r) {
		farm.Email = ""
	}

	return farm, nil
}

func (f *FarmAPI) deleteNodeFromFarm(r *http.Request) (interface{}, mw.Response) {
	farmID := getFarmID(r.Context())

	var nodeAPI NodeAPI
	nodeID := mux.Vars(r)["node_id"]

	db := mw.Database(r)
	node, err := nodeAPI.Get(r.Context(), db, nodeID, false)
	if err != nil {
		return nil, mw.NotFound(err)
	}

	if node.FarmId != int64(farmID) {
		return nil, mw.Forbidden(fmt.Errorf("only the farm owner of this node can delete this node"))
	}

	err = nodeAPI.Delete(r.Context(), db, nodeID)
	if err != nil {
		return nil, mw.Error(err)
	}

	return nil, mw.Ok()
}

func (f *FarmAPI) addFarmIPs(r *http.Request) (interface{}, mw.Response) {
	// Get the farm from the middleware context
	farmID := getFarmID(r.Context())

	var info []struct {
		IP schema.IPCidr `json:"address"`
		GW net.IP        `json:"gateway"`
	}

	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		return nil, mw.BadRequest(err)
	}

	db := mw.Database(r)
	for _, entry := range info {
		ip, err := entry.IP.Validate()
		if err != nil {
			return nil, mw.BadRequest(err)
		}

		err = f.PushIP(r.Context(), db, farmID, ip, entry.GW)
		if err != nil {
			return nil, mw.BadRequest(err)
		}
	}
	return nil, mw.Ok()
}

func (f *FarmAPI) deleteFarmIps(r *http.Request) (interface{}, mw.Response) {
	// Get the farm from the middleware context
	farmID := getFarmID(r.Context())

	var info []schema.IPCidr
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		return nil, mw.BadRequest(err)
	}

	db := mw.Database(r)
	for _, ip := range info {
		ip, err := ip.Validate()
		if err != nil {
			return nil, mw.BadRequest(err)
		}

		err = f.RemoveIP(r.Context(), db, farmID, ip)
		if err != nil {
			return nil, mw.BadRequest(err)
		}
	}
	return nil, mw.Ok()
}

// LoadFarmMiddleware verifies if the farmer who wants to update information
// about it's farm is indeed that farmer
func LoadFarmMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sid := mux.Vars(r)["farm_id"]

		id, err := strconv.ParseInt(sid, 10, 64)
		if err != nil {
			mw.BadRequest(err)
			return
		}

		db := mw.Database(r)

		farmerID := httpsig.KeyIDFromContext(r.Context())
		requestFarmerID, err := strconv.ParseInt(farmerID, 10, 64)
		if err != nil {
			mw.BadRequest(err)
			return
		}

		var filter directory.FarmFilter
		filter = filter.WithID(schema.ID(id)).WithOwner(requestFarmerID)
		farm, err := filter.Get(r.Context(), db)
		if err != nil {
			mw.NotFound(err)
			return
		}

		// Store the farm object in the request context for later usage
		ctx := context.WithValue(r.Context(), farmKey{}, farm.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getFarmID(ctx context.Context) schema.ID {
	return ctx.Value(farmKey{}).(schema.ID)
}

func (f *FarmAPI) getFarmPricing(r *http.Request) (interface{}, mw.Response) {
	sid := mux.Vars(r)["farm_id"]

	id, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return nil, mw.BadRequest(err)
	}

	db := mw.Database(r)

	farm, err := f.GetByID(r.Context(), db, id)
	if err != nil {
		return nil, mw.NotFound(err)
	}

	// hide the email of the farm for any non authenticated user
	if !f.isAuthenticated(r) {
		farm.Email = ""
	}

	return farm, nil
}
