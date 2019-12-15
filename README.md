# airports
A simple Go program to get the airports in a given area. 
It needs the user coordinates and the distance in degrees to search. 

After installing the Go runime just type:
$ go run airportsAround.go and follow the instructions. 

Alternatively pass the follow flags when running the command:
 $ go run airportsAround.go -lon XX -lat YY -d DD

 Where XX is your longitude, YY is your latitude and DD is the search radius (in degrees)


oldAirportdb.go in the old folder is an alternative algorithm that requires max and min coordinates pairs to search for airports. 


# Changelog 
Added support for invoking with flags 

# To do 

Opptimize the flags code, maybe with an new function. 