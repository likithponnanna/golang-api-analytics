package main

import (
	"fmt"
	"github.com/pborman/getopt/v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"sort"
	"strconv"
	"time"
)

/**
Struct contains the required values for each request entity.
 */
type responseEntity struct {
	time         time.Duration
	responseCode int
	responseSize int
}

/**
Function to get server Response
Returns a server response and start time.
 */
func getServerResponse(url string) (*http.Response, time.Time) {
	spaceClient := http.Client{
		Timeout: time.Second * 2,
	}

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("User-Agent", "spacecount-tutorial")

	start := time.Now()
	response, getErr := spaceClient.Do(request)
	if getErr != nil {
		log.Fatal(getErr)
	}

	return response, start

}


/**
Function to get total request time and response size.
Returns request time and response size.
 */
func performTimedUtils(response *http.Response, startTime time.Time) (time.Duration, bool, int) {
	requestTime := time.Since(startTime)

	if response.Body != nil {
		requestTime := time.Since(startTime)

		dump, err := httputil.DumpResponse(response, true)
		if err != nil {
			log.Fatal(err)
		}

		defer response.Body.Close()

		return requestTime, true, len(dump)
	}

	return requestTime, false, 0
}


/**
Function to print response from single url.
 */
func printBody(response *http.Response) {

	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	bodyString := string(body)
	fmt.Println("Response Body: ", bodyString)

}


/**
Function to get command line arguments.
Returns url and number of request.
 */
func runNumberedRequests() (string, int) {

	urlStr := getopt.StringLong("url", 'u',
		"https://www.alphavantage.co/query?function=TIME_SERIES_DAILY_ADJUSTED&symbol=IBM&apikey=demo",
		" The url to fetch! ")
	profileStr := getopt.StringLong("profile", 'p', "0",
		" Number of fetches required. \n Default url"+
		"is https://www.alphavantage.co/query?function=TIME_SERIES_DAILY_ADJUSTED&symbol=IBM&apikey=demo " +
		"if not explicitly set with --url flag")
	optHelp := getopt.BoolLong("help", 0, "Help")
	getopt.Parse()

	if *optHelp {
		getopt.Usage()
		os.Exit(0)
	}

	requestsNo, err := strconv.ParseInt(*profileStr, 10, 64)
	if err != nil {
		log.Fatal("Invalid request number : Err : ", err)
	}

	return *urlStr, int(requestsNo)

}

