package core

import (
	"fmt"
	"log"
	"io/ioutil"
	"encoding/xml"
	"bytes"
	"os"
	"strings"
	"unicode"
	"strconv"
	"os/exec"
	"bufio"
	"crypto/rand"
)


const (
	StdLen = 16
	UUIDLen = 20
)

var GOHOME = os.ExpandEnv("$GOPATH") + "/src/"
var available_methods []string
var	int_methods  []string
var	api_methods  []string
var	int_mappings []string
var StdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

func NewLen(length int) string {
	return NewLenChars(length, StdChars)
}

// NewLenChars returns a new random string of the provided length, consisting
// of the provided byte slice of allowed characters (maximum 256).
func NewLenChars(length int, chars []byte) string {
	if length == 0 {
		return ""
	}
	clen := len(chars)
	if clen < 2 || clen > 256 {
		panic("uniuri: wrong charset length for NewLenChars")
	}
	maxrb := 255 - (256 % clen)
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.
	i := 0
	for {
		if _, err := rand.Read(r); err != nil {
			panic("uniuri: error reading random bytes: " + err.Error())
		}
		for _, rb := range r {
			c := int(rb)
			if c > maxrb {
				// Skip this number to avoid modulo bias.
				continue
			}
			b[i] = chars[c%clen]
			i++
			if i == length {
				return string(b)
			}
		}
	}
}

