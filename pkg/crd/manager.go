/*
 * Copyright contributors to the Hyperledger Fabric Operator project
 *
 * SPDX-License-Identifier: Apache-2.0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at:
 *
 * 	  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package crd

import (
	"path/filepath"

	"github.com/IBM-Blockchain/fabric-operator/pkg/util"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	CACRD      = "./config/crd/bases/ibp_v1alpha1_ibpca.yaml"
	PeerCRD    = "./config/crd/bases/ibp_v1alpha1_ibppeer.yaml"
	OrdererCRD = "./config/crd/bases/ibp_v1alpha1_ibporderer.yaml"
	ConsoleCRD = "./config/crd/bases/ibp_v1alpha1_ibpconsole.yaml"

	CACRDFile      = "ibp_v1alpha1_ibpca.yaml"
	PEERCRDFIle    = "ibp_v1alpha1_ibppeer.yaml"
	ORDERERCRDFILE = "ibp_v1alpha1_ibporderer.yaml"
	CONSOLECRDFILE = "ibp_v1alpha1_ibpconsole.yaml"
)

var log = logf.Log.WithName("crd_manager")

//go:generate counterfeiter -o mocks/client.go -fake-name Client . Client

type Client interface {
	CreateCRD(crd *extv1.CustomResourceDefinition) (*extv1.CustomResourceDefinition, error)
}

type Manager struct {
	Client Client
	crds   []*extv1.CustomResourceDefinition
}

func GetCRDList() []string {
	return []string{CACRD, PeerCRD, OrdererCRD, ConsoleCRD}
}

func GetCRDListFromDir(dir string) []string {

	crds := []string{
		filepath.Join(dir, CACRDFile),
		filepath.Join(dir, PEERCRDFIle),
		filepath.Join(dir, ORDERERCRDFILE),
		filepath.Join(dir, CONSOLECRDFILE),
	}

	return crds
}

func NewManager(c Client, files ...string) (*Manager, error) {
	m := &Manager{
		Client: c,
	}
	for _, file := range files {
		crd, err := util.GetCRDFromFile(file)
		if err != nil {
			return nil, err
		}
		m.crds = append(m.crds, crd)
	}
	return m, nil
}

func (m *Manager) Create() error {
	log.Info("Create CRDs")
	for _, crd := range m.crds {
		log.Info("Creating", "CRD", crd.Name)
		_, err := m.Client.CreateCRD(crd)
		if err != nil {
			return err
		}
	}
	return nil
}
