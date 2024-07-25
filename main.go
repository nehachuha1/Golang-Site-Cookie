package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"math/rand"
	"net/http"
	"slices"
	"strings"
	"time"
)

type User struct {
	ID       uint8
	Username string
	Password string
}

type MyHandler struct {
	users    map[string]*User
	sessions []string
}

func NewMyHandler() *MyHandler {
	return &MyHandler{
		users:    map[string]*User{},
		sessions: []string{},
	}
}

func RemoveElement(sl []string, val string) []string {
	var valInd uint64
	if len(sl) > 0 {
		for i := 0; i < len(sl); i++ {
			if sl[i] == val {
				valInd = uint64(i)
				break
			}
		}
		sl[valInd] = sl[len(sl)-1]
		sl = sl[:len(sl)-1]
	}
	return sl
}

func RandomResponseID() string {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
	var returnSlice []string
	for i := 0; i < 12; i++ {
		n := rand.Intn(len(chars))
		returnSlice = append(returnSlice, string(chars[n]))
	}
	return strings.Join(returnSlice, "")
}

func (api *MyHandler) Login(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("AuthID")

	if !errors.Is(err, http.ErrNoCookie) && !slices.Contains(api.sessions, cookie.Value) {
		w.Write([]byte("There's some error while tryed to create cookie. Redirecting...\n"))
		time.Sleep(time.Second * 5)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	newCookie := &http.Cookie{
		Name:    "AuthID",
		Value:   RandomResponseID(),
		Expires: time.Now().Add(5 * time.Minute),
	}
	http.SetCookie(w, newCookie)
	api.sessions = append(api.sessions, newCookie.Value)
	w.Write([]byte("New cookie has successfully created!\n"))
	time.Sleep(3 * time.Second)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (api *MyHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("AuthID")
	if errors.Is(err, http.ErrNoCookie) || !slices.Contains(api.sessions, cookie.Value) {
		w.Write([]byte("There's no cookie AuthID\n"))
	}

	cookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, cookie)
	w.Write([]byte("You successfully logged out!\n"))
	api.sessions = RemoveElement(api.sessions, cookie.Value)

}

func (api *MyHandler) MainPage(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("AuthID")
	if errors.Is(err, http.ErrNoCookie) {
		w.Write([]byte("You are not registered\n"))
		time.Sleep(time.Second * 5)
		return
	}
	w.Write([]byte("You are currently logged in!\n"))
}

func main() {
	r := mux.NewRouter()

	api := NewMyHandler()
	r.HandleFunc("/", api.MainPage).Methods("GET")
	r.HandleFunc("/login", api.Login).Methods("GET")
	r.HandleFunc("/logout", api.Logout).Methods("GET")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Printf("Error by starting server: %v", err)
	}
}
