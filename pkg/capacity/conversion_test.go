package capacity

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/threefoldtech/tfexplorer/models/generated/workloads"
)

func TestCloudUnitsFromResourceUnits(t *testing.T) {
	tests := []struct {
		rsu workloads.RSU
		cu  float64
		su  float64
	}{
		{
			rsu: workloads.RSU{
				CRU: 1,
				MRU: 1,
			},
			cu: 0.25,
			su: 0,
		},
		{
			rsu: workloads.RSU{
				CRU: 2,
				MRU: 4,
			},
			cu: 1,
			su: 0,
		},
		{
			rsu: workloads.RSU{
				CRU: 4,
				MRU: 8,
			},
			cu: 2,
			su: 0,
		},
		{
			rsu: workloads.RSU{
				CRU: 4,
				MRU: 64,
			},
			cu: 2,
			su: 0,
		},
		{
			rsu: workloads.RSU{
				CRU: 4,
				MRU: 32,
			},
			cu: 2,
			su: 0,
		},
		{
			rsu: workloads.RSU{
				SRU: 120,
				HRU: 1200,
			},
			cu: 0,
			su: 1.4,
		},
		{
			rsu: workloads.RSU{
				SRU: 40,
				HRU: 1000,
			},
			cu: 0,
			su: 0.967,
		},
		{
			rsu: workloads.RSU{
				SRU: 1200,
			},
			cu: 0,
			su: 4,
		},
		{
			rsu: workloads.RSU{
				HRU: 12000,
			},
			cu: 0,
			su: 10,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.rsu), func(t *testing.T) {
			cu, su, _ := CloudUnitsFromResourceUnits(tt.rsu)
			assert.Equal(t, tt.cu, cu, "wrong number of cu")
			assert.Equal(t, tt.su, su, "wrong number of su")
		})
	}
}
