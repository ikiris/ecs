package ecs

import (
	"math"
)

const (
	InvalidEntity Id = 0 // Represents the default entity Id, which is invalid
	firstEntity   Id = 1
	MaxEntity     Id = math.MaxUint32
)

// World is the main data-holder. You usually pass it to other functions to do things.
type World struct {
	nextId       Id
	minId, maxId Id // This is the range of Ids returned by NewId
	arch         map[Id]archetypeId
	engine       *archEngine
}

// Creates a new world
func NewWorld() *World {
	return &World{
		nextId: firstEntity + 1,
		minId:  firstEntity + 1,
		maxId:  MaxEntity,
		arch:   make(map[Id]archetypeId),
		engine: newArchEngine(),
	}
}

// Sets an range of Ids that the world will use when creating new Ids. Potentially helpful when you have multiple worlds and don't want their Id space to collide.
// Deprecated: This API is tentative. It may be better to just have the user create Ids as they see fit
func (w *World) SetIdRange(min, max Id) {
	if min <= firstEntity {
		panic("max must be greater than 1")
	}
	if max <= firstEntity {
		panic("max must be greater than 1")
	}
	if min > max {
		panic("min must be less than max!")
	}

	w.minId = min
	w.maxId = max
}

// Creates a new Id which can then be used to create an entity
func (w *World) NewId() Id {
	if w.nextId < w.minId {
		w.nextId = w.minId
	}

	id := w.nextId

	if w.nextId == w.maxId {
		w.nextId = w.minId
	} else {
		w.nextId++
	}
	return id
}

// func (w *World) Count(anything ...any) int {
// 	return w.engine.Count(anything...)
// }

// func (w *World) Print(amount int) {
// 	fmt.Println("--- World ---")
// 	fmt.Printf("nextId: %d\n", w.nextId)

// 	// max := amount
// 	// for id, archId := range w.arch {
// 	// 	fmt.Printf("id(%d) -> archId(%d)\n", id, archId)
// 	// 	max--; if max <= 0 { break }
// 	// }

// 	// w.engine.Print(amount)
// }

// // A debug function for describing the current state of memory allocations in the ECS
// func (w *World) DescribeMemory() {
// 	fmt.Println("--- World ---")
// 	fmt.Printf("nextId: %d\n", w.nextId)
// 	fmt.Printf("Active Ent Count: %d\n", len(w.arch))
// 	for archId, lookup := range w.engine.lookup {
// 		efficiency := 100 * (1.0 - float64(len(lookup.holes))/float64(len(lookup.id)))
// 		fmt.Printf("Lookup[%d] = {len(index)=%d, len(id)=%d, len(holes)=%d} | Efficiency=%.2f%%\n", archId, len(lookup.index), len(lookup.id), len(lookup.holes), efficiency)
// 	}
// }

// TODO - Note: This function is not safe inside Maps or view iteraions
// TODO - make this loop-safe by:
// 1. Read the entire entity into an entity object
// 2. Call loop-safe delete method on that ID (which tags it somehow to indicate it needs to be cleaned up)
// 3. Modify the entity object by removing the requested components
// 4. Write the entity object to the destination archetype
// 4.a If the destination archetype is currently locked/flagged to indicate we are looping over it then wait for the lock release before writing the entity
// 4.b When creating Maps and Views we need to lock each archId that needs to be processed. Notably this guarantees that all "Writes" to this archetypeId will be done AFTER the lambda has processed - Meaning that we won't execute the same entity twice.
// 4.b.i When creating a view I may need like a "Close" method or "end" or something otherwise I'm not sure how to unlock the archId for modification
// Question: Why not write directly to holes if possible?

// Writes components to the entity specified at id. This API can potentially break if you call it inside of a loop. Specifically, if you cause the archetype of the entity to change by writing a new component, then the loop may act in mysterious ways.
// Deprecated: This API is tentative, I might replace it with something similar to bevy commands to alleviate the above concern
func Write(world *World, id Id, comp ...Component) {
	world.Write(id, comp...)
}

func (world *World) Write(id Id, comp ...Component) {
	if len(comp) <= 0 { return } // Do nothing if there are no components

	archId, ok := world.arch[id]
	if ok {
		newarchetypeId := world.engine.rewriteArch(archId, id, comp...)
		world.arch[id] = newarchetypeId
	} else {
		// Id does not yet exist, we need to add it for the first time
		archId = world.engine.GetarchetypeId(comp...)
		world.arch[id] = archId

		// Write all components to that archetype
		// TODO - Push this inward for efficiency?
		for i := range comp {
			comp[i].write(world.engine, archId, id)
		}
	}
}

// Reads a specific component of the entity specified at id.
// Returns true if the entity was found and had that component, else returns false.
// Deprecated: This API is tentative, I'm trying to improve the QueryN construct so that it can capture this usecase.
func Read[T any](world *World, id Id) (T, bool) {
	var ret T
	archId, ok := world.arch[id]
	if !ok {
		return ret, false
	}

	return readArch[T](world.engine, archId, id)
}

// Reads a pointer to the component of the entity at the specified id.
// Returns true if the entity was found and had that component, else returns false.
// This pointer is short lived and can become invalid if any other entity changes in the world
// Deprecated: This API is tentative, I'm trying to improve the QueryN construct so that it can capture this usecase.
func ReadPtr[T any](world *World, id Id) *T {
	archId, ok := world.arch[id]
	if !ok {
		return nil
	}

	return readPtrArch[T](world.engine, archId, id)
}

// This is safe for maps and loops
// 1. This deletes the high level id -> archId lookup
// 2. This creates a "hole" in the archetype list
// Returns true if the entity was deleted, else returns false if the entity does not exist (or was already deleted)

// Deletes the entire entity specified by the id
// This can be called inside maps and loops, it will delete the entity immediately.
// Returns true if the entity exists and was actually deleted, else returns false
func Delete(world *World, id Id) bool {
	archId, ok := world.arch[id]
	if !ok {
		return false
	}

	delete(world.arch, id)

	world.engine.TagForDeletion(archId, id)
	// Note: This was the old, more direct way, but isn't loop safe
	// - world.engine.DeleteAll(archId, id)
	return true
}

// Returns true if the entity exists in the world else it returns false
func (world *World) Exists(id Id) bool {
	_, ok := world.arch[id]
	return ok
}
