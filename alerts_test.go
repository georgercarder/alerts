package alerts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type animal interface {
	speak() string
}

type cat struct{}

func (c *cat) speak() string {
	return "meow"
}

type dog struct{}

func (d *dog) speak() string {
	return "woof"
}

type bird struct{}

func (b *bird) speak() string {
	return "tweet"
}

func TestAlerts(t *testing.T) {

	//go G_Alerts().SendAlert(chanName, data)

	alertNames := []string{"cat", "dog", "bird"}

	var subs [](<-chan interface{})
	for _, an := range alertNames {
		s, err := G_Alerts().NewSubscription(an)
		if err != nil {
			assert.Equal(t, true, false, err.Error())
		}
		subs = append(subs, s)
	}
	doneCH := make(chan bool)
	go func() {
		for {
			a := <-subs[0]
			c := a.(*cat)
			assert.Equal(t, "meow", c.speak(), "cat err")
			doneCH <- true
		}
	}()
	go func() {
		for {
			a := <-subs[1]
			d := a.(*dog)
			assert.Equal(t, "woof", d.speak(), "dog err")
			doneCH <- true
		}
	}()
	go func() {
		for {
			a := <-subs[2]
			b := a.(*bird)
			assert.Equal(t, "tweet", b.speak(), "bird err")
			doneCH <- true
		}
	}()
	go G_Alerts().SendAlert(alertNames[0], &cat{})
	go G_Alerts().SendAlert(alertNames[1], &dog{})
	go G_Alerts().SendAlert(alertNames[2], &bird{})
	numFinished := 0
	for {
		select {
		case <-doneCH:
			numFinished++
			if numFinished < 3 {
				continue
			}
			return
		}
	}
}
