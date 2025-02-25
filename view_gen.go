package ecs

// --------------------------------------------------------------------------------
// - View 1
// --------------------------------------------------------------------------------

// Represents a view of data in a specific world. Provides access to the components specified in the generic block
type View1[A any] struct {
	world  *World
	filter filterList

	storageA componentSliceStorage[A]
}

// Creates a View for the specified world with the specified component filters.
func Query1[A any](world *World, filters ...Filter) *View1[A] {

	storageA := getStorage[A](world.engine)

	var AA A

	comps := []componentId{

		name(AA),
	}
	filterList := newFilterList(comps, filters...)
	filterList.regenerate(world)

	v := &View1[A]{
		world:  world,
		filter: filterList,

		storageA: storageA,
	}
	return v
}

// Reads a pointer to the underlying component at the specified id.
// Read will return even if the specified id doesn't match the filter list
// Read will return the value if it exists, else returns nil.
// If you execute any ecs.Write(...) or ecs.Delete(...) this pointer may become invalid.
func (v *View1[A]) Read(id Id) *A {
	if id == InvalidEntity {
		return nil
	}

	archId, ok := v.world.arch[id]
	if !ok {
		return nil
	}
	lookup, ok := v.world.engine.lookup[archId]
	if !ok {
		panic("LookupList is missing!")
	}
	index, ok := lookup.index[id]
	if !ok {
		return nil
	}

	var retA *A

	sliceA, ok := v.storageA.slice[archId]
	if ok {
		retA = &sliceA.comp[index]
	}

	return retA
}

// Maps the lambda function across every entity which matched the specified filters.
func (v *View1[A]) MapId(lambda func(id Id, a *A)) {
	v.filter.regenerate(v.world)

	var sliceA *componentSlice[A]
	var compA []A
	var retA *A

	for _, archId := range v.filter.archIds {

		sliceA, _ = v.storageA.slice[archId]

		lookup, ok := v.world.engine.lookup[archId]
		if !ok {
			panic("LookupList is missing!")
		}
		ids := lookup.id

		// TODO - this flattened version causes a mild performance hit. But the other one combinatorially explodes. I also cant get BCE to work with it. See option 2 for higher performance.

		compA = nil
		if sliceA != nil {
			compA = sliceA.comp
		}

		retA = nil
		for idx := range ids {
			if ids[idx] == InvalidEntity {
				continue
			} // Skip if its a hole

			if compA != nil {
				retA = &compA[idx]
			}
			lambda(ids[idx], retA)
		}

		// 	// Option 2 - This is faster but has a combinatorial explosion problem
		// 	if compA == nil && compB == nil {
		// 		return
		// 	} else if compA != nil && compB == nil {
		// 		if len(ids) != len(compA) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], &compA[i], nil)
		// 		}
		// 	} else if compA == nil && compB != nil {
		// 		if len(ids) != len(compB) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], nil, &compB[i])
		// 		}
		// 	} else if compA != nil && compB != nil {
		// 		if len(ids) != len(compA) || len(ids) != len(compB) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], &compA[i], &compB[i])
		// 		}
		// 	}
	}

	// Original - doesn't handle optional
	// for _, archId := range v.filter.archIds {
	// 	aSlice, ok := v.storageA.slice[archId]
	// 	if !ok { continue }
	// 	bSlice, ok := v.storageB.slice[archId]
	// 	if !ok { continue }

	// 	lookup, ok := v.world.engine.lookup[archId]
	// 	if !ok { panic("LookupList is missing!") }

	// 	ids := lookup.id
	// 	aComp := aSlice.comp
	// 	bComp := bSlice.comp
	// 	if len(ids) != len(aComp) || len(ids) != len(bComp) {
	// 		panic("ERROR - Bounds don't match")
	// 	}
	// 	for i := range ids {
	// 		if ids[i] == InvalidEntity { continue }
	// 		lambda(ids[i], &aComp[i], &bComp[i])
	// 	}
	// }
}

// Deprecated: This API is a tentative alternative way to map
func (v *View1[A]) MapSlices(lambda func(id []Id, a []A)) {
	v.filter.regenerate(v.world)

	id := make([][]Id, 0)

	sliceListA := make([][]A, 0)

	for _, archId := range v.filter.archIds {

		sliceA, ok := v.storageA.slice[archId]
		if !ok {
			continue
		}

		lookup, ok := v.world.engine.lookup[archId]
		if !ok {
			panic("LookupList is missing!")
		}

		id = append(id, lookup.id)

		sliceListA = append(sliceListA, sliceA.comp)
	}

	for idx := range id {
		lambda(id[idx],
			sliceListA[idx],
		)
	}
}

// --------------------------------------------------------------------------------
// - View 2
// --------------------------------------------------------------------------------

// Represents a view of data in a specific world. Provides access to the components specified in the generic block
type View2[A, B any] struct {
	world  *World
	filter filterList

	storageA componentSliceStorage[A]
	storageB componentSliceStorage[B]
}

// Creates a View for the specified world with the specified component filters.
func Query2[A, B any](world *World, filters ...Filter) *View2[A, B] {

	storageA := getStorage[A](world.engine)
	storageB := getStorage[B](world.engine)

	var AA A
	var BB B

	comps := []componentId{

		name(AA),
		name(BB),
	}
	filterList := newFilterList(comps, filters...)
	filterList.regenerate(world)

	v := &View2[A, B]{
		world:  world,
		filter: filterList,

		storageA: storageA,
		storageB: storageB,
	}
	return v
}

// Reads a pointer to the underlying component at the specified id.
// Read will return even if the specified id doesn't match the filter list
// Read will return the value if it exists, else returns nil.
// If you execute any ecs.Write(...) or ecs.Delete(...) this pointer may become invalid.
func (v *View2[A, B]) Read(id Id) (*A, *B) {
	if id == InvalidEntity {
		return nil, nil
	}

	archId, ok := v.world.arch[id]
	if !ok {
		return nil, nil
	}
	lookup, ok := v.world.engine.lookup[archId]
	if !ok {
		panic("LookupList is missing!")
	}
	index, ok := lookup.index[id]
	if !ok {
		return nil, nil
	}

	var retA *A
	var retB *B

	sliceA, ok := v.storageA.slice[archId]
	if ok {
		retA = &sliceA.comp[index]
	}
	sliceB, ok := v.storageB.slice[archId]
	if ok {
		retB = &sliceB.comp[index]
	}

	return retA, retB
}

// Maps the lambda function across every entity which matched the specified filters.
func (v *View2[A, B]) MapId(lambda func(id Id, a *A, b *B)) {
	v.filter.regenerate(v.world)

	var sliceA *componentSlice[A]
	var compA []A
	var retA *A

	var sliceB *componentSlice[B]
	var compB []B
	var retB *B

	for _, archId := range v.filter.archIds {

		sliceA, _ = v.storageA.slice[archId]
		sliceB, _ = v.storageB.slice[archId]

		lookup, ok := v.world.engine.lookup[archId]
		if !ok {
			panic("LookupList is missing!")
		}
		ids := lookup.id

		// TODO - this flattened version causes a mild performance hit. But the other one combinatorially explodes. I also cant get BCE to work with it. See option 2 for higher performance.

		compA = nil
		if sliceA != nil {
			compA = sliceA.comp
		}
		compB = nil
		if sliceB != nil {
			compB = sliceB.comp
		}

		retA = nil
		retB = nil
		for idx := range ids {
			if ids[idx] == InvalidEntity {
				continue
			} // Skip if its a hole

			if compA != nil {
				retA = &compA[idx]
			}
			if compB != nil {
				retB = &compB[idx]
			}
			lambda(ids[idx], retA, retB)
		}

		// 	// Option 2 - This is faster but has a combinatorial explosion problem
		// 	if compA == nil && compB == nil {
		// 		return
		// 	} else if compA != nil && compB == nil {
		// 		if len(ids) != len(compA) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], &compA[i], nil)
		// 		}
		// 	} else if compA == nil && compB != nil {
		// 		if len(ids) != len(compB) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], nil, &compB[i])
		// 		}
		// 	} else if compA != nil && compB != nil {
		// 		if len(ids) != len(compA) || len(ids) != len(compB) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], &compA[i], &compB[i])
		// 		}
		// 	}
	}

	// Original - doesn't handle optional
	// for _, archId := range v.filter.archIds {
	// 	aSlice, ok := v.storageA.slice[archId]
	// 	if !ok { continue }
	// 	bSlice, ok := v.storageB.slice[archId]
	// 	if !ok { continue }

	// 	lookup, ok := v.world.engine.lookup[archId]
	// 	if !ok { panic("LookupList is missing!") }

	// 	ids := lookup.id
	// 	aComp := aSlice.comp
	// 	bComp := bSlice.comp
	// 	if len(ids) != len(aComp) || len(ids) != len(bComp) {
	// 		panic("ERROR - Bounds don't match")
	// 	}
	// 	for i := range ids {
	// 		if ids[i] == InvalidEntity { continue }
	// 		lambda(ids[i], &aComp[i], &bComp[i])
	// 	}
	// }
}

// Deprecated: This API is a tentative alternative way to map
func (v *View2[A, B]) MapSlices(lambda func(id []Id, a []A, b []B)) {
	v.filter.regenerate(v.world)

	id := make([][]Id, 0)

	sliceListA := make([][]A, 0)
	sliceListB := make([][]B, 0)

	for _, archId := range v.filter.archIds {

		sliceA, ok := v.storageA.slice[archId]
		if !ok {
			continue
		}
		sliceB, ok := v.storageB.slice[archId]
		if !ok {
			continue
		}

		lookup, ok := v.world.engine.lookup[archId]
		if !ok {
			panic("LookupList is missing!")
		}

		id = append(id, lookup.id)

		sliceListA = append(sliceListA, sliceA.comp)
		sliceListB = append(sliceListB, sliceB.comp)
	}

	for idx := range id {
		lambda(id[idx],
			sliceListA[idx], sliceListB[idx],
		)
	}
}

// --------------------------------------------------------------------------------
// - View 3
// --------------------------------------------------------------------------------

// Represents a view of data in a specific world. Provides access to the components specified in the generic block
type View3[A, B, C any] struct {
	world  *World
	filter filterList

	storageA componentSliceStorage[A]
	storageB componentSliceStorage[B]
	storageC componentSliceStorage[C]
}

// Creates a View for the specified world with the specified component filters.
func Query3[A, B, C any](world *World, filters ...Filter) *View3[A, B, C] {

	storageA := getStorage[A](world.engine)
	storageB := getStorage[B](world.engine)
	storageC := getStorage[C](world.engine)

	var AA A
	var BB B
	var CC C

	comps := []componentId{

		name(AA),
		name(BB),
		name(CC),
	}
	filterList := newFilterList(comps, filters...)
	filterList.regenerate(world)

	v := &View3[A, B, C]{
		world:  world,
		filter: filterList,

		storageA: storageA,
		storageB: storageB,
		storageC: storageC,
	}
	return v
}

// Reads a pointer to the underlying component at the specified id.
// Read will return even if the specified id doesn't match the filter list
// Read will return the value if it exists, else returns nil.
// If you execute any ecs.Write(...) or ecs.Delete(...) this pointer may become invalid.
func (v *View3[A, B, C]) Read(id Id) (*A, *B, *C) {
	if id == InvalidEntity {
		return nil, nil, nil
	}

	archId, ok := v.world.arch[id]
	if !ok {
		return nil, nil, nil
	}
	lookup, ok := v.world.engine.lookup[archId]
	if !ok {
		panic("LookupList is missing!")
	}
	index, ok := lookup.index[id]
	if !ok {
		return nil, nil, nil
	}

	var retA *A
	var retB *B
	var retC *C

	sliceA, ok := v.storageA.slice[archId]
	if ok {
		retA = &sliceA.comp[index]
	}
	sliceB, ok := v.storageB.slice[archId]
	if ok {
		retB = &sliceB.comp[index]
	}
	sliceC, ok := v.storageC.slice[archId]
	if ok {
		retC = &sliceC.comp[index]
	}

	return retA, retB, retC
}

// Maps the lambda function across every entity which matched the specified filters.
func (v *View3[A, B, C]) MapId(lambda func(id Id, a *A, b *B, c *C)) {
	v.filter.regenerate(v.world)

	var sliceA *componentSlice[A]
	var compA []A
	var retA *A

	var sliceB *componentSlice[B]
	var compB []B
	var retB *B

	var sliceC *componentSlice[C]
	var compC []C
	var retC *C

	for _, archId := range v.filter.archIds {

		sliceA, _ = v.storageA.slice[archId]
		sliceB, _ = v.storageB.slice[archId]
		sliceC, _ = v.storageC.slice[archId]

		lookup, ok := v.world.engine.lookup[archId]
		if !ok {
			panic("LookupList is missing!")
		}
		ids := lookup.id

		// TODO - this flattened version causes a mild performance hit. But the other one combinatorially explodes. I also cant get BCE to work with it. See option 2 for higher performance.

		compA = nil
		if sliceA != nil {
			compA = sliceA.comp
		}
		compB = nil
		if sliceB != nil {
			compB = sliceB.comp
		}
		compC = nil
		if sliceC != nil {
			compC = sliceC.comp
		}

		retA = nil
		retB = nil
		retC = nil
		for idx := range ids {
			if ids[idx] == InvalidEntity {
				continue
			} // Skip if its a hole

			if compA != nil {
				retA = &compA[idx]
			}
			if compB != nil {
				retB = &compB[idx]
			}
			if compC != nil {
				retC = &compC[idx]
			}
			lambda(ids[idx], retA, retB, retC)
		}

		// 	// Option 2 - This is faster but has a combinatorial explosion problem
		// 	if compA == nil && compB == nil {
		// 		return
		// 	} else if compA != nil && compB == nil {
		// 		if len(ids) != len(compA) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], &compA[i], nil)
		// 		}
		// 	} else if compA == nil && compB != nil {
		// 		if len(ids) != len(compB) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], nil, &compB[i])
		// 		}
		// 	} else if compA != nil && compB != nil {
		// 		if len(ids) != len(compA) || len(ids) != len(compB) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], &compA[i], &compB[i])
		// 		}
		// 	}
	}

	// Original - doesn't handle optional
	// for _, archId := range v.filter.archIds {
	// 	aSlice, ok := v.storageA.slice[archId]
	// 	if !ok { continue }
	// 	bSlice, ok := v.storageB.slice[archId]
	// 	if !ok { continue }

	// 	lookup, ok := v.world.engine.lookup[archId]
	// 	if !ok { panic("LookupList is missing!") }

	// 	ids := lookup.id
	// 	aComp := aSlice.comp
	// 	bComp := bSlice.comp
	// 	if len(ids) != len(aComp) || len(ids) != len(bComp) {
	// 		panic("ERROR - Bounds don't match")
	// 	}
	// 	for i := range ids {
	// 		if ids[i] == InvalidEntity { continue }
	// 		lambda(ids[i], &aComp[i], &bComp[i])
	// 	}
	// }
}

