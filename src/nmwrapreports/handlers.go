package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"os"
	"reflect"
	//"reflect"
	//"crypto/sha512"
	"archive/zip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jonas-p/go-shp"
	"github.com/jung-kurt/gofpdf"
	gdal "github.com/lukeroth/gdal_go"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	//"reflect"
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	_ "github.com/go-sql-driver/mysql"
	mathrand "math/rand"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

type Userpack struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserpack struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type PassPack struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

//FireStations type is for marshalling the output from arcgis server
type FireStations struct {
	DisplayFieldName string `json:"displayFieldName"`
	Features         []struct {
		Attributes struct {
			ADDRESS  string `json:"ADDRESS"`
			CITY     string `json:"CITY"`
			INSTNAME string `json:"INST_NAME"`
			OBJECTID int    `json:"OBJECTID"`
		} `json:"attributes"`
	} `json:"features"`
	FieldAliases struct {
		ADDRESS  string `json:"ADDRESS"`
		CITY     string `json:"CITY"`
		INSTNAME string `json:"INST_NAME"`
		OBJECTID string `json:"OBJECTID"`
	} `json:"fieldAliases"`
	Fields []struct {
		Alias string `json:"alias"`
		Name  string `json:"name"`
		Type  string `json:"type"`
	} `json:"fields"`
}

//CommunitesatRisk type is for marshalling the output from arcgis server
type CommunitesatRisk struct {
	DisplayFieldName string `json:"displayFieldName"`
	Features         []struct {
		Attributes struct {
			NAME       string `json:"NAME"`
			NAME1      string `json:"NAME_1"`
			OBJECTID12 int    `json:"OBJECTID_12"`
			Rate2016   string `json:"Rate_2016"`
		} `json:"attributes"`
	} `json:"features"`
	FieldAliases struct {
		NAME       string `json:"NAME"`
		NAME1      string `json:"NAME_1"`
		OBJECTID12 string `json:"OBJECTID_12"`
		Rate2016   string `json:"Rate_2016"`
	} `json:"fieldAliases"`
	Fields []struct {
		Alias string `json:"alias"`
		Name  string `json:"name"`
		Type  string `json:"type"`
	} `json:"fields"`
}

//IncorporatedCityBoundaries type is for marshalling the output from arcgis server
type IncorporatedCityBoundaries struct {
	DisplayFieldName string `json:"displayFieldName"`
	Features         []struct {
		Attributes struct {
			GEOID10     string  `json:"GEOID10"`
			NAME10      string  `json:"NAME10"`
			NAMELSAD10  string  `json:"NAMELSAD10"`
			OBJECTID    int     `json:"OBJECTID"`
			ShapeArea   float64 `json:"Shape_Area"`
			ShapeLength float64 `json:"Shape_Length"`
		} `json:"attributes"`
	} `json:"features"`
	FieldAliases struct {
		GEOID10     string `json:"GEOID10"`
		NAME10      string `json:"NAME10"`
		NAMELSAD10  string `json:"NAMELSAD10"`
		OBJECTID    string `json:"OBJECTID"`
		ShapeArea   string `json:"Shape_Area"`
		ShapeLength string `json:"Shape_Length"`
	} `json:"fieldAliases"`
	Fields []struct {
		Alias string `json:"alias"`
		Name  string `json:"name"`
		Type  string `json:"type"`
	} `json:"fields"`
}

//VegetationTreatments type is for marshalling the output from arcgis server
type VegetationTreatments struct {
	DisplayFieldName      string `json:"displayFieldName"`
	ExceededTransferLimit bool   `json:"exceededTransferLimit"`
	Features              []struct {
		Attributes struct {
			AcreUS      float64 `json:"Acre_US"`
			Agency      string  `json:"Agency"`
			Description string  `json:"Description"`
			LandOwner   string  `json:"Land_Owner"`
			NameProj    string  `json:"Name_Proj"`
			OBJECTID    int     `json:"OBJECTID"`
			Partners    string  `json:"Partners"`
			ProjectType string  `json:"Project_Type"`
			ShapeArea   float64 `json:"Shape_Area"`
			ShapeLength float64 `json:"Shape_Length"`
			TargetSpec  string  `json:"Target_Spec"`
			TypeProj    string  `json:"Type_Proj"`
			YearCal     string  `json:"Year_Cal"`
		} `json:"attributes"`
	} `json:"features"`
	FieldAliases struct {
		AcreUS      string `json:"Acre_US"`
		Agency      string `json:"Agency"`
		Description string `json:"Description"`
		LandOwner   string `json:"Land_Owner"`
		NameProj    string `json:"Name_Proj"`
		OBJECTID    string `json:"OBJECTID"`
		Partners    string `json:"Partners"`
		ProjectType string `json:"Project_Type"`
		ShapeArea   string `json:"Shape_Area"`
		ShapeLength string `json:"Shape_Length"`
		TargetSpec  string `json:"Target_Spec"`
		TypeProj    string `json:"Type_Proj"`
		YearCal     string `json:"Year_Cal"`
	} `json:"fieldAliases"`
	Fields []struct {
		Alias string `json:"alias"`
		Name  string `json:"name"`
		Type  string `json:"type"`
	} `json:"fields"`
}

