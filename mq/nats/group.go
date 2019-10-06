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
	handlers pubsub.HandlersChain
}

// Group -
func (group *Group) Group(relativePath string, handlers ...pubsub.HandlerFunc) pubsub.Group {
	return &Group{
		handlers: group.combineHandlers(handlers),
		basePath: group.calculateAbsolutePath(relativePath),
		broker:   group.broker,
	}
}

// Use -
func (group *Group) Use(middlewares ...pubsub.HandlerFunc) pubsub.IRoutes {
	return nil
}

func (group *Group) combineHandlers(handlers pubsub.HandlersChain) pubsub.HandlersChain {
	finalSize := len(group.handlers) + len(handlers)
	if finalSize >= int(abortIndex) {
		panic("too many handlers")
	}

	mergedHandlers := make(pubsub.HandlersChain, finalSize)

	copy(mergedHandlers, group.handlers)
	copy(mergedHandlers[len(group.handlers):], handlers)

	return mergedHandlers
}

func (group *Group) calculateAbsolutePath(relativePath string) string {
	return joinPaths(group.basePath, relativePath)
}

func (group *Group) returnObj() pubsub.IRoutes {
	if group.root {
		return group.broker
	}
	return group
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
