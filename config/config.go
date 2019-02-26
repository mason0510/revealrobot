package config


import (
"encoding/json"
)
//固定地址保存在数据库中
type Config struct {
	Debug                bool   `json:"debug"`
	RemoteAddress        string `json:"remote_address"`
	RedisHost            string `json:"redis_host"`
	RedisPass            string `json:"redis_pass"`
	Crypto               string `json:"crypto"`
	Use                  string `json:"Use"`
	Port                 string `json:"port"`
	Allow                string `json:"allow"`
	RobotAccount         string `json:"robot_account"`
	RobotOpenAccount     string `json:"robot_open_account"`
	RobotPrivate         string `json:"robot_private"`
	RobotOpenPrivate     string `json:"robot_open_private"`
	ContracAccount       string `json:"contrac_account"`
	MysqlConn            string `json:"mysql_conn"`
	Arena                int    `json:"arena"`
	Tablename            string `json:"tablename"`
	EosPermission        string `json:"eos_permission"`
	EosTablename         string `json:"eos_tablename"`
	Testnode             string `json:"testnode"`
	TestrevealKey        string `json:"testrevealKey"`
	TestactorAccountName string `json:"testactorAccountName"`
	TestactorAccountKey  string `json:"testactorAccountKey"`
	Node                 string `json:"node"`
	RevealKey            string `json:"revealKey"`
	ActorAccountName     string `json:"actorAccountName"`
	ActorAccountKey      string `json:"actorAccountKey"`
	TimeUrl              string `json:timeUrl`
}
var (
	C *Config
)

func InitConfig(data []byte) error {
	C = new(Config)
	return json.Unmarshal(data, &C)
}

func Port() string {
	if C.Port == "" {
		return ":9879"
	} else {
		return C.Port
	}
}

func Allow() string {
	if C.Allow == "" {
		return "*"
	} else {
		return C.Allow
	}
}