// Deprecated: This API is a tentative alternative way to map
func (v *View3[A, B, C]) MapSlices(lambda func(id []Id, a []A, b []B, c []C)) {
	v.filter.regenerate(v.world)

	id := make([][]Id, 0)

	sliceListA := make([][]A, 0)
	sliceListB := make([][]B, 0)
	sliceListC := make([][]C, 0)

	for _, archId := range v.filter.archIds {

		sliceA, ok := v.storageA.slice[archId]
		if !ok {
			continue
		}
		sliceB, ok := v.storageB.slice[archId]
		if !ok {
			continue
		}
		sliceC, ok := v.storageC.slice[archId]
		if !ok {
			continue
		}

		lookup, ok := v.world.engine.lookup[archId]
		if !ok {
			panic("LookupList is missing!")
		}

		id = append(id, lookup.id)

		sliceListA = append(sliceListA, sliceA.comp)
		sliceListB = append(sliceListB, sliceB.comp)
		sliceListC = append(sliceListC, sliceC.comp)
	}

	for idx := range id {
		lambda(id[idx],
			sliceListA[idx], sliceListB[idx], sliceListC[idx],
		)
	}
}

// --------------------------------------------------------------------------------
// - View 4
// --------------------------------------------------------------------------------

// Represents a view of data in a specific world. Provides access to the components specified in the generic block
type View4[A, B, C, D any] struct {
	world  *World
	filter filterList

	storageA componentSliceStorage[A]
	storageB componentSliceStorage[B]
	storageC componentSliceStorage[C]
	storageD componentSliceStorage[D]
}

// Creates a View for the specified world with the specified component filters.
func Query4[A, B, C, D any](world *World, filters ...Filter) *View4[A, B, C, D] {

	storageA := getStorage[A](world.engine)
	storageB := getStorage[B](world.engine)
	storageC := getStorage[C](world.engine)
	storageD := getStorage[D](world.engine)

	var AA A
	var BB B
	var CC C
	var DD D

	comps := []componentId{

		name(AA),
		name(BB),
		name(CC),
		name(DD),
	}
	filterList := newFilterList(comps, filters...)
	filterList.regenerate(world)

	v := &View4[A, B, C, D]{
		world:  world,
		filter: filterList,

		storageA: storageA,
		storageB: storageB,
		storageC: storageC,
		storageD: storageD,
	}
	return v
}

// Reads a pointer to the underlying component at the specified id.
// Read will return even if the specified id doesn't match the filter list
// Read will return the value if it exists, else returns nil.
// If you execute any ecs.Write(...) or ecs.Delete(...) this pointer may become invalid.
func (v *View4[A, B, C, D]) Read(id Id) (*A, *B, *C, *D) {
	if id == InvalidEntity {
		return nil, nil, nil, nil
	}

	archId, ok := v.world.arch[id]
	if !ok {
		return nil, nil, nil, nil
	}
	lookup, ok := v.world.engine.lookup[archId]
	if !ok {
		panic("LookupList is missing!")
	}
	index, ok := lookup.index[id]
	if !ok {
		return nil, nil, nil, nil
	}

	var retA *A
	var retB *B
	var retC *C
	var retD *D

	sliceA, ok := v.storageA.slice[archId]
	if ok {
		retA = &sliceA.comp[index]
	}
	sliceB, ok := v.storageB.slice[archId]
	if ok {
		retB = &sliceB.comp[index]
	}
	sliceC, ok := v.storageC.slice[archId]
	if ok {
		retC = &sliceC.comp[index]
	}
	sliceD, ok := v.storageD.slice[archId]
	if ok {
		retD = &sliceD.comp[index]
	}

	return retA, retB, retC, retD
}

// Maps the lambda function across every entity which matched the specified filters.
func (v *View4[A, B, C, D]) MapId(lambda func(id Id, a *A, b *B, c *C, d *D)) {
	v.filter.regenerate(v.world)

	var sliceA *componentSlice[A]
	var compA []A
	var retA *A

	var sliceB *componentSlice[B]
	var compB []B
	var retB *B

	var sliceC *componentSlice[C]
	var compC []C
	var retC *C

	var sliceD *componentSlice[D]
	var compD []D
	var retD *D

	for _, archId := range v.filter.archIds {

		sliceA, _ = v.storageA.slice[archId]
		sliceB, _ = v.storageB.slice[archId]
		sliceC, _ = v.storageC.slice[archId]
		sliceD, _ = v.storageD.slice[archId]

		lookup, ok := v.world.engine.lookup[archId]
		if !ok {
			panic("LookupList is missing!")
		}
		ids := lookup.id

		// TODO - this flattened version causes a mild performance hit. But the other one combinatorially explodes. I also cant get BCE to work with it. See option 2 for higher performance.

		compA = nil
		if sliceA != nil {
			compA = sliceA.comp
		}
		compB = nil
		if sliceB != nil {
			compB = sliceB.comp
		}
		compC = nil
		if sliceC != nil {
			compC = sliceC.comp
		}
		compD = nil
		if sliceD != nil {
			compD = sliceD.comp
		}

		retA = nil
		retB = nil
		retC = nil
		retD = nil
		for idx := range ids {
			if ids[idx] == InvalidEntity {
				continue
			} // Skip if its a hole

			if compA != nil {
				retA = &compA[idx]
			}
			if compB != nil {
				retB = &compB[idx]
			}
			if compC != nil {
				retC = &compC[idx]
			}
			if compD != nil {
				retD = &compD[idx]
			}
			lambda(ids[idx], retA, retB, retC, retD)
		}

		// 	// Option 2 - This is faster but has a combinatorial explosion problem
		// 	if compA == nil && compB == nil {
		// 		return
		// 	} else if compA != nil && compB == nil {
		// 		if len(ids) != len(compA) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], &compA[i], nil)
		// 		}
		// 	} else if compA == nil && compB != nil {
		// 		if len(ids) != len(compB) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], nil, &compB[i])
		// 		}
		// 	} else if compA != nil && compB != nil {
		// 		if len(ids) != len(compA) || len(ids) != len(compB) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], &compA[i], &compB[i])
		// 		}
		// 	}
	}

	// Original - doesn't handle optional
	// for _, archId := range v.filter.archIds {
	// 	aSlice, ok := v.storageA.slice[archId]
	// 	if !ok { continue }
	// 	bSlice, ok := v.storageB.slice[archId]
	// 	if !ok { continue }

	// 	lookup, ok := v.world.engine.lookup[archId]
	// 	if !ok { panic("LookupList is missing!") }

	// 	ids := lookup.id
	// 	aComp := aSlice.comp
	// 	bComp := bSlice.comp
	// 	if len(ids) != len(aComp) || len(ids) != len(bComp) {
	// 		panic("ERROR - Bounds don't match")
	// 	}
	// 	for i := range ids {
	// 		if ids[i] == InvalidEntity { continue }
	// 		lambda(ids[i], &aComp[i], &bComp[i])
	// 	}
	// }
}

// Deprecated: This API is a tentative alternative way to map
func (v *View4[A, B, C, D]) MapSlices(lambda func(id []Id, a []A, b []B, c []C, d []D)) {
	v.filter.regenerate(v.world)

	id := make([][]Id, 0)

	sliceListA := make([][]A, 0)
	sliceListB := make([][]B, 0)
	sliceListC := make([][]C, 0)
	sliceListD := make([][]D, 0)

	for _, archId := range v.filter.archIds {

		sliceA, ok := v.storageA.slice[archId]
		if !ok {
			continue
		}
		sliceB, ok := v.storageB.slice[archId]
		if !ok {
			continue
		}
		sliceC, ok := v.storageC.slice[archId]
		if !ok {
			continue
		}
		sliceD, ok := v.storageD.slice[archId]
		if !ok {
			continue
		}

		lookup, ok := v.world.engine.lookup[archId]
		if !ok {
			panic("LookupList is missing!")
		}

		id = append(id, lookup.id)

		sliceListA = append(sliceListA, sliceA.comp)
		sliceListB = append(sliceListB, sliceB.comp)
		sliceListC = append(sliceListC, sliceC.comp)
		sliceListD = append(sliceListD, sliceD.comp)
	}

	for idx := range id {
		lambda(id[idx],
			sliceListA[idx], sliceListB[idx], sliceListC[idx], sliceListD[idx],
		)
	}
}

// --------------------------------------------------------------------------------
// - View 5
// --------------------------------------------------------------------------------

// Represents a view of data in a specific world. Provides access to the components specified in the generic block
type View5[A, B, C, D, E any] struct {
	world  *World
	filter filterList

	storageA componentSliceStorage[A]
	storageB componentSliceStorage[B]
	storageC componentSliceStorage[C]
	storageD componentSliceStorage[D]
	storageE componentSliceStorage[E]
}

// Creates a View for the specified world with the specified component filters.
func Query5[A, B, C, D, E any](world *World, filters ...Filter) *View5[A, B, C, D, E] {

	storageA := getStorage[A](world.engine)
	storageB := getStorage[B](world.engine)
	storageC := getStorage[C](world.engine)
	storageD := getStorage[D](world.engine)
	storageE := getStorage[E](world.engine)

	var AA A
	var BB B
	var CC C
	var DD D
	var EE E

	comps := []componentId{

		name(AA),
		name(BB),
		name(CC),
		name(DD),
		name(EE),
	}
	filterList := newFilterList(comps, filters...)
	filterList.regenerate(world)

	v := &View5[A, B, C, D, E]{
		world:  world,
		filter: filterList,

		storageA: storageA,
		storageB: storageB,
		storageC: storageC,
		storageD: storageD,
		storageE: storageE,
	}
	return v
}

// Reads a pointer to the underlying component at the specified id.
// Read will return even if the specified id doesn't match the filter list
// Read will return the value if it exists, else returns nil.
// If you execute any ecs.Write(...) or ecs.Delete(...) this pointer may become invalid.
func (v *View5[A, B, C, D, E]) Read(id Id) (*A, *B, *C, *D, *E) {
	if id == InvalidEntity {
		return nil, nil, nil, nil, nil
	}

	archId, ok := v.world.arch[id]
	if !ok {
		return nil, nil, nil, nil, nil
	}
	lookup, ok := v.world.engine.lookup[archId]
	if !ok {
		panic("LookupList is missing!")
	}
	index, ok := lookup.index[id]
	if !ok {
		return nil, nil, nil, nil, nil
	}

	var retA *A
	var retB *B
	var retC *C
	var retD *D
	var retE *E

	sliceA, ok := v.storageA.slice[archId]
	if ok {
		retA = &sliceA.comp[index]
	}
	sliceB, ok := v.storageB.slice[archId]
	if ok {
		retB = &sliceB.comp[index]
	}
	sliceC, ok := v.storageC.slice[archId]
	if ok {
		retC = &sliceC.comp[index]
	}
	sliceD, ok := v.storageD.slice[archId]
	if ok {
		retD = &sliceD.comp[index]
	}
	sliceE, ok := v.storageE.slice[archId]
	if ok {
		retE = &sliceE.comp[index]
	}

	return retA, retB, retC, retD, retE
}

