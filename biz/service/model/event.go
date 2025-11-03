package model

type EventReq struct {
	Uid            string
	EventName      string
	EventOrganizer string
	AwardTime      int64
}

type Event struct {
	EventId        string
	EventLevel     string
	RecognizeId    string
	AwardLevel     string
	AwardContent   string
	Uid            string
	EventName      string
	EventOrganizer string
	MaterialUrl    string
	MaterialStatus string
	AutoExtracted  bool
	AwardTime      string
	CreateAT       int64
	UpdateAT       int64
	DeleteAT       int64
}
