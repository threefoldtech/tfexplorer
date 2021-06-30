package client

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/threefoldtech/tfexplorer/models/generated/directory"
	"github.com/threefoldtech/tfexplorer/schema"
	"github.com/threefoldtech/zos/pkg/capacity"
	"github.com/threefoldtech/zos/pkg/capacity/dmi"
)

type (
	httpDirectory struct {
		*httpClient
	}

	httpNodeIter struct {
		cl       *httpDirectory
		proofs   bool
		page     int
		size     int
		cache    []directory.Node
		cacheIdx int
		finished bool
	}

	httpFarmIter struct {
		cl       *httpDirectory
		page     int
		size     int
		cache    []directory.Farm
		cacheIdx int
		finished bool
	}
)

func (d *httpDirectory) FarmRegister(farm directory.Farm) (schema.ID, error) {
	var output struct {
		ID schema.ID `json:"id"`
	}

	_, err := d.post(d.url("farms"), farm, &output, http.StatusCreated)
	return output.ID, err
}

func (d *httpDirectory) FarmUpdate(farm directory.Farm) error {
	_, err := d.put(d.url("farms", fmt.Sprintf("%d", farm.ID)), farm, nil, http.StatusOK)
	return err
}

func (d *httpDirectory) FarmList(tid schema.ID, name string, page *Pager) (farms []directory.Farm, err error) {
	query := url.Values{}
	page.apply(query)
	if tid > 0 {
		query.Set("owner", fmt.Sprint(tid))
	}
	if len(name) != 0 {
		query.Set("name", name)
	}
	_, err = d.get(d.url("farms"), query, &farms, http.StatusOK)
	return
}

func (d *httpDirectory) FarmGet(id schema.ID) (farm directory.Farm, err error) {
	_, err = d.get(d.url("farms", fmt.Sprint(id)), nil, &farm, http.StatusOK)
	return
}

func (d *httpDirectory) FarmAddIP(id schema.ID, ip directory.PublicIP) error {
	_, err := d.post(d.url("farms", fmt.Sprintf("%d/ip", id)), []directory.PublicIP{ip}, nil, http.StatusOK)
	return err
}
func (d *httpDirectory) FarmDeleteIP(id schema.ID, ipaddr schema.IPCidr) error {
	_, err := d.deleteWithBody(d.url("farms", fmt.Sprintf("%d/ip", id)), []schema.IPCidr{ipaddr}, nil, http.StatusOK)
	return err
}

func (d *httpDirectory) Farms(cacheSize int) FarmIter {
	// pages start at index 1
	return &httpFarmIter{cl: d, size: cacheSize, page: 1}
}

func (fi *httpFarmIter) Next() (*directory.Farm, error) {
	// check if there are still cached farms
	if fi.cacheIdx >= len(fi.cache) {
		if fi.finished {
			return nil, nil
		}
		// pull new data in cache
		pager := Page(fi.page, fi.size)
		farms, err := fi.cl.FarmList(0, "", pager)
		if err != nil {
			return nil, errors.Wrap(err, "could not get farms")
		}
		if len(farms) == 0 {
			// iteration finished, no more  farms
			return nil, nil
		}
		fi.cache = farms
		fi.cacheIdx = 0
		fi.page++
		if len(farms) < fi.size {
			fi.finished = true
		}
	}
	fi.cacheIdx++
	return &fi.cache[fi.cacheIdx-1], nil
}

func (d *httpDirectory) NodeRegister(node directory.Node) error {
	_, err := d.post(d.url("nodes"), node, nil, http.StatusCreated)
	return err
}

func (d *httpDirectory) NodeList(filter NodeFilter, pager *Pager) (nodes []directory.Node, err error) {
	query := url.Values{}
	pager.apply(query)
	filter.Apply(query)
	_, err = d.get(d.url("nodes"), query, &nodes, http.StatusOK)
	return
}

func (d *httpDirectory) NodeGet(id string, proofs bool) (node directory.Node, err error) {
	query := url.Values{}
	query.Set("proofs", fmt.Sprint(proofs))
	_, err = d.get(d.url("nodes", id), query, &node, http.StatusOK)
	return
}

func (d *httpDirectory) NodeSetInterfaces(id string, ifaces []directory.Iface) error {
	_, err := d.post(d.url("nodes", id, "interfaces"), ifaces, nil, http.StatusCreated)
	return err
}

func (d *httpDirectory) NodeSetPorts(id string, ports []uint) error {
	var input struct {
		P []uint `json:"ports"`
	}
	input.P = ports

	_, err := d.post(d.url("nodes", id, "ports"), input, nil, http.StatusOK)
	return err
}

