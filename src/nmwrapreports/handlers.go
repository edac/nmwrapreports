package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	gdal "github.com/hbarrett/gdal"
	"github.com/jonas-p/go-shp"
	"github.com/jung-kurt/gofpdf"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	//"reflect"
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	_ "github.com/go-sql-driver/mysql"
	mathrand "math/rand"
	"net/smtp"
	"net/url"
	"strconv"
	"strings"
	"time"
)

//   ___ _____ ___ _   _  ___ _____   ___  ___ ___ ___
//  / __|_   _| _ \ | | |/ __|_   _| |   \| __| __/ __|
//  \__ \ | | |   / |_| | (__  | |   | |) | _|| _|\__ \
//  |___/ |_| |_|_\\___/ \___| |_|   |___/|___|_| |___/

//defaultColorPalette for use in reports.
type defaultColorPalette struct{}

//frontpagedict for index that will never be used.
type frontpagedict struct {
	MESSAGE string
}

//JobStatus is the struct returned when requesting an ArcGIS GP job.
type JobStatus struct {
	JobID     string `json:"jobId"`
	JobStatus string `json:"jobStatus"`
}

//JobSuccess is the struct returned when checking the status an ArcGIS GP job.
type JobSuccess struct {
	JobID     string `json:"jobId"`
	JobStatus string `json:"jobStatus"`
	Results   struct {
		OutputZipFile struct {
			ParamURL string `json:"paramUrl"`
		} `json:"Output_Zip_File"`
	} `json:"results"`
	Inputs struct {
		LayersToClip struct {
			ParamURL string `json:"paramUrl"`
		} `json:"Layers_to_Clip"`
		AreaOfInterest struct {
			ParamURL string `json:"paramUrl"`
		} `json:"Area_of_Interest"`
		FeatureFormat struct {
			ParamURL string `json:"paramUrl"`
		} `json:"Feature_Format"`
		RasterFormat struct {
			ParamURL string `json:"paramUrl"`
		} `json:"Raster_Format"`
		SpatialReference struct {
			ParamURL string `json:"paramUrl"`
		} `json:"Spatial_Reference"`
	} `json:"inputs"`
	Messages []interface{} `json:"messages"`
}

//ZipStatus is the struct returned when the GP extract job is successful.
type ZipStatus struct {
	ParamName string `json:"paramName"`
	DataType  string `json:"dataType"`
	Value     struct {
		URL string `json:"url"`
	} `json:"value"`
}

//UserHistory for storing and requesing a users geometry requests.
type UserHistory struct {
	ID    string `json:"id"`
	Geom  Geom   `json:"geom"`
	Title string `json:"title"`
}

//Userpack is used to process auth
type Userpack struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//CreateUserpack has name instead of password for creating a new user. I could probably merge these two...
type CreateUserpack struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

//PassPack For processing if a user is logged in.
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
	Rings   [][][]float64 `json:"rings"`
	Title   string        `json:"title"`
	History bool          `json:"history"`
}

// GeoJSON for unpacking user supplied geom
type GeoJSON struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}

//   ___ ___ _  _ ___ ___    _   _
//  / __| __| \| | __| _ \  /_\ | |
// | (_ | _|| .` | _||   / / _ \| |__
//  \___|___|_|\_|___|_|_\/_/ \_\____|

//RandString creates "random" alphanumeric string
func RandString(n int) string {
	var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[mathrand.Intn(len(letterRunes))]
	}
	return string(b)
}

//randomFilename could probably be replaced with RandString...
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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

//DeleteHistory - Allows users to delete geom history
func DeleteHistory(w http.ResponseWriter, r *http.Request) {
	if IsLoggedIn(r) {
		fmt.Println("GETTING HERERE!")
		userdata, err := UserData(r)
		logErr(err)
		fmt.Println(userdata["id"])
		jsbody, err := ioutil.ReadAll(r.Body)
		var userhist UserHistory
		if err := json.Unmarshal(jsbody, &userhist); err != nil {
			log.Println(err)
		}
		fmt.Println(userhist.ID)

		db, err := sql.Open("mysql", dbuser+":"+dbpass+"@/"+dbname)
		if err != nil {

			log.Println(err)

		}
		defer db.Close()
		err = db.Ping()
		if err != nil {

			log.Println(err)

		}
		logErr(err)
		var userid int
		q := "SELECT userid FROM areasofinterest WHERE id='" + userhist.ID + "'"
		fmt.Println(q)
		err = db.QueryRow(q).Scan(&userid)
		logErr(err)

		fmt.Println("a")
		fmt.Println(strconv.Itoa(userid))
		if userdata["id"] == strconv.Itoa(userid) {
			DeleteQuery := "DELETE FROM areasofinterest WHERE id=\"" + userhist.ID + "\";"
			fmt.Println(DeleteQuery)
			_, err := db.Exec(DeleteQuery)
			if err != nil {
				log.Fatal(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, "Failed to delete id"+userhist.ID)
			} else {
				w.WriteHeader(http.StatusOK)
			}

		}

	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "You must be logged in!")
	}

}

