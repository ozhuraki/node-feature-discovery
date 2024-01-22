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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystemSource(t *testing.T) {
	assert.Equal(t, src.Name(), Name)

	// Check that GetLabels works with empty features
	src.features = nil
	l, err := src.GetLabels()

	assert.Nil(t, err, err)
	assert.Empty(t, l)

}
