package disk

import "syscall"

type SystemSpecs struct {
	TotalCapacity uint64
	FreeCapacity  uint64
	UsedCapacity  uint64
	DrivePath     string
}

func getSystemSpecs() (SystemSpecs, error) {
	var stat syscall.Statfs_t

	path, err := getAppSupportPath()
	if err != nil {
		return SystemSpecs{}, err
	}

	err = syscall.Statfs(path, &stat)
	if err != nil {
		return SystemSpecs{}, err
	}

	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bfree * uint64(stat.Bsize)
	used := total - free
	drive := path

	return SystemSpecs{
		TotalCapacity: total,
		FreeCapacity:  free,
		UsedCapacity:  used,
		DrivePath:     drive,
	}, err
}