//History - Return JSON of users geom history
func History(w http.ResponseWriter, r *http.Request) {
	if IsLoggedIn(r) {
		fmt.Println("start")
		userdata, err := UserData(r)
		logErr(err)
		fmt.Println(userdata)
		db, err := sql.Open("mysql", dbuser+":"+dbpass+"@/"+dbname)
		if err != nil {

			log.Println(err)

		}
		defer db.Close()
		err = db.Ping()
		if err != nil {

			log.Println(err)

		}
		logErr(err)
		q := "SELECT * FROM areasofinterest WHERE userid='" + userdata["id"] + "'"
		fmt.Println(q)
		rows, err := db.Query(q)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err)
		}
		var id int
		var userid int
		var geom string
		var title string
		var jsonText = []byte(`[]`)
		var userhistory []UserHistory
		if err := json.Unmarshal([]byte(jsonText), &userhistory); err != nil {
			log.Println(err)
		}

		for rows.Next() {

			err = rows.Scan(&id, &userid, &geom, &title)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, err)
			}
			var usergeom Geom
			if err := json.Unmarshal([]byte(geom), &usergeom); err != nil {
				log.Println(err)
			}
			userhistory = append(userhistory, UserHistory{ID: strconv.Itoa(id), Geom: usergeom, Title: title})
		}
		newjson, _ := json.Marshal(userhistory)
		w.Header().Set("Content-Type", "application/json")
		w.Write(newjson)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "You must be logged in!")
	}
}

//  ___ ___ ___  ___  ___ _____ ___
// | _ \ __| _ \/ _ \| _ \_   _/ __|
// |   / _||  _/ (_) |   / | | \__ \
// |_|_\___|_|  \___/|_|_\ |_| |___/

// POSTGeom is how the user generates reports, by putting geom...
func POSTGeom(w http.ResponseWriter, r *http.Request) {
	if IsLoggedIn(r) {
		jsbody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}
		var myGeom Geom
		json.Unmarshal(jsbody, &myGeom)
		fmt.Println(myGeom)
		myRings, _ := json.Marshal(myGeom.Rings)
		fmt.Println(string(myRings))

		fname, err := ReportGen(myGeom, r)
		logErr(err)
		fmt.Fprintln(w, fname)
	} else {

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "You must be logged in!")

	}
}

//POSTGeomForExtract - Route to pre-process geom for ExtractGen
func POSTGeomForExtract(w http.ResponseWriter, r *http.Request) {
	if IsLoggedIn(r) {
		jsbody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}
		var myGeom Geom
		json.Unmarshal(jsbody, &myGeom)
		//fmt.Println(myGeom)
		msg, err := ExtractGen(myGeom, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, msg)
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, msg)
		}
		// fmt.Fprintln(w, fname)
	} else {

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "You must be logged in!")

	}
}