func Process(template *gos,r string, web string, tmpl string) (local_string string) {
	// r = GOHOME + GoS Project
	arch := gosArch{}
	if template.Type == "webapp" {
	local_string = `package main 
import (`

	// if template.Type == "webapp" {
		
		net_imports := []string{"net/http", "time","github.com/gorilla/sessions","github.com/elazarl/go-bindata-assetfs","bytes","encoding/json" ,"fmt","html",  "html/template", "strings", "reflect", "unsafe"}
		/*
			Methods before so that we can create to correct delegate method for each object
		*/
		if template.Type == "bind" {
			net_imports = append(net_imports, "net/http/httptest")
		} else {
			net_imports = append(net_imports, "github.com/elazarl/go-bindata-assetfs")
		}

		for _,imp := range template.Methods.Methods {
			if !contains(available_methods, imp.Name) {
				available_methods = append(available_methods, imp.Name)
			}
		}
		apiraw := ``
		for _,imp := range template.Endpoints.Endpoints {
			if !contains(api_methods, imp.Method) {
				api_methods = append(api_methods, imp.Method)
			}
				meth := template.findMethod(imp.Method)
				apiraw += ` 
				if  r.URL.Path == "` + imp.Path +`" && r.Method == strings.ToUpper("` + imp.Type +`") { 
					` + meth.Method + `
					callmet = true
				}
				` 
			
		}
		timeline :=  ``
		for _,imp := range template.Timers.Timers {
			if !contains(api_methods, imp.Method) {
				api_methods = append(api_methods,imp.Method)
			}
			meth := template.findMethod(imp.Method)
			timeline += `
			` + imp.Name +` := time.NewTicker(time.` + imp.Unit + ` * ` + imp.Interval +`)
					    go func() {
					        for _ = range ` + imp.Name +`.C {
					           ` + meth.Method +`
					        }
					    }()
    `
		}

		
		fmt.Printf("APi Methods %v\n",api_methods)
		     netMa := 	`template.FuncMap{"js" : net_importjs,"css" : net_importcss,"sd" : net_sessionDelete,"sr" : net_sessionRemove,"sc": net_sessionKey,"ss" : net_sessionSet,"sso": net_sessionSetInt,"sgo" : net_sessionGetInt,"sg" : net_sessionGet,"form" : formval,"eq": equalz, "neq" : nequalz, "lte" : netlt`
           for _,imp := range available_methods {
           	if !contains(api_methods, imp) {
          		netMa += `,"` + imp + `" : net_` + imp
      		}
           }
           int_lok := []string{}

           for _,imp := range template.Header.Objects {
			//struct return and function
				
			if !contains(int_lok, imp.Name) {
				int_lok = append(int_lok,imp.Name)
			 	netMa += `,"` + imp.Name + `" : net_` + imp.Name
			}
			}


           for _,imp := range template.Templates.Templates {

				netMa += `,"` + imp.Name + `" : net_` + imp.Name
				netMa += `,"b` + imp.Name + `" : net_b` + imp.Name
				netMa += `,"c` + imp.Name + `" : net_c` + imp.Name
           }
           netMa += `}`

		for _,imp := range template.RootImports {
				//fmt.Println(imp)
			if !strings.Contains(imp.Src,".xml") {
					if imp.Download == "true" {
						fmt.Println("∑ Downloading Package " + imp.Src)
						RunCmd("go get " + imp.Src)
					}
					if  !contains(net_imports, imp.Src) {
						net_imports = append(net_imports, imp.Src)
					}
			}
		}

		fmt.Println(template.Methods.Methods[0].Name)

		for _,imp := range net_imports {
			local_string += `
			"` + imp + `"`
		}
		local_string += `
		)
				var store = sessions.NewCookieStore([]byte("` + template.Key +`"))

				func net_sessionGet(key string,s *sessions.Session) string {
					return s.Values[key].(string)
				}


				func net_sessionDelete(s *sessions.Session) string {
						//keys := make([]string, len(s.Values))

						//i := 0
						for k := range s.Values {
						   // keys[i] = k.(string)
						    net_sessionRemove(k.(string), s)
						    //i++
						}

					return ""
				}

				func net_sessionRemove(key string,s *sessions.Session) string {
					delete(s.Values, key)
					return ""
				}
				func net_sessionKey(key string,s *sessions.Session) bool {					
				 if _, ok := s.Values[key]; ok {
					    //do something here
				 		return true
					}

					return false
				}

				func net_sessionGetInt(key string,s *sessions.Session) interface{} {
					return s.Values[key]
				}

				func net_sessionSet(key string, value string,s *sessions.Session) string {
					 s.Values[key] = value
					 return ""
				}
				func net_sessionSetInt(key string, value interface{},s *sessions.Session) string {
					 s.Values[key] = value
					 return ""
				}

				func net_importcss(s string) string {
					return "<link rel=\"stylesheet\" href=\"" + s + "\" /> "
				}

				func net_importjs(s string) string {
					return "<script type=\"text/javascript\" src=\"" + s + "\" ></script> "
				}



				func formval(s string, r*http.Request) string {
					return r.FormValue(s)
				}
			
				func renderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, p *Page) {
				     filename :=  tmpl  + ".tmpl"
				    body, err := Asset(filename)
				    session, er := store.Get(r, "session-")

				 	if er != nil {
				           session,er = store.New(r,"session-")
				    }
				    p.Session = session
				    p.R = r
				    if err != nil {
				       fmt.Print(err)
				    } else {
				    t := template.New("PageWrapper")
				    t = t.Funcs(` + netMa + `)
				    t, _ = t.Parse(strings.Replace(strings.Replace(strings.Replace(BytesToString(body), "/{", "\"{",-1),"}/", "}\"",-1 ) ,"` + "`" + `", ` + "`" + `\"` + "`" +` ,-1) )
				    outp := new(bytes.Buffer)
				    error := t.Execute(outp, p)
				    if error != nil {
				    fmt.Print(error)
				    return
				    } 

				    p.Session.Save(r, w)

				    fmt.Fprintf(w, html.UnescapeString(outp.String()) )
				    }
				}

				func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
				  return func(w http.ResponseWriter, r *http.Request) {
				  	if !apiAttempt(w,r) {
				      fn(w, r, "")
				  	}
				  }
				} 

				func mHandler(w http.ResponseWriter, r *http.Request) {
				  	
				  	if !apiAttempt(w,r) {
				      handler(w, r, "")
				  	}
				  
				} 
				func mResponse(v interface{}) string {
					data,_ := json.Marshal(&v)
					return string(data)
				}
				func apiAttempt(w http.ResponseWriter, r *http.Request) bool {
					session, er := store.Get(r, "session-")
					response := ""
					if er != nil {
						session,_ = store.New(r, "session-")
					}
					callmet := false

					` + apiraw + `

					if callmet {
						session.Save(r,w)
						if response != "" {
							//Unmarshal json
							w.Header().Set("Access-Control-Allow-Origin", "*")
							w.Header().Set("Content-Type",  "application/json")
							w.Write([]byte(response))
						}
						return true
					}
					return false
				}

				func handler(w http.ResponseWriter, r *http.Request, context string) {
				  // fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
				  p,err := loadPage(r.URL.Path , context,r)
				  if err != nil {
				        http.Error(w, err.Error(), http.StatusInternalServerError)
				        return
				  }

				  if !p.isResource {
				        renderTemplate(w, r,  "` + web +`" + r.URL.Path, p)
				  } else {
				       w.Write(p.Body)
				  }
				}

				func loadPage(title string, servlet string,r *http.Request) (*Page,error) {
				    filename :=  "` +  web + `" + title + ".tmpl"
				    if title == "/" {
				    	filename = "` + web + `/index.tmpl"
				    }
				    body, err := Asset(filename)
				    if err != nil {
				      filename = "` + web + `" + title + ".html"
				      if title == "/" {
				    	filename = "` + web + `/index.html"
				    	}
				      body, err = Asset(filename)
				      if err != nil {
				         filename = "` + web + `" + title
				         body, err = Asset(filename)
				         if err != nil {
				            return nil, err
				         } else {
				          if strings.Contains(title, ".tmpl") || title == "/" {
				              return nil,nil
				          }
				          return &Page{Title: title, Body: body,isResource: true,request: nil}, nil
				         }
				      } else {
				         return &Page{Title: title, Body: body,isResource: true,request: nil}, nil
				      }
				    } 
				    //load custom struts
				    return &Page{Title: title, Body: body,isResource:false,request:r}, nil
				}
				func BytesToString(b []byte) string {
				    bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
				    sh := reflect.StringHeader{bh.Data, bh.Len}
				    return *(*string)(unsafe.Pointer(&sh))
				}
				func equalz(args ...interface{}) bool {
		    	    if args[0] == args[1] {
		        	return true;
				    }
				    return false;
				 }
				 func nequalz(args ...interface{}) bool {
				    if args[0] != args[1] {
				        return true;
				    }
				    return false;
				 }

				 func netlt(x,v float64) bool {
				    if x < v {
				        return true;
				    }
				    return false;
				 }
				 func netgt(x,v float64) bool {
				    if x > v {
				        return true;
				    }
				    return false;
				 }
				 func netlte(x,v float64) bool {
				    if x <= v {
				        return true;
				    }
				    return false;
				 }
				 func netgte(x,v float64) bool {
				    if x >= v {
				        return true;
				    }
				    return false;
				 }
				 type Page struct {
					    Title string
					    Body  []byte
					    request *http.Request
					    isResource bool
					    R *http.Request
					    Session *sessions.Session
					}`
					for _,imp := range template.Variables {
						local_string += `
						var ` + imp.Name + ` ` + imp.Type 
					}
		if template.Init_Func != "" {
			local_string += `
			func init(){
				` + template.Init_Func + `
			}`

		}

		//Lets Do structs
		for _,imp := range template.Header.Structs {
			if !contains(arch.objects, imp.Name) {
			fmt.Println("Processing Struct : " + imp.Name)
			arch.objects = append(arch.objects, imp.Name)
			local_string += `
			type ` + imp.Name + ` struct {`
			local_string += imp.Attributes
			local_string += `
			}`
			}
		}
		


		for _,imp := range template.Header.Objects {
			local_string += `
			type ` + imp.Name + ` ` + imp.Templ
		}
		

		//Create an object map
		for _,imp := range template.Header.Objects {
			//struct return and function
			fmt.Println("∑ Processing object :" + imp.Name)
			if !contains(available_methods, imp.Name) {
				//addcontructor
				available_methods = append(available_methods,imp.Name)
				int_methods = append(int_methods,imp.Name)
				local_string += `
				func  net_`+ imp.Name + `(args ...interface{}) (d ` + imp.Templ +`){
					if len(args) > 0 {
					jso := args[0].(string)
					var jsonBlob = []byte(jso)
					err := json.Unmarshal(jsonBlob, &d)
					if err != nil {
						fmt.Println("error:", err)
						return
					}
					return
					} else {
						d = ` + imp.Templ +`{} 
						return
					}
				}`	    

			}

			delegateMethods := strings.Split(imp.Methods,"\n")

			for _,im := range delegateMethods {
				
				if stripSpaces(im) != "" {
				fmt.Println(imp.Name + "->" + im)
				function_map := strings.Split(im, ")")

				

				if !contains(int_mappings, function_map[0] + imp.Templ) {
					int_mappings = append(int_mappings,function_map[0] + imp.Templ)
					funcsp := strings.Split(function_map[0],"(")
					meth := template.findMethod(stripSpaces(funcsp[0]))

					//process limits and keep local deritives
					if meth.Autoface == "" || meth.Autoface == "true"  {
						
						/*
							
						*/
						procc_funcs := true
						fmt.Println( )

						if meth.Limit != "" {
							if !contains(strings.Split(meth.Limit,","), imp.Name ){
								procc_funcs = false 
							}
						}

						if contains(api_methods, meth.Name){
							procc_funcs = false
						}
						
						objectName := meth.Object
						if objectName == "" {
							objectName = "object"
						}
						if procc_funcs {	
							if !contains(int_methods, stripSpaces(funcsp[0])) && meth.Name != "000" {
						int_methods = append(int_methods,stripSpaces(funcsp[0]))
						}
					  	local_string += `
					  	func  net_`+ stripSpaces(funcsp[0]) + `(` + strings.Trim(funcsp[1] + `, ` + objectName + ` ` +imp.Templ, ",") +`) ` + stripSpaces(function_map[1])
						if stripSpaces(function_map[1]) == "" {
							local_string += ` string`
						}

						local_string += ` {
									` + meth.Method

						if stripSpaces(function_map[1]) == "" {
							local_string += ` 
								return ""
							`
						}
						local_string += ` 
						}`



						if meth.Keeplocal == "false" || meth.Keeplocal == "" {
						local_string += `
						func (` + objectName + ` ` + imp.Templ +`) ` +  stripSpaces(funcsp[0]) + `(` + strings.Trim(funcsp[1], ",") +`) ` + stripSpaces(function_map[1])

							local_string += ` {
							` + meth.Method

						local_string +=  `
						}`
						}
						}
					}





				}
				}
			}

			//create Unused methods methods
			fmt.Println(int_methods)
			for _,imp := range available_methods {
				if !contains(int_methods,imp) && !contains(api_methods, imp)  {
					fmt.Println("Processing : " + imp)
						meth := template.findMethod(imp)
						addedit := false
						if meth.Returntype == "" {
							meth.Returntype = "string"
							addedit = true
						}
						local_string += `
						func net_` + meth.Name + `(args ...interface{}) ` + meth.Returntype + ` {
							`
						for k,nam := range strings.Split(meth.Variables,","){
							if nam != "" {
								local_string +=  nam + ` := ` + `args[` + strconv.Itoa(k) + `]
								`
							}
						}
						local_string += meth.Method
						if addedit {
						 local_string +=  `
						 return ""
						 `
						}
						local_string += `
						}` 
					}
			}
					for _,imp := range template.Templates.Templates {
				local_string += `
				func  net_`+ imp.Name + `(args ...interface{}) string {
					var d ` + imp.Struct + `
					if len(args) > 0 {
					jso := args[0].(string)
					var jsonBlob = []byte(jso)
					err := json.Unmarshal(jsonBlob, &d)
					if err != nil {
						fmt.Println("error:", err)
						return ""
					}
					} else {
						d = ` + imp.Struct +`{}
					}

					filename :=  "` + tmpl + `/` + imp.TemplateFile + `.tmpl"
    				body, er := Asset(filename)
    				if er != nil {
    					return ""
    				}
    				 output := new(bytes.Buffer) 
					t := template.New("` +  imp.Name + `")
    				t = t.Funcs(` + netMa +`)
				  	t, _ = t.Parse(strings.Replace(strings.Replace(strings.Replace(BytesToString(body), "/{", "\"{",-1),"}/", "}\"",-1 ) ,"` + "`" + `", ` + "`" + `\"` + "`" +` ,-1) )
			
				    error := t.Execute(output, &d)
				    if error != nil {
				    fmt.Print(error)
				    } 
					return html.UnescapeString(output.String())
				}`	    
					local_string += `
				func  net_b`+ imp.Name + `(d ` + imp.Struct +`) string {
					filename :=  "` + tmpl + `/` + imp.TemplateFile + `.tmpl"
    				body, er := Asset(filename)
    				if er != nil {
    					return ""
    				}
    				 output := new(bytes.Buffer) 
					t := template.New("` +  imp.Name + `")
    				t = t.Funcs(` + netMa +`)
				  	t, _ = t.Parse(strings.Replace(strings.Replace(strings.Replace(BytesToString(body), "/{", "\"{",-1),"}/", "}\"",-1 ) ,"` + "`" + `", ` + "`" + `\"` + "`" +` ,-1) )
			
				    error := t.Execute(output, &d)
				    if error != nil {
				    fmt.Print(error)
				    } 
					return html.UnescapeString(output.String())
				}`	    
				local_string += `
				func  net_c`+ imp.Name + `(args ...interface{}) (d ` + imp.Struct +`) {
					if len(args) > 0 {
					var jsonBlob = []byte(args[0].(string))
					err := json.Unmarshal(jsonBlob, &d)
					if err != nil {
						fmt.Println("error:", err)
						return 
					}
					} else {
						d = `+ imp.Struct +`{}
					}
    				return
				}`	    
			}

     
			//Methods have been added


			local_string += `
			func dummy_timer(){
				dg := time.Second *5
				fmt.Println(dg)
			}`

		
			local_string += `func main() {
				` + template.Main
				
					
					 local_string += `
					 ` + timeline +`
					 fmt.Printf("Listenning on Port %v\n", "` + template.Port +`")
					 http.HandleFunc( "/",  makeHandler(handler))
					 http.Handle("/dist/",  http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: "` + web +`"}))
					 http.ListenAndServe(":`+ template.Port +`", nil)
					 }`
					
		
			
				fmt.Println("Saving file to " + r + "/" + template.Output)
				 d1 := []byte(local_string)
             _ = ioutil.WriteFile(r + "/" + template.Output, d1,0644)
		    
		}

	} else if template.Type == "bind" {
		local_string = `package ` + template.Package + ` 
				
	import (`

				net_imports := []string{"time","os","bytes","encoding/json","fmt","html","html/template","io/ioutil","strings","reflect","unsafe","crypto/aes","crypto/cipher","crypto/rand","io","encoding/base64","errors"}

		/*
			Methods before so that we can create to correct delegate method for each object
		*/
		
		for _,imp := range template.Methods.Methods {
			if !contains(available_methods, imp.Name) {
				available_methods = append(available_methods, imp.Name)
			}
		}
		apiraw := ``
		for _,imp := range template.Endpoints.Endpoints {
			if !contains(api_methods, imp.Method) {
				api_methods = append(api_methods, imp.Method)
			}
				meth := template.findMethod(imp.Method)
				apiraw += ` 
				if  path == "` + imp.Path +`" && method == strings.ToUpper("` + imp.Type +`") { 
					` + meth.Method + `
					callmet = true
				}
				` 
			
		}
	
		fmt.Printf("APi Methods %v\n",api_methods)
		     netMa := 	`template.FuncMap{"js" : net_importjs,"css" : net_importcss,"sDelete" : deleteSession,"sRemove" : net_RemoveSessionKey,"sExist": net_SessionKeyExists,"sSet" : net_SetSessionKey,"sSetField": net_SetSessionField,"sGet" : net_GetSession,"sGetString" : net_GetSessionString, "sGetN" : net_GetSessionFloat,"Get" : paramGet,"eq": equalz, "neq" : nequalz, "lte" : netlt`
           for _,imp := range available_methods {
           	if !contains(api_methods, imp) {
          		netMa += `,"` + imp + `" : net_` + imp
      		}
           }
           int_lok := []string{}

           for _,imp := range template.Header.Objects {
			//struct return and function
				
			if !contains(int_lok, imp.Name) {
				int_lok = append(int_lok,imp.Name)
			 	netMa += `,"` + imp.Name + `" : net_` + imp.Name
			}
			}


           for _,imp := range template.Templates.Templates {

				netMa += `,"` + imp.Name + `" : net_` + imp.Name
				netMa += `,"b` + imp.Name + `" : net_b` + imp.Name
				netMa += `,"c` + imp.Name + `" : net_c` + imp.Name
           }
           netMa += `}`

		for _,imp := range template.RootImports {
				//fmt.Println(imp)
			if !strings.Contains(imp.Src,".xml") {
					if imp.Download == "true" {
						fmt.Println("∑ Downloading Package " + imp.Src)
						RunCmd("go get " + imp.Src)
					}
					if  !contains(net_imports, imp.Src) {
						net_imports = append(net_imports, imp.Src)
					}
			}
		}

		fmt.Println(template.Methods.Methods[0].Name)

		for _,imp := range net_imports {
			local_string += `
			"` + imp + `"`
		}
		local_string += `
		)
				
				var key = []byte("` + template.Key +`")

				func net_importcss(s string) string {
					return "<link rel=\"stylesheet\" href=\"" + s + "\" /> "
				}

				func net_importjs(s string) string {
					return "<script type=\"text/javascript\" src=\"" + s + "\" ></script> "
				}

			
					 type page struct {
					    Title string
					    Body  []byte
					 	Parameters map[string]interface{}
					    isResource bool
					}

				type session struct {
					Values map[string]interface{}
					//custom props
					` + template.Session + `
				
				}

				func paramGet(ke string,f map[string]interface{}) string {
					return f[ke].(string)
				}
			
				func dummy_timer(){
					dg := time.Second *5
					fmt.Println(dg)
				}

				func LoadUrl(path string,bod []byte,method string)[]byte { 
								
				body := new(bytes.Buffer)
				body.Write(bod)
				var f interface{}
				if bod != nil {
				_ = json.Unmarshal(bod, &f)
				}
				data,proceed := apiAttempt(path,method,bod)				
				if proceed {
					return data
				} else {

								 p,err := loadPage(path)
								  if err != nil {
								  	fmt.Println(err)
								        return []byte("Error ")
								  }

								  if !p.isResource {
								      p.Parameters = f.(map[string]interface{}) 
								      return   []byte(html.UnescapeString(string(renderTemplate("web" + path, p))))
								  } else {
								       return p.Body
								  }

					return bod
				}
								 
				}

				func net_SetSessionField(key string, arg interface{}) string {
					s := openSessionMap()
					s[key] = arg
					keepSessionMap(s)
					return ""
				}
				func net_SetSessionKey(key string, arg interface{}) string {
					s := openSession()
					s.Values[key] = arg
					keepSession(s)
					return ""
				}

				func net_SessionKeyExists(key string) bool {
					s := openSession()
					 if _, ok := s.Values[key]; ok {
					    //do something here
				 		return true
					}

					return false
				}

				func net_GetSession(key string) interface{} {
					s := openSession() 
					return s.Values[key]
				}
				func net_GetSessionString(key string) string {
					s := openSession() 
					return s.Values[key].(string)
				}
				func net_GetSessionFloat(key string) float64 {
					s := openSession() 
					return s.Values[key].(float64)
				}

				func net_RemoveSessionField(key string) string {
					s := openSessionMap()
					delete(s,key)
					//save here
					keepSessionMap(s)
					return ""
				}

				func net_RemoveSessionKey(key string) string {
					s := openSession()
					delete(s.Values,key)
					//save here
					keepSession(s)
					return ""
				}

				func deleteSession() string {
					os.Remove(os.TempDir() + "/session")
					return ""
				}

				func encrypt(text []byte) ([]byte, error) {
				    block, err := aes.NewCipher(key)
				    if err != nil {
				        return nil, err
				    }
				    b := base64.StdEncoding.EncodeToString(text)
				    ciphertext := make([]byte, aes.BlockSize+len(b))
				    iv := ciphertext[:aes.BlockSize]
				    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
				        return nil, err
				    }
				    cfb := cipher.NewCFBEncrypter(block, iv)
				    cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
				    return ciphertext, nil
				}

				func decrypt(text []byte) ([]byte, error) {
				    block, err := aes.NewCipher(key)
				    if err != nil {
				        return nil, err
				    }
				    if len(text) < aes.BlockSize {
				        return nil, errors.New("ciphertext too short")
				    }
				    iv := text[:aes.BlockSize]
				    text = text[aes.BlockSize:]
				    cfb := cipher.NewCFBDecrypter(block, iv)
				    cfb.XORKeyStream(text, text)
				    data, err := base64.StdEncoding.DecodeString(string(text))
				    if err != nil {
				        return nil, err
				    }
				    return data, nil
				}


				func openSession() session {
				  body, err := ioutil.ReadFile(os.TempDir() + "/session")
    				if err != nil {
    						s := session{Values:make(map[string]interface{})}
    						return s
    				}
    				var d session
    				data,_ := decrypt(body)
    				err = json.Unmarshal(data, &d)
					if err != nil {
						fmt.Println("error:", err)
						return session{}
					}
					return d
				}

				func openSessionMap() map[string]interface{} {
				  body, err := ioutil.ReadFile(os.TempDir() + "/session")
    				if err != nil {
    						s := make(map[string]interface{})
    						return s
    				}
    				var d interface{}
    				data,_ := decrypt(body)
    				err = json.Unmarshal(data, &d)
					if err != nil {
						fmt.Println("error:", err)
						return make(map[string]interface{})
					}
					return d.(map[string]interface{})
				}

			

				
				func keepSession(s session){
				
					data,er := encrypt([]byte(mResponse(s)))
					if er != nil {
						fmt.Println(er)
						return
					}
					err := ioutil.WriteFile(os.TempDir() + "/session", data,0644)
					if err != nil {
						fmt.Println(err)
					}
				}

					func keepSessionMap(s interface{}){
					fmt.Println(mResponse(s))
					data,er := encrypt([]byte(mResponse(s)))
					if er != nil {
						fmt.Println(er)
						return
					}
					err := ioutil.WriteFile(os.TempDir() + "/session", data,0644)
					if err != nil {
						fmt.Println(err)
					}
				}


				func renderTemplate(tmpl string, f*page) []byte {
				   filename :=  tmpl  + ".tmpl"
				   body, err := Asset(filename)
				   outp := new(bytes.Buffer)
				    if err != nil {
				       fmt.Print(err)
				    } else {
				    t := template.New("PageWrapper")
				    t = t.Funcs(` + netMa + `)
				      t, _ = t.Parse(strings.Replace(strings.Replace(strings.Replace(BytesToString(body), "/{", "\"{",-1),"}/", "}\"",-1 ) ,"` + "`" + `", ` + "`" + `\"` + "`" +` ,-1) )
				   
				    error := t.Execute(outp, f)
				    if error != nil {
				    fmt.Print(error)
				    return nil
				    } 

				   return outp.Bytes()

				    
				    }
				    return outp.Bytes()
				}

				func loadPage(title string) (*page,error) {
				    filename :=  "web" + title + ".tmpl"
				    if title == "/" {
				    	filename = "web/index.tmpl"
				    }
				    body, err := Asset(filename)
				    if err != nil {
				      filename = "web" + title + ".html"
				      if title == "/" {
				    	filename = "web/index.html"
				    	}
				      body, err = Asset(filename)
				      if err != nil {
				         filename = "web" + title
				         body, err = Asset(filename)
				         if err != nil {
				            return nil, err
				         } else {
				          if strings.Contains(title, ".tmpl") || title == "/" {
				              return nil,nil
				          }
				          return &page{Title: title, Body: body,isResource: true}, nil
				         }
				      } else {
				         return &page{Title: title, Body: body,isResource: true}, nil
				      }
				    } 
				    //load custom struts
				    return &page{Title: title, Body: body,isResource:false}, nil
				}
				func apiAttempt(path string, method string,bod []byte) ([]byte,bool) {
				//	session, er := store.Get(r, "session-")
					response := ""
					session := openSession()
					callmet := false
					var f interface{}
					if bod != nil {
					_ = json.Unmarshal(bod, &f)
					}

					` + apiraw + `
				

					if callmet {
						keepSession(session)
						
						if response != "" {
							
							return []byte(response),true
						}
					
					}
					return []byte(""),false
				} 


			
				func mResponse(v interface{}) string {
					data,_ := json.Marshal(&v)
					return string(data)
				}
				
			

			
				func BytesToString(b []byte) string {
				    bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
				    sh := reflect.StringHeader{bh.Data, bh.Len}
				    return *(*string)(unsafe.Pointer(&sh))
				}
				func equalz(args ...interface{}) bool {
		    	    if args[0] == args[1] {
		        	return true;
				    }
				    return false;
				 }
				 func nequalz(args ...interface{}) bool {
				    if args[0] != args[1] {
				        return true;
				    }
				    return false;
				 }

				 func netlt(x,v float64) bool {
				    if x < v {
				        return true;
				    }
				    return false;
				 }
				 func netgt(x,v float64) bool {
				    if x > v {
				        return true;
				    }
				    return false;
				 }
				 func netlte(x,v float64) bool {
				    if x <= v {
				        return true;
				    }
				    return false;
				 }
				 func netgte(x,v float64) bool {
				    if x >= v {
				        return true;
				    }
				    return false;
				 }
				`
					for _,imp := range template.Variables {
						local_string += `
						var ` + imp.Name + ` ` + imp.Type 
					}
		if template.Init_Func != "" {
			local_string += `
			func init(){
				` + template.Init_Func + `
			}`

		}

		//Lets Do structs
		for _,imp := range template.Header.Structs {
			if !contains(arch.objects, imp.Name) {
			fmt.Println("Processing Struct : " + imp.Name)
			arch.objects = append(arch.objects, imp.Name)
			local_string += `
			type ` + imp.Name + ` struct {`
			local_string += imp.Attributes
			local_string += `
			}`
			}
		}
		


		for _,imp := range template.Header.Objects {
			local_string += `
			type ` + imp.Name + ` ` + imp.Templ
		}
		

		//Create an object map
		for _,imp := range template.Header.Objects {
			//struct return and function
			fmt.Println("∑ Processing object :" + imp.Name)
			if !contains(available_methods, imp.Name) {
				//addcontructor
				available_methods = append(available_methods,imp.Name)
				int_methods = append(int_methods,imp.Name)
				local_string += `
				func  net_`+ imp.Name + `(args ...interface{}) (d ` + imp.Templ +`){
					if len(args) > 0 {
					jso := args[0].(string)
					var jsonBlob = []byte(jso)
					err := json.Unmarshal(jsonBlob, &d)
					if err != nil {
						fmt.Println("error:", err)
						return
					}
					return
					} else {
						d = ` + imp.Templ +`{} 
						return
					}
				}`	    

			}

			delegateMethods := strings.Split(imp.Methods,"\n")

			for _,im := range delegateMethods {
				
				if stripSpaces(im) != "" {
				fmt.Println(imp.Name + "->" + im)
				function_map := strings.Split(im, ")")

				

				if !contains(int_mappings, function_map[0] + imp.Templ) {
					int_mappings = append(int_mappings,function_map[0] + imp.Templ)
					funcsp := strings.Split(function_map[0],"(")
					meth := template.findMethod(stripSpaces(funcsp[0]))

					//process limits and keep local deritives
					if meth.Autoface == "" || meth.Autoface == "true"  {
						
						/*
							
						*/
						procc_funcs := true
						fmt.Println( )

						if meth.Limit != "" {
							if !contains(strings.Split(meth.Limit,","), imp.Name ){
								procc_funcs = false 
							}
						}

						if contains(api_methods, meth.Name){
							procc_funcs = false
						}
						
						objectName := meth.Object
						if objectName == "" {
							objectName = "object"
						}
						if procc_funcs {	
							if !contains(int_methods, stripSpaces(funcsp[0])) && meth.Name != "000" {
						int_methods = append(int_methods,stripSpaces(funcsp[0]))
						}
					  	local_string += `
					  	func  net_`+ stripSpaces(funcsp[0]) + `(` + strings.Trim(funcsp[1] + `, ` + objectName + ` ` +imp.Templ, ",") +`) ` + stripSpaces(function_map[1])
						if stripSpaces(function_map[1]) == "" {
							local_string += ` string`
						}

						local_string += ` {
									` + meth.Method

						if stripSpaces(function_map[1]) == "" {
							local_string += ` 
								return ""
							`
						}
						local_string += ` 
						}`



						if meth.Keeplocal == "false" || meth.Keeplocal == "" {
						local_string += `
						func (` + objectName + ` ` + imp.Templ +`) ` +  stripSpaces(funcsp[0]) + `(` + strings.Trim(funcsp[1], ",") +`) ` + stripSpaces(function_map[1])

							local_string += ` {
							` + meth.Method

						local_string +=  `
						}`
						}
						}
					}





				}
				}
			}

			//create Unused methods methods
			fmt.Println(int_methods)
			for _,imp := range available_methods {
				if !contains(int_methods,imp) && !contains(api_methods, imp)  {
					fmt.Println("Processing : " + imp)
						meth := template.findMethod(imp)
						addedit := false
						if meth.Returntype == "" {
							meth.Returntype = "string"
							addedit = true
						}
						local_string += `
						func net_` + meth.Name + `(args ...interface{}) ` + meth.Returntype + ` {
							`
						for k,nam := range strings.Split(meth.Variables,","){
							if nam != "" {
								local_string +=  nam + ` := ` + `args[` + strconv.Itoa(k) + `]
								`
							}
						}
						local_string += meth.Method
						if addedit {
						 local_string +=  `
						 return ""
						 `
						}
						local_string += `
						}` 
					}
			}
					for _,imp := range template.Templates.Templates {
				local_string += `
				func  net_`+ imp.Name + `(args ...interface{}) string {
					var d ` + imp.Struct + `
					if len(args) > 0 {
					jso := args[0].(string)
					var jsonBlob = []byte(jso)
					err := json.Unmarshal(jsonBlob, &d)
					if err != nil {
						fmt.Println("error:", err)
						return ""
					}
					} else {
						d = ` + imp.Struct +`{}
					}

					filename :=  "` + tmpl + `/` + imp.TemplateFile + `.tmpl"
    				body, er := Asset(filename)
    				if er != nil {
    					return ""
    				}
    				 output := new(bytes.Buffer) 
					t := template.New("` +  imp.Name + `")
    				t = t.Funcs(` + netMa +`)
				  	t, _ = t.Parse(strings.Replace(strings.Replace(strings.Replace(BytesToString(body), "/{", "\"{",-1),"}/", "}\"",-1 ) ,"` + "`" + `", ` + "`" + `\"` + "`" +` ,-1) )
			
				    error := t.Execute(output, &d)
				    if error != nil {
				    fmt.Print(error)
				    } 
					return html.UnescapeString(output.String())
				}`	    
					local_string += `
				func  net_b`+ imp.Name + `(d ` + imp.Struct +`) string {
					filename :=  "` + tmpl + `/` + imp.TemplateFile + `.tmpl"
    				body, er := Asset(filename)
    				if er != nil {
    					return ""
    				}
    				 output := new(bytes.Buffer) 
					t := template.New("` +  imp.Name + `")
    				t = t.Funcs(` + netMa +`)
				  	t, _ = t.Parse(strings.Replace(strings.Replace(strings.Replace(BytesToString(body), "/{", "\"{",-1),"}/", "}\"",-1 ) ,"` + "`" + `", ` + "`" + `\"` + "`" +` ,-1) )
			
				    error := t.Execute(output, &d)
				    if error != nil {
				    fmt.Print(error)
				    } 
					return html.UnescapeString(output.String())
				}`	    
				local_string += `
				func  net_c`+ imp.Name + `(args ...interface{}) (d ` + imp.Struct +`) {
					if len(args) > 0 {
					var jsonBlob = []byte(args[0].(string))
					err := json.Unmarshal(jsonBlob, &d)
					if err != nil {
						fmt.Println("error:", err)
						return 
					}
					} else {
						d = `+ imp.Struct +`{}
					}
    				return
				}`	    
			}

     
			//Methods have been added


			

	



			
				fmt.Println("Saving file to " + r + "/" + template.Output)
				 d1 := []byte(local_string)
             _ = ioutil.WriteFile(r + "/" + template.Output, d1,0644)
		    
		}
	}


	return
}

