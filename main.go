package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"personal-web/connection"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func main() {


	route := mux.NewRouter()
	connection.DatabaseConnect()

	route.PathPrefix("/public").Handler(http.StripPrefix("/public", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/",home).Methods("GET")
	route.HandleFunc("/home", home).Methods("GET")
	route.HandleFunc("/contact",contact).Methods("GET")
	route.HandleFunc("/project",project).Methods("GET")
	route.HandleFunc("/blog-detail/{id}",blogDetail).Methods("GET")
	route.HandleFunc("/form-project",AddProject).Methods("POST")
	route.HandleFunc("/form-contact",AddContact).Methods("POST")
	route.HandleFunc("/delete-blog/{id}",deleteBlog).Methods("GET")
	route.HandleFunc("/edit-project/{id}",editBlog).Methods("GET")
	route.HandleFunc("/submitedit/{id}",submitEdit).Methods("POST")
	// register
	route.HandleFunc("/form-register",formRegister).Methods("GET")
	route.HandleFunc("/submit-register",register).Methods("POST")
	// // login
	route.HandleFunc("/form-login",formLogin).Methods("GET")
	route.HandleFunc("/submit-login",login).Methods("POST")

	route.HandleFunc("/logout",logout).Methods("GET")
	

	fmt.Println("server running port 7000")
	http.ListenAndServe("localhost:7000",route)
}
type SessionData struct{
	IsLogin bool
	UserName string
	FlashData string
	
}
var Data = SessionData{}

func helloWorld(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Hello World"))
}
func home(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","text/html; charset=utf8")
	var tmpl, err = template.ParseFiles("home.html")

	if err != nil{
		w.Write([]byte("web tidak tersedia" + err.Error()))
		return
	}
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)
		for _, f1 := range fm {
			
			flashes = append(flashes, f1.(string))
		}
	}
	Data.FlashData = strings.Join(flashes, " ")

	data,err :=connection.Conn.Query(context.Background(),"SELECT id,name,description,duration FROM tb_projects")
	var result[]Project
	for data.Next(){
		var each = Project{}
		err:= data.Scan(&each.ID,&each.NamaProject,&each.Description,&each.Duration)
		if err != nil{
			fmt.Println(err.Error())
			return
		}
		result = append(result, each)
	}
	resData :=map[string]interface{}{
		"DataSession" : Data,
		"Blogs":result,
	}
	w.WriteHeader(http.StatusOK)

	tmpl.Execute(w,resData)
}
func contact(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","text/html; charset=utf8")
	var tmpl, err = template.ParseFiles("contact.html")

	if err != nil{
		w.Write([]byte("web tidak tersedia" + err.Error()))
		return
	}
	tmpl.Execute(w,nil)
}
func project(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","text/html; charset=utf8")
	var tmpl, err = template.ParseFiles("project.html")

	if err != nil{
		w.Write([]byte("web tidak tersedia" + err.Error()))
		return
	}
	tmpl.Execute(w,nil)
}
// var dataProject=[] Project{}



type Project struct{
	NamaProject string
	StartDate time.Time
	EndDate time.Time
	Description string
	NodeJs string
	VueJs string
	ReactJs string
	Java string
	Duration string
	ID int
	Format_Start_date string
	Format_End_date string
}
type User struct{
	ID  int
	Name string
	Email string
	Password string
}
func AddProject(w http.ResponseWriter,r *http.Request){
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	var namaProject = r.PostForm.Get("input-project")
	var startDate = r.PostForm.Get("input-start")
	var endDate = r.PostForm.Get("input-end")
	var description =r.PostForm.Get("input-description")
	// var nodeJs =	r.PostForm.Get("nodejs")
	// var vueJs = r.PostForm.Get("vuejs")
	// var reactJs = r.PostForm.Get("reactjs") 
	// var java = r.PostForm.Get("java")
	layout := "2006-01-02"
	startDateParse,_ := time.Parse(layout,startDate)
	endDateParse,_ := time.Parse(layout,endDate)

	hours := endDateParse.Sub(startDateParse).Hours()
	days := hours / 24
	weeks := math.Round(days / 7)
  	months := math.Round(days / 30)
 	years := math.Round(days / 365)

	var duration string
	

	if days >= 1 && days <= 6 {
		duration = strconv.Itoa(int(days)) + " days"
	} else if days >= 7 && days <= 29 {
		duration = strconv.Itoa(int(weeks)) + " weeks"
	} else if days >= 30 && days <= 364 {
		duration = strconv.Itoa(int(months)) + " months"
	} else if days >= 365 {
		duration = strconv.Itoa(int(years)) + " years"
	}

	_,err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_projects (name, description,start_date, end_date,duration) VALUES ($1, $2, $3, $4, $5)", namaProject, description, startDateParse, endDateParse, duration)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : "+ err.Error()))
		return
	}

	http.Redirect(w,r,"/home",http.StatusMovedPermanently)
}
func AddContact(w http.ResponseWriter,r *http.Request){
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Nama : " + r.PostForm.Get("input-nama"))
	fmt.Println("email : " + r.PostForm.Get("input-email"))
	fmt.Println("phone Number : " + r.PostForm.Get("input-phone"))
	fmt.Println("subject : " + r.PostForm.Get("input-subject"))
	fmt.Println("Description : " + r.PostForm.Get("input-description"))
	http.Redirect(w,r,"/home",http.StatusMovedPermanently)
}
func blogDetail(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("blog-detail.html")

	if err != nil {
		w.Write([]byte("message :" + err.Error()))
		return
	}
	var BlogDetail = Project{}
	id,_ := strconv.Atoi(mux.Vars(r)["id"])
	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, description,start_date,end_date,duration FROM tb_projects WHERE id = $1", id).Scan(&BlogDetail.ID, &BlogDetail.NamaProject, &BlogDetail.Description, &BlogDetail.StartDate, &BlogDetail.EndDate, &BlogDetail.Duration)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : "+ err.Error()))
		return
	}
	BlogDetail.Format_Start_date = BlogDetail.StartDate.Format("2 January 2006")
	BlogDetail.Format_End_date = BlogDetail.EndDate.Format("2 January 2006")

	data := map[string]interface{}{
		"Blog": BlogDetail,
	}
	tmpl.Execute(w,data)
}

