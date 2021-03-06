package workloads

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContainer_GetRSU(t *testing.T) {
	tests := []struct {
		capcity ContainerCapacity
		rsu     RSU
	}{
		{
			capcity: ContainerCapacity{
				Cpu:      1,
				Memory:   1024,
				DiskSize: 256,
				DiskType: DiskTypeSSD,
			},
			rsu: RSU{
				CRU: 1,
				MRU: 1,
				SRU: 0,
			},
		},
		{
			capcity: ContainerCapacity{
				Cpu:      1,
				Memory:   1024,
				DiskSize: 1024,
				DiskType: DiskTypeSSD,
			},
			rsu: RSU{
				CRU: 1,
				MRU: 1,
				SRU: 0,
			},
		},
		{
			capcity: ContainerCapacity{
				Cpu:      4,
				Memory:   2048,
				DiskSize: 10240,
				DiskType: DiskTypeHDD,
			},
			rsu: RSU{
				CRU: 4,
				MRU: 2,
				SRU: 0,
				HRU: 0,
			},
		},
		{
			capcity: ContainerCapacity{
				Cpu:      1,
				Memory:   200,
				DiskSize: 10000,
				DiskType: DiskTypeHDD,
			},
			rsu: RSU{
				CRU: 1,
				MRU: 0.1953,
				SRU: 0,
				HRU: 0,
			},
		},
		{
			capcity: ContainerCapacity{
				Cpu:      1,
				Memory:   200,
				DiskSize: 52224,
				DiskType: DiskTypeSSD,
			},
			rsu: RSU{
				CRU: 1,
				MRU: 0.1953,
				SRU: 1,
				HRU: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%+v", tt.capcity), func(t *testing.T) {
			rsu, err := tt.capcity.GetRSU()
			require.NoError(t, err)
			assert.Equal(t, tt.rsu, rsu)
		})
	}
}
