package gojenkins

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Create a new build for this job.
// Params can be nil.
func (jenkins *Jenkins) BuildWithQueueID(job Job, params url.Values) (int, error) {

	//var url string
	//if hasParams(job) {
	url := fmt.Sprintf("%sbuildWithParameters", job.Url)
	//} else {
	//url = fmt.Sprintf("%sbuild", job.Url)
	//}
	fmt.Println(url)
	resp, err := jenkins.postUrlResp(url, params, nil)

	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusCreated {
		return 0, fmt.Errorf("Bad status code %d", resp.StatusCode)
	}
	fmt.Println("Headers: ", resp.Header)

	location := resp.Header["Location"]
	if len(location) == 0 {
		return 0, fmt.Errorf("Could not parse location header: none set")
	}

	split := strings.Split(location[0], "/")
	if len(split) < 2 {
		return 0, fmt.Errorf("Could not parse location header: path not understood")
	}

	queueID, err := strconv.Atoi(split[len(split)-2])
	if err != nil {
		return 0, fmt.Errorf("Could not parse location header: invalid integer")
	}
	return queueID, nil
}

func (jenkins *Jenkins) GetBuildFromJobAndQueueID(job Job, queueID int) (Build, error) {
	u := fmt.Sprintf("%s/api/xml?tree=builds[id,url,number,result,queueId]&xpath=//build[queueId=%d]", job.Url, queueID)
	var build Build
	err := jenkins.getXml(u, url.Values{}, &build)
	return build, err
}