// Maps the lambda function across every entity which matched the specified filters.
func (v *View5[A, B, C, D, E]) MapId(lambda func(id Id, a *A, b *B, c *C, d *D, e *E)) {
	v.filter.regenerate(v.world)

	var sliceA *componentSlice[A]
	var compA []A
	var retA *A

	var sliceB *componentSlice[B]
	var compB []B
	var retB *B

	var sliceC *componentSlice[C]
	var compC []C
	var retC *C

	var sliceD *componentSlice[D]
	var compD []D
	var retD *D

	var sliceE *componentSlice[E]
	var compE []E
	var retE *E

	for _, archId := range v.filter.archIds {

		sliceA, _ = v.storageA.slice[archId]
		sliceB, _ = v.storageB.slice[archId]
		sliceC, _ = v.storageC.slice[archId]
		sliceD, _ = v.storageD.slice[archId]
		sliceE, _ = v.storageE.slice[archId]

		lookup, ok := v.world.engine.lookup[archId]
		if !ok {
			panic("LookupList is missing!")
		}
		ids := lookup.id

		// TODO - this flattened version causes a mild performance hit. But the other one combinatorially explodes. I also cant get BCE to work with it. See option 2 for higher performance.

		compA = nil
		if sliceA != nil {
			compA = sliceA.comp
		}
		compB = nil
		if sliceB != nil {
			compB = sliceB.comp
		}
		compC = nil
		if sliceC != nil {
			compC = sliceC.comp
		}
		compD = nil
		if sliceD != nil {
			compD = sliceD.comp
		}
		compE = nil
		if sliceE != nil {
			compE = sliceE.comp
		}

		retA = nil
		retB = nil
		retC = nil
		retD = nil
		retE = nil
		for idx := range ids {
			if ids[idx] == InvalidEntity {
				continue
			} // Skip if its a hole

			if compA != nil {
				retA = &compA[idx]
			}
			if compB != nil {
				retB = &compB[idx]
			}
			if compC != nil {
				retC = &compC[idx]
			}
			if compD != nil {
				retD = &compD[idx]
			}
			if compE != nil {
				retE = &compE[idx]
			}
			lambda(ids[idx], retA, retB, retC, retD, retE)
		}

		// 	// Option 2 - This is faster but has a combinatorial explosion problem
		// 	if compA == nil && compB == nil {
		// 		return
		// 	} else if compA != nil && compB == nil {
		// 		if len(ids) != len(compA) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], &compA[i], nil)
		// 		}
		// 	} else if compA == nil && compB != nil {
		// 		if len(ids) != len(compB) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], nil, &compB[i])
		// 		}
		// 	} else if compA != nil && compB != nil {
		// 		if len(ids) != len(compA) || len(ids) != len(compB) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], &compA[i], &compB[i])
		// 		}
		// 	}
	}

	// Original - doesn't handle optional
	// for _, archId := range v.filter.archIds {
	// 	aSlice, ok := v.storageA.slice[archId]
	// 	if !ok { continue }
	// 	bSlice, ok := v.storageB.slice[archId]
	// 	if !ok { continue }

	// 	lookup, ok := v.world.engine.lookup[archId]
	// 	if !ok { panic("LookupList is missing!") }

	// 	ids := lookup.id
	// 	aComp := aSlice.comp
	// 	bComp := bSlice.comp
	// 	if len(ids) != len(aComp) || len(ids) != len(bComp) {
	// 		panic("ERROR - Bounds don't match")
	// 	}
	// 	for i := range ids {
	// 		if ids[i] == InvalidEntity { continue }
	// 		lambda(ids[i], &aComp[i], &bComp[i])
	// 	}
	// }
}

// Deprecated: This API is a tentative alternative way to map
func (v *View5[A, B, C, D, E]) MapSlices(lambda func(id []Id, a []A, b []B, c []C, d []D, e []E)) {
	v.filter.regenerate(v.world)

	id := make([][]Id, 0)

	sliceListA := make([][]A, 0)
	sliceListB := make([][]B, 0)
	sliceListC := make([][]C, 0)
	sliceListD := make([][]D, 0)
	sliceListE := make([][]E, 0)

	for _, archId := range v.filter.archIds {

		sliceA, ok := v.storageA.slice[archId]
		if !ok {
			continue
		}
		sliceB, ok := v.storageB.slice[archId]
		if !ok {
			continue
		}
		sliceC, ok := v.storageC.slice[archId]
		if !ok {
			continue
		}
		sliceD, ok := v.storageD.slice[archId]
		if !ok {
			continue
		}
		sliceE, ok := v.storageE.slice[archId]
		if !ok {
			continue
		}

		lookup, ok := v.world.engine.lookup[archId]
		if !ok {
			panic("LookupList is missing!")
		}

		id = append(id, lookup.id)

		sliceListA = append(sliceListA, sliceA.comp)
		sliceListB = append(sliceListB, sliceB.comp)
		sliceListC = append(sliceListC, sliceC.comp)
		sliceListD = append(sliceListD, sliceD.comp)
		sliceListE = append(sliceListE, sliceE.comp)
	}

	for idx := range id {
		lambda(id[idx],
			sliceListA[idx], sliceListB[idx], sliceListC[idx], sliceListD[idx], sliceListE[idx],
		)
	}
}

// --------------------------------------------------------------------------------
// - View 6
// --------------------------------------------------------------------------------

// Represents a view of data in a specific world. Provides access to the components specified in the generic block
type View6[A, B, C, D, E, F any] struct {
	world  *World
	filter filterList

	storageA componentSliceStorage[A]
	storageB componentSliceStorage[B]
	storageC componentSliceStorage[C]
	storageD componentSliceStorage[D]
	storageE componentSliceStorage[E]
	storageF componentSliceStorage[F]
}

// Creates a View for the specified world with the specified component filters.
func Query6[A, B, C, D, E, F any](world *World, filters ...Filter) *View6[A, B, C, D, E, F] {

	storageA := getStorage[A](world.engine)
	storageB := getStorage[B](world.engine)
	storageC := getStorage[C](world.engine)
	storageD := getStorage[D](world.engine)
	storageE := getStorage[E](world.engine)
	storageF := getStorage[F](world.engine)

	var AA A
	var BB B
	var CC C
	var DD D
	var EE E
	var FF F

	comps := []componentId{

		name(AA),
		name(BB),
		name(CC),
		name(DD),
		name(EE),
		name(FF),
	}
	filterList := newFilterList(comps, filters...)
	filterList.regenerate(world)

	v := &View6[A, B, C, D, E, F]{
		world:  world,
		filter: filterList,

		storageA: storageA,
		storageB: storageB,
		storageC: storageC,
		storageD: storageD,
		storageE: storageE,
		storageF: storageF,
	}
	return v
}

// Reads a pointer to the underlying component at the specified id.
// Read will return even if the specified id doesn't match the filter list
// Read will return the value if it exists, else returns nil.
// If you execute any ecs.Write(...) or ecs.Delete(...) this pointer may become invalid.
func (v *View6[A, B, C, D, E, F]) Read(id Id) (*A, *B, *C, *D, *E, *F) {
	if id == InvalidEntity {
		return nil, nil, nil, nil, nil, nil
	}

	archId, ok := v.world.arch[id]
	if !ok {
		return nil, nil, nil, nil, nil, nil
	}
	lookup, ok := v.world.engine.lookup[archId]
	if !ok {
		panic("LookupList is missing!")
	}
	index, ok := lookup.index[id]
	if !ok {
		return nil, nil, nil, nil, nil, nil
	}

	var retA *A
	var retB *B
	var retC *C
	var retD *D
	var retE *E
	var retF *F

	sliceA, ok := v.storageA.slice[archId]
	if ok {
		retA = &sliceA.comp[index]
	}
	sliceB, ok := v.storageB.slice[archId]
	if ok {
		retB = &sliceB.comp[index]
	}
	sliceC, ok := v.storageC.slice[archId]
	if ok {
		retC = &sliceC.comp[index]
	}
	sliceD, ok := v.storageD.slice[archId]
	if ok {
		retD = &sliceD.comp[index]
	}
	sliceE, ok := v.storageE.slice[archId]
	if ok {
		retE = &sliceE.comp[index]
	}
	sliceF, ok := v.storageF.slice[archId]
	if ok {
		retF = &sliceF.comp[index]
	}

	return retA, retB, retC, retD, retE, retF
}

// Maps the lambda function across every entity which matched the specified filters.
func (v *View6[A, B, C, D, E, F]) MapId(lambda func(id Id, a *A, b *B, c *C, d *D, e *E, f *F)) {
	v.filter.regenerate(v.world)

	var sliceA *componentSlice[A]
	var compA []A
	var retA *A

	var sliceB *componentSlice[B]
	var compB []B
	var retB *B

	var sliceC *componentSlice[C]
	var compC []C
	var retC *C

	var sliceD *componentSlice[D]
	var compD []D
	var retD *D

	var sliceE *componentSlice[E]
	var compE []E
	var retE *E

	var sliceF *componentSlice[F]
	var compF []F
	var retF *F

	for _, archId := range v.filter.archIds {

		sliceA, _ = v.storageA.slice[archId]
		sliceB, _ = v.storageB.slice[archId]
		sliceC, _ = v.storageC.slice[archId]
		sliceD, _ = v.storageD.slice[archId]
		sliceE, _ = v.storageE.slice[archId]
		sliceF, _ = v.storageF.slice[archId]

		lookup, ok := v.world.engine.lookup[archId]
		if !ok {
			panic("LookupList is missing!")
		}
		ids := lookup.id

		// TODO - this flattened version causes a mild performance hit. But the other one combinatorially explodes. I also cant get BCE to work with it. See option 2 for higher performance.

		compA = nil
		if sliceA != nil {
			compA = sliceA.comp
		}
		compB = nil
		if sliceB != nil {
			compB = sliceB.comp
		}
		compC = nil
		if sliceC != nil {
			compC = sliceC.comp
		}
		compD = nil
		if sliceD != nil {
			compD = sliceD.comp
		}
		compE = nil
		if sliceE != nil {
			compE = sliceE.comp
		}
		compF = nil
		if sliceF != nil {
			compF = sliceF.comp
		}

		retA = nil
		retB = nil
		retC = nil
		retD = nil
		retE = nil
		retF = nil
		for idx := range ids {
			if ids[idx] == InvalidEntity {
				continue
			} // Skip if its a hole

			if compA != nil {
				retA = &compA[idx]
			}
			if compB != nil {
				retB = &compB[idx]
			}
			if compC != nil {
				retC = &compC[idx]
			}
			if compD != nil {
				retD = &compD[idx]
			}
			if compE != nil {
				retE = &compE[idx]
			}
			if compF != nil {
				retF = &compF[idx]
			}
			lambda(ids[idx], retA, retB, retC, retD, retE, retF)
		}

		// 	// Option 2 - This is faster but has a combinatorial explosion problem
		// 	if compA == nil && compB == nil {
		// 		return
		// 	} else if compA != nil && compB == nil {
		// 		if len(ids) != len(compA) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], &compA[i], nil)
		// 		}
		// 	} else if compA == nil && compB != nil {
		// 		if len(ids) != len(compB) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], nil, &compB[i])
		// 		}
		// 	} else if compA != nil && compB != nil {
		// 		if len(ids) != len(compA) || len(ids) != len(compB) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], &compA[i], &compB[i])
		// 		}
		// 	}
	}

	// Original - doesn't handle optional
	// for _, archId := range v.filter.archIds {
	// 	aSlice, ok := v.storageA.slice[archId]
	// 	if !ok { continue }
	// 	bSlice, ok := v.storageB.slice[archId]
	// 	if !ok { continue }

	// 	lookup, ok := v.world.engine.lookup[archId]
	// 	if !ok { panic("LookupList is missing!") }

	// 	ids := lookup.id
	// 	aComp := aSlice.comp
	// 	bComp := bSlice.comp
	// 	if len(ids) != len(aComp) || len(ids) != len(bComp) {
	// 		panic("ERROR - Bounds don't match")
	// 	}
	// 	for i := range ids {
	// 		if ids[i] == InvalidEntity { continue }
	// 		lambda(ids[i], &aComp[i], &bComp[i])
	// 	}
	// }
}

// Deprecated: This API is a tentative alternative way to map
func (v *View6[A, B, C, D, E, F]) MapSlices(lambda func(id []Id, a []A, b []B, c []C, d []D, e []E, f []F)) {
	v.filter.regenerate(v.world)

	id := make([][]Id, 0)

	sliceListA := make([][]A, 0)
	sliceListB := make([][]B, 0)
	sliceListC := make([][]C, 0)
	sliceListD := make([][]D, 0)
	sliceListE := make([][]E, 0)
	sliceListF := make([][]F, 0)

	for _, archId := range v.filter.archIds {

		sliceA, ok := v.storageA.slice[archId]
		if !ok {
			continue
		}
		sliceB, ok := v.storageB.slice[archId]
		if !ok {
			continue
		}
		sliceC, ok := v.storageC.slice[archId]
		if !ok {
			continue
		}
		sliceD, ok := v.storageD.slice[archId]
		if !ok {
			continue
		}
		sliceE, ok := v.storageE.slice[archId]
		if !ok {
			continue
		}
		sliceF, ok := v.storageF.slice[archId]
		if !ok {
			continue
		}

		lookup, ok := v.world.engine.lookup[archId]
		if !ok {
			panic("LookupList is missing!")
		}

		id = append(id, lookup.id)

		sliceListA = append(sliceListA, sliceA.comp)
		sliceListB = append(sliceListB, sliceB.comp)
		sliceListC = append(sliceListC, sliceC.comp)
		sliceListD = append(sliceListD, sliceD.comp)
		sliceListE = append(sliceListE, sliceE.comp)
		sliceListF = append(sliceListF, sliceF.comp)
	}

	for idx := range id {
		lambda(id[idx],
			sliceListA[idx], sliceListB[idx], sliceListC[idx], sliceListD[idx], sliceListE[idx], sliceListF[idx],
		)
	}
}

// --------------------------------------------------------------------------------
// - View 7
// --------------------------------------------------------------------------------

// Represents a view of data in a specific world. Provides access to the components specified in the generic block
type View7[A, B, C, D, E, F, G any] struct {
	world  *World
	filter filterList

	storageA componentSliceStorage[A]
	storageB componentSliceStorage[B]
	storageC componentSliceStorage[C]
	storageD componentSliceStorage[D]
	storageE componentSliceStorage[E]
	storageF componentSliceStorage[F]
	storageG componentSliceStorage[G]
}

// Creates a View for the specified world with the specified component filters.
func Query7[A, B, C, D, E, F, G any](world *World, filters ...Filter) *View7[A, B, C, D, E, F, G] {

	storageA := getStorage[A](world.engine)
	storageB := getStorage[B](world.engine)
	storageC := getStorage[C](world.engine)
	storageD := getStorage[D](world.engine)
	storageE := getStorage[E](world.engine)
	storageF := getStorage[F](world.engine)
	storageG := getStorage[G](world.engine)

	var AA A
	var BB B
	var CC C
	var DD D
	var EE E
	var FF F
	var GG G

	comps := []componentId{

		name(AA),
		name(BB),
		name(CC),
		name(DD),
		name(EE),
		name(FF),
		name(GG),
	}
	filterList := newFilterList(comps, filters...)
	filterList.regenerate(world)

	v := &View7[A, B, C, D, E, F, G]{
		world:  world,
		filter: filterList,

		storageA: storageA,
		storageB: storageB,
		storageC: storageC,
		storageD: storageD,
		storageE: storageE,
		storageF: storageF,
		storageG: storageG,
	}
	return v
}

