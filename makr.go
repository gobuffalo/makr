package makr

import (
	"context"
	"os"

	"github.com/pkg/errors"
)

// Data to be passed into generators
type Data map[string]interface{}

// ShouldFunc decides whether a generator should be run or not
type ShouldFunc func(Data) bool

// Runnable interface must be implemented to be considered a runnable generator
type Runnable interface {
	Run(string, Data) error
}

// Generator is the top level construct that holds all of the Runnables
type Generator struct {
	Runners []Runnable
	Should  ShouldFunc
	Data    Data
}

// New Generator
func New() *Generator {
	return &Generator{
		Runners: []Runnable{},
		Data:    Data{},
	}
}

// Add a Runnable generator to the list
func (g *Generator) Add(r Runnable) {
	g.Runners = append(g.Runners, r)
}

// Run all of the generators
func (g *Generator) Run(rootPath string, data Data) error {
	dd := Data{}
	for k, v := range data {
		dd[k] = v
	}
	for k, v := range g.Data {
		dd[k] = v
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return chdir(rootPath, func() error {
		if g.Should != nil {
			b := g.Should(dd)
			if !b {
				return nil
			}
		}
		err := os.MkdirAll(rootPath, 0755)
		if err != nil {
			return errors.WithStack(err)
		}
		err = os.Chdir(rootPath)
		if err != nil {
			return errors.WithStack(err)
		}
		for _, r := range g.Runners {
			select {
			case <-ctx.Done():
				break
			default:
				err := r.Run(rootPath, dd)
				if err != nil {
					return errors.WithStack(err)
				}
			}
		}
		return nil
	})
}

func chdir(path string, fn func() error) error {
	pwd, _ := os.Getwd()
	defer os.Chdir(pwd)
	os.Chdir(path)
	return fn()
}

var nullShould = func(data Data) bool {
	return true
}
