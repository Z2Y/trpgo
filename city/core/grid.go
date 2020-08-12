package core

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo/common"
)

const gridSize = 64

type Grid struct {
	ecs.BasicEntity
	*common.RenderComponent
	*common.SpaceComponent
}
