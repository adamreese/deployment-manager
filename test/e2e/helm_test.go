package e2e

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

var chart = "gs://areese-charts/replicatedservice-3.tgz"

func TestHelm(t *testing.T) {
	kube := NewKubeContext()
	helm := NewHelmContext(t)

	t.Log(kube.CurrentContext())
	t.Log(kube.Cluster())
	t.Log(kube.Server())

	if !kube.Running() {
		t.Fatal("Not connected to kubernetes")
	}

	if !helmRunning(helm) {
		t.Fatal("Helm is not installed")
	}

	helm.Host = fmt.Sprintf("%s%s", kube.Server(), apiProxy)
	t.Logf("Using host: %v", helm.Host)

	t.Log("Executing deployment list")
	helm.Run("deployment", "list")

	deploymentName := genName()

	t.Log("Executing deploy")
	helm.Run("deploy", "--name", deploymentName, chart)

	t.Log("Executing deployment list")
	helm.Run("deployment", "list")

	t.Log("Executing deployment delete")
	helm.Run("deployment", "delete", deploymentName)
}

func genName() string {
	return fmt.Sprintf("%d", rand.Uint32())
}

func helmRunning(h *HelmContext) bool {
	out := h.Run("server", "status").Stdout()
	return strings.Count(out, "Running") == 5
}
