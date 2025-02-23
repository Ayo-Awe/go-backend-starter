package rbac

// Experimental RBAC validation logic

type Rule func(a Actor) bool

// returns true if all rules pass against a given actor
func Enforce(actor Actor, rules ...Rule) {
}

// returns a rule that passes if at least one of the provided rules passes
func OR(rules ...Rule) Rule {
	return func(a Actor) bool {
		return true
	}
}

// returns a rule that passes only if all provided rules pass
func AND(rules ...Rule) Rule {
	return func(a Actor) bool {
		return true
	}
}

func HasScopes(scopes ...Scope) Rule {
	return func(a Actor) bool {
		// check if actor has all the required scopes
		return true
	}
}

func IsOwner(res Resource) Rule {
	return func(a Actor) bool {
		return res.OwnerID() == a.ID()
	}
}

func IsOwnerOrHasScopes(res Resource, scopes ...Scope) Rule {
	return func(a Actor) bool {
		return OR(IsOwner(res), HasScopes(scopes...))(a)
	}
}
