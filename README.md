## Godot-like signals library

![Build Status](https://github.com/quasilyte/gsignal/workflows/Go/badge.svg)
[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/quasilyte/gsignal)](https://pkg.go.dev/mod/github.com/quasilyte/gsignal)

### Overview

A [Godot](https://docs.godotengine.org/en/stable/getting_started/step_by_step/signals.html)-inspired signals library for Go.

**Key features:**

* Amortized zero allocations in most use cases
* Efficient `Connect`, `Emit`, `Disconnect`
* Generic-based API gives us type safety and convenience

Some games that were built with this library (this list is incomplete):

* [Roboden](https://quasilyte.itch.io/roboden)
* [Assemblox](https://quasilyte.itch.io/assemblox)
* [Decipherism](https://quasilyte.itch.io/decipherism)
* [Retrowave City](https://quasilyte.itch.io/retrowave-city)
* [Autotanks](https://quasilyte.itch.io/autotanks)
* [Sinecord](https://quasilyte.itch.io/sinecord)
* [Learn Georgian](https://quasilyte.itch.io/georgian-trainer)

Why bother and use something like this?

* It reduces the objects coupling
* It's an elegant event listener solution for Go
* Signals are a familiar concept (Godot, Phaser, Qt, ...)

### Installation

```bash
go get github.com/quasilyte/gsignal
```

### Quick Start

```go
package main

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

func main() {
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
}
```

### Introduction

This concept is borrowed from [Godot signals](https://docs.godotengine.org/en/stable/getting_started/step_by_step/signals.html), but it also resembles [Signals and Slots from Qt](https://doc.qt.io/qt-6/signalsandslots.html).

In `gsignal` terms:

* `Signal` = a field inside a struct
* `Slot` = a function (or method value) bound to a signal
* `Disconnect` = remove bound function
* `Emit` = call all bound functions

This library disconnects **disposed** objects automatically. This is convenient when you connect a scene object to some `Event`. When object goes away from a scene (becomes disposed), there is no need to call its event handler anymore.

### Thread Safety Notice

This library never does any synchronization on its own. It's implied that event emitters and their subscribers are executed inside the same goroutine.

This is possible in the game context, but it may not be as easy to enforce in some other applications.

Let's imagine that you want to do a background task in a game and provide a signal-style event for its completion. You spawn a goroutine for the task, but you keep the code outside look like it's still single-threaded. All concurrent communication should be incapsulated in the object owning the `Event` object. When this object knows that this concurrent job is completed, it should emit the event and notify all the subscribers.
