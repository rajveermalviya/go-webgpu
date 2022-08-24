package main

import "github.com/rajveermalviya/gamen/display"

func main() {
	d, err := display.NewDisplay()
	if err != nil {
		panic(err)
	}
	defer d.Destroy()

	w, err := display.NewWindow(d)
	if err != nil {
		panic(err)
	}
	defer w.Destroy()

	w.SetCloseRequestedCallback(func() {
		d.Destroy()
	})

	for {
		if !d.Wait() {
			break
		}

		// we will render here
	}
}
