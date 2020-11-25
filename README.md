# golang-website-analytics

# Manual Build Steps:



- Run go get command to install missing getopt package.
```console
foo@bar:~$ go get github.com/pborman/getopt/v2
```

- Run build command to build executable binaries based on the system.
```console
 foo@bar:~$ go build
```

### --url flag
- Run the executable binary with --url flag to get json response. For linux and darwin system is as follows:

```console
foo@bar:~$  ./goAPIAnalyzer --url 
https://www.alphavantage.co/query?function=TIME_SERIES_DAILY_ADJUSTED&symbol=IBM&apikey=demo
```
OR for windows
```console
foo@bar:~$  goAPIAnalyzer.exe --url 
https://www.alphavantage.co/query?function=TIME_SERIES_DAILY_ADJUSTED&symbol=IBM&apikey=demo
```
### --profile flag
- Run with --profile numberOfRequest for example --profile 4 for four requests. Default url is 
 for the alphavantage API
```console
foo@bar:~$  ./goAPIAnalyzer --profile 4
```
OR 

```console
foo@bar:~$  goAPIAnalyzer.exe --profile 4
```


### --help flag
- Use --help flag for the systems appropriately to display options.

```console
foo@bar:~$  ./goAPIAnalyzer --help
```



#
##



 



