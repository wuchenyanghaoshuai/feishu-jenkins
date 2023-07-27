package feishu

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type JobData struct {
	JobName            []string `json:"value"`
	ChangeType         []string `json:"changeType"`
	GitlabSourceBranch []string `json:"gitlabSourceBranch"`
	Status             string   `json:"status"`
}

func GetProjectInfo(instance string) []map[string]string {

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/approval/v4/instances/%s", instance)
	//url := "https://open.feishu.cn/open-apis/approval/v4/instances/F3CC8321-52C9-4CC4-8A5A-B47169AF6E03"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return nil
	}

	//使用GetTenantAccessToken()函数获取tenantAccessToken
	tenantAccessToken, err := GetTenantAccessToken()
	if err != nil {
		fmt.Println("Error getting tenant access token:", err)
		return nil
	}
	newtenantAccessToken := fmt.Sprintf("Bearer %s", tenantAccessToken)

	req.Header.Set("Authorization", newtenantAccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return nil
	}

	var result map[string]interface{}
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	fmt.Println("===============status==================")
	status := result["data"].(map[string]interface{})["status"]
	//;	fmt.Println("status:", status, result)
	res := fmt.Sprintf(result["data"].(map[string]interface{})["form"].(string))

	var jobData []JobData
	err = json.Unmarshal([]byte(res), &jobData)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil
	}

	var jobNames []string
	for _, jd := range jobData {
		if jd.JobName != nil {
			jobNames = append(jobNames, jd.JobName...)
		}
	}

	lastSecond := jobNames[len(jobNames)-2]
	fmt.Println("changeType:", lastSecond)
	last := jobNames[len(jobNames)-1]
	fmt.Println("gitlabSourceBranch:", last)

	returnresult := []map[string]string{}

	for i := 0; i < len(jobNames)-2; i++ {
		dataMap := map[string]string{
			"jobName":            "",
			"changeType":         "",
			"gitlabSourceBranch": "",
			"status":             "",
		}

		dataMap["jobName"] = jobNames[i]
		dataMap["gitlabSourceBranch"] = last
		dataMap["changeType"] = lastSecond
		dataMap["status"] = status.(string)
		returnresult = append(returnresult, dataMap)
	}

	fmt.Printf("%T  %[1]v \n", returnresult)
	return returnresult

}
