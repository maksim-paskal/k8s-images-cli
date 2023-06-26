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
package config

import (
	"encoding/json"
	"flag"
	"os"
)

type AppConfig struct {
	KubeConfigFile  *string
	ImageIgnoreFile *string
	LogLevel        *string
	Namespace       *string
	NamespaceShort  *string
	Image           *string
	ImagePullPolicy *string
	ImageArch       *string
	ShowArch        *bool
}

//nolint:gochecknoglobals
var appConfig = &AppConfig{
	Namespace:       flag.String("namespace", "", "filter by namespace"),
	NamespaceShort:  flag.String("n", "", "filter by namespace short"),
	LogLevel:        flag.String("logLevel", "INFO", "log level"),
	KubeConfigFile:  flag.String("kubeconfig", os.Getenv("KUBECONFIG"), "kubeconfig path(s) (comma separated)"),
	ImageIgnoreFile: flag.String("imageignorefile", ".imageignore", "ignore image file"),
	Image:           flag.String("image", "", "filter by image"),
	ImagePullPolicy: flag.String("imagepoolpolicy", "", "filter by ImagePullPolicy (Always or IfNotPresent or Never)"),
	ImageArch:       flag.String("arch", "", "filter by Arch (amd64 or arm64)"),
	ShowArch:        flag.Bool("showArch", false, "show image arch"),
}

func (a *AppConfig) String() string {
	b, err := json.Marshal(a)
	if err != nil {
		return err.Error()
	}

	return string(b)
}

func (a *AppConfig) GetNamespace() *string {
	if a.NamespaceShort != nil && len(*a.NamespaceShort) > 0 {
		return a.NamespaceShort
	}

	return a.Namespace
}

func Get() *AppConfig {
	return appConfig
}

//nolint:gochecknoglobals
var gitVersion = "dev"

func GetVersion() string {
	return gitVersion
}
