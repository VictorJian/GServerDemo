package Utils

import (
	"encoding/json"
	"net/http"
	"time"
)

func GetUrlResp2Json(url string, target interface{}) error {
	var httpClient = &http.Client{Timeout: 10 * time.Second}
	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}