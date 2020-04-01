package main

import (
	"html/template"
	"net/http"
	"fmt"	
	"strconv"
	"errors"
	"log"
)

var (
	tonnels map[int]*GameAdvanced = make(map[int]*GameAdvanced)
)

var (
	errSpaceInWord = errors.New("space in word")
	NotGameError   = errors.New("Key is not found")
	errLose        = errors.New("you are lose")
	errWin         = errors.New("you are win")
)

// Game ...
type Game struct {
	hp       int
	word     string
	wordUser string
	beside   map[string]bool
}

// New ...
func New(word string) *Game {
	var beside = make(map[string]bool)
	var wordUser string = ""
	for i := 0; i < len(word); i++ {
		wordUser += "_"
	}
	return &Game{word: word, hp: len(word), wordUser: wordUser, beside: beside}
}

func (s *Game) validate() error {
	for i := 0; i < len(s.word); i++ {
		if s.word[i] == ' ' {
			return errSpaceInWord
		}
	}
	return nil
}
func (s *Game) send(char string) (string, int) {
	// base method
	fmt.Println(s.hp)
	var saveWordUser string = s.wordUser
	if _, ok := s.beside[char]; ok {
		return s.wordUser, s.hp
	}
	var wordUserLocal string
	for i := 0; i < len(s.word); i++ {
		if (s.wordUser[i]=='_') {
			if s.word[i] == char[0] {
				wordUserLocal += string(char[0])
			} else{
				wordUserLocal += "_"
			}
		} else {
			wordUserLocal += string(s.word[i])
		}
		
	}
	if wordUserLocal==saveWordUser{
		s.hp--
	}

	s.wordUser = wordUserLocal
	s.beside[char] = true
	// init wordUserLocal
	return wordUserLocal, s.hp
}

// validate ...
func (s *Game) getClassStyle () string{
	if len(s.wordUser)-s.hp==11{
		return "los"
	} else if (s.word==s.wordUser){
		return "win"
	} else {
		return "undefined"
	}
}

// GameAdvanced ...
type GameAdvanced struct {
	privateKey int
	publicKey int
	game *Game
}

// NewGameAdvanced ...
func  NewGameAdvanced (privateKey int, publicKey int, game *Game) *GameAdvanced{
	return &GameAdvanced {
		privateKey: privateKey,
		publicKey: publicKey,
		game: game,
	}
}
// getPublicGame ...
func getPublicGame (publicKey int) (*Game, error){
	if val, ok := tonnels[publicKey]; ok {
		return val.game, nil
	} else {
		return nil, NotGameError
	}
}

// getPrivateGame ...
func  getPrivateGame (privateKey int) (*Game, error){
	for _, val := range tonnels {
		if val.privateKey==privateKey{
			return val.game, nil
		}
	}
	return nil, NotGameError
}

func (s *Game) handleSend(w http.ResponseWriter, r *http.Request)  {
	val := r.FormValue("char")
	log.Println("User enter ", val)
	rs, hp := s.send(val)
	hpStr := strconv.Itoa(hp)
	answer := rs+" "+hpStr
	fmt.Fprint(w, answer)
}

// Run ...
func Run(word string, ip string) {
	game := New(word)

	log.Println("Server is listening...")

	http.HandleFunc("/send", game.handleSend)
	log.Fatal(http.ListenAndServe(ip, nil))
}

// error
func errorUser(w http.ResponseWriter, r *http.Request){
	t, err := template.ParseFiles("templates/userError.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	
	t.ExecuteTemplate(w, "userError", nil)
}


// indexHandle
func indexHandle(w http.ResponseWriter, r *http.Request){
	r.ParseForm()

	for name, val := range r.Form{
		fmt.Println("|", "\t", name, " ", val)
	}
	fmt.Println("indexHandle")
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "index", nil)
}

