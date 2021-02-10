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
package main

import (
	"context"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func getPodsImages(clientset *kubernetes.Clientset) (map[string]v1.Pod, error) {
	images := make(map[string]v1.Pod)

	pods, err := clientset.CoreV1().Pods(*appConfig.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return images, errors.Wrap(err, "error get pods")
	}

	for _, pod := range pods.Items {
		log := log.WithFields(log.Fields{
			"pod":       pod.Name,
			"namespace": pod.Namespace,
		})

		for _, initContainer := range pod.Spec.InitContainers {
			log.Debug(initContainer.Image)

			if _, ok := images[initContainer.Image]; !ok {
				images[initContainer.Image] = pod
			}
		}

		for _, container := range pod.Spec.Containers {
			log.Debug(container.Image)

			if _, ok := images[container.Image]; !ok {
				images[container.Image] = pod
			}
		}
	}

	return images, nil
}
