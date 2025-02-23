package rbac

// Resource is a generic interface representing any data that requires restricted access
type Resource interface {
	ID() int64
	OwnerID() int64
}

// Actor is a generic interface representing a user trying to access a Resource
type Actor interface {
	ID() int64
	Role() string
}
