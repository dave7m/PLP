package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// global variables
var (
	fileName    string
	fullURLFile string
)

// Municipality struct for parsing json
type Municipality struct {
	MunicipalityNumber                                   int64   `json:"MunicipalityNumber"`
	MunicipalityName                                     string  `json:"MunicipalityName"`
	Canton                                               string  `json:"Canton"`
	Country                                              string  `json:"Country"`
	Scenario1_RoofsOnly_PotentialSolarElectricity_GWh    float64 `json:"Scenario1_RoofsOnly_PotentialSolarElectricity_GWh"`
	Scenario2_RoofsOnly_PotentialSolarElectricity_GWh    float64 `json:"Scenario2_RoofsOnly_PotentialSolarElectricity_GWh"`
	Scenario2_RoofsOnly_PotentialSolarHeat_GWh           float64 `json:"Scenario2_RoofsOnly_PotentialSolarHeat_GWh"`
	Scenario3_RoofsFacades_PotentialSolarElectricity_GWh float64 `json:"Scenario3_RoofsFacades_PotentialSolarElectricity_GWh"`
	Scenario4_RoofsFacades_PotentialSolarElectricity_GWh float64 `json:"Scenario4_RoofsFacades_PotentialSolarElectricity_GWh"`
	Scenario4_RoofsFacades_PotentialSolarHeat_GWh        float64 `json:"scenario_4___roofs_facades___potential_solar_heat___g_wh"`
	Factsheet                                            string  `json:"factsheet"`
	Methodology                                          string  `json:"methodology"`
}

func main() {
	downloadFile()
	parseFile()
}

// as in https://golangdocs.com/golang-download-files
func downloadFile() {
	fullURLFile = "http://www.uvek-gis.admin.ch/BFE/ogd/52/Solarenergiepotenziale_Gemeinden_Daecher_und_Fassaden.json"

	// Build fileName from fullPath
	fileURL, err := url.Parse(fullURLFile)
	if err != nil {
		log.Fatal(err)
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName = segments[len(segments)-1]

	// Create blank file
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	// Put content on file
	resp, err := client.Get(fullURLFile)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)
	_, err = io.Copy(file, resp.Body)

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
}

// as in https://tutorialedge.net/golang/parsing-json-with-golang/
func parseFile() {
	fileContent, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	defer func(fileContent *os.File) {
		err := fileContent.Close()
		if err != nil {
			panic(err)
		}
	}(fileContent)

	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(fileContent)

	// we initialize our Municipality array
	var municipalities []Municipality

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'municipalities' which we defined above
	err = json.Unmarshal(byteValue, &municipalities)
	if err != nil {
		println(err.Error())
		os.Exit(2) // this
	}

	// now that we have our municipality array, we can simply do our calculations
	var totalPotentialOfAllMunicipalities float64
	for i := 0; i < len(municipalities); i++ {
		totalPotentialOfAllMunicipalities += municipalities[i].Scenario3_RoofsFacades_PotentialSolarElectricity_GWh
	}
	fmt.Printf("Total Potential of all Municipalities is %f GWh\n", totalPotentialOfAllMunicipalities)
	tl := selectKLargest(municipalities, 3)
	// tl is the third-largest by Scenario3
	fmt.Printf("%s in %s has the 3rd largest Potential of %f\n",
		tl.MunicipalityName,
		tl.Canton,
		tl.Scenario3_RoofsFacades_PotentialSolarElectricity_GWh)
}

// partial selection sort until k-th element
func selectKLargest(municipalities []Municipality, k int) Municipality {
	for i := 0; i < k; i++ {
		maxIndex := i
		maxValue := municipalities[i].Scenario3_RoofsFacades_PotentialSolarElectricity_GWh
		for j := i + 1; j < len(municipalities); j++ {
			if municipalities[j].Scenario3_RoofsFacades_PotentialSolarElectricity_GWh > maxValue {
				maxIndex = j
				maxValue = municipalities[j].Scenario3_RoofsFacades_PotentialSolarElectricity_GWh
			}
		}
		municipalities[i], municipalities[maxIndex] = municipalities[maxIndex], municipalities[i]
	}
	return municipalities[k-1]
}