func RunFile(root string,file string){
	fmt.Println("∑ Running " + root + "/" + file );
	exe_cmd("go run " + root + "/" +  file )
}

func RunCmd(cmd string){
	exe_cmd(cmd)
}

func exe_cmd(cmd string) {
  
    parts := strings.Fields(cmd)
  	fmt.Println(cmd)
    var out *exec.Cmd
    if len(parts) == 5 {
    	fmt.Println("Match")
    	out = exec.Command(parts[0],parts[1],parts[2],parts[3],parts[4])	
    } else if len(parts) == 4 {
    		out = exec.Command(parts[0],parts[1],parts[2],parts[3])	
    }else if len(parts) > 2 {
    	out = exec.Command(parts[0],parts[1],parts[2])
	} else if len(parts) == 1 {
		out = exec.Command(parts[0])
	} else {
		out = exec.Command(parts[0],parts[1])
	}
    stdout, err := out.StdoutPipe()
    if err != nil {
        fmt.Println("error occured")
        fmt.Printf("%s", err)
    }
    out.Start()
	r := bufio.NewReader(stdout)
	t := false
	for !t {
	if out.ProcessState != nil {
	if out.ProcessState.Exited() {
		t = true
		fmt.Println("Ω Sub Process finished")
	}
	} else {
		t = true
		fmt.Println("Ω Sub Process finished")
	}
	line, _, _ := r.ReadLine()
	if string(line) != "" {
    fmt.Println(string(line))
	}
    }
}

