package main

import (
	"context"
	"fmt"
	"prerequisite/multiple"
	"prerequisite/repository"
	"prerequisite/single"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	// dfs.Main()
	// fmt.Println()
	// kahn.Main()
	// fmt.Println()
	// all.Main()

	driver, err := neo4j.NewDriverWithContext("bolt://localhost:7687",
		neo4j.BasicAuth("neo4j", "password", ""))
	if err != nil {
		panic(err.Error())
	}
	err = driver.VerifyConnectivity(context.Background())
	if err != nil {
		panic(err.Error())
	}

	repo := repository.NewRepo(driver)

	edgelist, err := repo.GetEdgeList()
	if err != nil {
		panic(err.Error())
	}

	lookUp, err := repo.GetLookUp()
	if err != nil {
		panic(err.Error())
	}

	inDegree, err := repo.GetInDegrees()
	if err != nil {
		panic(err.Error())
	}

	inDegree2, err := repo.GetInDegrees()
	if err != nil {
		panic(err.Error())
	}

	// svc := single.NewRandomPathGenerate(edgelist, lookUp, inDegree)
	svc := single.NewMaximalPathGenerate(edgelist, lookUp, inDegree)
	result, err := svc.GetPath()

	if err != nil {
		panic(err.Error())
	}

	for _, v := range result {
		fmt.Printf("%3s (%0.2f)", v.Id, v.Weight)

		fmt.Print("| ")
	}

	fmt.Printf("%0.2f", result.GetDCG())
	fmt.Print("| ")

	fmt.Printf("%0.2f", result.GetIDCG())
	fmt.Print("| ")

	fmt.Printf("%0.2f", result.GetNDCG(result.GetDCG(), result.GetIDCG()))
	fmt.Print("| ")

	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()

	// --------------------------------------------------
	// return

	svc2 := multiple.NewAllGenerate(edgelist, lookUp, inDegree2)
	result2, err := svc2.GetPath()

	if err != nil {
		panic(err.Error())
	}

	for _, row := range result2 {
		for _, v := range row {
			fmt.Printf("%3s (%0.2f)", v.Id, v.Weight)

			fmt.Print("| ")
		}

		fmt.Printf("%0.2f", row.GetDCG())
		fmt.Print("| ")

		fmt.Printf("%0.2f", result.GetDCG())
		fmt.Print("| ")

		fmt.Printf("%0.2f", row.GetNDCG(row.GetDCG(), repository.ScoreIDCG(result.GetDCG())))
		fmt.Print("| ")

		fmt.Println()
	}
	fmt.Println()
}
