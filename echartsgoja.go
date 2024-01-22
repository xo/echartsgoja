// Package echartsgoja renders [Apache ECharts] as SVGs using the [goja]
// JavaScript runtime. Developed for use by [usql] for rendering charts.
//
// [Apache ECharts]: https://echarts.apache.org/examples/en/index.html
// [goja]: https://github.com/dop251/goja
// [usql]: https://github.com/xo/usql
package echartsgoja

import (
	"context"
	"embed"
	"fmt"
	"regexp"
	"sync"

	"github.com/dop251/goja"
	gojaparser "github.com/dop251/goja/parser"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/eventloop"
	"github.com/dop251/goja_nodejs/require"
)

// ECharts handles rendering [Apache ECharts] as SVGs, wrapping a [goja]
// runtime vm.
//
// [Apache ECharts]: https://echarts.apache.org/examples/en/index.html
// [goja]: https://github.com/dop251/goja
type ECharts struct {
	logger        func(...interface{})
	echartsJs     *goja.Program
	echartsgojaJs *goja.Program
	width         int
	height        int
	once          sync.Once
	err           error
}

// New creates a new echarts instance.
func New(opts ...Option) *ECharts {
	vm := &ECharts{
		width:  400,
		height: 300,
	}
	for _, o := range opts {
		o(vm)
	}
	return vm
}

// init initializes the echarts instance.
func (vm *ECharts) init() error {
	vm.once.Do(func() {
		if vm.echartsJs, vm.err = compileEmbeddedScript("echarts.min.js"); vm.err != nil {
			return
		}
		if vm.echartsgojaJs, vm.err = compileEmbeddedScript("echartsgoja.js"); vm.err != nil {
			return
		}
	})
	return vm.err
}

// run runs the embedded javascripts on the goja runtime, and exports a symbol
// name to v.
func (vm *ECharts) run(r *goja.Runtime, name string, v interface{}) error {
	if _, err := r.RunProgram(vm.echartsJs); err != nil {
		return err
	}
	if _, err := r.RunProgram(vm.echartsgojaJs); err != nil {
		return err
	}
	if name != "" {
		return r.ExportTo(r.Get(name), v)
	}
	return nil
}

// loop instantiates a new loop, and waits for it to terminate, or until the
// context is closed.
func (vm *ECharts) loop(ctx context.Context, f func(*goja.Runtime) error) error {
	// init
	if err := vm.init(); err != nil {
		return err
	}
	mod := console.RequireWithPrinter(vm)
	registry := new(require.Registry)
	registry.RegisterNativeModule(console.ModuleName, mod)
	loop := eventloop.NewEventLoop(
		eventloop.WithRegistry(registry),
	)
	errch := make(chan error, 1)
	go func() {
		loop.Run(func(r *goja.Runtime) {
			defer close(errch)
			errch <- f(r)
		})
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errch:
		return err
	}
}

// Version returns the embedded echarts version.
func (vm *ECharts) Version() (string, error) {
	if err := vm.init(); err != nil {
		return "", err
	}
	var f func() (string, error)
	if err := vm.run(goja.New(), "version", &f); err != nil {
		return "", err
	}
	return f()
}

// RenderOptions renders a echarts options as a SVG.
func (vm *ECharts) RenderOptions(ctx context.Context, opts string) (res string, err error) {
	defer func() {
		if e := recover(); e != nil {
			if ex, ok := e.(*goja.Exception); ok {
				res, err = "", ex.Unwrap()
			} else {
				res, err = "", fmt.Errorf("recovered from: %v", e)
			}
		}
	}()
	err = vm.loop(ctx, func(r *goja.Runtime) error {
		var render renderOptionsFunc
		if err := vm.run(r, "render_options", &render); err != nil {
			return err
		}
		var err error
		if res, err = render(vm.width, vm.height, opts); err != nil {
			return err
		}
		return nil
	})
	return
}

