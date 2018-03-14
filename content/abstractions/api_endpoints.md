---
title: "Web handlers"
date: 2018-03-14T15:00:10Z
---

### Defining a web service :

This section covers how to create an `end` tag (REST API endpoint). Your endpoints can be stateful as well as stateless. If a request's verb and path match a declared endpoint, the declared Go code will be ran.

Attributes of end tag :
path - URL path of endpoint
InnerXML - Go statements to be ran on endpoint execution.
type - This specifies the request verb ie: POST,GET,PUT,DELETE, star (to disregard request verb), f ( to execute as middleware)
id - Id of service. Used with open trace to find execution of service.

End tags are nested within `<endpoints>` in your GXML file.

The example below will declare a REST api endpoint called login. The tag defined will create endpoint : POST /index/api

	<gos>
	    ...
	    <endpoints>
	           <end path="/index/api" type="POST" >
	                // response sent to client will be {"Color":"#fff"}
	                response = mResponse(Button{Color:"#fff"})
	           </end>
	    </endpoints>
	    ...
	</gos>


List of variables available to your endpoint's Go code block :

	response - String response of api. The value set here will be returned as the endpoint's response.

	session - Current session (*sessions.Session from github.com/gorilla/sessions) of the request (if stateful endpoint). Use the session.Values (type map[string]interface{}) map to access and save data.

	r - *http.Request

	w - http.ResponseWriter

	span - opentracing.Span - this variable is only available when you run your project in development mode (with command gos --run). Use this function to log data to the tracer : span.LogEvent(event string) 

Useful functions :
	
	mResponse - will convert any Go interface/struct into a JSON string for output. The example below converts Button into a JSON string for output.