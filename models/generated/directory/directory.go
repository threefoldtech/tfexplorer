package directory

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/threefoldtech/tfexplorer/models/generated/phonebook"
	schema "github.com/threefoldtech/tfexplorer/schema"
)

type NodeResourcePrice struct {
	Currency PriceCurrencyEnum `bson:"currency" json:"currency"`
	Cru      float64           `bson:"cru" json:"cru"`
	Mru      float64           `bson:"mru" json:"mru"`
	Hru      float64           `bson:"hru" json:"hru"`
	Sru      float64           `bson:"sru" json:"sru"`
	Nru      float64           `bson:"nru" json:"nru"`
}

type NodeCloudUnitPrice struct {
	Currency PriceCurrencyEnum `bson:"currency" json:"currency"`
	CU       float64           `bson:"cu" json:"cu"`
	SU       float64           `bson:"su" json:"su"`
	NU       float64           `bson:"nu" json:"nu"`
	IPv4U    float64           `bson:"ipv4u" json:"ipv4u"`
}

func NewNodeCloudUnitPrice() NodeCloudUnitPrice {
	return NodeCloudUnitPrice{Currency: PriceCurrencyTFT}
}

type Farm struct {
	ID                  schema.ID           `bson:"_id" json:"id"`
	ThreebotID          int64               `bson:"threebot_id" json:"threebot_id"`
	IyoOrganization     string              `bson:"iyo_organization" json:"iyo_organization"`
	Name                string              `bson:"name" json:"name"`
	WalletAddresses     []WalletAddress     `bson:"wallet_addresses" json:"wallet_addresses"`
	Location            Location            `bson:"location" json:"location"`
	Email               schema.Email        `bson:"email" json:"email"`
	ResourcePrices      []NodeResourcePrice `bson:"resource_prices" json:"resource_prices"`
	PrefixZero          schema.IPRange      `bson:"prefix_zero" json:"prefix_zero"`
	IPAddresses         []PublicIP          `bson:"ipaddresses" json:"ipaddresses"`
	EnableCustomPricing bool                `bson:"enable_custom_pricing" json:"enable_custom_pricing"`
	FarmCloudUnitsPrice NodeCloudUnitPrice  `bson:"farm_cloudunits_price" json:"farm_cloudunits_price"`

	// Grid3 pricing enabled
	IsGrid3Compliant bool `bson:"is_grid3_compliant" json:"is_grid3_compliant"`
}
type FarmThreebotPrice struct {
	ThreebotID           int64              `bson:"threebot_id" json:"threebot_id"`
	FarmID               int64              `bson:"farm_id" json:"farm_id"`
	CustomCloudUnitPrice NodeCloudUnitPrice `bson:"custom_cloudunits_price" json:"custom_cloudunits_price"`
}

func NewFarm() (Farm, error) {
	const value = "{}"
	var object Farm
	if err := json.Unmarshal([]byte(value), &object); err != nil {
		return object, err
	}
	return object, nil
}

type WalletAddress = phonebook.WalletAddress

func NewNodeResourcePrice() (NodeResourcePrice, error) {
	const value = "{}"
	var object NodeResourcePrice
	if err := json.Unmarshal([]byte(value), &object); err != nil {
		return object, err
	}
	return object, nil
}

type Location struct {
	City      string  `bson:"city" json:"city"`
	Country   string  `bson:"country" json:"country"`
	Continent string  `bson:"continent" json:"continent"`
	Latitude  float64 `bson:"latitude" json:"latitude"`
	Longitude float64 `bson:"longitude" json:"longitude"`
}

func NewLocation() (Location, error) {
	const value = "{}"
	var object Location
	if err := json.Unmarshal([]byte(value), &object); err != nil {
		return object, err
	}
	return object, nil
}

