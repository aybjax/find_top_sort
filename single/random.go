package single

import (
	"math"
	"math/rand"
	"prerequisite/repository"

	dll "github.com/emirpasic/gods/lists/doublylinkedlist"
)

type randomPath struct {
	edgeList map[string]repository.Edges
	lookUp   map[string]string
	inDegree map[string]int
	dllist   *dll.List
	weight   float64
	result   repository.Edges
	//
	idWeight map[string]float64
}

func NewRandomPathGenerate(edgeList map[string]repository.Edges, lookUp map[string]string, inDegree map[string]int) *randomPath {
	return &randomPath{
		edgeList: edgeList,
		lookUp:   lookUp,
		inDegree: inDegree,
		dllist:   dll.New(),
		weight:   0,
		result:   make(repository.Edges, 0),
		idWeight: map[string]float64{},
	}
}

func (r *randomPath) GetPath() (repository.Edges, error) {
	if err := r.getRandomPath(); err != nil {
		return nil, err
	}

	return r.result, nil
}

func (r *randomPath) getRandomPath() error {
	if len(r.inDegree) == 0 && r.dllist.Empty() {
		return nil
	}

	for id, degree := range r.inDegree {
		if degree == 0 {
			r.dllist.Add(id)
			wt, _ := r.idWeight[id]
			r.weight += wt
			delete(r.inDegree, id)
		}
	}

	index := int(rand.Float64() * float64(r.weight))
	var chosen string

	for i, _ := range r.dllist.Values() {
		if i > index {
			chosenAny, _ := r.dllist.Get(i - 1)
			chosen, _ = chosenAny.(string)
			r.dllist.Remove(i - 1)
			wt, _ := r.idWeight[chosen]
			r.weight -= wt
			break
		}
	}

	var wt float64
	if chosen == "" {
		chosenAny, _ := r.dllist.Get(r.dllist.Size() - 1)
		chosen, _ = chosenAny.(string)
		r.dllist.Remove(r.dllist.Size() - 1)
		wt, _ = r.idWeight[chosen]
		r.weight -= wt
	}

	r.result = append(r.result, repository.Edge{
		Id:     chosen,
		Weight: math.Round(wt*100) / 100,
	})

	for _, edge := range r.edgeList[chosen] {
		r.inDegree[edge.Id] -= 1

		if r.inDegree[edge.Id] == 0 {
			existingWt, ok := r.idWeight[edge.Id]
			if !ok || existingWt < edge.Weight {
				r.idWeight[edge.Id] = edge.Weight
			}
		}
	}

	return r.getRandomPath()
}
