// Copyright 2021 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package spec

import (
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/cockroach/pkg/roachprod/vm"
	"github.com/cockroachdb/cockroach/pkg/roachprod/vm/aws"
	"github.com/cockroachdb/cockroach/pkg/roachprod/vm/azure"
	"github.com/cockroachdb/cockroach/pkg/roachprod/vm/gce"
	"github.com/cockroachdb/cockroach/pkg/util/randutil"
	"github.com/cockroachdb/cockroach/pkg/util/timeutil"
	"github.com/cockroachdb/errors"
)

type fileSystemType int

const (
	// Since ext4 is the default of 0, it isn't being
	// used anywhere in the code. Therefore, it isn't
	// added as a const here since it is unused, and
	// leads to a lint error.

	// Zfs file system.
	Zfs fileSystemType = 1
)

type MemPerCPU int

const (
	Auto MemPerCPU = iota
	Standard
	High
	Low
)

func (m MemPerCPU) String() string {
	switch m {
	case Auto:
		return "auto"
	case Standard:
		return "standard"
	case High:
		return "high"
	case Low:
		return "low"
	}
	return "unknown"
}

// ClusterSpec represents a test's description of what its cluster needs to
// look like. It becomes part of a clusterConfig when the cluster is created.
type ClusterSpec struct {
	// TODO(#104029): We should remove the Cloud field; the tests now specify
	// their compatible clouds.
	Cloud string

	// TODO(radu): defaultInstanceType is the default machine type used (unless
	// overridden by GCE.MachineType or AWS.MachineType); it does not belong in
	// the spec.
	defaultInstanceType string

	// TODO(radu): defaultZones is the default zones specification (unless
	// overridden by GCE.Zones or AWS.Zones); it does not belong in the spec.
	defaultZones string

	Arch      vm.CPUArch // CPU architecture; auto-chosen if left empty
	NodeCount int
	// CPUs is the number of CPUs per node.
	CPUs                 int
	Mem                  MemPerCPU
	SSDs                 int
	RAID0                bool
	VolumeSize           int
	PreferLocalSSD       bool
	Geo                  bool
	Lifetime             time.Duration
	ReusePolicy          clusterReusePolicy
	TerminateOnMigration bool
	UbuntuVersion        vm.UbuntuVersion

	// FileSystem determines the underlying FileSystem
	// to be used. The default is ext4.
	FileSystem fileSystemType

	RandomlyUseZfs bool

	GatherCores bool

	// GCE-specific arguments. These values apply only on clusters instantiated on GCE.
	GCE struct {
		MachineType    string
		MinCPUPlatform string
		VolumeType     string
		Zones          string
	}

	// AWS-specific arguments. These values apply only on clusters instantiated on AWS.
	AWS struct {
		MachineType string
		// VolumeThroughput is the min provisioned EBS volume throughput.
		VolumeThroughput int
		Zones            string
	}
}

// MakeClusterSpec makes a ClusterSpec.
func MakeClusterSpec(cloud string, instanceType string, nodeCount int, opts ...Option) ClusterSpec {
	spec := ClusterSpec{Cloud: cloud, defaultInstanceType: instanceType, NodeCount: nodeCount}
	defaultOpts := []Option{CPU(4), nodeLifetime(12 * time.Hour), ReuseAny()}
	for _, o := range append(defaultOpts, opts...) {
		o(&spec)
	}
	return spec
}

// ClustersCompatible returns true if the clusters are compatible, i.e. the test
// asking for s2 can reuse s1.
func ClustersCompatible(s1, s2 ClusterSpec) bool {
	s1.Lifetime = 0
	s2.Lifetime = 0
	return s1 == s2
}

// String implements fmt.Stringer.
func (s ClusterSpec) String() string {
	str := fmt.Sprintf("n%dcpu%d", s.NodeCount, s.CPUs)
	switch s.Mem {
	case Standard:
		str += "sm"
	case High:
		str += "hm"
	case Low:
		str += "lm"
	}
	if s.Geo {
		str += "-Geo"
	}
	return str
}

// checks if an AWS machine supports SSD volumes
func awsMachineSupportsSSD(machineType string) bool {
	typeAndSize := strings.Split(machineType, ".")
	if len(typeAndSize) == 2 {
		// All SSD machine types that we use end in 'd or begins with i3 (e.g. i3, i3en).
		return strings.HasPrefix(typeAndSize[0], "i3") || strings.HasSuffix(typeAndSize[0], "d")
	}
	return false
}