// Reads a pointer to the underlying component at the specified id.
// Read will return even if the specified id doesn't match the filter list
// Read will return the value if it exists, else returns nil.
// If you execute any ecs.Write(...) or ecs.Delete(...) this pointer may become invalid.
func (v *View7[A, B, C, D, E, F, G]) Read(id Id) (*A, *B, *C, *D, *E, *F, *G) {
	if id == InvalidEntity {
		return nil, nil, nil, nil, nil, nil, nil
	}

	archId, ok := v.world.arch[id]
	if !ok {
		return nil, nil, nil, nil, nil, nil, nil
	}
	lookup, ok := v.world.engine.lookup[archId]
	if !ok {
		panic("LookupList is missing!")
	}
	index, ok := lookup.index[id]
	if !ok {
		return nil, nil, nil, nil, nil, nil, nil
	}

	var retA *A
	var retB *B
	var retC *C
	var retD *D
	var retE *E
	var retF *F
	var retG *G

	sliceA, ok := v.storageA.slice[archId]
	if ok {
		retA = &sliceA.comp[index]
	}
	sliceB, ok := v.storageB.slice[archId]
	if ok {
		retB = &sliceB.comp[index]
	}
	sliceC, ok := v.storageC.slice[archId]
	if ok {
		retC = &sliceC.comp[index]
	}
	sliceD, ok := v.storageD.slice[archId]
	if ok {
		retD = &sliceD.comp[index]
	}
	sliceE, ok := v.storageE.slice[archId]
	if ok {
		retE = &sliceE.comp[index]
	}
	sliceF, ok := v.storageF.slice[archId]
	if ok {
		retF = &sliceF.comp[index]
	}
	sliceG, ok := v.storageG.slice[archId]
	if ok {
		retG = &sliceG.comp[index]
	}

	return retA, retB, retC, retD, retE, retF, retG
}

// Maps the lambda function across every entity which matched the specified filters.
func (v *View7[A, B, C, D, E, F, G]) MapId(lambda func(id Id, a *A, b *B, c *C, d *D, e *E, f *F, g *G)) {
	v.filter.regenerate(v.world)

	var sliceA *componentSlice[A]
	var compA []A
	var retA *A

	var sliceB *componentSlice[B]
	var compB []B
	var retB *B

	var sliceC *componentSlice[C]
	var compC []C
	var retC *C

	var sliceD *componentSlice[D]
	var compD []D
	var retD *D

	var sliceE *componentSlice[E]
	var compE []E
	var retE *E

	var sliceF *componentSlice[F]
	var compF []F
	var retF *F

	var sliceG *componentSlice[G]
	var compG []G
	var retG *G

	for _, archId := range v.filter.archIds {

		sliceA, _ = v.storageA.slice[archId]
		sliceB, _ = v.storageB.slice[archId]
		sliceC, _ = v.storageC.slice[archId]
		sliceD, _ = v.storageD.slice[archId]
		sliceE, _ = v.storageE.slice[archId]
		sliceF, _ = v.storageF.slice[archId]
		sliceG, _ = v.storageG.slice[archId]

		lookup, ok := v.world.engine.lookup[archId]
		if !ok {
			panic("LookupList is missing!")
		}
		ids := lookup.id

		// TODO - this flattened version causes a mild performance hit. But the other one combinatorially explodes. I also cant get BCE to work with it. See option 2 for higher performance.

		compA = nil
		if sliceA != nil {
			compA = sliceA.comp
		}
		compB = nil
		if sliceB != nil {
			compB = sliceB.comp
		}
		compC = nil
		if sliceC != nil {
			compC = sliceC.comp
		}
		compD = nil
		if sliceD != nil {
			compD = sliceD.comp
		}
		compE = nil
		if sliceE != nil {
			compE = sliceE.comp
		}
		compF = nil
		if sliceF != nil {
			compF = sliceF.comp
		}
		compG = nil
		if sliceG != nil {
			compG = sliceG.comp
		}

		retA = nil
		retB = nil
		retC = nil
		retD = nil
		retE = nil
		retF = nil
		retG = nil
		for idx := range ids {
			if ids[idx] == InvalidEntity {
				continue
			} // Skip if its a hole

			if compA != nil {
				retA = &compA[idx]
			}
			if compB != nil {
				retB = &compB[idx]
			}
			if compC != nil {
				retC = &compC[idx]
			}
			if compD != nil {
				retD = &compD[idx]
			}
			if compE != nil {
				retE = &compE[idx]
			}
			if compF != nil {
				retF = &compF[idx]
			}
			if compG != nil {
				retG = &compG[idx]
			}
			lambda(ids[idx], retA, retB, retC, retD, retE, retF, retG)
		}

		// 	// Option 2 - This is faster but has a combinatorial explosion problem
		// 	if compA == nil && compB == nil {
		// 		return
		// 	} else if compA != nil && compB == nil {
		// 		if len(ids) != len(compA) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], &compA[i], nil)
		// 		}
		// 	} else if compA == nil && compB != nil {
		// 		if len(ids) != len(compB) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], nil, &compB[i])
		// 		}
		// 	} else if compA != nil && compB != nil {
		// 		if len(ids) != len(compA) || len(ids) != len(compB) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], &compA[i], &compB[i])
		// 		}
		// 	}
	}

	// Original - doesn't handle optional
	// for _, archId := range v.filter.archIds {
	// 	aSlice, ok := v.storageA.slice[archId]
	// 	if !ok { continue }
	// 	bSlice, ok := v.storageB.slice[archId]
	// 	if !ok { continue }

	// 	lookup, ok := v.world.engine.lookup[archId]
	// 	if !ok { panic("LookupList is missing!") }

	// 	ids := lookup.id
	// 	aComp := aSlice.comp
	// 	bComp := bSlice.comp
	// 	if len(ids) != len(aComp) || len(ids) != len(bComp) {
	// 		panic("ERROR - Bounds don't match")
	// 	}
	// 	for i := range ids {
	// 		if ids[i] == InvalidEntity { continue }
	// 		lambda(ids[i], &aComp[i], &bComp[i])
	// 	}
	// }
}

// Deprecated: This API is a tentative alternative way to map
func (v *View7[A, B, C, D, E, F, G]) MapSlices(lambda func(id []Id, a []A, b []B, c []C, d []D, e []E, f []F, g []G)) {
	v.filter.regenerate(v.world)

	id := make([][]Id, 0)

	sliceListA := make([][]A, 0)
	sliceListB := make([][]B, 0)
	sliceListC := make([][]C, 0)
	sliceListD := make([][]D, 0)
	sliceListE := make([][]E, 0)
	sliceListF := make([][]F, 0)
	sliceListG := make([][]G, 0)

	for _, archId := range v.filter.archIds {

		sliceA, ok := v.storageA.slice[archId]
		if !ok {
			continue
		}
		sliceB, ok := v.storageB.slice[archId]
		if !ok {
			continue
		}
		sliceC, ok := v.storageC.slice[archId]
		if !ok {
			continue
		}
		sliceD, ok := v.storageD.slice[archId]
		if !ok {
			continue
		}
		sliceE, ok := v.storageE.slice[archId]
		if !ok {
			continue
		}
		sliceF, ok := v.storageF.slice[archId]
		if !ok {
			continue
		}
		sliceG, ok := v.storageG.slice[archId]
		if !ok {
			continue
		}

		lookup, ok := v.world.engine.lookup[archId]
		if !ok {
			panic("LookupList is missing!")
		}

		id = append(id, lookup.id)

		sliceListA = append(sliceListA, sliceA.comp)
		sliceListB = append(sliceListB, sliceB.comp)
		sliceListC = append(sliceListC, sliceC.comp)
		sliceListD = append(sliceListD, sliceD.comp)
		sliceListE = append(sliceListE, sliceE.comp)
		sliceListF = append(sliceListF, sliceF.comp)
		sliceListG = append(sliceListG, sliceG.comp)
	}

	for idx := range id {
		lambda(id[idx],
			sliceListA[idx], sliceListB[idx], sliceListC[idx], sliceListD[idx], sliceListE[idx], sliceListF[idx], sliceListG[idx],
		)
	}
}

// --------------------------------------------------------------------------------
// - View 8
// --------------------------------------------------------------------------------

// Represents a view of data in a specific world. Provides access to the components specified in the generic block
type View8[A, B, C, D, E, F, G, H any] struct {
	world  *World
	filter filterList

	storageA componentSliceStorage[A]
	storageB componentSliceStorage[B]
	storageC componentSliceStorage[C]
	storageD componentSliceStorage[D]
	storageE componentSliceStorage[E]
	storageF componentSliceStorage[F]
	storageG componentSliceStorage[G]
	storageH componentSliceStorage[H]
}

// Creates a View for the specified world with the specified component filters.
func Query8[A, B, C, D, E, F, G, H any](world *World, filters ...Filter) *View8[A, B, C, D, E, F, G, H] {

	storageA := getStorage[A](world.engine)
	storageB := getStorage[B](world.engine)
	storageC := getStorage[C](world.engine)
	storageD := getStorage[D](world.engine)
	storageE := getStorage[E](world.engine)
	storageF := getStorage[F](world.engine)
	storageG := getStorage[G](world.engine)
	storageH := getStorage[H](world.engine)

	var AA A
	var BB B
	var CC C
	var DD D
	var EE E
	var FF F
	var GG G
	var HH H

	comps := []componentId{

		name(AA),
		name(BB),
		name(CC),
		name(DD),
		name(EE),
		name(FF),
		name(GG),
		name(HH),
	}
	filterList := newFilterList(comps, filters...)
	filterList.regenerate(world)

	v := &View8[A, B, C, D, E, F, G, H]{
		world:  world,
		filter: filterList,

		storageA: storageA,
		storageB: storageB,
		storageC: storageC,
		storageD: storageD,
		storageE: storageE,
		storageF: storageF,
		storageG: storageG,
		storageH: storageH,
	}
	return v
}

// Reads a pointer to the underlying component at the specified id.
// Read will return even if the specified id doesn't match the filter list
// Read will return the value if it exists, else returns nil.
// If you execute any ecs.Write(...) or ecs.Delete(...) this pointer may become invalid.
func (v *View8[A, B, C, D, E, F, G, H]) Read(id Id) (*A, *B, *C, *D, *E, *F, *G, *H) {
	if id == InvalidEntity {
		return nil, nil, nil, nil, nil, nil, nil, nil
	}

	archId, ok := v.world.arch[id]
	if !ok {
		return nil, nil, nil, nil, nil, nil, nil, nil
	}
	lookup, ok := v.world.engine.lookup[archId]
	if !ok {
		panic("LookupList is missing!")
	}
	index, ok := lookup.index[id]
	if !ok {
		return nil, nil, nil, nil, nil, nil, nil, nil
	}

	var retA *A
	var retB *B
	var retC *C
	var retD *D
	var retE *E
	var retF *F
	var retG *G
	var retH *H

	sliceA, ok := v.storageA.slice[archId]
	if ok {
		retA = &sliceA.comp[index]
	}
	sliceB, ok := v.storageB.slice[archId]
	if ok {
		retB = &sliceB.comp[index]
	}
	sliceC, ok := v.storageC.slice[archId]
	if ok {
		retC = &sliceC.comp[index]
	}
	sliceD, ok := v.storageD.slice[archId]
	if ok {
		retD = &sliceD.comp[index]
	}
	sliceE, ok := v.storageE.slice[archId]
	if ok {
		retE = &sliceE.comp[index]
	}
	sliceF, ok := v.storageF.slice[archId]
	if ok {
		retF = &sliceF.comp[index]
	}
	sliceG, ok := v.storageG.slice[archId]
	if ok {
		retG = &sliceG.comp[index]
	}
	sliceH, ok := v.storageH.slice[archId]
	if ok {
		retH = &sliceH.comp[index]
	}

	return retA, retB, retC, retD, retE, retF, retG, retH
}

