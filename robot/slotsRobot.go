package revealrobot

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"strconv"
)

type SlotsActiveGameTable struct {
	Rows []struct {
		Id      eos.Uint64      `json:"id"`
		Player  eos.Name        `json:"player"`
		Referer eos.Name        `json:"referer"`
		Price   eos.Asset       `json:"price"`
		Result  eos.Uint64      `json:"result"`
		Seed    eos.Checksum256 `json:"seed"`
		Time    string          `json:"time"`
	} `json:"rows"`
	More bool `json:"more"`
}

type SlotsRevealData struct {
	BetId     eos.Uint64    `json:"game_id"`
	Signature ecc.Signature `json:"signature"`
}

type SlotsRobot struct {
	name     string
	config   *ServerConfig
	services *Services
}

func (r *SlotsRobot) run() {
	c := NewWithSecond()
	spec := "*/1 * * * * ?"
	_, _ = c.AddFunc(spec, func() {
		body, err := getTableRows(r.config.node, r.name, "activegame")
		fmt.Println("取得赌注数据: ", "body")
		if err == nil {
			var list SlotsActiveGameTable
			err = json.Unmarshal(body, &list)

			if err == nil {
				for _, row := range list.Rows {
					if row.Result == 65535 {
						fmt.Println("=======================老虎机开奖 " + strconv.Itoa(int(row.Id)) + "=============================")
						r.pushAction(row.Id, row.Seed)
					}
				}
			}
		}
	})
	c.Start()
}

func (r *SlotsRobot) pushAction(gameId eos.Uint64, seed eos.Checksum256) {
	keys, err := r.services.digestSigner.AvailableKeys(context.Background())
	digest, err := hex.DecodeString(seed.String())
	sig, err := r.services.digestSigner.SignDigest(digest, keys[0])
	data := SlotsRevealData{gameId, sig}
	action := eos.Action{
		Account: eos.AccountName(r.name),
		Name:    "reveal",
		Authorization: []eos.PermissionLevel{
			{Actor: eos.AccountName(r.config.actorAccountName), Permission: eos.PN("active")}, //owner active
		},
		ActionData: eos.NewActionData(&data),
	}

	tx := eos.NewTransaction([]*eos.Action{&action}, &r.services.txOpts)
	signedTx, packedTx, err := r.services.api.SignTransaction(context.Background(), tx, r.services.txOpts.ChainID, eos.CompressionNone)
	if err == nil {
		_, err = json.MarshalIndent(signedTx, "", "")
		if err == nil {
			_, err = json.Marshal(packedTx)
			if err == nil {
				response, err := r.services.api.PushTransaction(context.Background(), packedTx)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(response)
				}
			}
		}
	}
}

func (r *SlotsRobot) run1() {
	fmt.Println(r.config.node)
	body, err := getTableRows(r.config.node, r.name, "activegame")
	if err == nil {
		var list SlotsActiveGameTable
		err = json.Unmarshal(body, &list)

		if err == nil {
			for _, row := range list.Rows {
				if row.Result == 65535 {
					fmt.Println("=======================老虎机开奖 " + strconv.Itoa(int(row.Id)) + "=============================")
					r.pushAction(row.Id, row.Seed)
				}
			}
		}
	}
}