func getAWSOpts(
	machineType string, zones []string, volumeSize, ebsThroughput int, localSSD bool,
) vm.ProviderOpts {
	opts := aws.DefaultProviderOpts()
	if volumeSize != 0 {
		opts.DefaultEBSVolume.Disk.VolumeSize = volumeSize
	}
	if ebsThroughput != 0 {
		opts.DefaultEBSVolume.Disk.Throughput = ebsThroughput
	}
	if localSSD {
		opts.SSDMachineType = machineType
	} else {
		opts.MachineType = machineType
	}
	if len(zones) != 0 {
		opts.CreateZones = zones
	}
	return opts
}

func getGCEOpts(
	machineType string,
	zones []string,
	volumeSize, localSSDCount int,
	localSSD bool,
	RAID0 bool,
	terminateOnMigration bool,
	minCPUPlatform string,
	arch vm.CPUArch,
	volumeType string,
) vm.ProviderOpts {
	opts := gce.DefaultProviderOpts()
	opts.MachineType = machineType
	if arch == vm.ArchARM64 {
		// ARM64 machines don't support minCPUPlatform.
		opts.MinCPUPlatform = ""
	} else if minCPUPlatform != "" {
		opts.MinCPUPlatform = minCPUPlatform
	}
	if volumeSize != 0 {
		opts.PDVolumeSize = volumeSize
	}
	if len(zones) != 0 {
		opts.Zones = zones
	}
	opts.SSDCount = localSSDCount
	if localSSD && localSSDCount > 0 {
		// NB: As the default behavior for _roachprod_ (at least in AWS/GCP) is
		// to mount multiple disks as a single store using a RAID 0 array, we
		// must explicitly ask for multiple stores to be enabled, _unless_ the
		// test has explicitly asked for RAID0.
		opts.UseMultipleDisks = !RAID0
	}
	opts.TerminateOnMigration = terminateOnMigration
	if volumeType != "" {
		opts.PDVolumeType = volumeType
	}

	return opts
}

func getAzureOpts(machineType string, zones []string) vm.ProviderOpts {
	opts := azure.DefaultProviderOpts()
	opts.MachineType = machineType
	if len(zones) != 0 {
		opts.Locations = zones
	}
	return opts
}

// RoachprodClusterConfig contains general roachprod cluster configuration that
// does not depend on the test. It is used in conjunction with ClusterSpec to
// determine the final configuration.
type RoachprodClusterConfig struct {
	// UseIOBarrierOnLocalSSD is set if we don't want to mount local SSDs with the
	// `-o nobarrier` flag.
	UseIOBarrierOnLocalSSD bool

	// PreferredArch is the preferred CPU architecture; it is not guaranteed
	// (depending on cloud and on other requirements on machine type).
	PreferredArch vm.CPUArch
}