// Maps the lambda function across every entity which matched the specified filters.
func (v *View8[A, B, C, D, E, F, G, H]) MapId(lambda func(id Id, a *A, b *B, c *C, d *D, e *E, f *F, g *G, h *H)) {
	v.filter.regenerate(v.world)

	var sliceA *componentSlice[A]
	var compA []A
	var retA *A

	var sliceB *componentSlice[B]
	var compB []B
	var retB *B

	var sliceC *componentSlice[C]
	var compC []C
	var retC *C

	var sliceD *componentSlice[D]
	var compD []D
	var retD *D

	var sliceE *componentSlice[E]
	var compE []E
	var retE *E

	var sliceF *componentSlice[F]
	var compF []F
	var retF *F

	var sliceG *componentSlice[G]
	var compG []G
	var retG *G

	var sliceH *componentSlice[H]
	var compH []H
	var retH *H

	for _, archId := range v.filter.archIds {

		sliceA, _ = v.storageA.slice[archId]
		sliceB, _ = v.storageB.slice[archId]
		sliceC, _ = v.storageC.slice[archId]
		sliceD, _ = v.storageD.slice[archId]
		sliceE, _ = v.storageE.slice[archId]
		sliceF, _ = v.storageF.slice[archId]
		sliceG, _ = v.storageG.slice[archId]
		sliceH, _ = v.storageH.slice[archId]

		lookup, ok := v.world.engine.lookup[archId]
		if !ok {
			panic("LookupList is missing!")
		}
		ids := lookup.id

		// TODO - this flattened version causes a mild performance hit. But the other one combinatorially explodes. I also cant get BCE to work with it. See option 2 for higher performance.

		compA = nil
		if sliceA != nil {
			compA = sliceA.comp
		}
		compB = nil
		if sliceB != nil {
			compB = sliceB.comp
		}
		compC = nil
		if sliceC != nil {
			compC = sliceC.comp
		}
		compD = nil
		if sliceD != nil {
			compD = sliceD.comp
		}
		compE = nil
		if sliceE != nil {
			compE = sliceE.comp
		}
		compF = nil
		if sliceF != nil {
			compF = sliceF.comp
		}
		compG = nil
		if sliceG != nil {
			compG = sliceG.comp
		}
		compH = nil
		if sliceH != nil {
			compH = sliceH.comp
		}

		retA = nil
		retB = nil
		retC = nil
		retD = nil
		retE = nil
		retF = nil
		retG = nil
		retH = nil
		for idx := range ids {
			if ids[idx] == InvalidEntity {
				continue
			} // Skip if its a hole

			if compA != nil {
				retA = &compA[idx]
			}
			if compB != nil {
				retB = &compB[idx]
			}
			if compC != nil {
				retC = &compC[idx]
			}
			if compD != nil {
				retD = &compD[idx]
			}
			if compE != nil {
				retE = &compE[idx]
			}
			if compF != nil {
				retF = &compF[idx]
			}
			if compG != nil {
				retG = &compG[idx]
			}
			if compH != nil {
				retH = &compH[idx]
			}
			lambda(ids[idx], retA, retB, retC, retD, retE, retF, retG, retH)
		}

		// 	// Option 2 - This is faster but has a combinatorial explosion problem
		// 	if compA == nil && compB == nil {
		// 		return
		// 	} else if compA != nil && compB == nil {
		// 		if len(ids) != len(compA) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], &compA[i], nil)
		// 		}
		// 	} else if compA == nil && compB != nil {
		// 		if len(ids) != len(compB) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], nil, &compB[i])
		// 		}
		// 	} else if compA != nil && compB != nil {
		// 		if len(ids) != len(compA) || len(ids) != len(compB) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], &compA[i], &compB[i])
		// 		}
		// 	}
	}

	// Original - doesn't handle optional
	// for _, archId := range v.filter.archIds {
	// 	aSlice, ok := v.storageA.slice[archId]
	// 	if !ok { continue }
	// 	bSlice, ok := v.storageB.slice[archId]
	// 	if !ok { continue }

	// 	lookup, ok := v.world.engine.lookup[archId]
	// 	if !ok { panic("LookupList is missing!") }

	// 	ids := lookup.id
	// 	aComp := aSlice.comp
	// 	bComp := bSlice.comp
	// 	if len(ids) != len(aComp) || len(ids) != len(bComp) {
	// 		panic("ERROR - Bounds don't match")
	// 	}
	// 	for i := range ids {
	// 		if ids[i] == InvalidEntity { continue }
	// 		lambda(ids[i], &aComp[i], &bComp[i])
	// 	}
	// }
}

// Deprecated: This API is a tentative alternative way to map
func (v *View8[A, B, C, D, E, F, G, H]) MapSlices(lambda func(id []Id, a []A, b []B, c []C, d []D, e []E, f []F, g []G, h []H)) {
	v.filter.regenerate(v.world)

	id := make([][]Id, 0)

	sliceListA := make([][]A, 0)
	sliceListB := make([][]B, 0)
	sliceListC := make([][]C, 0)
	sliceListD := make([][]D, 0)
	sliceListE := make([][]E, 0)
	sliceListF := make([][]F, 0)
	sliceListG := make([][]G, 0)
	sliceListH := make([][]H, 0)

	for _, archId := range v.filter.archIds {

		sliceA, ok := v.storageA.slice[archId]
		if !ok {
			continue
		}
		sliceB, ok := v.storageB.slice[archId]
		if !ok {
			continue
		}
		sliceC, ok := v.storageC.slice[archId]
		if !ok {
			continue
		}
		sliceD, ok := v.storageD.slice[archId]
		if !ok {
			continue
		}
		sliceE, ok := v.storageE.slice[archId]
		if !ok {
			continue
		}
		sliceF, ok := v.storageF.slice[archId]
		if !ok {
			continue
		}
		sliceG, ok := v.storageG.slice[archId]
		if !ok {
			continue
		}
		sliceH, ok := v.storageH.slice[archId]
		if !ok {
			continue
		}

		lookup, ok := v.world.engine.lookup[archId]
		if !ok {
			panic("LookupList is missing!")
		}

		id = append(id, lookup.id)

		sliceListA = append(sliceListA, sliceA.comp)
		sliceListB = append(sliceListB, sliceB.comp)
		sliceListC = append(sliceListC, sliceC.comp)
		sliceListD = append(sliceListD, sliceD.comp)
		sliceListE = append(sliceListE, sliceE.comp)
		sliceListF = append(sliceListF, sliceF.comp)
		sliceListG = append(sliceListG, sliceG.comp)
		sliceListH = append(sliceListH, sliceH.comp)
	}

	for idx := range id {
		lambda(id[idx],
			sliceListA[idx], sliceListB[idx], sliceListC[idx], sliceListD[idx], sliceListE[idx], sliceListF[idx], sliceListG[idx], sliceListH[idx],
		)
	}
}

// --------------------------------------------------------------------------------
// - View 9
// --------------------------------------------------------------------------------

// Represents a view of data in a specific world. Provides access to the components specified in the generic block
type View9[A, B, C, D, E, F, G, H, I any] struct {
	world  *World
	filter filterList

	storageA componentSliceStorage[A]
	storageB componentSliceStorage[B]
	storageC componentSliceStorage[C]
	storageD componentSliceStorage[D]
	storageE componentSliceStorage[E]
	storageF componentSliceStorage[F]
	storageG componentSliceStorage[G]
	storageH componentSliceStorage[H]
	storageI componentSliceStorage[I]
}

// Creates a View for the specified world with the specified component filters.
func Query9[A, B, C, D, E, F, G, H, I any](world *World, filters ...Filter) *View9[A, B, C, D, E, F, G, H, I] {

	storageA := getStorage[A](world.engine)
	storageB := getStorage[B](world.engine)
	storageC := getStorage[C](world.engine)
	storageD := getStorage[D](world.engine)
	storageE := getStorage[E](world.engine)
	storageF := getStorage[F](world.engine)
	storageG := getStorage[G](world.engine)
	storageH := getStorage[H](world.engine)
	storageI := getStorage[I](world.engine)

	var AA A
	var BB B
	var CC C
	var DD D
	var EE E
	var FF F
	var GG G
	var HH H
	var II I

	comps := []componentId{

		name(AA),
		name(BB),
		name(CC),
		name(DD),
		name(EE),
		name(FF),
		name(GG),
		name(HH),
		name(II),
	}
	filterList := newFilterList(comps, filters...)
	filterList.regenerate(world)

	v := &View9[A, B, C, D, E, F, G, H, I]{
		world:  world,
		filter: filterList,

		storageA: storageA,
		storageB: storageB,
		storageC: storageC,
		storageD: storageD,
		storageE: storageE,
		storageF: storageF,
		storageG: storageG,
		storageH: storageH,
		storageI: storageI,
	}
	return v
}

// Reads a pointer to the underlying component at the specified id.
// Read will return even if the specified id doesn't match the filter list
// Read will return the value if it exists, else returns nil.
// If you execute any ecs.Write(...) or ecs.Delete(...) this pointer may become invalid.
func (v *View9[A, B, C, D, E, F, G, H, I]) Read(id Id) (*A, *B, *C, *D, *E, *F, *G, *H, *I) {
	if id == InvalidEntity {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil
	}

	archId, ok := v.world.arch[id]
	if !ok {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil
	}
	lookup, ok := v.world.engine.lookup[archId]
	if !ok {
		panic("LookupList is missing!")
	}
	index, ok := lookup.index[id]
	if !ok {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil
	}

	var retA *A
	var retB *B
	var retC *C
	var retD *D
	var retE *E
	var retF *F
	var retG *G
	var retH *H
	var retI *I

	sliceA, ok := v.storageA.slice[archId]
	if ok {
		retA = &sliceA.comp[index]
	}
	sliceB, ok := v.storageB.slice[archId]
	if ok {
		retB = &sliceB.comp[index]
	}
	sliceC, ok := v.storageC.slice[archId]
	if ok {
		retC = &sliceC.comp[index]
	}
	sliceD, ok := v.storageD.slice[archId]
	if ok {
		retD = &sliceD.comp[index]
	}
	sliceE, ok := v.storageE.slice[archId]
	if ok {
		retE = &sliceE.comp[index]
	}
	sliceF, ok := v.storageF.slice[archId]
	if ok {
		retF = &sliceF.comp[index]
	}
	sliceG, ok := v.storageG.slice[archId]
	if ok {
		retG = &sliceG.comp[index]
	}
	sliceH, ok := v.storageH.slice[archId]
	if ok {
		retH = &sliceH.comp[index]
	}
	sliceI, ok := v.storageI.slice[archId]
	if ok {
		retI = &sliceI.comp[index]
	}

	return retA, retB, retC, retD, retE, retF, retG, retH, retI
}

// Maps the lambda function across every entity which matched the specified filters.
func (v *View9[A, B, C, D, E, F, G, H, I]) MapId(lambda func(id Id, a *A, b *B, c *C, d *D, e *E, f *F, g *G, h *H, i *I)) {
	v.filter.regenerate(v.world)

	var sliceA *componentSlice[A]
	var compA []A
	var retA *A

	var sliceB *componentSlice[B]
	var compB []B
	var retB *B

	var sliceC *componentSlice[C]
	var compC []C
	var retC *C

	var sliceD *componentSlice[D]
	var compD []D
	var retD *D

	var sliceE *componentSlice[E]
	var compE []E
	var retE *E

	var sliceF *componentSlice[F]
	var compF []F
	var retF *F

	var sliceG *componentSlice[G]
	var compG []G
	var retG *G

	var sliceH *componentSlice[H]
	var compH []H
	var retH *H

	var sliceI *componentSlice[I]
	var compI []I
	var retI *I

	for _, archId := range v.filter.archIds {

		sliceA, _ = v.storageA.slice[archId]
		sliceB, _ = v.storageB.slice[archId]
		sliceC, _ = v.storageC.slice[archId]
		sliceD, _ = v.storageD.slice[archId]
		sliceE, _ = v.storageE.slice[archId]
		sliceF, _ = v.storageF.slice[archId]
		sliceG, _ = v.storageG.slice[archId]
		sliceH, _ = v.storageH.slice[archId]
		sliceI, _ = v.storageI.slice[archId]

		lookup, ok := v.world.engine.lookup[archId]
		if !ok {
			panic("LookupList is missing!")
		}
		ids := lookup.id

		// TODO - this flattened version causes a mild performance hit. But the other one combinatorially explodes. I also cant get BCE to work with it. See option 2 for higher performance.

		compA = nil
		if sliceA != nil {
			compA = sliceA.comp
		}
		compB = nil
		if sliceB != nil {
			compB = sliceB.comp
		}
		compC = nil
		if sliceC != nil {
			compC = sliceC.comp
		}
		compD = nil
		if sliceD != nil {
			compD = sliceD.comp
		}
		compE = nil
		if sliceE != nil {
			compE = sliceE.comp
		}
		compF = nil
		if sliceF != nil {
			compF = sliceF.comp
		}
		compG = nil
		if sliceG != nil {
			compG = sliceG.comp
		}
		compH = nil
		if sliceH != nil {
			compH = sliceH.comp
		}
		compI = nil
		if sliceI != nil {
			compI = sliceI.comp
		}

		retA = nil
		retB = nil
		retC = nil
		retD = nil
		retE = nil
		retF = nil
		retG = nil
		retH = nil
		retI = nil
		for idx := range ids {
			if ids[idx] == InvalidEntity {
				continue
			} // Skip if its a hole

			if compA != nil {
				retA = &compA[idx]
			}
			if compB != nil {
				retB = &compB[idx]
			}
			if compC != nil {
				retC = &compC[idx]
			}
			if compD != nil {
				retD = &compD[idx]
			}
			if compE != nil {
				retE = &compE[idx]
			}
			if compF != nil {
				retF = &compF[idx]
			}
			if compG != nil {
				retG = &compG[idx]
			}
			if compH != nil {
				retH = &compH[idx]
			}
			if compI != nil {
				retI = &compI[idx]
			}
			lambda(ids[idx], retA, retB, retC, retD, retE, retF, retG, retH, retI)
		}

		// 	// Option 2 - This is faster but has a combinatorial explosion problem
		// 	if compA == nil && compB == nil {
		// 		return
		// 	} else if compA != nil && compB == nil {
		// 		if len(ids) != len(compA) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], &compA[i], nil)
		// 		}
		// 	} else if compA == nil && compB != nil {
		// 		if len(ids) != len(compB) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], nil, &compB[i])
		// 		}
		// 	} else if compA != nil && compB != nil {
		// 		if len(ids) != len(compA) || len(ids) != len(compB) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], &compA[i], &compB[i])
		// 		}
		// 	}
	}

	// Original - doesn't handle optional
	// for _, archId := range v.filter.archIds {
	// 	aSlice, ok := v.storageA.slice[archId]
	// 	if !ok { continue }
	// 	bSlice, ok := v.storageB.slice[archId]
	// 	if !ok { continue }

	// 	lookup, ok := v.world.engine.lookup[archId]
	// 	if !ok { panic("LookupList is missing!") }

	// 	ids := lookup.id
	// 	aComp := aSlice.comp
	// 	bComp := bSlice.comp
	// 	if len(ids) != len(aComp) || len(ids) != len(bComp) {
	// 		panic("ERROR - Bounds don't match")
	// 	}
	// 	for i := range ids {
	// 		if ids[i] == InvalidEntity { continue }
	// 		lambda(ids[i], &aComp[i], &bComp[i])
	// 	}
	// }
}

