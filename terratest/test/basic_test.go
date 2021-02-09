package test

import (
	"fmt"
        "path/filepath"
        "strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
        "github.com/gruntwork-io/terratest/modules/k8s"
        "github.com/gruntwork-io/terratest/modules/random"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
)

const mallformedPacket = ";; Warning: Message parser reports malformed message packet."

func TestBasicExample(t *testing.T) {
        t.Parallel()

	var coreDNSPods []corev1.Pod

        // Path to the Kubernetes resource config we will test
        kubeResourcePath, err := filepath.Abs("../example/dnsendpoint.yaml")
        require.NoError(t, err)
	brokenEndpoint, err := filepath.Abs("../example/dnsendpoint_broken.yaml")
        require.NoError(t, err)

        // To ensure we can reuse the resource config on the same cluster to test different scenarios, we setup a unique
        // namespace for the resources for this test.
        // Note that namespaces must be lowercase.
        namespaceName := fmt.Sprintf("coredns-test-%s", strings.ToLower(random.UniqueId()))

        options := k8s.NewKubectlOptions("", "", namespaceName)
	mainNsOptions := k8s.NewKubectlOptions("", "", "coredns")
	podFilter := metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=coredns",
	}

        k8s.CreateNamespace(t, options, namespaceName)

        defer k8s.DeleteNamespace(t, options, namespaceName)

        defer k8s.KubectlDelete(t, options, kubeResourcePath)

        k8s.KubectlApply(t, options, kubeResourcePath)

	k8s.WaitUntilNumPodsCreated(t, mainNsOptions, podFilter, 1, 60, 1*time.Second)

	coreDNSPods = k8s.ListPods(t, mainNsOptions, podFilter)

	for _, pod := range coreDNSPods {
		k8s.WaitUntilPodAvailable(t, mainNsOptions, pod.Name, 60, 1*time.Second)
	}
	actualIP, _ := Dig(t, "localhost", 1053, "host1.example.org")

	assert.Contains(t, actualIP, "1.2.3.4")

	// check for NODATA replay on non labeled endpoints
	emptyIP, _ := Dig(t, "localhost", 1053, "host3.example.org")
	assert.NotContains(t, emptyIP,"1.2.3.4" )

	// Validate artificial(broken) DNS doesn't break CoreDNS
	k8s.KubectlApply(t, options, brokenEndpoint)
	brokenData1, _ := Dig(t, "localhost", 1053, "broken1.example.org")
	assert.Contains(t, brokenData1, mallformedPacket)
	brokenData2, _ := Dig(t, "localhost", 1053, "broken2.example.org")
	assert.Contains(t, brokenData2, mallformedPacket)
	// We still able to get "healthy" records
	currentIP, _ := Dig(t, "localhost", 1053, "host1.example.org")
	assert.Contains(t, currentIP, "1.2.3.4")
}
