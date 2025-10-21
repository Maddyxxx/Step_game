package models

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"reflect"
	"strings"
)

// GetDataFromDB - Получение данных из бд
func GetDataFromDB(db *sqlx.DB, tableName string) []interface{} {
	tx, _ := db.Begin()
	defer tx.Commit()
	rows, err := tx.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
	if err != nil {
		log.Printf("GetDataFromDB error: %v, failed tx.Query, tableName %s ", err, tableName)
		return nil
	}

	var structList []interface{} // решение с динамическим типом очень спорное.
	for rows.Next() {
		structType, columns := getStructPoints(tableName)
		if err := rows.Scan(columns...); err != nil {
			log.Printf("Error scan columns error during GetDataFromDB: %v", err)
		}
		structList = append(structList, structType)
	}
	return structList
}

// Сохранение данных в бд
func InsertDataToDB(db *sqlx.DB, structType interface{}) error {
	tx, _ := db.Begin()
	defer tx.Commit()

	tableName, args := getStructFields(structType)
	query := fmt.Sprintf("INSERT INTO %s VALUES (%s?)", tableName, strings.Repeat("?,", len(args)-1))
	insert, _ := tx.Prepare(query)
	defer insert.Close()

	_, err := insert.Exec(args...)
	if err != nil {
		log.Printf("InsertDataToDB error %s, tableName: %v", err, tableName)
		return err
	}
	return nil
}

// jsonToMap - Convert JSON data to map[string]interface{}
func jsonToMap(jsonData []byte) map[string]interface{} {
	context := map[string]interface{}{}
	if err := json.Unmarshal(jsonData, &context); err != nil {
		log.Printf("GetStateFromDB error, json.Unmarshal, jsonData:  %v, %v", string(jsonData), err)
	}
	return context
}

// GetStateFromDB - Получение состояния пользователя из бд
func (u *UserState) GetData(db *sqlx.DB) error {
	defer log.Printf("GetStateFromDB %v", u.ChatID)
	tx, _ := db.Begin()
	defer tx.Commit()

	getState, err := tx.Prepare("SELECT * from userstate WHERE chat_id = ?")
	if err != nil {
		log.Printf("No state in db for %d", u.ChatID)
		return err
	}
	defer getState.Close()
	rows, _ := getState.Query(u.ChatID) //скипнул ошибку💩

	jsonData := []byte{}
	for rows.Next() {
		err = rows.Scan(&u.ChatID, &u.UserName, &u.ScenarioName, &u.StepName, &jsonData)
		if err != nil {
			continue
		}
	}
	if len(jsonData) != 0 {
		u.Context = jsonToMap(jsonData)
		return nil
	} else {
		log.Printf("No state in db for %d", u.ChatID)
		return fmt.Errorf("nil data")
	}
}

// getUserState - возвращает имя таблицы и поля структуры UserState
func getUserState(args []interface{}, u UserState) (string, []interface{}) {
	context, _ := json.Marshal(u.Context) //скипнул ошибку💩
	args = append(args, u.ChatID, u.UserName, u.ScenarioName, u.StepName, context)
	return "userstate", args
}

// getRequest - возвращает имя таблицы и поля структуры Request
func getRequest(args []interface{}, r Request) (string, []interface{}) {
	args = append(args, r.Date, r.UserName, r.Operation, r.Result)
	return "requests", args
}

// getStructFields - возвращает имя таблицы и поля структуры structType
func getStructFields(structType interface{}) (string, []interface{}) {
	var args []interface{}
	switch structType := structType.(type) {
	case UserState:
		return getUserState(args, structType)
	case Request:
		return getRequest(args, structType)

	default:
		log.Print("Error: undefined structType")
		return "", args
	}
}

// getStructPoints - возвращает указатели структуры structType
func getStructPoints(tableName string) (interface{}, []interface{}) {
	var structType interface{}
	switch tableName {
	case "userstate":
		structType = &UserState{}

	case "requests":
		structType = &Request{}

	}
	s := reflect.ValueOf(structType).Elem()
	numCols := s.NumField()
	points := make([]interface{}, numCols)
	for i := 0; i < numCols; i++ {
		field := s.Field(i)
		points[i] = field.Addr().Interface()
	}
	return structType, points
}

func (u UserState) InsertData(db *sqlx.DB) error {
	return InsertDataToDB(db, u)
}

func (r Request) InsertData(db *sqlx.DB) error {
	return InsertDataToDB(db, r)
}

func (u UserState) UpdateData(db *sqlx.DB) {
	tx, _ := db.Begin()
	defer tx.Commit()
	// кажется, лучше позже открыть db.Begin() и не задерживать  tx.Commit()

	updateData, _ := tx.Prepare("UPDATE userstate SET step_name = ?, context = ? WHERE chat_id = ?")
	defer updateData.Close()

	context, _ := json.Marshal(u.Context)
	_, err := updateData.Exec(u.StepName, context, u.ChatID)
	if err != nil {
		log.Printf("UpdateData for userstate error, failed updateData.Query, %v", err)
	}
}

// DeleteData - Удаление состояния пользователя из бд
func (u UserState) DeleteData(db *sqlx.DB) {
	tx, _ := db.Begin()
	defer tx.Commit()

	deleteState, _ := tx.Prepare("DELETE FROM userstate WHERE chat_id = ?")
	defer deleteState.Close()

	if _, err := deleteState.Exec(u.ChatID); err != nil {
		log.Printf("deleteState error, %v", err)
	}
}
