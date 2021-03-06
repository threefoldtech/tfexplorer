package workloads

import (
	"bytes"
	"fmt"
)

var _ Workloader = (*Volume)(nil)
var _ Capaciter = (*Volume)(nil)

type Volume struct {
	ReservationInfo `bson:",inline"`

	Size int64          `bson:"size" json:"size"`
	Type VolumeTypeEnum `bson:"type" json:"type"`
}

func (v *Volume) GetRSU() (RSU, error) {
	switch v.Type {
	case VolumeTypeHDD:
		return RSU{
			HRU: float64(v.Size),
		}, nil
	case VolumeTypeSSD:
		return RSU{
			SRU: float64(v.Size),
		}, nil
	}
	return RSU{}, fmt.Errorf("volume type %s not supported", v.Type.String())
}

func (v *Volume) SignatureChallenge() ([]byte, error) {
	ric, err := v.ReservationInfo.SignatureChallenge()
	if err != nil {
		return nil, err
	}
	b := bytes.NewBuffer(ric)
	fmt.Fprintf(b, "%d", v.Size)
	fmt.Fprintf(b, "%s", v.Type.String())

	return b.Bytes(), nil
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
