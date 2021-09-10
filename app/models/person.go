package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const DobDateFormat = "02/01/2006" // DD/MM/YYYY

type JSONDate time.Time

func (date JSONDate) MarshalJSON() ([]byte, error) {
	d := fmt.Sprintf("\"%s\"", time.Time(date).Format(DobDateFormat))
	return []byte(d), nil
}

func (date *JSONDate) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	d, err := time.Parse(DobDateFormat, s)
	*date = JSONDate(d)
	return err
}

func (date JSONDate) MarshalText() (string, error) {
	d := fmt.Sprintf("%s", time.Time(date).Format(DobDateFormat))
	return d, nil
}

type Person struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Address string `json:"address"`
	DateOfBirth JSONDate `json:"dateOfBirth"`
}

func (p *Person) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Person) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &p); err != nil {
		return err
	}
	return nil
}
