package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildStringFromSlice(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		arr := []string{"noexec", "nosuid", "mode=50"}

		assert.Equal(t, "noexec,nosuid,mode=50", BuildStringFromSlice(arr))

	})
}

func TestNewConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		data := []struct {
			Mount string
			Mtype string
		}{
			{
				Mount: "/boot",
				Mtype: "/xfs",
			},
			{
				Mount: "/",
				Mtype: "/ext4",
			},
		}
		for _, d := range data {
			cnf := NewConfigWithOptions(
				WithConfigMount(d.Mount),
				WithConfigFSType(d.Mtype),
			)

			assert.Equal(t, cnf.GetMountPoint(), d.Mount)
			assert.Equal(t, cnf.GetFileSystemType(), d.Mtype)
		}
	})

	t.Run("failure", func(t *testing.T) {
		data := []struct {
			Mount string
			Mtype string
		}{
			{
				Mount: "/boot",
				Mtype: "/xfs",
			},
			{
				Mount: "/",
				Mtype: "/ext4",
			},
		}
		for _, d := range data {
			cnf := NewConfigWithOptions(
				WithConfigMount(d.Mtype),
				WithConfigFSType(d.Mount),
			)
			assert.NotEqual(t, cnf.GetMountPoint(), d.Mount)
			assert.NotEqual(t, cnf.GetFileSystemType(), d.Mtype)
		}
	})
}

func TestConfig_AddOptions(t *testing.T) {
	data := []struct {
		Mount string
		Mtype string
	}{
		{
			Mount: "/boot",
			Mtype: "/xfs",
		},
		{
			Mount: "/",
			Mtype: "/ext4",
		},
	}

	t.Run("add one option success", func(t *testing.T) {
		for _, d := range data {
			cnf := NewConfigWithOptions(
				WithConfigMount(d.Mount),
				WithConfigFSType(d.Mtype),
			)

			assert.Equal(t, cnf.GetMountPoint(), d.Mount)
			assert.Equal(t, cnf.GetFileSystemType(), d.Mtype)

			opt := "noexec"
			cnf.AddOption(opt)

			assert.Equal(t, len(cnf.GetOptions()), 1)
			assert.Equal(t, cnf.GetOptions()[0], opt)
		}
	})

	t.Run("add two options success", func(t *testing.T) {
		for _, d := range data {
			cnf := NewConfigWithOptions(
				WithConfigMount(d.Mount),
				WithConfigFSType(d.Mtype),
			)

			assert.Equal(t, cnf.GetMountPoint(), d.Mount)
			assert.Equal(t, cnf.GetFileSystemType(), d.Mtype)

			opts := []string{"noexec", "nosuid"}

			cnf.AddOptions(opts)

			assert.Equal(t, len(cnf.GetOptions()), 2)
			assert.Equal(t, cnf.GetOptions()[0], opts[0])
			assert.Equal(t, cnf.GetOptions()[1], opts[1])
		}
	})
}

func TestConfig_GetMountDevice(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		assert.Equal(t, GetMountDevice("/dev/sda", ""), "/dev/sda")
		assert.Equal(t, GetMountDevice("/dev/sda", "/var/home"), "/dev/sda")
		assert.Equal(t, GetMountDevice("1", "/var/home"), "1")
		assert.Equal(t, GetMountDevice("test.com", "/var/home"), "test.com:/var/home")
		assert.Equal(t, GetMountDevice("9test.com", "/var/home"), "9test.com:/var/home")
		assert.Equal(t, GetMountDevice("10.20.100.2", "/var/home"), "10.20.100.2:/var/home")
		assert.Equal(t, GetMountDevice("127.0.0.1", "/var/home"), "127.0.0.1:/var/home")
	})
}

func TestNewConfigFromMap(t *testing.T) {
	t.Run("success ", func(t *testing.T) {
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
		assert.NoError(t, err)
		assert.Equal(t, cnf.GetMountPoint(), "/home")
		assert.Equal(t, cnf.GetFileSystemType(), "nfs")
		assert.Equal(t, cnf.GetMountDevice(), "192.168.4.6:/var/nfs/home")
		assert.Equal(t, len(cnf.GetOptions()), 2)
	})
}

func TestNewConfigs(t *testing.T) {
	t.Run("success ", func(t *testing.T) {

		mm := make(map[string]interface{})

		//init first kv
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

		mm["192.168.4.6"] = data

		//init second kv
		data2 := make(map[string]interface{})
		data2["mount"] = "/home"
		data2["type"] = "xfs"

		mm["/dev/sda1"] = data2
		cnfs, err := NewConfigs(mm)
		assert.NoError(t, err)

		for _, v := range cnfs {
			if v.Source == "/dev/sda1" {
				assert.Equal(t, v.GetMountPoint(), data2["mount"])
				assert.Equal(t, v.GetFileSystemType(), data2["type"])
				assert.Equal(t, v.GetMountDevice(), "/dev/sda1")
			} else {
				assert.Equal(t, v.GetMountPoint(), "/home")
				assert.Equal(t, v.GetFileSystemType(), "nfs")
				assert.Equal(t, v.GetMountDevice(), "192.168.4.6:/var/nfs/home")
				assert.Equal(t, len(v.GetOptions()), 2)
			}

		}
	})
}
