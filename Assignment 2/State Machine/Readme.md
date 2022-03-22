##Readme
https://zetcode.com/golang/flag/

First, copy the path of your .machine file. If you try to run it without the path variable, 
the program will complain. So ```go run .``` will not work. Also, don't try ```go run main.go```.

The correct usage is:\
```go run . -path=YOUR_PATH```\
```go run . -path YOUR_PATH```\
```go run . --path=YOUR_PATH```\
```go run . --path YOUR_PATH```

If you want to create a nice, readable version of your .machine file, you can just add the flag
```-c```.
The newly created file shout be 'serialized-machine.machine'. 

In order to run the tests, just type ```go test```, add teh flag ```-v``` for more details.

The Parser only fails, if you have declared two states with identical state names, if you have declared
two transitions with identical state-action pair, or if the file cannot be found at your given path.

The Validator checks for multiple start states, if it has more than one end state, and if there is 
a run from the start state to any end state (This is required for the program to end).

The Serializer creates a new File and prints the internal structure into this file.

The program ends prematurely if you end up in a sink state (there is no outgoing transition), or if you type
``quit``.


