package multiple

import (
	"errors"
	"math"
	"prerequisite/repository"

	"github.com/jinzhu/copier"
)

// FIXME panics => print result before panicing
type all struct {
	edgeList map[string]repository.Edges
	lookUp   map[string]string
	inDegree map[string]weightedInDegree
	result   []repository.Edges
}

type weightedInDegree struct {
	count         int
	currentWeight float64
}

func NewAllGenerate(edgeList map[string]repository.Edges, lookUp map[string]string, simpleInDegree map[string]int) *all {
	inDegree := map[string]weightedInDegree{}

	for id, count := range simpleInDegree {
		inDegree[id] = weightedInDegree{
			count: count,
		}
	}

	return &all{
		edgeList: edgeList,
		lookUp:   lookUp,
		inDegree: inDegree,
		result:   make([]repository.Edges, 0),
	}
}

func (r *all) GetPath() ([]repository.Edges, error) {

	if err := r.getAll(r.edgeList, r.lookUp, r.inDegree, make(repository.Edges, 0)); err != nil {
		return nil, err
	}

	return r.result, nil
}

func (r *all) getAll(edgeList map[string]repository.Edges, lookUp map[string]string,
	inDegree map[string]weightedInDegree, result repository.Edges) error {

	var detected bool

	if len(inDegree) == 0 {
		r.result = append(r.result, result)

		return nil
	}

	for id, degree := range inDegree {
		if degree.count == 0 {
			newResult := make(repository.Edges, 0)
			copier.Copy(&newResult, &result)
			newInDegree := make(map[string]weightedInDegree)
			copier.Copy(&newInDegree, &inDegree)
			detected = true
			newResult = append(newResult, repository.Edge{
				Id:     id,
				Weight: math.Round(degree.currentWeight*100) / 100,
			})
			delete(newInDegree, id)
			for _, edge := range edgeList[id] {
				tmp := newInDegree[edge.Id]
				tmp.count -= 1
				tmp.currentWeight = edge.Weight
				newInDegree[edge.Id] = tmp
			}

			r.getAll(edgeList, lookUp, newInDegree, newResult)
		}
	}

	if !detected {
		return errors.New("Cycle detected")
	}

	return nil
}
