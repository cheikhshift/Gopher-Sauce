# GoS Mobile Tutorial

![enter image description here](http://s8.postimg.org/g84idl23p/header.png)

This tutorial covers writing a mobile application for IOS with Go lang and HTML. The application will add 50 to a session stored integer each time on application startup. GoS requires Xcode to run the apps and handle any other external processes your app may need to undergo.

The steps are as follows :

 2. [Create Package](#create-package)
 3. [Create Xcode Project](#create-xcode-project)
 4. [Edit Gos config](#edit-gos-config)
 4. [Setup index file](#setup-index)
 5. [Create button template](#create-button-template)
 6. [Gos compiler](#gos-compiler)
 7. [Run the app](#run-the-app)


### Requirements

 - Xcode 6 or later
 - [Go](https://golang.org/doc/install) (workspace set and functional)
 - [GoS](readme.md)

# Create Package
*There is a built in option for GoS to create a new Go package and add a web root directory, template directory and default GoS template file. To perform this, use the make option. Use the function below to create the package.

	$GOPATH/bin/gos make new/demo/pkg
 

# Create Xcode Project
Fire up Xcode and create a single view application. 

![enter image description here](http://s1.postimg.org/ylg1sqj0f/Screen_Shot_2015_12_21_at_8_42_14_AM.png)

To make the command line calls to inject Go lang into the project easier, have Xcode place the project within your newly created Go package

![enter image description here](http://s13.postimg.org/gefkv5i6f/Screen_Shot_2015_12_21_at_8_42_48_AM.png)


# Edit Gos Config
The generated xml file from the command ran initially needs a bit of adjusting. To make GoS compile this package and use it as an IOS app you will need to update the `deploy` tag and set it to bind.

*This is the upper portion of the gos.xml file: 

	<?xml version="1.0" encoding="UTF-8"?>
	<gos>
		<!--Stating the deployment type GoS should compile -->
		<!-- Curent valid types are webapp,shell and bind -->
		<!-- Shell = cli, sort of a GoS(Ghost) in the Shell -->
		<deploy>SETTHIS</deploy>
		<package>mymobile</package>
		
		<!-- Using import within different tags will have different results -->
		<!-- We going to make the goPkg Mongo Db Driver available to our application -->
		<!-- Using <import/> within the <go/> tag is similar to using the import call within a .go file -->
		<!-- To be less dramating, GoS will skip packages that it has already imported -->
		
		<!-- Go File output name -->
		<output>server_out.go</output>
		<!-- exported session fields available to Session -->
		...
 
 Update deploy and replace `SETTHIS` with `bind`

# Setup index file
The index file is crucial for the initial view of the app, otherwise a blank page will appear. The file extension can be html or tmpl, however the tmpl file is loaded over any match html file; And the file must be within your package's web directory.

## Setup basic HTML 5 page

The first part to setting up the index file is coding blank index.tmpl file with title bootstrap 101... .

	<!DOCTYPE html>
	<html lang="en">
	  <head>
	      <meta charset="utf-8">
	    <meta http-equiv="X-UA-Compatible" content="IE=edge">
	    <meta name="viewport" content="width=device-width, initial-scale=1">
	    <!-- The above 3 meta tags *must* come first in the head; any other head content must come *after* these tags -->
	    <title>Bootstrap 101 Template</title>
	
	    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
	    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
	    <!--[if lt IE 9]>
	      <script src="https://oss.maxcdn.com/html5shiv/3.7.2/html5shiv.min.js"></script>
	      <script src="https://oss.maxcdn.com/respond/1.4.2/respond.min.js"></script>
	    <![endif]-->
	  </head>
	  <body>
	     
	
	    <!-- jQuery (necessary for Bootstrap's JavaScript plugins) -->
	    <script src="https://ajax.googleapis.com/ajax/libs/jquery/1.11.3/jquery.min.js"></script>
	    <!-- Include all compiled plugins (below), or include individual files as needed -->
	   
	  </body>
	</html> 

## Using sessions to count
Notice that the starting page we have has no programming capability, making it a static web page stuck in the 2000s. To make the app more dynamic (literally), a number will be incremented by 50 and then written to the HTML. Since storedInt was never set within this app's session, the session function `sGetN` will return a 0. The code below will be added after opening the HTML `head` tag.

	    {{ $use := "storedInt" | sGetN  }}
        {{ $use := $use | a 50 }}
        {{sSet "storedInt" $use }}

	
- *`$use` is initialized, `sGetN` means `SessionGetNumber` A.k.A get a session stored number within this app. Once `$use` is manipulated it is stored again using `sSet`, SessionSet.
 
## Showing the count

Once we have the code in to add 50 on each refresh, the next step is being able to show where the value is at. The `sGet` (SessionGet) function will be used to get the stored number. We'll slip this tag right after the body tag.

	    <h2>Counting {{sGet "storedInt" }}</h2>


# Create button template
The template folder is useful for creating reusable HTML code within your project. In this example a bootstrap button template is used to reduce the amount of code required to use the button tag within HMTL. 

## Creating template file
First lets create a file at PACKAGENAME/tmpl/button.tmpl with this as its content : 

	<input class="btn btn-default" type="submit" value="Submit {{sGet "storedInt" }}"/>

The template file will load a session variable on its use as well.

## Create linkage and struct

Open up the `gos.xml` file and within the template tag add this line :

	<template name="Button" tmpl="button" struct="Button" /> 
Scroll up to add this struct within the header tag

	<struct name="Button">
		Color string
	</struct>

## Using it within HTML
To insert the button open two braces and enter the value entered for `name` attribute of the template. 

	{{Button}}

# GoS compiler
Now that all the files needed for the minimal app to run,  run the GoS cmd again with these params to compile this tutorial : `GoTetst2` is the relative Xcode project folder in the same path as the Go package.
	
		$GOPATH/bin/gos export new/demo/pkg gos.xml web tmpl GoTetst2

# Run the app
Now its like you just wrote serious C, run the app with Xcode and grab a cup of coffee, more to come with location and BLE support.

![enter image description here](http://s24.postimg.org/63ynabh9h/Screen_Shot_2015_12_21_at_9_40_52_AM.png)

![enter image description here](http://s4.postimg.org/avlwx161p/Simulator_Screen_Shot_Dec_21_2015_9_39_46_AM.png)

