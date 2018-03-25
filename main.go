package main

import (
	"CamServer/db"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Settings struct {
	SqlServer string `json:"sqlserver"`
	CertName  string `json:"certname"`
	KeyName   string `json:"keyname"`
	Port      string `json:"port"`
	DbUser    string `json:"dbuser"`
	DbPword   string `json:"dbpword"`
	Db        string `json:"db"`
}

var appSettings = Settings{}

//Read settings file and populate structure
func getSettings() (err error) {

	settingsStr, err := ioutil.ReadFile("./settings.json")

	err = json.Unmarshal(settingsStr, &appSettings)

	if err != nil {
		log.Printf("Error = %s\n", err)
		return err
	}

	return err

}

//Build the Sql Server connection string from the settings file
func getConnectionString() string {

	var connString = fmt.Sprintf("server=%s;user id=%s;password=%s;encrypt=%s;database=%s",
		appSettings.SqlServer, appSettings.DbUser, appSettings.DbPword,
		"disable", appSettings.Db)

	return connString
}

//Process the new traffic cam image
func addCameraImage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	cameraId := p.ByName("CamId")

	//read image data from body of request
	image, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: Error = %s\n", err.Error())
		io.WriteString(w, err.Error())

	} else {
		err := db.WriteCameraImage(cameraId, image, getConnectionString())
		if err != nil {
			io.WriteString(w, err.Error())
		} else {
			io.WriteString(w, "OK")
		}
	}

}

func main() {

	err := getSettings()
	if err != nil {
		log.Fatal("Error reading settings file: " + err.Error())
	}

	r := httprouter.New()

	r.POST("/AddPhoto/:CamId/", addCameraImage)

	log.Fatal(http.ListenAndServeTLS(appSettings.Port,
		appSettings.CertName, appSettings.KeyName, r))

}
