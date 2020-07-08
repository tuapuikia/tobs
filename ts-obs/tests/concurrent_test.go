package tests

import (
	"testing"
	"time"

	"ts-obs/cmd"
)

func changeRelease(t testing.TB) {
	if RELEASE_NAME == "test1" {
		RELEASE_NAME = "test2"
	} else if RELEASE_NAME == "test2" {
		RELEASE_NAME = "test1"
	} else {
		t.Fatalf("Unexpected release name %v", RELEASE_NAME)
	}

	if NAMESPACE == "test1" {
		NAMESPACE = "test2"
	} else if NAMESPACE == "test2" {
		NAMESPACE = "test1"
	} else {
		t.Fatalf("Unexpected namespace %v", RELEASE_NAME)
	}
}

func TestConcurrent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent tests")
	}

	testUninstall(t, "", "", false)

	oldname := RELEASE_NAME
	oldspace := NAMESPACE

	RELEASE_NAME = "test1"
	NAMESPACE = "test1"

	testInstall(t, "", "", "")
	changeRelease(t)
	testInstall(t, "", "", "")

	TestGrafana(t)
	TestMetrics(t)
	TestPortForward(t)
	TestPrometheus(t)
	TestTimescale(t)

	changeRelease(t)
	TestGrafana(t)
	TestMetrics(t)
	TestPortForward(t)
	TestPrometheus(t)
	TestTimescale(t)

	testUninstall(t, "", "", false)
	changeRelease(t)
	testUninstall(t, "", "", false)

	RELEASE_NAME = oldname
	NAMESPACE = oldspace

	testInstall(t, "", "", "")

	time.Sleep(10 * time.Second)

	t.Logf("Waiting for pods to initialize...")
	pods, err := cmd.KubeGetAllPods(NAMESPACE, RELEASE_NAME)
	if err != nil {
		t.Logf("Error getting all pods")
		t.Fatal(err)
	}

	for _, pod := range pods {
		err = cmd.KubeWaitOnPod(NAMESPACE, pod.Name)
		if err != nil {
			t.Logf("Error while waiting on pod")
			t.Fatal(err)
		}
	}

	time.Sleep(30 * time.Second)
}
