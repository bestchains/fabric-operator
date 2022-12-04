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

package override

import (
	"context"

	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	"github.com/IBM-Blockchain/fabric-operator/pkg/k8s/controllerclient"
	"github.com/IBM-Blockchain/fabric-operator/pkg/manager/resources/clusterrolebinding"
	"k8s.io/apimachinery/pkg/types"
)

type Override struct {
	Client controllerclient.Client

	SubjectKind clusterrolebinding.SubjectKind
}

func (o *Override) GetSubjectKind() clusterrolebinding.SubjectKind {
	if o.SubjectKind == "" {
		return clusterrolebinding.ServiceAccount
	}
	return o.SubjectKind
}

func (o *Override) GetOrganization(member current.NamespacedName) (*current.Organization, error) {
	organization := &current.Organization{}
	if err := o.Client.Get(context.TODO(), types.NamespacedName{Name: member.Name, Namespace: member.Namespace}, organization); err != nil {
		return nil, err
	}

	return organization, nil
}
