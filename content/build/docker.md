---
title: "Docker"
date: 2018-03-14T15:01:11Z
---

#### How to build docker image
GopherSauce will generate a docker file each time you build/export a project.

Run the following command to ensure your project's dependencies. This command requires dep installed.

	dep init

Build docker image with file.

	docker build -t ImageName . 