// Deprecated: This API is a tentative alternative way to map
func (v *View9[A, B, C, D, E, F, G, H, I]) MapSlices(lambda func(id []Id, a []A, b []B, c []C, d []D, e []E, f []F, g []G, h []H, i []I)) {
	v.filter.regenerate(v.world)

	id := make([][]Id, 0)

	sliceListA := make([][]A, 0)
	sliceListB := make([][]B, 0)
	sliceListC := make([][]C, 0)
	sliceListD := make([][]D, 0)
	sliceListE := make([][]E, 0)
	sliceListF := make([][]F, 0)
	sliceListG := make([][]G, 0)
	sliceListH := make([][]H, 0)
	sliceListI := make([][]I, 0)

	for _, archId := range v.filter.archIds {

		sliceA, ok := v.storageA.slice[archId]
		if !ok {
			continue
		}
		sliceB, ok := v.storageB.slice[archId]
		if !ok {
			continue
		}
		sliceC, ok := v.storageC.slice[archId]
		if !ok {
			continue
		}
		sliceD, ok := v.storageD.slice[archId]
		if !ok {
			continue
		}
		sliceE, ok := v.storageE.slice[archId]
		if !ok {
			continue
		}
		sliceF, ok := v.storageF.slice[archId]
		if !ok {
			continue
		}
		sliceG, ok := v.storageG.slice[archId]
		if !ok {
			continue
		}
		sliceH, ok := v.storageH.slice[archId]
		if !ok {
			continue
		}
		sliceI, ok := v.storageI.slice[archId]
		if !ok {
			continue
		}

		lookup, ok := v.world.engine.lookup[archId]
		if !ok {
			panic("LookupList is missing!")
		}

		id = append(id, lookup.id)

		sliceListA = append(sliceListA, sliceA.comp)
		sliceListB = append(sliceListB, sliceB.comp)
		sliceListC = append(sliceListC, sliceC.comp)
		sliceListD = append(sliceListD, sliceD.comp)
		sliceListE = append(sliceListE, sliceE.comp)
		sliceListF = append(sliceListF, sliceF.comp)
		sliceListG = append(sliceListG, sliceG.comp)
		sliceListH = append(sliceListH, sliceH.comp)
		sliceListI = append(sliceListI, sliceI.comp)
	}

	for idx := range id {
		lambda(id[idx],
			sliceListA[idx], sliceListB[idx], sliceListC[idx], sliceListD[idx], sliceListE[idx], sliceListF[idx], sliceListG[idx], sliceListH[idx], sliceListI[idx],
		)
	}
}

// --------------------------------------------------------------------------------
// - View 10
// --------------------------------------------------------------------------------

// Represents a view of data in a specific world. Provides access to the components specified in the generic block
type View10[A, B, C, D, E, F, G, H, I, J any] struct {
	world  *World
	filter filterList

	storageA componentSliceStorage[A]
	storageB componentSliceStorage[B]
	storageC componentSliceStorage[C]
	storageD componentSliceStorage[D]
	storageE componentSliceStorage[E]
	storageF componentSliceStorage[F]
	storageG componentSliceStorage[G]
	storageH componentSliceStorage[H]
	storageI componentSliceStorage[I]
	storageJ componentSliceStorage[J]
}

// Creates a View for the specified world with the specified component filters.
func Query10[A, B, C, D, E, F, G, H, I, J any](world *World, filters ...Filter) *View10[A, B, C, D, E, F, G, H, I, J] {

	storageA := getStorage[A](world.engine)
	storageB := getStorage[B](world.engine)
	storageC := getStorage[C](world.engine)
	storageD := getStorage[D](world.engine)
	storageE := getStorage[E](world.engine)
	storageF := getStorage[F](world.engine)
	storageG := getStorage[G](world.engine)
	storageH := getStorage[H](world.engine)
	storageI := getStorage[I](world.engine)
	storageJ := getStorage[J](world.engine)

	var AA A
	var BB B
	var CC C
	var DD D
	var EE E
	var FF F
	var GG G
	var HH H
	var II I
	var JJ J

	comps := []componentId{

		name(AA),
		name(BB),
		name(CC),
		name(DD),
		name(EE),
		name(FF),
		name(GG),
		name(HH),
		name(II),
		name(JJ),
	}
	filterList := newFilterList(comps, filters...)
	filterList.regenerate(world)

	v := &View10[A, B, C, D, E, F, G, H, I, J]{
		world:  world,
		filter: filterList,

		storageA: storageA,
		storageB: storageB,
		storageC: storageC,
		storageD: storageD,
		storageE: storageE,
		storageF: storageF,
		storageG: storageG,
		storageH: storageH,
		storageI: storageI,
		storageJ: storageJ,
	}
	return v
}

// Reads a pointer to the underlying component at the specified id.
// Read will return even if the specified id doesn't match the filter list
// Read will return the value if it exists, else returns nil.
// If you execute any ecs.Write(...) or ecs.Delete(...) this pointer may become invalid.
func (v *View10[A, B, C, D, E, F, G, H, I, J]) Read(id Id) (*A, *B, *C, *D, *E, *F, *G, *H, *I, *J) {
	if id == InvalidEntity {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil
	}

	archId, ok := v.world.arch[id]
	if !ok {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil
	}
	lookup, ok := v.world.engine.lookup[archId]
	if !ok {
		panic("LookupList is missing!")
	}
	index, ok := lookup.index[id]
	if !ok {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil
	}

	var retA *A
	var retB *B
	var retC *C
	var retD *D
	var retE *E
	var retF *F
	var retG *G
	var retH *H
	var retI *I
	var retJ *J

	sliceA, ok := v.storageA.slice[archId]
	if ok {
		retA = &sliceA.comp[index]
	}
	sliceB, ok := v.storageB.slice[archId]
	if ok {
		retB = &sliceB.comp[index]
	}
	sliceC, ok := v.storageC.slice[archId]
	if ok {
		retC = &sliceC.comp[index]
	}
	sliceD, ok := v.storageD.slice[archId]
	if ok {
		retD = &sliceD.comp[index]
	}
	sliceE, ok := v.storageE.slice[archId]
	if ok {
		retE = &sliceE.comp[index]
	}
	sliceF, ok := v.storageF.slice[archId]
	if ok {
		retF = &sliceF.comp[index]
	}
	sliceG, ok := v.storageG.slice[archId]
	if ok {
		retG = &sliceG.comp[index]
	}
	sliceH, ok := v.storageH.slice[archId]
	if ok {
		retH = &sliceH.comp[index]
	}
	sliceI, ok := v.storageI.slice[archId]
	if ok {
		retI = &sliceI.comp[index]
	}
	sliceJ, ok := v.storageJ.slice[archId]
	if ok {
		retJ = &sliceJ.comp[index]
	}

	return retA, retB, retC, retD, retE, retF, retG, retH, retI, retJ
}

// Maps the lambda function across every entity which matched the specified filters.
func (v *View10[A, B, C, D, E, F, G, H, I, J]) MapId(lambda func(id Id, a *A, b *B, c *C, d *D, e *E, f *F, g *G, h *H, i *I, j *J)) {
	v.filter.regenerate(v.world)

	var sliceA *componentSlice[A]
	var compA []A
	var retA *A

	var sliceB *componentSlice[B]
	var compB []B
	var retB *B

	var sliceC *componentSlice[C]
	var compC []C
	var retC *C

	var sliceD *componentSlice[D]
	var compD []D
	var retD *D

	var sliceE *componentSlice[E]
	var compE []E
	var retE *E

	var sliceF *componentSlice[F]
	var compF []F
	var retF *F

	var sliceG *componentSlice[G]
	var compG []G
	var retG *G

	var sliceH *componentSlice[H]
	var compH []H
	var retH *H

	var sliceI *componentSlice[I]
	var compI []I
	var retI *I

	var sliceJ *componentSlice[J]
	var compJ []J
	var retJ *J

	for _, archId := range v.filter.archIds {

		sliceA, _ = v.storageA.slice[archId]
		sliceB, _ = v.storageB.slice[archId]
		sliceC, _ = v.storageC.slice[archId]
		sliceD, _ = v.storageD.slice[archId]
		sliceE, _ = v.storageE.slice[archId]
		sliceF, _ = v.storageF.slice[archId]
		sliceG, _ = v.storageG.slice[archId]
		sliceH, _ = v.storageH.slice[archId]
		sliceI, _ = v.storageI.slice[archId]
		sliceJ, _ = v.storageJ.slice[archId]

		lookup, ok := v.world.engine.lookup[archId]
		if !ok {
			panic("LookupList is missing!")
		}
		ids := lookup.id

		// TODO - this flattened version causes a mild performance hit. But the other one combinatorially explodes. I also cant get BCE to work with it. See option 2 for higher performance.

		compA = nil
		if sliceA != nil {
			compA = sliceA.comp
		}
		compB = nil
		if sliceB != nil {
			compB = sliceB.comp
		}
		compC = nil
		if sliceC != nil {
			compC = sliceC.comp
		}
		compD = nil
		if sliceD != nil {
			compD = sliceD.comp
		}
		compE = nil
		if sliceE != nil {
			compE = sliceE.comp
		}
		compF = nil
		if sliceF != nil {
			compF = sliceF.comp
		}
		compG = nil
		if sliceG != nil {
			compG = sliceG.comp
		}
		compH = nil
		if sliceH != nil {
			compH = sliceH.comp
		}
		compI = nil
		if sliceI != nil {
			compI = sliceI.comp
		}
		compJ = nil
		if sliceJ != nil {
			compJ = sliceJ.comp
		}

		retA = nil
		retB = nil
		retC = nil
		retD = nil
		retE = nil
		retF = nil
		retG = nil
		retH = nil
		retI = nil
		retJ = nil
		for idx := range ids {
			if ids[idx] == InvalidEntity {
				continue
			} // Skip if its a hole

			if compA != nil {
				retA = &compA[idx]
			}
			if compB != nil {
				retB = &compB[idx]
			}
			if compC != nil {
				retC = &compC[idx]
			}
			if compD != nil {
				retD = &compD[idx]
			}
			if compE != nil {
				retE = &compE[idx]
			}
			if compF != nil {
				retF = &compF[idx]
			}
			if compG != nil {
				retG = &compG[idx]
			}
			if compH != nil {
				retH = &compH[idx]
			}
			if compI != nil {
				retI = &compI[idx]
			}
			if compJ != nil {
				retJ = &compJ[idx]
			}
			lambda(ids[idx], retA, retB, retC, retD, retE, retF, retG, retH, retI, retJ)
		}

		// 	// Option 2 - This is faster but has a combinatorial explosion problem
		// 	if compA == nil && compB == nil {
		// 		return
		// 	} else if compA != nil && compB == nil {
		// 		if len(ids) != len(compA) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], &compA[i], nil)
		// 		}
		// 	} else if compA == nil && compB != nil {
		// 		if len(ids) != len(compB) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], nil, &compB[i])
		// 		}
		// 	} else if compA != nil && compB != nil {
		// 		if len(ids) != len(compA) || len(ids) != len(compB) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], &compA[i], &compB[i])
		// 		}
		// 	}
	}

	// Original - doesn't handle optional
	// for _, archId := range v.filter.archIds {
	// 	aSlice, ok := v.storageA.slice[archId]
	// 	if !ok { continue }
	// 	bSlice, ok := v.storageB.slice[archId]
	// 	if !ok { continue }

	// 	lookup, ok := v.world.engine.lookup[archId]
	// 	if !ok { panic("LookupList is missing!") }

	// 	ids := lookup.id
	// 	aComp := aSlice.comp
	// 	bComp := bSlice.comp
	// 	if len(ids) != len(aComp) || len(ids) != len(bComp) {
	// 		panic("ERROR - Bounds don't match")
	// 	}
	// 	for i := range ids {
	// 		if ids[i] == InvalidEntity { continue }
	// 		lambda(ids[i], &aComp[i], &bComp[i])
	// 	}
	// }
}

// Deprecated: This API is a tentative alternative way to map
func (v *View10[A, B, C, D, E, F, G, H, I, J]) MapSlices(lambda func(id []Id, a []A, b []B, c []C, d []D, e []E, f []F, g []G, h []H, i []I, j []J)) {
	v.filter.regenerate(v.world)

	id := make([][]Id, 0)

	sliceListA := make([][]A, 0)
	sliceListB := make([][]B, 0)
	sliceListC := make([][]C, 0)
	sliceListD := make([][]D, 0)
	sliceListE := make([][]E, 0)
	sliceListF := make([][]F, 0)
	sliceListG := make([][]G, 0)
	sliceListH := make([][]H, 0)
	sliceListI := make([][]I, 0)
	sliceListJ := make([][]J, 0)

	for _, archId := range v.filter.archIds {

		sliceA, ok := v.storageA.slice[archId]
		if !ok {
			continue
		}
		sliceB, ok := v.storageB.slice[archId]
		if !ok {
			continue
		}
		sliceC, ok := v.storageC.slice[archId]
		if !ok {
			continue
		}
		sliceD, ok := v.storageD.slice[archId]
		if !ok {
			continue
		}
		sliceE, ok := v.storageE.slice[archId]
		if !ok {
			continue
		}
		sliceF, ok := v.storageF.slice[archId]
		if !ok {
			continue
		}
		sliceG, ok := v.storageG.slice[archId]
		if !ok {
			continue
		}
		sliceH, ok := v.storageH.slice[archId]
		if !ok {
			continue
		}
		sliceI, ok := v.storageI.slice[archId]
		if !ok {
			continue
		}
		sliceJ, ok := v.storageJ.slice[archId]
		if !ok {
			continue
		}

		lookup, ok := v.world.engine.lookup[archId]
		if !ok {
			panic("LookupList is missing!")
		}

		id = append(id, lookup.id)

		sliceListA = append(sliceListA, sliceA.comp)
		sliceListB = append(sliceListB, sliceB.comp)
		sliceListC = append(sliceListC, sliceC.comp)
		sliceListD = append(sliceListD, sliceD.comp)
		sliceListE = append(sliceListE, sliceE.comp)
		sliceListF = append(sliceListF, sliceF.comp)
		sliceListG = append(sliceListG, sliceG.comp)
		sliceListH = append(sliceListH, sliceH.comp)
		sliceListI = append(sliceListI, sliceI.comp)
		sliceListJ = append(sliceListJ, sliceJ.comp)
	}

	for idx := range id {
		lambda(id[idx],
			sliceListA[idx], sliceListB[idx], sliceListC[idx], sliceListD[idx], sliceListE[idx], sliceListF[idx], sliceListG[idx], sliceListH[idx], sliceListI[idx], sliceListJ[idx],
		)
	}
}

