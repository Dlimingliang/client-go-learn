package main

import (
	"flag"
	"fmt"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	fmt.Println("kubeconfig的值: ",*kubeconfig)
	flag.Parse()

	//加载kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(any(err.Error()))
	}

	config.APIPath = "api"
	//pod的group为空字符串
	config.GroupVersion = &corev1.SchemeGroupVersion
	//指定序列化工具
	config.NegotiatedSerializer = scheme.Codecs

	//构建rest客户端
	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		panic(any(err.Error()))
	}

	//返回结果
	result := &corev1.PodList{}

	err = restClient.Get().
		Namespace("kube-system").
		Resource("pods").
		//指定大小限制和序列化工具
		VersionedParams(&metav1.ListOptions{Limit: 100}, scheme.ParameterCodec).
		Do().
		Into(result)
	if err != nil {
		panic(any(err.Error()))
	}

	// 表头
	fmt.Printf("namespace\t status\t\t name\n")
	for _, item := range result.Items {
		fmt.Printf("%v\t %v\t %v\n", item.Name, item.Status.Phase, item.Name)
	}
}
