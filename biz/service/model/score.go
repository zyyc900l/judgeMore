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

type RecognizedEvent struct {
	RecognizedEventId   string
	College             string
	RecognizedEventName string
	Organizer           string
	RecognizedEventTime string
	RelatedMajors       string
	ApplicableMajors    string
	RecognitionBasis    string
	RecognizedLevel     string
	IsActive            bool
	CreateAT            int64
	UpdateAT            int64
	DeleteAT            int64
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
type ViewRecognizedRewardReq struct {
	EventName         *string
	OrganizerName     *string
	RecognizedEventId *string
}

func (p *ViewRecognizedRewardReq) IsSetEventName() bool {
	return p.EventName != nil
}

func (p *ViewRecognizedRewardReq) IsSetOrganizerName() bool {
	return p.OrganizerName != nil
}

func (p *ViewRecognizedRewardReq) IsSetRecognizedEventId() bool {
	return p.RecognizedEventId != nil
}

var DEFAULT string

func (p *ViewRecognizedRewardReq) GetEventName() string {
	if !p.IsSetEventName() {
		return DEFAULT
	}
	return *p.EventName
}

func (p *ViewRecognizedRewardReq) GetOrganizerName() string {
	if !p.IsSetOrganizerName() {
		return DEFAULT
	}
	return *p.OrganizerName
}

func (p *ViewRecognizedRewardReq) GetRecognizedEventId() string {
	if !p.IsSetRecognizedEventId() {
		return DEFAULT
	}
	return *p.RecognizedEventId
}
