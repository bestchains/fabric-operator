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

package v1beta1

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/IBM-Blockchain/fabric-operator/pkg/util"
)

// Component is a custom type that enumerates all the components (containers)
type Component string

const (
	INIT       Component = "INIT"
	CA         Component = "CA"
	ORDERER    Component = "ORDERER"
	PEER       Component = "PEER"
	GRPCPROXY  Component = "GRPCPROXY"
	FLUENTD    Component = "FLUENTD"
	DIND       Component = "DIND"
	COUCHDB    Component = "COUCHDB"
	CCLAUNCHER Component = "CCLAUNCHER"
	ENROLLER   Component = "ENROLLER"
	HSMDAEMON  Component = "HSMDAEMON"
)

func (crn *CRN) String() string {
	return fmt.Sprintf("crn:%s:%s:%s:%s:%s:%s:%s:%s:%s",
		crn.Version, crn.CName, crn.CType, crn.Servicename, crn.Location, crn.AccountID, crn.InstanceID, crn.ResourceType, crn.ResourceID)
}

func (catls *CATLS) GetBytes() ([]byte, error) {
	return util.Base64ToBytes(catls.CACert)
}

func (e *Enrollment) GetCATLSBytes() ([]byte, error) {
	if e.CATLS != nil {
		return e.CATLS.GetBytes()
	}
	return nil, errors.New("no CA TLS certificate set")
}

func (nsn NamespacedName) String() string {
	return nsn.Namespace + "-" + nsn.Name
}

func (annotations *BlockchainAnnotationList) Marshal() ([]byte, error) {
	return json.Marshal(annotations)
}

func (annotations *BlockchainAnnotationList) Unmarshal(raw []byte) error {
	if len(raw) == 0 {
		return nil
	}
	return json.Unmarshal(raw, annotations)
}

func (annotations *BlockchainAnnotationList) GetAnnotation(k string) (BlockchainAnnotation, error) {
	if annotations == nil {
		return BlockchainAnnotation{}, errors.New("nil annotation list")
	}
	annotation, ok := annotations.List[k]
	if !ok {
		return BlockchainAnnotation{}, errors.New("annotation not exist")
	}
	return annotation, nil
}
func (annotations *BlockchainAnnotationList) SetOrUpdateAnnotation(k string, annotation BlockchainAnnotation) (bool, error) {
	if annotations == nil {
		return false, errors.New("nil annotation list")
	}
	_, ok := annotations.List[k]
	annotations.List[k] = annotation

	annotations.LastAppliedTime = time.Now().String()

	return ok, nil
}

func (annotations *BlockchainAnnotationList) DeleteAnnotation(k string) error {
	if annotations == nil {
		return errors.New("nil annotation list")
	}
	_, ok := annotations.List[k]
	if !ok {
		return nil
	}
	delete(annotations.List, k)
	annotations.LastAppliedTime = time.Now().String()
	return nil
}
