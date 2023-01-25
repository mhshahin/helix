package kubernetes

import (
	"log"

	"github.com/mhshahin/helix/pkg/handler"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type EventWatcher struct {
	informer cache.SharedInformer
	stopper  chan struct{}
	fn       handler.EventHandler
}

// NewEventWatcher ...
func NewEventWatcher(config *rest.Config, namespace string, fn handler.EventHandler) *EventWatcher {
	clientset := kubernetes.NewForConfigOrDie(config)
	factory := informers.NewSharedInformerFactoryWithOptions(clientset, 0, informers.WithNamespace(namespace))
	informer := factory.Core().V1().Events().Informer()

	watcher := &EventWatcher{
		informer: informer,
		stopper:  make(chan struct{}),
		fn:       fn,
	}

	informer.AddEventHandler(watcher)
	informer.SetWatchErrorHandler(func(r *cache.Reflector, err error) {
		log.Println(err)
	})

	return watcher
}

func (ew *EventWatcher) OnAdd(obj interface{}) {
	event := obj.(*corev1.Event)
	ew.fn(event)
}
func (ew *EventWatcher) OnUpdate(oldObj, newObj interface{}) {}
func (ew *EventWatcher) OnDelete(obj interface{})            {}

func (ew *EventWatcher) Start() {
	go ew.informer.Run(ew.stopper)
}

func (ew *EventWatcher) Stop() {
	ew.stopper <- struct{}{}
	close(ew.stopper)
}
