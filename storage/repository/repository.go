package repository

import "github.com/charmbracelet/log"

var queries *Queries

func GetQueries() *Queries {
	if queries == nil {
		log.Fatal("Queries is nil, usage of uninitialized package repository")
	}
	return queries
}
func Setup(conn DBTX) {
	queries = New(conn)
}
