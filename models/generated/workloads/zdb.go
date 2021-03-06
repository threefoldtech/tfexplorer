package workloads

import (
	"bytes"
	"fmt"
)

var _ Workloader = (*ZDB)(nil)
var _ Capaciter = (*ZDB)(nil)

type ZDB struct {
	ReservationInfo `bson:",inline"`

	Size            int64             `bson:"size" json:"size"`
	Mode            ZDBModeEnum       `bson:"mode" json:"mode"`
	Password        string            `bson:"password" json:"password"`
	DiskType        DiskTypeEnum      `bson:"disk_type" json:"disk_type"`
	Public          bool              `bson:"public" json:"public"`
	StatsAggregator []StatsAggregator `bson:"stats_aggregator" json:"stats_aggregator"`
}

func (z *ZDB) GetRSU() (RSU, error) {
	switch z.DiskType {
	case DiskTypeHDD:
		return RSU{
			HRU: float64(z.Size),
		}, nil
	case DiskTypeSSD:
		return RSU{
			SRU: float64(z.Size),
		}, nil
	}
	return RSU{}, fmt.Errorf("ZDB disk type %s not supported", z.DiskType.String())
}

func (z *ZDB) SignatureChallenge() ([]byte, error) {
	ric, err := z.ReservationInfo.SignatureChallenge()
	if err != nil {
		return nil, err
	}

	b := bytes.NewBuffer(ric)
	if _, err := fmt.Fprintf(b, "%d", z.Size); err != nil {
		return nil, err
	}
	if _, err := fmt.Fprintf(b, "%s", z.Mode.String()); err != nil {
		return nil, err
	}
	if _, err := fmt.Fprintf(b, "%s", z.Password); err != nil {
		return nil, err
	}
	if _, err := fmt.Fprintf(b, "%s", z.DiskType.String()); err != nil {
		return nil, err
	}
	if _, err := fmt.Fprintf(b, "%t", z.Public); err != nil {
		return nil, err
	}
	for _, s := range z.StatsAggregator {
		if err := s.SigningEncode(b); err != nil {
			return nil, err
		}
	}

	return b.Bytes(), nil
}

type DiskTypeEnum uint8

const (
	DiskTypeHDD DiskTypeEnum = iota
	DiskTypeSSD
)

func (e DiskTypeEnum) String() string {
	switch e {
	case DiskTypeHDD:
		return "hdd"
	case DiskTypeSSD:
		return "ssd"
	}
	return "UNKNOWN"
}

type ZDBModeEnum uint8

const (
	ZDBModeSeq ZDBModeEnum = iota
	ZDBModeUser
)

func (e ZDBModeEnum) String() string {
	switch e {
	case ZDBModeSeq:
		return "seq"
	case ZDBModeUser:
		return "user"
	}
	return "UNKNOWN"
}
