package main

import (
	"fmt"
	"gosqlwebapp/backendsql"
	"gosqlwebapp/front"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
)

//PageVariables - GUI variables that change on webpages.
type PageVariables struct {
	Output        string
	SqlQuery      string
	Outputresults string
}

var store string

//EnsureDir - Ensures that a DIR exists
func ensureDir(dirName string) error {
	err := os.MkdirAll(dirName, os.ModeDir)
	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}

func writeToFile(filename string, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return file.Sync()
}

func mainpage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./ui/mainpage.gtpl")
	if err != nil { // if there is an error
		log.Print("template executing error: ", err) //log it
	}
	HomePageVars := PageVariables{}
	err2 := t.Execute(w, HomePageVars)
	if err2 != nil { // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

func start(w http.ResponseWriter, r *http.Request) {
	store = ""
	t, err := template.ParseFiles("./ui/mainpage.gtpl")
	if err != nil { // if there is an error
		log.Print("template executing error: ", err) //log it
	}
	HomePageVars := PageVariables{}
	r.ParseForm()
	sqlLogin := ""
	sqlQuery := ""
	if len(r.Form["sqlLogin"]) != 0 && r.Form["sqlLogin"][0] != "" {
		sqlLogin += r.Form["sqlLogin"][0]
	}
	if len(r.Form["sqlQuery"]) != 0 && r.Form["sqlQuery"][0] != "" { //tablequery
		sqlQuery += r.Form["sqlQuery"][0]
		HomePageVars.SqlQuery = r.Form["sqlQuery"][0]
		count, results, err := backendsql.SQLquery(sqlLogin, sqlQuery)
		if err != nil {
			HomePageVars.Output = "Error: \n" + err.Error()
			err22 := t.Execute(w, HomePageVars)
			if err22 != nil { // if there is an error
				log.Print("template executing error: ", err) //log it
			}
			return
		}
		fmt.Printf("Sucessfully pulled %d rows.\n", count)
		text := ""
		for index := range results {
			text += results[index] + "\n"
		}
		store = text
		HomePageVars.Output = text
	}
	if len(r.Form["tablequery"]) != 0 && r.Form["sqlQuery"][0] == "" {
		sqlQuery += "SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = N'" + r.Form["tablequery"][0] + "'"
		HomePageVars.SqlQuery = sqlQuery
		count, results, err := backendsql.SQLquery(sqlLogin, sqlQuery)
		if err != nil {
			HomePageVars.Output = "Error: \n" + err.Error()
			err22 := t.Execute(w, HomePageVars)
			if err22 != nil { // if there is an error
				log.Print("template executing error: ", err) //log it
			}
			return
		}
		fmt.Printf("Sucessfully pulled %d rows.\n", count)
		text := ""
		for index := range results {
			if index != 0 {
				text += results[index] + "\n"
			}
		}
		store = text
		HomePageVars.Output = text
	}
	err2 := t.Execute(w, HomePageVars)
	if err2 != nil { // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

func save(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./ui/mainpage.gtpl")
	if err != nil { // if there is an error
		log.Print("template executing error: ", err) //log it
	}
	HomePageVars := PageVariables{}
	filename := ""
	r.ParseForm()
	if len(r.Form["outputname"]) != 0 && r.Form["outputname"][0] != "" {
		filename += r.Form["outputname"][0] + ".csv"
		fmt.Println(filename + " has been saved!")
		writeerror := writeToFile("./outputsaves/"+filename, store)
		if writeerror != nil {
			fmt.Println(writeerror)
		}
	} else {
		fmt.Println("Filename was empty")
		HomePageVars.Outputresults = "...Filename was empty!"
	}
	err2 := t.Execute(w, HomePageVars)
	if err2 != nil { // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

func main() {
	ensureDir("outputsaves")
	http.HandleFunc("/", mainpage)
	http.HandleFunc("/ui/start", start)
	http.HandleFunc("/ui/save", save)
	front.BuildPage("mainpage",
		`<title>GO QL - Front-end SQL</title>
</head>

<body>
    <style>
        input {
            display: inline-block;
            float: left;
            margin-right: 20px;
            background-color: #00FFFF
        }
    </style>
    <a href="https://imgbb.com/"><img src="https://i.ibb.co/dDg4kWH/images.png" alt="SQL-LOGO" border="0"></a>
    <p><b>GO QL (Go) - Front-end SQL interface</b></p>
    <p>GO QL</p>
    <p>A simple front-end interface for TSQL.</p>
    <form action="/ui/start" method="POST">
        <p> <input type="text" name="sqlLogin" placeholder="Server=localhost;Database=master;Trusted_Connection=True;" size="50"><label> <- Type in the SQL login details here (i.e Server=localhost;Database=master;Trusted_Connection=True;)</label></p>
		<p> <input type="text" name="tablequery" placeholder="table name" size="50"><label> <- Type in table name here to grab column information for queries</label></p>
		<textarea name="sqlQuery" cols="90" rows="20" placeholder"Type your SQL query here">{{.SqlQuery}}</textarea>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<textarea name="output" cols="90" rows="20">{{.Output}}</textarea><br />
        <br>
		<button type="submit" style="background-color: #00FFFF;" value="submitquery">Submit query</button>
		</form>

		 <form action="/ui/save" method="POST">
		 <input type="text" name="outputname" placeholder="filename" size="50"><label><- Type name of the CSV file you want to save</label><p>
		<button type="submit" style="background-color: #00FFFF;" value="saveoutput">Save output to CSV</button>
		{{.Outputresults}}
		</form>

	 <br/>
    <br>
</body>`)
	serverPort := front.LaunchServer()
	err2 := http.ListenAndServe(serverPort, nil) // setting listening port
	if err2 != nil {
		fmt.Println(err2)
		log.Fatal("ListenAndServe: ", err2)
	}
}
