package event

type Dispatcher struct {
	subscribers []*Subscriber
}

func (d *Dispatcher) SetSubscribers(subscribers []*Subscriber) {
	for i := range subscribers {
		d.subscribers = append(d.subscribers, subscribers[i])
	}
}

func (d *Dispatcher) AddSubscriber(subscriber *Subscriber) {
	d.subscribers = append(d.subscribers, subscriber)
}

func (d *Dispatcher) Dispatch(eventName string) error {
	for i := range d.subscribers {
		err := d.subscribers[i].On(eventName)
		if err != nil {
			return err
		}
	}

	return nil
}
