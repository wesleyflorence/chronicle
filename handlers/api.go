// Package handlers defines route handlers
package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jomei/notionapi"
	"github.com/wesleyflorence/chronicle/components"
	"github.com/wesleyflorence/chronicle/notion"
)

// MedicineEntry parses medicine form and stores values in notion
func MedicineEntry(w http.ResponseWriter, r *http.Request, client *notionapi.Client, medicinePageID string) {
	type Payload struct {
		Medicine string
		Size     int
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

	if (r.Form.Get("milkOfMagnesia")) != "" {
		size, err := strconv.Atoi(r.Form.Get("milkOfMagnesia"))
		if err != nil {
			log.Printf("Error parsing milkOfMagnesia value: %v", err)
			http.Error(w, "Invalid value for milkOfMagnesia", http.StatusBadRequest)
			return
		}
		payload.Medicine = "Milk of Magnesia"
		payload.Size = size
	}

	page, err := notion.AppendMedicineEntry(client, medicinePageID, payload.Medicine, payload.Size, payload.Note)
	if err != nil {
		log.Printf("Error appending medicine entry: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	loc, _ := time.LoadLocation("America/Los_Angeles")
	created := page.CreatedTime.In(loc).Format("2006-01-02 03:04PM")
	doseProp, ok := page.Properties["Dose"].(*notionapi.TitleProperty)
	if !ok {
		log.Printf("Error unwrapping Dose returned from page: %v", err)
		http.Error(w, "Error unwrapping Dose returned from page", http.StatusInternalServerError)
		return
	}
	dose := doseProp.Title[0].Text.Content
	component := components.MedSuccess(payload.Medicine, dose, created)
	component.Render(r.Context(), w)
}

// DigestionEntry parses medicine form and stores values in notion
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

	loc, _ := time.LoadLocation("America/Los_Angeles")
	created := page.CreatedTime.In(loc).Format("2006-01-02 03:04PM")
	component := components.DigSuccess(created)
	component.Render(r.Context(), w)
}
