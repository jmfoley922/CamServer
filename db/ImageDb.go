package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

//Error logger
func LogError(errText string, connString string) error {

	conn, err := sql.Open("mssql", connString)
	if err != nil {
		fmt.Printf("LogAtmSeq error: %s\n", err.Error())
		return err
	}

	defer conn.Close()
	var sql1 = "insert into db_cameraImages_errorlog (data,date())values(?)"
	_, err = conn.Exec(sql1, errText)

	if err != nil {
		fmt.Printf("LogError error: %s\n", err.Error())

		return err
	}

	return err
}

//Write out red light violation image
func WriteCameraImage(cameraId string, imageData []byte, connString string) error {

	SqlDb, err := sql.Open("mssql", connString)
	if err != nil {
		log.Printf("Error opening db: %s\n", err.Error())
		return err
	}

	defer SqlDb.Close()
	var sql1 = "insert into db_cameraImages(cameraId, imageData) values (" +
		"?, ?)"

	_, err = SqlDb.Exec(sql1, cameraId, imageData)
	if err != nil {
		log.Printf("Error adding atm image: %s\n", err.Error())
		LogError(err.Error(), connString)
		return err
	}

	return err
}
