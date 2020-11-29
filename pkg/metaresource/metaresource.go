package metaresource

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"

	"github.com/giantswarm/microerror"
	"k8s.io/apiserver/pkg/apis/audit"
)

type MetaResource struct {
	Kind       string
	Name       string
	Namespace  string
	APIGroup   string
	APIVersion string
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

func (m *MetaResource) isTargetResource(event audit.Event) bool {
	return event.ObjectRef.Resource == m.Kind &&
		event.ObjectRef.Name == m.Name &&
		event.ObjectRef.Namespace == m.Namespace &&
		event.ObjectRef.APIGroup == m.APIGroup &&
		event.ObjectRef.APIVersion == m.APIVersion
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
