package main

import (
	"database/sql"
	"strings"

	sqlx "github.com/jmoiron/sqlx"
)

func getDB() (*sql.DB, error) {
	const file string = "database.sqlite3"
	return sql.Open("sqlite", file)
}

func getWorks() ([]Work, error) {
	db, e := getDB()

	if e != nil {
		return []Work{}, e
	}

	rows, err := db.Query("select WorkID, Title from Works")

	if err != nil {
		return []Work{}, err
	}

	defer rows.Close()

	var results []Work

	for rows.Next() {
		var t string
		var id string
		err := rows.Scan(&id, &t)

		if err != nil {
			return results, err
		}

		w := Work{WorkID: id, Title: t}

		results = append(results, w)
	}

	return results, nil

}

func getChars() ([]Character, error) {
	db, e := getDB()

	if e != nil {
		return []Character{}, e
	}

	rows, err := db.Query("select CharID, CharName from Characters")

	if err != nil {
		return []Character{}, e
	}

	defer rows.Close()

	var results []Character

	for rows.Next() {
		var charId string
		var charName string

		err := rows.Scan(&charId, &charName)

		if err != nil {
			return results, err
		}

		c := Character{Name: charName, CharID: charId}

		results = append(results, c)
	}

	return results, nil
}

func prepareQueryAndArgs(query SearchQuery) (string, []interface{}, error) {
	rawQuery := "select w.Title,highlight(par_fts, 0, '<mark>', '</mark>') as text, c.CharName from par_fts s join Paragraphs p on s.rowid = p.ParagraphID join Works w on w.WorkID = p.WorkID join Characters c on c.CharID = p.CharID where s.text MATCH ? "

	var args []any
	args = append(args, query.QueryText)

	if len(query.WorkIds) > 0 || len(query.CharIds) > 0 {
		if len(query.WorkIds) > 0 {
			rawQuery = rawQuery + "and p.WorkID in (?)"
			args = append(args, query.WorkIds)
		}

		if len(query.CharIds) > 0 {
			rawQuery = rawQuery + "and c.CharID in (?)"
			args = append(args, query.CharIds)
		}

		sqlQuery, args, err := sqlx.In(rawQuery, args...)
		if err != nil {
			return "", nil, err
		}

		return sqlQuery, args, nil

	} else {
		return rawQuery, nil, nil
	}
}

func executeFTS(query SearchQuery) ([]SearchResult, error) {
	var results []SearchResult

	db, e := getDB()

	if e != nil {
		return results, e
	}

	sqlQuery, args, err := prepareQueryAndArgs(query)

	if err != nil {
		return results, err
	}

	var rows *sql.Rows

	if len(query.WorkIds) > 0 || len(query.CharIds) > 0 {
		rows, err = db.Query(sqlQuery, args...)
	} else {
		rows, err = db.Query(sqlQuery, query.QueryText)
	}

	if err != nil {
		print(err)
		return results, err
	}

	defer rows.Close()

	for rows.Next() {
		var text string
		var title string
		var charName string

		err := rows.Scan(&title, &text, &charName)

		if err != nil {
			return results, err
		}

		r := SearchResult{Text: strings.ReplaceAll(text, "[p]", ""), Work: title, Character: charName}
		// print(r)
		// fmt.Println(valast.String(r))
		results = append(results, r)
	}

	return results, nil
}
