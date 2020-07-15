package main

import (
	"fmt"
	"strings"

	"github.com/threefoldtech/tfexplorer/models/generated/workloads"
	"github.com/threefoldtech/tfexplorer/provision/builders"

	"github.com/urfave/cli"
)

func generateVolume(c *cli.Context) error {
	s := c.Int64("size")
	t := strings.ToLower(c.String("type"))

	if t != workloads.DiskTypeHDD.String() && t != workloads.DiskTypeSSD.String() {
		return fmt.Errorf("volume type can only hdd or ssd")
	}

	if s < 1 { //TODO: upper bound ?
		return fmt.Errorf("size cannot be less then 1")
	}

	var volumeType workloads.VolumeTypeEnum
	if t == workloads.DiskTypeHDD.String() {
		volumeType = workloads.VolumeTypeEnum(workloads.VolumeTypeHDD)
	} else if t == workloads.DiskTypeSSD.String() {
		volumeType = workloads.VolumeTypeEnum(workloads.VolumeTypeSSD)
	}

	volumeBuilder := builders.NewVolumeBuilder(c.String("node"), s, volumeType)
	volumeBuilder.WithSize(s)

	if c.Int64("poolID") != 0 {
		volumeBuilder.WithPoolID(c.Int64("poolID"))
	}

	return writeWorkload(c.GlobalString("output"), volumeBuilder.Build())
}
