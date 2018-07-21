package main

var (
	//% is repaced with the "from" or "to" depot name
	basicTransform = Transform{
		From: "//%/",
		To:   "//%/",
	}
	//transforms are indexed by the <db name>:<field index> where index enumerates from 0
	//fields are documented here: https://www.perforce.com/perforce/doc.current/schema/
	Transforms = map[string]Transform{
		"db.archmap:0": basicTransform, "db.archmap:1": basicTransform,
		"db.depot:0":  Transform{"%", "%"},
		"db.domain:0": Transform{"%", "%"},
	}
)

type Transform struct {
	From string
	To   string
}
