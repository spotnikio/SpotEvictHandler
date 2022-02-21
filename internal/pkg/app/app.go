package app

import (
	"SpotEvictHandler/internal/pkg/constants"
	"SpotEvictHandler/internal/pkg/controller"
	"SpotEvictHandler/internal/pkg/handler"
	"SpotEvictHandler/internal/pkg/kube"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"time"

	//"github.com/prometheus/client_golang/prometheus"

	logrustash "github.com/bshuster-repo/logrus-logstash-hook"
	gas "github.com/firstrow/goautosocket"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type ServerParameters struct {
	port int // webhook server port
}

var parameters ServerParameters

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "SpotEvictHandler",
		Short: "Kubenretes controller to handle eviton of spoy unstance by the big pharms",
		Run:   Run,
	}

	// Add parmaeters to the Command
	cmd.Flags().IntVar(&parameters.port, "port", 8080, "A count of random letters")

	return cmd
}

func Run(cmd *cobra.Command, args []string) {
	err := configureLogging()
	if err != nil {
		logrus.Warn(err)
	}

	clientset, err := kube.GetKubernetesClient()
	if err != nil {
		logrus.Fatal(err)
	}

	// creating a new inctance of the conroller
	c, _ := controller.NewController(clientset)
	stop := make(chan struct{})
	defer close(stop)
	controller.RunController(c, stop)

	logrus.Infof("pods count is: %d", len(c.PodsIndexer.List()))

	//set timer to check every second uf there is an evict event
	ticker := time.NewTicker(time.Duration(1) * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				handler.Run(c)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", health(c))
	go http.ListenAndServe(":"+strconv.Itoa(parameters.port), mux)

	// Wait forever
	select {}

}

func health(c *controller.Controller) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if c.IsInitilize {
			fmt.Fprintf(w, "Controller is initilized successfuly")
			return
		}
		http.Error(w, "Controller is no initilized", http.StatusInternalServerError)
	}

}

func configureLogging() error {
	switch constants.LogstashHostName {
	case "":
		logrus.Infof("lgstash not configured")
	default:
		conn, err := gas.Dial("tcp", fmt.Sprintf("%s:%s", constants.LogstashHostName, constants.LogstashPort))
		if err != nil {
			log.Fatal(err)
		}
		hook := logrustash.New(conn, logrustash.DefaultFormatter(logrus.Fields{"AppName": "SpotEvictHandler"}))
		logrus.AddHook(hook)

	}
	return nil
}
