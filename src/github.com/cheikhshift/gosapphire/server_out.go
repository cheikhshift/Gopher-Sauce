package main 
import (
			"net/http"
			"os"
			"bytes"
			"encoding/json"
			"fmt"
			"io/ioutil"
			"html"
			"html/template"
			"strings"
			"reflect"
			"unsafe"
		)
				func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
				     filename :=  tmpl  + ".tmpl"
				    body, err := ioutil.ReadFile(filename)
				    if err != nil {
				       fmt.Print(err)
				    } else {
				    t := template.New("PageWrapper")
				    t = t.Funcs(netMap)
				    t, _ = t.Parse(strings.Replace(strings.Replace(strings.Replace(BytesToString(body), "/{", "\"{",-1),"}/", "}\"",-1 ) ,"`", `\"` ,-1) )
				    outp := new(bytes.Buffer)
				    error := t.Execute(outp, p)
				    if error != nil {
				    fmt.Print(error)
				    } 
				    fmt.Fprintf(w, html.UnescapeString(outp.String()) )
				    }
				}

				func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
				  return func(w http.ResponseWriter, r *http.Request) {
				      fn(w, r, "")
				  }
				} 

				func handler(w http.ResponseWriter, r *http.Request, context string) {
				  // fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
				  p,err := loadPage(r.URL.Path , context,r)
				  if err != nil {
				        http.Error(w, err.Error(), http.StatusInternalServerError)
				        return
				  }

				  if !p.isResource {
				        renderTemplate(w,  "web" + r.URL.Path, p)
				  } else {
				       w.Write(p.Body)
				  }
				}

				func loadPage(title string, servlet string,r *http.Request) (*Page,error) {
				    filename :=  "web" + title + ".tmpl"
				    body, err := ioutil.ReadFile(filename)
				    if err != nil {
				      filename = "web" + title + ".html"
				      body, err = ioutil.ReadFile(filename)
				      if err != nil {
				         filename = "web" + title
				         body, err = ioutil.ReadFile(filename)
				         if err != nil {
				            return nil, err
				         } else {
				          if strings.Contains(title, ".tmpl") {
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
					    s *map[string]interface{}
					}
			func init(){
				
		fmt.Println("Not up yet....?")
	
			}
			type DemoChild struct {
			SomeOtherAttr string
		
			}
			type DemoGos struct {
			SomeAttr string
			Child *DemoChild
		
			}
			type Bootstrap_alert struct {
			Strong string
			Text string
			Type string
		
			}
			type Button struct {
			Color string
		
			}
				func  net_myDemoObject(jso string) (d DemoGos){
					var jsonBlob = []byte(jso)
					err := json.Unmarshal(jsonBlob, &d)
					if err != nil {
						fmt.Println("error:", err)
						return
					}
					return
				}
					  	func  net_Hackfmt(save string, object DemoGos)  string {
									 
								return ""
							 
						}
						func (object DemoGos) Hackfmt(save string)  {
							
						}
					  	func  net_WhatsMyAttrLength( object DemoGos) string {
									
				return ""
		 
						}
						func (object DemoGos) WhatsMyAttrLength() string {
							
				return ""
		
						}
						func net_sendEmail(args ...interface{}) bool {
							to := args[0]
								from := args[1]
								
			   	fmt.Println("Send Email -> " + to.(string) + " ->" + from.(string))
				return true
		
						}
						func net_WhatsMyAttr(args ...interface{}) string {
							save := args[0]
								end := args[1]
								 
			gb := save.(string) + end.(string)
			fmt.Printf(gb)
		
						 return ""
						 
						}
						func net_login(args ...interface{}) string {
							 
			
		
						 return ""
						 
						}
				func  net_Button(jso string) string {
					var d Button
					var jsonBlob = []byte(jso)
					err := json.Unmarshal(jsonBlob, &d)
					if err != nil {
						fmt.Println("error:", err)
						return ""
					}

					filename :=  "tmpl/button.tmpl"
    				body, er := ioutil.ReadFile(filename)
    				if er != nil {
    					return ""
    				}
    				 output := new(bytes.Buffer) 
					t := template.New("Button")
    				t = t.Funcs(netMa)
				    t, _ = t.Parse(BytesToString(body))

				    error := t.Execute(output, &d)
				    if error != nil {
				    fmt.Print(error)
				    } 
					return output.String()
				}
				func  net_Bootstrap_alert(jso string) string {
					var d Bootstrap_alert
					var jsonBlob = []byte(jso)
					err := json.Unmarshal(jsonBlob, &d)
					if err != nil {
						fmt.Println("error:", err)
						return ""
					}

					filename :=  "tmpl/Bootstrap/alert.tmpl"
    				body, er := ioutil.ReadFile(filename)
    				if er != nil {
    					return ""
    				}
    				 output := new(bytes.Buffer) 
					t := template.New("Bootstrap_alert")
    				t = t.Funcs(netMa)
				    t, _ = t.Parse(BytesToString(body))

				    error := t.Execute(output, &d)
				    if error != nil {
				    fmt.Print(error)
				    } 
					return output.String()
				}
           var netMa = template.FuncMap{"eq": equalz, "neq" : nequalz, "lte" : netlt,"WhatsMyAttrLength" : net_WhatsMyAttrLength,"sendEmail" : net_sendEmail,"WhatsMyAttr" : net_WhatsMyAttr,"login" : net_login,"myDemoObject" : net_myDemoObject}
           var netMap = template.FuncMap{"eq": equalz, "neq" : nequalz, "lte" : netlt,"WhatsMyAttrLength" : net_WhatsMyAttrLength,"sendEmail" : net_sendEmail,"WhatsMyAttr" : net_WhatsMyAttr,"login" : net_login,"myDemoObject" : net_myDemoObject,"Button" : net_Button,"Bootstrap_alert" : net_Bootstrap_alert}
			func main() {
				

	
					 os.Chdir("/Users/Adrian/gosapphire/src/github.com/cheikhshift/gosapphire")
					 fmt.Printf("Listenning on Port %v\n", "8080")
					 http.HandleFunc( "/",  makeHandler(handler))
					 http.Handle("/dist/", http.StripPrefix("", http.FileServer(http.Dir("web"))))
					 http.ListenAndServe(":8080", nil)
			}