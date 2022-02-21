package annotations

import (
	"SpotEvictHandler/internal/pkg/constants"
	"SpotEvictHandler/internal/pkg/types"
	"strconv"

	coreV1 "k8s.io/api/core/v1"
)

func GetSpotnikAnnotationsValueFromPod(pod *coreV1.Pod) types.Annotations {
	annotations := types.Annotations{}

	// get anottations from pod if exists

	// annotations for eviction of pod
	if val, ok := pod.ObjectMeta.Annotations[constants.EvictPodAnnotation]; ok {
		value, _ := strconv.ParseBool(val)
		annotations.EvictPod = value
	} else {
		value, _ := strconv.ParseBool(constants.EvictPodAnnotationDefault)
		annotations.EvictPod = value
	}

	// annotations for pod internal webhook
	if val, ok := pod.ObjectMeta.Annotations[constants.SendPodWebhookAnnotation]; ok {
		value, _ := strconv.ParseBool(val)
		if value {
			annotations.SendPodWebhook = value
		} else {
			value, _ := strconv.ParseBool(constants.SendPodWebhookAnnotationDefault)
			annotations.SendPodWebhook = value
		}
		annotations.PodWebhookPath = constants.PodWebhookPath
		annotations.PodWebhookPort = constants.PodWebhookPort
	}

	// annotations for external webhook
	if val, ok := pod.ObjectMeta.Annotations[constants.SendExternalWebhookAnnotation]; ok {
		value, _ := strconv.ParseBool(val)
		if value {
			annotations.SendExternalWebhook = value
		} else {
			value, _ := strconv.ParseBool(constants.SendExternalWebhookAnnotationDefault)
			annotations.SendExternalWebhook = value
		}
		annotations.ExternalWebhookURL = constants.ExternalWebhookURL
	}

	return annotations

}
