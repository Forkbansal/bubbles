package list

import (
	"github.com/muesli/reflow/ansi"
	"strings"
	"testing"
)

// test is a shorthand and will be converted to proper testModels
// with genModels
type test struct {
	vWidth   int
	vHeight  int
	items    []string
	shouldBe string
}

type testModel struct {
	model        Model
	shouldBe     string
	afterMethode string
}

// TestViewBounds is use to make sure that the Renderer String
// NEVER leaves the bounds since then it could mess with the layout.
func TestViewBounds(t *testing.T) {
	for _, testM := range genModels(genTestModels()) {
		for i, line := range strings.Split(testM.model.View(), "\n") {
			lineWidth := ansi.PrintableRuneWidth(line)
			width := testM.model.Width
			if lineWidth > width {
				t.Errorf("The line:\n\n%s\n%s^\n\n is %d chars longer than the Viewport width.", line, strings.Repeat(" ", width-1), lineWidth-width)
			}
			if i > testM.model.Height {
				t.Error("There are more lines produced from the View() than the Viewport height")
			}
		}
	}
}

// TestGoldenSamples checks the View's string result against a knowen string (golden sample)
// Because there is no margin for diviations, if the test fails, lock also if the "golden sample" is sane.
func TestGoldenSamples(t *testing.T) {
	for _, testM := range genModels(genTestModels()) {
		actual := testM.model.View()
		expected := testM.shouldBe
		if actual != expected {
			t.Errorf("expected Output:\n\n%s\n\nactual Output:\n\n%s\n\n", expected, actual)
		}
	}
}

// TestPanic is also a golden sampling, but for cases that should panic.
func TestPanic(t *testing.T) {
	for _, testM := range genModels(genPanicTests()) {
		panicRes := make(chan interface{})
		go func(resChan chan<- interface{}) {
			defer func() { resChan <- recover() }() // Why does this Yield "%!s(<nil>)"?
			testM.model.View()
		}(panicRes)
		actual := <-panicRes
		expected := testM.shouldBe
		if actual != expected {
			t.Errorf("expected panic Output:\n\n%s\n\nactual Output:\n\n%s\n\n", expected, actual)
		}
	}
}

// TestDynamic tests the view output after a movement/view-changing method
func TestDynamic(t *testing.T) {
	for _, test := range genDynamicModels() {
		actual := test.model.View()
		expected := test.shouldBe
		if actual != expected {
			t.Errorf("expected Output, after Methode '%s' called:\n\n%s\n\nactual Output:\n\n%s\n\n", test.afterMethode, expected, actual)
		}
	}
}

// genModels embeds the fields from the rawModels into an actual model
func genModels(rawLists []test) []testModel {
	processedList := make([]testModel, len(rawLists))
	for i, list := range rawLists {
		m := NewModel()
		m.Height = list.vHeight
		m.Width = list.vWidth
		m.AddItems(MakeStringerList(list.items))
		newItem := testModel{model: m, shouldBe: list.shouldBe}
		processedList[i] = newItem
	}
	return processedList
}

// small helper function to generate simple test cases.
// for more elaborate ones append them afterwards.
func genTestModels() []test {
	return []test{
		// The default has abs linenumber and this seperator enabled
		// so that even if the terminal does not support colors
		// all propertys are still distinguishable.
		{
			240,
			80,
			[]string{
				"",
			},
			"\x1b[7m0 ╭>\x1b[0m",
		},
		// if exceding the boards and softwrap (at word bounderys are possible
		// wrap there. Dont increment the item number because its still the same item.
		{
			10,
			2,
			[]string{
				"robert frost",
			},
			"\x1b[7m0 ╭>robert\x1b[0m\n\x1b[7m  │ frost\x1b[0m",
		},
		{
			10,
			10,
			[]string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
			"\x1b[7m0 ╭>\x1b[0m\n" +
				`1 ╭ 
2 ╭ 
3 ╭ 
4 ╭ 
5 ╭ 
6 ╭ 
7 ╭ 
8 ╭ 
9 ╭ `,
		},
	}
}

// genPanicTests generats test cases that should panic with the shouldBe string
func genPanicTests() []test {
	return []test{
		// no width to display -> panic
		{
			0,
			1,
			[]string{""},
			"Can't display with zero width or hight of Viewport",
		},
		// no height to display -> panic
		{
			1,
			0,
			[]string{""},
			"Can't display with zero width or hight of Viewport",
		},
		// no item to display -> panic TODO handel/think-about this case
		//{
		//	1,
		//	1,
		//	[]string{},
		//	"",
		//},
	}
}

// genDynamicModels generats test cases for dynamic actions like movement, sorting, resizing
func genDynamicModels() []testModel {
	moveBottom := NewModel()
	moveBottom.Width = 10
	moveBottom.Height = 10
	moveBottom.AddItems(MakeStringerList([]string{"", "", "", ""}))
	moveBottom.Bottom()
	moveDown := NewModel()
	moveDown.Height = 50
	moveDown.Width = 80
	moveDown.AddItems(MakeStringerList([]string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""}))
	moveDown.curIndex = 45 // set cursor next to line Offset Border so that the down move, should move the hole visible area.
	moveDown.Move(1)
	return []testModel{
		{model: moveBottom,
			shouldBe:     "0 ╭ \n1 ╭ \n2 ╭ \n\x1b[7m3 ╭>\x1b[0m",
			afterMethode: "Bottom",
		},
		{model: moveDown,
			shouldBe: ` 1 ╭ 
 2 ╭ 
 3 ╭ 
 4 ╭ 
 5 ╭ 
 6 ╭ 
 7 ╭ 
 8 ╭ 
 9 ╭ 
10 ╭ 
11 ╭ 
12 ╭ 
13 ╭ 
14 ╭ 
15 ╭ 
16 ╭ 
17 ╭ 
18 ╭ 
19 ╭ 
20 ╭ 
21 ╭ 
22 ╭ 
23 ╭ 
24 ╭ 
25 ╭ 
26 ╭ 
27 ╭ 
28 ╭ 
29 ╭ 
30 ╭ 
31 ╭ 
32 ╭ 
33 ╭ 
34 ╭ 
35 ╭ 
36 ╭ 
37 ╭ 
38 ╭ 
39 ╭ 
40 ╭ 
41 ╭ 
42 ╭ 
43 ╭ 
44 ╭ 
45 ╭ ` +
				"\n\x1b[7m46 ╭>\x1b[0m\n" +
				`47 ╭ 
48 ╭ 
49 ╭ 
50 ╭ `,
			afterMethode: "Down",
		},
	}
}

func TestCheckBorder(t *testing.T) {
	m := NewModel()
	m.AddItems(MakeStringerList([]string{"", "", "", ""}))
	if !m.CheckWithinBorder(0) {
		t.Errorf("zero is not out of border")
	}
	if !m.CheckWithinBorder(len(m.listItems) - 1) {
		t.Errorf("lasitem is not out of border")
	}
	if m.CheckWithinBorder(-1) {
		t.Errorf("-1 is out of border")
	}
	if m.CheckWithinBorder(len(m.listItems)) {
		t.Errorf("len(list) is out of border")
	}
}