type Node struct {
	ID                schema.ID      `bson:"_id" json:"id"`
	NodeId            string         `bson:"node_id" json:"node_id"`
	HostName          string         `bson:"hostname" json:"hostname"`
	NodeIdV1          string         `bson:"node_id_v1" json:"node_id_v1"`
	FarmId            int64          `bson:"farm_id" json:"farm_id"`
	OsVersion         string         `bson:"os_version" json:"os_version"`
	Created           schema.Date    `bson:"created" json:"created"`
	Updated           schema.Date    `bson:"updated" json:"updated"`
	Uptime            int64          `bson:"uptime" json:"uptime"`
	Address           string         `bson:"address" json:"address"`
	Location          Location       `bson:"location" json:"location"`
	TotalResources    ResourceAmount `bson:"total_resources" json:"total_resources"`
	UsedResources     ResourceAmount `bson:"used_resources" json:"used_resources"`
	ReservedResources ResourceAmount `bson:"reserved_resources" json:"reserved_resources"`
	Workloads         WorkloadAmount `bson:"workloads" json:"workloads"`
	Proofs            []Proof        `bson:"proofs" json:"proofs"`
	Ifaces            []Iface        `bson:"ifaces" json:"ifaces"`
	PublicConfig      *PublicIface   `bson:"public_config,omitempty" json:"public_config"`
	FreeToUse         bool           `bson:"free_to_use" json:"free_to_use"`
	Approved          bool           `bson:"approved" json:"approved"`
	PublicKeyHex      string         `bson:"public_key_hex" json:"public_key_hex"`
	WgPorts           []int64        `bson:"wg_ports" json:"wg_ports"`
	Deleted           bool           `bson:"deleted" json:"deleted"`
	Reserved          bool           `bson:"reserved" json:"reserved"`
}

func NewNode() (Node, error) {
	const value = "{\"approved\": false, \"public_key_hex\": \"\"}"
	var object Node
	if err := json.Unmarshal([]byte(value), &object); err != nil {
		return object, err
	}
	return object, nil
}

type Iface struct {
	Name       string            `bson:"name" json:"name"`
	Addrs      []schema.IPRange  `bson:"addrs" json:"addrs"`
	Gateway    []net.IP          `bson:"gateway" json:"gateway"`
	MacAddress schema.MacAddress `bson:"macaddress" json:"macaddress"`
}

func NewIface() (Iface, error) {
	const value = "{}"
	var object Iface
	if err := json.Unmarshal([]byte(value), &object); err != nil {
		return object, err
	}
	return object, nil
}

type PublicIface struct {
	Master  string         `bson:"master" json:"master"`
	Type    IfaceTypeEnum  `bson:"type" json:"type"`
	Ipv4    schema.IPRange `bson:"ipv4" json:"ipv4"`
	Ipv6    schema.IPRange `bson:"ipv6" json:"ipv6"`
	Gw4     net.IP         `bson:"gw4" json:"gw4"`
	Gw6     net.IP         `bson:"gw6" json:"gw6"`
	Version int64          `bson:"version" json:"version"`
}

// Validate check if all the value of the object are valid
func (p PublicIface) Validate() error {
	if p.Master == "" {
		return fmt.Errorf("master field cannot be empty")
	}
	if len(p.Master) > 16 {
		return fmt.Errorf("master field should contain the name of a network interface. A network interface cannot be longer than 16 characters")
	}

	if p.Type != IfaceTypeMacvlan {
		return fmt.Errorf("type can only be of type macvlan")
	}

	if p.Ipv4.IP == nil || p.Ipv4.Mask == nil {
		return fmt.Errorf("ipv4 cannot be empty")
	}

	if p.Gw4 == nil {
		return fmt.Errorf("gw4 cannot be empty")
	}

	if p.Ipv6.IP == nil || p.Ipv6.Mask == nil {
		return fmt.Errorf("ipv6 cannot be empty")
	}

	if p.Gw6 == nil {
		return fmt.Errorf("gw6 cannot be empty")
	}

	if p.Ipv4.IP.To4() == nil {
		return fmt.Errorf("%s is not a valid IPv4 address", p.Ipv4.IP.String())
	}

	_, bits := p.Ipv4.Mask.Size()
	if bits != 32 {
		return fmt.Errorf("%s is not a valid IPv4 net mask", p.Ipv4.Mask.String())
	}

	if p.Gw4.To4() == nil {
		return fmt.Errorf("%s is not a valid IPv4 address", p.Gw4.String())
	}

	if p.Ipv6.IP.To4() != nil {
		return fmt.Errorf("%s is not a valid IPv6 address", p.Ipv6.IP.String())
	}

	_, bits = p.Ipv6.Mask.Size()
	if bits != 128 {
		return fmt.Errorf("%s is not a valid IPv6 net mask", p.Ipv6.Mask.String())
	}

	if p.Gw6.To4() != nil {
		return fmt.Errorf("%s is not a valid IPv6 address", p.Gw6.String())
	}

	return nil
}

type ResourceAmount struct {
	Cru uint64  `bson:"cru" json:"cru"`
	Mru float64 `bson:"mru" json:"mru"`
	Hru float64 `bson:"hru" json:"hru"`
	Sru float64 `bson:"sru" json:"sru"`
}

// Diff returns r - v
func (r *ResourceAmount) Diff(v ResourceAmount) ResourceAmount {
	return ResourceAmount{
		Cru: r.Cru - v.Cru,
		Mru: r.Mru - v.Mru,
		Hru: r.Mru - v.Hru,
		Sru: r.Sru - v.Sru,
	}
}

