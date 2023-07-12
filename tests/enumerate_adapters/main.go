package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func main() {
	instance := wgpu.CreateInstance(nil)
	defer instance.Release()

	adapter := findSuitableAdapter(instance)
	if adapter == nil {
		panic("couldn't find a suitable graphics adapter")
	}
	defer adapter.Release()

	fmt.Printf("selected: %s\n", prettify(adapter.GetProperties()))
}

func findSuitableAdapter(instance *wgpu.Instance) (selectedAdapter *wgpu.Adapter) {
	adapters := instance.EnumerateAdapters(nil)
	for i, adapter := range adapters {
		// logic for selecting the adapter goes here
		// can be done by checking properties & limits
		props := adapter.GetProperties()
		fmt.Printf("%d: %s\n", i, prettify(props))

		// for demostrating purposes selecting first one we get
		if i == 0 {
			selectedAdapter = adapter
			// DO NOT break after getting your desired adapter
			// because remaining other adapters MUST BE released.
		} else {
			// make sure to release un-selected adapters
			adapter.Release()
		}
	}
	return
}

func prettify(v any) string {
	var sb strings.Builder
	enc := json.NewEncoder(&sb)
	enc.SetIndent("", "\t")
	if err := enc.Encode(v); err != nil {
		panic(err)
	}
	return sb.String()
}
