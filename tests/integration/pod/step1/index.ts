// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";

const pod = new k8s.core.v1.Pod("pod-test", {
  metadata: {
    name: "pod-test",
  },
  spec: {
    containers: [
      {name: "nginx", image: "nginx:1.13-alpine"},
    ],
  },
})
