package model

type ScoreRule struct {
	RuleId            string  `json:"rule_id"`
	RecognizedEventId string  `json:"recognized_event_id"`
	EventLevel        string  `json:"event_level"`
	EventWeight       float64 `json:"score_weight"`
	Integral          int64   `json:"integral"`
	RuleDesc          string  `json:"rule_desc"`
	IsEditable        bool    `json:"is_editable"`
	AwardLevel        string  `json:"award_level"`
	AwardLevelWeight  float64 `json:"award_level_weight"`
	CreateAT          int64
	UpdateAT          int64
	DeleteAT          int64
}

type ScoreRecord struct {
	ResultId      string  `json:"result_id"`
	EventId       string  `json:"event_id"`
	UserId        string  `json:"user_id"`
	RuleId        string  `json:"rule_id"`
	AppealId      string  `json:"appeal_id"`
	FinalIntegral float64 `json:"final_integral"`
	Status        string  `json:"status"`
	CreateAT      int64
	UpdateAT      int64
	DeleteAT      int64
}
