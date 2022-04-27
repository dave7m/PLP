##Readme

For the third Assignment, the action names of the transitions from auto-forwarding states should in my
opinion match with the action names of the states that called them. E.g In rolling.machine, the transition 
"Rolling" has the same action name (slide) as the auto-forwarding state that called it (called by rolling,
which itself was called by slide). However, with this logic, the doggo.machine would not work, thus I 
decided to ignore any action names in transitions from auto-forwarding states. You'll find in the validator
some commented code that shows I would have tested for this.

\
Also, the parser ignores invalid states, and does not try to recover, so @!abc{hi} is not valid and ignored, @!abc2{hi} 
is valid, @abc{hi} is not valid and ignored.

\
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


