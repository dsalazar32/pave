package provider

import (
	"encoding/json"
	"fmt"
	"github.com/go-chef/chef"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"time"
)

type Chef struct {
	*chef.Client
	filePath string
}

type DBags struct {
	Items []chef.DataBagItem
}

func (p Chef) Read() (string, error) {
	bag, item, err := parseCloudStorageUrl(p.filePath)
	if err != nil {
		return "", err
	}

	var lr *chef.DataBagListResult
	var cErr error

	items := []chef.DataBagItem{}

	if item == "" {
		lr, cErr = p.Client.DataBags.ListItems(bag)
		if cErr != nil {
			return "", err
		}
		fmt.Println(lr)
	}

	// Grabs all of the items in a databag and collects it
	if item == "@all" {
		lr, cErr = p.Client.DataBags.ListItems(bag)
		if cErr != nil {
			return "", err
		}

		for k := range *lr {
			i, err := p.Client.DataBags.GetItem(bag, k)
			if err != nil {
				return "", err
			}
			items = append(items, i)
		}
		// Grabs the items in a databag and collects it
	} else {
		lr, cErr = p.Client.DataBags.ListItems(bag)
		if cErr != nil {
			return "", err
		}

		i, err := p.Client.DataBags.GetItem(bag, item)
		if err != nil {
			return "", err
		}
		items = append(items, i)
	}

	var itemS string
	if len(items) > 0 {
		j, err := json.Marshal(DBags{items})
		if err != nil {
			return "", err
		}
		itemS = string(j)
	}

	return itemS, nil
}

func (p Chef) Write(s string) error {
	panic("implement me")
}

func init() {
	Constructors[CHF] = &ProviderSpec{
		New:         NewChef,
		description: "Chef provider is used to interact with chefs data bags",
	}
}

// TODO: setup the proper configuration file
func NewChef(infile string) (Provider, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	client, err := chef.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &Chef{
		Client:   client,
		filePath: infile,
	}, nil
}

func LoadConfig() (*chef.Config, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	cfg := make(map[interface{}]interface{})
	b, err := ioutil.ReadFile(filepath.Join(usr.HomeDir, ".pave", "config.yml"))
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}

	t := func(i interface{}) map[interface{}]interface{} {
		return i.(map[interface{}]interface{})
	}

	c := t(cfg)["chef"]
	p, err := homedir.Expand(t(c)["key"].(string))
	if err != nil {
		return nil, err
	}

	k, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	return &chef.Config{
		Name:    t(c)["name"].(string),
		BaseURL: t(c)["base_url"].(string),
		SkipSSL: t(c)["skip_ssl"].(bool),
		Timeout: time.Duration(t(c)["timeout"].(int)) * time.Second,
		Key:     string(k),
	}, nil
}
