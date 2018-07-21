// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

package ints

import (
	"testing"

	"github.com/pulumi/pulumi-kubernetes/pkg/openapi"
	"github.com/pulumi/pulumi/pkg/resource"
	"github.com/pulumi/pulumi/pkg/testing/integration"
	"github.com/stretchr/testify/assert"
)

// TestPod tests booting
func TestPod(t *testing.T) {
	integration.ProgramTest(t, &integration.ProgramTestOptions{
		Dir:          "step1",
		Dependencies: []string{"@pulumi/pulumi", "@pulumi/kubernetes"},
		Quick:        true,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			assert.NotNil(t, stackInfo.Deployment)
			assert.Equal(t, 2, len(stackInfo.Deployment.Resources))
			stackRes := stackInfo.Deployment.Resources[0]
			assert.Equal(t, resource.RootStackType, stackRes.URN.Type())

			//
			// Assert pod is successfully created.
			//

			pod := stackInfo.Deployment.Resources[1]
			assert.Equal(t, "pod-test", string(pod.URN.Name()))

			// Status is "Running"
			phase, _ := openapi.Pluck(pod.Outputs, "live", "status", "phase")
			assert.Equal(t, "Running", phase)

			// Status "Ready" is "True".
			conditions, _ := openapi.Pluck(pod.Outputs, "live", "status", "conditions")
			ready := conditions.([]interface{})[1].(map[string]interface{})
			readyType, _ := ready["type"]
			assert.Equal(t, "Ready", readyType)
			readyStatus, _ := ready["status"]
			assert.Equal(t, "True", readyStatus)

			// Container is called "nginx" and uses image "nginx:1.13-alpine".
			containerStatuses, _ := openapi.Pluck(pod.Outputs, "live", "status", "containerStatuses")
			containerStatus := containerStatuses.([]interface{})[0].(map[string]interface{})
			containerName, _ := containerStatus["name"]
			assert.Equal(t, "nginx", containerName)
			image, _ := containerStatus["image"]
			assert.Equal(t, "nginx:1.13-alpine", image)
		},
		EditDirs: []integration.EditDir{
			{
				Dir:      "step2",
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					assert.NotNil(t, stackInfo.Deployment)
					assert.Equal(t, 2, len(stackInfo.Deployment.Resources))
					stackRes := stackInfo.Deployment.Resources[0]
					assert.Equal(t, resource.RootStackType, stackRes.URN.Type())

					//
					// Assert pod is successfully replaced with the new version, running nginx:1.15-alpine.
					//

					pod := stackInfo.Deployment.Resources[1]
					assert.Equal(t, "pod-test", string(pod.URN.Name()))

					// Status is "Running"
					phase, _ := openapi.Pluck(pod.Outputs, "live", "status", "phase")
					assert.Equal(t, "Running", phase)

					// Status "Ready" is "True".
					conditions, _ := openapi.Pluck(pod.Outputs, "live", "status", "conditions")
					ready := conditions.([]interface{})[1].(map[string]interface{})
					readyType, _ := ready["type"]
					assert.Equal(t, "Ready", readyType)
					readyStatus, _ := ready["status"]
					assert.Equal(t, "True", readyStatus)

					// Container is called "nginx" and uses image "nginx:1.13-alpine".
					containerStatuses, _ := openapi.Pluck(pod.Outputs, "live", "status", "containerStatuses")
					containerStatus := containerStatuses.([]interface{})[0].(map[string]interface{})
					containerName, _ := containerStatus["name"]
					assert.Equal(t, "nginx", containerName)
					image, _ := containerStatus["image"]
					assert.Equal(t, "nginx:1.15-alpine", image)
				},
			},
		},
	})
}
