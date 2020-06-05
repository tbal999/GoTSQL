package backendsql

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/howeyc/gopass"
)

var db *sql.DB

func appendSlice(x []string, y *[][]string) {
	slice := *y
	slice = append(slice, x)
	*y = slice
}

func printSlice(y [][]string) {
	for i := range y {
		fmt.Println(y[i])
	}
}

func columnheaders(r interface{}) string {
	text := ""
	reflection := reflect.ValueOf(r)
	for i := 0; i < reflection.NumField(); i++ {
		if reflection.Field(i).String() == "<[]driver.Value Value>" {
			x := reflection.Field(i)
			for ii := 0; ii < x.Len(); ii++ {
				if ii != x.Len()-1 {
					text += "COLUMN_" + strconv.Itoa(ii+1) + ","
				} else {
					text += "COLUMN_" + strconv.Itoa(ii+1)
				}
			}
		}
	}
	return text
}

func grabrows(r interface{}) string {
	Scanner := bufio.NewScanner(os.Stdin)
	text := ""
	reflection := reflect.ValueOf(r)
	for i := 0; i < reflection.NumField(); i++ {
		if reflection.Field(i).String() == "<[]driver.Value Value>" {
			x := reflection.Field(i)
			for ii := 0; ii < x.Len(); ii++ {
				switch x.Index(ii).Elem().Type().String() {
				case "string":
					if ii != x.Len()-1 {
						if x.Index(ii).Elem().String() == "<invalid Value>" {
							text += "" + ","
						} else {
							text += x.Index(ii).Elem().String() + ","
						}
					} else {
						if x.Index(ii).Elem().String() == "<invalid Value>" {
							text += "" + ","
						} else {
							text += x.Index(ii).Elem().String()
						}
					}
				case "time.Time":
					if ii != x.Len()-1 {
						if x.Index(ii).Elem().String() == "<invalid Value>" {
							text += "" + ","
						} else {
							//text += x.Index(ii).Elem().Format("2008-02-28")
						}
					} else {
						if x.Index(ii).Elem().String() == "<invalid Value>" {
							text += "" + ","
						} else {
							//text += x.Index(ii).Elem().Format("2008-02-28")
						}
					}
				case "int64":
					if ii != x.Len()-1 {
						if x.Index(ii).Elem().String() == "<invalid Value>" {
							text += "" + ","
						} else {
							text += strconv.FormatInt(x.Index(ii).Elem().Int(), 10) + ","
						}
					} else {
						if x.Index(ii).Elem().String() == "<invalid Value>" {
							text += "" + ","
						} else {
							text += strconv.FormatInt(x.Index(ii).Elem().Int(), 10)
						}
					}
				case "float64":
					if ii != x.Len()-1 {
						if x.Index(ii).Elem().String() == "<invalid Value>" {
							text += "" + ","
						} else {
							text += strconv.FormatFloat(x.Index(ii).Elem().Float(), 'f', 6, 64) + ","
						}
					} else {
						if x.Index(ii).Elem().String() == "<invalid Value>" {
							text += "" + ","
						} else {
							text += strconv.FormatFloat(x.Index(ii).Elem().Float(), 'f', 6, 64)
						}
					}
				}
				Scanner.Scan()

			}
		}
	}
	return text
}

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
	firstrow := true
	for rows.Next() {
		if firstrow == true {
			result = append(result, columnheaders(*rows))
			result = append(result, grabrows(*rows))
			firstrow = false
		} else {
			result = append(result, grabrows(*rows))
		}
		count++
	}
	return count, result, nil
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

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
