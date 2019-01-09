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
	url := fmt.Sprintf("%sbuildWithParameters", job.Url)
	resp, err := jenkins.postUrlResp(url, params, nil)

	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusCreated {
		return 0, fmt.Errorf("Bad status code %d", resp.StatusCode)
	}

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
	var build Build
	u := fmt.Sprintf("%sapi/xml?tree=builds[id,url,number,result,queueId]&xpath=//build[queueId=%d]", job.Url, queueID)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return build, err
	}

	fmt.Println(u)
	resp, err := jenkins.sendRequest(req)
	if err != nil {
		return build, err
	}

	err = jenkins.parseXmlResponse(resp, &build)

	return build, err
}
