package directives

import "fmt"

type Manager struct {
	plugins []string
}

func NewDirectivesManager(plugins []string) *Manager {
	d := new(Manager)
	d.plugins = make([]string, len(plugins))
	copy(d.plugins, plugins)
	return d
}

func (d *Manager) Remove(plugin string) {
	var x []string
	for _, v := range d.plugins {
		if v != plugin {
			x = append(x, v)
		}
	}
	d.plugins = x
}

func (d *Manager) InsertBefore(plugin, insertBefore string) {
	var x []string
	for _, v := range d.plugins {
		if v == insertBefore {
			x = append(x, plugin)
		}
		x = append(x, v)
	}
	if len(d.plugins) != len(x)-1 {
		panic(fmt.Sprintf("%v doesn't exist", insertBefore))
	}
	d.plugins = x
}

func (d *Manager) Get() []string {
	return d.plugins
}
