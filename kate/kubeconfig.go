package kate

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/nospof/secretretriever/tools"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func GetKubeConfig() *kubernetes.Clientset {
	inOut := os.Getenv("IN_K8S")

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	if inOut == "IN" {
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	}
	kate, err := kubernetes.NewForConfig(config)
	tools.CheckIfError(err)
	return kate
}
