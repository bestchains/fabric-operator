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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1beta1 "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeVotes implements VoteInterface
type FakeVotes struct {
	Fake *FakeIbp
	ns   string
}

var votesResource = schema.GroupVersionResource{Group: "ibp.com", Version: "", Resource: "votes"}

var votesKind = schema.GroupVersionKind{Group: "ibp.com", Version: "", Kind: "Vote"}

// Get takes name of the vote, and returns the corresponding vote object, and an error if there is any.
func (c *FakeVotes) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.Vote, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(votesResource, c.ns, name), &v1beta1.Vote{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Vote), err
}

// List takes label and field selectors, and returns the list of Votes that match those selectors.
func (c *FakeVotes) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.VoteList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(votesResource, votesKind, c.ns, opts), &v1beta1.VoteList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.VoteList{ListMeta: obj.(*v1beta1.VoteList).ListMeta}
	for _, item := range obj.(*v1beta1.VoteList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested votes.
func (c *FakeVotes) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(votesResource, c.ns, opts))

}

// Create takes the representation of a vote and creates it.  Returns the server's representation of the vote, and an error, if there is any.
func (c *FakeVotes) Create(ctx context.Context, vote *v1beta1.Vote, opts v1.CreateOptions) (result *v1beta1.Vote, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(votesResource, c.ns, vote), &v1beta1.Vote{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Vote), err
}

// Update takes the representation of a vote and updates it. Returns the server's representation of the vote, and an error, if there is any.
func (c *FakeVotes) Update(ctx context.Context, vote *v1beta1.Vote, opts v1.UpdateOptions) (result *v1beta1.Vote, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(votesResource, c.ns, vote), &v1beta1.Vote{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Vote), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeVotes) UpdateStatus(ctx context.Context, vote *v1beta1.Vote, opts v1.UpdateOptions) (*v1beta1.Vote, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(votesResource, "status", c.ns, vote), &v1beta1.Vote{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Vote), err
}

// Delete takes name of the vote and deletes it. Returns an error if one occurs.
func (c *FakeVotes) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(votesResource, c.ns, name), &v1beta1.Vote{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeVotes) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(votesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1beta1.VoteList{})
	return err
}

// Patch applies the patch and returns the patched vote.
func (c *FakeVotes) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.Vote, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(votesResource, c.ns, name, pt, data, subresources...), &v1beta1.Vote{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Vote), err
}