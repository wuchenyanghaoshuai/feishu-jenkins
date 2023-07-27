package jenkins

import (
	"context"
	"github.com/bndr/gojenkins"
	"log"
	"net/http"
)

type BuildRequest struct {
	JobName            string `json:"jobName"`
	GitlabSourceBranch string `json:"gitlabSourceBranch"`
	ChangeType         string `json:"changeType"`
}

func BuildHandler(jobName, changeTpe, gitlabSourceBranch string) {
	// 创建 HTTP 客户端
	httpClient := &http.Client{}
	// 创建一个空的上下文对象
	ctx := context.Background()
	// 创建 Jenkins 实例
	jenkins, err := gojenkins.CreateJenkins(httpClient, "http://192.168.3.100:8080", "用户名", "密码").Init(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}

	// 获取指定 Job 的信息
	job, err := jenkins.GetJob(ctx, jobName)
	if err != nil {
		log.Fatal(err)
		return
	}

	// 为指定 Job 和分支构建
	params := map[string]string{"CHANGE_TYPE": changeTpe, "gitlabSourceBranch": gitlabSourceBranch}
	build, err := job.InvokeSimple(ctx, params)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("Build %d for job '%s' and branch '%s' is in queue", build, jobName, gitlabSourceBranch)

}
