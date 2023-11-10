package single

import (
	"errors"
	"math"
	"prerequisite/repository"

	pq "github.com/emirpasic/gods/queues/priorityqueue"
	"github.com/emirpasic/gods/utils"
)

// FIXME panics => print result before panicing
type MaximalPath struct {
	edgeList map[string]repository.Edges
	lookUp   map[string]string
	inDegree map[string]int
	result   repository.Edges
	//
	idWeight map[string]float64
	_pq      *pq.Queue
}

// Element is an entry in the priority queue
type vertex struct {
	id     string
	weight float64
}

// Comparator function (sort by vertex's priority value in descending order)
func byPriority(a, b interface{}) int {
	priorityA := a.(vertex).weight
	priorityB := b.(vertex).weight
	return utils.Float64Comparator(priorityB, priorityA) // "-" descending order
}

func NewMaximalPathGenerate(edgeList map[string]repository.Edges, lookUp map[string]string, inDegree map[string]int) *MaximalPath {
	return &MaximalPath{
		edgeList: edgeList,
		lookUp:   lookUp,
		inDegree: inDegree,
		result:   make(repository.Edges, 0),
		_pq:      pq.NewWith(byPriority),
		idWeight: map[string]float64{},
	}
}

func (r *MaximalPath) GetPath() (repository.Edges, error) {
	if err := r.getMaximalPath(); err != nil {
		return nil, err
	}

	return r.result, nil
}

func (r *MaximalPath) getMaximalPath() error {
	if len(r.inDegree) == 0 && r._pq.Empty() {
		return nil
	}

	for id, degree := range r.inDegree {
		if degree == 0 {
			wt, _ := r.idWeight[id]
			r._pq.Enqueue(vertex{
				id:     id,
				weight: wt,
			})
			delete(r.inDegree, id)
		}
	}

	// TODO see and change random
	if r._pq.Empty() && !r._pq.Empty() {
		return errors.New("Cycle detected")
	}

	chosenAny, _ := r._pq.Dequeue()
	chosen, _ := chosenAny.(vertex)

	r.result = append(r.result, repository.Edge{
		Id:     chosen.id,
		Weight: math.Round(chosen.weight*100) / 100,
	})

	for _, edge := range r.edgeList[chosen.id] {
		r.inDegree[edge.Id] -= 1

		if r.inDegree[edge.Id] == 0 {
			existingWt, ok := r.idWeight[edge.Id]
			if !ok || existingWt < edge.Weight {
				r.idWeight[edge.Id] = edge.Weight
			}
		}
	}

	return r.getMaximalPath()
}
