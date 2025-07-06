package errors

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type FilterParams struct {
	ProjectID   string
	Fingerprint string
	TimeFrom    int64
	TimeTo      int64
	Search      string
}

type GetAllParams struct {
	FilterParams
	SortOrder string `validate:"omitempty,oneof=asc desc"`
	Limit     int
	Offset    int
}

type Create struct {
	Time       int64        `json:"time" validate:"required" example:"1704067200000" format:"int64"`
	Message    string       `json:"message" validate:"required" example:"Division by zero in calculate()"`
	Stacktrace *interface{} `json:"stacktrace" validate:"required"`
	File       string       `json:"file" validate:"required" example:"/var/www/app/index.php"`
	Line       int          `json:"line" validate:"required" example:"15"`
	// Context can be any JSON value
	// @Schema(
	//   oneOf={
	//     string,
	//     object,
	//     array,
	//     number,
	//     boolean,
	//     null
	//   },
	//   example={"key":"value"}
	// )
	Context *interface{} `json:"context"`
	IP      *string      `json:"ip,omitempty" example:"192.168.1.1"`
	URL     *string      `json:"url,omitempty" example:"https://example.com/api/v1/calculate"`
	Method  *string      `json:"method,omitempty" example:"POST"`
	// @Schema(
	//   type = "object",
	//   example = `{"Content-Type": "application/json", "Authorization": "Bearer token"}`
	// )
	Headers *map[string]interface{} `json:"headers,omitempty"`
	// @Schema(
	//   type = "object",
	//   example = `{"page": 1, "limit": 10}`
	// )
	QueryParams *map[string]interface{} `json:"queryParams,omitempty"`
	// @Schema(
	//   type = "object",
	//   example = `{"a": 5, "b": 0}`
	// )
	BodyParams *map[string]interface{} `json:"bodyParams,omitempty"`
	// @Schema(
	//   type = "object",
	//   example = `{"sessionId": "abc123", "theme": "dark"}`
	// )
	Cookies *map[string]interface{} `json:"cookies,omitempty"`
	// @Schema(
	//   type = "object",
	//   example = `{"userId": 123, "role": "admin"}`
	// )
	Session *map[string]interface{} `json:"session,omitempty"`
	// @Schema(
	//   type = "object",
	//   example = `{"avatar": "avatar.jpg", "size": 1024}`
	// )
	Files *map[string]interface{} `json:"files,omitempty"`
	// @Schema(
	//   type = "object",
	//   example = `{"APP_ENV": "production", "DB_HOST": "db.example.com"}`
	// )
	Env       *map[string]interface{} `json:"env,omitempty"`
	ProjectID string                  `json:"-"`
}

type Update struct {
	Message    string       `json:"message" validate:"required" example:"Error message"`
	Stacktrace *interface{} `json:"stacktrace" validate:"required"`
	File       string       `json:"file" validate:"required" example:"index.php"`
	Line       int          `json:"line" validate:"required" example:"1"`
	// Context can be any JSON value
	// @Schema(
	//   oneOf={
	//     string,
	//     object,
	//     array,
	//     number,
	//     boolean,
	//     null
	//   },
	//   example={"key":"value"}
	// )
	Context *interface{} `json:"context"`
}

type Entity struct {
	ID         string       `json:"id" example:"a08929b5-d4f0-4ceb-9cfe-bb4fc05b030c"`
	Message    string       `json:"message" validate:"required" example:"Error: Division by zero"`
	Stacktrace *interface{} `json:"stacktrace" validate:"required"`
	File       string       `json:"file" validate:"required" example:"/var/www/index.php"`
	Line       int          `json:"line" validate:"required" example:"15"`
	// Context can be any JSON value
	// @Schema(
	//   oneOf={
	//     string,
	//     object,
	//     array,
	//     number,
	//     boolean,
	//     null
	//   },
	//   example={"key":"value"}
	// )
	Context     *interface{}            `json:"context"`
	IP          *string                 `json:"ip" example:"192.168.1.1"`
	URL         *string                 `json:"url" example:"https://example.com/api/v1/calculate"`
	Method      *string                 `json:"method" example:"POST"`
	Headers     *map[string]interface{} `json:"headers"`
	QueryParams *map[string]interface{} `json:"queryParams"`
	BodyParams  *map[string]interface{} `json:"bodyParams"`
	Cookies     *map[string]interface{} `json:"cookies"`
	Session     *map[string]interface{} `json:"session"`
	Files       *map[string]interface{} `json:"files"`
	Env         *map[string]interface{} `json:"env"`
	Time        int64                   `json:"time" example:"1704067200000"` // Unix timestamp in milliseconds
}

type EntityList struct {
	Count int      `json:"count"`
	Items []Entity `json:"items"`
}

type Stats struct {
	Last24h int `json:"last24h"`
	Last7d  int `json:"last7d"`
	Last30d int `json:"last30d"`
}
