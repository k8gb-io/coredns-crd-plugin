package k8scrd

import (
	"os"

	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type configType string

const (
	inClusterConfigType configType = "inCluster"
	localConfigType     configType = "local"
)

// config factory provides k8s config based on argument.
func configFactory(ct configType) (cfg *restclient.Config, rct configType, err error) {
	rct = ct
	switch ct {
	case localConfigType:
		kubeconfig := os.Getenv("KUBECONFIG")
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	case inClusterConfigType:
		fallthrough
	default:
		cfg, err = restclient.InClusterConfig()
		rct = inClusterConfigType
	}
	return cfg, rct, err
}
