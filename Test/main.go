package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"
	"time"
)

type Main struct {
	X      int
	Y      int
	Height int
	Width  int
}
type Input []struct {
	X      int
	Y      int
	Height int
	Width  int
}
type Bird struct {
	Main  Main
	Input Input
}
type Identification struct {
	X      int
	Y      int
	Width  int
	Height int
	Time   string
}

func isRectangleOverlap(rec1 []int, rec2 []int) bool {
	x1, x2, y1, y2 := rec1[0], rec1[0]+rec1[2], rec1[1], rec1[1]+rec1[3]
	x12, x22, y12, y22 := rec2[0], rec2[0]+rec2[2], rec2[1], rec2[1]+rec2[3]
	return x1 < x22 && x12 < x2 && y1 < y22 && y12 < y2

}

var bird Bird
var bird1 Bird
var idents []Identification
var idents1 []Identification
var idents2 []Identification

// Compile templates on start of the application
var templates = template.Must(template.ParseFiles("./upload.html"))
var t Bird

//Display the named template
func display(w http.ResponseWriter, page string, data interface{}) {
	templates.ExecuteTemplate(w, page+".html", data)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)

	//err = json.Unmarshal(body, &t)
	//if err != nil {
	//panic(err)
	//}
	//fmt.Println("shayad okeye")
	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	//defer file.Close()
	fmt.Println("Shoroo Be Khoondan Kon")
	//body, _ := file.Read()

	//fmt.Println("alan chi?")
	//log.Println(string(body))
	/*
		json.Unmarshal([]byte(body), &bird)
		fmt.Println("%v", bird.Main.X)
		fmt.Println("%v", bird.Main.Y)
		fmt.Println("%v", bird.Main.Width)
		fmt.Println("%v", bird.Main.Height)
	*/
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create file
	dst, err := os.Create("new.json") //(handler.Filename)
	defer dst.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//var bird Bird
	//bird, err := json.NewDecoder(r.Body)
	birdJson, err := ioutil.ReadFile("./new.json")
	if err != nil {
		fmt.Print(err)
	}
	// Do something with the Person struct...

	//log.Println(t.Test)

	//fmt.Fprintf(w, "Person: %+v", t)
	json.Unmarshal([]byte(birdJson), &bird)
	//fmt.Println("%v", bird.Main.Width)
	//----------------------------------------------------------------------------
	///***************************
	idents1json, _ := ioutil.ReadFile("./test.json")

	//result, _ := json.MarshalIndent(idents1json, "", "")
	json.Unmarshal([]byte(idents1json), &idents1)
	//fmt.Fprint(w, result)
	for i := range idents1 {
		idents = append(idents, Identification{X: idents1[i].X, Y: idents1[i].Y, Width: idents1[i].Width, Height: idents1[i].Height, Time: idents1[i].Time})

	}
	///***************************
	dt := time.Now().Format("01-02-2006 15:04:05")
	for i := range bird.Input {
		if bird.Main.Width > 0 && bird.Main.Height > 0 && bird.Input[i].Width > 0 && bird.Input[i].Height > 0 {
			if isRectangleOverlap([]int{bird.Main.X, bird.Main.Y, bird.Main.Width, bird.Main.Height}, []int{bird.Input[i].X, bird.Input[i].Y, bird.Input[i].Width, bird.Input[i].Height}) { // Condition
				fmt.Println("Overlapped") // Clause

				idents = append(idents, Identification{X: bird.Input[i].X, Y: bird.Input[i].Y, Width: bird.Input[i].Width, Height: bird.Input[i].Height, Time: dt})
				file, _ := json.MarshalIndent(idents, "", " ")
				//_ = ioutil.WriteFile("test.json", file, 0644)
				//fmt.Println(string(file))

				//**************************************************************
				err := ioutil.WriteFile("test.json", file, 0777)
				// handle this error
				if err != nil {
					// print it out
					fmt.Println(err)
				}
				data, err := ioutil.ReadFile("test.json")
				if err != nil {
					fmt.Println(err)
				}

				fmt.Print(string(data))

				f, err := os.OpenFile("test.json", os.O_APPEND|os.O_WRONLY, 0600)
				if err != nil {
					panic(err)
				}
				defer f.Close()

				//if _, err = f.WriteString("new data that wasn't there originally\n"); err != nil {
				//panic(err)
				//}

				data, err = ioutil.ReadFile("test.json")
				if err != nil {
					fmt.Println(err)
				}

				fmt.Print(string(data))
				//**************************************************************
			} else {
				fmt.Println("Not Overlapped")
			}
		}

	}

	//fmt.Println("TEST------------")
	// iterating it
	for _, v := range idents {
		fmt.Println(v)
	}
	fmt.Println()
	fmt.Fprintf(w, "Successfully Uploaded File\n")
	//xx()
}

/*func xx() {
	decoder := json.NewDecoder(r.Body)
	var t Bird
	err1 := decoder.Decode(&t)
	if err1 != nil {
		panic(err1)
	}
}
*/

//-----------------------------------------------------------------------------
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		display(w, "upload", nil)
	case "POST":
		uploadFile(w, r)
	}
}
func getserver(w http.ResponseWriter, req *http.Request) {
	//log.Println(req.URL)
	fmt.Fprint(w, "Thats output\n")
	idents1json, _ := ioutil.ReadFile("./test.json")

	//result, _ := json.MarshalIndent(idents1json, "", "")
	json.Unmarshal([]byte(idents1json), &idents1)
	//fmt.Fprint(w, result)
	for i := range idents1 {
		idents2 = append(idents2, Identification{X: idents1[i].X, Y: idents1[i].Y, Width: idents1[i].Width, Height: idents1[i].Height, Time: idents1[i].Time})

	}
	//for _, v := range idents2 {
	//fmt.Fprint(w, v)
	//fmt.Println("HHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHh")
	//fmt.Println(v)
	//}
	file, _ := json.MarshalIndent(idents2, "", " ")
	//_ = ioutil.WriteFile("test.json", file, 0644)
	fmt.Fprint(w, string(file))

}
func main() {
	// Upload route
	http.HandleFunc("/upload", uploadHandler)

	//--------
	http.HandleFunc("/recieve", getserver)
	//--------

	//Listen on port 8080
	http.ListenAndServe(":8080", nil)
}
