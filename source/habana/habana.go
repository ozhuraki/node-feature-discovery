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

package habana

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"k8s.io/klog/v2"

	nfdv1alpha1 "sigs.k8s.io/node-feature-discovery/api/nfd/v1alpha1"
	"sigs.k8s.io/node-feature-discovery/pkg/utils"
	"sigs.k8s.io/node-feature-discovery/pkg/utils/hostpath"
	"sigs.k8s.io/node-feature-discovery/source"
)

var habanaAIFields = [...]string{
	"device.count",
	"device.family",
	"distro.type",
	"hfd.timestamp",
	"kernel.version",
	"product.name",
	"sys.vendor",
	// "product.serial",
}

const (
	Name            = "habana.ai"
	HabanaAIFeature = "habana.ai"
	HabanaAIVendor  = "1da3"
	HabanaAIDevice  = "1020"
)

// Implements the FeatureSource and LabelSource interfaces.
type habanaSource struct {
	features *nfdv1alpha1.Features
}

// Singleton source instance
var (
	src habanaSource
	_   source.FeatureSource = &src
	_   source.LabelSource   = &src
)

func (s *habanaSource) Name() string { return Name }

// Priority method of the LabelSource interface
func (s *habanaSource) Priority() int { return 0 }

// GetLabels method of the LabelSource interface
func (s *habanaSource) GetLabels() (source.FeatureLabels, error) {
	labels := source.FeatureLabels{}
	features := s.GetFeatures()

	for _, key := range habanaAIFields {
		if value, exists := features.Attributes[HabanaAIFeature].Elements[key]; exists {
			labels[key] = value
		}
	}

	return labels, nil
}

// 60:00.0 Processing accelerators: Device 1da3:1000 (rev 01)

// Read a single PCI device attribute
// A PCI attribute in this context, maps to the corresponding sysfs file
func readSinglePciAttribute(devPath string, attrName string) (string, error) {
	data, err := os.ReadFile(filepath.Join(devPath, attrName))
	if err != nil {
		return "", fmt.Errorf("failed to read device attribute %s: %w", attrName, err)
	}
	// Strip whitespace and '0x' prefix
	attrVal := strings.TrimSpace(strings.TrimPrefix(string(data), "0x"))

	if attrName == "class" && len(attrVal) > 4 {
		// Take four first characters, so that the programming
		// interface identifier gets stripped from the raw class code
		attrVal = attrVal[0:4]
	}
	return attrVal, nil
}

// Read PCI vendor and device information
func readPciDevInfo(devPath string) (*nfdv1alpha1.InstanceFeature, error) {
	attrs := make(map[string]string)
	for _, attr := range []string{"vendor", "device"} {
		attrVal, err := readSinglePciAttribute(devPath, attr)
		if err != nil {
			return nil, fmt.Errorf("failed to read device %s: %w", attr, err)
		}
		attrs[attr] = attrVal
	}
	return nfdv1alpha1.NewInstanceFeature(attrs), nil
}

// detectPci detects available PCI devices and retrieves their device attributes.
// An error is returned if reading any of the mandatory attributes fails.
func detectPci() ([]nfdv1alpha1.InstanceFeature, error) {
	sysfsBasePath := hostpath.SysfsDir.Path("bus/pci/devices")

	devices, err := os.ReadDir(sysfsBasePath)
	if err != nil {
		return nil, err
	}

	// Iterate over devices
	devInfo := make([]nfdv1alpha1.InstanceFeature, 0, len(devices))
	for _, device := range devices {
		info, err := readPciDevInfo(filepath.Join(sysfsBasePath, device.Name()))
		if err != nil {
			klog.ErrorS(err, "failed to read PCI device info")
			continue
		}

		devInfo = append(devInfo, *info)
	}

	return devInfo, nil
}

// Discover method of the FeatureSource interface
func (s *habanaSource) Discover() error {
	s.features = nfdv1alpha1.NewFeatures()

	s.features.Attributes[HabanaAIFeature] = nfdv1alpha1.NewAttributeFeatures(nil)

	devs, err := detectPci()
	if err != nil {
		return fmt.Errorf("failed to detect PCI devices: %s", err.Error())
	}

	numOfDevices := 0

	for _, dev := range devs {
		attrs := dev.Attributes
		vendor := strings.ToLower(attrs["vendor"])
		device := strings.ToLower(attrs["device"])

		if strings.HasPrefix(string(vendor), strings.ToLower(HabanaAIVendor)) && strings.HasPrefix(string(device), strings.ToLower(HabanaAIDevice)) {
			numOfDevices++
		}
	}

	if numOfDevices > 0 {
		s.features.Attributes[HabanaAIFeature].Elements["device.count"] = fmt.Sprintf("%d", numOfDevices)
	}

	s.features.Attributes[HabanaAIFeature].Elements["device.count"] = fmt.Sprintf("%d", numOfDevices)

	version, err := os.ReadFile("/proc/sys/kernel/osrelease")
	if err != nil {
		return err
	}

	s.features.Attributes[HabanaAIFeature].Elements["kernel.version"] = strings.TrimSpace(string(version))

	vendor, err := os.ReadFile(hostpath.SysfsDir.Path("devices/virtual/dmi/id/sys_vendor"))
	if err != nil {
		return err
	}

	s.features.Attributes[HabanaAIFeature].Elements["sys.vendor"] = strings.TrimSpace(string(vendor))

	lsb, err := parseLSBRelease()
	if err != nil {
		return err
	}

	if distroType, ok := lsb["DISTRIB_ID"]; ok {
		s.features.Attributes[HabanaAIFeature].Elements["distro.type"] = distroType
	}

	klog.V(3).InfoS("discovered features", "featureSource", s.Name(), "features", utils.DelayedDumper(s.features))

	return nil
}

func parseLSBRelease() (map[string]string, error) {
	release := map[string]string{}

	f, err := os.Open(hostpath.EtcDir.Path("lsb-release"))
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`^(?P<key>\w+)=(?P<value>.+)`)

	// Read line-by-line
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		if m := re.FindStringSubmatch(line); m != nil {
			release[m[1]] = strings.Trim(m[2], `"'`)
		}
	}

	return release, nil
}

// GetFeatures method of the FeatureSource Interface
func (s *habanaSource) GetFeatures() *nfdv1alpha1.Features {
	if s.features == nil {
		s.features = nfdv1alpha1.NewFeatures()
	}
	return s.features
}

func init() {
	source.Register(&src)
}
