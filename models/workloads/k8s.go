package workloads

import (
	"bytes"
	"fmt"
	"net"
)

var _ Workloader = (*K8S)(nil)
var _ Capaciter = (*K8S)(nil)

type K8S struct {
	contract Contract
	state    State

	Size            int64             `bson:"size" json:"size"`
	ClusterSecret   string            `bson:"cluster_secret" json:"cluster_secret"`
	NetworkId       string            `bson:"network_id" json:"network_id"`
	Ipaddress       net.IP            `bson:"ipaddress" json:"ipaddress"`
	MasterIps       []net.IP          `bson:"master_ips" json:"master_ips"`
	SshKeys         []string          `bson:"ssh_keys" json:"ssh_keys"`
	StatsAggregator []StatsAggregator `bson:"stats_aggregator" json:"stats_aggregator"`
}

// Contract implements the Workloader interface
func (k *K8S) Contract() Contract { return k.contract }

// State implements the Workloader interface
func (k *K8S) State() *State { return &k.state }

// GetRSU implements the Capaciter interface
func (k *K8S) GetRSU() RSU {
	switch k.Size {
	case 1:
		return RSU{
			CRU: 1,
			MRU: 2,
			SRU: 50,
		}
	case 2:
		return RSU{
			CRU: 2,
			MRU: 4,
			SRU: 100,
		}
	}
	return RSU{}
}

func (k *K8S) SignatureChallenge() ([]byte, error) {
	ric, err := k.contract.SignatureChallenge()
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

	return b.Bytes(), nil
}