package main

import (
	"encoding/json"
	"fmt"

	//"log"
	"strconv"

	//"io"
	"net/http"
	//"net/url"

	"RestfulAPI.com/Client"
	"github.com/gorilla/mux"
)

func GetProperties(w http.ResponseWriter, r *http.Request) {

	Properties := Client.GetProperties()

	json.NewEncoder(w).Encode(Properties)

}

func PostProperty(w http.ResponseWriter, r *http.Request) {

	var prop Client.Property

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&prop)

	fmt.Println(err)

	fmt.Println(prop)

	Client.InsertProperty(prop)

	w.WriteHeader(http.StatusCreated) // 201 Created

}

func GetPropertiesById(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	propid := params["propid"]
	intVar, _ := strconv.Atoi(propid)

	proper := Client.GetPropertiesById(intVar)

	if proper == (Client.Property{}) {
		http.Error(w, "Wrong propid++++++", http.StatusNotFound)

		w.WriteHeader(http.StatusNotFound)
	} else {
		json.NewEncoder(w).Encode(proper)
	}

}

func FilterProperties(w http.ResponseWriter, r *http.Request) {
	//From URL
	// params := mux.Vars(r)

	// params1 := params["City"]

	// params2 := params["Bedrooms"]

	// intVar, _ := strconv.Atoi(params2)

	// fmt.Println(params1,params2)

	// w.Header().Set("Content-Type", "application/json")

	//From Query Params

	param1 := r.URL.Query().Get("city")
	param2 := r.URL.Query().Get("bedrooms")

	intVar, _ := strconv.Atoi(param2)
	fmt.Println(param1, intVar)

	var Proper Client.Property
	c := Client.FilterProperties(param1, intVar)
	json.NewEncoder(w).Encode(c)

	fmt.Println(Proper)

}

func PutProperty(w http.ResponseWriter, r *http.Request) {

	var prop Client.Property

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&prop)

	fmt.Println(err)

	fmt.Println(prop)

	params := mux.Vars(r)

	propid := params["propid"]

	intVar, _ := strconv.Atoi(propid)
	Client.PutProperty(intVar, prop)

}

func DeleteProperty(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	propid := params["propid"]

	intVar, _ := strconv.Atoi(propid)

	Client.DeleteBooksById(intVar)

}

func BasicAuthMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		fmt.Println("username: ", user)
		fmt.Println("password: ", pass)
		if !ok || !checkUsernameAndPassword(user, pass) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password for this site"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}
		handler(w, r)
	}
}

func checkUsernameAndPassword(username, password string) bool {
	return username == "admin" && password == "admin@123"
}


func Getdocument(w http.ResponseWriter, r *http.Request) {
	
}	

func main() {

	Router := mux.NewRouter()

	Router.HandleFunc("/properties", BasicAuthMiddleware(GetProperties)).Methods("GET")
	//Router.HandleFunc("/properties", GetProperties).Methods("GET")

	Router.HandleFunc("/properties/{propid}", GetPropertiesById).Methods("GET")

	Router.HandleFunc("/properties1", FilterProperties).Methods("GET")
	//Router.HandleFunc("/properties/{City}/{Bedrooms}", FilterProperties).Methods("GET")  //from URL

	Router.HandleFunc("/properties", PostProperty).Methods("POST")

	Router.HandleFunc("/properties/{propid}", DeleteProperty).Methods("DELETE")

	Router.HandleFunc("/properties/{propid}", PutProperty).Methods("PUT")

	Router.HandleFunc("/ftdr/software/document-engine", Getdocument).Methods("GET")

	http.ListenAndServe(":9999", Router)
}
