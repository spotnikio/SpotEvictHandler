package types

type Annotations struct {
	EvictPod            bool
	SendPodWebhook      bool
	PodWebhookPath      string
	PodWebhookPort      string
	SendExternalWebhook bool
	ExternalWebhookURL  string
}

type NodePatchStringValue struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value bool   `json:"value"`
}

type AzureEventList struct {
	EventType string `json:"EventType"`
}

type AzureEvent struct {
	Arr []AzureEventList `json:"Events"`
}

type AmazonEvent struct {
	Action string `json:"action"`
	Time   string `json:"time"`
}

type EventToSend struct {
	ComputerName string `json:"computerName"`
}
