package errorsgroup

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
	SearchQuery string
}

type GetAllParams struct {
	FilterParams
	SortOrder string `validate:"omitempty,oneof=asc desc"`
	Limit     int
	Offset    int
}

type Entity struct {
	ID          string `json:"id" example:"a08929b5-d4f0-4ceb-9cfe-bb4fc05b030c"`
	Message     string `json:"message" validate:"required" example:"Error message"`
	Stacktrace  string `json:"stacktrace" validate:"required" example:"Stacktrace"`
	File        string `json:"file" validate:"required" example:"index.php"`
	Line        int    `json:"line" validate:"required" example:"1"`
	FirstSeenAt int    `json:"firstSeenAt" example:"1704067200"`
	LastSeenAt  int    `json:"lastSeenAt" example:"1704067200"`
	Counter     int    `json:"counter" example:"18"`
}

type EntityList struct {
	Count int      `json:"count"`
	Items []Entity `json:"items"`
}
