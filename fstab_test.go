package main

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestNewFstabEntry(t *testing.T) {
	entry := []struct {
		Device               string
		MountPoint           string
		FileSystemType       string
		Options              string
		BackupOperation      int
		FileSystemCheckOrder int
	}{
		{
			Device:               "/dev/mapper/rhel-root",
			MountPoint:           "/",
			FileSystemType:       "xfs",
			Options:              "defaults",
			BackupOperation:      0,
			FileSystemCheckOrder: 0,
		},
		{
			Device:               "/dev/sda1",
			MountPoint:           "swap",
			FileSystemType:       "swap",
			Options:              "defaults",
			BackupOperation:      0,
			FileSystemCheckOrder: 0,
		}, {
			Device:               "devpts",
			MountPoint:           "/dev/pts",
			FileSystemType:       "devpts",
			Options:              "mode=0620,gid=5",
			BackupOperation:      0,
			FileSystemCheckOrder: 0,
		},
	}
	t.Run("Success create", func(t *testing.T) {

		for _, ent := range entry {
			fsEntry := NewFstabEntry(ent.Device, ent.MountPoint, ent.FileSystemType, ent.Options, ent.BackupOperation, ent.FileSystemCheckOrder)

			assert.Equal(t, fsEntry.Device, ent.Device)
			assert.Equal(t, fsEntry.MountPoint, ent.MountPoint)
			assert.Equal(t, fsEntry.FileSystemType, ent.FileSystemType)
			assert.Equal(t, fsEntry.Options, ent.Options)
			assert.Equal(t, fsEntry.BackupOperation, ent.BackupOperation)
			assert.Equal(t, fsEntry.FileSystemCheckOrder, fsEntry.FileSystemCheckOrder)
		}
	})

	t.Run("Entry valid", func(t *testing.T) {
		for _, ent := range entry {
			fsEntry := NewFstabEntry(ent.Device, ent.MountPoint, ent.FileSystemType, ent.Options, ent.BackupOperation, ent.FileSystemCheckOrder)

			assert.True(t, fsEntry.IsMountPointValid())
			assert.True(t, fsEntry.IsFileSystemTypeValid())
			assert.True(t, fsEntry.IsBackupOperationValid())
			assert.True(t, fsEntry.IsFileSystemCheckOrderValid())
			assert.True(t, fsEntry.IsValid())

		}
	})

	t.Run("Entry invalid", func(t *testing.T) {
		invalidEntries := []struct {
			Device               string
			MountPoint           string
			FileSystemType       string
			Options              string
			BackupOperation      int
			FileSystemCheckOrder int
		}{
			{
				Device:               "/dev/sda1",
				MountPoint:           "swing",
				FileSystemType:       "swap1",
				Options:              "defaults",
				BackupOperation:      -1,
				FileSystemCheckOrder: 4,
			},
		}

		for _, ent := range invalidEntries {
			fsEntry := NewFstabEntry(ent.Device, ent.MountPoint, ent.FileSystemType, ent.Options, ent.BackupOperation, ent.FileSystemCheckOrder)
			log.Println(fsEntry.MountPoint, fsEntry.FileSystemType)
			assert.False(t, fsEntry.IsMountPointValid())
			assert.False(t, fsEntry.IsFileSystemTypeValid())
			assert.False(t, fsEntry.IsBackupOperationValid())
			assert.False(t, fsEntry.IsFileSystemCheckOrderValid())
			assert.False(t, fsEntry.IsValid())
		}
	})
}

func TestNewFstabEntryFromConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		data := make(map[string]interface{})
		data["mount"] = "/home"
		data["export"] = "/var/nfs/home"
		data["type"] = "nfs"
		opts := []string{"noexec", "nosuid"}
		opti := make([]interface{}, len(opts))
		for i := 0; i < len(opts); i++ {
			opti[i] = opts[i]
		}
		data["options"] = opti
		cnf, err := NewConfigFromMapData("192.168.4.6", data)
		assert.NoError(t, err, nil)

		ent := NewFstabLineFromConfig(*cnf)
		assert.Equal(t, ent.FileSystemCheckOrder, 0)
		assert.Equal(t, ent.FileSystemType, cnf.GetFileSystemType())
		assert.Equal(t, ent.MountPoint, cnf.GetMountPoint())
		assert.Equal(t, ent.Device, cnf.GetMountDevice())
		assert.Equal(t, ent.Options, cnf.GenerateOptionString())
		assert.Equal(t, ent.GenerateFstabEntryString(), "192.168.4.6:/var/nfs/home /home nfs noexec,nosuid 0 0")
	})

	t.Run("error. Require string for options field", func(t *testing.T) {
		data := make(map[string]interface{})
		data["mount"] = "/home"
		data["export"] = "/var/nfs/home"
		data["type"] = "nfs"
		data["options"] = "noexec"
		_, err := NewConfigFromMapData("192.168.4.6", data)
		assert.Error(t, err, nil)
	})
}
