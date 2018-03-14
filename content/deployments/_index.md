---
title: "Deployments"
date: 2018-03-14T15:00:28Z
---

Deploy specifies the manner that Gophersauce should build your application. You will find this tag within your project's GXML file. Use `webapp` to have GoS generate a webserver for you. To export this package to an app use `package` instead of webapp. This tag should always be within the root of the <gos/> tag.

#### Available Deploy types :
1. webapp : Build your project as a monolith application.
2. package : Export your library to other GopherSauce and Go projects.
3. faas : Deploy your end tags as OpenFaaS functions.

*This tag is required

	<gos>
	    ...
	    <deploy>webapp</deploy>
	    ...
	</gos>



