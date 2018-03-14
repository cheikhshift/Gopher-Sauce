---
title: "TLS"
date: 2018-03-14T15:00:32Z
---

### Secure application with TLS 1.2 (RFC 5246)
Here is a guide to building a web application served over TLS with a [GopherSauce](http://gophersauce.com) project.


#### Requirements
1. GopherSauce project. (Make the project root your terminal's working directory)
2. HTTPS certificate & Key to use with application.


### Step 0
You can skip this step if you already have HTTPS certificate and key files. Read the part of this [Gist](https://gist.github.com/denji/12b3a568f092ab951456) to quickly (in my opinion) generate HTTPS certificate and key files.
### Step 1
Update your`gos.gxml` root XML tag (AKA `<gos>`) . Add attribute `https-cert` with value specifying a path to your certificate file.

### Step 2
Update your`gos.gxml` root XML tag (AKA `<gos>`) . Add attribute `https-key` with value specifying a path to your key file.

### Step 3
Update your`gos.gxml` port tag. Change it to `443`.

### Step 4
Run your project with TLS with the following command :

	gos --run