// RenderScript renders a echarts script as a SVG.
func (vm *ECharts) RenderScript(ctx context.Context, src string) (res string, err error) {
	// compile
	var p *goja.Program
	if p, err = compile("script", fmt.Sprintf(renderScript, cleanRE.ReplaceAllString(src, ""))); err != nil {
		return
	}
	defer func() {
		if e := recover(); e != nil {
			if ex, ok := e.(*goja.Exception); ok {
				res, err = "", ex.Unwrap()
			} else {
				res, err = "", fmt.Errorf("recovered from: %v", e)
			}
		}
	}()
	var opts string
	err = vm.loop(ctx, func(r *goja.Runtime) error {
		if err := vm.run(r, "", nil); err != nil {
			return err
		}
		v, err := r.RunProgram(p)
		if err != nil {
			return err
		}
		buf, err := v.ToObject(r).MarshalJSON()
		if err != nil {
			return err
		}
		opts = string(buf)
		return nil
	})
	switch {
	case err != nil:
		return "", err
	case opts == "", opts == "{}":
		return "", ErrCouldNotExecuteScript
	}
	return vm.RenderOptions(ctx, opts)
}

// Log satisfies the [goja.Printer] interface.
func (vm *ECharts) Log(s string) {
	vm.log([]string{"LOG", s})
}

// Warn satisfies the [goja.Printer] interface.
func (vm *ECharts) Warn(s string) {
	vm.log([]string{"WARN", s})
}

// Error satisfies the [goja.Printer] interface.
func (vm *ECharts) Error(s string) {
	vm.log([]string{"ERROR", s})
}

// log is the script callback for logging a message.
func (vm *ECharts) log(s []string) {
	if vm.logger != nil {
		v := make([]interface{}, len(s))
		for i, ss := range s {
			v[i] = ss
		}
		vm.logger(v...)
	}
}

// Option is a echarts option.
type Option func(*ECharts)

// WithLogger is a echarts option to set the logger.
func WithLogger(logger func(...interface{})) Option {
	return func(vm *ECharts) {
		vm.logger = logger
	}
}

// WithWidthHeight is a echarts option to set the width, height.
func WithWidthHeight(width, height int) Option {
	return func(vm *ECharts) {
		vm.width, vm.height = width, height
	}
}

// Error is a error.
type Error string

// Errors.
const (
	// ErrCouldNotExecuteScript is the could not execute script error.
	ErrCouldNotExecuteScript Error = "could not execute script"
)

// Error satisfies the [error] interface.
func (err Error) Error() string {
	return string(err)
}

// compile compiles a goja program.
func compile(name, src string) (*goja.Program, error) {
	prg, err := goja.Parse(name, src, gojaparser.WithDisableSourceMaps)
	if err != nil {
		return nil, fmt.Errorf("unable to parse %s: %w", name, err)
	}
	p, err := goja.CompileAST(prg, true)
	if err != nil {
		return nil, fmt.Errorf("unable to compile %s: %w", name, err)
	}
	return p, nil
}

// compileEmbeddedScript compiles the embedded script as a goja program.
func compileEmbeddedScript(name string) (*goja.Program, error) {
	buf, err := jsScripts.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("unable to load %s: %w", name, err)
	}
	return compile(name, string(buf))
}

// cleanRE is a regexp for removing export {} from scripts.
var cleanRE = regexp.MustCompile(`export\s+\{\};`)

// renderScript is the render script.
const renderScript = `(function() {
  let exports = {};
  let app = {};
  let option = {};
  %s
  return option;
})();`

// renderOptionsFunc is the signature for the render_options func.
type renderOptionsFunc func(width, height int, opts string) (string, error)

// versionTxt is the embedded version.txt.
//
//go:embed version.txt
var versionTxt string

// jsScripts are the embedded echarts javascripts.
//
//go:embed *.js
var jsScripts embed.FS
