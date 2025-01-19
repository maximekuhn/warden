package permissions

// Policy is associated to a user plan and represent what is available
// to do for this specific plan.
type Policy string

const (
	PolicyCreateServer     Policy = "policy.createServer"
	PolicyListServers      Policy = "policy.listServers"
	PolicyGetServerDetails Policy = "policy.getServerDetails"
)