// --------------------------------------------------------------------------------
// - View 11
// --------------------------------------------------------------------------------

// Represents a view of data in a specific world. Provides access to the components specified in the generic block
type View11[A, B, C, D, E, F, G, H, I, J, K any] struct {
	world  *World
	filter filterList

	storageA componentSliceStorage[A]
	storageB componentSliceStorage[B]
	storageC componentSliceStorage[C]
	storageD componentSliceStorage[D]
	storageE componentSliceStorage[E]
	storageF componentSliceStorage[F]
	storageG componentSliceStorage[G]
	storageH componentSliceStorage[H]
	storageI componentSliceStorage[I]
	storageJ componentSliceStorage[J]
	storageK componentSliceStorage[K]
}

// Creates a View for the specified world with the specified component filters.
func Query11[A, B, C, D, E, F, G, H, I, J, K any](world *World, filters ...Filter) *View11[A, B, C, D, E, F, G, H, I, J, K] {

	storageA := getStorage[A](world.engine)
	storageB := getStorage[B](world.engine)
	storageC := getStorage[C](world.engine)
	storageD := getStorage[D](world.engine)
	storageE := getStorage[E](world.engine)
	storageF := getStorage[F](world.engine)
	storageG := getStorage[G](world.engine)
	storageH := getStorage[H](world.engine)
	storageI := getStorage[I](world.engine)
	storageJ := getStorage[J](world.engine)
	storageK := getStorage[K](world.engine)

	var AA A
	var BB B
	var CC C
	var DD D
	var EE E
	var FF F
	var GG G
	var HH H
	var II I
	var JJ J
	var KK K

	comps := []componentId{

		name(AA),
		name(BB),
		name(CC),
		name(DD),
		name(EE),
		name(FF),
		name(GG),
		name(HH),
		name(II),
		name(JJ),
		name(KK),
	}
	filterList := newFilterList(comps, filters...)
	filterList.regenerate(world)

	v := &View11[A, B, C, D, E, F, G, H, I, J, K]{
		world:  world,
		filter: filterList,

		storageA: storageA,
		storageB: storageB,
		storageC: storageC,
		storageD: storageD,
		storageE: storageE,
		storageF: storageF,
		storageG: storageG,
		storageH: storageH,
		storageI: storageI,
		storageJ: storageJ,
		storageK: storageK,
	}
	return v
}

// Reads a pointer to the underlying component at the specified id.
// Read will return even if the specified id doesn't match the filter list
// Read will return the value if it exists, else returns nil.
// If you execute any ecs.Write(...) or ecs.Delete(...) this pointer may become invalid.
func (v *View11[A, B, C, D, E, F, G, H, I, J, K]) Read(id Id) (*A, *B, *C, *D, *E, *F, *G, *H, *I, *J, *K) {
	if id == InvalidEntity {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil
	}

	archId, ok := v.world.arch[id]
	if !ok {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil
	}
	lookup, ok := v.world.engine.lookup[archId]
	if !ok {
		panic("LookupList is missing!")
	}
	index, ok := lookup.index[id]
	if !ok {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil
	}

	var retA *A
	var retB *B
	var retC *C
	var retD *D
	var retE *E
	var retF *F
	var retG *G
	var retH *H
	var retI *I
	var retJ *J
	var retK *K

	sliceA, ok := v.storageA.slice[archId]
	if ok {
		retA = &sliceA.comp[index]
	}
	sliceB, ok := v.storageB.slice[archId]
	if ok {
		retB = &sliceB.comp[index]
	}
	sliceC, ok := v.storageC.slice[archId]
	if ok {
		retC = &sliceC.comp[index]
	}
	sliceD, ok := v.storageD.slice[archId]
	if ok {
		retD = &sliceD.comp[index]
	}
	sliceE, ok := v.storageE.slice[archId]
	if ok {
		retE = &sliceE.comp[index]
	}
	sliceF, ok := v.storageF.slice[archId]
	if ok {
		retF = &sliceF.comp[index]
	}
	sliceG, ok := v.storageG.slice[archId]
	if ok {
		retG = &sliceG.comp[index]
	}
	sliceH, ok := v.storageH.slice[archId]
	if ok {
		retH = &sliceH.comp[index]
	}
	sliceI, ok := v.storageI.slice[archId]
	if ok {
		retI = &sliceI.comp[index]
	}
	sliceJ, ok := v.storageJ.slice[archId]
	if ok {
		retJ = &sliceJ.comp[index]
	}
	sliceK, ok := v.storageK.slice[archId]
	if ok {
		retK = &sliceK.comp[index]
	}

	return retA, retB, retC, retD, retE, retF, retG, retH, retI, retJ, retK
}

// Maps the lambda function across every entity which matched the specified filters.
func (v *View11[A, B, C, D, E, F, G, H, I, J, K]) MapId(lambda func(id Id, a *A, b *B, c *C, d *D, e *E, f *F, g *G, h *H, i *I, j *J, k *K)) {
	v.filter.regenerate(v.world)

	var sliceA *componentSlice[A]
	var compA []A
	var retA *A

	var sliceB *componentSlice[B]
	var compB []B
	var retB *B

	var sliceC *componentSlice[C]
	var compC []C
	var retC *C

	var sliceD *componentSlice[D]
	var compD []D
	var retD *D

	var sliceE *componentSlice[E]
	var compE []E
	var retE *E

	var sliceF *componentSlice[F]
	var compF []F
	var retF *F

	var sliceG *componentSlice[G]
	var compG []G
	var retG *G

	var sliceH *componentSlice[H]
	var compH []H
	var retH *H

	var sliceI *componentSlice[I]
	var compI []I
	var retI *I

	var sliceJ *componentSlice[J]
	var compJ []J
	var retJ *J

	var sliceK *componentSlice[K]
	var compK []K
	var retK *K

	for _, archId := range v.filter.archIds {

		sliceA, _ = v.storageA.slice[archId]
		sliceB, _ = v.storageB.slice[archId]
		sliceC, _ = v.storageC.slice[archId]
		sliceD, _ = v.storageD.slice[archId]
		sliceE, _ = v.storageE.slice[archId]
		sliceF, _ = v.storageF.slice[archId]
		sliceG, _ = v.storageG.slice[archId]
		sliceH, _ = v.storageH.slice[archId]
		sliceI, _ = v.storageI.slice[archId]
		sliceJ, _ = v.storageJ.slice[archId]
		sliceK, _ = v.storageK.slice[archId]

		lookup, ok := v.world.engine.lookup[archId]
		if !ok {
			panic("LookupList is missing!")
		}
		ids := lookup.id

		// TODO - this flattened version causes a mild performance hit. But the other one combinatorially explodes. I also cant get BCE to work with it. See option 2 for higher performance.

		compA = nil
		if sliceA != nil {
			compA = sliceA.comp
		}
		compB = nil
		if sliceB != nil {
			compB = sliceB.comp
		}
		compC = nil
		if sliceC != nil {
			compC = sliceC.comp
		}
		compD = nil
		if sliceD != nil {
			compD = sliceD.comp
		}
		compE = nil
		if sliceE != nil {
			compE = sliceE.comp
		}
		compF = nil
		if sliceF != nil {
			compF = sliceF.comp
		}
		compG = nil
		if sliceG != nil {
			compG = sliceG.comp
		}
		compH = nil
		if sliceH != nil {
			compH = sliceH.comp
		}
		compI = nil
		if sliceI != nil {
			compI = sliceI.comp
		}
		compJ = nil
		if sliceJ != nil {
			compJ = sliceJ.comp
		}
		compK = nil
		if sliceK != nil {
			compK = sliceK.comp
		}

		retA = nil
		retB = nil
		retC = nil
		retD = nil
		retE = nil
		retF = nil
		retG = nil
		retH = nil
		retI = nil
		retJ = nil
		retK = nil
		for idx := range ids {
			if ids[idx] == InvalidEntity {
				continue
			} // Skip if its a hole

			if compA != nil {
				retA = &compA[idx]
			}
			if compB != nil {
				retB = &compB[idx]
			}
			if compC != nil {
				retC = &compC[idx]
			}
			if compD != nil {
				retD = &compD[idx]
			}
			if compE != nil {
				retE = &compE[idx]
			}
			if compF != nil {
				retF = &compF[idx]
			}
			if compG != nil {
				retG = &compG[idx]
			}
			if compH != nil {
				retH = &compH[idx]
			}
			if compI != nil {
				retI = &compI[idx]
			}
			if compJ != nil {
				retJ = &compJ[idx]
			}
			if compK != nil {
				retK = &compK[idx]
			}
			lambda(ids[idx], retA, retB, retC, retD, retE, retF, retG, retH, retI, retJ, retK)
		}

		// 	// Option 2 - This is faster but has a combinatorial explosion problem
		// 	if compA == nil && compB == nil {
		// 		return
		// 	} else if compA != nil && compB == nil {
		// 		if len(ids) != len(compA) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], &compA[i], nil)
		// 		}
		// 	} else if compA == nil && compB != nil {
		// 		if len(ids) != len(compB) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], nil, &compB[i])
		// 		}
		// 	} else if compA != nil && compB != nil {
		// 		if len(ids) != len(compA) || len(ids) != len(compB) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], &compA[i], &compB[i])
		// 		}
		// 	}
	}

	// Original - doesn't handle optional
	// for _, archId := range v.filter.archIds {
	// 	aSlice, ok := v.storageA.slice[archId]
	// 	if !ok { continue }
	// 	bSlice, ok := v.storageB.slice[archId]
	// 	if !ok { continue }

	// 	lookup, ok := v.world.engine.lookup[archId]
	// 	if !ok { panic("LookupList is missing!") }

	// 	ids := lookup.id
	// 	aComp := aSlice.comp
	// 	bComp := bSlice.comp
	// 	if len(ids) != len(aComp) || len(ids) != len(bComp) {
	// 		panic("ERROR - Bounds don't match")
	// 	}
	// 	for i := range ids {
	// 		if ids[i] == InvalidEntity { continue }
	// 		lambda(ids[i], &aComp[i], &bComp[i])
	// 	}
	// }
}

// Deprecated: This API is a tentative alternative way to map
func (v *View11[A, B, C, D, E, F, G, H, I, J, K]) MapSlices(lambda func(id []Id, a []A, b []B, c []C, d []D, e []E, f []F, g []G, h []H, i []I, j []J, k []K)) {
	v.filter.regenerate(v.world)

	id := make([][]Id, 0)

	sliceListA := make([][]A, 0)
	sliceListB := make([][]B, 0)
	sliceListC := make([][]C, 0)
	sliceListD := make([][]D, 0)
	sliceListE := make([][]E, 0)
	sliceListF := make([][]F, 0)
	sliceListG := make([][]G, 0)
	sliceListH := make([][]H, 0)
	sliceListI := make([][]I, 0)
	sliceListJ := make([][]J, 0)
	sliceListK := make([][]K, 0)

	for _, archId := range v.filter.archIds {

		sliceA, ok := v.storageA.slice[archId]
		if !ok {
			continue
		}
		sliceB, ok := v.storageB.slice[archId]
		if !ok {
			continue
		}
		sliceC, ok := v.storageC.slice[archId]
		if !ok {
			continue
		}
		sliceD, ok := v.storageD.slice[archId]
		if !ok {
			continue
		}
		sliceE, ok := v.storageE.slice[archId]
		if !ok {
			continue
		}
		sliceF, ok := v.storageF.slice[archId]
		if !ok {
			continue
		}
		sliceG, ok := v.storageG.slice[archId]
		if !ok {
			continue
		}
		sliceH, ok := v.storageH.slice[archId]
		if !ok {
			continue
		}
		sliceI, ok := v.storageI.slice[archId]
		if !ok {
			continue
		}
		sliceJ, ok := v.storageJ.slice[archId]
		if !ok {
			continue
		}
		sliceK, ok := v.storageK.slice[archId]
		if !ok {
			continue
		}

		lookup, ok := v.world.engine.lookup[archId]
		if !ok {
			panic("LookupList is missing!")
		}

		id = append(id, lookup.id)

		sliceListA = append(sliceListA, sliceA.comp)
		sliceListB = append(sliceListB, sliceB.comp)
		sliceListC = append(sliceListC, sliceC.comp)
		sliceListD = append(sliceListD, sliceD.comp)
		sliceListE = append(sliceListE, sliceE.comp)
		sliceListF = append(sliceListF, sliceF.comp)
		sliceListG = append(sliceListG, sliceG.comp)
		sliceListH = append(sliceListH, sliceH.comp)
		sliceListI = append(sliceListI, sliceI.comp)
		sliceListJ = append(sliceListJ, sliceJ.comp)
		sliceListK = append(sliceListK, sliceK.comp)
	}

	for idx := range id {
		lambda(id[idx],
			sliceListA[idx], sliceListB[idx], sliceListC[idx], sliceListD[idx], sliceListE[idx], sliceListF[idx], sliceListG[idx], sliceListH[idx], sliceListI[idx], sliceListJ[idx], sliceListK[idx],
		)
	}
}

// --------------------------------------------------------------------------------
// - View 12
// --------------------------------------------------------------------------------

