package app

import (
	"os/user"
	"path/filepath"
	"io/ioutil"
	"encoding/json"
)

var configPath string
const defaultConfigName = "default"

type Config struct {
	Name string
	path string
	UiPath string
}

func init() {
	usr, err := user.Current()
	if err != nil {
		panic("can not get current user.")
	}
	configPath = filepath.Join(usr.HomeDir, ".shuffler", "config")
}

func (c *Config)Path() string {
	if c.path == "" {
		return filepath.Join(configPath, c.Name + ".json")
	}
	return c.path
}

func (c *Config)Load(name string) error {
	if name==""{
		name=defaultConfigName
	}
	c.Name=name
	data, err := ioutil.ReadFile(c.Path())
	if err != nil {
		return  err
	}
	if err := json.Unmarshal(data, c); err != nil {
		return err
	}
	return  nil
}

func (c *Config)Save() error {
	if c.Name==""{
		c.Name = defaultConfigName
	}
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(c.Path(), data, 777); err != nil {
		return err
	}
	return nil
}
