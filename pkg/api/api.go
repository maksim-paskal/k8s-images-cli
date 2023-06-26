/*
Copyright paskal.maksim@gmail.com
Licensed under the Apache License, Version 2.0 (the "License")
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package api

import (
	"context"
	"regexp"

	"github.com/maksim-paskal/k8s-images-cli/pkg/config"
	"github.com/maksim-paskal/k8s-images-cli/pkg/imageignore"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func getKubernetesClient(kubeConfigFile string) (*kubernetes.Clientset, error) {
	var (
		kubeconfig *rest.Config
		err        error
	)

	if len(kubeConfigFile) > 0 {
		kubeconfig, err = clientcmd.BuildConfigFromFlags("", kubeConfigFile)
		if err != nil {
			return nil, errors.Wrap(err, "clientcmd.BuildConfigFromFlags")
		}
	} else {
		kubeconfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, errors.Wrap(err, "rest.InClusterConfig")
		}
	}

	clientset, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		return nil, errors.Wrap(err, "kubernetes.NewForConfig")
	}

	return clientset, nil
}

func logImage(logEntry *log.Entry, image *Image, imageignore *imageignore.Type) {
	log := logEntry.WithFields(log.Fields{
		"ImageArch":       image.ImageArch,
		"ImagePullPolicy": image.ImagePullPolicy,
	})

	log.Debug(image.ImageName)

	printInfo := false

	if len(*config.Get().Image) > 0 {
		isMatch, err := regexp.MatchString(*config.Get().Image, image.ImageName)
		if err != nil {
			log.WithError(err).Error("error in regexp.MatchString")
		}

		printInfo = isMatch
	}

	if len(*config.Get().ImagePullPolicy) > 0 && *config.Get().ImagePullPolicy == image.ImagePullPolicy {
		printInfo = true
	}

	if len(*config.Get().ImageArch) > 0 && *config.Get().ImageArch == image.ImageArch {
		printInfo = true
	}

	if printInfo && !imageignore.Match(image.ImageName) {
		log.Info(image.ImageName)
	}
}

type Image struct {
	ImageName       string
	ImagePullPolicy string
	PodName         string
	ImageArch       string
}

func GetPodsImages(ctx context.Context, kubeConfigFile string, imageignore *imageignore.Type) (map[string]*Image, error) { //nolint:lll
	images := make(map[string]*Image)

	clientset, err := getKubernetesClient(kubeConfigFile)
	if err != nil {
		return nil, errors.Wrap(err, "error creating kubernetes client")
	}

	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "error get nodes")
	}

	nodeArch := make(map[string]string)
	for _, node := range nodes.Items {
		nodeArch[node.Name] = node.Labels["kubernetes.io/arch"]
	}

	pods, err := clientset.CoreV1().Pods(*config.Get().GetNamespace()).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "error get pods")
	}

	for _, pod := range pods.Items {
		log := log.WithFields(log.Fields{
			"kubeConfigFile": kubeConfigFile,
			"pod":            pod.Name,
			"namespace":      pod.Namespace,
		})

		image := Image{
			PodName:   pod.Name,
			ImageArch: nodeArch[pod.Spec.NodeName],
		}

		for _, initContainer := range pod.Spec.InitContainers {
			image.ImageName = initContainer.Image
			image.ImagePullPolicy = string(initContainer.ImagePullPolicy)

			logImage(log, &image, imageignore)

			if _, ok := images[initContainer.Image]; !ok {
				images[initContainer.Image] = &Image{
					PodName: pod.Name,
				}
			}
		}

		for _, container := range pod.Spec.Containers {
			image.ImageName = container.Image
			image.ImagePullPolicy = string(container.ImagePullPolicy)

			logImage(log, &image, imageignore)

			if _, ok := images[container.Image]; !ok {
				images[container.Image] = &image
			}
		}
	}

	return images, nil
}
