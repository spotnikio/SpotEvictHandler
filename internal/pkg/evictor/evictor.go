package evictor

import (
	"SpotEvictHandler/internal/pkg/controller"
	"SpotEvictHandler/internal/pkg/types"
	"SpotEvictHandler/internal/pkg/util"
	"SpotEvictHandler/internal/pkg/util/annotations"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	coreV1 "k8s.io/api/core/v1"
	policy "k8s.io/api/policy/v1beta1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func Evictor(c *controller.Controller) {
	start := time.Now()

	//corson the node befor start to evict pods
	err := util.CordonNode(c.Client, c.NodeName)
	if err != nil {
		logrus.Fatal("canot evict the node: %s", err)
		return
	}

	plist := util.GetPodFilterdByNode(c.PodsIndexer, c.NodeName)
	for _, p := range plist.Items {
		annotations := annotations.GetSpotnikAnnotationsValueFromPod(&p)
		eventData := types.EventToSend{ComputerName: c.NodeName}
		if annotations.SendPodWebhook {
			b, _ := json.Marshal(eventData)
			logrus.Infof("event send to: ", p.ObjectMeta.Name)
			postReqBody := bytes.NewBuffer(b)
			http.Post(fmt.Sprintf("http://%s:%s/%s", p.Status.PodIP, annotations.PodWebhookPort, annotations.PodWebhookPath), "application/json", postReqBody)
		}
		if annotations.SendExternalWebhook {
			b, _ := json.Marshal(eventData)
			logrus.Infof("event send to: ", p.ObjectMeta.Name)
			postReqBody := bytes.NewBuffer(b)
			http.Post(annotations.ExternalWebhookURL, "application/json", postReqBody)
		}
		if annotations.EvictPod {
			EvictPod(c.Client, &p)
		}
		time.Sleep(10 * time.Millisecond)
	}

	elapsed := time.Since(start)
	logrus.Infof("Evictor: evictor time is -> %s", elapsed)
}

func EvictPod(client kubernetes.Interface, p *coreV1.Pod) {
	client.PolicyV1beta1().Evictions(p.ObjectMeta.Namespace).Evict(context.TODO(), &policy.Eviction{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      p.ObjectMeta.Name,
			Namespace: p.ObjectMeta.Namespace,
		},
	})
}
