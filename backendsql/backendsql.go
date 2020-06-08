package backendsql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
)

//Pointer to DB struct
var db *sql.DB

//SQLquery - Submit a SQL username/password/access details (access) and query (q) and receive three outputs: number of rows, a 1D slice of query result, error.
func SQLquery(access, q string) (int, []string, error) {
	//connString := `Server=localhost;Database=master;Trusted_Connection=True;`
	var err error
	db, err = sql.Open("sqlserver", access)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Connected!\n")
	count, result, err := read(q)
	return count, result, err
}

//read takes in a TSQL query and passes it to a TSQL server using the go-mssqldb driver.
func read(query string) (int, []string, error) {
	result := []string{}
	ctx := context.Background()
	err := db.PingContext(ctx)
	if err != nil {
		return -1, result, err
	}
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return -1, result, err
	}
	defer rows.Close()
	var count int
	cols, err := rows.Columns()
	if err != nil {
		return -1, result, err
	}
	columnHeaders := strings.Join(cols, ",")
	firstrow := true
	for rows.Next() {
		if firstrow == true {
			result = append(result, columnHeaders)
			result = append(result, grabrows(*rows))
			firstrow = false
		} else {
			result = append(result, grabrows(*rows))
		}
		count++
	}
	return count, result, nil
}

//grabrows takes the type *sql.Rows interface and uses reflection to collect the items.
//The field of this interface we need is a '<[]driver.Value Value>' and contains the row values so this is where we want to collect the rows.
//For further info check the sql docs.
func grabrows(r interface{}) string {
	text := ""
	reflection := reflect.ValueOf(r)
	for i := 0; i < reflection.NumField(); i++ {
		if reflection.Field(i).String() == "<[]driver.Value Value>" {
			x := reflection.Field(i)
			for ii := 0; ii < x.Len(); ii++ {
				if ii != x.Len()-1 {
					text += typeswitch(x.Index(ii), false)
				} else {
					text += typeswitch(x.Index(ii), true)
				}
			}
		}
	}
	return text
}

//typeswitch returns a string that you append to the 'text' variable in the grabrows function.
//It switches on reflect.Value type casted as a string and converts the item into a string to be returned.
//switch x.Elem().Type() didn't seem to work unfortunately.
func typeswitch(x reflect.Value, lastitem bool) string {
	var text = ""
	switch lastitem {
	case false:
		switch x.IsNil() {
		case false:
			switch x.Elem().Type().String() {
			case "string":
				text += x.Elem().String() + ","
			case "bool":
				text += strconv.FormatBool(x.Elem().Bool()) + ","
			case "time.Time":
				text += "(Change_to_VARCHAR)," //Reflecting to time.Time value doesn't seem to work as far as I can see.
				//timestring, _ := reflect.ValueOf(x.Index(ii).Elem()).Interface().(time.Time)
				//fmt.Println(timestring)
			case "int64":
				text += strconv.FormatInt(x.Elem().Int(), 10) + ","
			case "float64":
				text += strconv.FormatFloat(x.Elem().Float(), 'f', 6, 64) + ","
			default:
				text += x.Elem().Type().String() + "_type_(Change_to_VARCHAR)," //Pick up any other types here.
			}
		case true:
			text += ","
		}
	case true:
		switch x.IsNil() {
		case false:
			switch x.Elem().Type().String() {
			case "string":
				text += x.Elem().String()
			case "bool":
				text += strconv.FormatBool(x.Elem().Bool())
			case "time.Time":
				text += "(Change_to_VARCHAR)" //Reflecting to time.Time value doesn't seem to work as far as I can see.
				//timestring, _ := reflect.ValueOf(x.Index(ii).Elem()).Interface().(time.Time)
				//fmt.Println(timestring)
			case "int64":
				text += strconv.FormatInt(x.Elem().Int(), 10)
			case "float64":
				text += strconv.FormatFloat(x.Elem().Float(), 'f', 6, 64)
			default:
				text += x.Elem().Type().String() + "_type_(Change_to_VARCHAR)" //Pick up any other types here.
			}
		}
	}
	return text
}
