package errors

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

type Create struct {
	Time       int64  `json:"time" validate:"required" example:"1704067200000" format:"int64"`
	Message    string `json:"message" validate:"required" example:"Error message"`
	Stacktrace string `json:"stacktrace" validate:"required" example:"Stacktrace"`
	File       string `json:"file" validate:"required" example:"index.php"`
	Line       int    `json:"line" validate:"required" example:"1"`
	Context    string `json:"context" example:"message context"`
	ProjectID  string
}

type Update struct {
	Message    string `json:"message" validate:"required" example:"Error message"`
	Stacktrace string `json:"stacktrace" validate:"required" example:"Stacktrace"`
	File       string `json:"file" validate:"required" example:"index.php"`
	Line       int    `json:"line" validate:"required" example:"1"`
	Context    string `json:"context" example:"message context"`
}

type Entity struct {
	ID         string `json:"id" example:"a08929b5-d4f0-4ceb-9cfe-bb4fc05b030c"`
	Message    string `json:"message" validate:"required" example:"Error message"`
	Stacktrace string `json:"stacktrace" validate:"required" example:"Stacktrace"`
	File       string `json:"file" validate:"required" example:"index.php"`
	Line       int    `json:"line" validate:"required" example:"1"`
	Context    string `json:"context" example:"message context"`
	Time       int64  `json:"time" example:"1704067200000"`
}

type EntityList struct {
	Count int      `json:"count"`
	Items []Entity `json:"items"`
}
