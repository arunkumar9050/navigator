// package hook is used to ensure execution of a set of pre-start, post-start,
// pre-stop and post-stop hooks. The transitioning between states is gated by
// the consumer of the package, through use of the Hooks.Transition function.
package hook

import (
	"fmt"
	"sync"

	"github.com/golang/glog"

	"github.com/jetstack/navigator/pkg/apis/navigator/v1alpha1"
)

type Interface interface {
	// The name of the Hook
	Name() string
	// Execute the action with the given Pilot resource
	Execute(pilot *v1alpha1.Pilot) error
}

func New(name string, fn func(*v1alpha1.Pilot) error) Interface {
	return &hookAdapter{
		name: name,
		fn:   fn,
	}
}

type hookAdapter struct {
	name string
	fn   func(pilot *v1alpha1.Pilot) error
}

var _ Interface = &hookAdapter{}

func (h *hookAdapter) Name() string {
	return h.name
}

func (h *hookAdapter) Execute(pilot *v1alpha1.Pilot) error {
	return h.fn(pilot)
}

type Hooks struct {
	// PreStart are hooks to be run before the application starts
	PreStart         []Interface
	executedPreStart map[string]struct{}
	// PostStart are hooks to be run after the application starts
	PostStart         []Interface
	executedPostStart map[string]struct{}
	// PreStop are hooks to be run before the application stops
	PreStop         []Interface
	executedPreStop map[string]struct{}
	// PostStop are hooks to be run after the application stops
	PostStop         []Interface
	executedPostStop map[string]struct{}

	lock sync.Mutex
}

// CurrentPhase returns the current phase
func (h *Hooks) CurrentPhase() v1alpha1.PilotPhase {
	switch {
	case len(h.executedPostStop) == len(h.PostStop):
		return v1alpha1.PilotPhasePostStop
	case len(h.executedPreStop) == len(h.PreStop):
		return v1alpha1.PilotPhasePreStop
	case len(h.executedPostStart) == len(h.PostStart):
		return v1alpha1.PilotPhasePostStart
	case len(h.executedPreStart) == len(h.PreStart):
		return v1alpha1.PilotPhasePreStart
	default:
		return ""
	}
}

func (h *Hooks) Transition(p v1alpha1.PilotPhase, pilot *v1alpha1.Pilot) error {
	h.lock.Lock()
	defer h.lock.Unlock()
	if pilot == nil {
		return fmt.Errorf("pilot resource must not be nil")
	}
	var hooks []Interface
	executed := map[string]struct{}{}
	switch p {
	case v1alpha1.PilotPhasePreStart:
		hooks = h.PreStart
		if h.executedPreStart == nil {
			h.executedPreStart = executed
		}
		executed = h.executedPreStart
	case v1alpha1.PilotPhasePostStart:
		hooks = h.PostStart
		if h.executedPostStart == nil {
			h.executedPostStart = executed
		}
		executed = h.executedPostStart
	case v1alpha1.PilotPhasePreStop:
		hooks = h.PreStop
		if h.executedPreStop == nil {
			h.executedPreStop = executed
		}
		executed = h.executedPreStop
	case v1alpha1.PilotPhasePostStop:
		hooks = h.PostStop
		if h.executedPostStop == nil {
			h.executedPostStop = executed
		}
		executed = h.executedPostStop
	default:
		return fmt.Errorf("invalid phase: %q", p)
	}
	for _, hook := range hooks {
		if _, ok := executed[hook.Name()]; ok {
			glog.V(4).Infof("Skipping already executed hook for %q in phase %q", hook.Name(), p)
			continue
		}
		glog.V(4).Infof("Executing %s hook '%s'", p, hook.Name())
		if err := hook.Execute(pilot); err != nil {
			glog.V(4).Infof("Error executing %s hook %q: %s", p, hook.Name(), err.Error())
			return fmt.Errorf("error executing %s hook %q: %s", p, hook.Name(), err.Error())
		}
		glog.V(4).Infof("Executed %s hook %q", p, hook.Name())
		executed[hook.Name()] = struct{}{}
	}
	return nil
}
