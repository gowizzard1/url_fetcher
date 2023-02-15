## Truvity-url-assignment
truvity-url-assignment is a command-line application for measuring the size of the response bodies of a list of URLs.

## Installation
To install the application, clone the repository and run the following command:

```go
go install ./...
```
This will install the awesomeProject command in your $GOBIN directory.

### Usage
```css
awesomeProject [flags] [url1] [url2] ... [urlN]
```
### Flags:

-workers: Number of worker goroutines (default 10)
-timeout: HTTP request timeout (default 5 seconds)
-max-idle: Maximum number of idle connections (default 100)
-idle-timeout: Idle connection timeout (default 30 seconds)
Example usage:

```bash
awesomeProject -workers=20 -timeout=10s https://google.com https://github.com
```
This will output the size of the response body for each URL in the list.

### Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

### License
MIT