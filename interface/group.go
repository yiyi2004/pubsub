package pubsub

// Group -
type Group interface {
	IRoutes
	Group(relativePath string, handlers ...HandlerFunc) Group
}

// IRoutes -
type IRoutes interface {
	Use(middlewares ...HandlerFunc) IRoutes
}

var _ Group = &DefaultGroup{}

// DefaultGroup -
type DefaultGroup struct {
	Root     bool
	BashPath string
	Broker   Broker
	Handlers HandlersChain
}

// Group -
func (group *DefaultGroup) Group(relativePath string, handlers ...HandlerFunc) Group {
	return nil
}

// Use -
func (group *DefaultGroup) Use(middlewares ...HandlerFunc) IRoutes {
	return nil
}
