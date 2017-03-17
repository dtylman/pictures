package progressbar

import "encoding/json"

//Status represents a status that can be shown in a progress bar
type Status struct {
	Percentage int    `json:"percentage"`
	Text       string `json:"text"`
	Done       bool   `json:"done"`
}

//ToJSON returns the status as json string
func (s *Status) ToJSON() (string, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