// createIndex
func createIndex(w http.ResponseWriter, r *http.Request){
	r.ParseForm()

	for name, val := range r.Form{
		fmt.Println("|", "\t", name, " ", val)
	}
	fmt.Println("createIndex")
	// create ...
	t, err := template.ParseFiles("templates/create.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "create", nil)

}

// // createPost
// func createAfter(w http.ResponseWriter, r *http.Request){
// }


// createBack
func createBack(w http.ResponseWriter, r *http.Request){
	r.ParseForm()

	for name, val := range r.Form{
		fmt.Println("|", "\t",name, " ", val)
	}
	fmt.Println("createBack")
	publicKey := r.FormValue("publicKey")	

	privateKey := r.FormValue("privateKey")

	publicKeyInt, err  := strconv.Atoi(publicKey)
	if err != nil {
		createIndex(w, r)
	}
	privateKeyInt, err := strconv.Atoi(privateKey)


	if err != nil {
		createIndex(w, r)
	} else {
		word := r.FormValue("word")
		game := New(word)
		gameAdvanced := NewGameAdvanced(privateKeyInt, publicKeyInt, game)


		tonnels[publicKeyInt]=gameAdvanced

		r.URL.Path = "/get/"
		getInfoWordIndex(w, r)

	}

}

// join
func join(w http.ResponseWriter, r *http.Request){
	fmt.Println("join")
	r.ParseForm()

	for name, val := range r.Form{
		fmt.Println("|", "\t",name, " ", val)
	}
	t, err := template.ParseFiles("templates/join.html")
	if err != nil {
		log.Fatal(err) 
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "join", nil)
}

// joinBack
func joinBack(w http.ResponseWriter, r *http.Request){
	r.ParseForm()

	for name, val := range r.Form{
		fmt.Println("|", "\t",name, " ", val)
	}
	fmt.Println("joinBack")
	publicKey := r.FormValue("publicKey")
	publicKeyInt, err := strconv.Atoi(publicKey)
	if err != nil {
		join(w, r)
	} else {
		if _, ok := tonnels[publicKeyInt]; ok{
			getTonnel(w, r)
			return
		}
		join(w, r)	
	}
}


type ctx struct {
	Word string
	Key int
	Class string
	Hp int
}

// NewCtx ...
func NewCtx (Key int, Hp int, Class string, Word string) *ctx{
	return &ctx {
		Key: Key,
		Hp: Hp,
		Class: Class,
		Word: Word,	
	}
}

// getTonnel
func getTonnel(w http.ResponseWriter, r *http.Request){
	fmt.Println("getTonnel")

	r.ParseForm()

	for name, val := range r.Form{
		fmt.Println("|", "\t",name, " ", val)
	}
	publicKey := r.FormValue("publicKey")
	publicKeyInt, _ := strconv.Atoi(publicKey)
	game, err  := getPublicGame(publicKeyInt)
	if err != nil {
		join(w, r)
	} else {
		t, err := template.ParseFiles("templates/tonnel.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}

		context := NewCtx ( publicKeyInt, 0, game.getClassStyle(), game.wordUser,  ) 
		t.ExecuteTemplate(w, "tonnel", context)		
	}


}

// send
func send(w http.ResponseWriter, r *http.Request){
	fmt.Println("send")

	r.ParseForm()

	for name, val := range r.Form{
		fmt.Println("|", "\t",name, " ", val)
	}

	publicKey := r.FormValue("publicKey")
	publicKeyInt, err := strconv.Atoi(publicKey)
	if err != nil {
		getTonnel(w, r)
	}
	game, err  := getPublicGame(publicKeyInt)
	if err != nil {
		getTonnel(w, r)
	}
	char := r.FormValue("char")
	if(char==""){
		getTonnel(w, r)
		return
	}
	game.send(char)

	t, err := template.ParseFiles("templates/tonnel.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	fmt.Println(game.hp)

	context := NewCtx ( publicKeyInt,  len(game.wordUser)-game.hp, game.getClassStyle(), game.wordUser,  ) 
	fmt.Println(context)
	t.ExecuteTemplate(w, "tonnel", context)
}

// getInfoWord
func getInfoWord(w http.ResponseWriter, r *http.Request){
	privateKey := r.FormValue("privateKey")
	privateKeyInt, err := strconv.Atoi(privateKey)
	if err != nil {
		fmt.Fprint(w, err.Error())
	}

	game, err := getPrivateGame(privateKeyInt)
	if err != nil {
		errorUser(w, r)
	}

	fmt.Fprint(w, game.wordUser)
}

// getInfoWord
func getInfoHp(w http.ResponseWriter, r *http.Request){
	privateKey := r.FormValue("privateKey")
	privateKeyInt, err := strconv.Atoi(privateKey)
	if err != nil {
		fmt.Fprint(w, err.Error())
	}

	game, err := getPrivateGame(privateKeyInt)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Fprint(w, len(game.word)-game.hp)
}

// getInfoClass
func getInfoClass(w http.ResponseWriter, r *http.Request){
	privateKey := r.FormValue("privateKey")
	privateKeyInt, err := strconv.Atoi(privateKey)
	if err != nil {
		fmt.Fprint(w, err.Error())
	}

	game, err := getPrivateGame(privateKeyInt)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Fprint(w, game.getClassStyle())
}

// CtxForGetInfoWordIndex ...
type CtxForGetInfoWordIndex struct {
	Class 		string
	PrivateKey 	string
	Hp 			int
	Word 		string
}

// NewCtxForGetInfoWordIndex ...
func  NewCtxForGetInfoWordIndex (PrivateKey string) *CtxForGetInfoWordIndex{
	privateKeyInt, _:= strconv.Atoi(PrivateKey)

	game, _ := getPrivateGame(privateKeyInt)
	return &CtxForGetInfoWordIndex{
		PrivateKey: PrivateKey,
		Class: game.getClassStyle(),
		Word: game.wordUser,
		Hp: len(game.word)-game.hp,
	}
}

// getInfoWordIndex
func getInfoWordIndex(w http.ResponseWriter, r *http.Request){
	t, err := template.ParseFiles("templates/wordAdmin.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	privateKey := r.FormValue("privateKey")

	_, err = strconv.Atoi(privateKey)
	if err != nil {
		errorUser(w, r)
	}
	t.ExecuteTemplate(w, "wordAdmin", NewCtxForGetInfoWordIndex(privateKey))
}



// main ...
func main () {
	fmt.Println("Server is listining...")

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	http.HandleFunc("/", indexHandle)
	http.HandleFunc("/create/", createIndex)
	http.HandleFunc("/createBack/", createBack)
	http.HandleFunc("/join/", join)
	http.HandleFunc("/joinBack/", joinBack)
	http.HandleFunc("/getTonnel/", getTonnel)
	http.HandleFunc("/send/",send)
	http.HandleFunc("/getInfoWord/", getInfoWord)
	http.HandleFunc("/getInfoWordIndex/", getInfoWordIndex)
	http.HandleFunc("/getInfoClass/", getInfoClass)
	http.HandleFunc("/getInfoHp/", getInfoHp)

	log.Fatal(http.ListenAndServe(":3000", nil))
}