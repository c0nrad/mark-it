package main

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	Database              = "markit"
	UsersCollection       = "users"
	EventsCollection      = "events"
	CommentsCollection    = "comments"
	ChatsCollection       = "chats"
	TeamsCollection       = "team"
	AttachmentsCollection = "attachments"
)

var Session *mgo.Session

type User struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	TS time.Time

	Name     string
	Email    string
	Password string

	ProfilePic  string
	Description string
}

type Team struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	TS time.Time

	Name        string
	Description string
	Body        string
}

type Chat struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	Team bson.ObjectId
	TS   time.Time

	Author     string
	ProfilePic string
	Message    string
}

type Attachment struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	Team bson.ObjectId

	TS time.Time

	Author string
	Title  string
	Link   string
}

type Event struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	Team bson.ObjectId

	TS time.Time

	Author string
	Type   string
	Data   string
}

type Comment struct {
	ID     bson.ObjectId `bson:"_id,omitempty"`
	Parent int
	User   int
	TS     time.Time

	Author     string
	ProfilePic string

	Url  string
	Path string
	Body string
}

func MongoURI() string {
	return "mongodb://localhost"
}

func ConnectToMongo() {
	session, err := mgo.Dial(MongoURI())
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)

	index := mgo.Index{
		Key:        []string{"provider", "username"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err = session.DB(Database).C(UsersCollection).EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	Session = session
}

func GetUserByEmail(email string) (*User, error) {
	session := Session.Copy()
	defer session.Close()

	var user User
	err := session.DB(Database).C(UsersCollection).Find(bson.M{"email": email}).One(&user)
	return &user, err
}

func GetUserByLogin(email, password string) (*User, error) {
	session := Session.Copy()
	defer session.Close()

	var user User
	err := session.DB(Database).C(UsersCollection).Find(bson.M{"email": email, "password": password}).One(&user)
	return &user, err
}

func InsertEvent(event Event) (*Event, error) {
	session := Session.Copy()
	defer session.Close()

	event.ID = bson.NewObjectId()
	err := session.DB(Database).C(EventsCollection).Insert(event)
	return &event, err
}

func InsertComment(comment Comment) (*Comment, error) {
	session := Session.Copy()
	defer session.Close()

	comment.ID = bson.NewObjectId()
	err := session.DB(Database).C(CommentsCollection).Insert(comment)
	return &comment, err
}

func InsertChat(chat Chat) (*Chat, error) {
	session := Session.Copy()
	defer session.Close()

	chat.ID = bson.NewObjectId()
	err := session.DB(Database).C(ChatsCollection).Insert(chat)
	return &chat, err
}

func InsertUser(user User) (*User, error) {
	session := Session.Copy()
	defer session.Close()

	user.ID = bson.NewObjectId()
	err := session.DB(Database).C(ChatsCollection).Insert(user)
	return &user, err
}

func GetTeamById(teamId string) (*Team, error) {
	session := Session.Copy()
	defer session.Close()

	if !bson.IsObjectIdHex(teamId) {
		return nil, errors.New("not valid teamId" + teamId)
	}

	var team Team
	err := session.DB(Database).C(TeamsCollection).FindId(bson.ObjectIdHex(teamId)).One(&team)
	return &team, err
}

func UpdateTeam(team Team) error {
	session := Session.Copy()
	defer session.Close()

	err := session.DB(Database).C(TeamsCollection).UpdateId(team.ID, team)
	return err
}

func GetComments() ([]Comment, error) {
	session := Session.Copy()
	defer session.Close()

	var comments []Comment
	err := session.DB(Database).C(CommentsCollection).Find(bson.M{}).All(&comments)
	return comments, err
}

func GetAttachmentsByTeam(teamId string) ([]Attachment, error) {
	session := Session.Copy()
	defer session.Close()

	var attachments []Attachment
	err := session.DB(Database).C(AttachmentsCollection).Find(bson.M{"team": bson.ObjectIdHex(teamId)}).All(&attachments)
	return attachments, err
}

func GetTeams() ([]Team, error) {
	session := Session.Copy()
	defer session.Close()

	var teams []Team
	err := session.DB(Database).C(TeamsCollection).Find(bson.M{}).All(&teams)
	return teams, err
}

func GetChatsByTeam(teamId string) ([]Chat, error) {
	session := Session.Copy()
	defer session.Close()

	var chats []Chat
	err := session.DB(Database).C(ChatsCollection).Find(bson.M{"team": bson.ObjectIdHex(teamId)}).All(&chats)
	return chats, err
}

func GetEventsByTeam(teamId string) ([]Event, error) {
	session := Session.Copy()
	defer session.Close()

	var events []Event
	err := session.DB(Database).C(EventsCollection).Find(bson.M{"team": bson.ObjectIdHex(teamId)}).All(&events)
	return events, err
}

func GetUsersByTeam(teamId string) ([]User, error) {
	session := Session.Copy()
	defer session.Close()

	var users []User
	err := session.DB(Database).C(UsersCollection).Find(bson.M{}).All(&users)
	return users, err
}

func AddFakeData() {
	Session.DB(Database).DropDatabase()

	var Users = []User{
		User{bson.NewObjectId(), time.Now(), "Stuart Larsen", "c0nrad@c0nrad.io", "a", "https://avatars3.githubusercontent.com/u/1901151?v=3&s=460", "Stuart is the lead programmer at Mark.It. He manages the technical day to day activites"},
		User{bson.NewObjectId(), time.Now(), "Katie Honadle", "katie.honadle@gmail.com", "a", "https://static.wixstatic.com/media/a905fd_76d88cea66d64e56a47221b59b381f21.jpg/v1/fill/w_870,h_836,al_c,q_85,usm_0.66_1.00_0.01/a905fd_76d88cea66d64e56a47221b59b381f21.jpg", "Katie has a wealth of experience in Marketing, Social Media, and Human Resources. She's currently VP of Marketing, Product, and Design at Mark.It."},
	}
	for _, user := range Users {
		Session.DB(Database).C(UsersCollection).Insert(user)
	}

	var Teams = []Team{
		Team{bson.NewObjectId(), time.Now(), "Party Planning Committee", "partypartyparty", "Swag"},
	}
	for _, team := range Teams {
		Session.DB(Database).C(TeamsCollection).Insert(team)
	}

	var Chats = []Chat{
		Chat{bson.NewObjectId(), Teams[0].ID, time.Now(), "Stuart", "https://avatars3.githubusercontent.com/u/1901151?v=3&s=460", "Did you come up with a name for the project?"},
		Chat{bson.NewObjectId(), Teams[0].ID, time.Now(), "Katie", "https://static.wixstatic.com/media/a905fd_76d88cea66d64e56a47221b59b381f21.jpg/v1/fill/w_870,h_836,al_c,q_85,usm_0.66_1.00_0.01/a905fd_76d88cea66d64e56a47221b59b381f21.jpg", "Yeah, I'm thinking MarkIt!"},
		Chat{bson.NewObjectId(), Teams[0].ID, time.Now(), "Stuart", "https://avatars3.githubusercontent.com/u/1901151?v=3&s=460", "That's an awesome name!"},
	}
	for _, chat := range Chats {
		Session.DB(Database).C(ChatsCollection).Insert(chat)
	}

	var Attachments = []Attachment{
		Attachment{bson.NewObjectId(), Teams[0].ID, time.Now(), "Stuart", "How to Use Slack.pdf", "https://google.com"},
		Attachment{bson.NewObjectId(), Teams[0].ID, time.Now(), "Katie", "Design Idea", "https://yahoo.com"},
	}
	for _, attachment := range Attachments {
		Session.DB(Database).C(AttachmentsCollection).Insert(attachment)
	}

	var Events = []Event{
		Event{bson.NewObjectId(), Teams[0].ID, time.Now(), "Stuart", "comment", "https://amazon.com"},
		Event{bson.NewObjectId(), Teams[0].ID, time.Now(), "Katie", "attachment", "Party Planning Committee"},
	}

	for _, event := range Events {
		Session.DB(Database).C(EventsCollection).Insert(event)
	}
}
