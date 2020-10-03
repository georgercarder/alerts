package alerts

import (
	"fmt"
	"sync"
	"time"

	mi "github.com/georgercarder/mod_init"
)

func G_Alerts() (a *Alerts) {
	aa, err := modInitializeAlerts.Get()
	if err != nil {
		//fmt.Fprintf(os.stdErr, err) // TODO
		return
	}
	return aa.(*Alerts)
}

const ModInitTimeout = 3 * time.Second // plenty of time

var modInitializeAlerts = mi.NewModInit(NewAlertsHub,
	ModInitTimeout, fmt.Errorf("*Alerts init error"))

type Alerts struct {
	sync.RWMutex
	Name2Chan map[string]*InterfaceChan
}

func NewAlertsHub() (a interface{}) { // *Alerts
	aa := new(Alerts)
	aa.Init()
	a = aa
	return
}

func (a *Alerts) SendAlert(chanName string, data interface{}) {
	if a.Name2Chan[chanName] == nil {
		a.Lock()
		channel := make(chan interface{})
		a.newChanLocked(chanName, channel)
		a.Unlock()
	}
	a.RLock()
	a.Name2Chan[chanName].CH <- data
	a.RUnlock()
	return
}

func (a *Alerts) Init() {
	a.Lock()
	defer a.Unlock()
	a.Name2Chan = make(map[string]*InterfaceChan)
}

func (a *Alerts) newChanLocked(
	chanName string, CH chan interface{}) (er error) {
	ch := new(InterfaceChan)
	ch.Init(CH)
	a.Name2Chan[chanName] = ch
	return
}

func (a *Alerts) NewSubscription(
	chanName string) (subCH <-chan interface{}, er error) {
	a.Lock()
	defer a.Unlock()
	if a.Name2Chan[chanName] == nil {
		channel := make(chan interface{})
		a.newChanLocked(chanName, channel)
	}
	ch := a.Name2Chan[chanName]
	subCH = ch.NewSubscription()
	return
}

type Chan interface {
	NewSubscription()
}

type InterfaceChan struct {
	sync.RWMutex
	CH            (chan interface{})
	Subscriptions [](chan interface{})
}

func (c *InterfaceChan) Init(CH chan interface{}) {
	c.Lock()
	defer c.Unlock()
	c.CH = CH
	c.fanout()
	return
}

// this is effectively creates a fanout on InterfaceChan.CH
func (c *InterfaceChan) NewSubscription() <-chan interface{} {
	c.Lock()
	subscription := make(chan interface{})
	c.Subscriptions = append(c.Subscriptions, subscription)
	c.Unlock()
	return subscription
}

func (c *InterfaceChan) fanout() {
	go func() {
		for {
			select {
			case s := <-c.CH:
				c.RLock()
				for i := 0; i < len(c.Subscriptions); i++ {
					c.Subscriptions[i] <- s
				}
				c.RUnlock()
			}
		}
	}()
	return
}
