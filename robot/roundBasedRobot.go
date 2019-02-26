package revealrobot

import (
	"revealrobot/utils/bet"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/robfig/cron"
	"time"
)

type RoundStatus struct {
	roundId int
	status  int
	endTime int64
	seed    string
}

type RevealData struct {
	GameId    eos.Uint64    `json:"game_id"`
	Signature ecc.Signature `json:"signature"`
}

type RoundBasedRobot struct {
	name          string
	count         int
	networkOffset int64
	status        RoundStatus

	config   *ServerConfig
	services *Services
}

func (r *RoundBasedRobot) run() {
	c := cron.New()
	spec := "*/1 * * * * ?"
	err := c.AddFunc(spec, func() {
		var currentTime = time.Now().UTC().Unix() + r.networkOffset

		if currentTime >= r.status.endTime {
			r.status = r.getStatus()
		}
		//下注
		if r.status.status == 1 {
			//获取下注时间
			return
		}

		if r.status.status == 2 {
			//获取下注时间
			r.bettime(currentTime)
		}

		r.services.refresh(currentTime)
	})
	if err != nil {
		fmt.Println(err)
	}
	c.Start()
}

func (r *RoundBasedRobot) bettime(currentTime int64) {
	r.count++
	var openTime = r.status.endTime - currentTime

	//网路赌注
	if openTime != 0 && openTime%10 == 0 {
		fmt.Println("===== ", r.name, " 本轮游戏结束剩余 ", openTime, "，网络时间:", currentTime, "，结束时间", r.status.endTime, "======")
	}

	if openTime <= 0 && openTime%2 == 0 {
		fmt.Println("==========================================", r.name, " 游戏开奖=======================================================", r.count)
		r.makeActions()
	}
}

func (r *RoundBasedRobot) makeActions() {
	res, err := r.pushAction()
	if err != nil {
		fmt.Println(err)
		// 有时候开奖， 重试一次
		res, err = r.pushAction()
	}
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}

func (r *RoundBasedRobot) pushAction() (string, error) {
	keys, err := r.services.digestSigner.AvailableKeys()
	digest, err := hex.DecodeString(r.status.seed)
	sig, err := r.services.digestSigner.SignDigest(digest, keys[0])
	data := RevealData{eos.Uint64(r.status.roundId), sig}
	action := eos.Action{
		Account: eos.AccountName(r.name),
		Name:    "reveal",
		Authorization: []eos.PermissionLevel{
			{Actor: eos.AccountName(r.config.actorAccountName), Permission: eos.PN("active")}, //owner active
		},
		ActionData: eos.NewActionData(&data),
	}
	tx := eos.NewTransaction([]*eos.Action{&action}, &r.services.txOpts)
	//fmt.Println(tx)
	signedTx, packedTx, err := r.services.api.SignTransaction(tx, r.services.txOpts.ChainID, eos.CompressionNone)
	if err != nil {
		return "", err
	}
	_, err = json.MarshalIndent(signedTx, "", "")
	if err != nil {
		return "", err
	}
	_, err = json.Marshal(packedTx)
	response, err := r.services.api.PushTransaction(packedTx)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("response", response)
	return "", nil
}

func (r *RoundBasedRobot) getStatus() RoundStatus {
	body, err := getTableRows(r.config.node, r.name, "activegame")
	if err == nil {
		var list bet.Playtable
		err = json.Unmarshal(body, &list)

		if err == nil {
			for _, row := range list.Rows {
				return RoundStatus{row.ID, row.Status, int64(row.EndTime), row.Seed}
			}
		}
	}
	return RoundStatus{0, 0, 0, ""}
}