// Represents a view of data in a specific world. Provides access to the components specified in the generic block
type View12[A, B, C, D, E, F, G, H, I, J, K, L any] struct {
	world  *World
	filter filterList

	storageA componentSliceStorage[A]
	storageB componentSliceStorage[B]
	storageC componentSliceStorage[C]
	storageD componentSliceStorage[D]
	storageE componentSliceStorage[E]
	storageF componentSliceStorage[F]
	storageG componentSliceStorage[G]
	storageH componentSliceStorage[H]
	storageI componentSliceStorage[I]
	storageJ componentSliceStorage[J]
	storageK componentSliceStorage[K]
	storageL componentSliceStorage[L]
}

// Creates a View for the specified world with the specified component filters.
func Query12[A, B, C, D, E, F, G, H, I, J, K, L any](world *World, filters ...Filter) *View12[A, B, C, D, E, F, G, H, I, J, K, L] {

	storageA := getStorage[A](world.engine)
	storageB := getStorage[B](world.engine)
	storageC := getStorage[C](world.engine)
	storageD := getStorage[D](world.engine)
	storageE := getStorage[E](world.engine)
	storageF := getStorage[F](world.engine)
	storageG := getStorage[G](world.engine)
	storageH := getStorage[H](world.engine)
	storageI := getStorage[I](world.engine)
	storageJ := getStorage[J](world.engine)
	storageK := getStorage[K](world.engine)
	storageL := getStorage[L](world.engine)

	var AA A
	var BB B
	var CC C
	var DD D
	var EE E
	var FF F
	var GG G
	var HH H
	var II I
	var JJ J
	var KK K
	var LL L

	comps := []componentId{

		name(AA),
		name(BB),
		name(CC),
		name(DD),
		name(EE),
		name(FF),
		name(GG),
		name(HH),
		name(II),
		name(JJ),
		name(KK),
		name(LL),
	}
	filterList := newFilterList(comps, filters...)
	filterList.regenerate(world)

	v := &View12[A, B, C, D, E, F, G, H, I, J, K, L]{
		world:  world,
		filter: filterList,

		storageA: storageA,
		storageB: storageB,
		storageC: storageC,
		storageD: storageD,
		storageE: storageE,
		storageF: storageF,
		storageG: storageG,
		storageH: storageH,
		storageI: storageI,
		storageJ: storageJ,
		storageK: storageK,
		storageL: storageL,
	}
	return v
}

// Reads a pointer to the underlying component at the specified id.
// Read will return even if the specified id doesn't match the filter list
// Read will return the value if it exists, else returns nil.
// If you execute any ecs.Write(...) or ecs.Delete(...) this pointer may become invalid.
func (v *View12[A, B, C, D, E, F, G, H, I, J, K, L]) Read(id Id) (*A, *B, *C, *D, *E, *F, *G, *H, *I, *J, *K, *L) {
	if id == InvalidEntity {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil
	}

	archId, ok := v.world.arch[id]
	if !ok {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil
	}
	lookup, ok := v.world.engine.lookup[archId]
	if !ok {
		panic("LookupList is missing!")
	}
	index, ok := lookup.index[id]
	if !ok {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil
	}

	var retA *A
	var retB *B
	var retC *C
	var retD *D
	var retE *E
	var retF *F
	var retG *G
	var retH *H
	var retI *I
	var retJ *J
	var retK *K
	var retL *L

	sliceA, ok := v.storageA.slice[archId]
	if ok {
		retA = &sliceA.comp[index]
	}
	sliceB, ok := v.storageB.slice[archId]
	if ok {
		retB = &sliceB.comp[index]
	}
	sliceC, ok := v.storageC.slice[archId]
	if ok {
		retC = &sliceC.comp[index]
	}
	sliceD, ok := v.storageD.slice[archId]
	if ok {
		retD = &sliceD.comp[index]
	}
	sliceE, ok := v.storageE.slice[archId]
	if ok {
		retE = &sliceE.comp[index]
	}
	sliceF, ok := v.storageF.slice[archId]
	if ok {
		retF = &sliceF.comp[index]
	}
	sliceG, ok := v.storageG.slice[archId]
	if ok {
		retG = &sliceG.comp[index]
	}
	sliceH, ok := v.storageH.slice[archId]
	if ok {
		retH = &sliceH.comp[index]
	}
	sliceI, ok := v.storageI.slice[archId]
	if ok {
		retI = &sliceI.comp[index]
	}
	sliceJ, ok := v.storageJ.slice[archId]
	if ok {
		retJ = &sliceJ.comp[index]
	}
	sliceK, ok := v.storageK.slice[archId]
	if ok {
		retK = &sliceK.comp[index]
	}
	sliceL, ok := v.storageL.slice[archId]
	if ok {
		retL = &sliceL.comp[index]
	}

	return retA, retB, retC, retD, retE, retF, retG, retH, retI, retJ, retK, retL
}

// Maps the lambda function across every entity which matched the specified filters.
func (v *View12[A, B, C, D, E, F, G, H, I, J, K, L]) MapId(lambda func(id Id, a *A, b *B, c *C, d *D, e *E, f *F, g *G, h *H, i *I, j *J, k *K, l *L)) {
	v.filter.regenerate(v.world)

	var sliceA *componentSlice[A]
	var compA []A
	var retA *A

	var sliceB *componentSlice[B]
	var compB []B
	var retB *B

	var sliceC *componentSlice[C]
	var compC []C
	var retC *C

	var sliceD *componentSlice[D]
	var compD []D
	var retD *D

	var sliceE *componentSlice[E]
	var compE []E
	var retE *E

	var sliceF *componentSlice[F]
	var compF []F
	var retF *F

	var sliceG *componentSlice[G]
	var compG []G
	var retG *G

	var sliceH *componentSlice[H]
	var compH []H
	var retH *H

	var sliceI *componentSlice[I]
	var compI []I
	var retI *I

	var sliceJ *componentSlice[J]
	var compJ []J
	var retJ *J

	var sliceK *componentSlice[K]
	var compK []K
	var retK *K

	var sliceL *componentSlice[L]
	var compL []L
	var retL *L

	for _, archId := range v.filter.archIds {

		sliceA, _ = v.storageA.slice[archId]
		sliceB, _ = v.storageB.slice[archId]
		sliceC, _ = v.storageC.slice[archId]
		sliceD, _ = v.storageD.slice[archId]
		sliceE, _ = v.storageE.slice[archId]
		sliceF, _ = v.storageF.slice[archId]
		sliceG, _ = v.storageG.slice[archId]
		sliceH, _ = v.storageH.slice[archId]
		sliceI, _ = v.storageI.slice[archId]
		sliceJ, _ = v.storageJ.slice[archId]
		sliceK, _ = v.storageK.slice[archId]
		sliceL, _ = v.storageL.slice[archId]

		lookup, ok := v.world.engine.lookup[archId]
		if !ok {
			panic("LookupList is missing!")
		}
		ids := lookup.id

		// TODO - this flattened version causes a mild performance hit. But the other one combinatorially explodes. I also cant get BCE to work with it. See option 2 for higher performance.

		compA = nil
		if sliceA != nil {
			compA = sliceA.comp
		}
		compB = nil
		if sliceB != nil {
			compB = sliceB.comp
		}
		compC = nil
		if sliceC != nil {
			compC = sliceC.comp
		}
		compD = nil
		if sliceD != nil {
			compD = sliceD.comp
		}
		compE = nil
		if sliceE != nil {
			compE = sliceE.comp
		}
		compF = nil
		if sliceF != nil {
			compF = sliceF.comp
		}
		compG = nil
		if sliceG != nil {
			compG = sliceG.comp
		}
		compH = nil
		if sliceH != nil {
			compH = sliceH.comp
		}
		compI = nil
		if sliceI != nil {
			compI = sliceI.comp
		}
		compJ = nil
		if sliceJ != nil {
			compJ = sliceJ.comp
		}
		compK = nil
		if sliceK != nil {
			compK = sliceK.comp
		}
		compL = nil
		if sliceL != nil {
			compL = sliceL.comp
		}

		retA = nil
		retB = nil
		retC = nil
		retD = nil
		retE = nil
		retF = nil
		retG = nil
		retH = nil
		retI = nil
		retJ = nil
		retK = nil
		retL = nil
		for idx := range ids {
			if ids[idx] == InvalidEntity {
				continue
			} // Skip if its a hole

			if compA != nil {
				retA = &compA[idx]
			}
			if compB != nil {
				retB = &compB[idx]
			}
			if compC != nil {
				retC = &compC[idx]
			}
			if compD != nil {
				retD = &compD[idx]
			}
			if compE != nil {
				retE = &compE[idx]
			}
			if compF != nil {
				retF = &compF[idx]
			}
			if compG != nil {
				retG = &compG[idx]
			}
			if compH != nil {
				retH = &compH[idx]
			}
			if compI != nil {
				retI = &compI[idx]
			}
			if compJ != nil {
				retJ = &compJ[idx]
			}
			if compK != nil {
				retK = &compK[idx]
			}
			if compL != nil {
				retL = &compL[idx]
			}
			lambda(ids[idx], retA, retB, retC, retD, retE, retF, retG, retH, retI, retJ, retK, retL)
		}

		// 	// Option 2 - This is faster but has a combinatorial explosion problem
		// 	if compA == nil && compB == nil {
		// 		return
		// 	} else if compA != nil && compB == nil {
		// 		if len(ids) != len(compA) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], &compA[i], nil)
		// 		}
		// 	} else if compA == nil && compB != nil {
		// 		if len(ids) != len(compB) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], nil, &compB[i])
		// 		}
		// 	} else if compA != nil && compB != nil {
		// 		if len(ids) != len(compA) || len(ids) != len(compB) {
		// 			panic("ERROR - Bounds don't match")
		// 		}
		// 		for i := range ids {
		// 			if ids[i] == InvalidEntity { continue }
		// 			lambda(ids[i], &compA[i], &compB[i])
		// 		}
		// 	}
	}

	// Original - doesn't handle optional
	// for _, archId := range v.filter.archIds {
	// 	aSlice, ok := v.storageA.slice[archId]
	// 	if !ok { continue }
	// 	bSlice, ok := v.storageB.slice[archId]
	// 	if !ok { continue }

	// 	lookup, ok := v.world.engine.lookup[archId]
	// 	if !ok { panic("LookupList is missing!") }

	// 	ids := lookup.id
	// 	aComp := aSlice.comp
	// 	bComp := bSlice.comp
	// 	if len(ids) != len(aComp) || len(ids) != len(bComp) {
	// 		panic("ERROR - Bounds don't match")
	// 	}
	// 	for i := range ids {
	// 		if ids[i] == InvalidEntity { continue }
	// 		lambda(ids[i], &aComp[i], &bComp[i])
	// 	}
	// }
}

// Deprecated: This API is a tentative alternative way to map
func (v *View12[A, B, C, D, E, F, G, H, I, J, K, L]) MapSlices(lambda func(id []Id, a []A, b []B, c []C, d []D, e []E, f []F, g []G, h []H, i []I, j []J, k []K, l []L)) {
	v.filter.regenerate(v.world)

	id := make([][]Id, 0)

	sliceListA := make([][]A, 0)
	sliceListB := make([][]B, 0)
	sliceListC := make([][]C, 0)
	sliceListD := make([][]D, 0)
	sliceListE := make([][]E, 0)
	sliceListF := make([][]F, 0)
	sliceListG := make([][]G, 0)
	sliceListH := make([][]H, 0)
	sliceListI := make([][]I, 0)
	sliceListJ := make([][]J, 0)
	sliceListK := make([][]K, 0)
	sliceListL := make([][]L, 0)

	for _, archId := range v.filter.archIds {

		sliceA, ok := v.storageA.slice[archId]
		if !ok {
			continue
		}
		sliceB, ok := v.storageB.slice[archId]
		if !ok {
			continue
		}
		sliceC, ok := v.storageC.slice[archId]
		if !ok {
			continue
		}
		sliceD, ok := v.storageD.slice[archId]
		if !ok {
			continue
		}
		sliceE, ok := v.storageE.slice[archId]
		if !ok {
			continue
		}
		sliceF, ok := v.storageF.slice[archId]
		if !ok {
			continue
		}
		sliceG, ok := v.storageG.slice[archId]
		if !ok {
			continue
		}
		sliceH, ok := v.storageH.slice[archId]
		if !ok {
			continue
		}
		sliceI, ok := v.storageI.slice[archId]
		if !ok {
			continue
		}
		sliceJ, ok := v.storageJ.slice[archId]
		if !ok {
			continue
		}
		sliceK, ok := v.storageK.slice[archId]
		if !ok {
			continue
		}
		sliceL, ok := v.storageL.slice[archId]
		if !ok {
			continue
		}

		lookup, ok := v.world.engine.lookup[archId]
		if !ok {
			panic("LookupList is missing!")
		}

		id = append(id, lookup.id)

		sliceListA = append(sliceListA, sliceA.comp)
		sliceListB = append(sliceListB, sliceB.comp)
		sliceListC = append(sliceListC, sliceC.comp)
		sliceListD = append(sliceListD, sliceD.comp)
		sliceListE = append(sliceListE, sliceE.comp)
		sliceListF = append(sliceListF, sliceF.comp)
		sliceListG = append(sliceListG, sliceG.comp)
		sliceListH = append(sliceListH, sliceH.comp)
		sliceListI = append(sliceListI, sliceI.comp)
		sliceListJ = append(sliceListJ, sliceJ.comp)
		sliceListK = append(sliceListK, sliceK.comp)
		sliceListL = append(sliceListL, sliceL.comp)
	}

	for idx := range id {
		lambda(id[idx],
			sliceListA[idx], sliceListB[idx], sliceListC[idx], sliceListD[idx], sliceListE[idx], sliceListF[idx], sliceListG[idx], sliceListH[idx], sliceListI[idx], sliceListJ[idx], sliceListK[idx], sliceListL[idx],
		)
	}
}
