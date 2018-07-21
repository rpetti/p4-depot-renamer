package main

import (
	"fmt"
)

var (
	//% is repaced with the "from" or "to" depot name
	pathTransform = Transform{
		From: "^//%s/",
		To:   "//%s/",
	}
	identTransform = Transform{
		From: "^%s$",
		To:   "%s",
	}
	//Transforms are indexed by the <db name>:<field index> where field index enumerates from 0
	//fields are documented here: https://www.perforce.com/perforce/doc.current/schema/
	Transforms = map[string]Transform{
		"db.archmap:0":    pathTransform,
		"db.archmap:1":    pathTransform,
		"db.depot:0":      identTransform,
		"db.depot:3":      Transform{From: "^%s/", To: "%s/"},
		"db.domain:0":     identTransform,
		"db.excl:0":       pathTransform,
		"db.graphperm:0":  identTransform,
		"db.have:1":       pathTransform,
		"db.have.pt:1":    pathTransform,
		"db.have.rp:1":    pathTransform,
		"db.haveg:2":      pathTransform,
		"db.haveview:4":   pathTransform,
		"db.integed:0":    pathTransform,
		"db.integed:1":    pathTransform,
		"db.integtx:0":    pathTransform,
		"db.integtx:1":    pathTransform,
		"db.label:1":      pathTransform,
		"db.locks:0":      pathTransform,
		"db.locksg:0":     pathTransform,
		"db.protect:6":    pathTransform,
		"db.resolve:0":    pathTransform,
		"db.resolve:1":    pathTransform,
		"db.resolve:8":    pathTransform,
		"db.resolveg:0":   pathTransform,
		"db.resolveg:1":   pathTransform,
		"db.resolvex:0":   pathTransform,
		"db.resolvex:1":   pathTransform,
		"db.resolvex:8":   pathTransform,
		"db.rev:0":        pathTransform,
		"db.rev:11":       pathTransform,
		"db.revbx:0":      pathTransform,
		"db.revbx:11":     pathTransform,
		"db.revcx:1":      pathTransform,
		"db.revdx:0":      pathTransform,
		"db.revdx:11":     pathTransform,
		"db.revhx:0":      pathTransform,
		"db.revhx:11":     pathTransform,
		"db.review:3":     pathTransform,
		"db.revpx:0":      pathTransform,
		"db.revpx:11":     pathTransform,
		"db.revsh:0":      pathTransform,
		"db.revsh:11":     pathTransform,
		"db.revsx:0":      pathTransform,
		"db.revsx:11":     pathTransform,
		"db.revtx:0":      pathTransform,
		"db.revtx:11":     pathTransform,
		"db.revux:0":      pathTransform,
		"db.revux:11":     pathTransform,
		"db.sendq:3":      pathTransform,
		"db.sendq:10":     pathTransform,
		"db.template:7":   pathTransform,
		"db.templatesx:8": pathTransform,
		"db.tamplatewx:8": pathTransform,
		"db.trigger:3":    pathTransform,
		"db.trigger:4":    pathTransform,
		"db.view:4":       pathTransform,
		"db.view.rp:4":    pathTransform,
		"db.working:1":    pathTransform,
		"db.working:17":   pathTransform,
		"db.workingg:1":   pathTransform,
		"db.workingg:17":  pathTransform,
		"db.workingx:0":   Transform{From: "^//([0-9]+)/%s/", To: "//$1/%s/"},
		"db.workingx:1":   pathTransform,
		"db.workingx:17":  pathTransform,
	}
)

//Transform defines a specific transformation
type Transform struct {
	From string
	To   string
}

func transformer(t <-chan JournalLine, out chan<- JournalLine, fromDepot string, toDepot string) {
	for input := range t {
		if input.EndOfFile {
			out <- input
			break
		}
		if input.Parsed {
			for idx, _ := range input.RowElems {
				//appy transforms to each element if they have a
				//transform defined for their table and column index
				key := fmt.Sprintf("%s:%d", input.Table, idx)
				if val, ok := Transforms[key]; ok {
					input.RowElems[idx].applyTransform(
						fmt.Sprintf(val.From, fromDepot),
						fmt.Sprintf(val.To, toDepot),
					)
				}
			}
		}
		out <- input
	}
}