// WatershedsHUC8 type is for marshalling the output from arcgis server
type WatershedsHUC8 struct {
	DisplayFieldName string `json:"displayFieldName"`
	Features         []struct {
		Attributes struct {
			AREAACRES                int         `json:"AREAACRES"`
			AREASQKM                 float64     `json:"AREASQKM"`
			AcresBurned2006_2016     string      `json:"Acres_Burned_2006_2016"`
			AcresBurned2006_2016X    interface{} `json:"Acres_Burned_2006_2016_X"`
			AcresBurned2006_2016Y    interface{} `json:"Acres_Burned_2006_2016_Y"`
			EquationRanking          int         `json:"Equation_Ranking"`
			EquationResult           int         `json:"Equation_Result"`
			F1IntermixAcres          string      `json:"F1_Intermix_Acres"`
			F1IntermixWUIStructures  int         `json:"F1_Intermix__WUI__Structures"`
			F2InterfaceAcres         string      `json:"F2_Interface_Acres"`
			F2InterfaceAcresX        interface{} `json:"F2_Interface_Acres_X"`
			F2InterfaceAcresY        interface{} `json:"F2_Interface_Acres_Y"`
			F2InterfaceWUIStructures int         `json:"F2_Interface__WUI__Structures"`
			FinalRankVeryHighHigh    int         `json:"Final_Rank_Very_High___High"`
			FireRank                 int         `json:"Fire_Rank"`
			FloodRankT               string      `json:"FloodRankT"`
			GNISID                   int         `json:"GNIS_ID"`
			HUC8                     string      `json:"HUC8"`
			// HUC8_1                   int         `json:"HUC8_1"`
			// HUC8_12                  int         `json:"HUC8_12"`
			HUC8No                   int         `json:"HUC8_no"`
			HighRisk                 int         `json:"High_Risk"`
			InterfacePercent         int         `json:"Interface_Percent"`
			IntermixPercent          int         `json:"Intermix_Percent"`
			LOADDATE                 int         `json:"LOADDATE"`
			LowRisk                  int         `json:"Low_Risk"`
			MediumRisk               int         `json:"Medium_Risk"`
			NAME                     string      `json:"NAME"`
			NAME1                    string      `json:"NAME_1"`
			NoCommunitiesInWatershed interface{} `json:"No_Communities_in_Watershed"`
			OBJECTID1                int         `json:"OBJECTID_1"`
			OtherAcres               int         `json:"Other_Acres"`
			OtherWUI                 int         `json:"Other__WUI_"`
			PopulationWithinHUC8     int         `json:"Population_Within_HUC_8"`
			PopulationRanking        int         `json:"Population_ranking"`
			RankText                 string      `json:"Rank_Text"`
			RankingStructuresInWUI   int         `json:"Ranking_Structures_in_WUI"`
			SOURCEFEAT               string      `json:"SOURCEFEAT"`
			STATES                   string      `json:"STATES"`
			STATES1                  string      `json:"STATES_1"`
			ShapeArea                float64     `json:"Shape_Area"`
			ShapeLeng                float64     `json:"Shape_Leng"`
			ShapeLength              float64     `json:"Shape_Length"`
			TNCPriority4             int         `json:"TNC_Priority_4"`
			TNCPriority5             int         `json:"TNC_Priority_5"`
			TNCRanking               int         `json:"TNC_Ranking"`
			TotalAcresWUI            string      `json:"Total_Acres_WUI"`
			TotalStructuresHUC       int         `json:"Total_Structures__HUC_"`
			TotalStructuresInWUI     int         `json:"Total_Structures_in_WUI"`
			TreatedAcres             string      `json:"Treated_Acres"`
			WildfireCount2006_2016   int         `json:"Wildfire_Count_2006_2016"`
			Collected                int         `json:"collected"`
			Ranknosme                int         `json:"ranknosme"`
		} `json:"attributes"`
	} `json:"features"`
	FieldAliases struct {
		AREAACRES                string `json:"AREAACRES"`
		AREASQKM                 string `json:"AREASQKM"`
		AcresBurned2006_2016     string `json:"Acres_Burned_2006_2016"`
		AcresBurned2006_2016X    string `json:"Acres_Burned_2006_2016_X"`
		AcresBurned2006_2016Y    string `json:"Acres_Burned_2006_2016_Y"`
		EquationRanking          string `json:"Equation_Ranking"`
		EquationResult           string `json:"Equation_Result"`
		F1IntermixAcres          string `json:"F1_Intermix_Acres"`
		F1IntermixWUIStructures  string `json:"F1_Intermix__WUI__Structures"`
		F2InterfaceAcres         string `json:"F2_Interface_Acres"`
		F2InterfaceAcresX        string `json:"F2_Interface_Acres_X"`
		F2InterfaceAcresY        string `json:"F2_Interface_Acres_Y"`
		F2InterfaceWUIStructures string `json:"F2_Interface__WUI__Structures"`
		FinalRankVeryHighHigh    string `json:"Final_Rank_Very_High___High"`
		FireRank                 string `json:"Fire_Rank"`
		FloodRankT               string `json:"FloodRankT"`
		GNISID                   string `json:"GNIS_ID"`
		HUC8                     string `json:"HUC8"`
		// HUC8_1                   string `json:"HUC8_1"`
		// HUC8_12                  string `json:"HUC8_12"`
		HUC8No                   string `json:"HUC8_no"`
		HighRisk                 string `json:"High_Risk"`
		InterfacePercent         string `json:"Interface_Percent"`
		IntermixPercent          string `json:"Intermix_Percent"`
		LOADDATE                 string `json:"LOADDATE"`
		LowRisk                  string `json:"Low_Risk"`
		MediumRisk               string `json:"Medium_Risk"`
		NAME                     string `json:"NAME"`
		NAME1                    string `json:"NAME_1"`
		NoCommunitiesInWatershed string `json:"No_Communities_in_Watershed"`
		OBJECTID1                string `json:"OBJECTID_1"`
		OtherAcres               string `json:"Other_Acres"`
		OtherWUI                 string `json:"Other__WUI_"`
		PopulationWithinHUC8     string `json:"Population_Within_HUC_8"`
		PopulationRanking        string `json:"Population_ranking"`
		RankText                 string `json:"Rank_Text"`
		RankingStructuresInWUI   string `json:"Ranking_Structures_in_WUI"`
		SOURCEFEAT               string `json:"SOURCEFEAT"`
		STATES                   string `json:"STATES"`
		STATES1                  string `json:"STATES_1"`
		ShapeArea                string `json:"Shape_Area"`
		ShapeLeng                string `json:"Shape_Leng"`
		ShapeLength              string `json:"Shape_Length"`
		TNCPriority4             string `json:"TNC_Priority_4"`
		TNCPriority5             string `json:"TNC_Priority_5"`
		TNCRanking               string `json:"TNC_Ranking"`
		TotalAcresWUI            string `json:"Total_Acres_WUI"`
		TotalStructuresHUC       string `json:"Total_Structures__HUC_"`
		TotalStructuresInWUI     string `json:"Total_Structures_in_WUI"`
		TreatedAcres             string `json:"Treated_Acres"`
		WildfireCount2006_2016   string `json:"Wildfire_Count_2006_2016"`
		Collected                string `json:"collected"`
		Ranknosme                string `json:"ranknosme"`
	} `json:"fieldAliases"`
	Fields []struct {
		Alias string `json:"alias"`
		Name  string `json:"name"`
		Type  string `json:"type"`
	} `json:"fields"`
}

// County type is for marshalling the output from arcgis server
type County struct {
	DisplayFieldName string `json:"displayFieldName"`
	Features         []struct {
		Attributes struct {
			NAME      string `json:"NAME"`
			NAMELSAD  string `json:"NAMELSAD"`
			OBJECTID1 int    `json:"OBJECTID_1"`
		} `json:"attributes"`
	} `json:"features"`
	FieldAliases struct {
		NAME      string `json:"NAME"`
		NAMELSAD  string `json:"NAMELSAD"`
		OBJECTID1 string `json:"OBJECTID_1"`
	} `json:"fieldAliases"`
	Fields []struct {
		Alias string `json:"alias"`
		Name  string `json:"name"`
		Type  string `json:"type"`
	} `json:"fields"`
}

//Geom starting struct for geometry
type Geom struct {
	Rings [][][]float64 `json:"rings"`
	Title string        `json:"title"`
}

