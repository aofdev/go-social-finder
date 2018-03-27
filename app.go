package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

//access token social
var facebookAccessToken = ""
var igAccessToken = ""

type Facebook struct {
	Name string `json:"name"`
	Page string `json:"id"`
}

type Instagram struct {
	Name string `json:"username"`
	Page string `json:"Id"`
}

var data string
var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {

	http.HandleFunc("/", index)
	http.ListenAndServe(":8000", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalln(err)
	}
	if r.Method == "GET" {
		home := "home"
		tpl.ExecuteTemplate(w, "home.gohtml", home)
	} else {
		r.ParseForm()
		str := cleanLink(r.FormValue("link"))
		if strings.Contains(str, "facebook") {
			fmt.Println("link facebook: ", str)
			data = getFacebookPageId(str, facebookAccessToken)
			fmt.Println("id facebook: ", data)
			facebook_id := struct {
				fb_id string
			}{
				data,
			}
			tpl.ExecuteTemplate(w, "home.gohtml", facebook_id)
		} else if strings.Contains(str, "instagram") {
			fmt.Println("link instagram: ", str)
			data = getInstagramPageId(str, igAccessToken)
			fmt.Println("id instagram: ", data)
			instagram_id := struct {
				ig_id string
			}{
				data,
			}
			tpl.ExecuteTemplate(w, "home.gohtml", instagram_id)
		} else if strings.Contains(str, "youtube") {
			fmt.Println("link youtube: ", str)
		} else if strings.Contains(str, "twitter") {
			fmt.Println("link twitter: ", str)
		} else {
			fmt.Println("not found")
			not := struct {
				Not string
			}{
				"not found",
			}
			tpl.ExecuteTemplate(w, "notfound.gohtml", not)
		}

	}

}

func getFacebookPageId(link string, accessToken string) string {
	var getPage string
	page := strings.Replace(link, "https://www.facebook.com/", "", -1)
	pageName := strings.Replace(page, "/", "", -1)
	req := "https://graph.facebook.com/v2.11/" + pageName + "?access_token=" + accessToken
	jsonData := myRequest(req)
	res := Facebook{}
	json.Unmarshal([]byte(jsonData), &res)
	getPage = res.Page
	return getPage
}

func getInstagramPageId(link string, accessToken string) string {
	var getPage string
	page := strings.Replace(link, "https://www.instagram.com/", "", -1)
	pageName := strings.Replace(page, "/", "", -1)
	req := "https://api.instagram.com/v1/users/search?q=" + pageName + "&count=1&access_token=" + accessToken
	jsonData := myRequest(req)
	res := Instagram{}
	json.Unmarshal([]byte(jsonData), &res)
	fmt.Println(jsonData)
	getPage = res.Page
	return getPage
}

func myRequest(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	return string(body)
}

func cleanLink(link string) string {
	t := strings.TrimSpace(link)
	var re = regexp.MustCompile(`\?.*`)
	s := re.ReplaceAllString(t, "$1")
	return s
}
