# Portain 

A small linux CLI tool for diagnosing and managing port conflicts 

## Features
- TCP port availability detection
- Linux process inspection
- Command-line interface

## Initial scope of commands: 
- port check 3000 //determine if the port is available or not 
- port inspect 3000 //identifies which process occupies the port 
- port find 3000 // find the next available port 
- port free 3000 //safely terminate the process to free the port after confirmation

Tech used: Golang 


