package main

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"strconv"
	"strings"
	"time"
)

//VERSION is and exported variable so the handelers can use it.
var VERSION string

//CODE ... I don't remember. Do we even need this?
var CODE string

// CODENAME is like a major version string
var CODENAME string

//altpath is used if you need an alternate path for some web servers.
var altpath string

//dbuser pass and name and secret to be pulled from config
var dbuser string
var dbpass string
var dbname string
var secret string

// var tokenmap map[string]string
type tokenz map[string]string

var tokenmap tokenz

func init() {
	//init tokenmap for later use.
	s := make(tokenz, len("tokenstring"))
	s["email"] = "tokenstring"
	tokenmap = s

}

// keys and claims for use in jwt tokens/cookies
var mySigningKey []byte
var token *jwt.Token
var claims jwt.MapClaims

//UserClaims for jwt
type UserClaims struct {
	Admin bool   `json:"admin"`
	Name  string `json:"name"`
	EMail string `json:"email"`
	jwt.StandardClaims
}

func main() {
	//if no args are given then load defaults.
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) > 0 {
		if argsWithoutProg[0] == "install" {
			if _, err := os.Stat("/etc/nmwrapreports/"); os.IsNotExist(err) {
				pathErr := os.MkdirAll("/etc/nmwrapreports/", 0777)
				if pathErr != nil {
					fmt.Println(pathErr)
				}
				d1 := []byte("#Log files location\nLogDir = \"/var/log/\"\n\n#the server socket info\nIP = \"localhost\"\nPort = \"888\"")
				err := ioutil.WriteFile("/etc/nmwrapreports/nmwrapreports.conf", d1, 0644)
				if err != nil {
					fmt.Println(err)
				}
				os.OpenFile("/var/log/nmwrapreports.log", os.O_RDONLY|os.O_CREATE, 0666)
			}
		} else {
			fmt.Println("Unknown param")
		}
	} else {
		VERSION = "1.5"
		CODENAME = "water"
		var configf = ReadConfig() //this is in config.go

		altpath = configf.AltPath
		dbuser = configf.DBUser
		dbpass = configf.DBPass
		dbname = configf.DBName
		secret = configf.Secret
		mySigningKey = []byte(secret)
		token = jwt.New(jwt.SigningMethodHS256)
		claims = token.Claims.(jwt.MapClaims)

		//cleanup job ran every hour.
		ticker := time.NewTicker(time.Hour * 1)
		go func() {
			for t := range ticker.C {
				cleanup(configf.DownloadWindow)
				_ = t
			}
		}()

		//extract job mailer run every minute.
		tickermin := time.NewTicker(time.Minute * 1)
		go func() {
			for t := range tickermin.C {
				ExtractJobs()
				_ = t
			}
		}()

		listensocket := configf.IP + ":" + configf.Port
		router := NewRouter()
		r := handlers.LoggingHandler(os.Stdout, router)
		log.Println("server running on " + listensocket)

		log.Fatal(http.ListenAndServe(listensocket, r))

	}
}

//Append is a function for appending slices
func Append(slice []string, items ...string) []string {
	for _, item := range items {
		slice = Extend(slice, item)
	}
	return slice
}

//Extend is an easy wat to grow a slice.
func Extend(slice []string, element string) []string {
	n := len(slice)
	if n == cap(slice) {
		// Slice is full; must grow.
		// We double its size and add 1, so if the size is zero we still grow.
		newSlice := make([]string, len(slice), 2*len(slice)+1)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0 : n+1]
	slice[n] = element
	return slice
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

//cleanup function to delete files in tmp.

func cleanup(window string) {

	var cut, _ = strconv.Atoi(window)
	var cutoff = time.Duration(cut) * time.Hour

	log.Println("cleanup running")
	fileInfo, err := ioutil.ReadDir("/tmp")
	if err != nil {
		log.Fatal(err.Error())
	}
	now := time.Now()
	for _, info := range fileInfo {
		if strings.HasSuffix(info.Name(), "pdf") {
			if diff := now.Sub(info.ModTime()); diff > cutoff {
				fmt.Printf("Deleting %s which is %s old\n", info.Name(), diff)
				var err = os.Remove("/tmp/" + info.Name())
				if err != nil {
					log.Println("Error trying to remove " + info.Name())
					log.Println(err)
				}
			}
		}
	}
}

func logErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
