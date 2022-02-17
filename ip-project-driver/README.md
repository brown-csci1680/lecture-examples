# Example IP node usage

This is an example framework for structuring command-line parsing in the IP driver.  

This example does not implement any of the network-related commands, it only
demonstrates how you can think about parsing command-line input and sending it
to handlers for individual commands.  In addition, it uses `libreadline`, a C
library that provides nice command-line history, eg. pressing
"up" will show your last command.  


To compile the example, run `make`.  
