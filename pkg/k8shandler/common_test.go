package k8shandler

import (
	"reflect"
	"testing"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	api "github.com/openshift/elasticsearch-operator/pkg/apis/logging/v1"
)

var (
	commonCpuValue = resource.MustParse("500m")
	commonMemValue = resource.MustParse("2Gi")

	nodeCpuValue = resource.MustParse("600m")
	nodeMemValue = resource.MustParse("3Gi")

	defaultTestCpuLimit   = resource.MustParse(defaultCPULimit)
	defaultTestCpuRequest = resource.MustParse(defaultCPURequest)
	defaultTestMemLimit   = resource.MustParse(defaultMemoryLimit)
	defaultTestMemRequest = resource.MustParse(defaultMemoryRequest)
)

/*
  Resource scenarios:
  1. common has a limit and request
     node sets new request and limit
     -> use node settings

  2. common has limit and request
     node has no settings
     -> use common settings

  3. common has no settings
     node has no settings
     -> use defaults

  4. common has limit and request
     node sets request only
     -> use common limit and node request

  5. common has limit and request
     node sets limt only
     -> use common request and node limit

  6. common has request only
     node has no settings
     -> request and limit set to be same and use common settings

  7. common has limit only
     node has no settings
     -> request and limit set to be same and use common settings

  8. common has no settings
     node only has request
     -> request and limit set to be same and use node settings

  9. common has no settings
     node only has limit
     -> request and limit set to be same and use node settings

 10. common has limit
     node has request
     -> use common limit and node request

 11. common has request
     node has limit
     -> use common request and node limit
*/

