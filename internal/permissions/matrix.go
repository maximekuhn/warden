package permissions

var roleToActions = map[Role][]Action{
	RoleViewer: {},
	RoleAdmin:  {ActionStartServer, ActionStopServer},
}

var planToPolicies = map[Plan][]Policy{
	PlanFree: {PolicyListServers, PolicyGetServerDetails},
	PlanPro:  {PolicyListServers, PolicyGetServerDetails, PolicyCreateServer},
}
