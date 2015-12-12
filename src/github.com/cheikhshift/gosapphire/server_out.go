package main 
import (
			"net/http"
			"github.com/elazarl/go-bindata-assetfs"
			"time"
			"github.com/gorilla/sessions"
			"bytes"
			"encoding/json"
			"fmt"
			"html"
			"html/template"
			"strings"
			"reflect"
			"unsafe"
		)
				var store = sessions.NewCookieStore([]byte("something-secretive-is-what-a-gorrilla-needs"))

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
				    t = t.Funcs(template.FuncMap{"js" : net_importjs,"css" : net_importcss,"sd" : net_sessionDelete,"sr" : net_sessionRemove,"sc": net_sessionKey,"ss" : net_sessionSet,"sso": net_sessionSetInt,"sgo" : net_sessionGetInt,"sg" : net_sessionGet,"form" : formval,"eq": equalz, "neq" : nequalz, "lte" : netlt,"WhatsMyAttrLength" : net_WhatsMyAttrLength,"sendEmail" : net_sendEmail,"WhatsMyAttr" : net_WhatsMyAttr,"Button" : net_Button,"bButton" : net_bButton,"cButton" : net_cButton,"Bootstrap_alert" : net_Bootstrap_alert,"bBootstrap_alert" : net_bBootstrap_alert,"cBootstrap_alert" : net_cBootstrap_alert})
				    t, _ = t.Parse(strings.Replace(strings.Replace(strings.Replace(BytesToString(body), "/{", "\"{",-1),"}/", "}\"",-1 ) ,"`", `\"` ,-1) )
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

					 
				if  r.URL.Path == "/index/api" && r.Method == strings.ToUpper("POST") { 
					 
			//login function
			response = mResponse(Button{Color:"#fff"})
			fmt.Println("Login!! -> " + session.Values["username"].(string))
		
					callmet = true
				}
				

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
				        renderTemplate(w, r,  "web" + r.URL.Path, p)
				  } else {
				       w.Write(p.Body)
				  }
				}

				func loadPage(title string, servlet string,r *http.Request) (*Page,error) {
				    filename :=  "web" + title + ".tmpl"
				    body, err := Asset(filename)
				    if err != nil {
				      filename = "web" + title + ".html"
				      body, err = Asset(filename)
				      if err != nil {
				         filename = "web" + title
				         body, err = Asset(filename)
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
					    R *http.Request
					    Session *sessions.Session
					}
						var FreeServer string
			func init(){
				
		fmt.Println("Not up yet....?\n")
	
			}
			type DemoChild struct {
			SomeOtherAttr string
		
			}
			type Basic_response struct {
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
			type myDemoObject DemoGos
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
								
			   	go fmt.Println("Send Email -> " + to.(string) + " ->" + from.(string))
				return true
		
						}
						func net_WhatsMyAttr(args ...interface{}) string {
							save := args[0]
								end := args[1]
								 
			gb := save.(string) + end.(string)
			fmt.Printf(gb)
		
						 return ""
						 
						}
				func  net_Button(args ...interface{}) string {
					var d Button
					if len(args) > 0 {
					jso := args[0].(string)
					var jsonBlob = []byte(jso)
					err := json.Unmarshal(jsonBlob, &d)
					if err != nil {
						fmt.Println("error:", err)
						return ""
					}
					} else {
						d = Button{}
					}

					filename :=  "tmpl/bootstrap/button.tmpl"
    				body, er := Asset(filename)
    				if er != nil {
    					return ""
    				}
    				 output := new(bytes.Buffer) 
					t := template.New("Button")
    				t = t.Funcs(template.FuncMap{"js" : net_importjs,"css" : net_importcss,"sd" : net_sessionDelete,"sr" : net_sessionRemove,"sc": net_sessionKey,"ss" : net_sessionSet,"sso": net_sessionSetInt,"sgo" : net_sessionGetInt,"sg" : net_sessionGet,"form" : formval,"eq": equalz, "neq" : nequalz, "lte" : netlt,"WhatsMyAttrLength" : net_WhatsMyAttrLength,"sendEmail" : net_sendEmail,"WhatsMyAttr" : net_WhatsMyAttr,"Button" : net_Button,"bButton" : net_bButton,"cButton" : net_cButton,"Bootstrap_alert" : net_Bootstrap_alert,"bBootstrap_alert" : net_bBootstrap_alert,"cBootstrap_alert" : net_cBootstrap_alert})
				  	t, _ = t.Parse(strings.Replace(strings.Replace(strings.Replace(BytesToString(body), "/{", "\"{",-1),"}/", "}\"",-1 ) ,"`", `\"` ,-1) )
			
				    error := t.Execute(output, &d)
				    if error != nil {
				    fmt.Print(error)
				    } 
					return html.UnescapeString(output.String())
				}
				func  net_bButton(d Button) string {
					filename :=  "tmpl/bootstrap/button.tmpl"
    				body, er := Asset(filename)
    				if er != nil {
    					return ""
    				}
    				 output := new(bytes.Buffer) 
					t := template.New("Button")
    				t = t.Funcs(template.FuncMap{"js" : net_importjs,"css" : net_importcss,"sd" : net_sessionDelete,"sr" : net_sessionRemove,"sc": net_sessionKey,"ss" : net_sessionSet,"sso": net_sessionSetInt,"sgo" : net_sessionGetInt,"sg" : net_sessionGet,"form" : formval,"eq": equalz, "neq" : nequalz, "lte" : netlt,"WhatsMyAttrLength" : net_WhatsMyAttrLength,"sendEmail" : net_sendEmail,"WhatsMyAttr" : net_WhatsMyAttr,"Button" : net_Button,"bButton" : net_bButton,"cButton" : net_cButton,"Bootstrap_alert" : net_Bootstrap_alert,"bBootstrap_alert" : net_bBootstrap_alert,"cBootstrap_alert" : net_cBootstrap_alert})
				  	t, _ = t.Parse(strings.Replace(strings.Replace(strings.Replace(BytesToString(body), "/{", "\"{",-1),"}/", "}\"",-1 ) ,"`", `\"` ,-1) )
			
				    error := t.Execute(output, &d)
				    if error != nil {
				    fmt.Print(error)
				    } 
					return html.UnescapeString(output.String())
				}
				func  net_cButton(l string) (d Button) {
					
					
					var jsonBlob = []byte(l)
					err := json.Unmarshal(jsonBlob, &d)
					if err != nil {
						fmt.Println("error:", err)
						return 
					}
    				return
				}
				func  net_Bootstrap_alert(args ...interface{}) string {
					var d Bootstrap_alert
					if len(args) > 0 {
					jso := args[0].(string)
					var jsonBlob = []byte(jso)
					err := json.Unmarshal(jsonBlob, &d)
					if err != nil {
						fmt.Println("error:", err)
						return ""
					}
					} else {
						d = Bootstrap_alert{}
					}

					filename :=  "tmpl/bootstrap/alert.tmpl"
    				body, er := Asset(filename)
    				if er != nil {
    					return ""
    				}
    				 output := new(bytes.Buffer) 
					t := template.New("Bootstrap_alert")
    				t = t.Funcs(template.FuncMap{"js" : net_importjs,"css" : net_importcss,"sd" : net_sessionDelete,"sr" : net_sessionRemove,"sc": net_sessionKey,"ss" : net_sessionSet,"sso": net_sessionSetInt,"sgo" : net_sessionGetInt,"sg" : net_sessionGet,"form" : formval,"eq": equalz, "neq" : nequalz, "lte" : netlt,"WhatsMyAttrLength" : net_WhatsMyAttrLength,"sendEmail" : net_sendEmail,"WhatsMyAttr" : net_WhatsMyAttr,"Button" : net_Button,"bButton" : net_bButton,"cButton" : net_cButton,"Bootstrap_alert" : net_Bootstrap_alert,"bBootstrap_alert" : net_bBootstrap_alert,"cBootstrap_alert" : net_cBootstrap_alert})
				  	t, _ = t.Parse(strings.Replace(strings.Replace(strings.Replace(BytesToString(body), "/{", "\"{",-1),"}/", "}\"",-1 ) ,"`", `\"` ,-1) )
			
				    error := t.Execute(output, &d)
				    if error != nil {
				    fmt.Print(error)
				    } 
					return html.UnescapeString(output.String())
				}
				func  net_bBootstrap_alert(d Bootstrap_alert) string {
					filename :=  "tmpl/bootstrap/alert.tmpl"
    				body, er := Asset(filename)
    				if er != nil {
    					return ""
    				}
    				 output := new(bytes.Buffer) 
					t := template.New("Bootstrap_alert")
    				t = t.Funcs(template.FuncMap{"js" : net_importjs,"css" : net_importcss,"sd" : net_sessionDelete,"sr" : net_sessionRemove,"sc": net_sessionKey,"ss" : net_sessionSet,"sso": net_sessionSetInt,"sgo" : net_sessionGetInt,"sg" : net_sessionGet,"form" : formval,"eq": equalz, "neq" : nequalz, "lte" : netlt,"WhatsMyAttrLength" : net_WhatsMyAttrLength,"sendEmail" : net_sendEmail,"WhatsMyAttr" : net_WhatsMyAttr,"Button" : net_Button,"bButton" : net_bButton,"cButton" : net_cButton,"Bootstrap_alert" : net_Bootstrap_alert,"bBootstrap_alert" : net_bBootstrap_alert,"cBootstrap_alert" : net_cBootstrap_alert})
				  	t, _ = t.Parse(strings.Replace(strings.Replace(strings.Replace(BytesToString(body), "/{", "\"{",-1),"}/", "}\"",-1 ) ,"`", `\"` ,-1) )
			
				    error := t.Execute(output, &d)
				    if error != nil {
				    fmt.Print(error)
				    } 
					return html.UnescapeString(output.String())
				}
				func  net_cBootstrap_alert(l string) (d Bootstrap_alert) {
					
					
					var jsonBlob = []byte(l)
					err := json.Unmarshal(jsonBlob, &d)
					if err != nil {
						fmt.Println("error:", err)
						return 
					}
    				return
				}
			func main() {
				

	
					 
			PublicName := time.NewTicker(time.Second * 60)
					    go func() {
					        for _ = range PublicName.C {
					            
			//login function
			
			fmt.Println("Clean up resources")
		
					        }
					    }()
    
					 fmt.Printf("Listenning on Port %v\n", "8080")
					 http.HandleFunc( "/",  makeHandler(handler))
					 http.Handle("/dist/",  http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: "web"}))
					 http.ListenAndServe(":8080", nil)
			}