func deleteBlog(w http.ResponseWriter,r *http.Request){
	id,_ := strconv.Atoi(mux.Vars(r)["id"])
	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id = $1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : "+ err.Error()))
		return
	}


	http.Redirect(w,r,"/home",http.StatusFound)
}
func editBlog(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("project-edit.html")
	if err != nil {
		w.Write([]byte("message :" + err.Error()))
		return
	}

	var BlogDetail = Project{}
	id,_ := strconv.Atoi(mux.Vars(r)["id"])
	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, description FROM tb_projects WHERE id = $1", id).Scan(&BlogDetail.ID, &BlogDetail.NamaProject, &BlogDetail.Description)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : "+ err.Error()))
		return
	}
	
	data := map[string]interface{}{
		"EDIT": BlogDetail,
	}
	tmpl.Execute(w,data)
}
func submitEdit(w http.ResponseWriter, r *http.Request){
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	id,_ := strconv.Atoi(mux.Vars(r)["id"])
	
	
	var namaProject = r.PostForm.Get("input-project")
	// var startDate = r.PostForm.Get("input-start")
	// var endDate = r.PostForm.Get("input-end")
	var description =r.PostForm.Get("input-description")
	// nodejs := r.PostForm.Get("nodejs")
	// golang := r.PostForm.Get("golang")
	// reactjs := r.PostForm.Get("reactjs")
	// vuejs := r.PostForm.Get("vuejs")

	// layout := "2006-01-02"
	// startDateParse,_ := time.Parse(layout,startDate)
	// endDateParse,_ := time.Parse(layout,endDate)

	// hours := endDateParse.Sub(startDateParse).Hours()
	// days := hours / 24
	// weeks := math.Round(days / 7)
  	// months := math.Round(days / 30)
 	// years := math.Round(days / 365)

	// var duration string
	

	// if years > 0{
	// 	duration = strconv.FormatFloat(years,'f',0,64) + "years"
	// }else if months > 0 {
	// 	duration = strconv.FormatFloat(months, 'f', 0, 64) + " Months"
	// }else if weeks > 0 {
	// 	duration = strconv.FormatFloat(weeks,'f',0,64) + "weeks"
	// } else if days > 0 {
	// 	duration = strconv.FormatFloat(days, 'f', 0, 64) + " Days"
	// } else if hours > 0 {
	// 	duration = strconv.FormatFloat(hours, 'f', 0, 64) + " Hours"
	// } else {
	// 	duration = "0 Days"
	// }
	_,err = connection.Conn.Exec(context.Background(), "UPDATE tb_projects SET name = $1, description = $2 WHERE id = $3", namaProject, description, id)



	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : "+ err.Error()))
		return
	}


	http.Redirect(w,r,"/home",http.StatusMovedPermanently)

}
func formRegister(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("form-register.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var name = r.PostForm.Get("inputName")
	var email = r.PostForm.Get("inputEmail")
	var password = r.PostForm.Get("inputPassword")

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	// fmt.Println(passwordHash)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_user(name, email, password) VALUES ($1, $2, $3)", name, email, passwordHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/form-login", http.StatusMovedPermanently)
}

func formLogin(w http.ResponseWriter, r *http.Request){
	w.Header().Set("contet-Type","text/html;charset=utf-8")
	var tmpl,err = template.ParseFiles("form-login.html")
	if err != nil{
		w.Write([]byte("message : "))
	}
	tmpl.Execute (w,nil)
}
func login(w http.ResponseWriter, r *http.Request){
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	var email = r.PostForm.Get("inputEmail")
	var password = r.PostForm.Get("inputPassword")

	user := User{}

	// mengambil data email, dan melakukan pengecekan email
	err = connection.Conn.QueryRow(context.Background(),
		"SELECT * FROM tb_user WHERE email=$1", email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {
		fmt.Println("Email belum terdaftar")
		http.Redirect(w, r, "/form-login", http.StatusMovedPermanently)
	
		return
	}

	// melakukan pengecekan password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println("Password salah")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message : Password salah " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	session.Values["Name"] = user.Name
	session.Values["Email"] = user.Email
	session.Values["IsLogin"] = true
	session.Options.MaxAge = 10800 // 3 JAM

	session.AddFlash("succesfull","message")
	session.Save(r,w)
	

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func logout (w http.ResponseWriter, r *http.Request){
	fmt.Println("logout")
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session,_ := store.Get(r,"SESSION_KEY")
	session.Options.MaxAge= -1
	session.Save(r,w)

	http.Redirect(w,r,"/form-login",http.StatusSeeOther)
}