// RoachprodOpts returns the opts to use when calling `roachprod.Create()`
// in order to create the cluster described in the spec.
func (s *ClusterSpec) RoachprodOpts(
	params RoachprodClusterConfig,
) (vm.CreateOpts, vm.ProviderOpts, error) {
	useIOBarrier := params.UseIOBarrierOnLocalSSD
	arch := params.PreferredArch

	createVMOpts := vm.DefaultCreateOpts()
	// N.B. We set "usage=roachtest" as the default, custom label for billing tracking.
	createVMOpts.CustomLabels = map[string]string{"usage": "roachtest"}
	createVMOpts.ClusterName = "" // Will be set later.
	if s.Lifetime != 0 {
		createVMOpts.Lifetime = s.Lifetime
	}
	switch s.Cloud {
	case Local:
		createVMOpts.VMProviders = []string{s.Cloud}
		// remaining opts are not applicable to local clusters
		return createVMOpts, nil, nil
	case AWS, GCE, Azure:
		createVMOpts.VMProviders = []string{s.Cloud}
	default:
		return vm.CreateOpts{}, nil, errors.Errorf("unsupported cloud %v", s.Cloud)
	}

	if s.Cloud != GCE && s.Cloud != AWS {
		if s.VolumeSize != 0 {
			return vm.CreateOpts{}, nil, errors.Errorf("specifying volume size is not yet supported on %s", s.Cloud)
		}
	}
	if s.Cloud != GCE {
		if s.SSDs != 0 {
			return vm.CreateOpts{}, nil, errors.Errorf("specifying SSD count is not yet supported on %s", s.Cloud)
		}
	}

	createVMOpts.GeoDistributed = s.Geo
	createVMOpts.Arch = string(arch)
	ssdCount := s.SSDs

	machineType := s.defaultInstanceType
	switch s.Cloud {
	case AWS:
		if s.AWS.MachineType != "" {
			machineType = s.AWS.MachineType
		}
	case GCE:
		if s.GCE.MachineType != "" {
			machineType = s.GCE.MachineType
		}
	}

	if s.CPUs != 0 {
		// Default to the user-supplied machine type, if any.
		// Otherwise, pick based on requested CPU count.
		var selectedArch vm.CPUArch

		if machineType == "" {
			// If no machine type was specified, choose one
			// based on the cloud and CPU count.
			var err error
			switch s.Cloud {
			case AWS:
				machineType, selectedArch, err = SelectAWSMachineType(s.CPUs, s.Mem, s.PreferLocalSSD && s.VolumeSize == 0, arch)
			case GCE:
				machineType, selectedArch = SelectGCEMachineType(s.CPUs, s.Mem, arch)
			case Azure:
				machineType, err = SelectAzureMachineType(s.CPUs, s.Mem, s.PreferLocalSSD)
			}

			if err != nil {
				return vm.CreateOpts{}, nil, err
			}
			if selectedArch != "" && selectedArch != arch {
				// TODO(srosenberg): we need a better way to monitor the rate of this mismatch, i.e.,
				// other than grepping cluster creation logs.
				fmt.Printf("WARN: requested arch %s for machineType %s, but selected %s\n", arch, machineType, selectedArch)
				createVMOpts.Arch = string(selectedArch)
			}
		}

		// Local SSD can only be requested
		// - if configured to prefer doing so,
		// - if no particular volume size is requested, and,
		// - on AWS, if the machine type supports it.
		// - on GCE, if the machine type is not ARM64.
		if s.PreferLocalSSD && s.VolumeSize == 0 && (s.Cloud != AWS || awsMachineSupportsSSD(machineType)) &&
			(s.Cloud != GCE || selectedArch != vm.ArchARM64) {
			// Ensure SSD count is at least 1 if UseLocalSSD is true.
			if ssdCount == 0 {
				ssdCount = 1
			}
			createVMOpts.SSDOpts.UseLocalSSD = true
			createVMOpts.SSDOpts.NoExt4Barrier = !useIOBarrier
		} else {
			createVMOpts.SSDOpts.UseLocalSSD = false
		}
	}

	if s.FileSystem == Zfs {
		if s.Cloud != GCE {
			return vm.CreateOpts{}, nil, errors.Errorf(
				"node creation with zfs file system not yet supported on %s", s.Cloud,
			)
		}
		createVMOpts.SSDOpts.FileSystem = vm.Zfs
	} else if s.RandomlyUseZfs && s.Cloud == GCE {
		rng, _ := randutil.NewPseudoRand()
		if rng.Float64() <= 0.2 {
			createVMOpts.SSDOpts.FileSystem = vm.Zfs
		}
	}

	zonesStr := s.defaultZones
	switch s.Cloud {
	case AWS:
		if s.AWS.Zones != "" {
			zonesStr = s.AWS.Zones
		}
	case GCE:
		if s.GCE.Zones != "" {
			zonesStr = s.GCE.Zones
		}
	}
	var zones []string
	if zonesStr != "" {
		zones = strings.Split(zonesStr, ",")
		if !s.Geo {
			zones = zones[:1]
		}
	}

	if createVMOpts.Arch == string(vm.ArchFIPS) && !(s.Cloud == GCE || s.Cloud == AWS) {
		return vm.CreateOpts{}, nil, errors.Errorf(
			"FIPS not yet supported on %s", s.Cloud,
		)
	}
	var providerOpts vm.ProviderOpts
	switch s.Cloud {
	case AWS:
		providerOpts = getAWSOpts(machineType, zones, s.VolumeSize, s.AWS.VolumeThroughput,
			createVMOpts.SSDOpts.UseLocalSSD)
	case GCE:
		providerOpts = getGCEOpts(machineType, zones, s.VolumeSize, ssdCount,
			createVMOpts.SSDOpts.UseLocalSSD, s.RAID0, s.TerminateOnMigration,
			s.GCE.MinCPUPlatform, vm.ParseArch(createVMOpts.Arch), s.GCE.VolumeType,
		)
	case Azure:
		providerOpts = getAzureOpts(machineType, zones)
	}

	return createVMOpts, providerOpts, nil
}

// Expiration is the lifetime of the cluster. It may be destroyed after
// the expiration has passed.
func (s *ClusterSpec) Expiration() time.Time {
	l := s.Lifetime
	if l == 0 {
		l = 12 * time.Hour
	}
	return timeutil.Now().Add(l)
}
