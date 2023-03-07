package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	srdb "github.com/surrealdb/surrealdb.go"
	"math/rand"
)

type person struct {
	Id   string
	Name string
	Age  int
}

type simulationOptions struct {
	batchSize  int
	iterations int
}

func generateLotsOfData(opts simulationOptions, db *srdb.DB) {
	for i := 0; i < opts.iterations; i++ {
		batch := make([]person, opts.batchSize, opts.batchSize)
		for b := 0; b < opts.batchSize; b++ {
			name := uuid.NewString()
			age := rand.Int()
			seqID := (opts.batchSize * i) + b
			batch[b] = person{Id: fmt.Sprintf("%d", seqID), Name: name, Age: age}
		}
		jsonArray, err := json.Marshal(batch)
		if err != nil {
			panic(err)
		}
		var result int
		query := fmt.Sprintf("INSERT INTO table %s TIMEOUT 5s", jsonArray)
		fmt.Printf("[%d/%d] query = %s\n", i+1, opts.iterations, query)
		_, err = db.Query(query, &result)
		if err != nil {
			panic(err)
		}
	}
}

func aSmallFunction(opts simulationOptions, db *srdb.DB) {
	var result person
	max := opts.batchSize * opts.iterations
	id := rand.Intn(max)
	fmt.Printf("Performing select for %d\n", id)
	_, err := db.Query(fmt.Sprintf("SELECT * FROM table:%d TIMEOUT 5s", id), &result)
	if err != nil {
		panic(err)
	}
}

func main() {
	db, err := srdb.New("ws://localhost:8000/rpc")
	if err != nil {
		panic(err)
	}
	fmt.Println("Before signup")

	//go func() {
	signin, err := db.Signin(map[string]interface{}{
		"user": "root",
		"pass": "root",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Finished signin %+v", signin)

	fmt.Println("Before use")
	_, err = db.Use("testns", "testdb")
	opts := simulationOptions{
		iterations: 1000,
		batchSize:  10,
	}
	go generateLotsOfData(opts, db)
	for true {
		aSmallFunction(opts, db)
	}
	fmt.Println("Program finished")
}
