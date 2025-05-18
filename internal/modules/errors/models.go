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
	Time        int64   `json:"time" validate:"required" example:"1704067200000" format:"int64"`
	Message     string  `json:"message" validate:"required" example:"Division by zero in calculate()"`
	Stacktrace  string  `json:"stacktrace" validate:"required" example:"at index.php:15, at main(), at calculate()"`
	File        string  `json:"file" validate:"required" example:"/var/www/app/index.php"`
	Line        int     `json:"line" validate:"required" example:"15"`
	Context     *string `json:"context" example:"{\"userId\": 123, \"action\": \"calculate\", \"input\": {\"a\": 5, \"b\": 0}}"`
	IP          *string `json:"ip,omitempty" example:"192.168.1.1"`
	URL         *string `json:"url,omitempty" example:"https://example.com/api/v1/calculate"`
	Method      *string `json:"method,omitempty" example:"POST"`
	Headers     *string `json:"headers,omitempty" example:"{\"Content-Type\": \"application/json\", \"Authorization\": \"Bearer token\"}"`
	QueryParams *string `json:"queryParams,omitempty" example:"{\"page\": 1, \"limit\": 10}"`
	BodyParams  *string `json:"bodyParams,omitempty" example:"{\"a\": 5, \"b\": 0}"`
	Cookies     *string `json:"cookies,omitempty" example:"{\"sessionId\": \"abc123\", \"theme\": \"dark\"}"`
	Session     *string `json:"session,omitempty" example:"{\"userId\": 123, \"role\": \"admin\"}"`
	Files       *string `json:"files,omitempty" example:"{\"avatar\": \"avatar.jpg\", \"size\": 1024}"`
	Env         *string `json:"env,omitempty" example:"{\"APP_ENV\": \"production\", \"DB_HOST\": \"db.example.com\"}"`
	ProjectID   string  `json:"-"`
}

type Update struct {
	Message    string  `json:"message" validate:"required" example:"Error message"`
	Stacktrace string  `json:"stacktrace" validate:"required" example:"Stacktrace"`
	File       string  `json:"file" validate:"required" example:"index.php"`
	Line       int     `json:"line" validate:"required" example:"1"`
	Context    *string `json:"context" example:"message context"`
}

type Entity struct {
	ID          string  `json:"id" example:"a08929b5-d4f0-4ceb-9cfe-bb4fc05b030c"`
	Message     string  `json:"message" validate:"required" example:"Error: Division by zero"`
	Stacktrace  string  `json:"stacktrace" validate:"required" example:"at index.php:15, at main()"`
	File        string  `json:"file" validate:"required" example:"/var/www/index.php"`
	Line        int     `json:"line" validate:"required" example:"15"`
	Context     *string `json:"context" example:"{\"userId\": 123, \"action\": \"calculate\"}"`
	IP          *string `json:"ip" example:"192.168.1.1"`
	URL         *string `json:"url" example:"https://example.com/api/v1/calculate"`
	Method      *string `json:"method" example:"POST"`
	Headers     *string `json:"headers" example:"{\"Content-Type\": \"application/json\", \"Authorization\": \"Bearer token\"}"`
	QueryParams *string `json:"queryParams" example:"{\"page\": 1, \"limit\": 10}"`
	BodyParams  *string `json:"bodyParams" example:"{\"a\": 5, \"b\": 0}"`
	Cookies     *string `json:"cookies" example:"{\"sessionId\": \"abc123\", \"theme\": \"dark\"}"`
	Session     *string `json:"session" example:"{\"userId\": 123, \"role\": \"admin\"}"`
	Files       *string `json:"files" example:"{\"avatar\": \"avatar.jpg\", \"size\": 1024}"`
	Env         *string `json:"env" example:"{\"APP_ENV\": \"production\", \"DB_HOST\": \"db.example.com\"}"`
	Time        int64   `json:"time" example:"1704067200000"` // Unix timestamp in milliseconds
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
