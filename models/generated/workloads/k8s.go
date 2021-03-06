package workloads

import (
	"bytes"
	"fmt"
	"net"

	schema "github.com/threefoldtech/tfexplorer/schema"
)

var _ Workloader = (*K8S)(nil)
var _ Capaciter = (*K8S)(nil)

type K8SCustomSize struct {
	CRU int64   `bson:"cru" json:"cru" `
	MRU float64 `bson:"mru" json:"mru" `
	SRU float64 `bson:"sru" json:"sru" `
}

type K8S struct {
	ReservationInfo `bson:",inline"`

	Size                  int64             `bson:"size" json:"size"`
	ClusterSecret         string            `bson:"cluster_secret" json:"cluster_secret"`
	NetworkId             string            `bson:"network_id" json:"network_id"`
	Ipaddress             net.IP            `bson:"ipaddress" json:"ipaddress"`
	MasterIps             []net.IP          `bson:"master_ips" json:"master_ips"`
	SshKeys               []string          `bson:"ssh_keys" json:"ssh_keys"`
	StatsAggregator       []StatsAggregator `bson:"stats_aggregator" json:"stats_aggregator"`
	PublicIP              schema.ID         `bson:"public_ip" json:"public_ip"`
	DatastoreEndpoint     string            `bson:"datastore_endpoint" json:"datastore_endpoint"`
	DisableDefaultIngress bool              `bson:"disable_default_ingress" json:"disable_default_ingress"`
	CustomSize            K8SCustomSize     `bson:"custom_size" json:"custom_size"`
}

var k8sSize = map[int64]RSU{
	1: {
		CRU: 1,
		MRU: 2,
		SRU: 50,
	},
	2: {
		CRU: 2,
		MRU: 4,
		SRU: 100,
	},
	3: {
		CRU: 2,
		MRU: 8,
		SRU: 25,
	},
	4: {
		CRU: 2,
		MRU: 5,
		SRU: 50,
	},
	5: {
		CRU: 2,
		MRU: 8,
		SRU: 200,
	},
	6: {
		CRU: 4,
		MRU: 16,
		SRU: 50,
	},
	7: {
		CRU: 4,
		MRU: 16,
		SRU: 100,
	},
	8: {
		CRU: 4,
		MRU: 16,
		SRU: 400,
	},
	9: {
		CRU: 8,
		MRU: 32,
		SRU: 100,
	},
	10: {
		CRU: 8,
		MRU: 32,
		SRU: 200,
	},
	11: {
		CRU: 8,
		MRU: 32,
		SRU: 800,
	},
	12: {
		CRU: 1,
		MRU: 64,
		SRU: 200,
	},
	13: {
		CRU: 1,
		MRU: 64,
		SRU: 400,
	},
	14: {
		CRU: 1,
		MRU: 64,
		SRU: 800,
	},
	15: {
		CRU: 1,
		MRU: 2,
		SRU: 25,
	},
	16: {
		CRU: 2,
		MRU: 4,
		SRU: 50,
	},
	17: {
		CRU: 4,
		MRU: 8,
		SRU: 50,
	},
	18: {
		CRU: 1,
		MRU: 1,
		SRU: 25,
	},
}

func (k *K8S) GetRSU() (RSU, error) {
	if k.Size == -1 {
		return RSU{
			CRU: k.CustomSize.CRU,
			MRU: k.CustomSize.MRU,
			SRU: k.CustomSize.SRU,
		}, nil
	}

	rsu, ok := k8sSize[k.Size]
	if !ok {
		return RSU{}, fmt.Errorf("K8S VM size %d is not supported", k.Size)
	}
	return rsu, nil
}

func (k *K8S) SignatureChallenge() ([]byte, error) {
	ric, err := k.ReservationInfo.SignatureChallenge()
	if err != nil {
		return nil, err
	}

	b := bytes.NewBuffer(ric)
	if _, err := fmt.Fprintf(b, "%d", k.Size); err != nil {
		return nil, err
	}
	if _, err := fmt.Fprintf(b, "%s", k.ClusterSecret); err != nil {
		return nil, err
	}
	if _, err := fmt.Fprintf(b, "%s", k.NetworkId); err != nil {
		return nil, err
	}
	if _, err := fmt.Fprintf(b, "%s", k.Ipaddress.String()); err != nil {
		return nil, err
	}
	for _, ip := range k.MasterIps {
		if _, err := fmt.Fprintf(b, "%s", ip.String()); err != nil {
			return nil, err
		}
	}
	for _, key := range k.SshKeys {
		if _, err := fmt.Fprintf(b, "%s", key); err != nil {
			return nil, err
		}
	}
	if _, err := fmt.Fprintf(b, "%d", k.PublicIP); err != nil {
		return nil, err
	}
	if _, err := fmt.Fprintf(b, "%s", k.DatastoreEndpoint); err != nil {
		return nil, err
	}
	if _, err := fmt.Fprintf(b, "%t", k.DisableDefaultIngress); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
