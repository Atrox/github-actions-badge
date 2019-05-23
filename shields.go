package main

type Endpoint struct {
	SchemaVersion int    `json:"schemaVersion,omitempty"`
	Label         string `json:"label,omitempty"`
	Message       string `json:"message,omitempty"`
	Color         string `json:"color,omitempty"`
	LabelColor    string `json:"labelColor,omitempty"`
	IsError       bool   `json:"isError,omitempty"`
	NamedLogo     string `json:"namedLogo,omitempty"`
	LogoSVG       string `json:"logoSvg,omitempty"`
	LogoColor     string `json:"logoColor,omitempty"`
	LogoWidth     int    `json:"logoWidth,omitempty"`
	LogoPosition  string `json:"logoPosition,omitempty"`
	Style         string `json:"style,omitempty"`
	CacheSeconds  int    `json:"cache_seconds,omitempty"`
}

func NewEndpoint() *Endpoint {
	return &Endpoint{
		SchemaVersion: 1,
		CacheSeconds:  300,

		Label:     "GitHub Actions",
		NamedLogo: "github",
	}
}

func (e *Endpoint) Success() {
	e.Color = "success"
	e.Message = "success"
}

func (e *Endpoint) Neutral() {
	e.Color = "success"
	e.Message = "neutral"
}

func (e *Endpoint) Pending() {
	e.Color = "yellow"
	e.Message = "pending"
	e.IsError = true
}

func (e *Endpoint) Failure() {
	e.Color = "critical"
	e.Message = "failure"
	e.IsError = true
}

func (e *Endpoint) Cancelled() {
	e.Color = "inactive"
	e.Message = "cancelled"
	e.IsError = true
}

func (e *Endpoint) TimedOut() {
	e.Color = "critical"
	e.Message = "timed out"
	e.IsError = true
}

func (e *Endpoint) ActionRequired() {
	e.Color = "critical"
	e.Message = "action required"
	e.IsError = true
}

func (e *Endpoint) ServerError() {
	e.Color = "inactive"
	e.Message = "server error"
	e.IsError = true
}
