---
title: "Methods"
date: 2018-03-14T14:59:00Z
---

This section covers how to extend template functionality with pipelines ( methods ) .

Attributes of template tag :

	name - Specifies the name of the function. Please keep in mind, that usage of this function within your package requires `Net` prepended to the declared name.
	If you desire to write methods for strict usage outside of template files, I recommend creating another Go file within your project directory with package name main.

	return - The return type of the function. Ie: string or a custom interface `DemoGos` 

	var - This is a comma delimited string of the function parameters. Say we need a function that has two parameters: param1 of type string and param 2 of type string. It will be declared as `param1 string,param2 string`.

	**InnerXml - This contains your Go statements.
	

Test case : The example below will declare a function that sends emails and returns a bool.

    <gos>
        ...
            <methods>
            ...
                <func name="sendEmail" var="to string,from string" return="bool">
                    fmt.Println("Send Email " ,to  , from)
                    return true
            	</func>
            ...
        </methods>
        ...
    </gos>
 
Generated package function :
    
    func NetsendEmail(to string, from string) bool

Use the following snippet within any .tmpl file to send an email.

	{{ if sendEmail "to" "from" }}
	          
		Email sent!

	{{else}}
	         
		Email error.

	{{end}}
