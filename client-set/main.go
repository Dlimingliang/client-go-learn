package main

import (
	"flag"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/utils/pointer"
	"path/filepath"
)

const (
	Clean          = "clean"
	Create         = "create"
	Namespace      = "client-set"
	DeploymentName = "client-set-deployment"
	ServiceName    = "client-test-service"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, "config"), "absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	fmt.Println("kubeconfig的值: ", *kubeconfig)

	operate := flag.String("operate", Clean, fmt.Sprintf("operate type : %s or %s", Create, Clean))
	flag.Parse()
	fmt.Printf("operation is %v\n", *operate)

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

	if *operate == Create {
		createNamespace(clientSet)
		createDeployment(clientSet)
		createService(clientSet)
	} else {
		clean(clientSet)
	}
}

func clean(clientSet *kubernetes.Clientset) {
	emptyDeleteOptions := metav1.DeleteOptions{}
	//删除service
	if err := clientSet.CoreV1().Services(Namespace).Delete(ServiceName, &emptyDeleteOptions); err != nil {
		panic(any(err.Error()))
	}
	//删除deployment
	if err := clientSet.AppsV1().Deployments(Namespace).Delete(DeploymentName, &emptyDeleteOptions); err != nil {
		panic(any(err.Error()))
	}
	//删除namespace
	if err := clientSet.CoreV1().Namespaces().Delete(Namespace, &emptyDeleteOptions); err != nil {
		panic(any(err.Error()))
	}
}

func createNamespace(clientSet *kubernetes.Clientset) {
	namespaceClient := clientSet.CoreV1().Namespaces()
	namespace, err := namespaceClient.Create(&corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: Namespace,
		},
	})
	if err != nil {
		panic(any(err.Error()))
	}
	fmt.Printf("Create namespace %s \n", namespace.Name)
}

func createService(clientSet *kubernetes.Clientset) {
	serviceClient := clientSet.CoreV1().Services(Namespace)
	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: ServiceName,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Port:     8080,
					NodePort: 30080,
				},
			},
			Selector: map[string]string{
				"app": "tomcat",
			},
			Type: corev1.ServiceTypeNodePort,
		},
	}
	result, err := serviceClient.Create(&service)
	if err != nil {
		panic(any(err.Error()))
	}
	fmt.Printf("Create service %s \n", result.Name)
}

func createDeployment(clientSet *kubernetes.Clientset) {
	deploymentClient := clientSet.AppsV1().Deployments(Namespace)
	deployment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: DeploymentName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "tomcat",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "tomcat",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "tomcat",
							Image: "tomcat:8.0.18-jre8",
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 8080,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							ImagePullPolicy: corev1.PullIfNotPresent,
						},
					},
				},
			},
		},
	}
	result, err := deploymentClient.Create(&deployment)
	if err != nil {
		panic(any(err.Error()))
	}
	fmt.Printf("Create deployment %s \n", result.Name)
}
