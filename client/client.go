package client

import (
	"crypto/ed25519"
	"fmt"
	"net/url"

	"github.com/threefoldtech/tfexplorer/models/generated/directory"
	"github.com/threefoldtech/tfexplorer/models/generated/phonebook"
	"github.com/threefoldtech/tfexplorer/models/generated/workloads"
	"github.com/threefoldtech/tfexplorer/pkg/capacity/types"
	wrklds "github.com/threefoldtech/tfexplorer/pkg/workloads"
	"github.com/threefoldtech/tfexplorer/schema"
	"github.com/threefoldtech/zos/pkg/capacity"
	"github.com/threefoldtech/zos/pkg/capacity/dmi"
)

type (
	// Client structure
	Client struct {
		Phonebook Phonebook
		Directory Directory
		Workloads Workloads
	}

	// NodeIter iterator over all nodes
	//
	// If the iterator has finished, a nil error and nil node pointer is returned
	NodeIter interface {
		Next() (*directory.Node, error)
	}

	// FarmIter iterator over all farms
	//
	// If the iterator has finished, a nil error and nil farm pointer is returned
	FarmIter interface {
		Next() (*directory.Farm, error)
	}

	// Directory API interface
	Directory interface {
		FarmRegister(farm directory.Farm) (schema.ID, error)
		FarmUpdate(farm directory.Farm) error
		FarmList(tid schema.ID, name string, page *Pager) (farms []directory.Farm, err error)
		FarmGet(id schema.ID) (farm directory.Farm, err error)
		Farms(cacheSize int) FarmIter
		FarmAddIP(id schema.ID, ip directory.PublicIP) error
		FarmDeleteIP(id schema.ID, ip schema.IPCidr) error

		GatewayRegister(Gateway directory.Gateway) error
		GatewayList(tid schema.ID, name string, page *Pager) (farms []directory.Gateway, err error)
		GatewayGet(id string) (farm directory.Gateway, err error)
		GatewayUpdateUptime(id string, uptime uint64) error
		GatewayUpdateReservedResources(id string, resources directory.ResourceAmount, workloads directory.WorkloadAmount) error

		NodeRegister(node directory.Node) error
		NodeList(filter NodeFilter, pager *Pager) (nodes []directory.Node, err error)
		NodeGet(id string, proofs bool) (node directory.Node, err error)
		Nodes(cacheSize int, proofs bool) NodeIter

		NodeSetInterfaces(id string, ifaces []directory.Iface) error
		NodeSetPorts(id string, ports []uint) error
		NodeSetPublic(id string, pub directory.PublicIface) error
		NodeSetFreeToUse(id string, free bool) error

		//TODO: this method call uses types from zos that is not generated
		//from the schema. Which is wrong imho.
		NodeSetCapacity(
			id string,
			resources directory.ResourceAmount,
			dmiInfo dmi.DMI,
			disksInfo capacity.Disks,
			hypervisor []string,
		) error

		NodeUpdateUptime(id string, uptime uint64) error
		NodeUpdateUsedResources(id string, resources directory.ResourceAmount, workloads directory.WorkloadAmount) error
	}

	// Phonebook interface
	Phonebook interface {
		Create(user phonebook.User) (schema.ID, error)
		List(name, email string, page *Pager) (output []phonebook.User, err error)
		Get(id schema.ID) (phonebook.User, error)
		// Update() #TODO
		Validate(id schema.ID, message, signature string) (bool, error)
	}

	// Workloads interface
	Workloads interface {
		Create(reservation workloads.Workloader) (resp wrklds.ReservationCreateResponse, err error)
		List(nextAction *workloads.NextActionEnum, customerTid int64, page *Pager) (reservation []workloads.Reservation, err error)
		Get(id schema.ID) (reservation workloads.Workloader, err error)

		SignProvision(id schema.ID, user schema.ID, signature string) error
		SignDelete(id schema.ID, user schema.ID, signature string) error

		PoolCreate(reservation types.Reservation) (resp wrklds.CapacityPoolCreateResponse, err error)
		PoolGet(poolID string) (result types.Pool, err error)
		PoolsGetByOwner(ownerID string) (result []types.Pool, err error)

		NodeWorkloads(nodeID string, from uint64) ([]workloads.Workloader, uint64, error)
		NodeWorkloadGet(gwid string) (result workloads.Workloader, err error)
		NodeWorkloadPutResult(nodeID, gwid string, result workloads.Result) error
		NodeWorkloadPutDeleted(nodeID, gwid string) error
	}

	// Identity is used by the client to authenticate to the explorer API
	Identity interface {
		// The unique ID as known by the explorer
		Identity() string
		// PrivateKey used to sign the requests
		PrivateKey() ed25519.PrivateKey
	}

	// Pager for listing
	Pager struct {
		p int
		s int
	}
)

func (p *Pager) apply(v url.Values) {
	if p == nil {
		return
	}

	if p.p < 1 {
		p.p = 1
	}

	if p.s == 0 {
		p.s = 10
	}

	v.Set("page", fmt.Sprint(p.p))
	v.Set("size", fmt.Sprint(p.s))
}

// Page returns a pager
func Page(page, size int) *Pager {
	return &Pager{p: page, s: size}
}

// NewClient creates a new client, if identity is not nil, it will be used
// to authenticate requests against the server
func NewClient(u string, id Identity) (*Client, error) {
	h, err := newHTTPClient(u, id)
	if err != nil {
		return nil, err
	}
	cl := &Client{
		Phonebook: &httpPhonebook{h},
		Directory: &httpDirectory{h},
		Workloads: &httpWorkloads{h},
	}

	return cl, nil
}