type GeoJSON struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}

type justGeom struct {
	Rings [][][]float64 `json:"rings"`
	Title string        `json:"title"`
}

func MakeSalt(n int) string {
	mathrand.Seed(time.Now().UnixNano())

	var saltrune = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = saltrune[mathrand.Intn(len(saltrune))]
	}
	return string(b)
}

func HashPass(pw string, salt string) string {
	hash := sha512.New()
	hash.Write([]byte(pw + salt))
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)
	return mdStr

}

func RandString(n int) string {
	var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[mathrand.Intn(len(letterRunes))]
	}
	return string(b)
}

// func checkCount(rows *sql.Rows) (count int) {
// 	for rows.Next() {
// 		err := rows.Scan(&count)
// 		logErr(err)
// 	}
// 	return count
// }

func GetCookieParts(r *http.Request) UserClaims {
	var cookie, err = r.Cookie("nmwrap")
	logErr(err)
	cookieparts := strings.Split(cookie.Value, ".")
	fmt.Println(cookieparts)
	data, _ := base64.StdEncoding.DecodeString(cookieparts[1] + "==")
	fmt.Println(string(data))
	var user UserClaims
	json.Unmarshal([]byte(string(data)), &user)
	fmt.Println(user.EMail)
	return user
}

func randomFilename() (s string, err error) {
	b := make([]byte, 8)
	_, err = rand.Read(b)
	if err != nil {
		return
	}
	s = fmt.Sprintf("%x", b)
	return
}

//Index is the front page of the app
func Index(w http.ResponseWriter, r *http.Request) {

	type frontpagedict struct {
		MESSAGE string
	}

	params := &frontpagedict{MESSAGE: "Welcome!"}
	t := template.New("frontpage")
	t, err := t.Parse(FrontpageTmpl)
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, params)
	if err != nil {
		log.Fatal(err)
	}
}

//Version shows the version.
func Version(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, VERSION)
}

//GetPOSTGeom is in case anyone wants to know what the route does.
func GetPOSTGeom(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This route requires a POST to do anything ")
}

//Login function
func Login(w http.ResponseWriter, r *http.Request) {
	//tokenmap := make(map[string]string)
	// connectstring := dbuser + ":" + dbpass + "@/" + dbname
	// log.Println(connectstring)
	db, err := sql.Open("mysql", dbuser+":"+dbpass+"@/"+dbname)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println(err)
		return
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println(err)
		return
	}

	// sha_512 := sha512.New()
	jsbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println(err)
		return
	}

	var user Userpack
	json.Unmarshal([]byte(string(jsbody)), &user)
	rows, err := db.Query("SELECT * FROM users WHERE email='" + user.Email + "'")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println(err)
		return
	}
	var id int
	var name string
	var email string
	var userssalt string
	var hash string
	var admin bool

	for rows.Next() {

		err = rows.Scan(&id, &name, &email, &userssalt, &hash, &admin)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)
		}

	}
	log.Println(id)
	log.Println(name)
	log.Println(userssalt)
	log.Println(hash)

	//first we generate the proposed hash...
	submittedhashed := HashPass(user.Password, userssalt)
	if submittedhashed == hash {
		//password is good
		log.Println("pass is good")
		claims["admin"] = admin
		claims["name"] = name
		claims["email"] = email
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
		tokenString, _ := token.SignedString(mySigningKey)
		http.SetCookie(w, &http.Cookie{
			Name:       "nmwrap",
			Value:      tokenString,
			Path:       "/",
			RawExpires: "0",
		})
		tokenmap[email] = tokenString
		fmt.Println(tokenmap)

	} else {

		w.WriteHeader(http.StatusUnauthorized)
	}

}
func CreateUser(w http.ResponseWriter, r *http.Request) {
	if IsLoggedIn(r) {
		// connectstring := dbuser + ":" + dbpass + "@/" + dbname
		// log.Println(connectstring)
		db, err := sql.Open("mysql", dbuser+":"+dbpass+"@/"+dbname)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)
			return
		}
		defer db.Close()
		err = db.Ping()
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)
			return
		}

		jsbody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)

			return
		}

		var user CreateUserpack
		json.Unmarshal([]byte(string(jsbody)), &user)
		log.Println(user.Email)
		log.Println(user.Name)
		randstring, _ := randomFilename()
		salt := MakeSalt(200)
		hashedpass := HashPass(randstring, salt)
		TheQuery := "INSERT INTO users (name,email,salt,hash,admin)  VALUES (\"" + user.Name + "\",\"" + user.Email + "\",\"" + salt + "\",\"" + hashedpass + "\",false);"
		log.Println(TheQuery)
		_, err = db.Exec(TheQuery)
		if err != nil {
			log.Fatal(err)
		} else {
			ResetToken := RandString(200)
			TheQuery := "INSERT INTO passwordtokens (email,token)  VALUES (\"" + user.Email + "\",\"" + ResetToken + "\");"
			log.Println(TheQuery)
			_, err := db.Exec(TheQuery)
			if err != nil {
				log.Fatal(err)
			} else {
				w.WriteHeader(http.StatusOK)
				// SendMail(string(user.Email), ResetToken)
				eMessage := "To:" + user.Email + "\r\n" +
					"Subject: NMWRAP Account creation \r\n" +
					"\r\n" +
					"Your account has been created.\r\n" +
					"Set your password by using the link below.\r\n" +
					"https://nmwrap.org/?token=" + ResetToken + "\r\n"

				SendMail(user.Email, eMessage)
				fmt.Fprintln(w, "mail sent")
			}
			// c, err := smtp.Dial("edacmail.unm.edu:25")
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// defer c.Close()
			// // Set the sender and recipient.
			// c.Mail("nmwrap@edac.unm.edu")
			// c.Rcpt(user.Email)

			// // Send the email body.
			// wc, err := c.Data()
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// defer wc.Close()
			// buf := bytes.NewBufferString("To:" + user.Email + "\r\n" +
			// 	"Subject: NMWRAP Account\r\n" +
			// 	"\r\n" +
			// 	"This is the email body.\r\n")
			// //"Hello, " + user.Name)
			// if _, err = buf.WriteTo(wc); err != nil {
			// 	log.Fatal(err)
			// }
			// log.Println(result)
			// log.Println(err)
			fmt.Fprintln(w, "Account "+user.Email+"created.")
		}
	}
}

