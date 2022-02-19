package constants

import "SpotEvictHandler/internal/pkg/util"

var (
	EvictPodAnnotation string = util.Getenv("EVICT_POD_ANNOTATION", "spotnik.io/enabled")

	EvictPodAnnotationDefault string = util.Getenv("EVICT_POD_ANNOTATION_DEFAULT", "true")

	SendPodWebhookAnnotation string = util.Getenv("SEND_POD_WEBHOOK_ANNOTATION", "spotnik.io/sendpodwebhook")

	SendPodWebhookAnnotationDefault string = util.Getenv("SEND_POD_WEBHOOK_ANNOTATION_DEFAULT", "false")

	PodWebhookPath string = util.Getenv("POD_WEBHOOK_PATH", "/")

	PodWebhookPort string = util.Getenv("POD_WEBHOOK_PORT", "80")

	SendExternalWebhookAnnotation string = util.Getenv("SEND_EXTERNAL_WEBHOOK_ANNOTATION", "spotnik.io/sendexternalwebhook")

	SendExternalWebhookAnnotationDefault string = util.Getenv("SEND_EXTERNAL_WEBHOOK_ANNOTATION_DEFAULT", "false")

	ExternalWebhookURL string = util.Getenv("EXTERNAL_WEBHOOK_URL", "")

	LogstashHostName string = util.Getenv("LogConnectionString", "")

	LogstashPort string = util.Getenv("LogPort", "9999")

	CloudProviderName string = util.Getenv("CLOUD_PROVIDER", "azure")
)
