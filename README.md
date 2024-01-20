# echartsgoja

`echartsgoja` renders [Apache ECharts][echarts-examples] as SVGs using the
[`goja`][goja] JavaScript runtime. Developed for use by [`usql`][usql] for
rendering charts.

[Overview][] | [TODO][] | [About][]

[Overview]: #overview "Overview"
[TODO]: #todo "TODO"
[About]: #about "About"

[![Unit Tests][echartsgoja-ci-status]][echartsgoja-ci]
[![Go Reference][goref-echartsgoja-status]][goref-echartsgoja]
[![Discord Discussion][discord-status]][discord]

[echartsgoja-ci]: https://github.com/xo/echartsgoja/actions/workflows/test.yml
[echartsgoja-ci-status]: https://github.com/xo/echartsgoja/actions/workflows/test.yml/badge.svg
[goref-echartsgoja]: https://pkg.go.dev/github.com/xo/echartsgoja
[goref-echartsgoja-status]: https://pkg.go.dev/badge/github.com/xo/echartsgoja.svg
[discord]: https://discord.gg/yJKEzc7prt "Discord Discussion"
[discord-status]: https://img.shields.io/discord/829150509658013727.svg?label=Discord&logo=Discord&colorB=7289da&style=flat-square "Discord Discussion"

## Overview

Install in the usual Go fashion:

```sh
$ go get github.com/xo/echartsgoja@latest
```

Then use like the following:

```go
package echartsgoja_test

import (
	"context"
	"log"
	"os"

	"github.com/xo/echartsgoja"
)

func Example() {
	echarts := echartsgoja.New()
	svg, err := echarts.RenderOptions(context.Background(), simpleOpts)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile("simple.svg", []byte(svg), 0o644); err != nil {
		log.Fatal(err)
	}
	// Output:
}

const simpleOpts = `{
  "xAxis": {
    "type": "category",
    "boundaryGap": false,
    "data": ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"]
  },
  "yAxis": {
    "type": "value"
  },
  "series": [{
    "data": [820, 932, 901, 934, 1290, 1330, 1320],
    "type": "line",
    "areaStyle": {}
  }]
}`
```

## TODO

- Rewrite as native Go

## About

`echartsgoja` was written primarily to support these projects:

- [usql][usql] - a universal command-line interface for SQL databases

Users of this package may find the [`github.com/xo/resvg`][resvg] package
helpful in rendering generated SVGs.

[echarts-examples]: https://echarts.apache.org/examples/en/index.html
[usql]: https://github.com/xo/usql
[resvg]: https://github.com/xo/resvg
[goja]: https://github.com/dop251/goja
