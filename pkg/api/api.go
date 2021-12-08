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
	corev1 "k8s.io/api/core/v1"
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

func logContainerInfo(logEntry *log.Entry, container corev1.Container, imageignore *imageignore.Type) {
	log := logEntry.WithFields(log.Fields{
		"container":       container.Name,
		"ImagePullPolicy": container.ImagePullPolicy,
	})

	log.Debug(container.Image)

	printInfo := false

	if len(*config.Get().Image) > 0 {
		isMatch, err := regexp.MatchString(*config.Get().Image, container.Image)
		if err != nil {
			log.WithError(err).Error("error in regexp.MatchString")
		}

		printInfo = isMatch
	}

	if len(*config.Get().ImagePullPolicy) > 0 && *config.Get().ImagePullPolicy == string(container.ImagePullPolicy) {
		printInfo = true
	}

	if printInfo && !imageignore.Match(container.Image) {
		log.Info(container.Image)
	}
}

func GetPodsImages(kubeConfigFile string, imageignore *imageignore.Type) (map[string]corev1.Pod, error) {
	images := make(map[string]corev1.Pod)

	clientset, err := getKubernetesClient(kubeConfigFile)
	if err != nil {
		return images, errors.Wrap(err, "error creating kubernetes client")
	}

	pods, err := clientset.CoreV1().Pods(*config.Get().Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return images, errors.Wrap(err, "error get pods")
	}

	for _, pod := range pods.Items {
		log := log.WithFields(log.Fields{
			"kubeConfigFile": kubeConfigFile,
			"pod":            pod.Name,
			"namespace":      pod.Namespace,
		})

		for _, initContainer := range pod.Spec.InitContainers {
			logContainerInfo(log, initContainer, imageignore)

			if _, ok := images[initContainer.Image]; !ok {
				images[initContainer.Image] = pod
			}
		}

		for _, container := range pod.Spec.Containers {
			logContainerInfo(log, container, imageignore)

			if _, ok := images[container.Image]; !ok {
				images[container.Image] = pod
			}
		}
	}

	return images, nil
}
