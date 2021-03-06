package main

import (
	"fmt"
	"git.code.oa.com/going/going/log"
	"git.code.oa.com/rainbow/golang-sdk/confapi"
	"git.code.oa.com/rainbow/golang-sdk/config"
	"git.code.oa.com/rainbow/golang-sdk/keep"
	"git.code.oa.com/rainbow/golang-sdk/types"
	"git.code.oa.com/rainbow/golang-sdk/watch"
)

type RainbowClient struct {
	ConnectStr      string
	AppID           string
	Groups          string //...string 代表什么？？
	EnvName         string
	UsingLocalCache bool
	UsingFileCache  bool
	client          *confapi.ConfAPI
}

func NewRainbowApi(appId, group string) (*RainbowClient, error) {
	rainbow := &RainbowClient{}
	rainbow.ConnectStr = "http://api.rainbow.oa.com:8080"
	rainbow.EnvName = "Default"
	var err error
	rainbow.client, err = confapi.New(
		types.ConnectStr(rainbow.ConnectStr),
		// types.ConnectStr("cl5://65026305:65536"), // 也可以使用cl5
		// types.ConnectStr("polaris://65026305:65536"), // 也可以使用polaris
		types.IsUsingLocalCache(true),
		types.IsUsingFileCache(true),

		// 预拉取这个appid, 环境(Env)下的 group
		types.AppID(appId),
		types.Groups(group),
	)
	return rainbow, err
}

func (rainbow *RainbowClient) SetGetOpts() []types.AssignGetOption {
	getOpts := make([]types.AssignGetOption, 0)
	getOpts = append(getOpts, types.WithAppID(rainbow.AppID))
	getOpts = append(getOpts, types.WithGroup(rainbow.Groups))
	getOpts = append(getOpts, types.WithEnvName(rainbow.EnvName))
	return getOpts
}

func (rainbow *RainbowClient) GetKey(key string, getOpts []types.AssignGetOption) (string, error) {
	val, err := rainbow.client.Get(key, getOpts...)
	if err != nil {
		log.Error("[rainbow.Get]%s\n", err.Error())
		return val, err
	}
	return val, nil
}

func (rainbow *RainbowClient) GetGroup(getOpts []types.AssignGetOption) (keep.Group, error) {
	// get group
	gval, err := rainbow.client.GetGroup(getOpts...)
	if err != nil {
		log.Error("[rainbow.Get]%s\n", err.Error())
		return gval, err
	}
	return gval, nil
}

// watchCallBack watch call back
func watchCallBack(oldVal watch.Result, newVal []*config.KeyValueItem) error {
	log.Debug("\n---------------------\n")
	log.Debug("rainbow old value:%+v\n", oldVal)
	log.Debug("rainbow new value:")
	for i := 0; i < len(newVal); i++ {
		log.Debug("%+v", *newVal[i])
	}
	log.Debug("\n---------------------\n")
	return nil
}

func (rainbow *RainbowClient) WatchGroup() {
	var watch = watch.Watcher{
		GetOptions: types.GetOptions{
			AppID:   rainbow.AppID,
			Group:   rainbow.Groups,
			EnvName: rainbow.EnvName,
		},
		CB: watchCallBack,
	}
	rainbow.client.AddWatcher(watch)
}

func main() {
	rainbow, err := NewRainbowApi("appid", "test")
	if err != nil {
		fmt.Println(err)
		return
	}
	get, err := rainbow.GetKey("story", rainbow.SetGetOpts())
	if err != nil {
		return
	}
	fmt.Println(get)
}
