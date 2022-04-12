package main

import (
	"flag"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	k8sv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const (
	Namespace = "my-test"
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
	namespaceClient := clientSet.CoreV1().Namespaces()

	//createNamespace(namespaceClient)
	updateNamespace(namespaceClient)

	//deleteNamespace(namespaceClient)

	//listNamespace(namespaceClient)
	//getNamespace(namespaceClient)
}

func getNamespace(namespaceClient k8sv1.NamespaceInterface) {
	namespace, err := namespaceClient.Get(Namespace, metav1.GetOptions{})
	if err != nil {
		panic(any(err.Error()))
	}
	fmt.Println(namespace)
}

func listNamespace(namespaceClient k8sv1.NamespaceInterface) {
	list, err := namespaceClient.List(metav1.ListOptions{})
	if err != nil {
		panic(any(err.Error()))
	}
	for _, item := range list.Items {
		fmt.Println(item.Name)
	}
}

func createNamespace(namespaceClient k8sv1.NamespaceInterface) {
	namespace, err := namespaceClient.Create(&corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: Namespace,
			Labels: map[string]string{
				"cluster": "123",
			},
			Annotations: map[string]string{
				"cluster": "123",
			},
		},
	})
	if err != nil {
		panic(any(err.Error()))
	}
	fmt.Printf("Create namespace %s \n", namespace.Name)
}

func deleteNamespace(namespaceClient k8sv1.NamespaceInterface) {
	deletePolicy := metav1.DeletePropagationForeground
	err := namespaceClient.Delete(Namespace, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		panic(any(err.Error()))
	}
	fmt.Printf("Delete namespace %s \n", Namespace)
}

func updateNamespace(namespaceClient k8sv1.NamespaceInterface) {
	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		namespace, err := namespaceClient.Get(Namespace, metav1.GetOptions{})
		if err != nil {
			panic(any(err.Error()))
		}
		namespace.Labels = map[string]string{
			"cluster": "789",
		}
		_, err = namespaceClient.Update(namespace)
		if err != nil {
			panic(any(err.Error()))
		}
		return err
	})
	if err != nil {
		panic(any(err.Error()))
	}
}
