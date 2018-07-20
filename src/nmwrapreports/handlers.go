package main

import (

	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	//"reflect"
	"github.com/gorilla/mux"
	"github.com/jung-kurt/gofpdf"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"

)

// FireStations type is for marshalling the output from arcgis server
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

type Geom struct {
	Rings [][][]float64 `json:"rings"`
	Title string        `json:"title"`
}

type justGeom struct {
	Rings [][][]float64 `json:"rings"`
	Title string        `json:"title"`
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




// POSTGeom is how the user generates reports, by putting geom...
func POSTGeom(w http.ResponseWriter, r *http.Request) {

	jsbody, err := ioutil.ReadAll(r.Body)
	// fmt.Println(string(jsbody))
	if err != nil {
		log.Println(err)
	}
	var myGeom Geom
	json.Unmarshal(jsbody, &myGeom)

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
	CARLow:=0
	CARMed:=0
	CARHigh:=0
	fname, _ := randomFilename()
	for layernum := 0; layernum < 6; layernum++ {
		queryurl := "https://edacarc.unm.edu/arcgis/rest/services/NMWRAP/NMWRAP/MapServer/" + strconv.Itoa(layernum) + "/query?where=&text=&objectIds=&time=&geometryType=esriGeometryPolygon&inSR=&spatialRel=esriSpatialRelIntersects&relationParam=&outFields=*&returnGeometry=false&returnTrueCurves=false&maxAllowableOffset=&geometryPrecision=&outSR=&returnIdsOnly=false&returnCountOnly=false&orderByFields=&groupByFieldsForStatistics=&outStatistics=&returnZ=false&returnM=false&gdbVersion=&returnDistinctValues=false&resultOffset=&resultRecordCount=&f=pjson&geometry="
		queryurl = queryurl + string(jsbody)
		queryurl=strings.Replace(queryurl, " ", "%20", -1)
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
			pdf.CellFormat(190, 15,"Fire stations" , "0", 1, "CM", false, 0, "")
			pdf.SetFont("Helvetica", "", 11)
			var myJSON FireStations
			json.Unmarshal(body, &myJSON)
			//fmt.Println("asaa")
            //fmt.Println(len(myJSON.Features))
			FireStationsBlurb := "There are "+strconv.Itoa(len(myJSON.Features))+" fire stations in this area. The proximity of fire stations is essential to an assessment of fire safety. Lorem ipsum dolor sit amet, consectetur adipiscing elit.\nInteger nec odio. Praesent libero. Sed cursus ante dapibus diam. Sed nisi. Nulla quis sem at nibh elementum imperdiet. Duis sagittis ipsum. Praesent mauris. Fusce nec tellus sed augue semper porta. Mauris massa. Vestibulum lacinia arcu eget nulla. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Curabitur sodales ligula in libero."
			lines := pdf.SplitLines([]byte(FireStationsBlurb), 190.0)
			_, lineHt := pdf.GetFontSize()

			for _, line := range lines {
				pdf.CellFormat(190.0, lineHt, string(line), "", 1, "TL", false, 0, "")
			}
			if len(myJSON.Features)>0{
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
			pdf.CellFormat(190, 15,"Communities At Risk" , "0", 1, "CM", false, 0, "")

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
			if len(myJSON.Features)>0{
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
					CARLow+= 1
				} else if element.Attributes.Rate2016 == "M" {
					rate = "Medium Risk"
					CARMed+= 1
				} else if element.Attributes.Rate2016 == "H" {
					rate = "High Risk"
					CARHigh+= 1
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
			carsummation:=0
			if CARLow>0{
				carsummation=carsummation+1
			}
			if CARMed>0{
				carsummation=carsummation+1
			}
			if CARHigh>0{
				carsummation=carsummation+1
			}
			if carsummation>0{
			pie := chart.PieChart{
				Title: "Communities At Risk",
				// ColorPalette : colorpal ,
				Width:  512,
				Height: 512,
				Values: []chart.Value{
					{Value: float64(CARHigh), Label: strconv.Itoa(CARHigh)+" High Risk"},
					{Value: float64(CARMed), Label: strconv.Itoa(CARMed)+" Medium Risk"},
					{Value: float64(CARLow), Label: strconv.Itoa(CARLow)+" Low Risk"},
					

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
			err = pie.Render(chart.PNG,buffer)
			if err != nil {
				fmt.Printf("Error rendering pie chart: %v\n", err)
			}
			piereader := bufio.NewReader(buffer)

			var options gofpdf.ImageOptions
			options.ImageType="PNG"
			fmt.Println(options)
			fmt.Println(pdf.GetPageSize())
			pdf.RegisterImageOptionsReader("piechart",options,piereader)
			


			whatwegot:=297.0
			whatweneed:=pdf.GetY()+128
			fmt.Println(whatwegot)
			fmt.Println(whatweneed)
			if whatweneed > whatwegot{
				extrapadding:=297.0-pdf.GetY()
				
				pdf.CellFormat(10, extrapadding, "", "0", 0, "", false, 0, "")
				fmt.Println("lol")
				
			}
			CurrentX:=pdf.GetX()
			CurrentY:=pdf.GetY()
			if pdf.Ok() {
				
				pdf.Image("piechart", CurrentX+31, CurrentY, 128, 128, false, "", 0, "")

					pdf.SetY(CurrentY+128)
			}
			fmt.Println(pdf.GetPageSize())
			fmt.Println(pdf.GetMargins())
			fmt.Println(CurrentY)
		}

		case 2:
			//IncorporatedCityBoundaries
			pdf.Ln(15)
			pdf.SetFont("Helvetica", "", 20)
			pdf.CellFormat(190, 15,"Incorporated City Boundaries" , "0", 1, "CM", false, 0, "")

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
			if len(myJSON.Features)>0{
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
			pdf.CellFormat(190, 15,"Vegetation Treatments" , "0", 1, "CM", false, 0, "")

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
			if len(myJSON.Features)>0{
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
			pdf.CellFormat(190, 15,"Watersheds HUC8" , "0", 1, "CM", false, 0, "")

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
			if len(myJSON.Features)>0{
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
			pdf.CellFormat(190, 15,"Counties" , "0", 1, "CM", false, 0, "")

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
			if len(myJSON.Features)>0{
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