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

package operatorerrors

import (
	"fmt"

	"github.com/go-logr/logr"
)

const (
	InvalidDeploymentCreateRequest = iota + 1
	InvalidDeploymentUpdateRequest
	InvalidServiceCreateRequest
	InvalidServiceUpdateRequest
	InvalidPVCCreateRequest
	InvalidPVCUpdateRequest
	InvalidConfigMapCreateRequest
	InvalidConfigMapUpdateRequest
	InvalidServiceAccountCreateRequest
	InvalidServiceAccountUpdateRequest
	InvalidRoleCreateRequest
	InvalidRoleUpdateRequest
	InvalidRoleBindingCreateRequest
	InvalidRoleBindingUpdateRequest
	InvalidPeerInitSpec
	InvalidOrdererType
	InvalidOrdererNodeCreateRequest
	InvalidOrdererNodeUpdateRequest
	InvalidOrdererInitSpec
	CAInitilizationFailed
	OrdererInitilizationFailed
	PeerInitilizationFailed
	MigrationFailed
	FabricPeerMigrationFailed
	FabricOrdererMigrationFailed
	InvalidCustomResourceCreateRequest
	FabricCAMigrationFailed
	OrganizationInitilizationFailed
	FederationInitilizationFailed
	InvalidClusterRoleCreateRequest
	InvalidClusterRoleUpdateRequest
	InvalidClusterRoleBindingCreateRequest
	InvalidClusterRoleBindingUpdateRequest
)

var (
	BreakingErrors = map[int]*struct{}{
		InvalidDeploymentCreateRequest:         nil,
		InvalidDeploymentUpdateRequest:         nil,
		InvalidServiceCreateRequest:            nil,
		InvalidServiceUpdateRequest:            nil,
		InvalidPVCCreateRequest:                nil,
		InvalidPVCUpdateRequest:                nil,
		InvalidConfigMapCreateRequest:          nil,
		InvalidConfigMapUpdateRequest:          nil,
		InvalidServiceAccountCreateRequest:     nil,
		InvalidServiceAccountUpdateRequest:     nil,
		InvalidRoleCreateRequest:               nil,
		InvalidRoleUpdateRequest:               nil,
		InvalidRoleBindingCreateRequest:        nil,
		InvalidRoleBindingUpdateRequest:        nil,
		InvalidPeerInitSpec:                    nil,
		InvalidOrdererType:                     nil,
		InvalidOrdererInitSpec:                 nil,
		CAInitilizationFailed:                  nil,
		OrdererInitilizationFailed:             nil,
		PeerInitilizationFailed:                nil,
		FabricPeerMigrationFailed:              nil,
		FabricOrdererMigrationFailed:           nil,
		InvalidCustomResourceCreateRequest:     nil,
		FabricCAMigrationFailed:                nil,
		OrganizationInitilizationFailed:        nil,
		FederationInitilizationFailed:          nil,
		InvalidClusterRoleCreateRequest:        nil,
		InvalidClusterRoleUpdateRequest:        nil,
		InvalidClusterRoleBindingCreateRequest: nil,
		InvalidClusterRoleBindingUpdateRequest: nil,
	}
)

type OperatorError struct {
	Code    int
	Message string
}

func (e *OperatorError) Error() string {
	return e.String()
}

func (e *OperatorError) String() string {
	return fmt.Sprintf("Code: %d - %s", e.Code, e.Message)
}

func New(code int, msg string) *OperatorError {
	return &OperatorError{
		Code:    code,
		Message: msg,
	}
}

func Wrap(err error, code int, msg string) *OperatorError {
	return &OperatorError{
		Code:    code,
		Message: fmt.Sprintf("%s: %s", msg, err.Error()),
	}
}

func IsBreakingError(err error, msg string, log logr.Logger) error {
	oerr := IsOperatorError(err)
	if oerr == nil {
		return err
	}
	_, breakingError := BreakingErrors[oerr.Code]
	if breakingError {
		if log != nil {
			log.Error(err, fmt.Sprintf("Breaking Error: %s", msg))
		}
		return nil
	}
	return err
}

func GetErrorCode(err error) int {
	oerr := IsOperatorError(err)
	if oerr == nil {
		return 0
	}

	return oerr.Code
}

type Causer interface {
	Cause() error
}

// GetCause gets the root cause of the error
func IsOperatorError(err error) *OperatorError {
	for err != nil {
		switch err.(type) {
		case *OperatorError:
			return err.(*OperatorError)
		case Causer:
			err = err.(Causer).Cause()
		default:
			return nil
		}
	}
	return nil
}
