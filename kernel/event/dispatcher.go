package event

import "sync"

// Data holds dispatched event data
type Data struct {
	Name  string
	Value interface{}
}

type dispatchData struct {
	event     *Data
	waitGroup *sync.WaitGroup
}

type dispatcher struct {
	eventChannel chan *dispatchData
	events       map[string][]func(*Data)
}

var dispatcherInstance *dispatcher

func Dispatch(name string, data interface{}) {
	dispatcherInstance.dispatch(name, data)
}

func DispatchSync(name string, data interface{}) {
	dispatcherInstance.dispatchSync(name, data)
}

func PrepareDispatcher(events map[string][]func(*Data)) {
	dispatcherInstance := &dispatcher{
		make(chan *dispatchData),
		events,
	}
	dispatcherInstance.run()
}

func (d *dispatcher) run() {
	go func() {
		for {
			go func(dat *dispatchData) {
				for _, listenerFunc := range d.events[dat.event.Name] {
					listenerFunc(dat.event)
				}
				if dat.waitGroup != nil {
					dat.waitGroup.Done()
				}
			}(<-d.eventChannel)
		}
	}()
}

func (d *dispatcher) dispatch(name string, data interface{}) {
	d.eventChannel <- &dispatchData{&Data{name, data}, nil}
}

func (d *dispatcher) dispatchSync(name string, data interface{}) {
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(1)
	d.eventChannel <- &dispatchData{&Data{name, data}, waitGroup}
	waitGroup.Wait()
}