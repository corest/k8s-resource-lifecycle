package metaresource

import (
	"bufio"
	"encoding/json"
	"os"
	"sort"
	"sync"

	"github.com/giantswarm/microerror"
	"k8s.io/apiserver/pkg/apis/audit"
)

type MetaResource struct {
	Kind      string
	Name      string
	Namespace string
	APIGroup  string
	Events    []audit.Event
}

func (m *MetaResource) FindEvents(f string, storeCh chan<- audit.Event, errCh chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(f)
	if err != nil {
		errCh <- microerror.Mask(err)
	}
	defer file.Close()

	var event audit.Event

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		err := json.Unmarshal(scanner.Bytes(), &event)
		if err != nil {
			errCh <- microerror.Mask(err)
		}

		if isResponseComplete(event.Stage) && isMutatingOperation(event.Verb) && m.isTargetResource(event) {
			storeCh <- *event.DeepCopy()
		}
	}
}

func (m *MetaResource) StoreEvent(event audit.Event) {
	m.Events = append(m.Events, event)
}

func (m *MetaResource) SortEvents() {
	sort.Slice(m.Events, func(i, j int) bool {
		return m.Events[i].StageTimestamp.Time.Before(m.Events[j].StageTimestamp.Time)
	})
}

func (m *MetaResource) isTargetResource(event audit.Event) bool {
	return event.ObjectRef.Resource == m.Kind &&
		event.ObjectRef.Name == m.Name &&
		event.ObjectRef.Namespace == m.Namespace &&
		event.ObjectRef.APIGroup == m.APIGroup
}

func isMutatingOperation(op string) bool {
	mutatingOperations := []string{"update", "create", "delete", "patch"}

	for _, s := range mutatingOperations {
		if s == op {
			return true
		}
	}

	return false
}

func isResponseComplete(stage audit.Stage) bool {
	return stage == audit.StageResponseComplete
}
