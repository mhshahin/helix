package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mhshahin/helix/pkg/config"
	"github.com/mhshahin/helix/pkg/handler"
	"github.com/mhshahin/helix/pkg/kubernetes"
	"github.com/mhshahin/helix/pkg/medium"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var confFile = flag.String("config", "config", "absolute path to the kubeconfig file")
var envFile = flag.String("env", "env.yaml", "env file absolute path")

func main() {
	flag.Parse()
	config.LoadConfig(*envFile)
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", *confFile)
	if err != nil {
		if err == rest.ErrNotInCluster {
			log.Panicln("provide a kubernetes config file or deploy in a cluster")
		}
		log.Panicln(err.Error())
	}

	m := medium.NewMedium()
	h := handler.NewHandler(m)
	w := kubernetes.NewEventWatcher(kubeConfig, config.Cfg.Namespace, h.OnEvent)

	w.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	gracefulExit := func() {
		defer close(c)
		w.Stop()
		fmt.Println("Exiting...")
	}

	sig := <-c
	fmt.Printf("Received %s signal to exit.\n", sig.String())
	gracefulExit()

}
