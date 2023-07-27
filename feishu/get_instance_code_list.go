package feishu

import (
	"encoding/json"
	"fieshu-jenkins/jenkins"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type InstanceListResponse struct {
	Code int `json:"code"`
	Data struct {
		HasMore          bool     `json:"has_more"`
		InstanceCodeList []string `json:"instance_code_list"`
		PageToken        string   `json:"page_token"`
	} `json:"data"`
	Msg string `json:"msg"`
}

func GetInstanceCodeList() {
	now := time.Now()
	// 构造当天的零点时间和23:59:59时间
	year, month, day := now.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
	endOfDay := time.Date(year, month, day, 23, 59, 59, 0, now.Location())

	// 转换为毫秒时间戳
	startTimestamp := startOfDay.UnixNano() / int64(time.Millisecond)
	endTimestamp := endOfDay.UnixNano() / int64(time.Millisecond)
	url := fmt.Sprintf("https://open.feishu.cn/open-apis/approval/v4/instances?approval_code=90E68816-635D-4AF5-AFB6-7C8EC7F617E4&end_time=%d&page_size=100&start_time=%d", endTimestamp, startTimestamp)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}

	//使用GetTenantAccessToken()函数获取tenantAccessToken
	tenantAccessToken, err := GetTenantAccessToken()
	if err != nil {
		fmt.Println("Error getting tenant access token:", err)
		return
	}
	newtenantAccessToken := fmt.Sprintf("Bearer %s", tenantAccessToken)

	req.Header.Set("Authorization", newtenantAccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return
	}

	var response InstanceListResponse
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		fmt.Println("解析响应失败:", err)
		return
	}
	for _, data := range response.Data.InstanceCodeList {
		fmt.Println("====================================================")
		//		fmt.Println(data)
		_, err := CreateRedisInstance("get", data)
		if err != nil {
			fmt.Println("获取key失败, redis里没有存这个key:", err)

			projectInfo := GetProjectInfo(data)

			for _, item := range projectInfo {
				fmt.Println("====================")
				fmt.Println("下面打印的是审批获取到的所有信息")
				fmt.Printf("JobName: %s, ChangeType: %s, GitlabSourceBranch: %s , Status: %s ,审批单实例ID: %s \n", item["jobName"], item["changeType"], item["gitlabSourceBranch"], item["status"], data)
				fmt.Println("====================")
				if item["status"] == "APPROVED" {
					fmt.Println("审批单已经通过了下面开始执行发版本程序")
					fmt.Printf("项目名称: %s, 发版环境: %s, 选择分之: %s , 审批单状态: %s \n", item["jobName"], item["changeType"], item["gitlabSourceBranch"], item["status"])
					_, err = CreateRedisInstance("set", data, "123")
					jenkins.BuildHandler(item["jobName"], item["changeType"], item["gitlabSourceBranch"])

				} else if item["status"] == "PENDING" {
					fmt.Println("单子正在审批中，请耐心等待")
				} else if item["status"] == "REJECTED" {
					fmt.Println("发版被拒绝,请找管理员确认原因")
					_, err = CreateRedisInstance("set", data, "123")
				} else {
					return
				}
			}

			//if projectInfo[0].Status == "APPROVED" {
			//	fmt.Printf("JobName: %s, ChangeType: %s, GitlabSourceBranch: %s , Status: %s ,Instance: %s \n", projectInfo[0].JobName, projectInfo[0].ChangeType, projectInfo[0].GitlabSourceBranch, projectInfo[0].Status, data)
			//	fmt.Println("这个instance已经被审批通过了,下面开始调用jenkins接口进行发版", data)
			//	_, err = CreateRedisInstance("set", data, "123")
			//	fmt.Println("已经成功的将key写入到redis")
			//	fmt.Println("====================================================")
			//	//使用获取到的三个参数去构建jenkins的job
			//	//jenkins.BuildHandler(jobName, changeTpe, gitlabSourceBranch)
			//	//instanceCodeList = append(instanceCodeList, data)
			//} else {
			//	fmt.Println("这个实例还没有被审批通过", data)
			//	fmt.Println("====================================================")
			//}
		} else {
			fmt.Println()
			fmt.Println()
			fmt.Println("这个审批已经处理过了", data)
			fmt.Println()
			fmt.Println()
		}
	}

}
