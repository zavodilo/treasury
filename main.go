package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"strconv"
	"treasury/entity"
	"treasury/server/response"
)

const (
	errCode    int    = 503
	okCode     int    = 200
	personType string = "Individual"
)

// TODO сделать асинхронную загрузку в бд
func main() {
	http.HandleFunc("/update", updateHandler)
	http.HandleFunc("/state", stateHandler)
	http.ListenAndServe(":8080", nil)
}

func stateHandler(w http.ResponseWriter, r *http.Request) {
	//TODO сделать получение статуса
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	err := update()

	if err != nil {
		resp = response.Response{
			Code:   errCode,
			Result: false,
			Info:   err.Error(),
		}
	}
	if err == nil {
		resp = response.Response{
			Code:   okCode,
			Result: true,
			Info:   "",
		}
	}
	JsonResponse(w, resp, errCode)
}

func update() error {
	url := "https://www.treasury.gov/ofac/downloads/sdn.xml"
	resp, err := http.Get(url)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//сохранение в файл для изучения
	//out, err := os.Create("sdn.xml")
	//checkErr(err)
	//defer out.Close()
	//io.Copy(out, resp.Body)
	if resp.StatusCode != 200 {
		return errors.New("StatusCode error: " + strconv.Itoa(resp.StatusCode))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	sdn := new(entity.Sdn)
	err = xml.Unmarshal(body, sdn)
	if err != nil {
		return err
	}

	var persons []*entity.Person
	for _, person := range sdn.SdnEntry {
		if person.SdnType == personType {
			helpPerson := &entity.Person{
				LastName:  person.LastName,
				FirstName: person.FirstName,
				Uid:       person.Uid,
			}
			persons = append(persons, helpPerson)
		}
	}

	//println(persons)
	return nil
}

func JsonResponse(w http.ResponseWriter, resp response.Response, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resp)
}
