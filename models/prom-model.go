package models

type AlertData struct {
	Status            string      `json:"status"`
	Alerts            []Alert     `json:"alerts"`
	CommonAnnotations Annotations `json:"commonAnnotations"`
	CommonLabels      Labels      `json:"commonLabels"`
	ExternalURL       string      `json:"externalURL"`
	GroupKey          string      `json:"groupKey"`
	GroupLabels       Labels      `json:"groupLabels"`
	Receiver          string      `json:"receiver"`
	TruncatedAlerts   int         `json:"truncatedAlerts"`
	Version           string      `json:"version"`
}

type Alert struct {
	Annotations  Annotations `json:"annotations"`
	Labels       Labels      `json:"labels"`
	StartsAt     string      `json:"startsAt"`
	EndsAt       string      `json:"endsAt"`
	Fingerprint  string      `json:"fingerprint"`
	GeneratorURL string      `json:"generatorURL"`
	Status       any
}

type Annotations struct {
	Description string `json:"description"`
	RunbookURL  string `json:"runbook_url"`
	Summary     string `json:"summary"`
}

type Labels struct {
	Alertname  string `json:"alertname"`
	Container  string `json:"container"`
	Controller string `json:"controller"`
	Endpoint   string `json:"endpoint"`
	Instance   string `json:"instance"`
	Job        string `json:"job"`
	Namespace  string `json:"namespace"`
	Pod        string `json:"pod"`
	Prometheus string `json:"prometheus"`
	Resource   string `json:"resource"`
	Service    string `json:"service"`
	Severity   string `json:"severity"`
	State      string `json:"state"`
}
