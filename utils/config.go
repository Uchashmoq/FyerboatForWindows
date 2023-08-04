package utils

import (
	"encoding/json"
	"errors"
	"os"
)

const (
	CONFIG_PATH = "./config.json"
)

type Node struct {
	Name       string
	ServerAddr string
	StaticKey  string
}
type Config struct {
	Nodes      map[string]*Node
	ClientAddr string
}

func (c *Config) Store(path string) error {
	fd, err := os.Create(path)
	if err != nil {
		return err
	}
	bytes, err1 := json.MarshalIndent(c, "", " ")
	if err1 != nil {
		return err1
	}
	_, err2 := fd.Write(bytes)
	if err2 != nil {
		return err2
	}
	fd.Close()
	return nil
}
func (c *Config) Load(path string) error {
	fd, err := os.Open(path)
	if err != nil {
		return err
	}
	b := make([]byte, 1024*1024*2)
	n, err1 := fd.Read(b)
	if err1 != nil {
		return err1
	}
	if n == 0 {
		return errors.New("file is empty")
	}
	json.Unmarshal(b[:n], c)
	fd.Close()
	return nil
}
