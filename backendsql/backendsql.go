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

var db *sql.DB

func grabrows(r interface{}) string {
	//Scanner := bufio.NewScanner(os.Stdin)
	text := ""
	reflection := reflect.ValueOf(r)
	for i := 0; i < reflection.NumField(); i++ {
		//fmt.Println(reflection.Field(i).Kind())
		if reflection.Field(i).String() == "<[]driver.Value Value>" {
			x := reflection.Field(i)
			for ii := 0; ii < x.Len(); ii++ {
				if ii != x.Len()-1 {
					switch x.Index(ii).IsNil() {
					case false:
						switch x.Index(ii).Elem().Type().String() {
						case "string":
							text += x.Index(ii).Elem().String() + ","
						case "bool":
							text += strconv.FormatBool(x.Index(ii).Elem().Bool()) + ","
						//case "time.Time":
						//	text += "DATETIME_type" //Haven't figured out how to convert this type into a string yet.
						case "int64":
							text += strconv.FormatInt(x.Index(ii).Elem().Int(), 10) + ","
						case "float64":
							text += strconv.FormatFloat(x.Index(ii).Elem().Float(), 'f', 6, 64) + ","
						default:
							text += "OTHER_type," //Pick up any other types here.
						}
					case true:
						text += ","
					}
				} else {
					switch x.Index(ii).IsNil() {
					case false:
						switch x.Index(ii).Elem().Type().String() {
						case "string":
							text += x.Index(ii).Elem().String()
						case "bool":
							text += strconv.FormatBool(x.Index(ii).Elem().Bool())
						//case "time.Time":
						//	text += "DATETIME_type" //Haven't figured out how to convert this type into a string yet.
						case "int64":
							text += strconv.FormatInt(x.Index(ii).Elem().Int(), 10)
						case "float64":
							text += strconv.FormatFloat(x.Index(ii).Elem().Float(), 'f', 6, 64)
						default:
							text += "OTHER_type" //Pick up other types here.
						}
					}
				}
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
	cols, errr := rows.Columns()
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
