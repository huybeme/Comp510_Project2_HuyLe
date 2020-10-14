package main

/*
	Comp510
	Project 2 part 1
		Take in an input and use input for a search query. Pull data from API and decode.

	Project 2 part 2
		Take data from API and create GUI based window instead printing on console.

*/

import (
	"encoding/json"
	"fmt"
	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

//JSON structs start here - for whatever reason, SongData.go not being reached, get undefined error
type SongData struct {
	Error    bool     `json:"error"`
	Message  string   `json:"message,omitempty"`
	Response Response `json:"response"`
}

type Response struct {
	Result []Result `json:"results"`
}

type Result struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

//JASON structs end here

// GUI structs start here
type ListElement struct {
	Element     widget.Clickable
	Title       string
	Description string
}

type ClassList struct {
	list     layout.List
	Items    []ListElement
	selected int
}

// GUI structs end here

var listControl ClassList
var songData SongData

func setupList(data SongData) {
	for i, value := range data.Response.Result { // replace unused variable when description is available
		dashIndex := strings.Index(value.Name, "-")
		song := value.Name[dashIndex+2:]
		listControl.Items = append(listControl.Items, ListElement{Title: song, Description: string(i)})
		listControl.list.Axis = layout.Vertical
	}
}

func printData(data SongData) {
	for i := 0; i < len(data.Response.Result); i++ {
		fmt.Println(data.Response.Result[i].Name)
	}
}

func querySearch(userInput string) SongData {
	api := "https://searchly.asuarez.dev/api/v1/song/search"
	queryURl := api + "?query=" + userInput
	// retrieve data from API
	response, err := http.Get(queryURl)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()     // defer executes at end of function
	if response.StatusCode != 200 { //		.Body implements io.Reader
		log.Fatal("Didn't get 200") // status code 200 is ok success status response code
	}
	// read data from api and assign data to variable to contain []byte values
	rawData, err := ioutil.ReadAll(response.Body) // read response variable which contains API data
	if err != nil {
		log.Fatal(err)
	}
	// decode []byte data into struct
	var data SongData
	json.Unmarshal(rawData, &data) // decode body of byte[] into interface

	return data
}

func main() {

	// maybe needs to be in eventLoop for searching directly on GUI
	fmt.Println("Enter an artist or song name: ")
	var userInput string
	userInput = "cloud"
	//fmt.Scanln(&userInput)
	songData = querySearch(userInput)
	setupList(songData)
	printData(songData)

	go startApp()
	app.Main()
}

func startApp() {
	defer os.Exit(0)
	mainWindow := app.NewWindow(app.Size(unit.Value{V: 1400}, unit.Value{V: 400}))
	err := eventLoop(mainWindow)
	if err != nil {
		log.Fatal(err)
	}
}

func eventLoop(mainWindow *app.Window) (err error) {
	appTheme = material.NewTheme(gofont.Collection())
	var operationsQ op.Ops
	for {
		event := <-mainWindow.Events()
		switch eventType := event.(type) {
		case system.DestroyEvent:
			return eventType.Err
		case system.FrameEvent:
			graphicsContext := layout.NewContext(&operationsQ, eventType)
			drawGUI(graphicsContext, appTheme)
			eventType.Frame(graphicsContext.Ops)
		}
	}
}

var (
	//entryLine = &widget.Editor{			// maybe out of scope of proj 2 but maybe implement search button when theres free time.
	//	SingleLine: true,
	//	Submit: true,
	//}
	appTheme *material.Theme
	//searchButton = new(widget.Clickable)
)

func drawGUI(gContext layout.Context, theme *material.Theme) layout.Dimensions {
	//in := layout.UniformInset(unit.Dp(8))

	retLayout := layout.Flex{Axis: layout.Horizontal}.Layout(gContext,
		layout.Rigid( // clickable
			func(gtx layout.Context) layout.Dimensions {
				border := widget.Border{Width: unit.Px(100)}
				return border.Layout(gtx, drawList(gContext, theme))
			}),
		layout.Flexed(1, drawDisplay(gContext, theme)), // description
	//layout.Rigid(												// entry bar
	//	func(gtx layout.Context) layout.Dimensions {
	//		e := material.Editor(appTheme, entryLine, "Enter a Search Entry")
	//		e.Font.Style = text.Italic
	//		border := widget.Border{Color: color.RGBA{A: 0xff}, CornerRadius: unit.Dp(8), Width: unit.Px(5)}
	//		return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
	//			return layout.UniformInset(unit.Dp(8)).Layout(gtx, e.Layout)
	//		})
	//	}),
	//layout.Rigid(func(gtx layout.Context) layout.Dimensions {	// button
	//	return in.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
	//		for searchButton.Clicked() {
	//			fmt.Println("you clicked the button")
	//		}
	//		return material.Button(appTheme, searchButton, "Search").Layout(gtx)
	//	})
	//}),
	)
	return retLayout
}

func drawDisplay(gContext layout.Context, theme *material.Theme) layout.Widget {
	return func(ctx layout.Context) layout.Dimensions {
		displayText := material.Body1(theme, listControl.Items[listControl.selected].Title) // change output here to description
		return layout.Center.Layout(ctx, displayText.Layout)                                // description currently empty
	}
}

func drawList(gContext layout.Context, theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return listControl.list.Layout(gtx, len(listControl.Items), selectItem)
	}
}

func selectItem(graphicsContext layout.Context, selectedItem int) layout.Dimensions {
	userSelection := &listControl.Items[selectedItem]
	if userSelection.Element.Clicked() {
		listControl.selected = selectedItem
	}
	var itemHeight int
	return layout.Stack{Alignment: layout.W}.Layout(graphicsContext,
		layout.Stacked(
			func(gtx layout.Context) layout.Dimensions {
				dimensions := material.Clickable(gtx, &userSelection.Element,
					func(gtx layout.Context) layout.Dimensions {
						return layout.UniformInset(unit.Sp(6)).
							Layout(gtx, material.H6(appTheme, userSelection.Title).Layout)
					})
				itemHeight = dimensions.Size.Y
				return dimensions
			}),
		layout.Stacked(
			func(gtx layout.Context) layout.Dimensions {
				if listControl.selected != selectedItem {
					return layout.Dimensions{}
				}
				paint.ColorOp{Color: appTheme.Color.Primary}.Add(gtx.Ops)
				highlightWidth := gtx.Px(unit.Dp(4))
				paint.PaintOp{Rect: f32.Rectangle{
					Max: f32.Point{
						X: float32(highlightWidth),
						Y: float32(itemHeight),
					}}}.Add(gtx.Ops)
				return layout.Dimensions{Size: image.Point{X: highlightWidth, Y: itemHeight}}
			}),
	)
}
