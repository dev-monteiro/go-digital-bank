package transport

import (
	"log"
	"net/http"
)

type HealthCont struct {
}

func NewHealthCont() *HealthCont {
	cont := &HealthCont{}
	http.HandleFunc("/health", cont.GetHealth)
	return cont
}

func (cont *HealthCont) GetHealth(resWr http.ResponseWriter, req *http.Request) {
	log.Println("[HealthCont] GetHealth")

	resWr.Header().Set("Content-Type", "application/json")

	if req.Method != http.MethodGet {
		resWr.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	resWr.WriteHeader(200)
}
