// Pre-config:
// git clone ....g3n-demo
// cd g3n-demo
// go install
// This is a simple model for your tests
package tests

import (
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/math32"
	"github.com/g3n/g3nd/app"
	"time"
)

// Sets the category and name of your test in the demos.Map
// The category name choosen here starts with a "|" so it shows as the
// last category in list. Change "model" to the name of your test.
func init() {
	app.DemoMap["|tests|.model"] = &testsModel{}
}

// This is your test object. You can store state here.
// By convention and to avoid conflict with other demo/tests name it
// using your test category and name.
type testsModel struct {
	grid *graphic.GridHelper // Pointer to a GridHelper created in 'Start'
}

// This method will be called once when the test is selected from the G3ND list.
// 'a' is a pointer to the G3ND application.
// It allows access to several methods such as a.Scene(), which returns the current scene,
// a.DemoPanel(), a.Camera(), a.Window() among others.
// You can build your scene adding your objects to the a.Scene()
func (t *testsModel) Start(a *app.App) {

	// Show axis helper
	ah := graphic.NewAxisHelper(1.0)
	a.Scene().Add(ah)

	// Creates a grid helper and saves its pointer in the test state
	t.grid = graphic.NewGridHelper(50, 1, &math32.Color{0.4, 0.4, 0.4})
	a.Scene().Add(t.grid)

	// Changes the camera position
	a.Camera().GetCamera().SetPosition(0, 4, 10)
	a.Camera().GetCamera().LookAt(&math32.Vector3{0, 0, 0})
}

// This method will be called at every frame
// You can animate your objects here.
func (t *testsModel) Update(a *app.App, deltaTime time.Duration) {

	// Rotate the grid, just for show.
	rps := float32(deltaTime.Seconds()) * 2 * math32.Pi
	t.grid.RotateY(rps * 0.05)
}

// Cleanup is called once at the end of the demo.
func (t *testsModel) Cleanup(a *app.App) {}
