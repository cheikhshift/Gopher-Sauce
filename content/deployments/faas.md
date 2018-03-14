---
title: "FaaS"
date: 2018-03-14T15:00:28Z
---


Here is a guide to help you build an OpenFaaS function with JSON responses. This function will convert the request body (JSON) to a GO struct. The function will then write a JSON response of an altered version of your request body.
PS. This is a lot faster with GopherSauce [Vim plugin](https://github.com/cheikhshift/vim-gos)

## Requirements

1. Install GopherSauce (`go get github.com/cheikhshift/gos`)
2. Docker running on host.
3. [OpenFaaS CLI](https://github.com/openfaas/faas).
4. [OpenFaaS Gateway (Link to setup-guide) ](https://github.com/openfaas/faas/blob/master/guide/deployment_swarm.md). Running & accessible OpenFaaS gateway.
5. Your `git` command setup & ready to perform commits.

## Assumptions
Your OpenFaaS gateway is at `http://localhost:8080`. TO update your gateway path update the `<gos>` tag within your `gos.gxml` file, add attribute `gateway="YOUR_GATEWAY"`. `gos.gxml` will be created after running command `gos --make`.

## Setup

### Github Repo
You will need a Git repository accessible online to successfully download your package's dependencies automatically. With this example I used this repository.
1. Create a new repository with local folder path template : `$GOPATH/src/{github.com || or other git service}/yourusername/foldername`

### GOS project
Change your working directory to the github repository you've just created.

### Create project
Run the following command to generate a new project. This will generate file `gos.gxml`

	gos --make

### Set deploy type 
Update the contents of `<deploy>` tag within your `gos.gxml` file, from `webapp` to `faas`.

### Import encode/json
Following your `<deploy>` tag add the following snippet to import Go pkg `encoding/json`

	<import src="encoding/json" />

## Define interface
Add the following snippet within the `<header>` tag of your `gos.gxml` file. This will define a new interface within GopherSauce

	<struct name="Testmodel" >
	//interface fields here
	</struct>

#### Add fields
Add snippet within newly placed `<struct>` tag. Declare FieldOne with type  `string`
		
	FieldOne string

Add snippet within newly placed `<struct>` tag. Declare FieldTwo with type  `int`

	FieldTwo int

Add snippet within newly placed `<struct>` tag. Declare FieldThree as an array of type  `string`

	FieldThree []string	
	
In the end your `<struct>` tag should look like this :

	<struct name="Testmodel" >
	//interface fields here
	FieldOne string
	FieldTwo int
	FieldThree []string
	</struct> 


## Create OpenFaaS function
Add the following snippet within the `<endpoints>` tag of your `gos.gxml` file. This will define a serverless function named GETTestJson on OpenFaaS.

	<end path="/test/json" type="GET" >
	</end>

#### Process json
1. Add the following snippet within your newly placed `<end>` tag. This will declare variable `t` with type `Testmodel`

	`var t Testmodel`	

2. Add the following snippet within your newly placed `<end>` tag. It will create a new `json.Decoder` with the input stream of your request's body. The body input stream is referred to as  `r.Body` (AKA FaaS request body). The variable `r` is available to  GO code within `<end>` tags. It refers to the current request (`*http.Request`).

	`decoder := json.NewDecoder(r.Body)`


3. Add the following snippet within your newly placed `<end>` tag. This will attempt to decode your request body into the interface of `t` (`Testmodel`). Your request body data will then be available with variable `t`. The variable `err` declared within the snippet will be used to check for errors.

	`err := decoder.Decode(&t)`
	
4. Add the following snippet within your newly placed `<end>` tag. Evaluate if error is `nil`, if so panic. Since we're serverless it is safe to panic :).

		if err != nil {
			panic(err)
		}

Add the following snippet within your newly placed `<end>` tag. This will update the field value of `FieldOne` to `"NewValue"`.

	t.FieldOne = "NewValue"

Add the following snippet within your newly placed `<end>` tag. The variable `response string` is available to `<end>` tags. Use this to set a JSON string response with the help of `func mResponse(v interface{})`.

	response = mResponse(t)

Your `<end>` tag should look like this :

	<end path="/test/json" type="GET" >
		//Golang code here
		 decoder := json.NewDecoder(r.Body)
		 var t Testmodel
		 err := decoder.Decode(&t)
		 if err != nil {
		    panic(err)
		 } 
		t.FieldOne = "NewValue"
		response = mResponse(t)
		//defered to keep body open if plan on using again.
		defer r.Body.Close()
	</end>


## Build project
Run your project with following command. (While in Github repository folder created earlier)

	gos --run


## Test function
Run the following command to invoke the function with `faas-cli`

	echo "{\"FieldOne\":\"Test\"}" | faas-cli invoke GETTestJson 
