"# GoTSQL" 

You can connect to Microsoft SQL servers via this web app, and submit dyanmic SQL queries and get results with it.
You can also save any results to CSV file contained within 'outputsaves' folder which is generated within the folder the application is kept.

The purpose of this tool is to show that you can use reflection to grab dynamic information without having to pre-write the schema in go. 

Bugs:
1) It can't parse date.time types so if you have a date/time column - cast it to varchar and it will parse.
If you can fix this go ahead! problem is in backendsql.go - row 52 and row 72.

2) In backendsql.go row 41 I use statement 'if reflection.Field(i).String() == "<[]driver.Value Value>"'
If somebody could refactor this so that it doesn't use a string to identify the correct struct field then that would be cool.
