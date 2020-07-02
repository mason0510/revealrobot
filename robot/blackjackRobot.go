package revealrobot

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
)

type BlackJackActionTable struct {
	Rows []struct {
		Id     eos.Uint64      `json:"game_id"`
		Action eos.Uint64      `json:"action"`
		Seed   eos.Checksum256 `json:"seed"`
	} `json:"rows"`
	More bool `json:"more"`
}

type BlackjackRevealData struct {
	BetId     eos.Uint64    `json:"id"`
	Signature ecc.Signature `json:"signature"`
}

type BlackjackRobot struct {
	name     string
	config   *ServerConfig
	services *Services
}

func (r *BlackjackRobot) run() {
	c := NewWithSecond()
	spec := "*/1 * * * * ?"
	_, _ = c.AddFunc(spec, func() {
		body, err := getTableRows(r.config.node, r.name, "actions")
		if err == nil {
			var list BlackJackActionTable
			err = json.Unmarshal(body, &list)

			if err == nil {
				for _, row := range list.Rows {
					fmt.Println("=======================21点发牌 " + string(row.Id) + "=============================")
					r.pushAction(row.Id, row.Seed)
				}
			}
		}
	})
	c.Start()
}

func (r *BlackjackRobot) pushAction(id eos.Uint64, seed eos.Checksum256) {
	keys, err := r.services.digestSigner.AvailableKeys(context.Background())
	digest, err := hex.DecodeString(seed.String())
	sig, err := r.services.digestSigner.SignDigest(digest, keys[0])
	data := BlackjackRevealData{id, sig}
	action := eos.Action{
		Account: eos.AccountName(r.name),
		Name:    "resolve",
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
