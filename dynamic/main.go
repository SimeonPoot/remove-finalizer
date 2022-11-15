package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	cm "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
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
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	certRes := schema.GroupVersionResource{Group: "cert-manager.io", Version: "v1", Resource: "certificates"}

	finalizer := "venafi-cleanup-operator"

	fmt.Println("hello")

	cert := cm.Certificate{}
	cert.Name = "example-com-2"
	fmt.Println(cert)

	result, err := client.Resource(certRes).Namespace(v1.NamespaceDefault).Get(context.TODO(), "example-com-2", metav1.GetOptions{})
	if err != nil {
		panic(fmt.Errorf("failed to get latest version of Deployment: %v", err))
	}

	fmt.Println("before: ", result.GetFinalizers())

	controllerutil.RemoveFinalizer(result, finalizer)

	fmt.Println("after: ", result.GetFinalizers())
	_, err = client.Resource(certRes).Namespace(v1.NamespaceDefault).Update(context.TODO(), result, metav1.UpdateOptions{})
	if err != nil {
		panic(fmt.Errorf("failed to get latest version of Deployment: %v", err))
	}

}
