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

package federation

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Federation update", func() {
	It("empty stack", func() {
		update := &Update{}
		Expect(update.GetUpdateStackWithTrues()).To(Equal("emptystack "))
	})
	It("full stack", func() {
		update := &Update{
			specUpdated:   true,
			memberUpdated: true,
		}
		Expect(update.GetUpdateStackWithTrues()).To(Equal("specUpdated memberUpdated "))
	})
})
