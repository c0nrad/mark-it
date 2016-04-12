package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var SessionStore = sessions.NewCookieStore([]byte("swagswagswagswag"))
var SessionName = "markit-session"

func main() {
	ConnectToMongo()
	// AddFakeData()

	router := mux.NewRouter()
	router.HandleFunc("/api/me", GetMeHandler)
	router.HandleFunc("/api/users/login", LoginHandler).Methods("POST")
	router.HandleFunc("/api/users/register", RegisterHandler).Methods("POST")

	router.HandleFunc("/api/team", GetTeamsHandler).Methods("GET")
	router.HandleFunc("/api/team/{teamId}", GetTeamHandler).Methods("GET")

	router.HandleFunc("/api/team/{teamId}", UpdateTeamHandler).Methods("PUT")

	router.HandleFunc("/api/team/{teamId}/members", GetMembersHandler).Methods("GET")

	router.HandleFunc("/api/team/{teamId}/chats", GetChatsHandler).Methods("GET")
	router.HandleFunc("/api/team/{teamId}/chats", NewChatsHandler).Methods("POST")

	router.HandleFunc("/api/team/comments", GetCommentsHandler).Methods("GET")
	router.HandleFunc("/api/team/comments", NewCommentHandler).Methods("POST")

	router.HandleFunc("/api/team/{teamId}/attachments", GetAttachmentsHandler)

	router.HandleFunc("/api/team/{teamId}/events", GetEventsHandler).Methods("GET")
	router.HandleFunc("/api/team/{teamId}/events", NewEventHandler).Methods("POST")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

	fmt.Println("[+] Listening on 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func GetCurrentUser(r *http.Request) (*User, error) {
	session, err := SessionStore.Get(r, SessionName)
	if err != nil {
		return nil, errors.New("invalid session store")
	}

	// Set some session values.
	email, okay := session.Values["email"].(string)
	if email == "" || !okay {
		return nil, errors.New("user not logged in")
	}

	user, err := GetUserByEmail(email)

	return user, err
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var newUser User
	err := decoder.Decode(&newUser)
	if err != nil {
		http.Error(w, "unable to decode user", 400)
		return
	}

	fmt.Printf("[+] LoginHandler %+v\n", newUser)

	user, err := InsertUser(newUser)
	if err != nil {
		http.Error(w, "unable to login user", 400)
		return
	}

	session, err := SessionStore.Get(r, SessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set some session values.
	session.Values["email"] = user.Email
	session.Save(r, w)

	json.NewEncoder(w).Encode(user)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var userAttempt User
	err := decoder.Decode(&userAttempt)
	if err != nil {
		http.Error(w, "unable to decode user", 400)
		return
	}

	fmt.Printf("[+] LoginHandler %+v\n", userAttempt)

	user, err := GetUserByLogin(userAttempt.Email, userAttempt.Password)
	if err != nil {
		http.Error(w, "unable to login user", 400)
		return
	}

	session, err := SessionStore.Get(r, SessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set some session values.
	session.Values["email"] = user.Email
	session.Save(r, w)

	json.NewEncoder(w).Encode(user)

}

func NewCommentHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var newComment Comment
	err := decoder.Decode(&newComment)
	if err != nil {
		http.Error(w, "unable to decode comment", 400)
		return
	}

	currentUser, err := GetCurrentUser(r)
	if err != nil {
		http.Error(w, "not logged in", 400)
		return
	}

	newComment.TS = time.Now()
	newComment.ProfilePic = currentUser.ProfilePic
	newComment.Author = currentUser.Name

	fmt.Printf("[+] NewCommentHandler: %+v\n", newComment)

	comment, err := InsertComment(newComment)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	newEvent := Event{TS: time.Now(), Author: comment.Author, Type: "comment", Data: comment.Body}
	_, err = InsertEvent(newEvent)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(comment)
}

func NewEventHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var newEvent Event
	err := decoder.Decode(&newEvent)
	if err != nil {
		http.Error(w, "unable to decode event", 400)
		return
	}

	currentUser, err := GetCurrentUser(r)
	if err != nil {
		http.Error(w, "not logged in", 400)
		return
	}

	newEvent.TS = time.Now()
	newEvent.Author = currentUser.Name
	// event.User = currentUser.ID

	event, err := InsertEvent(newEvent)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(event)
}

func GetMeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[+] GetMeHandler")

	user, err := GetCurrentUser(r)
	if err != nil {
		http.Error(w, "not logged in", 400)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func GetTeamHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	teamId := vars["teamId"]

	team, err := GetTeamById(teamId)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(team)
}

func UpdateTeamHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var team Team
	err := decoder.Decode(&team)
	if err != nil {
		http.Error(w, "unable to decode chat", 400)
		return
	}

	err = UpdateTeam(team)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(team)

}

func NewChatsHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var newChat Chat
	err := decoder.Decode(&newChat)
	if err != nil {
		http.Error(w, "unable to decode chat", 400)
		return
	}

	user, err := GetCurrentUser(r)
	if err != nil {
		http.Error(w, "not logged in", 400)
		return
	}

	vars := mux.Vars(r)
	teamId := vars["teamId"]

	newChat.Author = user.Name
	newChat.Team = bson.ObjectIdHex(teamId)
	newChat.TS = time.Now()
	newChat.ProfilePic = user.ProfilePic

	chat, err := InsertChat(newChat)
	if err != nil {
		fmt.Println(chat, newChat)
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(chat)
}

func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {

	comments, err := GetComments()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(comments)
}

func GetMembersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamId := vars["teamId"]

	users, err := GetUsersByTeam(teamId)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func GetChatsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamId := vars["teamId"]

	chats, err := GetChatsByTeam(teamId)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(chats)
}

func GetAttachmentsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamId := vars["teamId"]

	attachments, err := GetAttachmentsByTeam(teamId)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(attachments)
}

func GetEventsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamId := vars["teamId"]

	events, err := GetEventsByTeam(teamId)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(events)
}

func GetTeamsHandler(w http.ResponseWriter, r *http.Request) {

	teams, err := GetTeams()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(teams)
}
