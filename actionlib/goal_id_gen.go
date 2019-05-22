package actionlib

import (
	"fmt"
	"sync"

	"github.com/fetchrobotics/rosgo/ros"
)

type goalIdGenerator struct {
	goals      int
	goalsMutex sync.RWMutex
	nodeName   string
}

func newGoalIdGenerator(nodeName string) *goalIdGenerator {
	return &goalIdGenerator{
		nodeName: nodeName,
	}
}

func (g *goalIdGenerator) generateID() string {
	g.goalsMutex.Lock()
	defer g.goalsMutex.Unlock()

	g.goals++

	timeNow := ros.Now()
	return fmt.Sprintf("%s-%d-%d-%d", g.nodeName, g.goals, timeNow.Sec, timeNow.NSec)
}