func NewResourceAmount() (ResourceAmount, error) {
	const value = "{}"
	var object ResourceAmount
	if err := json.Unmarshal([]byte(value), &object); err != nil {
		return object, err
	}
	return object, nil
}

type WorkloadAmount struct {
	Network         uint16 `bson:"network" json:"network"`
	NetworkResource uint16 `bson:"network_resource" json:"network_resource"`
	Volume          uint16 `bson:"volume" json:"volume"`
	ZDBNamespace    uint16 `bson:"zdb_namespace" json:"zdb_namespace"`
	Container       uint16 `bson:"container" json:"container"`
	K8sVM           uint16 `bson:"k8s_vm" json:"k8s_vm"`
	GenericVM       uint16 `bson:"generic_vm" json:"generic_vm"`
	Proxy           uint16 `bson:"proxy" json:"proxy"`
	ReverseProxy    uint16 `bson:"reverse_proxy" json:"reverse_proxy"`
	Subdomain       uint16 `bson:"subdomain" json:"subdomain"`
	DelegateDomain  uint16 `bson:"delegate_domain" json:"delegate_domain"`
}

type Proof struct {
	Created      schema.Date            `bson:"created" json:"created"`
	HardwareHash string                 `bson:"hardware_hash" json:"hardware_hash"`
	DiskHash     string                 `bson:"disk_hash" json:"disk_hash"`
	Hardware     map[string]interface{} `bson:"hardware" json:"hardware"`
	Disks        map[string]interface{} `bson:"disks" json:"disks"`
	Hypervisor   []string               `bson:"hypervisor" json:"hypervisor"`
}

func NewProof() (Proof, error) {
	const value = "{}"
	var object Proof
	if err := json.Unmarshal([]byte(value), &object); err != nil {
		return object, err
	}
	return object, nil
}

type IfaceTypeEnum uint8

const (
	IfaceTypeMacvlan IfaceTypeEnum = iota
	IfaceTypeVlan
)

func (e IfaceTypeEnum) String() string {
	switch e {
	case IfaceTypeMacvlan:
		return "macvlan"
	case IfaceTypeVlan:
		return "vlan"
	}
	return "UNKNOWN"
}

type PriceCurrencyEnum uint8

const (
	PriceCurrencyEUR PriceCurrencyEnum = iota
	PriceCurrencyUSD
	PriceCurrencyTFT
	PriceCurrencyAED
	PriceCurrencyGBP
)

func (e PriceCurrencyEnum) String() string {
	switch e {
	case PriceCurrencyEUR:
		return "EUR"
	case PriceCurrencyUSD:
		return "USD"
	case PriceCurrencyTFT:
		return "TFT"
	case PriceCurrencyAED:
		return "AED"
	case PriceCurrencyGBP:
		return "GBP"
	}
	return "UNKNOWN"
}

type Gateway struct {
	ID             schema.ID      `bson:"_id" json:"id"`
	NodeId         string         `bson:"node_id" json:"node_id"`
	FarmId         int64          `bson:"farm_id" json:"farm_id"`
	OsVersion      string         `bson:"os_version" json:"os_version"`
	Created        schema.Date    `bson:"created" json:"created"`
	Updated        schema.Date    `bson:"updated" json:"updated"`
	Uptime         int64          `bson:"uptime" json:"uptime"`
	Address        string         `bson:"address" json:"address"`
	Location       Location       `bson:"location" json:"location"`
	PublicKeyHex   string         `bson:"public_key_hex" json:"public_key_hex"`
	Workloads      WorkloadAmount `bson:"workloads" json:"workloads"`
	ManagedDomains []string       `bson:"managed_domains" json:"managed_domains"`
	TcpRouterPort  int64          `bson:"tcp_router_port" json:"tcp_router_port"`
	DnsNameserver  []string       `bson:"dns_nameserver" json:"dns_nameserver"`
	FreeToUse      bool           `bson:"free_to_use" json:"free_to_use"`
}

// PublicIP structure
type PublicIP struct {
	Address       schema.IPCidr `bson:"address" json:"address"`
	Gateway       schema.IP     `bson:"gateway" json:"gateway"`
	ReservationID schema.ID     `bson:"reservation_id" json:"reservation_id"`
}

func (ip *PublicIP) Valid() error {
	if len(ip.Address.IP) == 0 || len(ip.Address.Mask) == 0 {
		return fmt.Errorf("invalid public ip address")
	}

	if len(ip.Gateway.IP) == 0 {
		return fmt.Errorf("invalid public ip gateway")
	}

	return nil
}
