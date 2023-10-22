package deprecated

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
)

func (h *handler) MDConvert(w http.ResponseWriter, r *http.Request) {

	md, _ := io.ReadAll(r.Body)
	rg := regexp.MustCompile(`(?:[\t ]*(?:\r?\n|\r))`)
	html := MdToHTML(md)
	str := string(html)
	result := rg.ReplaceAllString(str, "")
	html = []byte(result)
	w.Header().Set("Content-Type", "application/json")
	responseJSON := map[string]string{"msg": string(html)}
	json.NewEncoder(w).Encode(responseJSON)
}
