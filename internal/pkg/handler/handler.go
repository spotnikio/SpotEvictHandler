package handler

import (
	"SpotEvictHandler/internal/pkg/controller"
	"SpotEvictHandler/internal/pkg/evictor"
	"SpotEvictHandler/internal/pkg/types"
	"SpotEvictHandler/internal/pkg/util"

	"encoding/json"
	"errors"

	"github.com/sirupsen/logrus"
)

func Run(c *controller.Controller) {
	switch cloudProvider := c.CloudProvider; cloudProvider {
	case "azure":
		JsonEventsOutput, err := util.GetEvent("http://169.254.169.254/metadata/scheduledevents?api-version=2020-07-01", "Metadata", "true")
		if err != nil {
			logrus.Fatal(err)
			return
		}
		isPreempt, err := azureCheckPreempt(JsonEventsOutput)
		if err != nil {
			logrus.Fatal(err)
			return
		}
		if isPreempt {
			logrus.Infof("New evict event was catched: %s", c.NodeName)
			evictor.Evictor(c)
		}
		//wait peacefully to die, after saving all his family and friends
		select {}

	case "google":
		JsonEventsOutput, err := util.GetEvent("http://metadata.google.internal/computeMetadata/v1/instance/scheduling/preemptible", "Metadata-Flavor", "Google")
		if err != nil {
			logrus.Fatal(err)
			return
		}
		isPreempt, err := googleCheckPreempt(JsonEventsOutput)
		if err != nil {
			logrus.Fatal(err)
			return
		}
		if isPreempt {
			logrus.Infof("New evict event was catched: %s", c.NodeName)
			evictor.Evictor(c)
		}
		//wait peacefully to die, after saving all his family and friends
		select {}

	case "amazon":
		JsonEventsOutput, err := util.GetEvent("http://169.254.169.254/latest/meta-data/spot/instance-action", "Metadata", "true")
		if err != nil {
			logrus.Fatal(err)
			return
		}
		isPreempt, err := amazonCheckPreempt(JsonEventsOutput)
		if err != nil {
			logrus.Fatal(err)
			return
		}
		if isPreempt {
			logrus.Infof("New evict event was catched: %s", c.NodeName)
			evictor.Evictor(c)
		}
		//wait peacefully to die, after saving all his family and friends
		select {}

	default:
		logrus.Fatal("this Cloud Provider: %s is not supported", cloudProvider)
	}

}

// Function: Checks if a VM has an event called "preempt"
func azureCheckPreempt(JsonEventsOutput string) (bool, error) {
	event := types.AzureEvent{
		Arr: []types.AzureEventList{},
	}
	logrus.Infoln("JSON events output: " + JsonEventsOutput)
	err := json.Unmarshal([]byte(JsonEventsOutput), &event)
	if err != nil {
		logrus.Fatal(err)
		return false, errors.New("could not unmarshal")
	}

	for i := 0; i < len(event.Arr); i++ {
		if event.Arr[i].EventType == "Preempt" {
			EventFromJson := event.Arr[i].EventType
			logrus.Infoln("Event from JSON: " + EventFromJson)
			return true, nil
		} else {
			return false, nil
		}

	}
	return false, nil
}

func googleCheckPreempt(JsonEventsOutput string) (bool, error) {
	logrus.Infoln("JSON events output: " + JsonEventsOutput)
	if JsonEventsOutput == "TRUE" {
		return true, nil
	}
	return false, nil
}

func amazonCheckPreempt(JsonEventsOutput string) (bool, error) {
	event := types.AmazonEvent{}

	logrus.Infoln("JSON events output: " + JsonEventsOutput)
	err := json.Unmarshal([]byte(JsonEventsOutput), &event)
	if err != nil {
		logrus.Fatal(err)
		return false, errors.New("could not unmarshal")
	}

	if event.Action == "terminate" || event.Action == "stop" {
		return true, nil
	}
	return false, nil
}
