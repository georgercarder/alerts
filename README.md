## alerts 

### Fanout a channel across a golang application


### Usage

Across an application, subscribe to several alerts

```
package likes_cats_1

// ...

catSub1, err := G_Alerts().NewSubscription("cats")
// handle err

// ...

// handle alert elsewhere ..

for {
	a := <-catSub1 
	c := a.(*cat)
	sound := c.speak()
}

```

```
package likes_cats_2 

// ...

catSub2, err := G_Alerts().NewSubscription("cats")
// handle err

// ...

// handle alert elsewhere ..

for {
	a := <-catSub2 
	c := a.(*cat)
	sound := c.speak()
}
```

```
package dog_people_are_better

// ...

dogSub, err := G_Alerts().NewSubscription("dogs")
// handle err

// ...

// handle alert elsewhere ..

for {
	a := <-dogSub
	d := a.(*dog)
	sound := d.speak()
	assert.Equal(t, true, IsBetter(sound), "I should be correct.")
}
```

```
package alert_master

// ...

kitty := newCat()

// ...

	go G_Alerts().SendAlert("cats", kitty)

// ...

fido := newDog()

// ...

	go G_Alerts().SendAlert("dogs", fido)

```


### Please let me know in git issue if you have questions, comments, or would like to contribute. 
