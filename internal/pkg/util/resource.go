package util

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"

	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"SpotEvictHandler/internal/pkg/types"

	kubeTypes "k8s.io/apimachinery/pkg/types"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

func GetPodFilterdByNode(podsIndex cache.Indexer, nodeName string) coreV1.PodList {
	plist := coreV1.PodList{}

	for _, i := range podsIndex.List() {
		pod := i.(*coreV1.Pod)
		if pod.Spec.NodeName == nodeName && pod.ObjectMeta.OwnerReferences[0].Kind == "ReplicaSet" {
			plist.Items = append(plist.Items, *pod)
		}
	}
	return plist
}

func GetDeploymentFromIndexer(deploymentIndexer cache.Indexer, deploymentName string, namespace string) *appsV1.Deployment {

	for _, d := range deploymentIndexer.List() {
		deployment := d.(*appsV1.Deployment)
		if deployment.ObjectMeta.Name == deploymentName && deployment.ObjectMeta.Namespace == namespace {
			return deployment
		}
	}
	return nil
}

func GetNamespaceOfPod(namespacesIndex cache.Indexer, pod *coreV1.Pod) (*coreV1.Namespace, error) {

	for _, i := range namespacesIndex.List() {
		ns := i.(*coreV1.Namespace)
		if pod.ObjectMeta.Namespace == ns.ObjectMeta.Name {
			return ns, nil
		}
	}
	return nil, fmt.Errorf("namespace %s not found", pod.ObjectMeta.Namespace)
}

func GetNamespaceOfDeployment(namespacesIndex cache.Indexer, deployment *appsV1.Deployment) (*coreV1.Namespace, error) {

	for _, i := range namespacesIndex.List() {
		ns := i.(*coreV1.Namespace)
		if deployment.ObjectMeta.Namespace == ns.ObjectMeta.Name {
			return ns, nil
		}
	}
	return nil, fmt.Errorf("namespace %s not found", deployment.ObjectMeta.Namespace)
}

func CordonNode(client kubernetes.Interface, nodeName string) error {

	payload := []types.NodePatchStringValue{{
		Op:    "replace",
		Path:  "/spec/unschedulable",
		Value: true,
	}}
	payloadBytes, _ := json.Marshal(payload)
	_, err := client.CoreV1().Nodes().Patch(context.TODO(), nodeName, kubeTypes.JSONPatchType, payloadBytes, metaV1.PatchOptions{})
	return err
}

// Function: Read and return JSON
func GetEvent(url string, headerName string, headerValue string) (string, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set(headerName, headerValue)
	resp, err := client.Do(req)
	if err != nil {
		logrus.Fatal(err)
		return "", errors.New("Could not get json from url")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Fatal(err)
		return "", errors.New("Could not read json body")
	}

	sb := string(body)

	return sb, nil
}
