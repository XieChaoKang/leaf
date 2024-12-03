package log

type GameRecordLog struct {
	ID                    int64   `json:"id"`
	SubAnm                string  `json:"sub_anm"`
	Did                   string  `json:"did"`
	PkgName               string  `json:"pkg_name"`
	MatchId               string  `json:"match_id"`
	PlayType              string  `json:"play_type"`
	RoomType              string  `json:"room_type"`
	Score                 float64 `json:"score"`
	Uid                   string  `json:"uid"`
	PlayerBet             float64 `json:"player_bet"`
	Prize                 int64   `json:"prize"`
	Tax                   float64 `json:"tax"`
	FlagInfo              string  `json:"flag_info"`
	GameInfo              string  `json:"game_info"`
	Ext                   string  `json:"ext"`
	StarTime              int64   `json:"star_time"`
	EndTime               string  `json:"end_time"`
	Ctime                 string  `json:"ctime"`
	PersonalWithdrawValue int     `json:"personal_withdraw_value"`
	LogId                 string  `json:"log_id"`
	CurrencyType int `json:"currency_type"`
	PotWin int `json:"pot_win"`
}

type EventLog struct {
	ID                    int64  `json:"id"`
	Uid                   string `json:"uid"`
	DId                   string `json:"did"`
	EventName             string `json:"event_name"`
	PersonalWithdrawValue int    `json:"personal_withdraw_value"`
	EventValue            string `json:"event_value"`
}