func (d *httpDirectory) NodeSetPublic(id string, pub directory.PublicIface) error {
	_, err := d.post(d.url("nodes", id, "configure_public"), pub, nil, http.StatusCreated)
	return err
}

func (d *httpDirectory) NodeSetCapacity(
	id string,
	resources directory.ResourceAmount,
	dmiInfo dmi.DMI,
	disksInfo capacity.Disks,
	hypervisor []string) error {

	payload := struct {
		Capacity   directory.ResourceAmount `json:"capacity"`
		DMI        dmi.DMI                  `json:"dmi"`
		Disks      capacity.Disks           `json:"disks"`
		Hypervisor []string                 `json:"hypervisor"`
	}{
		Capacity:   resources,
		DMI:        dmiInfo,
		Disks:      disksInfo,
		Hypervisor: hypervisor,
	}

	_, err := d.post(d.url("nodes", id, "capacity"), payload, nil, http.StatusOK)
	return err
}

func (d *httpDirectory) NodeUpdateUptime(id string, uptime uint64) error {
	input := struct {
		U uint64 `json:"uptime"`
	}{
		U: uptime,
	}

	_, err := d.post(d.url("nodes", id, "uptime"), input, nil, http.StatusOK)
	return err
}

func (d *httpDirectory) NodeUpdateUsedResources(id string, resources directory.ResourceAmount, workloads directory.WorkloadAmount) error {
	input := struct {
		directory.ResourceAmount
		directory.WorkloadAmount
	}{
		resources,
		workloads,
	}
	_, err := d.post(d.url("nodes", id, "used_resources"), input, nil, http.StatusOK)
	return err
}

func (d *httpDirectory) NodeSetFreeToUse(id string, free bool) error {
	choice := struct {
		FreeToUse bool `json:"free_to_use"`
	}{FreeToUse: free}

	_, err := d.post(d.url("nodes", id, "configure_free"), choice, nil, http.StatusOK)
	return err
}

func (d *httpDirectory) Nodes(cacheSize int, proofs bool) NodeIter {
	// pages start at index 1
	return &httpNodeIter{cl: d, size: cacheSize, page: 1, proofs: proofs}
}

func (ni *httpNodeIter) Next() (*directory.Node, error) {
	// check if there are still cached nodes
	if ni.cacheIdx >= len(ni.cache) {
		if ni.finished {
			return nil, nil
		}
		// pull new data in cache
		pager := Page(ni.page, ni.size)
		filter := NodeFilter{}.WithProofs(ni.proofs)
		nodes, err := ni.cl.NodeList(filter, pager)
		if err != nil {
			return nil, errors.Wrap(err, "could not get nodes")
		}
		if len(nodes) == 0 {
			// no more nodes, iteration finished
			return nil, nil
		}
		ni.cache = nodes
		ni.cacheIdx = 0
		ni.page++
		if len(nodes) < ni.size {
			ni.finished = true
		}
	}
	ni.cacheIdx++
	return &ni.cache[ni.cacheIdx-1], nil
}
func (d *httpDirectory) GatewayRegister(Gateway directory.Gateway) error {
	_, err := d.post(d.url("gateways"), Gateway, nil, http.StatusCreated)
	return err
}

func (d *httpDirectory) GatewayList(tid schema.ID, name string, page *Pager) (Gateways []directory.Gateway, err error) {
	query := url.Values{}
	page.apply(query)
	if len(name) != 0 {
		query.Set("name", name)
	}
	_, err = d.get(d.url("gateways"), query, &Gateways, http.StatusOK)
	return
}

func (d *httpDirectory) GatewayGet(id string) (Gateway directory.Gateway, err error) {
	_, err = d.get(d.url("gateways", id), nil, &Gateway, http.StatusOK)
	return
}

func (d *httpDirectory) GatewayUpdateUptime(id string, uptime uint64) error {
	input := struct {
		U uint64 `json:"uptime"`
	}{
		U: uptime,
	}

	_, err := d.post(d.url("gateways", id, "uptime"), input, nil, http.StatusOK)
	return err
}

func (d *httpDirectory) GatewayUpdateReservedResources(id string, resources directory.ResourceAmount, workloads directory.WorkloadAmount) error {
	input := struct {
		directory.ResourceAmount
		directory.WorkloadAmount
	}{
		resources,
		workloads,
	}

	_, err := d.post(d.url("gateways", id, "reserved_resources"), input, nil, http.StatusOK)
	return err
}
