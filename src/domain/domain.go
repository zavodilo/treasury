package domain

import (
	"errors"
	"github.com/go-pg/pg/v10"
	"strings"
)

type Entry struct {
	Id        int64
	Uid       string
	LastName  string
	FirstName string
}
type State struct {
	Id   int64
	Info string
}

func GetPersonStrong(db *pg.DB, name string) (*Entry, error) {
	if name == "" {
		return nil, errors.New("name empty")
	}
	person := new(Entry)

	names := strings.Split(name, " ")

	query := db.Model(person)
	lenNames := len(names)
	if lenNames == 1 {
		query.Where("last_name = ?", names[0])
	}
	if lenNames == 2 {
		query.Where("last_name = ? AND first_name = ?", names[0], names[1])
	}
	if lenNames > 2 {
		nameHelp := names[1]
		for n := 2; n < lenNames; n++ {
			nameHelp += " " + names[n]
		}
		query.Where("last_name = ? AND first_name = ?", names[0], nameHelp)
	}
	err := query.First()
	return person, err
}

func GetPersonWeak(db *pg.DB, name string) (*[]Entry, error) {
	persons := new([]Entry)

	names := strings.Split(name, " ")
	for i, _ := range names {
		names[i] = "%" + names[i] + "%"
	}
	pgInNames := pg.In(names)
	err := db.Model(persons).
		Where("first_name IN (?) OR last_name IN (?)", pgInNames, pgInNames).
		Select()

	return persons, err
}

func GetAllPersons(db *pg.DB) (*[]Entry, error) {
	persons := new([]Entry)

	err := db.Model(persons).Select()

	return persons, err
}

func SetPerson(db *pg.DB, person *Entry) error {
	if person.Uid == "" {
		return errors.New("uid empty")
	}
	p, err := db.Model(person).Where("uid = ?", person.Uid).Update(&person)
	if err != nil {
		return err
	}
	if p.RowsAffected() == 0 {
		_, err = db.Model(person).Insert()
	}
	return err
}

func GetState(db *pg.DB) (*State, error) {
	state := new(State)

	err := db.Model(state).First()
	return state, err
}

func GetCount(db *pg.DB) (int, error) {
	person := new(Entry)
	return db.Model(person).Count()
}

func SetState(db *pg.DB, state *State) error {
	p, err := db.Model(state).Where("info <> ?", "''").Update(&state)
	if err != nil {
		return err
	}
	if p.RowsAffected() == 0 {
		_, err = db.Model(state).Insert()
	}
	return err
}
