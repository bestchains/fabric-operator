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

package user

import (
	"context"
	"fmt"

	iam "github.com/IBM-Blockchain/fabric-operator/api/iam/v1alpha1"
	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	"github.com/IBM-Blockchain/fabric-operator/pkg/k8s/controllerclient"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("gggg")

type UnaryAction string

const (
	Add    UnaryAction = "add"
	Remove UnaryAction = "remove"
)

// Reconcile targetUser by `organization` and its `id type`
// - set blockchain annotation
// - add blockchain organization label
func Reconcile(c controllerclient.Client, targetUser string, organization, enrollmentID string, idType IDType, action UnaryAction) error {
	var err error

	u, err := GetUser(c, targetUser)
	if err != nil {
		return err
	}

	switch action {
	case Add:
		err = ReconcileAdd(u, organization, enrollmentID, idType)
	case Remove:
		_, err = ReconcileRemove(u, organization, enrollmentID, idType)
	}

	if err != nil {
		return err
	}

	err = PatchUsers(c, *u)
	if err != nil {
		return err
	}

	return nil
}

func ReconcileAdd(u *iam.User, organization, enrollmentID string, idType IDType) error {
	var err error

	// add a organization label to targetUser
	if idType == CLIENT || idType == ADMIN {
		u.Labels[OrganizationLabel.String(organization)] = idType.String()
	}

	// set  annotation to current admin User
	annotationList := NewBlockchainAnnotationList()
	err = annotationList.Unmarshal([]byte(u.Annotations[BlockchainAnnotationKey]))
	if err != nil {
		return err
	}
	var id ID
	switch idType {
	case ADMIN:
		id = BuildAdminID(GetUsername(u))
	case ORDERER:
		id = BuildOrdererID(enrollmentID)
	case PEER:
		id = BuildPeerID(enrollmentID)
	default:
		id = BuildClientID(GetUsername(u))
	}

	var annotation BlockchainAnnotation
	annotation, err = annotationList.GetAnnotation(organization)
	if errors.Is(err, ErrAnnotationNotExist) {
		annotation = *NewBlockchainAnnotation(organization, id)
	}
	err = annotation.SetID(id)
	if err != nil {
		return err
	}
	err = annotationList.SetAnnotation(organization, annotation)
	if err != nil {
		return err
	}

	raw, err := annotationList.Marshal()
	if err != nil {
		return err
	}

	u.Annotations[BlockchainAnnotationKey] = string(raw)

	return nil
}

// ReconcileRemove remove organization's organization from User
func ReconcileRemove(u *iam.User, organization, enrollmentID string, idType IDType) (BlockchainAnnotation, error) {
	var err error

	// remove organization labels
	if idType == CLIENT || idType == ADMIN {
		delete(u.Labels, OrganizationLabel.String(organization))
	}

	// remove annotation under relevant organization
	annotationList := NewBlockchainAnnotationList()
	err = annotationList.Unmarshal([]byte(u.Annotations[BlockchainAnnotationKey]))
	if err != nil {
		return BlockchainAnnotation{}, err
	}
	// cache annotation for return
	annotation, err := annotationList.GetAnnotation(organization)
	if err != nil {
		return BlockchainAnnotation{}, errors.Wrap(err, "organization "+organization)
	}
	switch idType {
	case ORDERER, PEER:
		err = annotation.RemoveID(enrollmentID)
		if err == nil {
			_ = annotationList.SetAnnotation(organization, annotation)
		}
	default:
		// delete organization from user's annotation
		err = annotationList.DeleteAnnotation(organization)
		if err != nil && err != ErrAnnotationNotExist {
			return BlockchainAnnotation{}, err
		}
	}
	raw, err := annotationList.Marshal()
	if err != nil {
		return BlockchainAnnotation{}, err
	}
	u.Annotations[BlockchainAnnotationKey] = string(raw)

	return annotation, nil
}

// ReconcileTransfer reconcile the `Admin` transfer from `origin user` to `target user`
// transfer means:
//   - annotation will be removed from `origin user`
//   - annotation will be set to `target user`
func ReconcileTransfer(c controllerclient.Client, originUser string, targetUser string, organization string) error {
	var err error

	origin, err := GetUser(c, originUser)
	if err != nil {
		return errors.Wrap(err, "unable to find originUser")
	}

	target, err := GetUser(c, targetUser)
	if err != nil {
		return errors.Wrap(err, "unable to find targetUser")
	}

	annotation, err := ReconcileRemove(origin, organization, "", ADMIN)
	if err != nil {
		return err
	}

	// origin user's id should be removed
	err = annotation.RemoveID(originUser)
	if err != nil {
		return err
	}
	// target user's id should be set
	err = annotation.SetID(BuildAdminID(targetUser))
	if err != nil {
		return err
	}

	// add annotation to target admin
	toAnnotationList := NewBlockchainAnnotationList()
	toAnnotationList.SetAnnotation(organization, annotation)
	raw, err := toAnnotationList.Marshal()
	if err != nil {
		return err
	}
	target.Annotations[BlockchainAnnotationKey] = string(raw)

	// add organization label to target admin
	target.Labels[OrganizationLabel.String(organization)] = ADMIN.String()

	// patch users along with the changes
	err = PatchUsers(c, *origin, *target)
	if err != nil {
		return err
	}

	return nil
}

// GetUsername return `user.spec.name` as its username
func GetUsername(u *iam.User) string {
	return u.Spec.Name
}

// GetUser with label `t7d.io.username`(Token generated by IAM will use `user.spec.name` as username)
func GetUser(c controllerclient.Client, targetUser string) (*iam.User, error) {
	selector, err := UserSelector(targetUser)
	if err != nil {
		return nil, errors.Wrap(err, "invalid user selector")
	}
	ul, err := ListUsers(c, selector)
	if err != nil {
		return nil, err
	}
	if len(ul.Items) == 0 {
		return nil, errors.Errorf("user %s not found", targetUser)
	}
	if len(ul.Items) > 1 {
		return nil, errors.Errorf("user %s has %d items", targetUser, len(ul.Items))
	}
	return &ul.Items[0], nil
}

// ListUsers with a selector
func ListUsers(c controllerclient.Client, selector labels.Selector) (*iam.UserList, error) {
	userList := &iam.UserList{}
	err := c.List(context.TODO(), userList, &client.ListOptions{
		LabelSelector: selector,
	})
	for i, u := range userList.Items {
		log.Info(fmt.Sprintf("%d, name:%s, label:%s", i, u.Name, u.Labels))
	}
	userList1 := &iam.UserList{}
	_ = c.List(context.TODO(), userList1)
	for i, u := range userList1.Items {
		log.Info(fmt.Sprintf("full %d, name:%s, label:%s", i, u.Name, u.Labels))
	}

	selector1 := labels.SelectorFromSet(map[string]string{current.NETWORK_FEDERATION_LABEL: "federation-sample"})
	n1 := &current.NetworkList{}
	_ = c.List(context.TODO(), n1, &client.ListOptions{LabelSelector: selector1})
	log.Info(fmt.Sprintf("use selector to get network label %s, get data len %d", selector1.String(), len(n1.Items)))
	if len(n1.Items) > 0 {
		log.Info("use selector to get network label, get data, show name and selector")
		for _, n := range n1.Items {
			log.Info(fmt.Sprintf("show name %s and selector %v", n.Name, n.Labels))
		}
	}
	n2 := &current.NetworkList{}
	_ = c.List(context.TODO(), n2)
	log.Info(fmt.Sprintf("list network full, get data len %d", len(n2.Items)))
	if len(n2.Items) > 0 {
		log.Info("list full get data, show name and selector")
		for _, n := range n2.Items {
			log.Info(fmt.Sprintf("show name %s and selector %v", n.Name, n.Labels))
		}
	}
	selector2, err := labels.Parse(current.NETWORK_FEDERATION_LABEL + "=federation-sample")
	if err != nil {
		log.Error(err, "get error")
	}
	n3 := &current.NetworkList{}
	_ = c.List(context.TODO(), n3, client.MatchingLabels{current.NETWORK_FEDERATION_LABEL: "federation-sample", current.NETWORK_INITIATOR_LABEL: "org1"})
	log.Info(fmt.Sprintf("list network 3, get data len %d", len(n3.Items)))
	if len(n3.Items) > 0 {
		log.Info("list network 3 get data, show name and selector")
		for _, n := range n3.Items {
			log.Info(fmt.Sprintf("show name %s and selector %v", n.Name, n.Labels))
		}
		log.Info(fmt.Sprintf("1 selector %+v, 3 selector %+v", selector1, selector2))
	}
	if len(n2.Items) > 0 {
		log.Info(fmt.Sprintf("finally, 1 selector %+v, 3 selector %+v", selector1, labels.SelectorFromSet(n2.Items[0].Labels)))
	}
	if err != nil {
		return nil, err
	}
	return userList, nil
}

// PatchUsers do `k8s patch` on Users
func PatchUsers(c controllerclient.Client, users ...iam.User) error {
	var err error
	for _, u := range users {
		err = c.Patch(context.TODO(), &u, nil, controllerclient.PatchOption{
			Resilient: &controllerclient.ResilientPatch{
				Retry:    2,
				Into:     &iam.User{},
				Strategy: client.MergeFrom,
			},
		})

		if err != nil {
			return err
		}
	}
	return nil
}
