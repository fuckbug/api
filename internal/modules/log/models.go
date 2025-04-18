package log

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type FilterParams struct {
	ProjectID   string
	TimeFrom    int
	TimeTo      int
	LevelFilter string
	SearchQuery string
}

type GetAllParams struct {
	FilterParams
	SortOrder string `validate:"omitempty,oneof=asc desc"`
	Limit     int
	Offset    int
}

type Create struct {
	Time      int64  `json:"time" validate:"required" example:"1704067200000" format:"int64"`
	Level     string `json:"level" validate:"required,oneof=DEBUG INFO WARN ERROR"`
	Message   string `json:"message" validate:"required" example:"first log message"`
	Context   string `json:"context" example:"message context"`
	ProjectID string
}

type Update struct {
	Level   string `json:"level" validate:"omitempty,oneof=DEBUG INFO WARN ERROR"`
	Message string `json:"message" validate:"omitempty" example:"updated log message"`
	Context string `json:"context" example:"message context"`
}

type Entity struct {
	ID      string `json:"id" example:"a08929b5-d4f0-4ceb-9cfe-bb4fc05b030c"`
	Level   string `json:"level" example:"INFO"`
	Message string `json:"message" example:"first log message"`
	Context string `json:"context" example:"message context"`
	Time    int64  `json:"time" example:"1704067200000"`
}

type EntityList struct {
	Count int      `json:"count"`
	Items []Entity `json:"items"`
}
