//go:build js && wasm

package gojsdomframe

import (
	"fmt"
	"log"
	"syscall/js"
	"time"

	dom "honnef.co/go/js/dom/v2"
)

const serverURL = "http://127.0.0.1:8080"

var doc = dom.GetWindow().Document()

func importScript(body dom.Element, url string) {
	el := doc.CreateElement("script").(*dom.HTMLScriptElement)
	el.SetSrc(url)
	body.AppendChild(el)
}

func CreateFrame() {
	if doc == nil { // the JS way of testing for nil
		panic("unable to get 'document' object")
	}

	// try to get the 'body' property of JS from the global scope
	body := doc.GetElementsByTagName("body")[0]
	if body == nil {
		panic("unable to get 'body' object")
	}

	// importScript(body, "https://ajax.googleapis.com/ajax/libs/jquery/3.4.1/jquery.min.js")
	// importScript(body, "https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js")
	// importScript(body, "https://www.gstatic.com/charts/loader.js")

	root := CreateElementDiv("wrapper")
	root.AppendChild(NewHeaderDiv())
	root.AppendChild(NewMainDiv())
	root.AppendChild(NewAside1Div())
	root.AppendChild(NewFooterDiv())
	body.AppendChild(root)
}

func CreateElementDiv(classNames string) *dom.HTMLDivElement {
	el := doc.CreateElement("div").(*dom.HTMLDivElement)
	el.SetClass(classNames)
	return el
}

func NewMainDiv() *dom.HTMLDivElement {
	if doc == nil { // the JS way of testing for nil
		panic("unable to get 'document' object")
	}
	main := CreateElementDiv("main")
	main.AppendChild(newHello())
	main.AppendChild(newChart())

	return main
}

func NewAside1Div() *dom.HTMLDivElement {
	if doc == nil { // the JS way of testing for nil
		panic("unable to get 'document' object")
	}
	div := CreateElementDiv("aside aside-1")
	div.AppendChild(newNav())

	return div
}

func NewHeaderDiv() *dom.HTMLDivElement {
	if doc == nil { // the JS way of testing for nil
		panic("unable to get 'document' object")
	}
	div := CreateElementDiv("header")

	return div
}

func NewFooterDiv() *dom.HTMLDivElement {
	if doc == nil { // the JS way of testing for nil
		panic("unable to get 'document' object")
	}
	div := CreateElementDiv("footer")

	return div
}

func newHello() dom.Element {
	el := doc.CreateElement("h1")
	msg := fmt.Sprintf("Hello from WebAssembly! (made with 💖 Go) (Cliecked %v)", 0)
	el.SetInnerHTML(msg)

	cb := func(ev dom.Event) {

		fmt.Println("button clicked v2")
		msg := fmt.Sprintf("Hello from WebAssembly! (made with 💖 Go) (Cliecked %v)", 0)
		// element.SetAttribute("innerHTML", msg)
		el.SetInnerHTML(msg)
		//cb.Release() // release the function if the button will not be clicked again
	}
	el.AddEventListener("click", false, cb)

	return el
}

func InitializeData(this js.Value, inputs []js.Value) any {
	google := js.Global().Get("google")
	if google.IsUndefined() {
		log.Fatalf("Invalide google obj: %v", google)
	}
	// data = new google.visualization.DataTable();
	data := google.Get("visualization")
	if data.IsUndefined() {
		log.Fatalf("Invalide data1 obj: %v", google)
	}
	data = data.Get("DataTable").New()
	if data.IsUndefined() {
		log.Fatalf("Invalide data obj: %v", google)
	}
	// data.addColumn('datetime', 'Time');
	data.Call("addColumn", "datetime", "Time")
	// data.addColumn('number', '% CPU usage');
	data.Call("addColumn", "number", "% CPU usage")
	// return data;
	js.Global().Set("data", data)
	return nil
}

type Opt map[string]interface{}

func newChart() dom.Element {
	// google.charts.load('current', {packages: ['corechart', 'line']});
	google := js.Global().Get("google")
	if google.IsUndefined() {
		log.Fatalf("Invalide google obj: %v", google)
	}

	pkgOpt := map[string]interface{}{
		"packages": []interface{}{"corechart", "line"},
	}
	//pkgOpt := `{packages: ['corechart', 'line']}`
	google.Get("charts").Call("load", "current", pkgOpt)
	// google.setOnLoadCallback(InitializeData);
	initCb := js.FuncOf(InitializeData)
	google.Call("setOnLoadCallback", initCb)

	div := CreateElementDiv("display")
	div.SetID("display")

	el := doc.CreateElement("button")
	el.SetInnerHTML("Show")

	div.AppendChild(el)

	cb := func(ev dom.Event) {

		display := js.Global().Get("document").Call("getElementById", "display")
		if display.IsUndefined() || display.IsNull() {
			log.Fatalf("Invalide display obj: %v", display)
		}
		log.Printf("Show display obj: %v\n", display.Get("id"))
		// var chart = new google.visualization.LineChart(document.getElementById('display'));
		chart := google.Get("visualization").Get("LineChart").New(display)

		options := map[string]interface{}{
			"hAxis": map[string]interface{}{
				"title":  "Time",
				"format": "HH:mm:ss",
			},
			"vAxis": map[string]interface{}{
				"title":    "% CPU Usage",
				"minValue": 0,
			},
			"colors": []interface{}{"#a52714"},
			"crosshair": map[string]interface{}{
				"color":   "#000",
				"trigger": "selection",
			},
			"dateFormat": "HH:mm:ss",
		}
		// chart.draw(data, options);
		data := js.Global().Get("data")
		now := time.Now().Unix()
		// Get the Date object constructor from JavaScript
		dateConstructor := js.Global().Get("Date")
		// Return a new JS "Date" object with the time from the Go "now" variable
		// We're passing the UNIX timestamp to the "Date" constructor
		// Because JS uses milliseconds for UNIX timestamp, we need to multiply the timestamp by 1000
		rows := []interface{}{}
		var i int64
		for i = 0; i < 10; i++ {
			dt := dateConstructor.New((now + i*1000) * 1000)
			row := []interface{}{dt, 10 + i*10}
			rows = append(rows, row)
		}
		data.Call("addRows", rows)
		chart.Call("draw", data, options)
	}
	el.AddEventListener("click", false, cb)

	return div
}

type TaskID string
type TaskStatus string

type ReplyGetTaskStatus struct {
	Id     TaskID
	Status TaskStatus
}

func newNav() dom.Element {
	el := doc.CreateElement("nav")

	ul := doc.CreateElement("ul")
	el.AppendChild(ul)

	addNavLi := func(href, text string) {
		li := doc.CreateElement("li")
		ul.AppendChild(li)
		a1 := doc.CreateElement("a").(*dom.HTMLAnchorElement)
		a1.SetHref(href)
		a1.SetInnerHTML(text)
		li.AppendChild(a1)
	}
	addNavLi("#section1", "Section 1")
	addNavLi("#section2", "Section 2")

	return el
}