func Exe_Stall(cmd string) {
    fmt.Println(cmd)
    parts := strings.Fields(cmd)
    var out *exec.Cmd
    if len(parts) > 2 {
    	out = exec.Command(parts[0],parts[1],parts[2])
	} else if len(parts) == 1 {
		out = exec.Command(parts[0])
	} else {
		out = exec.Command(parts[0],parts[1])
	}
    stdout, err := out.StdoutPipe()
    if err != nil {
        fmt.Println("error occured")
        fmt.Printf("%s", err)
    }
    out.Start()
	r := bufio.NewReader(stdout)
	t := false
	for !t {
	line, _, _ := r.ReadLine()
	if string(line) != "" {
    fmt.Println(string(line))
	}
    }
}

func stripSpaces(str string) string {
    return strings.Map(func(r rune) rune {
        if unicode.IsSpace(r) {
            // if the character is a space, drop it
            return -1
        }
        // else keep it in the string
        return r
    }, str)
}


func (d*gos) findStruct(name string) Struct {
	for _, imp := range d.Header.Structs {
		if imp.Name == name {
			return imp
		}
	}
	return Struct{Name:"000"}
}

func (d*gos) findMethod(name string) Method {
	for _, imp := range d.Methods.Methods {
		if imp.Name == name {
			return imp
		}
	}
	return Method{Name:"000"}
}

