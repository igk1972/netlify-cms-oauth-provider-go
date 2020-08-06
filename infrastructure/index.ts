// Copyright 2016-2019, Pulumi Corporation.  All rights reserved.

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
    port: 80, 
    external: true,
    protocol: "HTTP"
});

// Step 3: Build and publish a Docker image to a private ECR registry.
const img = awsx.ecs.Image.fromPath("cms-oauth-img", "../");

// const randomRandomString = new random.RandomString("random", {
//     length: 16,
//     overrideSpecial: "/@£$",
//     special: true,
// });

// Step 4: Create a Fargate service task that can scale out.
const appService = new awsx.ecs.FargateService("app-svc", {
    cluster,
    taskDefinitionArgs: {
        container: {
            image: img,
            cpu: 102 /*10% of 1024*/,
            memory: 50 /*MB*/,
            portMappings: [ web ],
            environment: [
                { 
                    name: "HOST",
                    value: ":80"
                },
                { 
                    name: "SESSION_SECRET",
                    value: "choprodrIbroh41huwru"
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
