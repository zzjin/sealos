// Copyright © 2023 sealos.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package crdb

import (
	"errors"
	"time"

	"github.com/go-logr/logr"
	"github.com/labring/crdbase/utils"
	appsv1 "k8s.io/api/apps/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	pkg_client "sigs.k8s.io/controller-runtime/pkg/client"
	pkg_manager "sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	crdBaseName    = "crdbase"
	crdBaseURL     = "crdbase.sealos.io"
	crdBasePackage = "v0.1.0"

	crdApiVersion = "v1"

	crdbaseTagKey = "crdb"
)

var (
	ErrNotFound = errors.New("notFound")
)

var (
	providerName = crdBaseName + "-" + crdApiVersion
)

var (
	PollInterval = 1 * time.Second
	MaxWait      = 10 * time.Second
)

type CRDBaseConfig struct {
	Manager        pkg_manager.Manager
	GroupVersion   schema.GroupVersion
	ServiceAccount string
	Namespace      string
}

type CRDBase struct {
	CRDBaseConfig

	log logr.Logger

	client        pkg_client.Client     // client
	clientSet     *kubernetes.Clientset // raw client set
	dynamicClient dynamic.Interface     // dynamic client
}

// NewCRDBase create a new crd base object for future use
func NewCRDBase(conf CRDBaseConfig, log ...logr.Logger) (*CRDBase, error) {
	ret := &CRDBase{
		CRDBaseConfig: conf,
	}

	if len(log) > 0 {
		ret.log = log[0].WithName(crdBaseName)
	} else {
		ret.log = utils.NewNullLogger()
	}

	if conf.Manager != nil {
		managerConf := conf.Manager.GetConfig()

		client, err := pkg_client.New(managerConf, pkg_client.Options{})
		if err != nil {
			return nil, err
		}
		client = pkg_client.NewNamespacedClient(client, conf.Namespace)
		ret.client = client

		clientSet, err := kubernetes.NewForConfig(managerConf)
		if err != nil {
			return nil, err
		}
		ret.clientSet = clientSet

		dynamicClient, err := dynamic.NewForConfig(managerConf)
		if err != nil {
			return nil, err
		}
		ret.dynamicClient = dynamicClient
	}

	if err := initScheme(); err != nil {
		return nil, err
	}

	return ret, nil
}

func initScheme() error {
	sch := runtime.NewScheme()

	if err := rbacv1.AddToScheme(sch); err != nil {
		return err
	}

	if err := appsv1.AddToScheme(sch); err != nil {
		return err
	}

	if err := apiextv1.AddToScheme(sch); err != nil {
		return err
	}

	return nil
}

func (crdb *CRDBase) Clone() *CRDBase {
	return crdb
}
