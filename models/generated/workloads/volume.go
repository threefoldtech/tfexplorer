package workloads

import schema "github.com/threefoldtech/tfexplorer/schema"

type Volume struct {
	WorkloadId      int64             `bson:"workload_id" json:"workload_id"`
	NodeId          string            `bson:"node_id" json:"node_id"`
	Size            int64             `bson:"size" json:"size"`
	Type            VolumeTypeEnum    `bson:"type" json:"type"`
	StatsAggregator []StatsAggregator `bson:"stats_aggregator" json:"stats_aggregator"`
	FarmerTid       int64             `bson:"farmer_tid" json:"farmer_tid"`
	PoolId          int64             `bson:"pool_id" json:"pool_id"`

	Description             string         `bson:"description" json:"description"`
	Currencies              []string       `bson:"currencies" json:"currencies"`
	SigningRequestProvision SigningRequest `bson:"signing_request_provision" json:"signing_request_provision"`
	SigningRequestDelete    SigningRequest `bson:"signing_request_delete" json:"signing_request_delete"`
	ExpirationProvisioning  schema.Date    `bson:"expiration_provisioning" json:"expiration_provisioning"`

	ID                  schema.ID          `bson:"_id" json:"id"`
	Json                string             `bson:"json" json:"json"`
	CustomerTid         int64              `bson:"customer_tid" json:"customer_tid"`
	CustomerSignature   string             `bson:"customer_signature" json:"customer_signature"`
	NextAction          NextActionEnum     `bson:"next_action" json:"next_action"`
	SignaturesProvision []SigningSignature `bson:"signatures_provision" json:"signatures_provision"`
	SignatureFarmer     SigningSignature   `bson:"signature_farmer" json:"signature_farmer"`
	SignaturesDelete    []SigningSignature `bson:"signatures_delete" json:"signatures_delete"`
	Epoch               schema.Date        `bson:"epoch" json:"epoch"`
	Metadata            string             `bson:"metadata" json:"metadata"`
	Result              Result             `bson:"result" json:"result"`
	WorkloadType        WorkloadTypeEnum   `bson:"workload_type" json:"workload_type"`
}

func (v *Volume) WorkloadID() int64 {
	return v.WorkloadId
}

func (v *Volume) GetWorkloadType() WorkloadTypeEnum {
	return v.WorkloadType
}

func (v *Volume) GetID() schema.ID {
	return v.ID
}

func (v *Volume) SetID(id schema.ID) {
	v.ID = id
}

func (v *Volume) GetJson() string {
	return v.Json
}

func (v *Volume) GetCustomerTid() int64 {
	return v.CustomerTid
}

func (v *Volume) GetCustomerSignature() string {
	return v.CustomerSignature
}

func (v *Volume) GetNextAction() NextActionEnum {
	return v.NextAction
}

func (v *Volume) SetNextAction(next NextActionEnum) {
	v.NextAction = next
}

func (v *Volume) GetSignaturesProvision() []SigningSignature {
	return v.SignaturesProvision
}

func (v *Volume) PushSignatureProvision(signature SigningSignature) {
	v.SignaturesProvision = append(v.SignaturesProvision, signature)
}

func (v *Volume) GetSignatureFarmer() SigningSignature {
	return v.SignatureFarmer
}

func (v *Volume) SetSignatureFarmer(signature SigningSignature) {
	v.SignatureFarmer = signature
}

func (v *Volume) GetSignaturesDelete() []SigningSignature {
	return v.SignaturesDelete
}

func (v *Volume) PushSignatureDelete(signature SigningSignature) {
	v.SignaturesDelete = append(v.SignaturesDelete, signature)
}

func (v *Volume) GetEpoch() schema.Date {
	return v.Epoch
}

func (v *Volume) GetMetadata() string {
	return v.Metadata
}

func (v *Volume) GetResult() Result {
	return v.Result
}

func (v *Volume) SetResult(result Result) {
	v.Result = result
}

func (v *Volume) GetDescription() string {
	return v.Description
}

func (v *Volume) GetCurrencies() []string {
	return v.Currencies
}

func (v *Volume) GetSigningRequestProvision() SigningRequest {
	return v.SigningRequestProvision
}

func (v *Volume) GetSigningRequestDelete() SigningRequest {
	return v.GetSigningRequestDelete()
}

func (v *Volume) GetExpirationProvisioning() schema.Date {
	return v.ExpirationProvisioning
}

type VolumeTypeEnum uint8

const (
	VolumeTypeHDD VolumeTypeEnum = iota
	VolumeTypeSSD
)

func (e VolumeTypeEnum) String() string {
	switch e {
	case VolumeTypeHDD:
		return "HDD"
	case VolumeTypeSSD:
		return "SSD"
	}
	return "UNKNOWN"
}
