package main

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Config struct {
	Source               string
	Mount                string
	Type                 string
	Export               string
	Options              []string
	BackupOperation      int
	FileSystemCheckOrder int
}

func (c *Config) SetBackupOperation(s int) {
	c.BackupOperation = s

}

func (c *Config) SetFileSystemCheckOrder(s int) {
	c.BackupOperation = s

}

func (c *Config) SetSource(s string) {
	c.Source = s

}

func (c *Config) SetMount(m string) {
	c.Mount = m
}

func (c *Config) SetType(t string) {
	c.Type = t
}

func (c *Config) SetExport(export string) {
	c.Export = export
}

func (c *Config) AddOption(opt string) {
	if c.Options == nil {
		c.Options = make([]string, 0)
	}
	c.Options = append(c.Options, opt)
}

func (c *Config) AddOptions(opts []string) {
	if c.Options == nil {
		c.Options = make([]string, 0)
	}
	for _, v := range opts {
		c.Options = append(c.Options, v)
	}
}

func (c *Config) GetOptions() []string {
	return c.Options
}

func (c *Config) GenerateOptionString() string {
	if len(c.Options) == 0 {
		return "defaults"
	}
	return BuildStringFromSlice(c.Options)
}

func BuildStringFromSlice(arr []string) string {
	if len(arr) == 0 {
		return ""
	}
	opt := strings.Join(arr, ",")
	return opt
}

func (c *Config) GetMountPoint() string {
	return c.Mount
}

func (c *Config) GetBackupOperation() int {
	return c.BackupOperation
}

func (c *Config) GetFileSystemCheckOrder() int {
	return c.FileSystemCheckOrder
}

func (c *Config) GetFileSystemType() string {
	return c.Type
}

func (c *Config) GetMountDevice() string {
	return GetMountDevice(c.Source, c.Export)
}

func (c *Config) IsMountValid() bool {
	return c.Mount != ""
}

func (c *Config) IsFileSystemTypeValid() bool {
	return c.Type != ""
}

func (c *Config) IsValid() bool {
	return c.IsFileSystemTypeValid() && c.IsMountValid()

}

func GetMountDevice(source string, export string) string {
	var mountDevice string
	if CheckHost(source) {
		mountDevice = fmt.Sprintf("%s:%s", source, export)
	} else if CheckIPAddress(source) {
		mountDevice = fmt.Sprintf("%s:%s", source, export)
	} else {
		mountDevice = source
	}

	return mountDevice
}

func NewConfigFromMapData(source string, m map[string]interface{}) (*Config, error) {
	//parse mount field
	if m["mount"] == nil {
		return nil, errors.New("mount point not found")
	}
	mount, ok := m["mount"].(string)
	if !ok {
		return nil, errors.New("invalid format for mount field. Require string")
	}
	//parse export field
	export := ""
	if m["export"] != nil {
		e := m["export"]
		exp, ok := e.(string)
		if !ok {
			return nil, errors.New("invalid format for export field. Require string")
		} else {
			export = exp
		}
	}

	//parse type field
	if m["type"] == nil {
		return nil, errors.New("mount point not found")
	}
	fsType, ok := m["type"].(string)
	if !ok {
		return nil, errors.New("invalid format for type field. Require string")
	}

	//parse options
	options := make([]string, 0)
	if m["options"] != nil {
		op := m["options"]
		switch reflect.TypeOf(op).Kind() {
		case reflect.Slice:
			opts, ok := op.([]interface{})
			if !ok {
				return nil, errors.New("invalid format for options field. Require list of strings")
			}
			for _, v := range opts {
				opt, ok := v.(string)
				if !ok {
					return nil, errors.New("invalid format for options field. Require string")
				}
				options = append(options, opt)
			}
		default:
			return nil, errors.New("invalid format for options field. Require list of strings")
		}
	}

	// create new configuration
	conf := NewConfigWithOptions(
		WithConfigSource(source),
		WithConfigMount(mount),
		WithConfigFSType(fsType),
		WithConfigExport(export),
		WithConfigBackupOperation(0),
		WithConfigFileSystemCheckOrder(0),
		WithConfigOptions(options),
	)
	return conf, nil
}

func NewConfigs(mm map[string]interface{}) ([]*Config, error) {
	var configs []*Config
	for k, v := range mm {
		r, ok := v.(map[string]interface{})
		if !ok {
			return nil, errors.New("parse data error. Cannot cast value of mount point")
		}
		conf, err := NewConfigFromMapData(k, r)
		if err != nil {
			return nil, err
		}
		configs = append(configs, conf)
	}
	return configs, nil
}

type ConfigOption func(*Config)

func NewConfigWithOptions(options ...ConfigOption) *Config {
	config := &Config{}
	for _, opt := range options {
		opt(config)
	}
	return config
}

func WithConfigSource(source string) ConfigOption {
	return func(config *Config) {
		config.Source = source
	}
}
func WithConfigBackupOperation(bo int) ConfigOption {
	return func(config *Config) {
		config.BackupOperation = bo
	}
}
func WithConfigFileSystemCheckOrder(fsCheckOrder int) ConfigOption {

	return func(config *Config) {
		config.FileSystemCheckOrder = fsCheckOrder
	}
}

func WithConfigExport(export string) ConfigOption {
	return func(config *Config) {
		config.Export = export
	}
}

func WithConfigMount(mount string) ConfigOption {
	return func(config *Config) {
		config.Mount = mount
	}
}
func WithConfigFSType(fsType string) ConfigOption {
	return func(config *Config) {
		config.Type = fsType
	}
}

func WithConfigOptions(options []string) ConfigOption {
	return func(config *Config) {
		config.Options = options
	}
}
