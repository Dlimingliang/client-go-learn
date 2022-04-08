package main

import (
	"flag"
	"fmt"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	apischeme "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
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

	//实例化dynamicClient
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(any(err.Error()))
	}

	//获取非结构化数据
	gvr := apischeme.GroupVersionResource{
		Version:  "v1",
		Resource: "pods",
	}
	unstructObj, err := dynamicClient.Resource(gvr).Namespace("kube-system").List(metav1.ListOptions{Limit: 100})
	if err != nil {
		panic(any(err.Error()))
	}

	//将非结构化数据转换为结构化数据
	podList := &corev1.PodList{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructObj.UnstructuredContent(), podList)
	if err != nil {
		panic(any(err.Error()))
	}

	// 表头
	fmt.Printf("namespace\t status\t\t name\n")
	for _, item := range podList.Items {
		fmt.Printf("%v\t %v\t %v\n", item.Name, item.Status.Phase, item.Name)
	}
}
