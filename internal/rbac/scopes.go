package rbac

// map of roles to scopes
var scopesMap = map[string][]Scope{}

func GetScopesForRole(role string) []Scope {
	return scopesMap[role]
}
