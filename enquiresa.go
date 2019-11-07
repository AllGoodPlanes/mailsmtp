package main

import(

	"net/http"
	"net/smtp"
	"fmt"
	"errors"
	"log"
	"html/template"
	"io"
	"os"
	"strings"
	"compress/gzip"

)
//to test using smtp to send emails from an app
var enquiresconfirmTemplate = template.Must(template.ParseGlob("templates/enquiresconfirm.html"))

var user string
var pass string

const STATIC_URL string = "/static/"
const STATIC_ROOT string = "static/"

func main() {

	mux := http.NewServeMux()

	fmt.Println("Listening...")
	mux.HandleFunc("/", makeGzipHandler(Enquires))
	fs:= http.FileServer(http.Dir("static"))
	mux.Handle("/static/",http.StripPrefix("/static/", fs))
	http.ListenAndServe(GetPort(), mux)
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

type Context struct {
	Title  string
	Static string
	User   string
}

func (w gzipResponseWriter) Write(b []byte) (int, error){
	return w.Writer.Write(b)
}

func makeGzipHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r*http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"),"gzip") {
			fn(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr:=gzipResponseWriter{Writer:gz, ResponseWriter: w}
		fn(gzr, r)
	}
}


func Home(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	context := Context{Title: "Soria Consulting"}
	context.Static = STATIC_URL
	t := template.Must(template.ParseGlob("templates/home.html"))
	err := t.ExecuteTemplate(w, "home", context)
	if err != nil {
		log.Print("template executing error: ", err)
	}

}

type loginAuth struct {
  username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unkown fromServer")
		}
	}
	return nil, nil
}

func Enquires(w http.ResponseWriter, req *http.Request){
	w.Header().Set("Content-Type", "text/html")
	context := Context{Title: "Contact Us"}
	fmt.Println("method:", req.Method)
	if req.Method == "GET"{
	render (w, "enquires", context)
} else {
	req.ParseForm()

	email := req.FormValue("email")
	telephone := req.FormValue("telephone")
	location := req.FormValue("location")
	subjectheading := req.FormValue("subjectheading")
	enquiry := req.FormValue("enquiry")

fmt.Println(email+telephone+location+subjectheading+enquiry)

	message := []byte("To:" + "trymriverman@gmail.com" + "\r\n" +
		"Subject:" + subjectheading + "\r\n" +
		"\r\n" +
		"email:" + email + "\r\n" +
		"\r\n" +
		"Telephone:" + telephone + "\r\n" +
		"\r\n" +
		"Location:" + location + "\r\n" +
		"\r\n" +
		"Enquiry details:" + enquiry )
to := []string{"trymriverman@gmail.com"}
 auth := LoginAuth("admin@gocloudcoding.com","hotdAsp0t")
 err := smtp.SendMail("smtp.office365.com:587", auth, "admin@gocloudcoding.com", to, message)

 if err != nil{
	 log.Fatal(err)
 }

enquiresconfirmTemplate.ExecuteTemplate(w, "enquiresconfirm.html", nil)
}
}

func render(w http.ResponseWriter, tmpl string, context Context) {
	context.Static = STATIC_URL
	tmpl_list := []string{"templates/base.html",
		fmt.Sprintf("templates/%s.html", tmpl),
	}
	t, err := template.ParseFiles(tmpl_list...)
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err= t.ExecuteTemplate(w, "base", context)
	if err != nil {
		log.Print("template executing error: ", err)
	}

}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Info: No port detected in the environment, defaulting to :" + port)

	return ":" + port
}


