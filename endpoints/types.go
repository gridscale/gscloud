package endpoints

import (
	"encoding/json"
	"time"
)

const gsTimeLayout = "2006-01-02T15:04:05Z"

//GSTime is custom time type of gridscale
type GSTime struct {
	time.Time
}

//UnmarshalJSON custom unmarshaller for GSTime
func (t *GSTime) UnmarshalJSON(b []byte) error {
	var tstring string
	if err := json.Unmarshal(b, &tstring); err != nil {
		return err
	}
	parsedTime, err := time.Parse(gsTimeLayout, tstring)
	*t = GSTime{parsedTime}
	return err
}

//MarshalJSON custom marshaller for GSTime
func (t GSTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time.Format(gsTimeLayout))
}
