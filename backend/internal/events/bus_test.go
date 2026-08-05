package events

import "testing"

func TestBusScopesAndSequencesEventsByRun(t *testing.T) {
	bus := NewBus(2)
	first, unsubscribe := bus.Subscribe("run-1")
	defer unsubscribe()
	other, unsubscribeOther := bus.Subscribe("run-2")
	defer unsubscribeOther()

	bus.Publish(Event{RunID: "run-1", Type: RunStarted})
	bus.Publish(Event{RunID: "run-1", Type: RunCompleted})
	if event := <-first; event.ID != "1" || event.Type != RunStarted {
		t.Errorf("first event = %#v", event)
	}
	if event := <-first; event.ID != "2" || event.Type != RunCompleted {
		t.Errorf("second event = %#v", event)
	}
	select {
	case event := <-other:
		t.Errorf("other run received %#v", event)
	default:
	}
}

func TestBusDisconnectsSlowSubscriber(t *testing.T) {
	bus := NewBus(1)
	channel, _ := bus.Subscribe("run")
	bus.Publish(Event{RunID: "run", Type: RunStarted})
	bus.Publish(Event{RunID: "run", Type: RunCompleted})

	<-channel
	if _, open := <-channel; open {
		t.Error("slow subscriber channel remains open")
	}
}

func TestBusCloseClosesSubscriptions(t *testing.T) {
	bus := NewBus(1)
	channel, unsubscribe := bus.Subscribe("run")
	bus.Close()
	unsubscribe()
	bus.Close()

	if _, open := <-channel; open {
		t.Error("subscriber remains open after Close")
	}
}
