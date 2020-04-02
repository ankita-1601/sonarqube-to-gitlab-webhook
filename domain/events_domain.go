package domain

//Events struct from sonarqube
type Events struct {
	ServerURL   string            `json:"serverUrl"`
	TaskID      string            `json:"taskId"`
	Status      string            `json:"status"`
	AnalysedAt  string            `json:"analysedAt"`
	Revision    string            `json:"revision"`
	ChangedAt   string            `json:"changedAt"`
	Project     Project           `json:"project"`
	Branch      Branch            `json:"branch"`
	QualityGate QualityGate       `json:"qualityGate"`
	Properties  map[string]string `json:"properties"`
}

// Project struct
type Project struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

//Branch struct
type Branch struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	IsMain bool   `json:"isMain"`
	URL    string `json:"url"`
}

//Conditions struct
type Conditions struct {
	Metric         string `json:"metric"`
	Operator       string `json:"operator"`
	Value          string `json:"value,omitempty"`
	Status         string `json:"status"`
	ErrorThreshold string `json:"errorThreshold"`
}

//QualityGate struct
type QualityGate struct {
	Name       string       `json:"name"`
	Status     string       `json:"status"`
	Conditions []Conditions `json:"conditions"`
}

// //Properties struct
// type Properties struct {
// 	EnabledGitlabPost    bool
// 	DisableQualityReport bool
// }
