package handler

import (
	"log"

	"github.com/mhshahin/helix/pkg/config"
	"github.com/mhshahin/helix/pkg/medium"
	corev1 "k8s.io/api/core/v1"
)

type EventHandler func(event *corev1.Event)

type Handler struct {
	Medium *medium.Medium
}

// NewHandler...
func NewHandler(medium *medium.Medium) *Handler {
	return &Handler{
		Medium: medium,
	}
}

// OnEvent...
func (h *Handler) OnEvent(ev *corev1.Event) {
	if inRules(ev) {
		err := h.Medium.Matrix.SendEvent(ev)
		if err != nil {
			log.Println("there was an error in sending message to matrix:", err.Error())
			return
		}
	}
}

func inRules(ev *corev1.Event) bool {
	for _, rule := range config.Cfg.Rules {
		if ev.InvolvedObject.Kind == rule.Kind &&
			ev.Type == rule.Type {
			return true
		}
	}

	return false
}

// func (h *Handler) eventMatcher(event *corev1.Event) bool {
// 	for _, rule := range h.rules {
// 		if matchString()
// 	}
// 	return true
// }

// func matchString(pattern, s string) bool {
// 	matched, _ := regexp.MatchString(pattern, s)

// 	return matched
// }
