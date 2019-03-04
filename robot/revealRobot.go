package revealrobot

import (
	cf "../config"
	"../utils/stringhandler"
	"encoding/json"
	"fmt"
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"io/ioutil"
	"net/http"
	"time"
)

// ===========================     设置信息    ========================================
const IS_TEST = false

var (
	roundBasedGames [4]string
	diceGameName    string
	serverConfig    ServerConfig
)

type Timestamp struct {
	API  string   `json:"api"`
	V    string   `json:"v"`
	Ret  []string `json:"ret"`
	Data struct {
		T string `json:"t"`
	} `json:"data"`
}

func Init() {
	if IS_TEST {
		roundBasedGames = [4]string{"godappbaccar", "godappcbacca", "godapproulet", "godappredbla"}
		diceGameName = "godappdice12"
		serverConfig = ServerConfig{
			"https://api-kylin.eoslaomao.com",
			"5Kart9egqapRE6bXvSEr9sAaJecWxAvzZags9B831oab7TK29w7",
			"codemonkeyte",
			"5JqyADz2gRjct3E78pYHpxqEN6GbTWcv21GJ8SxExYZEmeQuoRA"}
	} else {
		serverConfig = ServerConfig{}
		serverConfig.node = cf.C.Node
		serverConfig.revealKey = cf.C.RevealKey
		serverConfig.actorAccountName = cf.C.ActorAccountName
		serverConfig.actorAccountKey = cf.C.ActorAccountKey

		roundBasedGames = [4]string{"baccarat.e", "dappbaccarat", "roulette.e", "warofstar.e"}
		diceGameName = "godice.e"
	}
}

func RevealRobot() {
	services := createServices(serverConfig)
	fmt.Println("RevealRobot", serverConfig)
	networkOffset := GetNetWorkOffset()
	for i := range roundBasedGames {
		robot := RoundBasedRobot{roundBasedGames[i], 0, networkOffset,
			RoundStatus{0, 0, 0, ""},
			&serverConfig,
			&services,
		}
		robot.run()
	}
	dice := DiceRobot{diceGameName, &serverConfig, &services}
	dice.run()
	select {}
}

func GetNetWorkOffset() int64 {
	timeResp, err := http.Get(cf.C.TimeUrl)
	if err != nil {
		fmt.Println("err", err)
	}
	s, err := ioutil.ReadAll(timeResp.Body)
	if err != nil {
		fmt.Println(err)
	}
	var data Timestamp
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal([]byte(s), &data)
	var value = (stringhandler.StringToInt(data.Data.T) / 1000) - time.Now().UTC().Unix()
	fmt.Println("===========网络延时===========", value)
	return value
}

func createServices(config ServerConfig) Services {
	digestSigner := *eos.NewKeyBag()
	_ = digestSigner.ImportPrivateKey(config.revealKey)

	api := eos.New(config.node)
	bag := eos.NewKeyBag()
	_ = bag.Add(config.actorAccountKey)
	key, _ := bag.AvailableKeys()
	api.SetCustomGetRequiredKeys(func(tx *eos.Transaction) (keys []ecc.PublicKey, e error) {
		return key, nil
	})
	api.SetSigner(bag)
	txOps, _ := getTxOps(api)
	return Services{*api, txOps, digestSigner, 0}
}

func getTxOps(api *eos.API) (eos.TxOptions, error) {
	opts := *&eos.TxOptions{}
	err := opts.FillFromChain(api)
	return opts, err
}
