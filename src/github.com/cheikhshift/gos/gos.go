package main

import (
	"github.com/cheikhshift/gos/core"
	"fmt"
	"os"
	"strings"
//	"io/ioutil"
)

var webroot string
var template_root string
var gos_root string
var GOHOME string



func main() {
	GOHOME = os.ExpandEnv("$GOPATH") + "/src/"
    	//fmt.Println( os.Args)
    if len(os.Args) > 1 {
    //args := os.Args[1:]
    	if os.Args[1] == "compile" || os.Args[1] == "run" {
    		GOHOME = GOHOME   + strings.Trim(os.Args[2],"/")
    		serverconfig := os.Args[3]
    		webroot = os.Args[4]
    		template_root = os.Args[5]
    		fmt.Println("∑ GoS Speed compiler ");
    		coreTemplate,err := core.LoadGos( GOHOME + "/" + serverconfig ); 
			if err != nil {
				fmt.Println(err)
				return 
			}

			//fmt.Println(coreTemplate.Methods.Methods)
			coreTemplate.WriteOut = false
			core.Process(coreTemplate,GOHOME, webroot,template_root);

			if os.Args[1] == "run" {
				core.RunFile(GOHOME, coreTemplate.Output)
			}
    	}

	} else { 
	
    fmt.Println("∑ Welcome to Gos v1.0")
	fmt.Println("To begin please tell us a bit about the gos project you wish to compile");
	fmt.Printf("We need the GoS package folder relative to your $GOPATH/src (%v)\n", GOHOME)
   	gosProject := "" 
   	serverconfig := ""

   	fmt.Scanln(&gosProject)
   	GOHOME = GOHOME  + "/" + strings.Trim(gosProject,"/")
   	fmt.Printf("We need your Gos Project config source (%v)\n", GOHOME)
   	fmt.Scanln(&serverconfig)
    //fmt.Println(GOHOME)
	webroot,template_root = core.DoubleInput("What is the name of your webroot's folder ?", "What is the name of your template folder? ") 
		fmt.Println("Are you ready to begin? ");
		if core.AskForConfirmation() {
			fmt.Println("ΩΩ Operation Started!!");
			coreTemplate,err := core.LoadGos( GOHOME + "/" + serverconfig ); 
			if err != nil {
				fmt.Println(err)
				return 
			}

			fmt.Println(coreTemplate)

		} else {
			fmt.Println("Operation Cancelled!!")
		}
	}

}