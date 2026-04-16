module server

go 1.25.5

require (
    github.com/gorilla/mux v1.8.1
    monitoring/api v0.0.0
)   

replace monitoring/api => ../api
