---
title: "Serving templates"
date: 2018-03-14T15:00:28Z
---

This section covers the usage of dynamic web pages.

# Create a page
To add a new page, create a new `page_name.tmpl` file within your Gophersauce project web root. If you used the CLI argument --make, your web root folder would be web (relative to your project's directory). To access your page visit hostname/page_name.

The Go struct literal of a page in your web root :

     type Page struct {
            R *http.Request
            Session *sessions.Session
     }

These variables are accessible to template web pages within your project's web root. Field R would be accessed with syntax `{{ .R }}`. For templates in your template (tmpl) folder, Gophersauce will use the explicitly defined (`<struct>`)interface for that template. If no interface is set, `NoStruct{}` will be used, hence no fields available.


### Template pipelines

This section covers the list of functions available within all your templates compiled using GopherSauce. Please keep in mind that the .Session and .R variable are only available to template files within your server web root.

1. js - will take its only input and add it as the src attribute of the html `<script/>` tag `<script src="var"></script>`
	- Usage : `{{js "dist/js/bootstrap.js"}}`

2. css - Will take its only input and add it as the href attribute of the html `<link/>` tag.

	- Usage : `{{css "dist/css/bootstrap.css" }}`

3. sd - Will delete the current page session

	- Usage : `{{.Session | sd }}`

4. sr - Will remove a specified session key

	- Usage : `{{.Session | sr "KeyName" }}`

5. sc - Will check to see if a session key exists

	- Usage : `{{.Session | sc "KeyName" }}`
6. ss - Will set a string value as a session variable.

	- Usage : `{{.Session | ss "KeyName" "Variable" }}`
7. sso - Will set a struct as a session variable.

	- Usage : This requires three steps to work.

		+ Add tag `<import src="encoding/gob">` to your GXML root.

		+ Register the struct types within your `<init/>` tag of your GXML file.

		    
		     `gob.Register(&Object{})`

		+ Once the steps above are completed you can now set the linked object as session variables : `{{.Session | sso "KeyName" $object_with_Object_struct }}`

8. sgo - Will retrieve a a session stored interface.

	- Usage : `{{$desiredObject = .Session | sgo "keyName" }}`

8. sg - Will retrieve a string stored as a session variable.

	- Usage : `{{$string := .Session | sg "KeyName" }}`

9. form - Will retrieve a request variable no matter how it is submitted.
	- Usage : `{{ $input = .R | form "Key" }}` . .R is a page variable with type http.Request from the Go lang package net/http

10. eq - Will compare two variables and return a bool of value true if they are equal
	- Usage : `{{if eq "Obj1" "Obj1" }} {{end}}`
11. neq - Will compare two variables and return a bool of value true if they are not equal.

	- Usage : `{{if neq "Obj1" "Obj2" }} {{end}}`
12. lte - Will see if the first number declared is less than or equal to than the second number declared, if this statement proves to be true it will return a bool with the value true.

	- Usage : `{{if lte 5 10 }} {{end}}`
13. lt - Will see if the first number declared is less than the second number declared, if this statement proves to be true it wil return a bool with value true.

	- Usage : `{{if lt 5 7 }} {{end}}`
14. gte - Will see if the first number declared is greater than or equal to the second number declared, if this statement proves to be true it will return a bool with value true

	- Usage : `{{if gte 5 2}} {{end}}`
15. gt - Will see if the first number declared is greater than or equal to the second number declared, if this statement proves to be true it will return a bool with value true

	- Usage : `{{if gt 5 2}} {{end}}`