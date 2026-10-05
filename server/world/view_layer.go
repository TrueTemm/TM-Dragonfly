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
	"sync"
)

type layer struct {
	nameTag           *string
	alwaysShowNameTag *bool
	scoreTag          *string
	visibility        VisibilityLevel
}

type ViewLayerUpdater interface {
	ViewLayerEntityChanged(entity Entity)
}

type viewLayerViewer interface {
	ViewLayer() *ViewLayer
}

type ViewLayer struct {
	mu       sync.RWMutex
	entities map[*EntityHandle]layer
	updater  ViewLayerUpdater
}

func NewViewLayer(updater ViewLayerUpdater) *ViewLayer {
	return &ViewLayer{
		entities: map[*EntityHandle]layer{},
		updater:  updater,
	}
}

func (v *ViewLayer) Entities() []*EntityHandle {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return slices.Collect(maps.Keys(v.entities))
}

func (v *ViewLayer) ViewNameTag(entity Entity, nameTag string) {
	v.update(entity, func(l *layer) {
		l.nameTag = &nameTag
	})
}

func (v *ViewLayer) ViewPublicNameTag(entity Entity) {
	v.update(entity, func(l *layer) {
		l.nameTag = nil
	})
}

func (v *ViewLayer) NameTag(entity Entity) (string, bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	nameTag := v.entities[entity.H()].nameTag
	if nameTag == nil {
		return "", false
	}
	return *nameTag, true
}

func (v *ViewLayer) ViewAlwaysShowNameTag(entity Entity, alwaysShow bool) {
	v.update(entity, func(l *layer) {
		l.alwaysShowNameTag = &alwaysShow
	})
}

func (v *ViewLayer) ViewPublicAlwaysShowNameTag(entity Entity) {
	v.update(entity, func(l *layer) {
		l.alwaysShowNameTag = nil
	})
}

func (v *ViewLayer) AlwaysShowNameTag(entity Entity) (bool, bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	alwaysShow := v.entities[entity.H()].alwaysShowNameTag
	if alwaysShow == nil {
		return false, false
	}
	return *alwaysShow, true
}

func (v *ViewLayer) ViewScoreTag(entity Entity, scoreTag string) {
	v.update(entity, func(l *layer) {
		l.scoreTag = &scoreTag
	})
}

func (v *ViewLayer) ViewPublicScoreTag(entity Entity) {
	v.update(entity, func(l *layer) {
		l.scoreTag = nil
	})
}

func (v *ViewLayer) ScoreTag(entity Entity) (string, bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	scoreTag := v.entities[entity.H()].scoreTag
	if scoreTag == nil {
		return "", false
	}
	return *scoreTag, true
}

func (v *ViewLayer) ViewVisibility(entity Entity, level VisibilityLevel) {
	v.update(entity, func(l *layer) {
		l.visibility = level
	})
}

func (v *ViewLayer) Visibility(entity Entity) VisibilityLevel {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return v.entities[entity.H()].visibility
}

func (v *ViewLayer) Remove(entity Entity) {
	if v.remove(entity) {
		v.refresh(entity)
	}
}

func (v *ViewLayer) remove(entity Entity) bool {
	handle := entity.H()

	v.mu.Lock()
	_, ok := v.entities[handle]
	delete(v.entities, handle)
	v.mu.Unlock()
	return ok
}

func (v *ViewLayer) update(entity Entity, update func(*layer)) {
	handle := entity.H()

	v.mu.Lock()
	l := v.entities[handle]
	update(&l)
	if l.empty() {
		delete(v.entities, handle)
	} else {
		v.entities[handle] = l
	}
	v.mu.Unlock()

	v.refresh(entity)
}

func (v *ViewLayer) Close() error {
	v.mu.Lock()
	defer v.mu.Unlock()

	clear(v.entities)
	return nil
}

func (l layer) empty() bool {
	return l.nameTag == nil && l.alwaysShowNameTag == nil && l.scoreTag == nil && l.visibility == PublicVisibility()
}

func (v *ViewLayer) refresh(entity Entity) {
	if v.updater != nil {
		v.updater.ViewLayerEntityChanged(entity)
	}
}
