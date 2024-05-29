# go-truecoach [![Go Reference](https://pkg.go.dev/badge/github.com/tobiassjosten/go-truecoach.svg)](https://pkg.go.dev/github.com/tobiassjosten/go-truecoach) [![Codacy Badge](https://app.codacy.com/project/badge/Grade/3d28f7a6c22c4d029e2309dde5b11edf)](https://app.codacy.com/gh/tobiassjosten/go-truecoach/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)

An unofficial Go library for integrating with [TrueCoach](https://truecoach.co/).

Table of contents:

*   [Installation](#installation)
*   [Quick start](#quick-start)
*   [Limitations](#limitations)
*   [Contribute](#contribute)

## Installation

1.  Download the module:

    go get -u github.com/tobiassjosten/go-truecoach

2.  Import it in your project:

    import "github.com/tobiassjosten/go-truecoach"

## Quick start

```go
package main

import (
	"fmt"
	"github.com/tobiassjosten/go-truecoach"
)

func main() {
	tc := truecoach.NewService("SuperSecret123!!!")

	clients, err := tc.Clients()
	if err != nil {
		panic(err)
	}

	for _, client := range clients {
		metrics, err := tc.Metrics(client.ID)
		if err != nil {
			panic(err)
		}

		for _, metric := range metrics {
			for _, assessment := range metric.Assessments {
				for _, sample := range assessment.Samples {
					fmt.Println(sample)
				}
			}
		}
	}
}
```

## Limitations

This library was developed by reverse engineering the TrueCoach API, based on their JavaScript frontend. As such, nothing is guaranteed to work and, even if it does, no future promises are made.

(Although I've been using it myself with zero problems for several years.)

## Contribute

Feel free to [create a ticket](https://github.com/tobiassjosten/go-truecoach/issues/new) if you want to discuss or suggest something. I'd love your input and will happily work with you to cover your use cases and explore your ideas for improvements.

Changes can be suggested directly by [creating a pull request](https://github.com/tobiassjosten/go-truecoach/compare) but I'd recommend starting an issue first, so you don't end up wasting your time with something I end up rejecting.

### Contributors

*   [Tobias Sjösten](https://github.com/tobiassjosten)
