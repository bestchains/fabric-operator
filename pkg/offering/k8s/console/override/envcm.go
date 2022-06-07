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
	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	"github.com/IBM-Blockchain/fabric-operator/pkg/manager/resources"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (o *Override) CM(object v1.Object, cm *corev1.ConfigMap, action resources.Action, options map[string]interface{}) error {
	instance := object.(*current.IBPConsole)
	switch action {
	case resources.Create:
		return o.CreateCM(instance, cm)
	case resources.Update:
		return o.UpdateCM(instance, cm)
	}

	return nil
}

func (o *Override) CreateCM(instance *current.IBPConsole, cm *corev1.ConfigMap) error {
	cm.Data["HOST_URL"] = "https://" + instance.Namespace + "-" + instance.Name + "-console" + "." + instance.Spec.NetworkInfo.Domain + ":443"

	err := o.CommonCM(instance, cm)
	if err != nil {
		return err
	}

	return nil
}
