package k8scrd

import (
	"context"

	restclient "k8s.io/client-go/rest"

	"github.com/AbsaOSS/k8s_crd/common/k8sctrl"
	dnsendpoint "github.com/AbsaOSS/k8s_crd/extdns"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
)

// KubeController stores the current runtime configuration and cache
type KubeController struct {
	client      dnsendpoint.ExtDNSInterface
	controllers []cache.SharedIndexInformer
	labelFilter string
	hasSynced   bool
}

// RunKubeController kicks off the k8s controllers
func RunKubeController(ctx context.Context, cfg *restclient.Config, filter string) (*k8sctrl.KubeController, error) {

	err := dnsendpoint.AddToScheme(scheme.Scheme)
	if err != nil {
		return nil, err
	}

	kubeClient, err := dnsendpoint.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	ctrl := k8sctrl.NewKubeController(ctx, kubeClient, filter)

	go ctrl.Run()

	return ctrl, nil

}
