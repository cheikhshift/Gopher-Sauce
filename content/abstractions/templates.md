---
title: "Templates"
date: 2018-03-14T14:58:44Z
---

This section covers how to declare reusable HTML templates.

Attributes of template tag : 
	
	name - This is the name of the template within your application. You can load a template by simply calling `{{<template name>}}` within any template in your project.

	tmpl - This specifies the path to the template file with the declared template root in mind. For example if your file was in `$GOPATH/PACKAGENAME/tmpl/file.tmpl` the tmpl attribute will be `file` because GoS will prepend the absolute path. 

	struct - This specifies the name of the Go struct literal to use with template. 
	In the example below, a struct will be declared as an alert.

The example is using a template file with path bootstrap/alert.

project.gxml

	<gos>
	    ...
	    <header>
	    	<!-- You may also declare your structs from another Go file within this directory. -->
	        <struct name="Bootstrap_alert">
	                Strong string
	                Text string
	                Type string
	        </struct> 
	        ...
	    </header>
	    ...
	    <templates>
	        ...
	         <template name="Bootstrap_alert" tmpl="bootstrap/alert" struct="Bootstrap_alert" /> 
	    </templates>
	</gos>

contents of tmpl/bootstrap/alert.tmpl :

	 <div class="alert alert-{{.Type}} alert-dismissible fade in" role="alert">
		<strong>{{.Strong}}</strong> <p>{{.Text}}</p>
	</div>


Now to use the template within other templates there are two ways of doing this. (With the example above in mind)

    <!-- No parameters -->
    {{Bootstrap_alert}}

    <!-- Add  `c` to your template name to initialize its struct; And add `b`(b<template name> $var) to your template name to output the html. dbOperation is intended to update your interface's fields. -->
    {{ $struct := cBootstrap_alert }}
    {{ $struct | dbOperation "queryvalue"  }}
    {{bBootstrap_alert $struct}}