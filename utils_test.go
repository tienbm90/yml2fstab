package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckFileSystemTypeValid(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		assert.True(t, CheckFileSystemTypeValid("nfs"))
		assert.True(t, CheckFileSystemTypeValid("ext4"))
	})

	t.Run("failure", func(t *testing.T) {
		assert.False(t, CheckFileSystemTypeValid("ffff"))
		assert.False(t, CheckFileSystemTypeValid("ext5"))
		assert.False(t, CheckFileSystemTypeValid("hdfs"))
	})
}

func TestCheckMountPointValid(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		assert.True(t, CheckMountPointValid("swap"))
		assert.True(t, CheckMountPointValid("/dev"))
		assert.True(t, CheckMountPointValid("/var"))
		assert.True(t, CheckMountPointValid("/home/hello"))
		assert.True(t, CheckMountPointValid("/home/backup"))
		assert.True(t, CheckMountPointValid("/dafadf545435"))
	})

	t.Run("invalid", func(t *testing.T) {
		assert.False(t, CheckMountPointValid("swap1"))
		assert.False(t, CheckMountPointValid("swing"))
		assert.False(t, CheckMountPointValid("faldshfjkha"))
	})
}

func TestCheckIPAddress(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		assert.True(t, CheckIPAddress("172.16.0.1"))
		assert.True(t, CheckIPAddress("127.0.0.1"))
		assert.True(t, CheckIPAddress("255.255.255.255"))
		assert.True(t, CheckIPAddress("0.0.0.0"))
		assert.True(t, CheckIPAddress("0.0.0.1"))
	})

	t.Run("invalid", func(t *testing.T) {
		assert.False(t, CheckIPAddress("127.0.0.a"))
		assert.False(t, CheckIPAddress("127.a.0.1"))
		assert.False(t, CheckIPAddress("127.%.0.1"))
		assert.False(t, CheckIPAddress("255.255.255.256"))
		assert.False(t, CheckIPAddress("255.255.255.-1"))
		assert.False(t, CheckIPAddress("255.255.255"))
		assert.False(t, CheckIPAddress("255.255"))
		assert.False(t, CheckIPAddress("255"))
		assert.False(t, CheckIPAddress("1"))
		assert.False(t, CheckIPAddress("-1.-1.-1.2"))
	})
}

func TestCheckHost(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		assert.True(t, CheckHost("international.com"))
		assert.True(t, CheckHost("international.io"))
		assert.True(t, CheckHost("www.international.io"))
		assert.True(t, CheckHost("www.9gag.io"))
		assert.True(t, CheckHost("9gag.io"))
		assert.True(t, CheckHost("localhost"))

	})

	t.Run("invalid", func(t *testing.T) {
		assert.False(t, CheckHost("9gag\\com"))
		assert.False(t, CheckHost("9gag_com"))
		assert.False(t, CheckHost("1"))
	})
}