//ReportGen - Generate pdf report from geom
func ReportGen(myGeom Geom, r *http.Request) (string, error) {
	user := GetCookieParts(r)
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Helvetica", "", 16)
	pdf.AddSpotColor("PANTONE 145 CVC", 0, 42, 100, 25)
	pdf.AddPage()
	pdf.SetMargins(10, 10, 10)
	pdf.SetFillSpotColor("PANTONE 145 CVC", 90)
	pdf.Rect(0, 0, 210, 20, "F")
	pdf.Image("/var/lib/nmwrapreports/ziafire.png", 2, 2, 16, 16, false, "", 0, "")
	pdf.SetFont("Helvetica", "", 35)
	pdf.WriteAligned(0, 0, myGeom.Title, "C")
	pdf.SetFont("Helvetica", "", 16)
	CARLow := 0
	CARMed := 0
	CARHigh := 0
	fname, _ := randomFilename()
	myGeomMarshal, _ := json.Marshal(myGeom)
	for layernum := 0; layernum < 6; layernum++ {
		queryurl := "https://edacarc.unm.edu/arcgis/rest/services/NMWRAP/NMWRAP/MapServer/" + strconv.Itoa(layernum) + "/query" //?where=&text=&objectIds=&time=&geometryType=esriGeometryPolygon&inSR=&spatialRel=esriSpatialRelIntersects&relationParam=&outFields=*&returnGeometry=false&returnTrueCurves=false&maxAllowableOffset=&geometryPrecision=&outSR=&returnIdsOnly=false&returnCountOnly=false&orderByFields=&groupByFieldsForStatistics=&outStatistics=&returnZ=false&returnM=false&gdbVersion=&returnDistinctValues=false&resultOffset=&resultRecordCount=&f=pjson&geometry="
		resp, err := http.PostForm(queryurl, url.Values{
			"f":                    {"pjson"},
			"geometry":             {string(myGeomMarshal)},
			"geometryType":         {"esriGeometryPolygon"},
			"outFields":            {"*"},
			"returnCountOnly":      {"false"},
			"returnDistinctValues": {"false"},
			"returnGeometry":       {"false"},
			"returnIdsOnly":        {"false"},
			"returnM":              {"false"},
			"returnTrueCurves":     {"false"},
			"returnZ":              {"false"},
			"spatialRel":           {"esriSpatialRelIntersects"}})

		if err != nil {
			log.Println(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
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
			type ColorPaletteRed interface {
				BackgroundColor() drawing.Color
				BackgroundStrokeColor() drawing.Color
				CanvasColor() drawing.Color
				CanvasStrokeColor() drawing.Color
				AxisStrokeColor() drawing.Color
				TextColor() drawing.Color
				GetSeriesColor(index int) drawing.Color
			}

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

				buffer := bytes.NewBuffer([]byte{})
				err = pie.Render(chart.PNG, buffer)
				if err != nil {
					fmt.Printf("Error rendering pie chart: %v\n", err)
				}
				piereader := bufio.NewReader(buffer)

				var options gofpdf.ImageOptions
				options.ImageType = "PNG"

				pdf.RegisterImageOptionsReader("piechart", options, piereader)

				whatwegot := 297.0
				whatweneed := pdf.GetY() + 128

				if whatweneed > whatwegot {
					extrapadding := 297.0 - pdf.GetY()

					pdf.CellFormat(10, extrapadding, "", "0", 0, "", false, 0, "")

				}
				CurrentX := pdf.GetX()
				CurrentY := pdf.GetY()
				if pdf.Ok() {

					pdf.Image("piechart", CurrentX+31, CurrentY, 128, 128, false, "", 0, "")

					pdf.SetY(CurrentY + 128)
				}

			}

		case 2:
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
			pdf.CellFormat(190, 15, "Counties with intersecting data", "0", 1, "CM", false, 0, "")

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

	err := pdf.OutputFileAndClose("/tmp/" + fname + ".pdf")
	if err != nil {
		log.Println(err)
	} else {
		db, err := sql.Open("mysql", dbuser+":"+dbpass+"@/"+dbname)
		if err != nil {
			log.Println(err)
		}
		defer db.Close()
		err = db.Ping()
		if err != nil {
			log.Println(err)
		}
		rows, err := db.Query("SELECT * FROM users WHERE email='" + user.EMail + "'")
		if err != nil {
			log.Println(err)
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
				log.Println(err)
			}

		}
		log.Println(json.Marshal(myGeom))
		log.Println(myGeom.History)
		if myGeom.History != true {
			geomtitle := myGeom.Title
			marstring, _ := json.Marshal(myGeom)
			AreaQuery := "INSERT INTO areasofinterest (userid,geom,title)  VALUES ('" + strconv.Itoa(id) + "','" + string(marstring) + "','" + geomtitle + "');"
			log.Println(AreaQuery)
			_, err = db.Exec(AreaQuery)
			if err != nil {
				log.Fatal(err)
			}
		}
		return fname, nil
	}
	return "Yikes", errors.New("End Of Function")
}

//GetReport shows the generated report.
func GetReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	fullpath := "/tmp/" + key + ".pdf"
	fname := vars["fname"]
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

func (dp defaultColorPalette) CanvasColorRed() drawing.Color {
	DefaultCanvasColor := ColorWhite
	return DefaultCanvasColor
}

//GetReportFromUpload - preprocesses geom for ReportGen
func GetReportFromUpload(w http.ResponseWriter, r *http.Request) {
	if IsLoggedIn(r) {
		AllowdShapeExtensions := []string{"cpg", "dbf", "prj", "sbn", "sbx", "shp", "shx"}
		file, _, err := r.FormFile("file")
		if err != nil {
			panic(err)
		}
		defer file.Close()
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
		ShapeName := ""
		defer reader.Close()
		dest := "/tmp/" + RandomFileName
		os.MkdirAll(dest, 755)
		for _, f := range reader.File {
			//	fmt.Println(f.Name)
			extension := strings.Split(f.Name, ".")
			path := dest + "/" + f.Name
			if stringInSlice(extension[1], AllowdShapeExtensions) {
				if extension[1] == "shp" {
					ShapeName = f.Name
				}

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
		Shapefile := "/tmp/" + RandomFileName + "/" + ShapeName
		myshape, err := shp.Open(Shapefile)
		if err != nil {
			log.Fatal(err)
		}
		defer myshape.Close()
		driver := gdal.OGRDriverByName("ESRI Shapefile")
		fmt.Println(driver)
		datasource, _ := driver.Open(Shapefile, 0)
		fmt.Println(datasource.LayerCount())
		layer := datasource.LayerByIndex(0)
		myfeature := layer.Feature(0)
		geom := myfeature.Geometry()
		spatialRef := gdal.CreateSpatialReference("PROJCS[\"WGS 84 / Pseudo-Mercator\",GEOGCS[\"WGS 84\",DATUM[\"WGS_1984\",SPHEROID[\"WGS 84\",6378137,298.257223563,AUTHORITY[\"EPSG\",\"7030\"]],AUTHORITY[\"EPSG\",\"6326\"]],PRIMEM[\"Greenwich\",0,AUTHORITY[\"EPSG\",\"8901\"]],UNIT[\"degree\",0.0174532925199433,AUTHORITY[\"EPSG\",\"9122\"]],AUTHORITY[\"EPSG\",\"4326\"]],PROJECTION[\"Mercator_1SP\"],PARAMETER[\"central_meridian\",0],PARAMETER[\"scale_factor\",1],PARAMETER[\"false_easting\",0],PARAMETER[\"false_northing\",0],UNIT[\"metre\",1,AUTHORITY[\"EPSG\",\"9001\"]],AXIS[\"X\",EAST],AXIS[\"Y\",NORTH],EXTENSION[\"PROJ4\",\"+proj=merc +a=6378137 +b=6378137 +lat_ts=0.0 +lon_0=0.0 +x_0=0.0 +y_0=0 +k=1.0 +units=m +nadgrids=@null +wktext  +no_defs\"],AUTHORITY[\"EPSG\",\"3857\"]]")
		geom.TransformTo(spatialRef)
		fmt.Println(geom.ToWKT())
		// fmt.Println(geom.ToGML())
		fmt.Println(geom.ToJSON())
		//fmt.Println(geom.ToKML())
		fmt.Println("ZZZ")
		fmt.Println(geom)
		var myGeoJSON GeoJSON
		json.Unmarshal([]byte(geom.ToJSON()), &myGeoJSON)
		//	fmt.Println(reflect.TypeOf(myGeoJSON.Coordinates[0]))
		var myGeom Geom
		myGeom.Title = r.FormValue("title")
		myGeom.Rings = myGeoJSON.Coordinates
		fmt.Println(myGeom)
		fname, err := ReportGen(myGeom, r)
		logErr(err)
		fmt.Println(fname)
		fmt.Fprintln(w, fname)

	} else {
		return
	}
}

//  ___   _ _____ _     _____  _______ ___    _   ___ _____ ___ ___  _  _
// |   \ /_\_   _/_\   | __\ \/ /_   _| _ \  /_\ / __|_   _|_ _/ _ \| \| |
// | |) / _ \| |/ _ \  | _| >  <  | | |   / / _ \ (__  | |  | | (_) | .` |
// |___/_/ \_\_/_/ \_\ |___/_/\_\ |_| |_|_\/_/ \_\___| |_| |___\___/|_|\_|

//ExtractMailer - redundant
func ExtractMailer(message string, recipient string) {
	emailbody := strings.Replace(DataExtractMessage, "$message", message, -1)
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: NMWRAP Data Extract Alert!\n"
	msg := []byte(subject + mime + emailbody)
	recpt := []string{recipient}
	err := smtp.SendMail("edacmail.unm.edu:25", nil, "nmwrap@edac.unm.edu", recpt, msg)
	if err != nil {
		log.Println(err)

	}
}

//ExtractJobs func to check status of jobs and alert.
func ExtractJobs() {
	fmt.Println("gigaextract")
	db, err := sql.Open("mysql", dbuser+":"+dbpass+"@/"+dbname)
	if err != nil {
		log.Println(err)
		ExtractMailer(err.Error(), "nmwrap@edac.unm.edu")
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {

		ExtractMailer(err.Error(), "nmwrap@edac.unm.edu")
	}

	rows, err := db.Query("SELECT email, jobid FROM extractjobs WHERE status='1'")

	if err != nil {
		log.Println(err)
		//Send mail to nmwrap@edac.unm.edu
		ExtractMailer(err.Error(), "nmwrap@edac.unm.edu")
	} else {

		var jobid string
		var email string

		for rows.Next() {
			err = rows.Scan(&email, &jobid)
			if err != nil {
				ExtractMailer(err.Error(), "nmwrap@edac.unm.edu")
				log.Println(err)
				//send admin mail
			} else {
				var netClient = &http.Client{
					Timeout: time.Second * 10,
				}
				url := "https://edacarc.unm.edu/arcgis/rest/services/NMWRAP/ExtractData/GPServer/Extract%20Data/jobs/" + jobid + "?f=pjson"
				resp, _ := netClient.Get(url)
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println(err)
					ExtractMailer(err.Error(), "nmwrap@edac.unm.edu")
				} else {
					var jobstat JobSuccess
					if err := json.Unmarshal(body, &jobstat); err != nil {
						log.Println(err)
						ExtractMailer(err.Error(), "nmwrap@edac.unm.edu")
					}
					if jobstat.JobStatus == "esriJobSucceeded" {
						outurl := "https://edacarc.unm.edu/arcgis/rest/services/NMWRAP/ExtractData/GPServer/Extract%20Data/jobs/" + jobid + "/results/Output_Zip_File?f=pjson"

						zipresp, _ := netClient.Get(outurl)
						zipbody, err := ioutil.ReadAll(zipresp.Body)
						if err != nil {
							fmt.Println("err")
							ExtractMailer(err.Error(), "nmwrap@edac.unm.edu")
						} else {
							var zipstatus ZipStatus
							if err := json.Unmarshal(zipbody, &zipstatus); err != nil {
								log.Println(err)
								ExtractMailer(err.Error(), "nmwrap@edac.unm.edu")
							} else {
								fmt.Println(zipstatus.Value.URL)

								emailbody := strings.Replace(DataExtractSucess, "$dlurl", zipstatus.Value.URL, -1)

								mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
								subject := "Subject: NMWRAP Data Extract\n"
								msg := []byte(subject + mime + emailbody)
								recpt := []string{email}
								err := smtp.SendMail("edacmail.unm.edu:25", nil, "nmwrap@edac.unm.edu", recpt, msg)
								if err != nil {
									log.Println(err)

								} else {
									log.Println("JobID " + jobid + " download info sent to " + email)

									TheQuery := "DELETE FROM extractjobs WHERE jobid=\"" + jobid + "\";"
									//log.Println(TheQuery)
									_, err := db.Exec(TheQuery)
									if err != nil {
										ExtractMailer(err.Error(), "nmwrap@edac.unm.edu")
									}
									//	log.Println("B")
									logErr(err)

								}

							}
						}

					} else if jobstat.JobStatus == "esriJobFailed" {
						ExtractMailer("I am sorry to inform you that your extract job has failed.", email)
						TheQuery := "DELETE FROM extractjobs WHERE jobid=\"" + jobid + "\";"
						_, err := db.Exec(TheQuery)
						if err != nil {
							ExtractMailer(err.Error(), "nmwrap@edac.unm.edu")
						}
						logErr(err)
					} else if jobstat.JobStatus == "esriJobTimedOut" {
						ExtractMailer("I am sorry to inform you that your extract job has timed out", email)
						TheQuery := "DELETE FROM extractjobs WHERE jobid=\"" + jobid + "\";"
						_, err := db.Exec(TheQuery)
						if err != nil {
							ExtractMailer(err.Error(), "nmwrap@edac.unm.edu")
						}
						logErr(err)
					} else if jobstat.JobStatus == "esriJobCancelled" {
						ExtractMailer("I am sorry to inform you that your extract job has been cancled by an admin.", email)
						TheQuery := "DELETE FROM extractjobs WHERE jobid=\"" + jobid + "\";"
						_, err := db.Exec(TheQuery)
						if err != nil {
							ExtractMailer(err.Error(), "nmwrap@edac.unm.edu")
						}
						logErr(err)
					} else if jobstat.JobStatus == "esriJobDeleted" {
						ExtractMailer("I am sorry to inform you that your extract job has been deleted.", email)
						TheQuery := "DELETE FROM extractjobs WHERE jobid=\"" + jobid + "\";"
						_, err := db.Exec(TheQuery)
						if err != nil {
							ExtractMailer(err.Error(), "nmwrap@edac.unm.edu")
						}
						logErr(err)
					}

				}

			}
		}
	}
}

//ExtractFromUpload - Pre-process geom from shp file for ExtractGen
func ExtractFromUpload(w http.ResponseWriter, r *http.Request) {
	if IsLoggedIn(r) {
		AllowdShapeExtensions := []string{"cpg", "dbf", "prj", "sbn", "sbx", "shp", "shx"}

		file, _, err := r.FormFile("file")
		if err != nil {
			panic(err)
		}
		defer file.Close()
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
		ShapeName := ""
		defer reader.Close()
		dest := "/tmp/" + RandomFileName
		os.MkdirAll(dest, 755)
		for _, f := range reader.File {
			extension := strings.Split(f.Name, ".")
			path := dest + "/" + f.Name
			if stringInSlice(extension[1], AllowdShapeExtensions) {
				if extension[1] == "shp" {
					ShapeName = f.Name
				}

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
		Shapefile := "/tmp/" + RandomFileName + "/" + ShapeName
		myshape, err := shp.Open(Shapefile)
		if err != nil {
			log.Fatal(err)
		}
		defer myshape.Close()
		driver := gdal.OGRDriverByName("ESRI Shapefile")
		fmt.Println(driver)
		datasource, _ := driver.Open(Shapefile, 0)
		fmt.Println(datasource.LayerCount())
		layer := datasource.LayerByIndex(0)
		myfeature := layer.Feature(0)
		geom := myfeature.Geometry()
		spatialRef := gdal.CreateSpatialReference("PROJCS[\"WGS 84 / Pseudo-Mercator\",GEOGCS[\"WGS 84\",DATUM[\"WGS_1984\",SPHEROID[\"WGS 84\",6378137,298.257223563,AUTHORITY[\"EPSG\",\"7030\"]],AUTHORITY[\"EPSG\",\"6326\"]],PRIMEM[\"Greenwich\",0,AUTHORITY[\"EPSG\",\"8901\"]],UNIT[\"degree\",0.0174532925199433,AUTHORITY[\"EPSG\",\"9122\"]],AUTHORITY[\"EPSG\",\"4326\"]],PROJECTION[\"Mercator_1SP\"],PARAMETER[\"central_meridian\",0],PARAMETER[\"scale_factor\",1],PARAMETER[\"false_easting\",0],PARAMETER[\"false_northing\",0],UNIT[\"metre\",1,AUTHORITY[\"EPSG\",\"9001\"]],AXIS[\"X\",EAST],AXIS[\"Y\",NORTH],EXTENSION[\"PROJ4\",\"+proj=merc +a=6378137 +b=6378137 +lat_ts=0.0 +lon_0=0.0 +x_0=0.0 +y_0=0 +k=1.0 +units=m +nadgrids=@null +wktext  +no_defs\"],AUTHORITY[\"EPSG\",\"3857\"]]")
		geom.TransformTo(spatialRef)
		fmt.Println(geom.ToWKT())
		// fmt.Println(geom.ToGML())
		fmt.Println(geom.ToJSON())
		//fmt.Println(geom.ToKML())
		fmt.Println("ZZZ")
		fmt.Println(geom)
		var myGeoJSON GeoJSON
		json.Unmarshal([]byte(geom.ToJSON()), &myGeoJSON)
		var myGeom Geom
		myGeom.Title = r.FormValue("title")
		myGeom.Rings = myGeoJSON.Coordinates

		msg, err := ExtractGen(myGeom, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, msg)
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, msg)
		}

	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "You must be logged in!")
	}
}

//ExtractGen - creates an extract job and adds job id to DB
func ExtractGen(myGeom Geom, r *http.Request) (string, error) {
	myGeomMarshal, _ := json.Marshal(myGeom)
	var usergeom Geom
	if err := json.Unmarshal([]byte(myGeomMarshal), &usergeom); err != nil {

		return "Could not Unmarshal JSON Request.", errors.New("Could not Unmarshal JSON Request")
	}
	myRings, _ := json.Marshal(myGeom.Rings)
	fmt.Println(string(myRings))
	thisuser, err := UserData(r)
	aoi := "{\"features\":[{\"geometry\":{\"rings\":" + string(myRings) + "}}]}"
	queryurl := "https://edacarc.unm.edu/arcgis/rest/services/NMWRAP/ExtractData/GPServer/Extract%20Data/submitJob"
	resp, err := http.PostForm(queryurl, url.Values{
		"f":                {"pjson"},
		"Area_of_Interest": {aoi},
	})

	if err != nil {
		log.Println(err)
		return "Failed to post geom to service.", errors.New("Failed to post geom to service")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "Failed to read response body.", err
	}

	var jobstat JobStatus
	fmt.Println(string(body))
	json.Unmarshal(body, &jobstat)
	if jobstat.JobStatus == "esriJobSubmitted" {
		fmt.Println(jobstat.JobID)
		//Put db stuff here and if err not nil then return:
		db, err := sql.Open("mysql", dbuser+":"+dbpass+"@/"+dbname)

		if err != nil {

			return "DB err", err
		}
		defer db.Close()
		err = db.Ping()
		if err != nil {

			return "DB Ping test failed", err
		}

		TheQuery := "INSERT INTO extractjobs (email,jobid,status)  VALUES (\"" + thisuser["email"] + "\",\"" + jobstat.JobID + "\",\"1\");"
		//log.Println(TheQuery)
		_, err = db.Exec(TheQuery)
		if err != nil {
			log.Fatal(err)
			return "Failed to insert job into DB.", err
		}
		log.Println(myGeom.History)
		if myGeom.History != true {
			geomtitle := myGeom.Title
			marstring, _ := json.Marshal(myGeom)
			AreaQuery := "INSERT INTO areasofinterest (userid,geom,title)  VALUES ('" + thisuser["id"] + "','" + string(marstring) + "','" + geomtitle + "');"
			log.Println(AreaQuery)
			_, err = db.Exec(AreaQuery)
			if err != nil {
				log.Fatal(err)
			}
		}
		return "Extract task submitted. An e-mail will be sent to " + thisuser["email"] + " with job status updates.", nil

	}

	return "Job failed to submit.", errors.New("Job failed to submit")

}

//    _  _   _ _____ _  _ ___ _  _ _____ ___ ___   _ _____ ___ ___  _  _
//   /_\| | | |_   _| || | __| \| |_   _|_ _/ __| /_\_   _|_ _/ _ \| \| |
//  / _ \ |_| | | | | __ | _|| .` | | |  | | (__ / _ \| |  | | (_) | .` |
// /_/ \_\___/  |_| |_||_|___|_|\_| |_| |___\___/_/ \_\_| |___\___/|_|\_|

// MakeSalt makes a "random" salt for salting passwords.
func MakeSalt(n int) string {
	mathrand.Seed(time.Now().UnixNano())

	var saltrune = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = saltrune[mathrand.Intn(len(saltrune))]
	}
	return string(b)
}

//HashPass hashes passwords with a salt.
func HashPass(pw string, salt string) string {
	hash := sha512.New()
	hash.Write([]byte(pw + salt))
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)
	return mdStr

}

