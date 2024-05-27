package freeswitch

import (
	"contactcenter-api/common/http"
	"contactcenter-api/common/log"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/luandnh/goesl"
)

type FreeswitchESLConfig struct {
	Address  string
	Port     int
	Password string
	Timeout  int
}

type FreeswitchESL struct {
	Client  *goesl.Client
	configs []FreeswitchESLConfig
}

type ESLResponse struct {
	Command  string         `json:"command"`
	Format   string         `json:"format"`
	Data     map[string]any `json:"data"`
	Status   string         `json:"status"`
	Response any            `json:"response"`
}

var ESLClient FreeswitchESL

func (e *FreeswitchESL) Connect() (*goesl.Client, error) {
	for _, config := range e.configs {
		client, err := goesl.NewClient(config.Address, config.Port, config.Password, config.Timeout)
		if err != nil {
			log.Error(err)
		}
		return client, err
	}
	return nil, errors.New("connect esl failed")
}

func (e *FreeswitchESL) ConnectToHost(host string) (*goesl.Client, error) {
	for _, config := range e.configs {
		if host == config.Address {
			client, err := goesl.NewClient(config.Address, config.Port, config.Password, config.Timeout)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			return client, err
		}
	}
	return nil, errors.New("connect " + host + " failed")
}

func NewFreeswitchESL(configs ...FreeswitchESLConfig) FreeswitchESL {
	return FreeswitchESL{
		configs: configs,
	}
}

func (e *FreeswitchESL) GetCurrentHostOfChannel(uuid string) (*goesl.Client, string, error) {
	cmd := "uuid_exists " + uuid
	for _, config := range e.configs {
		client, err := goesl.NewClient(config.Address, config.Port, config.Password, config.Timeout)
		if err == nil {
			if rawResponse, err := client.Api(cmd); err != nil {
				log.Error(err)
				return nil, "", err
			} else if string(rawResponse.Body) == "true" {
				return client, config.Address, nil
			} else {
				client.Close()
			}
		} else {
			log.Error(err)
		}
	}
	return nil, "", nil
}

func (e *FreeswitchESL) SendToAllHost(cmd string) error {
	var wg sync.WaitGroup
	for _, c := range e.configs {
		wg.Add(1)
		go func(config FreeswitchESLConfig) {
			defer wg.Done()
			client, err := goesl.NewClient(config.Address, config.Port, config.Password, config.Timeout)
			if err != nil {
				log.Error(err)
				return
			}
			defer client.Close()
			if rawResponse, err := client.Api(cmd); err != nil {
				log.Errorf("SendToAllHost : %s %s %s", config.Address, cmd, err.Error())
				return
			} else {
				log.Infof("SendToAllHost : %s %s", cmd, string(rawResponse.Body))
			}
		}(c)
	}
	wg.Wait()
	return nil
}

