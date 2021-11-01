# YML2FSTAB
## Description:
Suppose we have a YAML file describing the mount points of a system as follows.
```yaml
---
fstab:
  /dev/sda1:
    mount: /boot
    type: xfs
  /dev/sda2:
    mount: /
    type: ext4
  /dev/sdb1:
    mount: /var/lib/postgresql
    type: ext4
  192.168.4.5:
    mount: /home
    export: /var/nfs/home
    type: nfs
    options:
      - noexec
      - nosuid
```
#### Output
- Please create a utility to process that YAML and output an /etc/fstab file based on that yaml.
Please provide your code along with a README describing your ideas/approachs and the use.
You are free to use languages and environments regarding a system of your choice but the tasks
must run on a Linux system.
- Note: You should use a programming/scripting language to achieve this task, and not Ansible
or such.

## Solution

- Check the input file exists.
- Check input file has valid( yml format and not empty)
- Read file and create fstab configurations
- Write configurations into new temporary file
- Override /etc/fstab with new configurations

## Build and run
### Prerequisite
```text
Go sdk is installed on your system.
```
### Build:
```shell
go build -o yml2fstab
```
### Run:
#### Supported filesystems
- "ext"
- "ext2"
- "ext3"
- "ext4"
- "jfs"
- "reiserfs"
- "xfs"
- "btrfs"
- "swap"
- "iso9660"
- "nfs"
- "nfs4"
- "udf"
- "vfat"
- "devpts"


#### Note: Require sudo if you want to update /etc/fstab
```shell
./yml2fstab -in input.yml -out /tmp/fstab
```

## Test:
```shell
go test .

```

## Parameter:
```shell
in : Path to yml file. Default is input.yml
out : Path to output file. Default is /etc/fstab
tmp-file: Path to temp file. Default is /tmp/fstab.temp
```
## Third party lib:
- "gopkg.in/yaml.v3"
