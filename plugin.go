package main

import (
	"drone-message/model"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type (
	// Repo `repo base info`
	Repo struct {
		FullName string // repository full name
		ModName  string // repository module name
	}

	// Build `build info`
	Build struct {
		Status     string // providers the current build status
		Link       string // providers the current build link
		RepoName   string // docker repo
		Image      string // docker image name
		Tag        string // providers the current build tag
		Port       string // drone ci port
		StartAt    int64  // build start at (unix timestamp)
		FinishedAt int64  // build finish at (unix timestamp)
		Stage      string // drone stage
		Event      string // drone event
	}

	// Commit `commit info`
	Commit struct {
		Branch  string //  providers the branch for the current commit
		Link    string //  providers the http link to the current commit in the remote source code management system(e.g.GitHub)
		Message string //  providers the commit message for the current build
		Sha     string //  providers the commit sha for the current build
		Authors CommitAuthors
	}

	// CommitAuthors `commit author info`
	CommitAuthors struct {
		Avatar string //  providers the author avatar for the current commit
		Email  string //  providers the author email for the current commit
		Name   string //  providers the author name for the current commit
	}

	// Drone `drone info`
	Drone struct {
		Repo   Repo
		Build  Build
		Commit Commit
	}

	// Config `plugin private config`
	Config struct {
		Debug       bool
		AccessToken string
		IsAtALL     bool
		Mobiles     string
		Username    string
		MsgType     string
		Secret      string
	}

	// Extra `extra variables`
	Extra struct {
		Color   ExtraColor
		Pic     ExtraPic
		Db      ExtraDb
		LinkSha bool
	}

	// ExtraPic `extra config for pic`
	ExtraPic struct {
		WithPic       bool
		SuccessPicURL string
		FailurePicURL string
	}

	// ExtraColor `extra config for color`
	ExtraColor struct {
		WithColor    bool
		SuccessColor string
		FailureColor string
	}

	ExtraDb struct {
		DbLog      bool
		DbType     string
		DbName     string
		DbHost     string
		DbPort     int64
		DbUsername string
		DbPassword string
		DbTable    string
	}

	// Plugin `plugin all config`
	Plugin struct {
		Drone  Drone
		Config Config
		Extra  Extra
	}

	// Envfile
	Envfile struct {
		ConfigPkg string   `yaml:"configPkg"`
		CheckList []string `yaml:"checkList"`
		ImageList []string `yaml:"imageList"`
	}
)

// Exec `execute webhook`
func (p *Plugin) Exec() error {
	fmt.Printf("%+v", p.Config)
	if p.Config.Debug {
		for _, e := range os.Environ() {
			log.Println(e)
		}
	}
	var err error
	if p.Config.AccessToken == "" {
		msg := "missing dingtalk access token"
		return errors.New(msg)
	}

	if p.Config.Secret == "" {
		msg := "missing dingtalk secret"
		return errors.New(msg)
	}
	DbLog(p)

	if 6 > len(p.Drone.Commit.Sha) {
		return errors.New("commit sha cannot short than 6")
	}
	newWebhook := NewWebHook(p.Config.AccessToken, p.Config.Secret)
	mobiles := strings.Split(p.Config.Mobiles, ",")
	if p.Drone.Repo.ModName != "" {
		if checkModuleNmae(p.Drone.Repo.ModName) {
			switch strings.ToLower(p.Config.MsgType) {
			case "markdown":
				err = newWebhook.SendMarkdownMsg("新的构建通知", p.baseTpl(), p.Config.IsAtALL, mobiles...)
			case "text":
				err = newWebhook.SendTextMsg(p.baseTpl(), p.Config.IsAtALL, mobiles...)
			case "link":
				err = newWebhook.SendLinkMsg(p.Drone.Build.Status, p.baseTpl(), p.Drone.Commit.Authors.Avatar, p.Drone.Build.Link)
			default:
				msg := "not support message type"
				err = errors.New(msg)
			}

			if err == nil {
				log.Println("send message success!")
			}
		}
	} else {
		switch strings.ToLower(p.Config.MsgType) {
		case "markdown":
			err = newWebhook.SendMarkdownMsg("新的构建通知", p.baseTpl(), p.Config.IsAtALL, mobiles...)
		case "text":
			err = newWebhook.SendTextMsg(p.baseTpl(), p.Config.IsAtALL, mobiles...)
		case "link":
			err = newWebhook.SendLinkMsg(p.Drone.Build.Status, p.baseTpl(), p.Drone.Commit.Authors.Avatar, p.Drone.Build.Link)
		default:
			msg := "not support message type"
			err = errors.New(msg)
		}

		if err == nil {
			log.Println("send message success!")
		}
	}

	aa, _ := json.Marshal(p.Drone.Build)
	if p.Extra.Db.DbLog {
		_ = WriteLog(&model.PublishLog{
			CommitId:   p.Drone.Commit.Sha,
			CommitLink: p.Drone.Commit.Link,
			BuildLink:  p.Drone.Build.Link,
			Author:     p.Drone.Commit.Authors.Name,
			Branch:     p.Drone.Commit.Branch,
			Message:    p.Drone.Commit.Message,
			Event:      p.Drone.Build.Event,
			Remark:     string(aa),
		})
	}
	return err
}

// actionCard `output the tpl of actionCard`
func (p *Plugin) actionCardTpl() string {
	var tpl string

	// title
	title := strings.Title(p.Drone.Repo.FullName)
	// with color on title
	if p.Extra.Color.WithColor {
		title = fmt.Sprintf("<font color=%s>%s</font>", p.getColor(), title)
	}

	tpl = fmt.Sprintf("# %s \n", title)

	branch := fmt.Sprintf("> %s 分支", strings.Title(p.Drone.Commit.Branch))
	tpl += branch + "\n\n"

	// with pic
	if p.Extra.Pic.WithPic {
		tpl += fmt.Sprintf("![%s](%s)\n\n",
			p.Drone.Build.Status,
			p.getPicURL())
	}

	// commit message
	commitMsg := fmt.Sprintf("Commit 信息：%s", p.Drone.Commit.Message)
	if p.Extra.Color.WithColor {
		commitMsg = fmt.Sprintf("<font color=%s>%s</font>", p.getColor(), commitMsg)
	}
	tpl += commitMsg + "\n\n"
	// author info
	authorInfo := fmt.Sprintf("提交者：`%s(%s)`", p.Drone.Commit.Authors.Name, p.Drone.Commit.Authors.Email)
	tpl += authorInfo + "\n\n"

	// sha info
	commitSha := p.Drone.Commit.Sha
	if p.Extra.LinkSha {
		commitSha = fmt.Sprintf("[点击查看 Commit %s 信息](%s)", commitSha[:6], p.Drone.Commit.Link)
	}
	tpl += commitSha + "\n\n"

	// docker info
	log.Println(fmt.Sprintf("repo name:%s", p.Drone.Build.RepoName))
	log.Println(fmt.Sprintf("repo name:%s", p.Drone.Build.Image))

	return tpl
}

// markdownTpl `output the tpl of markdown`
func (p *Plugin) markdownTpl() string {
	var tpl string

	// title
	title := strings.Title(p.Drone.Repo.FullName)
	// with color on title
	if p.Extra.Color.WithColor {
		title = fmt.Sprintf("<font color=%s>%s</font>", p.getColor(), title)
	}

	tpl += fmt.Sprintf("# %s \n", title)

	// branch or tag
	var branch string
	if p.Drone.Build.Tag == "" {
		branch = fmt.Sprintf("> %s -> %s 分支", p.Drone.Build.Event, p.Drone.Commit.Branch)
	} else {
		branch = fmt.Sprintf("> 发布标签： %s", p.Drone.Build.Tag)
	}
	tpl += branch + "\n\n"

	// module name
	if p.Drone.Repo.ModName != "" {
		modname := fmt.Sprintf("模块：%s", strings.Title(p.Drone.Repo.ModName))
		tpl += modname + "\n\n"
	}

	// author info
	authorInfo := fmt.Sprintf("提交者：`%s(%s)`", p.Drone.Commit.Authors.Name, p.Drone.Commit.Authors.Email)
	tpl += authorInfo + "\n\n"

	// build take seconds info
	takeSeconds := fmt.Sprintf("构建耗时：`%ds`", p.Drone.Build.FinishedAt-p.Drone.Build.StartAt)
	tpl += takeSeconds + "\n\n"

	// stage info
	stageName := fmt.Sprintf("构建阶段：`%s`", p.Drone.Build.Stage)
	tpl += stageName + "\n\n"

	// with pic
	if p.Extra.Pic.WithPic {
		tpl += fmt.Sprintf("![%s](%s)\n\n",
			p.Drone.Build.Status,
			p.getPicURL())
	}

	// commit message
	commitMsg := fmt.Sprintf("Commit 信息：%s", p.Drone.Commit.Message)
	if p.Extra.Color.WithColor {
		commitMsg = fmt.Sprintf("<font color=%s>%s</font>", p.getColor(), commitMsg)
	}
	tpl += commitMsg + "\n\n"

	// sha info
	commitSha := p.Drone.Commit.Sha
	// commitSha[:6]
	if p.Extra.LinkSha {
		commitSha = fmt.Sprintf("[查看 Commit 信息](%s)", p.Drone.Commit.Link)
	}

	// build detail link
	// add port
	s := strings.Split(p.Drone.Build.Link, "//")
	s2 := strings.Split(s[1], "/")
	buildLink := fmt.Sprintf("%s//%s/%s/%s/%s", s[0], s2[0], s2[1], s2[2], s2[3])
	buildDetail := fmt.Sprintf("[查看构建信息](%s)", buildLink)

	var repos []string
	if p.Drone.Repo.ModName != "" {
		envfile := Envfile{}
		envfile.ReadYaml("./env.yaml")
		repos = envfile.ImageList
	}
	//var repos []string

	//读取文件
	//b, err := ioutil.ReadFile("repo.txt")
	//if err != nil {
	//	fmt.Println("ioutil ReadFile error: ", err)
	//}

	//imagepath := string(b)

	if p.Extra.Db.DbLog {
		dbMsg := fmt.Sprintf("数据库记录信息:%v", model.DbMsg)
		tpl += dbMsg + "\n\n"
	}
	// deploy link
	if len(repos) > 0 {
		//repos := strings.Split(string(b), ",")
		for _, reponame := range repos {
			fmt.Println("repo:", reponame)
			content := strings.Split(reponame, ":")[0]
			RepoName := strings.Split(content, "/")[1]
			Image := strings.Split(content, "/")[2]
			repoinfo := fmt.Sprintf("---\n\n > %s 镜像地址：`%s`", Image, reponame)
			deployUrl := fmt.Sprintf("https://devops.keking.cn/#/k8s/imagetag?namespace=%s&reponame=%s", RepoName, Image)
			deployLink := fmt.Sprintf("> [部署 「%s」 ](%s)", Image, deployUrl)
			tpl += repoinfo + "\n\n" + deployLink + "\n\n"
		}
	}

	tpl += "---\n\n\n\n" + commitSha + " | " + buildDetail + "\n\n"

	//fmt.Println(tpl)
	return tpl
}

func (p *Plugin) baseTpl() string {
	tpl := ""
	switch strings.ToLower(p.Config.MsgType) {
	case "markdown":
		tpl = p.markdownTpl()
	case "text":
		tpl = fmt.Sprintf(`[%s] %s
%s (%s)
@%s
%s (%s)
`,
			p.Drone.Build.Status,
			strings.TrimSpace(p.Drone.Commit.Message),
			p.Drone.Repo.FullName,
			p.Drone.Commit.Branch,
			p.Drone.Commit.Sha,
			p.Drone.Commit.Authors.Name,
			p.Drone.Commit.Authors.Email)
	case "link":
		tpl = fmt.Sprintf(`%s(%s) @%s %s(%s)`,
			p.Drone.Repo.FullName,
			p.Drone.Commit.Branch,
			p.Drone.Commit.Sha[:6],
			p.Drone.Commit.Authors.Name,
			p.Drone.Commit.Authors.Email)
	case "actioncard":
		tpl = p.actionCardTpl()
	}

	return tpl
}

/**
get picture url
*/
func (p *Plugin) getPicURL() string {
	pics := make(map[string]string)
	//  success picture url
	pics["success"] = "https://dci-file.bj.bcebos.com/fonchain-main/prod/image/0/comon/c663a410-521e-43b2-adbc-9fb00d5b0a91.png"
	if p.Extra.Pic.SuccessPicURL != "" {
		pics["success"] = p.Extra.Pic.SuccessPicURL
	}
	//  failure picture url
	pics["failure"] = "https://dci-file.bj.bcebos.com/fonchain-main/prod/image/0/comon/320f07c4-ea53-4980-b194-3e258e927549.png"
	if p.Extra.Pic.FailurePicURL != "" {
		pics["failure"] = p.Extra.Pic.FailurePicURL
	}

	url, ok := pics[p.Drone.Build.Status]
	if ok {
		return url
	}

	return ""
}

/**
get color for message title
*/
func (p *Plugin) getColor() string {
	colors := make(map[string]string)
	//  success color
	colors["success"] = "#008000"
	if p.Extra.Color.SuccessColor != "" {
		colors["success"] = "#" + p.Extra.Color.SuccessColor
	}
	//  failure color
	colors["failure"] = "#FF0000"
	if p.Extra.Color.FailureColor != "" {
		colors["failure"] = "#" + p.Extra.Color.FailureColor
	}

	color, ok := colors[p.Drone.Build.Status]
	if ok {
		return color
	}

	return ""
}

func (c *Envfile) ReadYaml(f string) {
	buffer, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatalf(err.Error())
	}
	err = yaml.Unmarshal(buffer, &c)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func checkModuleNmae(name string) bool {
	envfile := Envfile{}
	envfile.ReadYaml("./env.yaml")

	if len(envfile.CheckList) == 0 {
		fmt.Println("+ skip module package check")
	} else {
		modname := envfile.CheckList
		var whether bool
		for _, mod := range modname {
			if mod == name {
				whether = true
				continue
			}
		}
		if whether {
			fmt.Printf("+ Name matching succeeded, 「%s」 continue !\n", name)
			return true
		} else {
			fmt.Println("+ No matching name,jump step")
			return false
		}
	}
	return false
}
