# Uber Fx

This document serves as a quick reference for [Uber Fx](https://github.com/uber-go/fx) framework used in ASTA Fx. It
goes through parts of Fx that is required to understand ASTA Fx library.

## Prerequisite

You will need the following before getting started.

- Basic understanding on dependency injection pattern
- Understanding on Go struct and interfaces

## Introduction

Uber Fx is a dependency injection framework for Go. It helps to eliminate manual dependency management and facilitate
code reuse.

## Minimal Application

Below shows a minimal application with Fx framework. It starts up the application without any dependency provided.

```go
package main

import "go.uber.org/fx"

func main() {
	// fx.New() returns *fx.App pointing to a new Fx app
	app := fx.New()
	// (*fx.App).Run() starts up and blocks on signals channel
	// gracefully shuts down after receiving interrupt
	app.Run()
}
```

You will see output similar to the following.

```
[Fx] PROVIDE    fx.Lifecycle <= go.uber.org/fx.New.func1()
[Fx] PROVIDE    fx.Shutdowner <= go.uber.org/fx.(*App).shutdowner-fm()
[Fx] PROVIDE    fx.DotGraph <= go.uber.org/fx.(*App).dotGraph-fm()
[Fx] RUNNING
```

The application starts up and wait for interrupt signal. You can see the Fx framework provides three objects by default.
`fx.Lifecycle` will be useful for building our apps, and it will be covered in later sections.

## Providers

To add providers (dependencies), pass constructors to `fx.Provide()` to create `fx.Option`. Pass `fx.Option` as
parameters when creating the app.

Below is an example of setting up `zap.logger` in a Fx application.

```go
package main

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewLogger a simple constructor that constructs a zap.Logger
func NewLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		// ignore error handling for time being
		panic(err)
	}
	return logger
}

func main() {
	// fx.New() accepts zero or more fx.Option
	app := fx.New(
		// creates a new fx.Option with logger constructor
		fx.Provide(NewLogger),
	)
	app.Run()
}
```

You will see one more line of output compared to previous example.

```
[Fx] PROVIDE    *zap.Logger <= main.NewLogger()
```

You may also throw an error in the constructor. Fx will log the error and stop the app from creating.

```go
func NewLogger() (*zap.Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return logger, nil
}
```

## TODO invoke

## TODO lifecycle hooks

## TODO decorators

## TODO modules

## TODO parameter object and result object

## TODO annotations