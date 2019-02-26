package bet


type Playtable struct {
	Rows []struct {
		ID               int    `json:"id"`
		EndTime          int    `json:"end_time"`
		PlayerCards      []int  `json:"player_cards"`
		BankerCards      []int  `json:"banker_cards"`
		Symbol           string `json:"symbol"`
		Status           int    `json:"status"`
		LargestWinner    string `json:"largest_winner"`
		LargestWinAmount string `json:"largest_win_amount"`
		Seed             string `json:"seed"`
	} `json:"rows"`
	More bool `json:"more"`
}

type Bets struct {
	Rows []struct {
		ID      int    `json:"id"`
		GameID  int    `json:"game_id"`
		Player  string `json:"player"`
		Referer string `json:"referer"`
		Bet     string `json:"bet"`
		BetType int    `json:"bet_type"`
	} `json:"rows"`
	More bool `json:"more"`
}
type PlayerBet struct {
	ID      int    `json:"id"`
	GameID  int    `json:"game_id"`
	Player  string `json:"player"`
	Referer string `json:"referer"`
	Bet     string `json:"bet"`
	BetType int    `json:"bet_type"`
}
type PlayerAmount struct {
	Player string `json:"player"`
	Bet    string `json:"bet"`
}

type BetsAmount struct {
	Rows []struct {
		Player string `json:"player"`
		Bet    string `json:"bet"`
	} `json:"rows"`
	More bool `json:"more"`
}


