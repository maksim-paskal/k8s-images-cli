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
	"flag"
	"os"
)

type AppConfig struct {
	Version         string
	showVersion     *bool
	KubeConfigFile  *string
	ImageIgnoreFile *string
	LogLevel        *string
	Namespace       *string
	Image           *string
	ImagePullPolicy *string
}

//nolint:gochecknoglobals
var appConfig = &AppConfig{
	Version:         gitVersion,
	showVersion:     flag.Bool("version", false, "show version"),
	Namespace:       flag.String("namespace", "", "fileter by namespace"),
	LogLevel:        flag.String("logLevel", "INFO", "log level"),
	KubeConfigFile:  flag.String("kubeconfig", os.Getenv("KUBECONFIG"), "kubeconfig path comma separated"),
	ImageIgnoreFile: flag.String("imageignorefile", ".imageignore", "ignore image file"),
	Image:           flag.String("image", "", "pod info by image"),
	ImagePullPolicy: flag.String("imagepoolpolicy", "", "pod info by ImagePullPolicy (Always or IfNotPresent or Never)"),
}