//Logout function
func Logout(w http.ResponseWriter, r *http.Request) {
	if IsLoggedIn(r) {
		user := GetCookieParts(r)
		//fmt.Fprintln(w, tokenmap[user.EMail])
		delete(tokenmap, user.EMail)

		if val, ok := tokenmap[user.EMail]; ok {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, val+" is NOT logged out")
		} else {
			w.WriteHeader(http.StatusOK)

			fmt.Fprintln(w, "You are  logged out")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "You are not logged in. Please log in first.")
	}
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	jsbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println(err)

		return
	}
	var passpak PassPack
	json.Unmarshal([]byte(string(jsbody)), &passpak)

	db, err := sql.Open("mysql", dbuser+":"+dbpass+"@/"+dbname)
	logErr(err)
	var usermail string
	err = db.QueryRow("SELECT email FROM passwordtokens WHERE token=?", passpak.Token).Scan(&usermail)

	switch {
	case err == sql.ErrNoRows:
		w.WriteHeader(http.StatusInternalServerError)

	case err != nil:

		fmt.Fprintln(w, err)
	default:
		var usersalt string
		err = db.QueryRow("SELECT salt from users WHERE email =?", usermail).Scan(&usersalt)
		logErr(err)
		// log.Println(usersalt)
		newhash := HashPass(passpak.Password, usersalt)
		// log.Println(newhash)
		// err = db.QueryRow("UPDATE users SET hash =? WHERE email =?", newhash, usermail).Scan(&usersalt)
		_, err := db.Exec("UPDATE users SET hash =? WHERE email =?", newhash, usermail)
		if err != nil {
			log.Println("a")
			log.Println(err)
		} else {
			TheQuery := "DELETE FROM passwordtokens WHERE email=\"" + string(usermail) + "\";"
			log.Println(TheQuery)
			_, err := db.Exec(TheQuery)
			if err != nil {
				log.Println("B")
				log.Fatal(err)
			} else {

				w.WriteHeader(http.StatusOK)
			}
		}

	}

	// log.Println(passpak.Token)
	// log.Println(passpak.Password)
}
func SendMail(resetemail string, message string) {

	c, err := smtp.Dial("edacmail.unm.edu:25")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	// Set the sender and recipient.
	c.Mail("nmwrap@edac.unm.edu")
	c.Rcpt(resetemail)

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}
	defer wc.Close()
	buf := bytes.NewBufferString(message)
	//"Hello, " + user.Name)
	_, err = buf.WriteTo(wc)
	logErr(err)
}

//ResetPassword will create a temporary key that a user can use to reset password.
func ResetPassword(w http.ResponseWriter, r *http.Request) {
	resetemail, err := ioutil.ReadAll(r.Body)
	logErr(err)
	if !strings.Contains(string(resetemail), " ") && strings.Contains(string(resetemail), "@") {
		// fmt.Fprintln(w, string("good"))
		db, err := sql.Open("mysql", dbuser+":"+dbpass+"@/"+dbname)
		logErr(err)
		var userid string
		err = db.QueryRow("SELECT id FROM users WHERE email=?", string(resetemail)).Scan(&userid)

		switch {
		case err == sql.ErrNoRows:
			log.Printf("No user with that ID.")

			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "No user with that address")
			return
		case err != nil:
			log.Fatal(err)
			fmt.Fprintln(w, err)
			return
		default:

			TheQuery := "DELETE FROM passwordtokens WHERE email=\"" + string(resetemail) + "\";"
			log.Println(TheQuery)
			result1, err := db.Exec(TheQuery)
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Fprintln(w, result1)

				ResetToken := RandString(200)
				TheQuery := "INSERT INTO passwordtokens (email,token)  VALUES (\"" + string(resetemail) + "\",\"" + ResetToken + "\");"
				log.Println(TheQuery)
				_, err := db.Exec(TheQuery)
				if err != nil {
					log.Fatal(err)
				} else {
					w.WriteHeader(http.StatusOK)
					eMessage := "To:" + string(resetemail) + "\r\n" +
						"Subject: NMWRAP Password Reset \r\n" +
						"\r\n" +
						"Reset your password by using the link below.\r\n" +
						"https://nmwrap.org/?token=" + ResetToken + "\r\n"

					SendMail(string(resetemail), eMessage)
					fmt.Fprintln(w, "mail sent")
					//fmt.Println(result)
					// c, err := smtp.Dial("edacmail.unm.edu:25")
					// if err != nil {
					// 	log.Fatal(err)
					// }
					// defer c.Close()
					// // Set the sender and recipient.
					// c.Mail("nmwrap@edac.unm.edu")
					// c.Rcpt(string(resetemail))

					// // Send the email body.
					// wc, err := c.Data()
					// if err != nil {
					// 	log.Fatal(err)
					// }
					// defer wc.Close()
					// buf := bytes.NewBufferString("To:" + string(resetemail) + "\r\n" +
					// 	"Subject: NMWRAP Password Reset \r\n" +
					// 	"\r\n" +
					// 	"https://nmwrap.org/?token=" + ResetToken + "\r\n")
					// //"Hello, " + user.Name)
					// _, err = buf.WriteTo(wc)
					// logErr(err)
					// if err != nil {
					// 	log.Fatal(err)
					// } else {
					// 	w.WriteHeader(http.StatusOK)
					// 	fmt.Fprintln(w, "mail sent")
					// }
				}
			}
		}

	} else {
		fmt.Fprintln(w, string("bad"))

	}
	//	fmt.Fprintln(w, string(resetemail))
}

func logerr(err error) {
	if err != nil {
		log.Println(err)
	}
}
func IsLoggedIn(r *http.Request) bool {
	var cookie, err = r.Cookie("nmwrap")
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(cookie)
		logerr(err)
		//data, err := base64.StdEncoding.DecodeString(cookie.Value)
		logerr(err)
		token, err := jwt.ParseWithClaims(cookie.Value, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		logErr(err)
		if claims, ok := token.Claims.(*UserClaims); ok && token.Valid && cookie.Value == tokenmap[claims.EMail] {
			return true

		} else {
			return false
		}
	}
	return false
}

//LoggedIn is used for clients to take no action, but check if the server has logged them out.
func LoggedIn(w http.ResponseWriter, r *http.Request) {

	if IsLoggedIn(r) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Logged in!")
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "You are not logged in. Please log in first.")
	}

}

func CheckReset(w http.ResponseWriter, r *http.Request) {
	tok, err := ioutil.ReadAll(r.Body)
	logErr(err)
	token := string(tok)
	db, err := sql.Open("mysql", dbuser+":"+dbpass+"@/"+dbname)
	logErr(err)
	var usermail string
	err = db.QueryRow("SELECT email FROM passwordtokens WHERE token=?", token).Scan(&usermail)

	switch {
	case err == sql.ErrNoRows:
		logErr(err)
		fmt.Fprintln(w, "False")
	case err != nil:

		fmt.Fprintln(w, err)
	default:
		fmt.Fprintln(w, "True")
	}
}

