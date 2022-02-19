package controller

import (
	"SpotEvictHandler/internal/pkg/constants"
	"fmt"

	"github.com/sirupsen/logrus"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	utilRuntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"SpotEvictHandler/internal/pkg/util"
)

type Controller struct {
	Client        kubernetes.Interface
	PodsIndexer   cache.Indexer
	PodsInformer  cache.Controller
	NodeName      string
	IsInitilize   bool
	CloudProvider string
}

func NewController(client kubernetes.Interface) (*Controller, error) {

	c := Controller{
		Client:      client,
		IsInitilize: false,
	}

	// creat the watcher for pods and nodes
	podsListWatcher := cache.NewListWatchFromClient(client.CoreV1().RESTClient(), "pods", metaV1.NamespaceAll, fields.Everything())

	// create the informer for nodes and pods
	podsIndexer, podsInformer := cache.NewIndexerInformer(podsListWatcher, &coreV1.Pod{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc:    c.Add,
		UpdateFunc: c.Update,
		DeleteFunc: c.Delete,
	}, cache.Indexers{})

	//set the parameters on the Controller struct
	c.PodsInformer = podsInformer
	c.PodsIndexer = podsIndexer

	//get env name from environment variable
	c.NodeName = util.Getenv("POD_NODE_NAME", "")
	logrus.Infof("Node: %s is detected", c.NodeName)

	// get the cloud proviser
	c.CloudProvider = constants.CloudProviderName
	logrus.Infof("Cloud provider is %s", c.CloudProvider)

	return &c, nil
}

func RunController(c *Controller, stop chan struct{}) {

	logrus.Infof("Starting Controller")

	// Run the informer on go routines
	go c.PodsInformer.Run(stop)

	// Wait for all cache to be sync with the cluster
	if !cache.WaitForCacheSync(stop, c.PodsInformer.HasSynced) {
		utilRuntime.HandleError(fmt.Errorf("timed out waiting for pode caches to sync"))
	}
	if c.NodeName != "" && c.CloudProvider != "" {
		c.IsInitilize = true
	}

}

// Funcs to handle event from the cluster (we need them to be exist for the indexer to be updated)
func (c *Controller) Add(obj interface{}) {
	// fmt.Println("Add func")
}
func (c *Controller) Update(old interface{}, new interface{}) {
	// fmt.Println("Update func")
}
func (c *Controller) Delete(obj interface{}) {
	// fmt.Println("Delete func")
}
