package nats

import (
	"math"
	"path"
	"reflect"
	"runtime"

	pubsub "github.com/zhangce1999/pubsub/interface"
)

const abortIndex int8 = math.MaxInt8 / 2

// Group -
type Group struct {
	root     bool
	basePath string
	broker   *Broker
	Handlers pubsub.HandlersChain
}

// Group -
func (group *Group) Group(relativePath string, Handlers ...pubsub.HandlerFunc) pubsub.Group {
	return &Group{
		Handlers: group.combineHandlers(Handlers),
		basePath: group.calculateAbsolutePath(relativePath),
		broker:   group.broker,
	}
}

// Use -
func (group *Group) Use(middlewares ...pubsub.HandlerFunc) pubsub.IRoutes {
	return nil
}

func (group *Group) combineHandlers(Handlers pubsub.HandlersChain) pubsub.HandlersChain {
	finalSize := len(group.Handlers) + len(Handlers)
	if finalSize >= int(abortIndex) {
		panic("too many Handlers")
	}

	mergedHandlers := make(pubsub.HandlersChain, finalSize)

	copy(mergedHandlers, group.Handlers)
	copy(mergedHandlers[len(group.Handlers):], Handlers)

	return mergedHandlers
}

func (group *Group) calculateAbsolutePath(relativePath string) string {
	return joinPaths(group.basePath, relativePath)
}

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	appendSlash := lastChar(relativePath) == '/' && lastChar(finalPath) != '/'
	if appendSlash {
		return finalPath + "/"
	}
	return finalPath
}

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

func nameOfFunction(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func insertBackSlash(str string) string {
	finalPath := path.Join("/", str)
	appendSlash := lastChar(str) == '/' && lastChar(finalPath) != '/'
	if appendSlash {
		return finalPath + "/"
	}
	return finalPath
}

func (group *Group) returnObj() pubsub.IRoutes {
	if group.root {
		return group.broker
	}
	return group
}
