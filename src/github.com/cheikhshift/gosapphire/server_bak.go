package main

import (
   
    "net/http"
    "fmt"
    "io/ioutil"
    "encoding/xml"
    "html/template"
    "strings"
    "reflect"
    "unsafe"
)

var root string

//Full Page Declaration
type Page struct {
    Title string
    Body  []byte
    request *http.Request
    isResource bool
    s *map[string]interface{}
}

type server struct {
    XMLName    xml.Name `xml:"server"`
    Port string `xml:"port"`
    Package  string `xml:"package"`
}


type Demo struct {
  Para string
  other string
}

func New() *Demo {
  return &Demo{Para:"f"}
}

/*
auto - get ~~~
*/
func (s *Demo) canon() (string) {
  s.Para = "other"
  return ""
}
func canon() (string) {
 
  return ""
}
type vbn interface {
  canon() string
}


/*
  Definition of functions
*/

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



  func EmailDealWith(args ...interface{}) *Demo {
    return &Demo{Para: "Deugndeem"}
}

  func clog(args ...interface{}) error {
    object := args[0].(*Demo)
    object.Para = "deug"
    object.other = "srt"
    fmt.Printf("%v", object)
    return nil
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

func BytesToString(b []byte) string {
    bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
    sh := reflect.StringHeader{bh.Data, bh.Len}
    return *(*string)(unsafe.Pointer(&sh))
}


func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
     filename :=  tmpl + ".tmpl"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
       fmt.Print(err)
    } else {
    t := template.New("PageWrapper")
    t = t.Funcs(template.FuncMap{"New_DemoGos": EmailDealWith,"clog" : clog,"eq": equalz, "neq" : nequalz, "lte" : netlte  })
    t, _ = t.Parse(BytesToString(body))
    error := t.Execute(w, p)
    if error != nil {
    fmt.Print(error)
    } 
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

func main() {
  var i vbn
  i = New()
  //i = ki
  l := i.canon()
  fmt.Printf("%v,%v", i,l)

 fmt.Printf("Listenning on Port %v\n", "8080")
 http.HandleFunc( "/",  makeHandler(handler))
 http.ListenAndServe(":8080", nil)
}