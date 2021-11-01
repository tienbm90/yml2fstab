package main

import "fmt"

var FileSystemTypes = []string{"ext", "ext2", "ext3", "ext4", "jfs", "reiserfs", "xfs",
	"btrfs", "swap", "iso9660", "nfs", "nfs4", "udf", "vfat", "devpts"}

type FstabLine struct {
	//usually the given name or UUID of the mounted device
	Device string
	//designates the directory where the device is/will be mounted.
	MountPoint string

	// shows the type of filesystem in use.
	FileSystemType string
	// lists any active mount options. If using multiple options they must be separated by commas.
	Options string

	//1 = dump utility backup of a partition. 0 = no backup. This is an outdated backup method and should NOT be used.
	BackupOperation int

	//0 means that fsck will not check the filesystem. Numbers higher than this represent the check order. The root filesystem should be set to 1 and other partitions set to 2
	FileSystemCheckOrder int
}

//
func (ent *FstabLine) IsFileSystemTypeValid() bool {
	return CheckFileSystemTypeValid(ent.FileSystemType)
}

func (ent *FstabLine) IsBackupOperationValid() bool {
	if ent.BackupOperation < 0 || ent.BackupOperation > 1 {
		return false
	}
	return true
}

func (ent *FstabLine) IsFileSystemCheckOrderValid() bool {
	if ent.FileSystemCheckOrder < 0 || ent.FileSystemCheckOrder > 2 {
		return false
	}

	return true
}

func (ent *FstabLine) IsValid() bool {
	return ent.IsBackupOperationValid() &&
		ent.IsFileSystemCheckOrderValid() &&
		ent.IsFileSystemTypeValid() &&
		ent.IsMountPointValid()
}

func (ent *FstabLine) IsMountPointValid() bool {
	return CheckMountPointValid(ent.MountPoint)
}

func (ent *FstabLine) SetDevice(dvc string) {
	ent.Device = dvc
}

func (ent *FstabLine) SetMountPoint(mp string) {
	ent.MountPoint = mp
}

func (ent *FstabLine) SetFileSystemType(fst string) {
	ent.FileSystemType = fst
}

func (ent *FstabLine) SetOptions(opts string) {
	ent.Options = opts
}

func (ent *FstabLine) SetBackupOperation(bo int) {
	ent.BackupOperation = bo
}

func (ent *FstabLine) SetFileSystemCheckOrder(fsco int) {
	ent.FileSystemCheckOrder = fsco

}

func (ent *FstabLine) GenerateFstabEntryString() string {
	return fmt.Sprintf("%s %s %s %s %d %d",
		ent.Device,
		ent.MountPoint,
		ent.FileSystemType,
		ent.Options,
		ent.BackupOperation,
		ent.FileSystemCheckOrder)
}

func NewFstabEntry(device string,
	mountPoint string,
	fileSystemType string,
	options string,
	backupOperation int,
	fileSystemCheckOrder int) *FstabLine {
	return &FstabLine{
		Device:               device,
		MountPoint:           mountPoint,
		FileSystemType:       fileSystemType,
		Options:              options,
		BackupOperation:      backupOperation,
		FileSystemCheckOrder: fileSystemCheckOrder,
	}
}

func NewFstabLineFromConfig(c Config) *FstabLine {
	ent := NewFstabLineWithOptions(
		WithDevice(c.GetMountDevice()),
		WithMountPoint(c.GetMountPoint()),
		WithOptions(c.GenerateOptionString()),
		WithFileSystemType(c.GetFileSystemType()),
		WithBackupOperation(c.GetBackupOperation()),
		WithFileSystemCheckOrder(c.GetFileSystemCheckOrder()),
		)
	return ent
}

type FstabLineOption func(entry *FstabLine)

func NewFstabLineWithOptions(options ...FstabLineOption) *FstabLine {
	line := &FstabLine{}
	for _, opt := range options {
		opt(line)
	}
	return line
}

func WithDevice(device string) FstabLineOption {
	return func(line *FstabLine) {
		line.Device = device
	}
}
func WithMountPoint(mountpoint string) FstabLineOption {
	return func(line *FstabLine) {
		line.MountPoint = mountpoint
	}
}

func WithOptions(options string) FstabLineOption {
	return func(line *FstabLine) {
		line.Options = options
	}
}

func WithFileSystemType(fsType string) FstabLineOption {
	return func(line *FstabLine) {
		line.FileSystemType = fsType
	}
}

func WithFileSystemCheckOrder(fsco int) FstabLineOption {
	return func(line *FstabLine) {
		line.FileSystemCheckOrder = fsco
	}
}

func WithBackupOperation(bo int) FstabLineOption {
	return func(line *FstabLine) {
		line.BackupOperation = bo
	}
}
