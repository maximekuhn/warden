package permissions

var roleToActions = map[Role][]Action{
	RoleViewer: {ActionViewServer},
	RoleAdmin:  {ActionViewServer, ActionStartServer, ActionStopServer},
}

var planToPolicies = map[Plan][]Policy{
	PlanFree: {PolicyListServers, PolicyGetServerDetails},
	PlanPro:  {PolicyListServers, PolicyGetServerDetails, PolicyCreateServer},
}
