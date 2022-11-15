package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func main() {

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	finalizer := "finalizer.extensions/v1beta1"

	fmt.Println("hello")

	result, err := clientset.CoreV1().Pods(v1.NamespaceDefault).Get(context.TODO(), "nginx", metav1.GetOptions{})
	if err != nil {
		panic(fmt.Errorf("failed to get pod: %v", err))
	}

	fmt.Println("before: ", result.GetFinalizers())

	controllerutil.RemoveFinalizer(result, finalizer)

	fmt.Println("after: ", result.GetFinalizers())

	_, err = clientset.CoreV1().Pods(v1.NamespaceDefault).Update(context.TODO(), result, metav1.UpdateOptions{})
	if err != nil {
		panic(fmt.Errorf("failed to update pod: %v", err))
	}

	fmt.Println("goodbye")
}
