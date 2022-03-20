package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/exp/constraints"
	"golang.org/x/text/encoding/charmap"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	/*
		a := []int{6, 3, -6, -45, 7, 3, 9, 4, 35, 6, 8, 4, 24, 8, 5}
		fmt.Println("sorting ", a)
		quicksort(a)
		fmt.Println(a)

		y := []string{"Foo", "foo", "foobar", "FizzBuzz", "abc", "cbde", "aaaa", "zzz", "fses"}
		fmt.Println("sorting ", y)
		quicksort(y)
		fmt.Println(y)

		z := []float64{10.45, 3.141, -49, 25.24, 924.1, 4.5, 6.2, 9.5, -3.5}
		fmt.Println("sorting ", z)
		quicksort(z)
		fmt.Println(z)

		x := []uint{6, 3, 7, 3, 9, 4, 35, 6, 8, 4, 24, 8, 5}
		fmt.Println("sorting ", x)
		quicksort(x)
		fmt.Println(x)

		callExternalCommand()
	*/
	fileName := downloadFile()
	parseFile(fileName)
}

func parseFile(filename string) {
	fileContent, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer fileContent.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(fileContent)

	// we initialize our Municipality array
	var municipalities []Municipality

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	err = json.Unmarshal(byteValue, &municipalities)
	if err != nil {
		println(err.Error())
		os.Exit(2) // this
	}

	var totalPotentialOfAllMunicipalities float64
	for i := 0; i < len(municipalities); i++ {
		totalPotentialOfAllMunicipalities += municipalities[i].Scenario3_RoofsFacades_PotentialSolarElectricity_GWh
	}
	fmt.Printf("Total Potential of all Municipalities is: %f GWh\n", totalPotentialOfAllMunicipalities)
	tl := selectKLargest(municipalities, 3) // not correct yet
	fmt.Printf("%s in %s has the 3rd largest potential of %f\n",
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

var (
	fileName    string
	fullURLFile string
)

type Municipalities struct {
	Municipalities []Municipality
}

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

// https://golangdocs.com/golang-download-files
// https://tutorialedge.net/golang/parsing-json-with-golang/
func downloadFile() string {
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
	return fileName
}

func callExternalCommand() {
	os := runtime.GOOS
	var out []byte
	var err error
	switch os {
	case "windows":
		cmd := exec.Command("systeminfo")
		out, err = cmd.Output()
		if err != nil {
			fmt.Println("Error: ", err)
		}
		// not tested, might not work.
	case "linux", "darwin":
		cmd := exec.Command("uname -a")
		out, err = cmd.Output()
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	// because the returned byte array is not encoded in the right format, the following 5 lines take care of it
	// compare to https://www.reddit.com/r/golang/comments/9zsipj/help_osexec_output_on_nonenglish_windows_cmd/
	d := charmap.CodePage850.NewDecoder()
	output, err := d.Bytes(out)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(output))
}

// generics are only supported since Go 1.18, so make sure to have downloaded the recent version
// function arguments as in https://golangexample.com/generic-sort-for-slices-in-golang/
// constraints.ordered contains types that are comparable with <, <=, ==, !=, >= and >.
// see https://pkg.go.dev/golang.org/x/exp/constraints
func quicksort[E constraints.Ordered](list []E) {
	l := 0
	r := len(list) - 1
	quicksortHelper(list, l, r)
}

// this quicksort algorithm uses dual pivot, meaning in the partitioning we create two pivots instead of one
// the algorithm is semi efficient, because we do not gain very much from dual pivot, and the pivots are not chosen randomly, so the worst case does not improve (O(n^2)).
func quicksortHelper[E constraints.Ordered](list []E, l int, r int) {
	if r-l <= 0 {
		return
	}
	if list[l] > list[r] {
		list[l], list[r] = list[r], list[l]
	}
	p, q := partition(list, l, r)
	quicksortHelper(list, l, p-1)
	quicksortHelper(list, p+1, q-1)
	quicksortHelper(list, q+1, r)
}

func partition[E constraints.Ordered](list []E, lo int, hi int) (int, int) {
	l := lo + 1
	m := lo + 1
	g := hi
	for m < g {
		if list[m] < list[lo] {
			list[l], list[m] = list[m], list[l]
			l++
			m++
		} else if list[m] >= list[hi] {
			g--
			list[m], list[g] = list[g], list[m]
		} else {
			m++
		}
	}
	l--
	list[lo], list[l] = list[l], list[lo]
	list[hi], list[m] = list[m], list[hi]
	return l, m
}
