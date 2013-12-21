# A Very Simple Wiki

This repo provides a reference example of the [go web application tutorial](http://golang.org/doc/articles/wiki/)

### Using this repo
Clone the repo and take a look at the log.  
The log contains the individual steps and construction of the wiki concept

### Building 
use `go run wiki.go` to start the server locally

### Using the Wiki
access [http://localhost:8080/view/FrontPage](http://localhost:8080/view/FrontPage)  
edit pages [http://localhost:8080/edit/FrontPage](http://localhost:8080/edit/FrontPage)  

### Configuration
Listens on localhost:8080  
Wiki content is stored in data/  
HTML templates are stored in templates/