/*
 _____               _____
|_   _| __ _   _  __|_   _|__ _ __ ___  _ __ ___
  | || '__| | | |/ _ \| |/ _ \ '_ ` _ \| '_ ` _ \
  | || |  | |_| |  __/| |  __/ | | | | | | | | | |
  |_||_|   \__,_|\___||_|\___|_| |_| |_|_| |_| |_|

 _____ __  __       ____                               __ _
|_   _|  \/  |     |  _ \ _ __ __ _  __ _  ___  _ __  / _| |_   _
  | | | |\/| |_____| | | | '__/ _` |/ _` |/ _ \| '_ \| |_| | | | |
  | | | |  | |_____| |_| | | | (_| | (_| | (_) | | | |  _| | |_| |
  |_| |_|  |_|     |____/|_|  \__,_|\__, |\___/|_| |_|_| |_|\__, |
                                    |___/                   |___/

@author TrueTemm
@link   https://github.com/TrueTemm
TM-Dragonfly Project
*/

package world

import (
	"maps"
	"slices"

	"github.com/df-mc/dragonfly/server/block/cube"
)

type RedstoneUpdateCause uint8

const (
	RedstoneUpdateCauseBlockUpdate RedstoneUpdateCause = iota

	RedstoneUpdateCauseScheduledTick

	RedstoneUpdateCauseCompilerRebuild
)

type RedstoneUpdate struct {
	Pos cube.Pos

	ChangedNeighbour cube.Pos

	HasChangedNeighbour bool

	ChangedRedstoneRelevant bool

	Source cube.Pos

	HasSource bool

	Before Block

	After Block

	OldPower int

	NewPower int

	CurrentTick int64

	Cause RedstoneUpdateCause
}

type RedstonePowerSource interface {
	RedstonePower(pos cube.Pos, tx *Tx, face cube.Face) int
}

type RedstoneStrongPowerSource interface {
	RedstoneStrongPower(pos cube.Pos, tx *Tx, face cube.Face) int
}

type RedstoneWeakBlockPowerer interface {
	RedstoneWeaklyPowersBlocks() bool
}

type RedstonePowerRelayer interface {
	RedstoneSignalLoss(pos cube.Pos, tx *Tx) int
}

type RedstonePowerRelayerNeighbourer interface {
	RedstoneRelayerNeighbours(pos cube.Pos, tx *Tx) []cube.Pos
}

type RedstonePowerConsumer interface {
	RedstonePowerUpdate(pos cube.Pos, tx *Tx, power int) (after Block, changed bool)
}

type RedstonePowerPostUpdater interface {
	RedstonePowerPostUpdate(pos cube.Pos, tx *Tx, before, after Block, oldPower, newPower int)
}

type RedstonePowerAction interface {
	RedstonePowerAction(pos cube.Pos, tx *Tx, oldPower, newPower int)
}

type RedstonePowerContextAction interface {
	RedstonePowerActionUpdate(pos cube.Pos, tx *Tx, update RedstoneUpdate)
}

type RedstoneNonConductive interface {
	RedstoneNonConductive()
}

type redstoneEngine struct {
	currentTick       int64
	dirty             map[cube.Pos]redstoneDirty
	power             map[cube.Pos]int
	output            map[cube.Pos]int
	evaluating        map[cube.Pos]struct{}
	suppressedSources map[cube.Pos]int
	torchBurnout      map[cube.Pos]redstoneTorchBurnout
}

type redstoneDirty struct {
	changed                 cube.Pos
	hasChanged              bool
	changedRedstoneRelevant bool
	source                  cube.Pos
	hasSource               bool
	cause                   RedstoneUpdateCause
}

type redstoneTorchBurnout struct {
	offTicks             []int64
	burnedOut            bool
	pendingSelfTriggered bool
}

const (
	redstoneTorchBurnoutThreshold   = 8
	redstoneTorchBurnoutWindowTicks = 60
)

type redstoneGraph struct {
	nodes []redstoneNode
	edges []redstoneEdge
}

type redstoneNode struct {
	pos    cube.Pos
	source bool
	sink   bool
}

type redstoneEdge struct {
	from, to int
	weight   int
}

func newRedstoneEngine(tick int64) *redstoneEngine {
	return &redstoneEngine{
		currentTick: tick,
		dirty:       make(map[cube.Pos]redstoneDirty),
		power:       make(map[cube.Pos]int),
		output:      make(map[cube.Pos]int),
		evaluating:  make(map[cube.Pos]struct{}),
	}
}

func (e *redstoneEngine) invalidateAround(pos, changed cube.Pos, cause RedstoneUpdateCause, r cube.Range) {
	e.invalidateAroundWith(pos, redstoneDirty{changed: changed, hasChanged: true, source: changed, hasSource: true, cause: cause}, r)
}

