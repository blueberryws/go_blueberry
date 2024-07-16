package go_blueberry

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
)

func createPagesTable(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE page_data (
                        name TEXT PRIMARY KEY,
                        content TEXT NOT NULL
                )`)
	if err != nil {
		fmt.Println(err)
		log.Fatal("Unable to create page_data table.")
	}
}

func insertPageData[T any](name string, page T, db *sql.DB) (int64, error) {
	// serialize to JSON.
	pageJson, err := json.Marshal(page)
	result, err := db.Exec(
		`INSERT INTO page_data (name, content) VALUES (?,?);`, name, string(pageJson),
	)
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func updatePageData[T any](name string, page T, db *sql.DB) error {
	pageJson, err := json.Marshal(page)
	_, err = db.Exec(
		`UPDATE page_data set content=? where name=?;`, string(pageJson), name,
	)
	return err
}

func getPageData[T any](name string, db *sql.DB) (T, error) {
	rows, err := db.Query(
		`SELECT content FROM page_data WHERE name=?;`, name,
	)
	var pageData T
	if err != nil {
		return pageData, err
	}
	defer rows.Close()
	var content string
	for rows.Next() {
		err := rows.Scan(&content)
		if err != nil {
			return pageData, err
		}
		break
	}
	err = json.Unmarshal([]byte(content), &pageData)
	return pageData, err
}

func makePageDataGetter[T any](name string, db *sql.DB) func() T {
	return func() T {
		pageData, err := getPageData[T](name, db)
		if err != nil {
			log.Fatal(err)
		}
		return pageData
	}
}
