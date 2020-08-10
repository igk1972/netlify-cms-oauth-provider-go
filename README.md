# About the Project
This project is to provide an OAuth Provider Server for the [Pulumi's CMS project](https://github.com/pulumi/doc-cms). The CMS for Pulumi website is going to deploy on AWS rather than on Netlify. Netlify use the Netlify Identity Service which provides OAuth provider server. Based on [Netlify's instruction](https://www.netlifycms.org/docs/external-oauth-clients/) of customize this step we need to provide our own OAuth client.

Netlify-CMS oauth client sends token in form as Netlify service itself. This project is implementated in Go (golang), connected with AWS Fargate and configured AWS Route 53 domain and certification using Pulumi
Here are some reference:
- @igk1972's [OAuth provider](https://github.com/igk1972/netlify-cms-oauth-provider-go) for OAuth Provider and it's frontend
- pulumi's [hello fargate example](https://github.com/pulumi/examples/tree/master/aws-ts-hello-fargate) for connecting to AWS Fargate to adopt Docker setting in cloud
- pulumi's [static website example](https://github.com/pulumi/examples/tree/master/aws-ts-static-website) for configuring certificate and obtain a subdomain for the provider server


## File Structure
- ./infrastructure
  - Pulumi code with setting up AWS Fargate and the configuring certificate and domain
- ./main.go the code for the provider itself and it's front end

# Getting Start

### Step 1. Register OAuth Application in Github and Obtain Key and Secret
- Steps are provided using this link https://docs.netlify.com/visitor-access/oauth-provider-tokens/#setup-and-settings
- For the Home Page Url use https://doc-cms.pulumi-demos.net
- For the Authorization callback URL enter https://doc-cms-oauth.pulumi-demos.net

### Step 2. Fill in the pulumi configuration
1. On the root directory do this.

2. Get into the infrastructure folder and initialize a new stack
```bash
$ cd infrastructure
$ pulumi stack init oauth-provider
```

3. Set AWS Region
```bash
$ pulumi config set aws:region us-east-1
```

4. Set Target Domain of OAuth Provider
```bash
$ pulumi config set netlify-cms-oauth-provider-infrastructure:targetDomain doc-cms-oauth.pulumi-demos.net
```

5. Set the Github Key and Secret
- change the {YOUR_GITHUB_KEY} and {YOUR_GITHUB_SECRET} with the key and secret obtain from Step 1.
```bash
$ pulumi config set netlify-cms-oauth-provider-infrastructure:githubKey {YOUR_GITHUB_KEY}
$ pulumi config set --secret netlify-cms-oauth-provider-infrastructure:githubSecret
$ {YOUR_GITHUB_SECRET}
```
- `--secret` tag is used to hash the secret so on the stack configuration yaml file it won't be shown
- Don't do ` $ pulumi config set --secret netlify-cms-oauth-provider-infrastructure:githubKey {YOUR_GITHUB_SECRET} `
because it might cause the secret be stored inside the command memory
- To make sure if key and secret is right do
```bash
$ pulumi config get netlify-cms-oauth-provider-infrastructure:githubKey
$ pulumi config get netlify-cms-oauth-provider-infrastructure:githubSecret
```


6. Don't forget to update AWS token before next step!!!!

### Step 3. Running Infrastructure
```bash
$ pulumi up
```

### Step 4. Config CMS
You also need to add `base_url` to the backend section of your netlify-cms's config file. 

Go to the doc-cms repo which stores resource for CMS of Pulumi's website https://github.com/pulumi/doc-cms and on file public/config.yml add the base_url line with https://doc-cms-oauth.pulumi-demos.net

```
backend:
  name: github
  repo: user/repo   # Path to your Github repository
  branch: master    # Branch to update
  base_url: https://doc-cms-oauth.pulumi-demos.net # Path to ext auth provider
```

Then build use 
```bash
$ yarn build
```
and go to the infrastructure folder and do pulumi up to up date changes