/**
Function to return minimum and maximum response time .
Returns minimum and maximum response time.
*/
func minMaxTime(sliceProcess []time.Duration) (time.Duration, time.Duration) {
	var max = sliceProcess[0]
	var min = sliceProcess[0]
	for _, value := range sliceProcess {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}

/**
Function to return minimum and maximum buffer size.
Returns minimum and maximum buffer responses.
 */
func minMaxBufferSize(sliceProcess []int) (int, int) {
	var max = sliceProcess[0]
	var min = sliceProcess[0]
	for _, value := range sliceProcess {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}

/**
Function to calculate median time.
Returns median time.
 */
func medianTime(sliceProcess []time.Duration) time.Duration {
	if len(sliceProcess) == 1 {
		return sliceProcess[0]
	}
	floatTimeArray := []float64{}
	for _, timeSingle := range sliceProcess {
		floatTimeArray = append(floatTimeArray, float64(timeSingle))
	}

	sort.Float64s(floatTimeArray)
	midEle := len(floatTimeArray) / 2.0

	if midEle%2 == 0 {
		return time.Duration((floatTimeArray[midEle-1] + floatTimeArray[midEle]) / 2)
	}

	return time.Duration(floatTimeArray[midEle])

}

/**
Function to calculate mean time.
Returns the mean time.
 */
func meanTime(sliceProcess []time.Duration) time.Duration {
	total := 0.0
	for _, timeSingle := range sliceProcess {
		total += float64(timeSingle)
	}
	return time.Duration(total / float64(len(sliceProcess)))
}

/**
Function to calculate success percentage of requests.
Returns the success percentage.
 */
func percentSuccess(sliceProcess []int) float64 {
	totalUnsuccessfull := 0
	totalReqs := len(sliceProcess)
	for _, code := range sliceProcess {
		if code == 200 {
			totalUnsuccessfull++
		}
	}

	return float64(totalUnsuccessfull/totalReqs) * 100

}

/**
Function to remove duplicate error codes.
Returns an array without duplicate error codes.
 */
func getUniqueErrorCodes(sliceProcess []int) []int {
	keys := make(map[int]bool)
	var uniqueErrorCodes []int
	for _, code := range sliceProcess {
		if _, value := keys[code]; !value {
			keys[code] = true
			uniqueErrorCodes = append(uniqueErrorCodes, code)
		}
	}

	return uniqueErrorCodes

}


/**
Function to generate a list of error codes
Returns error code list with duplicates.
 */
func getErrorCodes(sliceProcess []int) []int {
	var errorCodesSlice []int
	for _, code := range sliceProcess {
		if code != 200 {
			errorCodesSlice = append(errorCodesSlice, code)
		}
	}

	return errorCodesSlice

}


/**
Function to format final answer
 */
func formatPrintProfiler(numOfReq int, fastestTime time.Duration, slowestTime time.Duration,
	meanTm time.Duration, medianTm time.Duration, percenRequestSuccess float64,
	uniqueErrorCodeSlice []int, minResponseSize int, maxResponseSize int) {

	fmt.Println("\n Number Of Requests: ", numOfReq)
	fmt.Println(" Fastest Time: ", fastestTime, "\n Slowest Time: ", slowestTime,
		"\n Mean Time: ", meanTm, "\n Median Time: ", medianTm, "\n Percent Success Requests: ",
		percenRequestSuccess)

	if len(uniqueErrorCodeSlice) == 0 {
		fmt.Println("Unique Error Code List is empty i.e only 200 response code encountered ",
			uniqueErrorCodeSlice)
	} else {
		fmt.Println("Unique Error Code List: ", uniqueErrorCodeSlice)
	}

	fmt.Println("Smallest Response Size in bytes: ", minResponseSize,
		"\n Largest Response Size in bytes: ", maxResponseSize)

}


/**
Function to get appropriate calculations for the final response.
 */
func processProfilerData(responseArray []responseEntity) {

	numOfReq := len(responseArray)
	var timingArray []time.Duration
	var responseCodes []int
	var responseSizes []int

	for _, individualResponse := range responseArray {
		timingArray = append(timingArray, individualResponse.time)
		responseCodes = append(responseCodes, individualResponse.responseCode)
		responseSizes = append(responseSizes, individualResponse.responseSize)
	}

	slowestTime, fastestTime := minMaxTime(timingArray)
	meanTm := meanTime(timingArray)
	medianTm := medianTime(timingArray)
	percenRequestSuccess := percentSuccess(responseCodes)
	errorCodeSlice := getErrorCodes(responseCodes)
	uniqueErrorCodeSlice := getUniqueErrorCodes(errorCodeSlice)
	minResponseSize, maxResponseSize := minMaxBufferSize(responseSizes)

	formatPrintProfiler(numOfReq, fastestTime, slowestTime, meanTm,
		medianTm, percenRequestSuccess, uniqueErrorCodeSlice,
		minResponseSize, maxResponseSize)

}

/**
Function to generate a response array containing each request element.
 */
func responseProfiler(numOfRequests int, url string) {
	var responseArray []responseEntity

	for i := 0; i < numOfRequests; i++ {

		response, startTime := getServerResponse(url)
		requestTime, flagUpdate, responseSz := performTimedUtils(response, startTime)
		if flagUpdate == true {
			entity := responseEntity{time: requestTime, responseCode: response.StatusCode, responseSize: responseSz}
			responseArray = append(responseArray, entity)
		}
	}

	processProfilerData(responseArray)

}

/**
Main Code Controller
 */
func main() {

	url, reqNo := runNumberedRequests()

	if reqNo == 0 {
		fmt.Println("URL:  ", url, " \n Number Of Request: ", reqNo+1)
		response, _ := getServerResponse(url)
		printBody(response)

	} else {
		fmt.Println("URL:  ", url)
		responseProfiler(reqNo, url)
	}

}
