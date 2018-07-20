# NMWRAP Reports

This app was created to fulfill the report creation functionality that is required for the NMWRAP deliverables. 

## Getting Started

```
git clone https://github.com/edac/nmwrapreports.git
```

See deployment for notes on how to deploy the project on a live system.

### Prerequisites

To compile this code, you will need the golang compiler and libraries.


```
https://golang.org/dl/
```

### Installing

Set your GOPATH environment variable to the location you cloned to

```
export GOPATH="/path/to/nmwrapreports"

```

Set your PATH environment variable to include the go binaries if this was not done as part of the golang install.

```
export PATH="$PATH:/usr/local/go/bin/"

```

Download all requirements from github.

```
go get nmwrapreports
```

Build the app

```
go install nmwrapreports
```

Copy the included nmwrapreports.conf file to etc

```
cp /path/to/nmwrapreports/nmwrapreports.conf /etc/nmwrapreports/nmwrapreports.conf 
```



## Running the service


```
/path/to/nmwrapreports/bin/nmwrapreports
```

## License

This project is licensed under the MIT License


