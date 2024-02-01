package main

import (
	"encoding/xml"
	"github.com/go-pg/pg/v10"
	"io"
	"net/http"
	"strconv"
	"strings"
	"treasury/src/domain"
	"treasury/src/driver/postgres"
	"treasury/src/entity"
	"treasury/src/server/response"
)

const (
	errCode    int    = 503
	okCode     int    = 200
	personType string = "Individual"
	sdnUrl     string = "https://www.treasury.gov/ofac/downloads/sdn.xml"
)

var Db *pg.DB

func main() {
	initDb()
	initServer()
}

func initServer() {
	http.HandleFunc("/update", updateHandler)
	http.HandleFunc("/state", stateHandler)
	http.HandleFunc("/count", countHandler)
	http.HandleFunc("/get_names", getNamesHandler)
	http.ListenAndServe(":8080", nil)
}

func initDb() {
	var err error
	Db, err = postgres.StartDB()
	if err != nil {
		panic(err)
	}
}

func getNamesHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		resp := response.EmptySearchResponse{
			Result: false,
			Info:   "Empty name",
		}
		response.JsonResponse(w, resp, errCode)
		return
	}
	typeSearch := r.URL.Query().Get("type")
	if typeSearch == "" {
		//resp := response.EmptySearchResponse{
		//	Result: false,
		//	Info:   "Empty type",
		//}
		//response.JsonResponse(w, resp, errCode)
		//return
		typeSearch = "strong"
	}

	typeSearch = strings.ToLower(typeSearch)
	if typeSearch == "strong" {
		person, err := domain.GetPersonStrong(Db, name)
		if err != nil {
			resp := response.EmptySearchResponse{
				Result: false,
				Info:   err.Error(),
			}
			response.JsonResponse(w, resp, errCode)
			return
		}
		response.JsonResponse(w, person, errCode)
	}
	if typeSearch == "weak" {
		persons, err := domain.GetPersonWeak(Db, name)
		if err != nil {
			resp := response.EmptySearchResponse{
				Result: false,
				Info:   err.Error(),
			}
			response.JsonResponse(w, resp, errCode)
			return
		}
		response.JsonResponse(w, persons, errCode)
	}

}

func countHandler(w http.ResponseWriter, r *http.Request) {
	resp := response.StateResponse{}
	count, err := domain.GetCount(Db)
	if err != nil {
		resp = response.StateResponse{
			Result: false,
			Info:   err.Error(),
		}
		response.JsonResponse(w, resp, errCode)
	} else {
		resp = response.StateResponse{
			Result: true,
			Info:   strconv.Itoa(count),
		}
		response.JsonResponse(w, resp, okCode)
	}
}

func stateHandler(w http.ResponseWriter, r *http.Request) {
	resp := response.StateResponse{}
	status, err := domain.GetState(Db)
	if err != nil {
		resp = response.StateResponse{
			Result: false,
			Info:   err.Error(),
		}
		response.JsonResponse(w, resp, errCode)
	} else {
		resp = response.StateResponse{
			Result: true,
			Info:   status.Info,
		}
		response.JsonResponse(w, resp, okCode)
	}
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	resp := response.UpdateResponse{}
	go update()

	resp = response.UpdateResponse{
		Code:   okCode,
		Result: true,
		Info:   "",
	}
	response.JsonResponse(w, resp, okCode)
}

func update() {
	domain.SetState(Db, &domain.State{Info: "updating"})

	resp, err := http.Get(sdnUrl)

	if err != nil {
		domain.SetState(Db, &domain.State{Info: err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		domain.SetState(Db, &domain.State{Info: "StatusCode error: " + strconv.Itoa(resp.StatusCode)})
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil && err != io.EOF {
		domain.SetState(Db, &domain.State{Info: err.Error()})
		return
	}
	sdn := new(entity.Sdn)
	err = xml.Unmarshal(body, sdn)
	if err != nil {
		domain.SetState(Db, &domain.State{Info: err.Error()})
		return
	}

	persons, err := domain.GetAllPersons(Db)
	var person domain.Entry
	for _, entry := range sdn.SdnEntry {
		if entry.SdnType != personType {
			continue
		}
		person = domain.Entry{}
		for _, person = range *persons {
			if person.Uid == entry.Uid {
				if person.LastName == entry.LastName && person.FirstName == entry.FirstName {
					person = domain.Entry{}
				}
				break
			}
		}
		if person.Uid == "" {
			helpPerson := &domain.Entry{
				LastName:  entry.LastName,
				FirstName: entry.FirstName,
				Uid:       entry.Uid,
			}
			err = domain.SetPerson(Db, helpPerson)
			if err != nil {
				domain.SetState(Db, &domain.State{Info: err.Error()})
				return
			}

		}
	}

	err = domain.SetState(Db, &domain.State{Info: "ok"})
	if err != nil {
		domain.SetState(Db, &domain.State{Info: err.Error()})
		return
	}
}