//GetCookieParts is used for parsing user supplied cookie claims.
func GetCookieParts(r *http.Request) UserClaims {
	var cookie, err = r.Cookie("nmwrap")
	logErr(err)
	cookieparts := strings.Split(cookie.Value, ".")
	//fmt.Println(cookieparts)
	data, _ := base64.StdEncoding.DecodeString(cookieparts[1] + "==")
	//fmt.Println(string(data))
	var user UserClaims
	json.Unmarshal([]byte(string(data)), &user)
	//fmt.Println(user.EMail)
	return user
}

//Login function
func Login(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("mysql", dbuser+":"+dbpass+"@/"+dbname)
	fmt.Println(reflect.TypeOf(db))
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

	//first we generate the proposed hash...
	submittedhashed := HashPass(user.Password, userssalt)
	if submittedhashed == hash {
		//password is good
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

	} else {

		w.WriteHeader(http.StatusUnauthorized)
	}

}

//CreateUser creates DB instance of user and sends pw reset to supplied email.
func CreateUser(w http.ResponseWriter, r *http.Request) {
	if IsLoggedIn(r) {
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
		randstring, _ := randomFilename()
		salt := MakeSalt(200)
		hashedpass := HashPass(randstring, salt)
		TheQuery := "INSERT INTO users (name,email,salt,hash,admin)  VALUES (\"" + user.Name + "\",\"" + user.Email + "\",\"" + salt + "\",\"" + hashedpass + "\",false);"
		_, err = db.Exec(TheQuery)
		if err != nil {
			log.Fatal(err)
		} else {
			ResetToken := RandString(200)
			TheQuery := "INSERT INTO passwordtokens (email,token)  VALUES (\"" + user.Email + "\",\"" + ResetToken + "\");"
			_, err := db.Exec(TheQuery)
			if err != nil {
				log.Fatal(err)
			} else {
				emailbodytmp1 := strings.Replace(NewAccountMail, "$name", user.Name, -1)
				emailbodytmp2 := strings.Replace(emailbodytmp1, "$email", user.Email, -1)
				emailbody := strings.Replace(emailbodytmp2, "$token", ResetToken, -1)
				mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
				subject := "Subject: NMWRAP Account!\n"
				msg := []byte(subject + mime + emailbody)
				recpt := []string{user.Email}
				err := smtp.SendMail("edacmail.unm.edu:25", nil, "nmwrap@edac.unm.edu", recpt, msg)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintln(w, "Account Creation Failed")
				} else {
					w.WriteHeader(http.StatusOK)
					fmt.Fprintln(w, "Account "+user.Email+"created.")
				}
			}
		}
	}
}

