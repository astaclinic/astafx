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

## Providing Dependencies

### Providers

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

You may also throw an error in the constructor. Fx will log the error and stop the app from creating if there is error
in the constructor.

```go
func NewLogger() (*zap.Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return logger, nil
}
```

### Dependency Injection

The constructors may accept arguments to specify their dependency. Fx will automatically inject required dependencies
to construct the object. For example, the logger constructor below requires a logger config which we need to provide
to Fx.

```go
package main

import (
	"encoding/json"
	"os"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	LogLevel zapcore.Level `json:"log_level"`
}

// NewConfig read config from filesystem and construct config object
func NewConfig() (*Config, error) {
	file, err := os.ReadFile("config.json")
	if err != nil {
		return nil, err
	}
	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// NewLogger accept config object and create logger
func NewLogger(config *Config) (*zap.Logger, error) {
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(config.LogLevel)
	logger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}
	return logger, nil
}

func main() {
	app := fx.New(
		fx.Provide(
			NewConfig, // provides config to Fx
			NewLogger, // Fx injects config to logger constructor
		),
	)
	app.Run()
}

```

Fx will automatically create a dependency graph and inject dependencies as needed.

### Invoke

One caveat of `fx.Provide()` is that the constructors are not invoked if there are no consumers. In the example below,
the log message "app starting" is not printed.

```go
package main

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	logger.Info("app starting")  // log message in constructor
	return logger
}

func main() {
	app := fx.New(
		fx.Provide(NewLogger),
	)
	app.Run()
}
```

To consume a dependency, `fx.Invoke()` should be used. Dependencies are only instantiated when needed.

```go
package main

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return logger
}

func LogOnStart(logger *zap.Logger) {
	logger.Info("app starting")
}

func main() {
	app := fx.New(
		fx.Provide(NewLogger),
		fx.Invoke(LogOnStart), // invoke LogOnStart with the provided logger
	)
	app.Run()
}
```

### Interface Types

For interface type, note that the dependencies are match with their interface type but not their underlying type. 
Also, pointers that implements certain interface will not be injected into the interface type. So the examples below
will not work.

```go
package main

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MyLog interface {
	Info(msg string, fields ...zap.Field)
}

func NewLogger() MyLog {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return logger
}

func LogOnStart(logger *zap.Logger) {
	logger.Info("app starting")
}

func main() {
	app := fx.New(
		fx.Provide(NewLogger),
		fx.Invoke(LogOnStart),
	)
	app.Run()
}
```

```
[Fx] ERROR              Failed to start: missing dependencies for function "main".LogOnStart
        /my-project/main.go:20:
missing type:
        - *zap.Logger (did you mean to use main.MyLog?)
```

```go
package main

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MyLog interface {
	Info(msg string, fields ...zap.Field)
}

func NewLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return logger
}

func LogOnStart(logger MyLog) {
	logger.Info("app starting")
}

func main() {
	app := fx.New(
		fx.Provide(NewLogger),
		fx.Invoke(LogOnStart),
	)
	app.Run()
}
```

```
[Fx] ERROR              Failed to start: missing dependencies for function "main".LogOnStart
        /my-project/main.go:20:
missing type:
        - main.MyLog (did you mean to use *zap.Logger?)
```

## Application Lifecycle

### Lifecycle Stages

The lifecycle of a Fx application has two high-level phases: initialization and execution. Both of these, in turn are
comprised of multiple steps.

#### Initialization

During initialization, Fx will:

- register all constructors passed to `fx.Provide`
- register all decorators passed to `fx.Decorate` (will be covered later)
- run all functions passed to `fx.Invoke`, calling constructors and decorators as needed

We have covered how to initialize dependencies in the above sections.

#### Execution

During execution, Fx will:

- run all startup hooks appended to the application by providers, decorators, and invoked functions
- wait for a signal to stop running
- run all shutdown hooks appended to the application

For long-lived service like http listeners, a startup hook and a shutdown hook should be added to start and stop the
listeners respectively.

### Lifecycle Hooks

To add lifecycle hooks, inject `fx.Lifecycle` and the service in an invoke function. Use `(*fx.Lifecycle).Append` to
add `fx.Hook` which has two fields `OnStart` and `OnStop` accepting `func(ctx context.Context) error`. The example below
shows how to use `http.Server` with lifecycle hook.

```go
package main

import (
	"context"
	"net"
	"net/http"

	"go.uber.org/fx"
)

func NewHTTPServer() *http.Server {
	return &http.Server{
		Addr: ":8080",
	}
}

func RunHttpServer(srv *http.Server, lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			go func() {
				err := srv.Serve(ln)
				if err != nil && err != http.ErrServerClosed {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
}

func main() {
	app := fx.New(
		fx.Provide(NewHTTPServer),
		fx.Invoke(RunHttpServer),
	)
	app.Run()
}
```

The above example shows how to use setup lifecycle hook for `http.Server`. In `OnStart` hook, the http server kicks off
and serve traffic. In `OnStop` hook, the http server shuts down.

```
[Fx] PROVIDE    *http.Server <= main.NewHTTPServer()
[Fx] PROVIDE    fx.Lifecycle <= go.uber.org/fx.New.func1()
[Fx] PROVIDE    fx.Shutdowner <= go.uber.org/fx.(*App).shutdowner-fm()
[Fx] PROVIDE    fx.DotGraph <= go.uber.org/fx.(*App).dotGraph-fm()
[Fx] INVOKE             main.RunHttpServer()
[Fx] HOOK OnStart               main.RunHttpServer.func1() executing (caller: main.RunHttpServer)
[Fx] HOOK OnStart               main.RunHttpServer.func1() called by main.RunHttpServer ran successfully in 335.833µs
[Fx] RUNNING
^C[Fx] INTERRUPT
[Fx] HOOK OnStop                main.RunHttpServer.func2() executing (caller: main.RunHttpServer)
[Fx] HOOK OnStop                main.RunHttpServer.func2() called by main.RunHttpServer ran successfully in 121.458µs
```

We can see the complete lifecycle in the log above. First, Fx register the provider `*http.Server`. Then, Fx invoke
`RunHttpServer()` to register lifecycle hooks. Afterwards, the application starts and `OnStart` hook is executed.
Finally, the `OnStop` hook is executed after receiving an interrupt signal.

### Minimal HTTP Server with Fx

## Advanced Fx Use Cases

### Decorators

### Modules

### Parameter object and Result object

### Annotations