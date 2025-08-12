package main

import (
	"context"
	"fmt"

	loftclient "github.com/loft-sh/api/v4/pkg/clientset/versioned"

	managementv1 "github.com/loft-sh/api/v4/pkg/apis/management/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	// List users using controller-runtime client
	ListUsersControllerRuntime()
	// List users using regular kube client
	ListUsersKubeClient()
}

var Scheme = runtime.NewScheme()

func init() {
	_ = clientgoscheme.AddToScheme(Scheme)
	_ = managementv1.AddToScheme(Scheme)
}

func ListUsersControllerRuntime() {
	// get kube config
	restConfig, err := ctrl.GetConfig()
	if err != nil {
		panic(err)
	}

	// create kube client
	kubeClient, err := client.New(restConfig, client.Options{Scheme: Scheme})
	if err != nil {
		panic(err)
	}

	// list users
	userList := &managementv1.UserList{}
	err = kubeClient.List(context.Background(), userList, &client.ListOptions{})
	if err != nil {
		panic(err)
	}

	// print users
	for _, user := range userList.Items {
		fmt.Println(user.Name)
	}
}

func ListUsersKubeClient() {
	// get kube config
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(clientcmd.NewDefaultClientConfigLoadingRules(), &clientcmd.ConfigOverrides{})
	if clientConfig == nil {
		panic("nil clientConfig")
	}

	// create kube client
	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		panic(err)
	}

	// create loft client
	loftClient, err := loftclient.NewForConfig(restConfig)
	if err != nil {
		panic(err)
	}

	// list users
	users, err := loftClient.ManagementV1().Users().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	// print users
	for _, user := range users.Items {
		fmt.Println(user.Name)
	}
}