func LoadGos(path string) (*gos,*Error) {
	fmt.Println("∑ loading " + path)
	v := &gos{}
	 body, err := ioutil.ReadFile(path)
    if err != nil {
    	return nil, &Error{code: 404,reason:"file not found! @ " + path}
    }

 	//obj := Error{}
 	//fmt.Println(obj);
    d := xml.NewDecoder(bytes.NewReader(body))
    d.Entity = map[string]string{
        "&": "&",
    }
    err = d.Decode(&v)
    if err != nil {
        fmt.Printf("error: %v", err)
        return nil,nil
    }
   	//process mergs
   	for _,imp := range v.RootImports {
   		//fmt.Println(imp.Src)
   		if strings.Contains(imp.Src,".xml") {
   			v.MergeWith(GOHOME + "/" + strings.Trim(imp.Src,"/"))
   		}
   	}

    return v,nil
}

func (d*gos) MergeWith(target string) {
	fmt.Println("∑ Merging " + target)
    imp,err := LoadGos(target)
    if err != nil {
    	fmt.Println(err)
    } else {
    
    for _,im := range imp.RootImports {
   	if strings.Contains(im.Src,".xml") {
   			imp.MergeWith(GOHOME + "/" + strings.Trim(im.Src,"/"))
   	}
   }

    d.RootImports = append(imp.RootImports,d.RootImports...)
    d.Header.Structs = append(imp.Header.Structs, d.Header.Structs...)
    d.Header.Objects = append(imp.Header.Objects, d.Header.Objects...)
    d.Methods.Methods = append(imp.Methods.Methods, d.Methods.Methods...)
    d.Timers.Timers = append(imp.Timers.Timers, d.Timers.Timers...)
    d.Templates.Templates = append(imp.Templates.Templates, d.Templates.Templates...)
    d.Endpoints.Endpoints = append(imp.Endpoints.Endpoints,d.Endpoints.Endpoints...)
	}			
}

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}



func DoubleInput(p1 string,p2 string) (r1 string, r2 string) {
    fmt.Println(p1)
    fmt.Scanln(&r1)
    fmt.Println(p2)
    fmt.Scanln(&r2)
    return
}


func AskForConfirmation() bool {
	var response string
		fmt.Println("Please type yes or no and then press enter:")
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return AskForConfirmation()
	}
}


func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

func getPy() string {
	return ``
}

func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}