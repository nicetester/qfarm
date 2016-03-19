package worker

import (
	"encoding/json"
	"log"

	"github.com/qfarm/qfarm/redis"
)

type Notifier struct {
	redis *redis.Service
}

func NewNotifier(redis *redis.Service) *Notifier {
	return &Notifier{redis: redis}
}

func (n *Notifier) SendEventWithPayload(repo, desc, eventType, payload string) {
	if n.redis == nil {
		log.Printf("WARNING: Redis is not configured. Skip sending event!")
		return
	}

	e := Event{Description: desc, Repo: repo, Type: eventType, Payload: payload}

	data, err := json.Marshal(e)
	if err != nil {
		log.Printf("Can't marshal event. Err: %v", err)
	} else {
		err := n.redis.Publish("events", data)
		if err != nil {
			log.Printf("Can't send event to subscribers. Err: %v", err)
		}
	}
}

func (n *Notifier) SendEvent(repo, desc, eventType string) {
	n.SendEventWithPayload(repo, desc, eventType, "")
}

type Event struct {
	Repo        string `json:"repo,omitempty"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
	Payload     string `json:"payload,omitempty"`
}

const (
	EventTypeAllDone      = "all-done"
	EventTypeDownloadDone = "download-done"

	EventTypeAligncheckDone  = "aligncheck-done"
	EventTypeDeadcodeDone    = "deadcode-done"
	EventTypeDuplDone        = "dupl-done"
	EventTypeErrcheckDone    = "errcheck-done"
	EventTypeGoconstDone     = "goconst-done"
	EventTypeGocycloDone     = "gocyclo-done"
	EventTypeGofmtDone       = "gofmt-done"
	EventTypeGoimportsDone   = "goimports-done"
	EventTypeGolintDone      = "golint-done"
	EventTypeGotypeDone      = "gotype-done"
	EventTypeIneffassignDone = "ineffassign-done"
	EventTypeInterfacerDone  = "interfacer-done"
	EventTypeLllDone         = "lll-done"
	EventTypeStructcheckDone = "structcheck-done"
	EventTypeTestDone        = "test-done"
	EventTypeTestifyDone     = "testify-done"
	EventTypeVarcheckDone    = "varcheck-done"
	EventTypeVetDone         = "vet-done"
	EventTypeVetshadowDone   = "vetshadow-done"
	EventTypeUnconvertDone   = "unconvert-done"
	EventTypeMetalinterErr   = "metalinter-error"

	EventTypeCoverageDone = "coverage-done"
	EventTypeCoverageErr  = "coverage-error"
	EventTypeError        = "error"

	EventTypeAlreadyAnalyzed = "already-analyzed"
)

var linterEventsMapping = map[string]string{
	"aligncheck":  EventTypeAligncheckDone,
	"deadcode":    EventTypeDeadcodeDone,
	"dupl":        EventTypeDuplDone,
	"errcheck":    EventTypeErrcheckDone,
	"goconst":     EventTypeGoconstDone,
	"gocyclo":     EventTypeGocycloDone,
	"gofmt":       EventTypeGofmtDone,
	"goimports":   EventTypeGoimportsDone,
	"golint":      EventTypeGolintDone,
	"gotype":      EventTypeGotypeDone,
	"ineffassign": EventTypeIneffassignDone,
	"interfacer":  EventTypeInterfacerDone,
	"lll":         EventTypeLllDone,
	"structcheck": EventTypeStructcheckDone,
	"test":        EventTypeTestDone,
	"testify":     EventTypeTestifyDone,
	"varcheck":    EventTypeVarcheckDone,
	"vet":         EventTypeVetDone,
	"vetshadow":   EventTypeVetshadowDone,
	"unconvert":   EventTypeUnconvertDone,
}
