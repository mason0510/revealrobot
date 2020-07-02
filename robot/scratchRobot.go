package revealrobot

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/robfig/cron"
	"strconv"
)

type ActiveCardTable struct {
	Rows []struct {
		Id       eos.Uint64      `json:"id"`
		Player   eos.Name        `json:"player"`
		Referer  eos.Name        `json:"referer"`
		CardType eos.Uint64      `json:"card_type"`
		Price    eos.Asset       `json:"price"`
		Reward   eos.Asset       `json:"reward"`
		Result   eos.Uint64      `json:"result"`
		Seed     eos.Checksum256 `json:"seed"`
		Time     string          `json:"time"`
	} `json:"rows"`
	More bool `json:"more"`
}

type ScratchRevealData struct {
	BetId     eos.Uint64    `json:"card_id"`
	Signature ecc.Signature `json:"signature"`
}

type ScratchRobot struct {
	name     string
	config   *ServerConfig
	services *Services
}

func (r *ScratchRobot) run() {
	c := cron.New()
	spec := "*/1 * * * * ?"
	c.AddFunc(spec, func() {
		body, err := getTableRows(r.config.node, r.name, "activecard")
		if err == nil {
			var list ActiveCardTable
			err = json.Unmarshal(body, &list)

			if err == nil {
				for _, row := range list.Rows {
					if row.Result == 0 {
						fmt.Println("=======================刮刮乐开奖 " + strconv.Itoa(int(row.Id)) + "=============================")
						r.pushAction(row.Id, row.Seed)
					}
				}
			}
		}
	})
	c.Start()
}

func (r *ScratchRobot) pushAction(cardId eos.Uint64, seed eos.Checksum256) {
	keys, err := r.services.digestSigner.AvailableKeys(context.Background())
	digest, err := hex.DecodeString(seed.String())
	sig, err := r.services.digestSigner.SignDigest(digest, keys[0])
	data := ScratchRevealData{cardId, sig}
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
