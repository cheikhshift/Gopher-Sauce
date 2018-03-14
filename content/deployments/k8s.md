---
title: "K8s"
date: 2018-03-14T15:00:32Z
---

Here is a guide on deploying a GopherSauce project to a local Kubernetes cluster.


#### Requirements
1. Docker installed and running.
2. Kubernetes setup.
3. Kubectl.
4. A [GopherSauce](http://gophersauce.com) or Go project.
5.  [Minikube](https://kubernetes.io/docs/tutorials/stateless-application/hello-minikube/) (Optional, guide references `minikube dashboard` command )

*Make sure the project you plan on building is the current working directory. 

## Step 1 : Run Kubeconfig
Run the following command to write a new kubernetes deployment file. While the command is running it will suggest a docker command to run. This command will generate a docker image, with name corresponding to your deployment file (with name `default-deployment.yaml`).

	kubeconfig	

The command will look for a `gos.gxml` by default with port information. You may also specify the port with command `kubeconfig <PORT>`

## Step 2 : Build docker image
If you plan on using local images make sure you're using a local docker registry.

Build a docker image of your project.

	docker build -t {folderName} .
	
## Step 3 :  Create a deployment
Run the following command to launch a new deployment.

	kubectl create -f default-deployment.yaml

## Step 4 : Create service
The following command will create a new load balancer service of your deployment.

	kubectl expose deployment {folderName}-deployment --type=LoadBalancer

Notes : Replace `{folderName}` with the name of the current working directory (AKA project folder name).

## Step 5 : Access service
Run the following Minikube command to access your service:

	 minikube service {folderName}-deployment

## Step 6 : Manage and monitor
Run the following command to open the very helpful kubernetes dashboard. From here you can manage anything kubernetes related. 

	minikube dashboard
	

Deploying Stateful applications : Create a StatefulSet with configuration file starter `default-statefulset.yaml`,  generated when you run command `kubeconfig`. This will also create :
- a new service.
- a persistent volume & claim.


