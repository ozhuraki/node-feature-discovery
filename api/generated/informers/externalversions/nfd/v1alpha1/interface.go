/*
Copyright 2024 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	internalinterfaces "sigs.k8s.io/node-feature-discovery/api/generated/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// NodeFeatures returns a NodeFeatureInformer.
	NodeFeatures() NodeFeatureInformer
	// NodeFeatureGroups returns a NodeFeatureGroupInformer.
	NodeFeatureGroups() NodeFeatureGroupInformer
	// NodeFeatureRules returns a NodeFeatureRuleInformer.
	NodeFeatureRules() NodeFeatureRuleInformer
	// NodeFeatureStatuses returns a NodeFeatureStatusInformer.
	NodeFeatureStatuses() NodeFeatureStatusInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// NodeFeatures returns a NodeFeatureInformer.
func (v *version) NodeFeatures() NodeFeatureInformer {
	return &nodeFeatureInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// NodeFeatureGroups returns a NodeFeatureGroupInformer.
func (v *version) NodeFeatureGroups() NodeFeatureGroupInformer {
	return &nodeFeatureGroupInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// NodeFeatureRules returns a NodeFeatureRuleInformer.
func (v *version) NodeFeatureRules() NodeFeatureRuleInformer {
	return &nodeFeatureRuleInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}

// NodeFeatureStatuses returns a NodeFeatureStatusInformer.
func (v *version) NodeFeatureStatuses() NodeFeatureStatusInformer {
	return &nodeFeatureStatusInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}
