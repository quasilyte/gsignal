package gsignal_test

import (
	"fmt"

	"github.com/quasilyte/gsignal"
)

type button struct {
	Name         string
	EventClicked gsignal.Event[*button]
}

func (b *button) Click() { b.EventClicked.Emit(b) }

type listener struct {
	disposed bool
}

func (l *listener) IsDisposed() bool { return l.disposed }

func (l *listener) onClick(b *button) {
	fmt.Println("listener on click")
}

func Example() {
	b := &button{Name: "example"}
	b.Click() // nothing happens, 0 connections

	i := 1
	b.EventClicked.Connect(nil, func(b *button) {
		fmt.Printf("%s clicked (%d)\n", b.Name, i)
		i++
	})
	b.Click() // prints "example clicked (1)" once
	b.Click() // prints "example clicked (2)" once; again

	l := &listener{}
	b.EventClicked.Connect(l, l.onClick)
	b.Click() // prints "example clicked (3)", then "listener on click"

	l.disposed = true // this will cause a disconnect
	b.Click()         // prints "example clicked (4)" once

	// Output:
	// example clicked (1)
	// example clicked (2)
	// example clicked (3)
	// listener on click
	// example clicked (4)
}
