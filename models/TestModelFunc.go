package models

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"testing"
)

func makeStates(db *sqlx.DB, n int) {
	for num := range n {
		state := UserState{
			ChatID:       int64(num),
			UserName:     fmt.Sprintf("%d", num*1234),
			ScenarioName: fmt.Sprintf("%d", num*58),
			StepName:     num,
			Context:      map[string]interface{}{},
		}
		state.InsertData(db)
	}
}

func TestUpdateStateFromDB(t *testing.T) {
	var db = InitDB("test.db")
	n := 10
	makeStates(db, n)

	var states, states2 []UserState

	for num := range n {
		var state UserState
		state.ChatID = int64(num)
		err := state.GetData(db)
		if err != nil {
			states = append(states, state)
		}
	}

	for _, state := range states {
		state.StepName *= 10
		state.UpdateData(db)
	}

	for num := range n {
		var state UserState
		state.ChatID = int64(num)
		err := state.GetData(db)
		if err != nil {
			states2 = append(states2, state)
		}

	}

	for num := range n {
		fmt.Println(num)
		if states[num].StepName == states2[num].StepName {
			t.Fatalf("StepName1 = %d, StepName2 %d , want diferent; GOT equal %d == %d",
				states[num].StepName, num*10, states[num].StepName, states2[num].StepName)
		}
	}
}

func makeReqs(db *sqlx.DB, n int) {
	for num := range n {
		req := Request{
			UserName:  fmt.Sprintf("%d", num*1234),
			Date:      fmt.Sprintf("%d", num*58),
			Operation: "testOperation",
			Result:    "testResult",
		}
		req.InsertData(db)
	}
}

func TestInsertDataToDB(t *testing.T) {
	var db = InitDB("test.db")
	n := 10

	for num := range n {
		req := Request{
			UserName:  fmt.Sprintf("%d", num*1234),
			Date:      fmt.Sprintf("%d", num*58),
			Operation: "testOperation",
			Result:    "testResult",
		}
		req.InsertData(db)
	}

	listData := GetDataFromDB(db, "requests")

	for _, data := range listData {
		req := data.(*Request)
		fmt.Println(data)
		fmt.Println(req.Operation, req.Date)
	}

	if len(listData) != n {
		t.Fatalf("Error, %v, len reqs = %d", listData, len(listData))
	}
}
