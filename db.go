package goscript

import "github.com/hashicorp/go-memdb"

type State struct {
	DomainEntity string
	Domain       string
	Entity       string
	State        string
	Attributes   map[string]interface{}
}

func memdbSchema() *memdb.DBSchema {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"state": {
				Name: "state",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "DomainEntity"},
					},
					"domain": {
						Name:    "domain",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Domain"},
					},
					"entity": {
						Name:    "entity",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Entity"},
					},

					"domain_entity": {
						Name:   "domain_entity",
						Unique: true,
						Indexer: &memdb.CompoundIndex{Indexes: []memdb.Indexer{
							&memdb.StringFieldIndex{Field: "Domain"},
							&memdb.StringFieldIndex{Field: "Entity"},
						}},
					},
				},
			},
		},
	}
	return schema
}
