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

package organization

// Update defines a list of elements that we detect spec updates on
type Update struct {
	adminOrCAUpdated bool
}

func (u *Update) AdminOrCAUpdated() bool {
	return u.adminOrCAUpdated
}

// GetUpdateStackWithTrues is a helper method to print updates that have been detected
func (u *Update) GetUpdateStackWithTrues() string {
	stack := ""

	if u.AdminOrCAUpdated() {
		stack += "adminOrCAUpdated "
	}

	if len(stack) == 0 {
		stack = "emptystack "
	}

	return stack
}
