package permissions

// Plan represents which subscription plan the user
// currently has.
// It is a 1..1 relation
type Plan string

const (
	PlanFree Plan = "plan.free"
	PlanPro  Plan = "plan.pro"
)
