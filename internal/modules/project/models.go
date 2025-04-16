package project

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type GetAllParams struct {
	SortOrder string `validate:"omitempty,oneof=asc desc"`
	Limit     int
	Offset    int
}

type Create struct {
	Name string `json:"name" validate:"required" example:"New project"`
}

type Update struct {
	Name string `json:"name" validate:"required" example:"New project"`
}

type Entity struct {
	ID   string `json:"id" example:"a08929b5-d4f0-4ceb-9cfe-bb4fc05b030c"`
	Name string `json:"name" example:"New project"`
}

type EntityList struct {
	Count int      `json:"count"`
	Items []Entity `json:"items"`
}
