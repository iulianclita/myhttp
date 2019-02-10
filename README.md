## Overview

Myhttp is a tool that fetches the content of multiple urls and prints each one of them on a separate line along with the corresponding MD5 hash of the HTML content located at the current url.

### Installation
First, please make sure you have Go installed. For more information follow this [link](https://golang.org/doc/install)

The next steps explain how to build myhttp from source
```sh
cd path/to/location
git clone https://github.com/iulianclita/myhttp.git
cd myhttp
go build
```

### Usage
After succesfully building the tool, use it like this
```sh
$> ./myhttp golang.org google.com
http://www.golang.org d1b40e2a2ba488a054186e4ed0733f9752f66949
http://google.com 9d8ec921bdd275fb2a605176582e08758eb60641

Use the -parallel flag to control the maximum number of parallel HTTP requests

$> ./myhttp -parallel 3 golang.org google.com
http://www.golang.org d1b40e2a2ba488a054186e4ed0733f9752f66949
http://google.com 9d8ec921bdd275fb2a605176582e08758eb60641
```

### Testing
To run all unit tests do the following:
```sh
cd /path/to/myhttp/source
go test -v ./...
```

To run integration tests add the integration tag:
```sh
cd /path/to/myhttp/source
go test -v -tags integration ./...
```