// 1
func TestResourcesCommonAndNodeDefined(t *testing.T) {

	commonRequirements := buildResource(
		commonCpuValue,
		commonCpuValue,
		commonMemValue,
		commonMemValue,
	)

	nodeRequirements := buildResource(
		nodeCpuValue,
		nodeCpuValue,
		nodeMemValue,
		nodeMemValue,
	)

	expected := buildResource(
		nodeCpuValue,
		nodeCpuValue,
		nodeMemValue,
		nodeMemValue,
	)

	actual := newResourceRequirements(nodeRequirements, commonRequirements)

	if !areResourcesSame(actual, expected) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

// 2
func TestResourcesNoNodeDefined(t *testing.T) {

	commonRequirements := buildResource(
		commonCpuValue,
		commonCpuValue,
		nodeMemValue,
		nodeMemValue,
	)

	nodeRequirements := v1.ResourceRequirements{}

	expected := buildResource(
		commonCpuValue,
		commonCpuValue,
		nodeMemValue,
		nodeMemValue,
	)

	actual := newResourceRequirements(nodeRequirements, commonRequirements)

	if !areResourcesSame(actual, expected) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

// 3
func TestResourcesNoCommonNoNodeDefined(t *testing.T) {

	commonRequirements := v1.ResourceRequirements{}

	nodeRequirements := v1.ResourceRequirements{}

	expected := buildResource(
		defaultTestCpuLimit,
		defaultTestCpuRequest,
		defaultTestMemLimit,
		defaultTestMemRequest,
	)

	actual := newResourceRequirements(nodeRequirements, commonRequirements)

	if !areResourcesSame(actual, expected) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

// 4
func TestResourcesCommonAndNodeRequestDefined(t *testing.T) {

	commonRequirements := buildResource(
		commonCpuValue,
		commonCpuValue,
		commonMemValue,
		commonMemValue,
	)

	nodeRequirements := buildResourceOnlyRequests(
		nodeCpuValue,
		nodeMemValue,
	)

	expected := buildResource(
		commonCpuValue,
		nodeCpuValue,
		commonMemValue,
		nodeMemValue,
	)

	actual := newResourceRequirements(nodeRequirements, commonRequirements)

	if !areResourcesSame(actual, expected) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

// 5
func TestResourcesCommonAndNodeLimitDefined(t *testing.T) {

	commonRequirements := buildResource(
		commonCpuValue,
		commonCpuValue,
		commonMemValue,
		commonMemValue,
	)

	nodeRequirements := buildResourceOnlyLimits(
		nodeCpuValue,
		nodeMemValue,
	)

	expected := buildResource(
		nodeCpuValue,
		commonCpuValue,
		nodeMemValue,
		commonMemValue,
	)

	actual := newResourceRequirements(nodeRequirements, commonRequirements)

	if !areResourcesSame(actual, expected) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

// 6
func TestResourcesCommonRequestAndNoNodeDefined(t *testing.T) {

	commonRequirements := buildResourceOnlyRequests(
		commonCpuValue,
		commonMemValue,
	)

	nodeRequirements := v1.ResourceRequirements{}

	expected := buildResource(
		commonCpuValue,
		commonCpuValue,
		commonMemValue,
		commonMemValue,
	)

	actual := newResourceRequirements(nodeRequirements, commonRequirements)

	if !areResourcesSame(actual, expected) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

// 7
func TestResourcesCommonLimitAndNoNodeDefined(t *testing.T) {

	commonRequirements := buildResourceOnlyLimits(
		commonCpuValue,
		commonMemValue,
	)

	nodeRequirements := v1.ResourceRequirements{}

	expected := buildResource(
		commonCpuValue,
		commonCpuValue,
		commonMemValue,
		commonMemValue,
	)

	actual := newResourceRequirements(nodeRequirements, commonRequirements)

	if !areResourcesSame(actual, expected) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

// 8
func TestResourcesNoCommonAndNodeRequestDefined(t *testing.T) {

	commonRequirements := v1.ResourceRequirements{}

	nodeRequirements := buildResourceOnlyRequests(
		nodeCpuValue,
		nodeMemValue,
	)

	expected := buildResource(
		nodeCpuValue,
		nodeCpuValue,
		nodeMemValue,
		nodeMemValue,
	)

	actual := newResourceRequirements(nodeRequirements, commonRequirements)

	if !areResourcesSame(actual, expected) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

// 9
func TestResourcesNoCommonAndNodeLimitDefined(t *testing.T) {

	commonRequirements := v1.ResourceRequirements{}

	nodeRequirements := buildResourceOnlyLimits(
		nodeCpuValue,
		nodeMemValue,
	)

	expected := buildResource(
		nodeCpuValue,
		nodeCpuValue,
		nodeMemValue,
		nodeMemValue,
	)

	actual := newResourceRequirements(nodeRequirements, commonRequirements)

	if !areResourcesSame(actual, expected) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

// 10
func TestResourcesCommonLimitAndNodeResourceDefined(t *testing.T) {

	commonRequirements := buildResourceOnlyLimits(
		commonCpuValue,
		commonMemValue,
	)

	nodeRequirements := buildResourceOnlyRequests(
		nodeCpuValue,
		nodeMemValue,
	)

	expected := buildResource(
		commonCpuValue,
		nodeCpuValue,
		commonMemValue,
		nodeMemValue,
	)

	actual := newResourceRequirements(nodeRequirements, commonRequirements)

	if !areResourcesSame(actual, expected) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

// 11
func TestResourcesCommonResourceAndNodeLimitDefined(t *testing.T) {

	commonRequirements := buildResourceOnlyRequests(
		commonCpuValue,
		commonMemValue,
	)

	nodeRequirements := buildResourceOnlyLimits(
		nodeCpuValue,
		nodeMemValue,
	)

	expected := buildResource(
		nodeCpuValue,
		commonCpuValue,
		nodeMemValue,
		commonMemValue,
	)

	actual := newResourceRequirements(nodeRequirements, commonRequirements)

	if !areResourcesSame(actual, expected) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func TestProxyContainerResourcesDefined(t *testing.T) {

	imageName := "openshift/oauth-proxy:1.1"
	clusterName := "elasticsearch"

	expectedCPU := resource.MustParse("100m")
	expectedMemory := resource.MustParse("64Mi")

	proxyContainer, err := newProxyContainer(imageName, clusterName)
	if err != nil {
		t.Errorf("Failed to populate Proxy container: %v", err)
	}

	if memoryLimit, ok := proxyContainer.Resources.Limits["memory"]; ok {
		if memoryLimit.Cmp(expectedMemory) != 0 {
			t.Errorf("Expected CPU limit %s but got %s", expectedMemory.String(), memoryLimit.String())
		}
	} else {
		t.Errorf("Proxy container is missing CPU limit. Expected limit %s", expectedMemory.String())
	}

	if cpuRequest, ok := proxyContainer.Resources.Requests["cpu"]; ok {
		if cpuRequest.Cmp(expectedCPU) != 0 {
			t.Errorf("Expected CPU request %s but got %s", expectedCPU.String(), cpuRequest.String())
		}
	} else {
		t.Errorf("Proxy container is missing CPU request. Expected request %s", expectedCPU.String())
	}

	if memoryLimit, ok := proxyContainer.Resources.Limits["memory"]; ok {
		if memoryLimit.Cmp(expectedMemory) != 0 {
			t.Errorf("Expected memory limit %s but got %s", expectedMemory.String(), memoryLimit.String())
		}
	} else {
		t.Errorf("Proxy container is missing memory limit. Expected limit %s", expectedMemory.String())
	}
}

func TestPodSpecHasTaintTolerations(t *testing.T) {

	expectedTolerations := []v1.Toleration{
		v1.Toleration{
			Key:      "node.kubernetes.io/disk-pressure",
			Operator: v1.TolerationOpExists,
			Effect:   v1.TaintEffectNoSchedule,
		},
	}

	podTemplateSpec := newPodTemplateSpec("test-node-name", "test-cluster-name", "test-namespace-name", api.ElasticsearchNode{}, api.ElasticsearchNodeSpec{}, map[string]string{}, map[api.ElasticsearchNodeRole]bool{}, nil)

	if !reflect.DeepEqual(podTemplateSpec.Spec.Tolerations, expectedTolerations) {
		t.Errorf("Exp. the tolerations to be %q but was %q", expectedTolerations, podTemplateSpec.Spec.Tolerations)
	}
}

func buildResource(cpuLimit, cpuRequest, memLimit, memRequest resource.Quantity) v1.ResourceRequirements {
	return v1.ResourceRequirements{
		Limits: v1.ResourceList{
			v1.ResourceCPU:    cpuLimit,
			v1.ResourceMemory: memLimit,
		},
		Requests: v1.ResourceList{
			v1.ResourceCPU:    cpuRequest,
			v1.ResourceMemory: memRequest,
		},
	}
}

func buildResourceOnlyRequests(cpuRequest, memRequest resource.Quantity) v1.ResourceRequirements {
	return v1.ResourceRequirements{
		Requests: v1.ResourceList{
			v1.ResourceCPU:    cpuRequest,
			v1.ResourceMemory: memRequest,
		},
	}
}

func buildResourceOnlyLimits(cpuLimit, memLimit resource.Quantity) v1.ResourceRequirements {
	return v1.ResourceRequirements{
		Limits: v1.ResourceList{
			v1.ResourceCPU:    cpuLimit,
			v1.ResourceMemory: memLimit,
		},
	}
}

func areResourcesSame(lhs, rhs v1.ResourceRequirements) bool {

	if !areQuantitiesSame(lhs.Limits.Cpu(), rhs.Limits.Cpu()) {
		return false
	}

	if !areQuantitiesSame(lhs.Requests.Cpu(), rhs.Requests.Cpu()) {
		return false
	}

	if !areQuantitiesSame(lhs.Limits.Memory(), rhs.Limits.Memory()) {
		return false
	}

	if !areQuantitiesSame(lhs.Requests.Memory(), rhs.Requests.Memory()) {
		return false
	}

	return true
}

func areQuantitiesSame(lhs, rhs *resource.Quantity) bool {
	return lhs.Cmp(*rhs) == 0
}
