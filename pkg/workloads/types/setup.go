package types

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Setup sets up indexes for types, must be called at least
// Onetime during the life time of the object
func Setup(ctx context.Context, db *mongo.Database) error {
	col := db.Collection(ReservationCollection)
	indexes := []mongo.IndexModel{
		{
			Keys: bson.M{"data_reservation.networks.network_resources.node_id": 1},
		},
	}

	types := []string{"containers", "volumes", "zdbs", "networks", "kubernetes",
		"proxies", "reverse_proxies", "subdomains", "domain_delegates", "gateway4to6", "capacity_pool"}
	for _, typ := range types {
		indexes = append(
			indexes,
			mongo.IndexModel{
				Keys: bson.M{fmt.Sprintf("data_reservation.%s.node_id", typ): 1},
			},
		)

	}
	indexes = append(indexes, mongo.IndexModel{Keys: bson.M{"next_action": 1}})
	indexes = append(indexes, mongo.IndexModel{Keys: bson.M{"customer_tid": 1}})

	if _, err := col.Indexes().CreateMany(ctx, indexes); err != nil {
		return err
	}

	col = db.Collection(queueCollection)
	indexes = []mongo.IndexModel{
		{
			Keys: bson.M{"node_id": 1},
		},
		{
			Keys: bson.M{"workload_id": 1},
		},
	}

	if _, err := col.Indexes().CreateMany(ctx, indexes); err != nil {
		return err
	}

	col = db.Collection(WorkloadCollection)
	indexes = []mongo.IndexModel{
		{
			Keys: bson.M{"node_id": 1},
		},
		{
			Keys: bson.M{"workload_id": 1},
		},
		{
			Keys: bson.M{"pool_id": 1},
		},
		{
			Keys: bson.M{"customer_tid": 1},
		},
		{
			Keys: bson.M{"next_action": 1},
		},
		{
			Keys: bson.M{"workload_type": 1},
		},
		{
			Keys: bson.M{"public_ip": 1},
		},
	}

	if _, err := col.Indexes().CreateMany(ctx, indexes); err != nil {
		return err
	}

	return nil
}