func (e *redstoneEngine) invalidateAroundBlockChange(pos cube.Pos, before, after Block, cause RedstoneUpdateCause, r cube.Range) {
	d := redstoneDirty{
		changed:                 pos,
		hasChanged:              true,
		changedRedstoneRelevant: isRedstoneRelevant(before) || isRedstoneRelevant(after),
		source:                  pos,
		hasSource:               true,
		cause:                   cause,
	}
	e.invalidateAroundWith(pos, d, r)
}

func (e *redstoneEngine) invalidateAroundWith(pos cube.Pos, d redstoneDirty, r cube.Range) {
	if e == nil || pos.OutOfBounds(r) {
		return
	}
	e.invalidate(pos, d, r)
	pos.Neighbours(func(neighbour cube.Pos) {
		e.invalidate(neighbour, d, r)
	}, r)
}

func (e *redstoneEngine) invalidate(pos cube.Pos, d redstoneDirty, r cube.Range) {
	if pos.OutOfBounds(r) {
		return
	}
	if existing, ok := e.dirty[pos]; ok {
		e.dirty[pos] = mergeRedstoneDirty(existing, d)
		return
	}
	e.dirty[pos] = d
}

func mergeRedstoneDirty(a, b redstoneDirty) redstoneDirty {
	if redstoneDirtyPriority(b) >= redstoneDirtyPriority(a) {
		b.changedRedstoneRelevant = a.changedRedstoneRelevant || b.changedRedstoneRelevant
		return b
	}
	a.changedRedstoneRelevant = a.changedRedstoneRelevant || b.changedRedstoneRelevant
	return a
}

func redstoneDirtyPriority(d redstoneDirty) int {
	switch d.cause {
	case RedstoneUpdateCauseBlockUpdate:
		return 3
	case RedstoneUpdateCauseScheduledTick:
		return 2
	case RedstoneUpdateCauseCompilerRebuild:
		return 1
	default:
		return 0
	}
}

func (e *redstoneEngine) removeChunk(chunkPos ChunkPos) {
	if e == nil {
		return
	}
	maps.DeleteFunc(e.dirty, func(pos cube.Pos, dirty redstoneDirty) bool {
		return chunkPosFromBlockPos(pos) == chunkPos ||
			(dirty.hasChanged && chunkPosFromBlockPos(dirty.changed) == chunkPos)
	})
	maps.DeleteFunc(e.power, func(pos cube.Pos, _ int) bool {
		return chunkPosFromBlockPos(pos) == chunkPos
	})
	maps.DeleteFunc(e.output, func(pos cube.Pos, _ int) bool {
		return chunkPosFromBlockPos(pos) == chunkPos
	})
	maps.DeleteFunc(e.torchBurnout, func(pos cube.Pos, _ redstoneTorchBurnout) bool {
		return chunkPosFromBlockPos(pos) == chunkPos
	})
}

func (e *redstoneEngine) forget(pos cube.Pos) {
	if e == nil {
		return
	}
	delete(e.power, pos)
	delete(e.output, pos)
	delete(e.evaluating, pos)
}

func (e *redstoneEngine) tick(tx *Tx, tick int64) {
	if e == nil || len(e.dirty) == 0 {
		return
	}
	e.currentTick = tick
	dirty := maps.Clone(e.dirty)
	clear(e.dirty)

	candidates := slices.Collect(maps.Keys(dirty))
	slices.SortFunc(candidates, compareBlockPos)

	graph := e.compile(tx, candidates)
	cancelledSources, checkedSources := e.updateGraphSources(tx, graph, dirty)
	previousSuppressed := e.suppressedSources
	e.suppressedSources = cancelledSources
	defer func() {
		e.suppressedSources = previousSuppressed
	}()

	powers := e.graphPower(tx, graph)
	for i, node := range graph.nodes {
		d := redstoneDirtyContext(dirty, node.pos)
		if node.sink {
			e.update(tx, node.pos, d, powers[i])
		}
	}
	for _, node := range graph.nodes {
		if _, ok := checkedSources[node.pos]; ok {
			continue
		}
		d := redstoneDirtyContext(dirty, node.pos)
		if node.source {
			e.updateSource(tx, node.pos, d)
		}
	}
}

func redstoneDirtyContext(dirty map[cube.Pos]redstoneDirty, pos cube.Pos) redstoneDirty {
	if d, ok := dirty[pos]; ok {
		return d
	}
	var (
		bestPos  cube.Pos
		best     redstoneDirty
		bestDist int
		ok       bool
	)
	for dirtyPos, d := range dirty {
		dist := redstoneManhattanDistance(pos, dirtyPos)
		if !ok || dist < bestDist || (dist == bestDist && compareBlockPos(dirtyPos, bestPos) < 0) {
			bestPos, best, bestDist, ok = dirtyPos, d, dist, true
		}
	}
	if !ok {
		return redstoneDirty{changed: pos, hasChanged: true, source: pos, hasSource: true, cause: RedstoneUpdateCauseCompilerRebuild}
	}
	return best
}