//Logout function
func Logout(w http.ResponseWriter, r *http.Request) {
	if IsLoggedIn(r) {
		user := GetCookieParts(r)
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

// ChangePassword checks if user supplied token exists. If so, then updates DB with a hash created from user supplied password.
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
		newhash := HashPass(passpak.Password, usersalt)
		_, err := db.Exec("UPDATE users SET hash =? WHERE email =?", newhash, usermail)
		if err != nil {
			log.Println("a")
			log.Println(err)
		} else {
			TheQuery := "DELETE FROM passwordtokens WHERE email=\"" + string(usermail) + "\";"
			_, err := db.Exec(TheQuery)
			if err != nil {
				log.Fatal(err)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		}

	}
}

//UserData is actual user data from the DB using the request cookie claims
func UserData(r *http.Request) (map[string]string, error) {
	//	fmt.Println("start UD")
	userinfo := map[string]string{}
	userclaims := GetCookieParts(r)
	db, err := sql.Open("mysql", dbuser+":"+dbpass+"@/"+dbname)
	if err != nil {

		log.Println(err)
		return userinfo, err
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {

		log.Println(err)
		return userinfo, err
	}

	rows, err := db.Query("SELECT * FROM users WHERE email='" + userclaims.EMail + "'")
	if err != nil {
		log.Println(err)
		return userinfo, err
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

			log.Println(err)
			return userinfo, err
		}
		//fmt.Println(name)

	}
	stradmin := strconv.FormatBool(admin)
	//fmt.Println(email)
	userinfo = map[string]string{
		"id":        strconv.Itoa(id),
		"name":      name,
		"email":     email,
		"userssalt": userssalt,
		"hash":      hash,
		"admin":     string(stradmin),
	}

	return userinfo, nil

}

//SendMail sends mail from nmwrap@edac.unm.edu using edacmail as smarthost. Some of this should be moved to config.
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
		db, err := sql.Open("mysql", dbuser+":"+dbpass+"@/"+dbname)
		logErr(err)
		var userid string
		var name string
		err = db.QueryRow("SELECT id, name FROM users WHERE email=?", string(resetemail)).Scan(&userid, &name)

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
			result1, err := db.Exec(TheQuery)
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Fprintln(w, result1)

				ResetToken := RandString(200)
				TheQuery := "INSERT INTO passwordtokens (email,token)  VALUES (\"" + string(resetemail) + "\",\"" + ResetToken + "\");"
				_, err := db.Exec(TheQuery)
				if err != nil {
					log.Fatal(err)
				} else {

					emailbodytmp1 := strings.Replace(PasswordReset, "$name", name, -1)
					emailbodytmp2 := strings.Replace(emailbodytmp1, "$email", string(resetemail), -1)
					emailbody := strings.Replace(emailbodytmp2, "$token", ResetToken, -1)

					mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
					subject := "Subject: NMWRAP Password Reset\n"
					msg := []byte(subject + mime + emailbody)
					recpt := []string{string(resetemail)}
					err := smtp.SendMail("edacmail.unm.edu:25", nil, "nmwrap@edac.unm.edu", recpt, msg)
					if err != nil {
						log.Println(err)
						w.WriteHeader(http.StatusInternalServerError)
						fmt.Fprintln(w, "Account Creation Failed")
					} else {
						w.WriteHeader(http.StatusOK)
						fmt.Fprintln(w, "Account "+string(resetemail)+"created.")
					}
				}
			}
		}

	} else {
		fmt.Fprintln(w, string("bad"))

	}

}

//IsLoggedIn - Checks if you are logged in.
func IsLoggedIn(r *http.Request) bool {
	var cookie, err = r.Cookie("nmwrap")
	if err != nil {
		log.Println(err)
	} else {

		token, err := jwt.ParseWithClaims(cookie.Value, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		logErr(err)
		if claims, ok := token.Claims.(*UserClaims); ok && token.Valid && cookie.Value == tokenmap[claims.EMail] {
			return true

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

//CheckReset -For front ends that want to check if token is valid before showing a user a PW reset form.
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
