package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/manticoresoftware/go-sdk/manticore"

	"mc_fts/src/server/analytics"
	"mc_fts/src/server/sqlstore"
	"mc_fts/src/server/templates"
)

var db *sql.DB

type homeItem struct {
	IncidentId string
	Fields     string
	Snippet    template.HTML
}

type incidentFull struct {
	IncidentId      string
	RuleId          string
	RuleDisplayName string
	RuleDescription string
	HostId          *string
	Host            *string
	LinkId          *string
	Link            *string
	Severity        string
	Status          string
	Snippet         template.HTML
}

type IncidentFilter struct {
	Severity string
	Status   string
}

func MCQuery(cl *manticore.Client, index string, query string, offset int32, filter IncidentFilter) ([]string, *time.Duration, error) {
	if filter.Severity != "" {
		query += fmt.Sprintf(" INCIDENT_SEVERITY_%v", strings.ToUpper(filter.Severity))
	}
	if filter.Status != "" {
		query += fmt.Sprintf(" INCIDENT_STATUS_%v", strings.ToUpper(filter.Status))
	}

	mcQuery := strings.Join(strings.Split(query, " "), "|")
	search := manticore.NewSearch(mcQuery, index, "")
	search.Limit = 100
	search.MaxMatches = 50000
	search.Offset = offset

	qres, err := cl.RunQuery(search)
	if err != nil {
		return nil, nil, err
	}

	results := make([]string, 0)

	for _, match := range qres.Matches {
		results = append(results, match.Attrs[0].(manticore.JsonOrStr).String())
	}

	return results, &qres.QueryTime, nil
}

func SearchIncidents(cl *manticore.Client, query string, filter IncidentFilter) ([]incidentFull, *time.Duration, error) {
	start := time.Now()

	incidentIDs, _, err := MCQuery(cl, "incidents", query, 0, filter)
	if err != nil {
		return []incidentFull{}, nil, err
	}

	incidentValues := make([]string, 0)
	for i, incidentID := range incidentIDs {
		incidentValues = append(incidentValues, fmt.Sprintf("(%v, '%v')", i, incidentID))
	}

	if len(incidentValues) == 0 {
		return []incidentFull{}, nil, nil
	}
	incidentValuesStr := strings.Join(incidentValues, ", \n")

	queryRaw := fmt.Sprintf(sqlstore.IncidentsList, incidentValuesStr)

	rows, err := db.Query(queryRaw)
	if err != nil {
		return []incidentFull{}, nil, err
	}
	defer rows.Close()
	results := make([]incidentFull, 0)
	for rows.Next() {
		var r incidentFull
		var snip string
		if err := rows.Scan(&r.IncidentId, &r.RuleId, &r.RuleDisplayName, &r.RuleDescription, &r.LinkId, &r.Link, &r.HostId, &r.Host, &r.Severity, &r.Status); err != nil {
			return nil, nil, err
		}
		r.Snippet = template.HTML(strings.Replace(snip, "\n", "<br>", -1))
		results = append(results, r)
	}

	qt := time.Since(start)

	return results, &qt, nil
}

func main() {
	words := []string{"layer", "opposite", "waist", "become", "address", "adult", "upper", "twelve", "card", "prefer", "patient", "concerning", "welcome", "bread", "connect", "beyond", "law", "northern", "more", "gray", "west", "except", "OK", "negative", "nation", "program", "plenty", "wine", "information", "produce", "animal", "smart", "fear", "lock", "upper", "physical", "beautiful", "truck", "steady", "card", "walk", "rock", "bear", "grass", "hand", "odd", "proof", "decrease", "represent", "over", "quiet", "solve", "require", "important", "inform", "nose", "very", "crowd", "third", "request", "woman", "practical", "invite", "adjective", "wake", "soon", "itself", "relation", "fork", "food", "average", "change", "well", "each", "quality", "supply", "point", "dollar", "child", "pound", "balance", "suddenly", "cook", "notice", "traffic", "recognize", "drunk", "toilet", "always", "say", "reason", "under", "forget", "replace", "medical", "clothes", "breast", "straight", "duck", "admit"}

	dburl := "postgresql://postgres@127.0.0.1:5432/nextdb?sslmode=disable"

	var err error
	if db, err = sql.Open("postgres", dburl); err != nil {
		log.Fatal(err)
	}

	cl := manticore.NewClient()
	cl.SetServer("127.0.0.1", 9313)
	_, err = cl.Open()
	if err != nil {
		panic(err)
	}

	//rows, err := db.Query(sqlstore.IncidentsFieldsListFull)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//defer rows.Close()
	//results := make([]homeItem, 0, 10)
	//for rows.Next() {
	//	var h homeItem
	//	if err := rows.Scan(&h.IncidentId, &h.Fields); err != nil {
	//		return
	//	}
	//
	//	results = append(results, h)
	//}
	//
	//for i, h := range results {
	//	qStr := fmt.Sprintf(
	//		`insert into incidents values(%v, '%v', '%v')`,
	//		i,
	//		h.IncidentId,
	//		h.Fields,
	//	)
	//	_, err = cl.Sphinxql(qStr)
	//	if err != nil {
	//		fmt.Println(err.Error())
	//		break
	//	}
	//	var percent = (float32(i+1) / float32(len(results))) * 100
	//	fmt.Printf("%v/100\n", percent)
	//}

	tplHome := template.Must(template.New(".").Parse(templates.TplStrHome))
	tplResults := template.Must(template.New(".").Parse(templates.TplStrResults))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		q := r.FormValue("q")
		sev := r.FormValue("severity")
		st := r.FormValue("status")
		filter := IncidentFilter{
			Severity: sev,
			Status:   st,
		}
		if q == "" && sev == "" && st == "" {
			rows, err := db.Query(sqlstore.IncidentsFieldsListPart)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			defer rows.Close()
			results := make([]homeItem, 0, 10)
			for rows.Next() {
				var r homeItem
				var snip string
				if err := rows.Scan(&r.IncidentId, &r.Fields); err != nil {
					http.Error(w, err.Error(), 404)
					return
				}
				r.Snippet = template.HTML(strings.Replace(snip, "\n", "<br>", -1))
				results = append(results, r)
			}

			tplHome.Execute(w, map[string]interface{}{
				"Results": results,
			})
			return
		}
		if q == "calculate" {
			searchTimes := make([]float32, 0)
			for i := 0; i < 1; i++ {
				fmt.Println(i+1, "/", 10)
				for _, word := range words {
					_, qtime, err := SearchIncidents(&cl, word, filter)
					if err != nil {
						http.Error(w, err.Error(), 500)
						return
					}
					if qtime != nil && qtime.Milliseconds() > 0 {
						searchTimes = append(searchTimes, float32(qtime.Milliseconds()))
					}
				}
			}

			if len(searchTimes) == 0 {
				panic("Empty")
			}

			a := analytics.Anal{}
			a.Show(searchTimes)

			return
		}

		results, qtime, err := SearchIncidents(&cl, q, filter)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
		if qtime != nil {
			fmt.Println(qtime.Milliseconds())
		}
		tplResults.Execute(w, map[string]interface{}{
			"Results": results,
			"Query":   q,
		})
	})

	const PORT = "4242"

	fmt.Println("Starting on", PORT, "PORT...")
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