func (d redstoneDirty) redstoneUpdate(pos cube.Pos, before Block, oldPower, newPower int, tick int64) RedstoneUpdate {
	return RedstoneUpdate{
		Pos:                     pos,
		ChangedNeighbour:        d.changed,
		HasChangedNeighbour:     d.hasChanged,
		ChangedRedstoneRelevant: d.changedRedstoneRelevant,
		Source:                  d.source,
		HasSource:               d.hasSource,
		Before:                  before,
		OldPower:                oldPower,
		NewPower:                newPower,
		CurrentTick:             tick,
		Cause:                   d.cause,
	}
}

func (d redstoneDirty) propagatedFrom(pos cube.Pos) redstoneDirty {
	d.changed = pos
	d.hasChanged = true
	d.changedRedstoneRelevant = true
	return d
}

func redstoneManhattanDistance(a, b cube.Pos) int {
	return abs(a[0]-b[0]) + abs(a[1]-b[1]) + abs(a[2]-b[2])
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func (e *redstoneEngine) compile(tx *Tx, candidates []cube.Pos) redstoneGraph {
	nodes := make([]redstoneNode, 0, len(candidates))
	seen := make(map[cube.Pos]struct{}, len(candidates)*2)
	for _, pos := range candidates {
		nodeCount := len(nodes)
		e.compileRegion(tx, pos, seen, &nodes)
		if len(nodes) == nodeCount {
			e.compileAdjacentRedstone(tx, pos, seen, &nodes)
		}
	}
	for i := 0; i < len(nodes); i++ {
		e.compileAdjacentRedstone(tx, nodes[i].pos, seen, &nodes)
	}
	slices.SortFunc(nodes, func(a, b redstoneNode) int {
		return compareBlockPos(a.pos, b.pos)
	})
	edges := e.compileEdges(tx, nodes)
	return redstoneGraph{nodes: nodes, edges: edges}
}

func (e *redstoneEngine) compileAdjacentRedstone(tx *Tx, pos cube.Pos, seen map[cube.Pos]struct{}, nodes *[]redstoneNode) {
	if b, ok := tx.World().blockLoaded(pos); ok && e.redstoneBlockMayConduct(tx, pos, b) {
		pos.Neighbours(func(neighbour cube.Pos) {
			if b, ok := tx.World().blockLoaded(neighbour); ok && isRedstoneRelevant(b) {
				e.compileRegion(tx, neighbour, seen, nodes)
			}
		}, tx.Range())
	}
	pos.Neighbours(func(neighbour cube.Pos) {
		b, ok := tx.World().blockLoaded(neighbour)
		if !ok {
			return
		}
		if isRedstoneRelevant(b) {
			e.compileRegion(tx, neighbour, seen, nodes)
		}
		if e.redstoneBlockMayConduct(tx, neighbour, b) {
			neighbour.Neighbours(func(conductedNeighbour cube.Pos) {
				if b, ok := tx.World().blockLoaded(conductedNeighbour); ok && isRedstoneRelevant(b) {
					e.compileRegion(tx, conductedNeighbour, seen, nodes)
				}
			}, tx.Range())
		}
	}, tx.Range())
}

func (e *redstoneEngine) redstoneBlockMayConduct(tx *Tx, pos cube.Pos, b Block) bool {
	for _, face := range cube.Faces() {
		if redstoneStrongPowerConductor(pos, b, tx, face) {
			return true
		}
	}
	return false
}

func (e *redstoneEngine) compileRegion(tx *Tx, pos cube.Pos, seen map[cube.Pos]struct{}, nodes *[]redstoneNode) {
	if _, ok := seen[pos]; ok || pos.OutOfBounds(tx.Range()) {
		return
	}
	queue := []cube.Pos{pos}
	for len(queue) != 0 {
		p := queue[0]
		queue = queue[1:]
		if _, ok := seen[p]; ok || p.OutOfBounds(tx.Range()) {
			continue
		}
		seen[p] = struct{}{}

		b, ok := tx.World().blockLoaded(p)
		if !ok {
			continue
		}
		source, consumer, action, relayer := classifyRedstoneBlock(b)
		if !source && !consumer && !action && !relayer {
			continue
		}
		*nodes = append(*nodes, redstoneNode{
			pos:    p,
			source: source,
			sink:   consumer || action,
		})
		if !relayer {
			continue
		}
		for _, neighbour := range e.redstoneRelayerConnectedPositions(tx, p, b) {
			if b, ok := tx.World().blockLoaded(neighbour); ok && isRedstoneRelevant(b) {
				queue = append(queue, neighbour)
			}
		}
	}
}

func (e *redstoneEngine) update(tx *Tx, pos cube.Pos, d redstoneDirty, newPower int) {
	b := tx.Block(pos)
	oldPower, newPower := e.power[pos], ClampRedstonePower(newPower)

	after, blockChanged := b, false
	if consumer, ok := b.(RedstonePowerConsumer); ok {
		after, blockChanged = consumer.RedstonePowerUpdate(pos, tx, newPower)
	}
	action, hasAction := b.(RedstonePowerAction)
	contextAction, hasContextAction := b.(RedstonePowerContextAction)

	update := d.redstoneUpdate(pos, b, oldPower, newPower, e.currentTick)
	if blockChanged {
		update.After = after
	}
	shouldRunAction := hasContextAction || (hasAction && oldPower != newPower)
	if oldPower != newPower || blockChanged || shouldRunAction {
		if !e.redstoneUpdateAllowed(tx, update) {
			return
		}
	}

	if !blockChanged && !shouldRunAction {
		storeRedstonePower(e.power, pos, newPower)
		return
	}

	if blockChanged {
		tx.SetBlock(pos, after, &SetOpts{DisableRedstoneUpdates: true})
		e.invalidateAroundWith(pos, d.propagatedFrom(pos), tx.Range())
		if postUpdater, ok := b.(RedstonePowerPostUpdater); ok {
			postUpdater.RedstonePowerPostUpdate(pos, tx, b, after, oldPower, newPower)
		}
	}
	if hasContextAction {
		contextAction.RedstonePowerActionUpdate(pos, tx, update)
	} else if shouldRunAction {
		action.RedstonePowerAction(pos, tx, oldPower, newPower)
	}
	if blockChanged || shouldRunAction {
		storeRedstonePower(e.power, pos, newPower)
	}
}

func (e *redstoneEngine) updateGraphSources(tx *Tx, graph redstoneGraph, dirty map[cube.Pos]redstoneDirty) (map[cube.Pos]int, map[cube.Pos]struct{}) {
	var cancelled map[cube.Pos]int
	var checked map[cube.Pos]struct{}
	for _, node := range graph.nodes {
		if !node.source {
			continue
		}
		b, ok := tx.World().blockLoaded(node.pos)
		if !ok {
			continue
		}
		if _, ok := b.(RedstonePowerRelayer); ok {
			continue
		}
		if checked == nil {
			checked = make(map[cube.Pos]struct{})
		}
		checked[node.pos] = struct{}{}

		d := redstoneDirtyContext(dirty, node.pos)
		if !e.updateSource(tx, node.pos, d) {
			if cancelled == nil {
				cancelled = make(map[cube.Pos]int)
			}
			cancelled[node.pos] = e.output[node.pos]
		}
	}
	return cancelled, checked
}

func (e *redstoneEngine) updateSource(tx *Tx, pos cube.Pos, d redstoneDirty) bool {
	b := tx.Block(pos)
	oldPower, newPower := e.output[pos], e.sourcePower(pos, tx)
	if oldPower == newPower {
		return true
	}
	update := d.redstoneUpdate(pos, b, oldPower, newPower, e.currentTick)
	if !e.redstoneUpdateAllowed(tx, update) {
		return false
	}
	storeRedstonePower(e.output, pos, newPower)
	e.invalidateAroundWith(pos, d.propagatedFrom(pos), tx.Range())
	return true
}

func storeRedstonePower(cache map[cube.Pos]int, pos cube.Pos, power int) {
	if power == 0 {
		delete(cache, pos)
		return
	}
	cache[pos] = power
}

func (e *redstoneEngine) directPower(pos cube.Pos, tx *Tx) int {
	power := 0
	for _, face := range cube.Faces() {
		power = max(power, e.directPowerFrom(pos, tx, face))
	}
	return power
}

func (e *redstoneEngine) directPowerFrom(pos cube.Pos, tx *Tx, face cube.Face) int {
	neighbour := pos.Side(face)
	if neighbour.OutOfBounds(tx.Range()) {
		return 0
	}
	b, ok := tx.World().blockLoaded(neighbour)
	if !ok {
		return 0
	}
	if source, ok := b.(RedstonePowerSource); ok {
		return ClampRedstonePower(e.redstonePower(source, neighbour, tx, face.Opposite()))
	}
	return 0
}

func (e *redstoneEngine) strongPower(pos cube.Pos, tx *Tx) int {
	power := 0
	for _, face := range cube.Faces() {
		power = max(power, e.strongPowerFrom(pos, tx, face))
	}
	return power
}

func (e *redstoneEngine) strongPowerFrom(pos cube.Pos, tx *Tx, face cube.Face) int {
	neighbour := pos.Side(face)
	if neighbour.OutOfBounds(tx.Range()) {
		return 0
	}
	b, ok := tx.World().blockLoaded(neighbour)
	if !ok {
		return 0
	}
	if source, ok := b.(RedstoneStrongPowerSource); ok {
		if power, ok := e.suppressedSources[neighbour]; ok {
			return ClampRedstonePower(power)
		}
		return ClampRedstonePower(source.RedstoneStrongPower(neighbour, tx, face.Opposite()))
	}
	return 0
}

func (e *redstoneEngine) conductedStrongPower(pos cube.Pos, tx *Tx) int {
	power := 0
	for _, face := range cube.Faces() {
		power = max(power, e.conductedStrongPowerFrom(pos, tx, face))
	}
	return power
}

func (e *redstoneEngine) conductedStrongPowerFrom(pos cube.Pos, tx *Tx, face cube.Face) int {
	conductorPos := pos.Side(face)
	if conductorPos.OutOfBounds(tx.Range()) {
		return 0
	}
	conductor, ok := tx.World().blockLoaded(conductorPos)
	if !ok || !redstoneStrongPowerConductor(conductorPos, conductor, tx, face.Opposite()) {
		return 0
	}
	power := 0
	for _, sourceFace := range cube.Faces() {
		power = max(power, e.strongPowerFrom(conductorPos, tx, sourceFace))
	}
	return power
}

func (e *redstoneEngine) weakBlockPower(pos cube.Pos, tx *Tx) int {
	power := 0
	for _, face := range cube.Faces() {
		power = max(power, e.weakBlockPowerFrom(pos, tx, face))
	}
	return power
}

func (e *redstoneEngine) weakBlockPowerFrom(pos cube.Pos, tx *Tx, face cube.Face) int {
	sourcePos := pos.Side(face)
	if sourcePos.OutOfBounds(tx.Range()) {
		return 0
	}
	b, ok := tx.World().blockLoaded(sourcePos)
	if !ok {
		return 0
	}
	if source, ok := b.(RedstonePowerSource); ok && e.redstoneWeaklyPowersBlocks(b) {
		return ClampRedstonePower(e.redstonePower(source, sourcePos, tx, face.Opposite()))
	}
	return 0
}

func (e *redstoneEngine) conductedWeakPower(pos cube.Pos, tx *Tx) int {
	power := 0
	for _, face := range cube.Faces() {
		power = max(power, e.conductedWeakPowerFrom(pos, tx, face))
	}
	return power
}

func (e *redstoneEngine) conductedWeakPowerFrom(pos cube.Pos, tx *Tx, face cube.Face) int {
	conductorPos := pos.Side(face)
	if conductorPos.OutOfBounds(tx.Range()) {
		return 0
	}
	conductor, ok := tx.World().blockLoaded(conductorPos)
	if !ok || !redstoneStrongPowerConductor(conductorPos, conductor, tx, face.Opposite()) {
		return 0
	}
	return e.weakBlockPower(conductorPos, tx)
}

func (e *redstoneEngine) conductedActivationPower(pos cube.Pos, tx *Tx) int {
	return max(e.conductedStrongPower(pos, tx), e.conductedWeakPower(pos, tx))
}

func (e *redstoneEngine) conductedActivationPowerFrom(pos cube.Pos, tx *Tx, face cube.Face) int {
	return max(e.conductedStrongPowerFrom(pos, tx, face), e.conductedWeakPowerFrom(pos, tx, face))
}

func (e *redstoneEngine) conductivePowerTo(pos cube.Pos, tx *Tx) int {
	b, ok := tx.World().blockLoaded(pos)
	if !ok || !RedstoneFullPowerConductor(pos, b, tx) {
		return 0
	}
	return max(e.strongPower(pos, tx), e.weakBlockPower(pos, tx))
}

func (e *redstoneEngine) acceptsDirectSourcePower(pos cube.Pos, tx *Tx) bool {
	b, ok := tx.World().blockLoaded(pos)
	if !ok {
		return true
	}
	if isRedstoneRelevant(b) {
		return true
	}
	return !RedstoneFullPowerConductor(pos, b, tx)
}

func (e *redstoneEngine) acceptsWeakConductedPower(pos cube.Pos, tx *Tx) bool {
	b, ok := tx.World().blockLoaded(pos)
	if !ok {
		return true
	}
	_, relayer := b.(RedstonePowerRelayer)
	return !relayer
}

func (e *redstoneEngine) redstoneWeaklyPowersBlocks(b Block) bool {
	weakBlockPowerer, ok := b.(RedstoneWeakBlockPowerer)
	return ok && weakBlockPowerer.RedstoneWeaklyPowersBlocks()
}

func (e *redstoneEngine) sourcePower(pos cube.Pos, tx *Tx) int {
	b, ok := tx.World().blockLoaded(pos)
	if !ok {
		return 0
	}
	source, ok := b.(RedstonePowerSource)
	if !ok {
		return 0
	}
	power := 0
	for _, face := range cube.Faces() {
		power = max(power, ClampRedstonePower(e.redstonePower(source, pos, tx, face)))
	}
	return power
}

func (e *redstoneEngine) graphPower(tx *Tx, graph redstoneGraph) []int {
	powers := make([]int, len(graph.nodes))
	if len(graph.nodes) == 0 {
		return powers
	}

	index := make(map[cube.Pos]int, len(graph.nodes))
	sources := make([]RedstonePowerSource, len(graph.nodes))
	relayers := make([]RedstonePowerRelayer, len(graph.nodes))
	edges := make([][]redstoneEdge, len(graph.nodes))
	for i, node := range graph.nodes {
		index[node.pos] = i
		if b, ok := tx.World().blockLoaded(node.pos); ok {
			sources[i], _ = b.(RedstonePowerSource)
			relayers[i], _ = b.(RedstonePowerRelayer)
		}
	}
	for _, edge := range graph.edges {
		edges[edge.from] = append(edges[edge.from], edge)
	}

	queue := make([]int, 0, len(graph.nodes))
	push := func(i, power int) {
		power = ClampRedstonePower(power)
		if power <= powers[i] {
			return
		}
		powers[i] = power
		queue = append(queue, i)
	}

	for i, node := range graph.nodes {
		if node.sink && relayers[i] == nil {
			push(i, e.conductedActivationPower(node.pos, tx))
			continue
		}
		push(i, e.conductedStrongPower(node.pos, tx))
	}

	for i, source := range sources {

		if source == nil || relayers[i] != nil {
			continue
		}
		pos := graph.nodes[i].pos
		for _, face := range cube.Faces() {
			j, ok := index[pos.Side(face)]
			if !ok {
				continue
			}
			power := ClampRedstonePower(e.redstonePower(source, pos, tx, face))
			push(j, power)
		}
	}

	for head := 0; head < len(queue); head++ {
		i := queue[head]
		if relayers[i] == nil {
			continue
		}
		for _, edge := range edges[i] {
			push(edge.to, powers[i]-edge.weight)
		}
	}
	return powers
}

func (e *redstoneEngine) powerTo(pos cube.Pos, tx *Tx) int {
	power := 0
	for _, face := range cube.Faces() {
		power = max(power, e.powerFrom(pos, tx, face))
	}
	if e.acceptsWeakConductedPower(pos, tx) {
		power = max(power, e.conductedActivationPower(pos, tx))
	} else {
		power = max(power, e.conductedStrongPower(pos, tx))
	}
	return ClampRedstonePower(power)
}

func (e *redstoneEngine) powerFrom(pos cube.Pos, tx *Tx, face cube.Face) int {
	power := e.conductedStrongPowerFrom(pos, tx, face)
	if e.acceptsWeakConductedPower(pos, tx) {
		power = e.conductedActivationPowerFrom(pos, tx, face)
	}
	type step struct {
		pos   cube.Pos
		from  cube.Face
		loss  int
		depth int
	}
	queue := []step{{pos: pos.Side(face), from: face.Opposite(), loss: 0, depth: 0}}
	seen := make(map[cube.Pos]int, 16)
	for len(queue) != 0 {
		s := queue[0]
		queue = queue[1:]
		if s.pos.OutOfBounds(tx.Range()) || s.loss >= 15 || s.depth >= 15 {
			continue
		}
		if s.pos == pos {
			continue
		}
		if loss, ok := seen[s.pos]; ok && loss <= s.loss {
			continue
		}
		seen[s.pos] = s.loss

		b, ok := tx.World().blockLoaded(s.pos)
		if !ok {
			continue
		}
		relayer, isRelayer := b.(RedstonePowerRelayer)

		if source, ok := b.(RedstonePowerSource); ok && !isRelayer && e.acceptsDirectSourcePower(pos, tx) {
			power = max(power, ClampRedstonePower(e.redstonePower(source, s.pos, tx, s.from)-s.loss))
		}
		if !isRelayer {
			continue
		}
		for _, next := range e.redstoneRelayerNeighbourPositions(tx, s.pos, b) {
			to := redstoneStepFace(s.pos, next)
			if to == s.from {
				continue
			}
			nextBlock, ok := tx.World().blockLoaded(next)
			if !ok {
				continue
			}
			loss := s.loss
			if _, nextRelayer := nextBlock.(RedstonePowerRelayer); nextRelayer {
				loss += max(relayer.RedstoneSignalLoss(s.pos, tx), 1)
			}
			if loss <= 15 {
				queue = append(queue, step{pos: next, from: to.Opposite(), loss: loss, depth: s.depth + 1})
			}
		}
	}
	return ClampRedstonePower(power)
}

func (e *redstoneEngine) torchBurnoutStatus(pos cube.Pos, currentTick int64) (burnedOut, recoverable bool) {
	data, ok := e.pruneTorchBurnout(pos, currentTick)
	if !ok {
		return false, true
	}
	if !data.burnedOut {
		return false, true
	}
	return true, len(data.offTicks) < redstoneTorchBurnoutThreshold
}

func (e *redstoneEngine) recordTorchTurnOff(pos cube.Pos, currentTick int64) bool {
	if e.torchBurnout == nil {
		e.torchBurnout = make(map[cube.Pos]redstoneTorchBurnout)
	}
	data, _ := e.pruneTorchBurnout(pos, currentTick)
	data.offTicks = append(data.offTicks, currentTick)
	if len(data.offTicks) >= redstoneTorchBurnoutThreshold {
		data.burnedOut = true
	}
	e.torchBurnout[pos] = data
	return data.burnedOut
}

func (e *redstoneEngine) clearTorchBurnout(pos cube.Pos) {
	if e == nil {
		return
	}
	delete(e.torchBurnout, pos)
}

func (e *redstoneEngine) markTorchSelfTriggered(pos cube.Pos) {
	if e == nil {
		return
	}
	if e.torchBurnout == nil {
		e.torchBurnout = make(map[cube.Pos]redstoneTorchBurnout)
	}
	data := e.torchBurnout[pos]
	data.pendingSelfTriggered = true
	e.torchBurnout[pos] = data
}

func (e *redstoneEngine) consumeTorchSelfTriggered(pos cube.Pos) bool {
	if e == nil || e.torchBurnout == nil {
		return false
	}
	data, ok := e.torchBurnout[pos]
	if !ok {
		return false
	}
	selfTriggered := data.pendingSelfTriggered
	data.pendingSelfTriggered = false
	if len(data.offTicks) == 0 && !data.burnedOut {
		delete(e.torchBurnout, pos)
	} else {
		e.torchBurnout[pos] = data
	}
	return selfTriggered
}

func (e *redstoneEngine) pruneTorchBurnout(pos cube.Pos, currentTick int64) (redstoneTorchBurnout, bool) {
	if e == nil || e.torchBurnout == nil {
		return redstoneTorchBurnout{}, false
	}
	data, ok := e.torchBurnout[pos]
	if !ok {
		return redstoneTorchBurnout{}, false
	}
	data.offTicks = slices.DeleteFunc(data.offTicks, func(tick int64) bool {
		return currentTick-tick >= redstoneTorchBurnoutWindowTicks
	})
	if len(data.offTicks) == 0 && !data.burnedOut && !data.pendingSelfTriggered {
		delete(e.torchBurnout, pos)
		return redstoneTorchBurnout{}, false
	}
	e.torchBurnout[pos] = data
	return data, true
}

func (e *redstoneEngine) redstonePower(source RedstonePowerSource, pos cube.Pos, tx *Tx, face cube.Face) int {
	if power, ok := e.suppressedSources[pos]; ok {
		return ClampRedstonePower(power)
	}
	if _, ok := e.evaluating[pos]; ok {
		return 0
	}
	e.evaluating[pos] = struct{}{}
	defer delete(e.evaluating, pos)
	return source.RedstonePower(pos, tx, face)
}

func (e *redstoneEngine) redstoneUpdateAllowed(tx *Tx, update RedstoneUpdate) bool {
	ctx := tx.Event()
	tx.World().Handler().HandleRedstoneUpdate(ctx, update)
	return !ctx.Cancelled()
}

type redstoneLightDiffuser interface {
	LightDiffusionLevel() uint8
}

func redstoneStrongPowerConductor(pos cube.Pos, b Block, tx *Tx, face cube.Face) bool {
	if !b.Model().FaceSolid(pos, face, tx) {
		return false
	}
	if _, ok := b.(RedstoneNonConductive); ok {
		return false
	}
	if diffuser, ok := b.(redstoneLightDiffuser); ok && diffuser.LightDiffusionLevel() == 0 {
		return false
	}
	return true
}

func RedstoneFullPowerConductor(pos cube.Pos, b Block, tx *Tx) bool {
	for _, face := range cube.Faces() {
		if !redstoneStrongPowerConductor(pos, b, tx, face) {
			return false
		}
	}
	return true
}

func (e *redstoneEngine) compileEdges(tx *Tx, nodes []redstoneNode) []redstoneEdge {
	index := make(map[cube.Pos]int, len(nodes))
	for i, node := range nodes {
		index[node.pos] = i
	}
	edges := make([]redstoneEdge, 0, len(nodes))
	for i, node := range nodes {
		b, loaded := tx.World().blockLoaded(node.pos)
		if !loaded {
			continue
		}
		relayer, ok := b.(RedstonePowerRelayer)
		if !ok {
			continue
		}
		for _, neighbour := range e.redstoneRelayerNeighbourPositions(tx, node.pos, b) {
			j, ok := index[neighbour]
			if !ok {
				continue
			}
			weight := 0
			if neighbourBlock, ok := tx.World().blockLoaded(neighbour); ok {
				if _, neighbourRelayer := neighbourBlock.(RedstonePowerRelayer); neighbourRelayer {
					weight = max(relayer.RedstoneSignalLoss(node.pos, tx), 1)
				}
			}
			edges = append(edges, redstoneEdge{from: i, to: j, weight: weight})
		}
	}
	slices.SortFunc(edges, compareRedstoneEdge)
	return edges
}

func (e *redstoneEngine) redstoneRelayerNeighbourPositions(tx *Tx, pos cube.Pos, b Block) []cube.Pos {
	if neighbourer, ok := b.(RedstonePowerRelayerNeighbourer); ok {
		neighbours := slices.Clone(neighbourer.RedstoneRelayerNeighbours(pos, tx))
		slices.SortFunc(neighbours, compareBlockPos)
		return neighbours
	}
	neighbours := make([]cube.Pos, 0, len(cube.Faces()))
	e.redstoneRelayerNeighbours(tx, pos, func(neighbour cube.Pos) {
		neighbours = append(neighbours, neighbour)
	})
	slices.SortFunc(neighbours, compareBlockPos)
	return neighbours
}

func (e *redstoneEngine) redstoneRelayerConnectedPositions(tx *Tx, pos cube.Pos, b Block) []cube.Pos {
	neighbours := e.redstoneRelayerNeighbourPositions(tx, pos, b)
	seen := make(map[cube.Pos]struct{}, len(neighbours)+8)
	for _, neighbour := range neighbours {
		seen[neighbour] = struct{}{}
	}
	for _, candidate := range redstoneRelayerIncomingCandidates(pos, tx.Range()) {
		if _, ok := seen[candidate]; ok {
			continue
		}
		candidateBlock, ok := tx.World().blockLoaded(candidate)
		if !ok {
			continue
		}
		if _, ok := candidateBlock.(RedstonePowerRelayer); !ok {
			continue
		}
		if slices.Contains(e.redstoneRelayerNeighbourPositions(tx, candidate, candidateBlock), pos) {
			neighbours = append(neighbours, candidate)
			seen[candidate] = struct{}{}
		}
	}
	slices.SortFunc(neighbours, compareBlockPos)
	return neighbours
}

func redstoneRelayerIncomingCandidates(pos cube.Pos, r cube.Range) []cube.Pos {
	candidates := make([]cube.Pos, 0, 26)
	for x := pos[0] - 1; x <= pos[0]+1; x++ {
		for y := pos[1] - 1; y <= pos[1]+1; y++ {
			for z := pos[2] - 1; z <= pos[2]+1; z++ {
				candidate := cube.Pos{x, y, z}
				if candidate == pos || candidate.OutOfBounds(r) {
					continue
				}
				candidates = append(candidates, candidate)
			}
		}
	}
	return candidates
}

func (e *redstoneEngine) redstoneRelayerNeighbours(tx *Tx, pos cube.Pos, f func(cube.Pos)) {
	for _, face := range cube.Faces() {
		neighbour := pos.Side(face)
		if !neighbour.OutOfBounds(tx.Range()) {
			f(neighbour)
		}
	}
}

func redstoneStepFace(from, to cube.Pos) cube.Face {
	dx, dy, dz := to[0]-from[0], to[1]-from[1], to[2]-from[2]
	switch {
	case dy > 0:
		return cube.FaceUp
	case dy < 0:
		return cube.FaceDown
	case dx > 0:
		return cube.FaceEast
	case dx < 0:
		return cube.FaceWest
	case dz > 0:
		return cube.FaceSouth
	case dz < 0:
		return cube.FaceNorth
	default:
		return cube.FaceUp
	}
}

func classifyRedstoneBlock(b Block) (source, consumer, action, relayer bool) {
	_, source = b.(RedstonePowerSource)
	_, consumer = b.(RedstonePowerConsumer)
	_, action = b.(RedstonePowerAction)
	if !action {
		_, action = b.(RedstonePowerContextAction)
	}
	_, relayer = b.(RedstonePowerRelayer)
	return
}

func isRedstoneRelevant(b Block) bool {
	source, consumer, action, relayer := classifyRedstoneBlock(b)
	return source || consumer || action || relayer
}

func compareBlockPos(a, b cube.Pos) int {
	if a[1] != b[1] {
		return a[1] - b[1]
	}
	if a[2] != b[2] {
		return a[2] - b[2]
	}
	return a[0] - b[0]
}

func compareRedstoneEdge(a, b redstoneEdge) int {
	if a.from != b.from {
		return a.from - b.from
	}
	if a.to != b.to {
		return a.to - b.to
	}
	return a.weight - b.weight
}

func ClampRedstonePower(power int) int {
	if power < 0 {
		return 0
	}
	if power > 15 {
		return 15
	}
	return power
}
