package action

import (
	//	"encoding/json"
	"fmt"
	//	"github.com/golang/glog"
	//	vmtmeta "github.com/pamelasanchezvi/communicator/metadata"
	//	"github.com/pamelasanchezvi/communicator/util"
	//	"github.com/vmturbo/vmturbo-go-sdk/sdk"
	"bytes"
	"io/ioutil"
	"net/http"
)

/*
type recordingTransport struct {
	req *http.Request
}

func (t *recordingTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	t.req = req
	return nil, errors.New("dummy impl")
}
*/

type migration struct {
	destination_node_id string
	task_ids            []string
}

func RequestMesosAction(mesosClient *MesosClient) (string, error) {

	baseUrl := "http://" + mesosClient.MesosMasterIP + ":" + mesosClient.MesosMasterPort + "/" + mesosClient.Action + "?"
	//fullUrl := baseUrl + "destination_node_id=32f951d7-52f8-4842-ae1f-eb8d7ec6ac94-S0&task_ids=basic-0.6432abd7-179f-11e6-9521-52540006b4aa"
	fmt.Println(" --> The full Url is ", baseUrl)
	/*
		str := "\"" + mesosClient.TaskId + "\""
		taskid := []string{str}                         //"basic-0.b34401b2-1844-11e6-bafb-52540006b4aa"}
		node := "\"" + mesosClient.DestinationId + "\"" //"32f951d7-52f8-4842-ae1f-eb8d7ec6ac94-S0"

		m := migration{node, taskid}
		fmt.Printf("payload is : %+v \n", m)
		b, err := json.Marshal(m)
		var jsonStr = []byte(b)
	*/
	var jsonStr = []byte(`{"destination_node_id":"` + mesosClient.DestinationId + `", "task_ids": ["` + mesosClient.TaskId + `"]}`)
	req, err := http.NewRequest("POST", baseUrl, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	/*	form := url.Values{}
		form.Set("destination_node_id", "32f951d7-52f8-4842-ae1f-eb8d7ec6ac94-S0")
		form.Set("task_ids", "basic-0.b34401b2-1844-11e6-bafb-52540006b4aa")

		fmt.Println(" ")
		resp, err := http.PostForm(baseUrl,form)
	*/
	if err != nil {
		fmt.Printf(" --> error %s \n", err)
	}
	defer resp.Body.Close()
	fmt.Printf("----> request is : %+v\n", req)
	fmt.Printf("response status %s Headers: %s \n", resp.Header, resp.Status)

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	//	fmt.Println("Get Succeed: %v", respContent)
	//	defer resp.Body.Close()

	return string(body), nil
}

func RequestPendingTasks(mesosClient *MesosClient) []*PendingTask {
	// 10.10.174.96:5555/GetPendingTasks
	baseUrl := "http://" + mesosClient.MesosMasterIP + ":" + mesosClient.MesosMasterPort + "/" + "GetPendingTasks"
	//fullUrl := baseUrl + "destination_node_id=32f951d7-52f8-4842-ae1f-eb8d7ec6ac94-S0&task_ids=basic-0.6432abd7-179f-11e6-9521-52540006b4aa"
	fmt.Println(" --> The full Url is ", baseUrl)
	req, err := http.NewRequest("GET", baseUrl, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf(" --> error %s \n", err)
	}
	var pendingTasks = make([]*PendingTask)
	byteContent := []byte(resp)
	err = json.Unmarshal(byteContent, &pendingTasks)
	if err != nil {
		fmt.Printf("JSON error in getPendingTasks %s", err)
	}
	var pendingTaskArray []*PendingTask
	pendingTaskArray = pendingTasks
	defer resp.Body.Close()
	return pendingTaskArray
}
