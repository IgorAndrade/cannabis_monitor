package globo

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
)

const PUBLISHER = "Globo"

type QueryResult struct {
	Publisher string `json:"publisher,omitempty"`
	Page      string `json:"page,omitempty"`
	Title     string `json:"title,omitempty"`
	Detail    string `json:"detail,omitempty"`
	Link      string `json:"link,omitempty"`
	Date      string `json:"Date,omitempty"`
}

func (q QueryResult) GetID() string {
	h := sha1.New()
	h.Write([]byte(q.Title))
	sha1Hash := hex.EncodeToString(h.Sum(nil))
	return sha1Hash
}

func (q QueryResult) String() string {
	if b, err := json.Marshal(q); err == nil {
		return string(b)
	}
	return ""
}
