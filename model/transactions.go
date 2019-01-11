package model

import (
	"fmt"

	"github.com/gocql/gocql"
)

func (d *Data) Get(session *gocql.Session, nameSpace string) error {

	// Query for data
	query := fmt.Sprintf("SELECT data FROM %s where name=?", nameSpace)
	iterable := session.Query(query, d.Name).Iter()

	// If no records found, keep Found field as false and return
	if iterable.NumRows() > 0 {
		d.Found = true
	} else {
		return nil
	}

	// If records found, add to Data field
	m := map[string]interface{}{}
	for iterable.MapScan(m) {
		d.Data = append(d.Data, m["data"])
		// clear map as required by gocql package
		m = map[string]interface{}{}
	}
	return nil
}
