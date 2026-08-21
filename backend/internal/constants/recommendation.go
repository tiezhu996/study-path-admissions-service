package constants

const (
	RecStatusDraft    = "draft"
	RecStatusSent     = "sent"
	RecStatusAccepted = "accepted"
)

func IsValidRecommendationStatus(s string) bool {
	return s == RecStatusDraft || s == RecStatusSent || s == RecStatusAccepted
}

func NextRecommendationStatuses(s string) []string {
	switch s {
	case RecStatusDraft:
		return []string{RecStatusSent}
	case RecStatusSent:
		return []string{RecStatusDraft}
	case RecStatusAccepted:
		return nil
	default:
		return nil
	}
}
