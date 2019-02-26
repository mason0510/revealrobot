package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/robfig/cron"
)

type ActiveBetsTable struct {
	Rows []struct {
		Id        eos.Uint64      `json:"id"`
		Player    eos.Name        `json:"player"`
		Referer   eos.Name        `json:"referer"`
		BetAmount eos.Uint64      `json:"bet_number"`
		Asset     eos.Asset       `json:"bet_asset"`
		Seed      eos.Checksum256 `json:"seed"`
		Time      string          `json:"time"`
	} `json:"rows"`
	More bool `json:"more"`
}

type DiceRevealData struct {
	BetId     eos.Uint64    `json:"bet_id"`
	Signature ecc.Signature `json:"signature"`
}

type DiceRobot struct {
	name     string
	config   *ServerConfig
	services *Services
}

func (r *DiceRobot) run() {
	c := cron.New()
	spec := "*/1 * * * * ?"
	_ = c.AddFunc(spec, func() {
		body, err := getTableRows(r.config.node, r.name, "activebets")
		if err == nil {
			var list ActiveBetsTable
			err = json.Unmarshal(body, &list)

			if err == nil {
				for _, row := range list.Rows {
					fmt.Println("=======================骰子开奖 " + string(row.Id) + "=============================")
					r.pushAction(row.Id, row.Seed)
				}
			}
		}
	})
	c.Start()
}

func (r *DiceRobot) pushAction(betId eos.Uint64, seed eos.Checksum256) {
	keys, err := r.services.digestSigner.AvailableKeys()
	digest, err := hex.DecodeString(seed.String())
	sig, err := r.services.digestSigner.SignDigest(digest, keys[0])
	data := DiceRevealData{betId, sig}
	action := eos.Action{
		Account: eos.AccountName(r.name),
		Name:    "reveal",
		Authorization: []eos.PermissionLevel{
			{Actor: eos.AccountName(r.config.actorAccountName), Permission: eos.PN("active")}, //owner active
		},
		ActionData: eos.NewActionData(&data),
	}

	tx := eos.NewTransaction([]*eos.Action{&action}, &r.services.txOpts)
	signedTx, packedTx, err := r.services.api.SignTransaction(tx, r.services.txOpts.ChainID, eos.CompressionNone)
	if err == nil {
		_, err = json.MarshalIndent(signedTx, "", "")
		if err == nil {
			_, err = json.Marshal(packedTx)
			if err == nil {
				response, err := r.services.api.PushTransaction(packedTx)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(response)
				}
			}
		}
	}
}
