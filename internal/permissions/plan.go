package permissions

import "fmt"

// Plan represents which subscription plan the user
// currently has.
// It is a 1..1 relation
type Plan string

const (
	PlanFree Plan = "plan.free"
	PlanPro  Plan = "plan.pro"
)

func PlanFromString(s string) (Plan, error) {
	if s == string(PlanFree) {
		return PlanFree, nil
	}
	if s == string(PlanPro) {
		return PlanPro, nil
	}
	return "", fmt.Errorf("unknown plan: '%s'", s)
}
