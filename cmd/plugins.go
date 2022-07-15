package main

type pluginManager struct {
	plugins []string
}

func newPluginManager(plugins []string) *pluginManager {
	d := new(pluginManager)
	d.plugins = make([]string, len(plugins))
	copy(d.plugins, plugins)
	return d
}

func (d *pluginManager) remove(plugin string) {
	var x []string
	for _, v := range d.plugins {
		if v != plugin {
			x = append(x, v)
		}
	}
	d.plugins = x
}

func (d *pluginManager) insertBefore(plugin, insertBefore string) {
	var x []string
	for _, v := range d.plugins {
		if v == insertBefore {
			x = append(x, plugin)
		}
		x = append(x, v)
	}
	d.plugins = x
}

func (d *pluginManager) get() []string {
	return d.plugins
}
