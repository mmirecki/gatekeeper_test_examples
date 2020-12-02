package main

import (
	"flag"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	master = flag.String("master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
)

func main() {
	var kubeconfig *string

	kubeconfigEnvVar := os.Getenv("KUBECONFIG")
	fmt.Println(kubeconfigEnvVar)
	kubeconfig = &kubeconfigEnvVar

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	createAM(config)
	listAM(config)
}

func listAM(config *rest.Config) {

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	amGVR := schema.GroupVersionResource{
		Group:    "mutations.gatekeeper.sh",
		Version:  "v1alpha1",
		Resource: "assignmetadata",
	}
	namClient := dynClient.Resource(amGVR)

	crds, err := namClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("===FOR====")
	for _, crd := range crds.Items {
		fmt.Println("CRD: ", crd)
	}
	fmt.Println("===DONE====")

	//fmt.Println(dynClient, namClient)
}

func createAM(config *rest.Config) {
	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	amGVR := schema.GroupVersionResource{
		Group:    "mutations.gatekeeper.sh",
		Version:  "v1alpha1",
		Resource: "assignmetadata",
	}
	amClient := dynClient.Resource(amGVR)

	amDefintion := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "mutations.gatekeeper.sh/v1alpha1",
			"kind":       "AssignMetadata",
			"metadata": map[string]interface{}{
				"name": "assignMetadata1",
			},
			"spec": map[string]interface{}{
				"location": "metadata.labels.x",
				"parameters": map[string]interface{}{
					"value": "testValue",
				},
			},
		},
	}

	result, err := amClient.Create(amDefintion, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("ERR: %+v \n\n", err)
		//panic(err.Error())
	}
	fmt.Println("RESULT: ", result)

}
