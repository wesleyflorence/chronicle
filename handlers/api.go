package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jomei/notionapi"
	"github.com/wesleyflorence/chronicle/notion"
)

func MedicineEntry(w http.ResponseWriter, r *http.Request, client *notionapi.Client, medicinePageID string) {
	type Payload struct {
		Medicine string
		Note     string
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing medicine entry: %v", err)
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}
	payload := Payload{
		Medicine: r.Form.Get("medicine"),
		Note:     r.Form.Get("note"),
	}

	page, err := notion.AppendMedicineEntry(client, medicinePageID, payload.Medicine, payload.Note)
	if err != nil {
		log.Printf("Error appending medicine entry: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	loc, _ := time.LoadLocation("Local")
	created := page.CreatedTime.In(loc).Format("2006-01-02 03:04PM")
	doseProp, ok := page.Properties["Dose"].(*notionapi.TitleProperty)
	if !ok {
		log.Printf("Error unwrapping Dose returned from page: %v", err)
		http.Error(w, "Error unwrapping Dose returned from page", http.StatusInternalServerError)
		return
	}
	dose := doseProp.Title[0].Text.Content
	response := fmt.Sprintf(`<div id="med-response-target" class="text-xs text-stone-600" hx-ext="remove-me"><div remove-me="5s">%s dose %s :: %s</div></div>`, payload.Medicine, dose, created)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(response))
}

func DigestionEntry(w http.ResponseWriter, r *http.Request, client *notionapi.Client, digestionDbID string) {
	type Payload struct {
		Bristol int
		Size    string
		Note    string
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing digestion entry: %v", err)
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}
	bristol, err := strconv.Atoi(r.Form.Get("bristol"))
	if err != nil {
		log.Printf("Error parsing bristol value: %v", err)
		http.Error(w, "Invalid value for Bristol", http.StatusBadRequest)
		return
	}
	payload := Payload{
		Bristol: bristol,
		Size:    r.Form.Get("size"),
		Note:    r.Form.Get("note"),
	}

	page, err := notion.AppendDigestionEntry(client, digestionDbID, payload.Bristol, payload.Size, payload.Note)
	if err != nil {
		log.Printf("Error appending digestion entry: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	loc, _ := time.LoadLocation("Local")
	created := page.CreatedTime.In(loc).Format("2006-01-02 03:04PM")
	response := fmt.Sprintf(`<div id="dig-response-target" class="text-xs text-stone-600" hx-ext="remove-me"><div remove-me="5s">New Entry :: %s</div></div>`, created)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(response))
}
