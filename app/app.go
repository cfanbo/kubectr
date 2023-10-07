package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/klog/v2"
)

var (
	ErrEmpty = errors.New("empty")
)

type containerType string

var (
	initContainer       containerType = "initContainer"
	normalContainer     containerType = "container"
	ephemeralContainers containerType = "ephemeralContainer"
)

type ContainerState string

// CtrOptions app options
type CtrOptions struct {
	configFlags *genericclioptions.ConfigFlags

	rawConfig              api.Config
	userSpecialedPodName   string
	userSpecifiedNamespace string
	version                bool
}

// Complete prepare for ctr
func (o *CtrOptions) Complete(cmd *cobra.Command, args []string) error {
	var err error
	o.rawConfig, err = o.configFlags.ToRawKubeConfigLoader().RawConfig()
	if err != nil {
		return err
	}

	// pod name
	o.userSpecialedPodName = args[0]

	// namespace
	o.userSpecifiedNamespace, err = cmd.Flags().GetString("namespace")
	if err != nil {
		return err
	}

	if o.userSpecifiedNamespace == "" {
		currentCtx := o.rawConfig.Contexts[o.rawConfig.CurrentContext]
		if currentCtx == nil {
			return errors.New("failed to read kubeconfig configuration")
		}
		o.userSpecifiedNamespace = currentCtx.Namespace
	}
	if o.userSpecifiedNamespace == "" {
		o.userSpecifiedNamespace = "default"
	}

	return nil
}

// Run execute
func (o *CtrOptions) Run() {
	namespace := o.userSpecifiedNamespace
	containers, err := o.containersForPod()
	if err != nil {
		if errors.Is(err, ErrEmpty) {
			fmt.Printf("No container found in %s namespace.\n", namespace)
			return
		}
		klog.V(2).ErrorS(err, "fetch all containers in "+namespace+" namespace")
		fmt.Println(err)
		return
	}

	render(containers)
}

func (o *CtrOptions) containersForPod() ([]Container, error) {
	restConfig, _ := o.configFlags.ToRESTConfig()
	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	pod, err := clientSet.CoreV1().Pods(o.userSpecifiedNamespace).Get(context.Background(), o.userSpecialedPodName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	containerTotal := len(pod.Spec.InitContainers) + len(pod.Spec.Containers) + len(pod.Spec.EphemeralContainers)
	if containerTotal == 0 {
		return nil, ErrEmpty
	}

	// container status
	statusMap := make(map[string]corev1.ContainerStatus, len(pod.Status.InitContainerStatuses)+len(pod.Status.ContainerStatuses))
	for _, containerStatus := range pod.Status.InitContainerStatuses {
		statusMap[containerStatus.Name] = containerStatus
	}
	for _, containerStatus := range pod.Status.ContainerStatuses {
		statusMap[containerStatus.Name] = containerStatus
	}
	for _, containerStatus := range pod.Status.EphemeralContainerStatuses {
		statusMap[containerStatus.Name] = containerStatus
	}

	containerList := make([]Container, 0, containerTotal)
	// initContainers
	for _, container := range pod.Spec.InitContainers {
		containerList = append(containerList, Container{
			Type:            initContainer,
			Name:            container.Name,
			Image:           container.Image,
			ImagePullPolicy: string(container.ImagePullPolicy),
			Ports:           parsePort(container.Ports),
			Status:          statusMap[container.Name],
		})
	}

	// containers
	for _, container := range pod.Spec.Containers {
		containerList = append(containerList, Container{
			Type:            normalContainer,
			Name:            container.Name,
			Image:           container.Image,
			ImagePullPolicy: string(container.ImagePullPolicy),
			Ports:           parsePort(container.Ports),
			Status:          statusMap[container.Name],
		})
	}

	// EphemeralContainers
	for _, ephemeral := range pod.Spec.EphemeralContainers {
		container := ephemeral.EphemeralContainerCommon
		containerList = append(containerList, Container{
			Type:            ephemeralContainers,
			Name:            container.Name,
			Image:           container.Image,
			ImagePullPolicy: string(container.ImagePullPolicy),
			Ports:           parsePort(container.Ports),
			Status:          statusMap[container.Name],
		})
	}

	return containerList, nil
}
