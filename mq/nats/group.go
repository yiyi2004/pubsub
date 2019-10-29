package nats

import (
	"math"
	"path"
	"reflect"
	"runtime"

	pubsub "github.com/zhangce1999/pubsub/interface"
)

const abortIndex int8 = math.MaxInt8 / 2

// group -
type group struct {
	root     bool
	basePath string
	broker   *Broker
	Handlers pubsub.HandlersChain
}

// group -
func (g *group) Group(relativePath string, Handlers ...pubsub.HandlerFunc) pubsub.Group {
	return &group{
		Handlers: g.combineHandlers(Handlers),
		basePath: g.calculateAbsolutePath(relativePath),
		broker:   g.broker,
	}
}

// Use -
func (g *group) Use(middlewares ...pubsub.HandlerFunc) pubsub.IRoutes {
	return nil
}

func (g *group) combineHandlers(Handlers pubsub.HandlersChain) pubsub.HandlersChain {
	finalSize := len(g.Handlers) + len(Handlers)
	if finalSize >= int(abortIndex) {
		panic("too many Handlers")
	}

	mergedHandlers := make(pubsub.HandlersChain, finalSize)

	copy(mergedHandlers, g.Handlers)
	copy(mergedHandlers[len(g.Handlers):], Handlers)

	return mergedHandlers
}

func (g *group) calculateAbsolutePath(relativePath string) string {
	return joinPaths(g.basePath, relativePath)
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

func (g *group) returnObj() pubsub.IRoutes {
	if g.root {
		return g.broker
	}
	return g
}
