"# GoTSQL" 

Can be placed on the web so you can connect to MS SQL servers online (as long as you can form the network connection!)

The point of this tool is to show that you can use reflection to grab all the information you need and convert it into CSV.

Bugs:
1) It can't parse date.time types so if you have a date/time column - cast it to varchar and it will parse.
if you want to fix it go ahead.

2) In backendsql.go row 41 I use statement 'if reflection.Field(i).String() == "<[]driver.Value Value>"'
If somebody could refactor this so that it doesn't use a string to identify the correct struct field then that would be cool.
