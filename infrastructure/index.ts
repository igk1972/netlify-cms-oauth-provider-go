// Copyright 2016-2020, Pulumi Corporation.  All rights reserved.

import * as awsx from "@pulumi/awsx";
import * as pulumi from "@pulumi/pulumi"
import { config } from "process";
import * as random from "@pulumi/random";

// Step 1: Read stack inputs for Github key and Github secret
const pulumiConfig = new pulumi.Config();
const githubKey = pulumiConfig.require("githubKey");
const githubSecret = pulumiConfig.requireSecret("githubSecret")

// Step 1: Create an ECS Fargate cluster.
const cluster = new awsx.ecs.Cluster("cluster");

// Step 2: Define the Networking for our service.
const alb = new awsx.elasticloadbalancingv2.ApplicationLoadBalancer(
    "net-lb", { external: true, securityGroups: cluster.securityGroups });

const web = alb.createListener("web", { 
    port: 443,
    external: true,
    protocol: "HTTPS",
    certificateArn: "arn:aws:acm:us-east-1:616138583583:certificate/607bd17c-9e6e-438a-a90e-a6a2cbfdc678"
});

const tg = alb.createTargetGroup("oauth-tg", {
    port: 80
});

new awsx.lb.ListenerRule("oauth-listener-rule", web, {
    actions: [{
        type: "forward",
        targetGroupArn: tg.targetGroup.arn,
    }],
    conditions: [{
        field: "path-pattern",
        values: "/*",
    }],
});

// Step 3: Build and publish a Docker image to a private ECR registry.
const img = awsx.ecs.Image.fromPath("cms-oauth-img", "../");

// Create a random string and also mark its `result` property as a secret,
// so it is not stored in plaintext in the stack's state.
const sessionSecretRandomString = new random.RandomPassword("random", {
    length: 32,
}, { additionalSecretOutputs: ["result"] });

// Step 4: Create a Fargate service task that can scale out.
const appService = new awsx.ecs.FargateService("app-svc", {
    cluster,
    taskDefinitionArgs: {
        container: {
            image: img,
            memory: 128 /*MB*/,
            portMappings: [ tg ],
            environment: [
                { 
                    name: "HOST",
                    value: ":80"
                },
                { 
                    name: "SESSION_SECRET",
                    value: sessionSecretRandomString.result
                },
                {
                    name: "GITHUB_KEY",
                    value: githubKey
                },
                {
                    name: "GITHUB_SECRET",
                    value: githubSecret
                }
            ]
        },
    },
    desiredCount: 1,
});

// Step 5: Export the Internet address for the service.
export const url = web.endpoint.hostname;
