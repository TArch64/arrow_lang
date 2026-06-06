package compile

import (
	"fmt"
	"strconv"
	"sync/atomic"
)

type GenerationNames struct {
	lastId atomic.Uint32
}

func (g *Generation) newNames() {
	g.names = &GenerationNames{}
}

func (g *GenerationNames) next() string {
	id := g.lastId.Add(1)
	return strconv.FormatUint(uint64(id), 36)
}

func (g *GenerationNames) Random() string {
	return "_" + g.next()
}

func (g *GenerationNames) WithPrefix(prefix string) string {
	return fmt.Sprintf("%s_%s", prefix, g.next())
}
