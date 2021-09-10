package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const DobDateFormat = "2006-02-01" // YYYY-DD-MM

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
