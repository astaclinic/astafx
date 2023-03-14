# ASTA Fx

This repository contains common [uber-fx](https://github.com/uber-go/fx) modules for microservices.

## Design Goals

The modules aims to reduce boilerplate in setting up external dependencies (e.g. databases, webservers) and
injecting them in application code. Each module includes constructors, decorators (if it consists additional setup
functions) and lifecycle hooks (if it is long-running). 

## Usage

### Prerequisite

Please go through [Fx quick guide](docs/uberfx.md) if you are not familar with Fx.

### Importing Default Module Set

To import the default set of modules, add `astafx.Module` to fx app.

```go
package app

import (
	"go.uber.org/fx"
	"github.com/astaclinic/astafx"
)

func New() *fx.App {
	app := fx.New(
		astafx.Module,
		// include other fx modules here
	)
	return app
}
```