func ExecJsonAPI(command, arguments string, parameters map[string]any) (any, error) {
	jsonCmd := map[string]any{
		"command": command,
		"data": map[string]any{
			"arguments": arguments,
		},
	}
	jsonCmdData := jsonCmd["data"].(map[string]any)
	if len(parameters) > 0 {
		for key, value := range parameters {
			jsonCmdData[key] = value
		}
	}
	cmdBytes, err := json.Marshal(jsonCmd)
	if err != nil {
		return nil, err
	}

	// r := regexp.MustCompile("\u005C")
	// cmd := r.ReplaceAllString(string(cmdBytes), "")
	cmd := "json " + string(cmdBytes)
	log.Debug(cmd)
	client, err := ESLClient.Connect()
	if err != nil {
		return nil, err
	}
	defer client.Close()
	rawResponse, err := client.Api(cmd)
	if err != nil {
		return nil, err
	}
	result := make(map[string]any)
	err = json.Unmarshal(rawResponse.Body, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ExecJsonBgAPI(command, arguments string) error {
	cmd := map[string]any{
		"command": command,
		"data": map[string]any{
			"arguments": arguments,
		},
	}
	cmdBytes, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	log.Debug(string(cmdBytes))
	client, err := ESLClient.Connect()
	if err != nil {
		return err
	}
	defer client.Close()
	err = client.BgApi(string(cmdBytes))
	return err
}

func ExecAPI(command string) (string, error) {
	client, err := ESLClient.Connect()
	if err != nil {
		return "", err
	}
	defer client.Close()
	rawResponse, err := client.Api(command)
	if err != nil {
		return "", err
	}
	return string(rawResponse.Body), nil
}

func ParseAPIResponse(data any) (ESLResponse, error) {
	tmp, err := json.Marshal(data)
	if err != nil {
		return ESLResponse{}, err
	}
	var result ESLResponse
	err = json.Unmarshal(tmp, &result)
	if err != nil {
		return ESLResponse{}, err
	}
	return result, nil
}

func (e *FreeswitchESL) ClearCallCenter() {
	var wg sync.WaitGroup
	wg.Add(len(e.configs))
	for _, cfg := range e.configs {
		go func(config FreeswitchESLConfig) {
			defer wg.Done()
			client, err := goesl.NewClient(config.Address, config.Port, config.Password, config.Timeout)
			if err != nil {
				log.Error(err)
				return
			}
			defer client.Close()
			keys := []string{"configuration.callcenter.conf"}
			if res, err := client.Api("switchname"); err != nil {
				log.Error(err)
			} else if len(string(res.Body)) > 0 {
				keys = append(keys, "configuration.callcenter.conf."+string(res.Body))
			}
			for _, key := range keys {
				customEvent := []string{
					"CUSTOM",
					"Event-Name: CUSTOM",
					"Event-Subclass: fusion::file",
					"API-Command: cache",
					"API-Command-Argument: delete " + key,
				}
				if _, err := client.SendEvent(customEvent); err != nil {
					log.Error(err)
				}
				key := "rm%20%2Fvar%2Fcache%2Ffusionpbx%2F" + key
				urlString := fmt.Sprintf("http://%s:%d/webapi/system?%s", config.Address, 7080, key)
				if res, err := http.GetWithBasicAuth(urlString, "freeswitch", config.Password); err != nil {
					log.Error(err)
				} else {
					log.Info(res.Status())
				}
			}
			res, err := client.Api("reloadxml")
			if err != nil {
				log.Error(err)
			} else {
				log.Infof("reload xml : %s", string(res.Body))
			}
		}(cfg)
	}
	wg.Wait()
}

func (e *FreeswitchESL) ClearExtensionDirectory(extension, context string) {
	for _, cfg := range e.configs {
		go func(config FreeswitchESLConfig) {
			key := "rm%20%2Fvar%2Fcache%2Ffusionpbx%2F" + fmt.Sprintf("directory.%s@%s", extension, context)
			urlString := fmt.Sprintf("http://%s:%d/webapi/system?%s", config.Address, 7080, key)
			if res, err := http.GetWithBasicAuth(urlString, "freeswitch", config.Password); err != nil {
				log.Error(err)
			} else {
				log.Info(res.Status())
			}
		}(cfg)
	}
}

func (e *FreeswitchESL) ClearCacheOfDomain(context string) {
	for _, cfg := range e.configs {
		go func(config FreeswitchESLConfig) {
			key := "rm%20%2Fvar%2Fcache%2Ffusionpbx%2F" + fmt.Sprintf("directory*@%s", context)
			urlString := fmt.Sprintf("http://%s:%d/webapi/system?%s", config.Address, 7080, key)
			if res, err := http.GetWithBasicAuth(urlString, "freeswitch", config.Password); err != nil {
				log.Error(err)
			} else {
				log.Info(res.Status())
			}
		}(cfg)
	}
}

func (e *FreeswitchESL) ClearCacheDialplan(context string) {
	for _, cfg := range e.configs {
		go func(config FreeswitchESLConfig) {
			key := "rm%20%2Fvar%2Fcache%2Ffusionpbx%2F" + fmt.Sprintf("dialplan.%s", context)
			urlString := fmt.Sprintf("http://%s:%d/webapi/system?%s", config.Address, 7080, key)
			if res, err := http.GetWithBasicAuth(urlString, "freeswitch", config.Password); err != nil {
				log.Error(err)
			} else {
				log.Info(res.Status())
			}
		}(cfg)
	}
}
