package main

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
	"strconv"
	"strings"
	"time"
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
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type SimilarData struct {
	Error    bool                      `json:"error"`
	Message  string                    `json:"message,omitempty"`
	Response map[string][]SimilarItems `json:"response"`
}

type SimilarItems struct {
	Artist    string  `json:"artist_name"`
	ArtistURL string  `json:"artist_url"`
	ID        int     `json:"id"`
	IndexID   int     `json:"index_id"`
	Lyrics    string  `json:"lyrics"`
	Percent   float32 `json:"percentage"`
	Song      string  `json:"song_name"`
	URL       string  `json:"song_url"`
}

//JASON structs end here

// GUI structs start here
type ListElement struct {
	Element widget.Clickable
	Title   string
	ID      int
}

type ClassList struct {
	list     layout.List
	Items    []ListElement
	selected int
}

type SimilarElements struct {
	Element widget.Clickable
	Artist  string
	Song    string
	Similar []SimilarItems // might not need this also
}

type SimilarList struct {
	list        layout.List
	Items       SimilarElements
	selectedNum int
}

//type LyricElements struct{
//	Lyrics		[]string
//}
//
//type LyricList struct{
//	list 		layout.List
//	Items		LyricElements
//	selectedNum	int
//}
// GUI structs end here

var listControl ClassList
var songData SongData
var similarControl SimilarList
var appTheme *material.Theme

func setupList(data SongData) {
	for i, value := range data.Response.Result {
		dashIndex := strings.Index(value.Name, "-")
		song := value.Name[dashIndex+2:]
		listControl.Items = append(listControl.Items, ListElement{Title: song, ID: data.Response.Result[i].ID})
		if i == len(data.Response.Result)-1 {
			fmt.Println("setupList complete")
		}
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
	err = json.Unmarshal(rawData, &data) // decode body of byte[] into interface
	fmt.Println("querySearch complete")
	return data
}

func similarQuery(id int) []SimilarItems {
	api := "https://searchly.asuarez.dev/api/v1/similarity/by_song"
	similarURL := api + "?song_id=" + strconv.Itoa(id)

	client := http.Client{
		Timeout: 15 * time.Second,
	}

	response, err := client.Get(similarURL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		log.Fatal("Didn't get 200")
	}
	rawData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var simData SimilarData
	json.Unmarshal(rawData, &simData)

	var values []SimilarItems
	for _, value := range simData.Response {
		values = append(values, value...)
	}
	return values
}

func main() {
	// search something
	fmt.Println("Enter an artist or song name: ")
	var userInput string
	userInput = "angel"
	//fmt.Scanln(&userInput)

	// retrieve data from API and compile data into listControl
	songData = querySearch(userInput)
	setupList(songData)

	go startApp()
	fmt.Println("GUI program started")
	app.Main()

}

func startApp() {
	defer os.Exit(0)
	mainWindow := app.NewWindow(app.Size(unit.Value{V: 1000}, unit.Value{V: 400}))
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

func drawGUI(gContext layout.Context, theme *material.Theme) layout.Dimensions {
	listControl.list.Axis = layout.Vertical
	similarControl.list.Axis = layout.Vertical

	retLayout := layout.Flex{Axis: layout.Horizontal}.Layout(gContext,
		layout.Rigid(drawList(gContext, theme)),
		layout.Rigid(drawSimilarList(gContext, theme)),
		//layout.Flexed(1, drawLyrics(gContext, theme)),
	)
	return retLayout
}

//func drawLyrics(gContext layout.Context, theme *material.Theme) layout.Widget{
//	return func(gtx layout.Context) layout.Dimensions {
//		gContext.Constraints.Max.Y = 100
//		return similarControl.list.Layout(gtx, len(similarControl.Items.Similar), selectLyrics)
//	}
//}
//
//func selectLyrics(gtx layout.Context, selectedItem int) layout.Dimensions{
//	userSelection := &similarControl.Items
//	var editor = new(widget.Editor)
//	if userSelection.Element.Clicked(){
//		similarControl.selectedNum = selectedItem
//	}
//	return material.Editor(appTheme, editor, userSelection.Similar[selectedItem].Lyrics).Layout(gtx)
//}

func drawSimilarList(gContext layout.Context, theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return similarControl.list.Layout(gtx, len(similarControl.Items.Similar), selectSimilarItem)
	}
}

func selectSimilarItem(graphicsContext layout.Context, selectedItem int) layout.Dimensions {
	userSelection := &similarControl.Items
	if userSelection.Element.Clicked() {
		similarControl.selectedNum = selectedItem
		fmt.Println("similar data (song name): " + userSelection.Similar[selectedItem].Song + " selected")
	}
	var itemHeight int
	return layout.Stack{Alignment: layout.NW}.Layout(graphicsContext,
		layout.Stacked(
			func(gtx layout.Context) layout.Dimensions {
				dimensions := material.Clickable(gtx, &userSelection.Element,
					func(gtx layout.Context) layout.Dimensions {
						return layout.UniformInset(unit.Sp(6)).
							Layout(gtx, material.H6(appTheme,
								"Artist: "+userSelection.Similar[selectedItem].Artist+"\nSong: "+userSelection.Similar[selectedItem].Song).Layout)
					})
				itemHeight = dimensions.Size.Y
				return dimensions
			}),
		layout.Stacked(
			func(gtx layout.Context) layout.Dimensions {
				if similarControl.selectedNum != selectedItem {
					return layout.Dimensions{}
				}
				paint.ColorOp{Color: appTheme.Color.Hint}.Add(gtx.Ops)
				highlightWidth := gtx.Px(unit.Dp(4))
				paint.PaintOp{Rect: f32.Rectangle{
					Max: f32.Point{
						X: float32(highlightWidth),
						Y: float32(itemHeight),
					}}}.Add(gtx.Ops)
				return layout.Dimensions{Size: image.Point{X: highlightWidth, Y: itemHeight}}
			},
		),
	)
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
		similarControl.Items.Similar = similarQuery(listControl.Items[selectedItem].ID)
		fmt.Println("song data: " + userSelection.Title + " selected")
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
			func(gtx layout.Context) layout.Dimensions { //another one of those 'glorious anonymous functions
				if listControl.selected != selectedItem {
					return layout.Dimensions{} //if not selected - don't do anything special
				}
				paint.ColorOp{Color: appTheme.Color.Primary}.Add(gtx.Ops) //add a paint operation
				highlightWidth := gtx.Px(unit.Dp(4))                      //lets make it 4 device independent pixals
				paint.PaintOp{Rect: f32.Rectangle{                        //paint a rectangle using 32 bit floats
					Max: f32.Point{
						X: float32(highlightWidth),
						Y: float32(itemHeight),
					}}}.Add(gtx.Ops)
				return layout.Dimensions{Size: image.Point{X: highlightWidth, Y: itemHeight}}
			},
		),
	)
}
