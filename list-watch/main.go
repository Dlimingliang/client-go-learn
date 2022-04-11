package main

import (
	"flag"
	"fmt"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, "config"), "absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	fmt.Println("kubeconfig的值: ", *kubeconfig)
	flag.Parse()

	//加载kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(any(err.Error()))
	}

	//实例化clientset
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(any(err.Error()))
	}

	deploymentClient := clientSet.AppsV1().Deployments(corev1.NamespaceAll)
	list, err := deploymentClient.List(metav1.ListOptions{})
	if err != nil {
		panic(any(err.Error()))
	}
	fmt.Println(list.ResourceVersion)
	for _, item := range list.Items {
		fmt.Printf("* %s (%d replicas %s)\n", item.Name, *item.Spec.Replicas, item.ResourceVersion)
	}
}
