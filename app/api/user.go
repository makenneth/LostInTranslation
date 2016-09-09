package api 
import (
  "log"
  "net/http"
  "time"
  "../database"
  "golang.org/x/crypto/bcrypt"
    "../token"
  )
type UserData struct {
  Username string `json:"username"`
  Password string `json:"password"`
}

type User struct {
  Id int `json:"id"`
  Username string `json:"username"`
}
var currentUser User;

func SetCurrentUser(u User){
  currentUser = u
}

func GetCurrentUser(w http.ResponseWriter, userToken string) User {
  if currentUser != (User{}){
    return currentUser
  }
  bool := FindCurrentUser(w, userToken)
  if !bool {
    return (User{})
  }
  return currentUser
}
func FindCurrentUser(w http.ResponseWriter, userToken string) bool {
  var (
    id int
    username string
    )
  newToken, _ := token.GenerateRandomToken(32)
  err := database.DBConn.QueryRow(`UPDATE users
    SET session_token = $1
    WHERE session_token = $2 
    returning id, username`, newToken, userToken).Scan(&id, &username)
  log.Println(err)
  if id == 0 || err != nil {
    return false
  }
  expiration := time.Now().Add(30 * 24 * time.Hour)
  cookie := &http.Cookie{Name: "sessiontokenLit", Value: newToken, Expires: expiration, Path: "/"}
  http.SetCookie(w, cookie)
  if err != nil {
    return false
  }
  SetCurrentUser(User{id, username});
  return true
}

func (u *UserData) checkPassword() error {
  var (
    passwordDigest string
    sessionToken string
  )
  err := database.DBConn.QueryRow(`SELECT password_digest, session_token 
    FROM users 
    WHERE username = $1`, u.Username).Scan(&passwordDigest, &sessionToken)

  if err != nil {
    return err
  }

  password, digest := []byte(u.Password), []byte(passwordDigest)
  err = bcrypt.CompareHashAndPassword(digest, password)

  return err
}

func (u *UserData) resetSessionToken() (int, string, error) {
  newToken, _ := token.GenerateRandomToken(32)
  var playerId int
  err := database.DBConn.QueryRow(`UPDATE users 
    SET session_token = $1
    WHERE username = $2 returning id`, newToken, u.Username).Scan(&playerId)

  return playerId, newToken, err
}

func (u *UserData) InsertUser() (string, int, error) { 

  password := []byte(u.Password)
  digest, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
  if err != nil {
    panic(err)
  }

  token, err := token.GenerateRandomToken(32)
  if err != nil {
    panic(err)
  }

  var newPlayerId int
  err = database.DBConn.QueryRow(`INSERT INTO 
    users (session_token, username, password_digest) 
    VALUES ($1, $2, $3) returning id`, token, u.Username, digest).Scan(&newPlayerId)

  return token, newPlayerId, err
}  