// POSTGeom is how the user generates reports, by putting geom...
func POSTGeom(w http.ResponseWriter, r *http.Request) {

	jsbody, err := ioutil.ReadAll(r.Body)
	// fmt.Println(string(jsbody))
	if err != nil {
		log.Println(err)
	}
	var myGeom Geom
	json.Unmarshal(jsbody, &myGeom)
	fmt.Println(myGeom)
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Helvetica", "", 16)
	pdf.AddSpotColor("PANTONE 145 CVC", 0, 42, 100, 25)
	pdf.AddPage()
	pdf.SetMargins(10, 10, 10)
	pdf.SetFillSpotColor("PANTONE 145 CVC", 90)
	pdf.Rect(0, 0, 210, 20, "F")
	pdf.Image("/var/lib/nmwrapreports/ziafire.png", 2, 2, 16, 16, false, "", 0, "")
	// pdf.Text(100, 10, myGeom.Title)
	//pdf.CellFormat(0, 15, "", "3", 1, "", false, 0, "")
	pdf.SetFont("Helvetica", "", 35)
	pdf.WriteAligned(0, 0, myGeom.Title, "C")
	// pdf.Text(100, 10, myGeom.Title)
	//pdf.CellFormat(190, 15,myGeom.Title , "0", 1, "CM", false, 0, "")
	pdf.SetFont("Helvetica", "", 16)
	CARLow := 0
	CARMed := 0
	CARHigh := 0
	fname, _ := randomFilename()
	for layernum := 0; layernum < 6; layernum++ {
		queryurl := "https://edacarc.unm.edu/arcgis/rest/services/NMWRAP/NMWRAP/MapServer/" + strconv.Itoa(layernum) + "/query?where=&text=&objectIds=&time=&geometryType=esriGeometryPolygon&inSR=&spatialRel=esriSpatialRelIntersects&relationParam=&outFields=*&returnGeometry=false&returnTrueCurves=false&maxAllowableOffset=&geometryPrecision=&outSR=&returnIdsOnly=false&returnCountOnly=false&orderByFields=&groupByFieldsForStatistics=&outStatistics=&returnZ=false&returnM=false&gdbVersion=&returnDistinctValues=false&resultOffset=&resultRecordCount=&f=pjson&geometry="
		queryurl = queryurl + string(jsbody)
		queryurl = strings.Replace(queryurl, " ", "%20", -1)
		fmt.Println(queryurl)
		resp, err := http.Get(queryurl)
		if err != nil {
			log.Println(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		// fmt.Println(string(body))
		if err != nil {
			log.Println(err)
		}
		switch lnum := layernum; lnum {
		case 0:
			pdf.Ln(15)
			pdf.SetFont("Helvetica", "", 20)
			pdf.CellFormat(190, 15, "Fire stations", "0", 1, "CM", false, 0, "")
			pdf.SetFont("Helvetica", "", 11)
			var myJSON FireStations
			json.Unmarshal(body, &myJSON)
			//fmt.Println("asaa")
			//fmt.Println(len(myJSON.Features))
			FireStationsBlurb := "There are " + strconv.Itoa(len(myJSON.Features)) + " fire stations in this area. The proximity of fire stations is essential to an assessment of fire safety. Lorem ipsum dolor sit amet, consectetur adipiscing elit.\nInteger nec odio. Praesent libero. Sed cursus ante dapibus diam. Sed nisi. Nulla quis sem at nibh elementum imperdiet. Duis sagittis ipsum. Praesent mauris. Fusce nec tellus sed augue semper porta. Mauris massa. Vestibulum lacinia arcu eget nulla. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Curabitur sodales ligula in libero."
			lines := pdf.SplitLines([]byte(FireStationsBlurb), 190.0)
			_, lineHt := pdf.GetFontSize()

			for _, line := range lines {
				pdf.CellFormat(190.0, lineHt, string(line), "", 1, "TL", false, 0, "")
			}
			if len(myJSON.Features) > 0 {
				pdf.CellFormat(0, 5, "", "3", 1, "", false, 0, "")
				pdf.SetFont("Helvetica", "", 12)
				pdf.CellFormat(45, 7, "Address", "1", 0, "C", true, 0, "")
				pdf.CellFormat(29, 7, "City", "1", 0, "C", true, 0, "")
				pdf.CellFormat(116, 7, "Name", "1", 0, "C", true, 0, "")
				pdf.Ln(-1)
				pdf.SetFont("Helvetica", "", 6)
			}
			for _, element := range myJSON.Features {
				// fmt.Println(element.Attributes.ADDRESS)
				pdf.CellFormat(45, 7, element.Attributes.ADDRESS, "1", 0, "", false, 0, "")
				pdf.CellFormat(29, 7, element.Attributes.CITY, "1", 0, "", false, 0, "")
				pdf.CellFormat(116, 7, element.Attributes.INSTNAME, "1", 0, "", false, 0, "")
				pdf.Ln(-1)
			}

		case 1:
			pdf.Ln(15)
			pdf.SetFont("Helvetica", "", 20)
			pdf.CellFormat(190, 15, "Communities At Risk", "0", 1, "CM", false, 0, "")

			pdf.Ln(5)

			var myJSON CommunitesatRisk
			pdf.SetFont("Helvetica", "", 11)
			json.Unmarshal(body, &myJSON)
			CommunitesatRiskBlurb := strconv.Itoa(len(myJSON.Features)) + " communites at risk were found in this area. Lorem ipsum dolor sit amet, consectetur adipiscing elit.\nInteger nec odio. Praesent libero. Sed cursus ante dapibus diam. Sed nisi. Nulla quis sem at nibh elementum imperdiet. Duis sagittis ipsum. Praesent mauris. Fusce nec tellus sed augue semper porta. Mauris massa. Vestibulum lacinia arcu eget nulla. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Curabitur sodales ligula in libero."
			lines := pdf.SplitLines([]byte(CommunitesatRiskBlurb), 190.0)
			_, lineHt := pdf.GetFontSize()
			for _, line := range lines {
				pdf.CellFormat(190.0, lineHt, string(line), "", 1, "TL", false, 0, "")
			}
			if len(myJSON.Features) > 0 {
				pdf.Ln(3)
				pdf.SetFont("Helvetica", "", 12)
				pdf.CellFormat(64, 7, "Name", "1", 0, "C", true, 0, "")
				pdf.CellFormat(63, 7, "County", "1", 0, "C", true, 0, "")
				pdf.CellFormat(63, 7, "Rate", "1", 0, "C", true, 0, "")
				pdf.Ln(-1)
			}
			pdf.SetFont("Helvetica", "", 6)

			json.Unmarshal(body, &myJSON)

			for _, element := range myJSON.Features {
				// fmt.Println(element.Attributes.NAME1)
				rate := ""
				if element.Attributes.Rate2016 == "L" {
					rate = "Low Risk"
					CARLow++
				} else if element.Attributes.Rate2016 == "M" {
					rate = "Medium Risk"
					CARMed++
				} else if element.Attributes.Rate2016 == "H" {
					rate = "High Risk"
					CARHigh++
				}
				pdf.CellFormat(64, 7, element.Attributes.NAME, "1", 0, "", false, 0, "")
				pdf.CellFormat(63, 7, element.Attributes.NAME1, "1", 0, "", false, 0, "")
				pdf.CellFormat(63, 7, rate, "1", 0, "", false, 0, "")
				pdf.Ln(-1)
			}
			fmt.Println(CARLow)
			fmt.Println(CARMed)
			fmt.Println(CARHigh)

			ColorGreen := drawing.Color{R: 0, G: 217, B: 101, A: 255}

			ColorRed := drawing.Color{R: 217, G: 0, B: 116, A: 255}

			ColorYellow := drawing.Color{R: 217, G: 210, B: 0, A: 255}

			DefaultColors := []drawing.Color{

				ColorRed,
				ColorYellow,
				ColorGreen,
			}
			type ColorPaletteRed interface {
				BackgroundColor() drawing.Color
				BackgroundStrokeColor() drawing.Color
				CanvasColor() drawing.Color
				CanvasStrokeColor() drawing.Color
				AxisStrokeColor() drawing.Color
				TextColor() drawing.Color
				GetSeriesColor(index int) drawing.Color
			}

			fmt.Println("defc colors")
			fmt.Println(DefaultColors)
			// var myColorPalette colors.ColorPalette{}

			var colorpal ColorPaletteRed
			fmt.Println("colorpal")
			fmt.Println(colorpal)
			carsummation := 0
			if CARLow > 0 {
				carsummation = carsummation + 1
			}
			if CARMed > 0 {
				carsummation = carsummation + 1
			}
			if CARHigh > 0 {
				carsummation = carsummation + 1
			}
			if carsummation > 0 {
				pie := chart.PieChart{
					Title: "Communities At Risk",
					// ColorPalette : colorpal ,
					Width:  512,
					Height: 512,
					Values: []chart.Value{
						{Value: float64(CARHigh), Label: strconv.Itoa(CARHigh) + " High Risk"},
						{Value: float64(CARMed), Label: strconv.Itoa(CARMed) + " Medium Risk"},
						{Value: float64(CARLow), Label: strconv.Itoa(CARLow) + " Low Risk"},
					},
				}
				fmt.Println("lol")
				fmt.Println(pie)

				// s := reflect.ValueOf(&pie).Elem()
				// typeOfT := s.Type()

				// for i := 0; i < s.NumField(); i++ {
				// 	f := s.Field(i)
				// 	fmt.Printf("%d: %s %s = %v\n", i,
				// 		typeOfT.Field(i).Name, f.Type(), f.Interface())
				// }

				var lol chart.ColorPalette
				fmt.Println(lol)
				buffer := bytes.NewBuffer([]byte{})
				err = pie.Render(chart.PNG, buffer)
				if err != nil {
					fmt.Printf("Error rendering pie chart: %v\n", err)
				}
				piereader := bufio.NewReader(buffer)

				var options gofpdf.ImageOptions
				options.ImageType = "PNG"
				fmt.Println(options)
				fmt.Println(pdf.GetPageSize())
				pdf.RegisterImageOptionsReader("piechart", options, piereader)

				whatwegot := 297.0
				whatweneed := pdf.GetY() + 128
				fmt.Println(whatwegot)
				fmt.Println(whatweneed)
				if whatweneed > whatwegot {
					extrapadding := 297.0 - pdf.GetY()

					pdf.CellFormat(10, extrapadding, "", "0", 0, "", false, 0, "")
					fmt.Println("lol")

				}
				CurrentX := pdf.GetX()
				CurrentY := pdf.GetY()
				if pdf.Ok() {

					pdf.Image("piechart", CurrentX+31, CurrentY, 128, 128, false, "", 0, "")

					pdf.SetY(CurrentY + 128)
				}
				fmt.Println(pdf.GetPageSize())
				fmt.Println(pdf.GetMargins())
				fmt.Println(CurrentY)
			}

		case 2:
			//IncorporatedCityBoundaries
			pdf.Ln(15)
			pdf.SetFont("Helvetica", "", 20)
			pdf.CellFormat(190, 15, "Incorporated City Boundaries", "0", 1, "CM", false, 0, "")

			pdf.Ln(5)

			pdf.CellFormat(0, 5, "", "3", 1, "", false, 0, "")
			var myJSON IncorporatedCityBoundaries
			pdf.SetFont("Helvetica", "", 11)
			json.Unmarshal(body, &myJSON)
			IncorporatedCityBoundariesBlurb := "Incorporated City Boundaries Blurb Lorem ipsum dolor sit amet, consectetur adipiscing elit.\nInteger nec odio. Praesent libero. Sed cursus ante dapibus diam. Sed nisi. Nulla quis sem at nibh elementum imperdiet. Duis sagittis ipsum. Praesent mauris. Fusce nec tellus sed augue semper porta. Mauris massa. Vestibulum lacinia arcu eget nulla. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Curabitur sodales ligula in libero."
			lines := pdf.SplitLines([]byte(IncorporatedCityBoundariesBlurb), 190.0)
			_, lineHt := pdf.GetFontSize()
			for _, line := range lines {
				pdf.CellFormat(190.0, lineHt, string(line), "", 1, "TL", false, 0, "")
			}
			if len(myJSON.Features) > 0 {
				pdf.Ln(3)
				pdf.SetFont("Helvetica", "", 12)
				pdf.CellFormat(95, 7, "Name", "1", 0, "C", true, 0, "")
				pdf.CellFormat(95, 7, "Area", "1", 0, "C", true, 0, "")

				pdf.Ln(-1)
			}
			pdf.SetFont("Helvetica", "", 6)

			json.Unmarshal(body, &myJSON)

			for _, element := range myJSON.Features {
				pdf.CellFormat(95, 7, element.Attributes.NAME10, "1", 0, "", false, 0, "")
				pdf.CellFormat(95, 7, strconv.FormatFloat(element.Attributes.ShapeArea, 'E', -1, 64), "1", 0, "", false, 0, "")

				pdf.Ln(-1)
			}
		case 3:
			//
			pdf.Ln(15)
			pdf.SetFont("Helvetica", "", 20)
			pdf.CellFormat(190, 15, "Vegetation Treatments", "0", 1, "CM", false, 0, "")

			pdf.Ln(5)

			var myJSON VegetationTreatments
			pdf.SetFont("Helvetica", "", 11)
			json.Unmarshal(body, &myJSON)
			VegetationTreatmentsBlurb := "Vegetation Treatments Blurb Lorem ipsum dolor sit amet, consectetur adipiscing elit.\nInteger nec odio. Praesent libero. Sed cursus ante dapibus diam. Sed nisi. Nulla quis sem at nibh elementum imperdiet. Duis sagittis ipsum. Praesent mauris. Fusce nec tellus sed augue semper porta. Mauris massa. Vestibulum lacinia arcu eget nulla. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Curabitur sodales ligula in libero."
			lines := pdf.SplitLines([]byte(VegetationTreatmentsBlurb), 190.0)
			_, lineHt := pdf.GetFontSize()
			for _, line := range lines {
				pdf.CellFormat(190.0, lineHt, string(line), "", 1, "TL", false, 0, "")
			}
			if len(myJSON.Features) > 0 {
				pdf.Ln(3)
				pdf.SetFont("Helvetica", "", 12)
				pdf.CellFormat(64, 7, "Description", "1", 0, "C", true, 0, "")
				pdf.CellFormat(63, 7, "NameProj", "1", 0, "C", true, 0, "")
				pdf.CellFormat(63, 7, "Partners", "1", 0, "C", true, 0, "")
				pdf.Ln(-1)
			}
			pdf.SetFont("Helvetica", "", 6)

			json.Unmarshal(body, &myJSON)

			for _, element := range myJSON.Features {
				pdf.CellFormat(64, 7, element.Attributes.Description, "1", 0, "", false, 0, "")
				pdf.CellFormat(63, 7, element.Attributes.NameProj, "1", 0, "", false, 0, "")
				pdf.CellFormat(63, 7, element.Attributes.Partners, "1", 0, "", false, 0, "")

				pdf.Ln(-1)
			}
		case 4:
			//
			pdf.Ln(15)
			pdf.SetFont("Helvetica", "", 20)
			pdf.CellFormat(190, 15, "Watersheds HUC8", "0", 1, "CM", false, 0, "")

			pdf.Ln(5)

			var myJSON WatershedsHUC8
			pdf.SetFont("Helvetica", "", 11)
			json.Unmarshal(body, &myJSON)
			WatershedsHUC8Blurb := "Watersheds HUC8 Blurb Lorem ipsum dolor sit amet, consectetur adipiscing elit.\nInteger nec odio. Praesent libero. Sed cursus ante dapibus diam. Sed nisi. Nulla quis sem at nibh elementum imperdiet. Duis sagittis ipsum. Praesent mauris. Fusce nec tellus sed augue semper porta. Mauris massa. Vestibulum lacinia arcu eget nulla. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Curabitur sodales ligula in libero."
			lines := pdf.SplitLines([]byte(WatershedsHUC8Blurb), 200.0)
			_, lineHt := pdf.GetFontSize()

			for _, line := range lines {
				pdf.CellFormat(190.0, lineHt, string(line), "", 1, "TL", false, 0, "")
			}
			if len(myJSON.Features) > 0 {
				pdf.Ln(3)
				pdf.SetFont("Helvetica", "", 12)
				pdf.CellFormat(64, 7, "Name", "1", 0, "C", true, 0, "")
				pdf.CellFormat(63, 7, "TNC Ranking", "1", 0, "C", true, 0, "")
				pdf.CellFormat(63, 7, "State", "1", 0, "C", true, 0, "")
				pdf.Ln(-1)
				pdf.SetFont("Helvetica", "", 6)
			}

			pdf.SetFont("Helvetica", "", 6)

			json.Unmarshal(body, &myJSON)

			for _, element := range myJSON.Features {
				pdf.CellFormat(64, 7, element.Attributes.NAME, "1", 0, "", false, 0, "")
				pdf.CellFormat(63, 7, strconv.Itoa(element.Attributes.TNCRanking), "1", 0, "", false, 0, "")
				pdf.CellFormat(63, 7, element.Attributes.STATES, "1", 0, "", false, 0, "")

				pdf.Ln(-1)
			}
		case 5:
			//County
			pdf.Ln(15)
			pdf.SetFont("Helvetica", "", 20)
			pdf.CellFormat(190, 15, "Counties", "0", 1, "CM", false, 0, "")

			pdf.Ln(5)

			var myJSON County
			pdf.SetFont("Helvetica", "", 11)
			json.Unmarshal(body, &myJSON)
			CountyBlurb := "County Blurb Lorem ipsum dolor sit amet, consectetur adipiscing elit.\nInteger nec odio. Praesent libero. Sed cursus ante dapibus diam. Sed nisi. Nulla quis sem at nibh elementum imperdiet. Duis sagittis ipsum. Praesent mauris. Fusce nec tellus sed augue semper porta. Mauris massa. Vestibulum lacinia arcu eget nulla. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Curabitur sodales ligula in libero."
			lines := pdf.SplitLines([]byte(CountyBlurb), 190.0)
			_, lineHt := pdf.GetFontSize()
			for _, line := range lines {
				pdf.CellFormat(190.0, lineHt, string(line), "", 1, "TL", false, 0, "")
			}
			if len(myJSON.Features) > 0 {
				pdf.Ln(3)
				pdf.SetFont("Helvetica", "", 12)
				pdf.CellFormat(190, 7, "County", "1", 0, "C", true, 0, "")

				pdf.Ln(-1)
			}
			pdf.SetFont("Helvetica", "", 6)

			json.Unmarshal(body, &myJSON)

			for _, element := range myJSON.Features {
				pdf.CellFormat(190, 7, element.Attributes.NAMELSAD, "1", 0, "", false, 0, "")

				pdf.Ln(-1)
			}
		default:
		}

	}

	err = pdf.OutputFileAndClose("/tmp/" + fname + ".pdf")
	if err != nil {
		log.Println(err)
	} else {
		fmt.Fprintln(w, fname)
	}

}

//GetReport shows the generated report.
func GetReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println(mux.Vars(r))
	key := vars["key"]
	fullpath := "/tmp/" + key + ".pdf"
	fname := vars["fname"]
	log.Println(fullpath)
	w.Header().Set("Content-Disposition", "attachment; filename="+fname)
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	http.ServeFile(w, r, fullpath)

}

var (
	// ColorWhite is white.
	ColorWhite = drawing.Color{R: 255, G: 255, B: 255, A: 255}
	// ColorBlue is the basic theme blue color.
	ColorBlue = drawing.Color{R: 0, G: 116, B: 217, A: 255}
	// ColorCyan is the basic theme cyan color.
	ColorCyan = drawing.Color{R: 0, G: 217, B: 210, A: 255}
	// ColorGreen is the basic theme green color.
	ColorGreen = drawing.Color{R: 0, G: 217, B: 101, A: 255}
	// ColorRed is the basic theme red color.
	ColorRed = drawing.Color{R: 217, G: 0, B: 116, A: 255}
	// ColorOrange is the basic theme orange color.
	ColorOrange = drawing.Color{R: 217, G: 101, B: 0, A: 255}
	// ColorYellow is the basic theme yellow color.
	ColorYellow = drawing.Color{R: 217, G: 210, B: 0, A: 255}
	// ColorBlack is the basic theme black color.
	ColorBlack = drawing.Color{R: 51, G: 51, B: 51, A: 255}
	// ColorLightGray is the basic theme light gray color.
	ColorLightGray = drawing.Color{R: 239, G: 239, B: 239, A: 255}

	// ColorAlternateBlue is a alternate theme color.
	ColorAlternateBlue = drawing.Color{R: 106, G: 195, B: 203, A: 255}
	// ColorAlternateGreen is a alternate theme color.
	ColorAlternateGreen = drawing.Color{R: 42, G: 190, B: 137, A: 255}
	// ColorAlternateGray is a alternate theme color.
	ColorAlternateGray = drawing.Color{R: 110, G: 128, B: 139, A: 255}
	// ColorAlternateYellow is a alternate theme color.
	ColorAlternateYellow = drawing.Color{R: 240, G: 174, B: 90, A: 255}
	// ColorAlternateLightGray is a alternate theme color.
	ColorAlternateLightGray = drawing.Color{R: 187, G: 190, B: 191, A: 255}

	// ColorTransparent is a transparent (alpha zero) color.
	ColorTransparent = drawing.Color{R: 1, G: 1, B: 1, A: 0}
)

type defaultColorPalette struct{}

func (dp defaultColorPalette) CanvasColorRed() drawing.Color {
	DefaultCanvasColor := ColorWhite
	return DefaultCanvasColor
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func GetReportFromUpload(w http.ResponseWriter, r *http.Request) {
	if IsLoggedIn(r) {
		AllowdShapeExtensions := []string{"cpg", "dbf", "prj", "sbn", "sbx", "shp", "shx"}

		//var Buf bytes.Buffer
		// in your case file would be fileupload
		fmt.Println(r.FormFile)
		file, header, err := r.FormFile("file")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		name := strings.Split(header.Filename, ".")
		fmt.Printf("File name %s\n", name[0])
		RandomFileName := RandString(10)
		ZipFile := "/tmp/" + RandomFileName + ".zip"
		out, err := os.Create(ZipFile)
		if err != nil {
			fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
			return
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			fmt.Fprintln(w, err)
		}
		reader, err := zip.OpenReader(ZipFile)
		if err != nil {

			log.Fatal(err)

		}

		defer reader.Close()
		dest := "/tmp/" + RandomFileName
		os.MkdirAll(dest, 755)
		for _, f := range reader.File {
			fmt.Println(f.Name)
			extension := strings.Split(f.Name, ".")
			path := dest + "/" + f.Name
			if stringInSlice(extension[1], AllowdShapeExtensions) {
				rc, err := f.Open()
				logErr(err)
				defer rc.Close()
				f, err := os.OpenFile(
					path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
				logErr(err)
				defer f.Close()

				_, err = io.Copy(f, rc)
				logErr(err)
			}
		}
		Shapefile := "/tmp/" + RandomFileName + "/" + name[0] + ".shp"
		myshape, err := shp.Open(Shapefile)
		if err != nil {
			log.Fatal(err)
		}
		defer myshape.Close()

		fields := myshape.Fields()
		fmt.Println(fields)

		for myshape.Next() {
			n, p := myshape.Shape()

			// print feature
			fmt.Println(reflect.TypeOf(p).Elem(), p.BBox())
			fmt.Println(n)

		}
		spatialRef := gdal.CreateSpatialReference("")
		spatialRef.FromEPSG(3857)
		srString, err := spatialRef.ToWKT()
		fmt.Println(srString)

		return
	} else {
		fmt.Println("lasdf")
		driver := gdal.OGRDriverByName("ESRI Shapefile")
		fmt.Println(driver)
		datasource, _ := driver.Open("/tmp/rfBd56ti2SMtYvSgD5xAV0YU99zampta7Z7S575KLkIZ9PYkL17LTlsVqMNTZyLKMIFSD2x28MlgPJ0SDZVHnHJPxMKi0tWxu3pQJ71N5GWfOIGTdSWXbRLGAwD1IkzuZ5G1pEDzqqm3sncCYry01AuHiK7FDcCc35S4IzoOjgm2v8KyBpNlS52DyhMEXiJev6e8bqQK/reportgeom.shp", 0)

		fmt.Println(datasource.LayerCount())
		layer := datasource.LayerByIndex(0)
		myfeature := layer.Feature(0)
		geom := myfeature.Geometry()
		spatialRef := gdal.CreateSpatialReference("PROJCS[\"WGS 84 / Pseudo-Mercator\",GEOGCS[\"WGS 84\",DATUM[\"WGS_1984\",SPHEROID[\"WGS 84\",6378137,298.257223563,AUTHORITY[\"EPSG\",\"7030\"]],AUTHORITY[\"EPSG\",\"6326\"]],PRIMEM[\"Greenwich\",0,AUTHORITY[\"EPSG\",\"8901\"]],UNIT[\"degree\",0.0174532925199433,AUTHORITY[\"EPSG\",\"9122\"]],AUTHORITY[\"EPSG\",\"4326\"]],PROJECTION[\"Mercator_1SP\"],PARAMETER[\"central_meridian\",0],PARAMETER[\"scale_factor\",1],PARAMETER[\"false_easting\",0],PARAMETER[\"false_northing\",0],UNIT[\"metre\",1,AUTHORITY[\"EPSG\",\"9001\"]],AXIS[\"X\",EAST],AXIS[\"Y\",NORTH],EXTENSION[\"PROJ4\",\"+proj=merc +a=6378137 +b=6378137 +lat_ts=0.0 +lon_0=0.0 +x_0=0.0 +y_0=0 +k=1.0 +units=m +nadgrids=@null +wktext  +no_defs\"],AUTHORITY[\"EPSG\",\"3857\"]]")
		//spatialRef.FromEPSG(3857)
		geom.TransformTo(spatialRef)
		fmt.Println(geom.ToWKT())
		fmt.Println(geom.ToGML())
		fmt.Println(geom.ToJSON())
		fmt.Println(geom.ToKML())
		var myGeoJSON GeoJSON
		json.Unmarshal([]byte(geom.ToJSON()), &myGeoJSON)
		fmt.Println(reflect.TypeOf(myGeoJSON.Coordinates[0]))
		var myGeom Geom
		myGeom.Title = "test"
		myGeom.Rings = myGeoJSON.Coordinates
		res1b, _ := json.Marshal(myGeom)
		fmt.Println(string(res1b))
		// for _, v := range myGeoJSON.Coordinates[0] {
		// 	fmt.Println(v[0])
		// 	fmt.Println(v[0])

		// }
		//NextFeature
		return
	}
}
