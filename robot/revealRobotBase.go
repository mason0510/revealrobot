package revealrobot

import (
	"fmt"
	"github.com/eoscanada/eos-go"
	"io/ioutil"
	"net/http"
	"strings"
)

type ServerConfig struct {
	node      string
	revealKey string
	actorAccountName string
	actorAccountKey  string
}

type Services struct {
	api          eos.API
	txOpts       eos.TxOptions
	digestSigner eos.KeyBag
	lastRefresh  int64
}

func getTableRows(node string, game string, table string) ([]byte, error) {
	url := node + "/v1/chain/get_table_rows"
	payload := strings.NewReader("{\n  \"scope\": \"" + game + "\",\n  \"code\": \"" + game + "\",\n  " +
		"\"table\": \"" + table + "\",\n  \"json\": \"true\",\n  \"limit\": 1000\n}")
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("Postman-Token", "602072f4-443d-4cc3-b163-695c3283cb55")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(res.Body)
}

func (s *Services) refresh(currentTime int64) {
	if currentTime-s.lastRefresh > 60 {
		s.lastRefresh = currentTime
		fmt.Println("==========================================更新网络配置==================================================")
		opts, err := getTxOps(&s.api)
		if err == nil {
			s.txOpts = opts
		}
	}
}
