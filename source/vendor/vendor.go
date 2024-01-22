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

package vendor

import (
	"os"

	"k8s.io/klog/v2"

	nfdv1alpha1 "sigs.k8s.io/node-feature-discovery/pkg/apis/nfd/v1alpha1"
	"sigs.k8s.io/node-feature-discovery/pkg/utils"
	"sigs.k8s.io/node-feature-discovery/source"
)

// Name of this feature source
const Name = "vendor"

const (
	VendorFeature = "vendor"
	NameFeature   = "name"
)

// systemSource implements the FeatureSource and LabelSource interfaces.
type vendorSource struct {
	features *nfdv1alpha1.Features
}

// Singleton source instance
var (
	src vendorSource
	_   source.FeatureSource = &src
	_   source.LabelSource   = &src
)

func (s *vendorSource) Name() string { return Name }

// Priority method of the LabelSource interface
func (s *vendorSource) Priority() int { return 0 }

// GetLabels method of the LabelSource interface
func (s *vendorSource) GetLabels() (source.FeatureLabels, error) {
	labels := source.FeatureLabels{}
	features := s.GetFeatures()

	if value, exists := features.Attributes[VendorFeature].Elements["name"]; exists {
		labels[VendorFeature] = value
	}

	return labels, nil
}

// Discover method of the FeatureSource interface
func (s *vendorSource) Discover() error {
	s.features = nfdv1alpha1.NewFeatures()

	// Get node name
	s.features.Attributes[NameFeature] = nfdv1alpha1.NewAttributeFeatures(nil)
	s.features.Attributes[NameFeature].Elements["nodename"] = utils.NodeName()

	// Get vendor information
	vendor, err := getVendor()
	if err != nil {
		klog.ErrorS(err, "failed to get vendor")
	} else {
		s.features.Attributes[VendorFeature] = nfdv1alpha1.NewAttributeFeatures(nil)
		s.features.Attributes[VendorFeature].Elements["name"] = vendor
	}

	klog.V(3).InfoS("discovered features", "featureSource", s.Name(), "features", utils.DelayedDumper(s.features))

	return nil
}

// GetFeatures method of the FeatureSource Interface
func (s *vendorSource) GetFeatures() *nfdv1alpha1.Features {
	if s.features == nil {
		s.features = nfdv1alpha1.NewFeatures()
	}
	return s.features
}

// Read /sys/devices/virtual/dmi/id/sys_vendor
func getVendor() (string, error) {
	vendor, err := os.ReadFile("/sys/devices/virtual/dmi/id/sys_vendor")
	if err != nil {
		return "", err
	}

	return string(vendor), nil
}

func init() {
	source.Register(&src)
}
