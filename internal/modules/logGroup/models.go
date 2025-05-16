package loggroup

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type FilterParams struct {
	ProjectID string
	TimeFrom  int64
	TimeTo    int64
	Level     string
	Search    string
}

type GetAllParams struct {
	FilterParams
	SortOrder string `validate:"omitempty,oneof=asc desc"`
	Limit     int
	Offset    int
}

type Entity struct {
	ID          string `json:"id" example:"a08929b5-d4f0-4ceb-9cfe-bb4fc05b030c"`
	Level       string `json:"level" example:"INFO"`
	Message     string `json:"message" validate:"required" example:"Log message"`
	FirstSeenAt int64  `json:"firstSeenAt" example:"1704067200"`
	LastSeenAt  int64  `json:"lastSeenAt" example:"1704067200"`
	Counter     int    `json:"counter" example:"18"`
}

type EntityList struct {
	Count int      `json:"count"`
	Items []Entity `json:"items"`
}
