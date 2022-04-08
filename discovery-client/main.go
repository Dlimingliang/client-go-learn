package main

import (
	"flag"
	"fmt"

	"path/filepath"

	apischeme "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
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

	//新建discoveryClient
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(any(err.Error()))
	}

	//获取所有分组和资源数据
	APIGroup, APIResourceListSlice, err := discoveryClient.ServerGroupsAndResources()
	if err != nil {
		panic(any(err.Error()))
	}

	// Group信息
	fmt.Printf("APIGroup :\n\n %v\n\n\n\n", APIGroup)

	for _, singleAPIResourceList := range APIResourceListSlice {
		groupVersionStr := singleAPIResourceList.GroupVersion
		gv, err := apischeme.ParseGroupVersion(groupVersionStr)
		if err != nil {
			panic(any(err.Error()))
		}
		fmt.Println("**********************************")
		fmt.Printf("GV string[%v]\nGV struct [%#v]\nresources: \n\n", groupVersionStr, gv)

		for _, resource := range singleAPIResourceList.APIResources {
			fmt.Printf("%v\n", resource.Name)
		}
	}
}
