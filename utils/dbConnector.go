package utils

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func UpdateWordInNeo4jDb(word, wikiHtmlString string, camDicHtmlString string) {
	ctx := context.Background()
	dbUri := "neo4j+s://f09cd40d.databases.neo4j.io"
	dbUser := "neo4j"
	dbPassword := "C-0-IUZS_B8BO5SHw2P5Q9LYo2D3-ekmbLj04Ivfa0w"

	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))

	fmt.Println("Connection established.")

	result, err := neo4j.ExecuteQuery(ctx, driver,
		"MERGE (w:Word {word: $word}) "+
			"ON CREATE set w.createAt = timestamp(), w.searchCount = 0 "+
			"ON MATCH set w.lastUpdateAt = timestamp() "+
			//"set w.searchCount = w.searchCount + 1 " +
			"set w.sourceHtmlWiki = $sourceHtmlWiki "+
			"set w.sourceHtmlCamDic = $sourceHtmlCamDic "+
			"RETURN w.word AS word ",
		map[string]any{
			"word":             word,
			"sourceHtmlWiki":   wikiHtmlString,
			"sourceHtmlCamDic": camDicHtmlString,
		}, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))

	if err != nil {
		panic(err)
	}

	fmt.Printf("The query `%v` returned %v records in %+v.\n",
		result.Summary.Query().Text(), len(result.Records),
		result.Summary.ResultAvailableAfter())

	defer driver.Close(ctx)

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}
}

func GetWordInNeo4jDb(word string) map[string]any {
	ctx := context.Background()
	dbUri := "neo4j+s://f09cd40d.databases.neo4j.io"
	dbUser := "neo4j"
	dbPassword := "C-0-IUZS_B8BO5SHw2P5Q9LYo2D3-ekmbLj04Ivfa0w"

	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))

	fmt.Println("Connection established.")

	//result2, err := neo4j.ExecuteQuery(ctx, driver,
	//	"MATCH (p:Person {name: $name}) RETURN p",
	//	map[string]any{
	//		"name": "Alice",
	//	}, neo4j.EagerResultTransformer,
	//	neo4j.ExecuteQueryWithDatabase("neo4j"))
	//if err != nil {
	//	panic(err)
	//}
	//
	//summary := result2.Summary
	//fmt.Printf("Created %v nodes in %+v.\n",
	//	summary.Counters().NodesCreated(),
	//	summary.ResultAvailableAfter())

	result, _ := neo4j.ExecuteQuery(ctx, driver,
		"Match (w:Word {word: $word}) "+
			"RETURN w AS word ",
		map[string]any{
			"word": word,
		}, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))

	// Loop through results and do something with them
	//for _, record := range result.Records {
	//	//word, _ := record.Get("word")
	//	//fmt.Println(record.AsMap()["word"])
	//}

	// Summary information
	fmt.Printf("The query `%v` returned %v records in %+v.\n",
		result.Summary.Query().Text(), len(result.Records),
		result.Summary.ResultAvailableAfter())

	defer driver.Close(ctx)

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}

	if len(result.Records) == 0 {
		return nil
	}
	return result.Records[0].AsMap()
}
