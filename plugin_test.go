package main

import (
	"testing"
)

func TestPlugin(t *testing.T) {
	var err error
	p := Plugin{}
	//p.Config.Debug = true
	//FIXME
	p.Config.AccessToken = "4d18a039145d00dad6fcf4829def90dea9c02c2e125cd2fe5239fd18ccaba55c"
	p.Config.Secret = "SEC25781791183b500cda5929fd365d6befb6d7b13705dc02253224a04b9f609e5e"
	p.Drone.Commit.Sha = "1234566"
	p.Config.MsgType = "text"
	p.Config.MsgType = "link"
	p.Extra.Color.WithColor = true
	p.Extra.Color.FailureColor = "#555555"
	p.Extra.Color.SuccessColor = "#222222"
	p.Extra.Pic.WithPic = true
	p.Extra.Pic.FailurePicURL = ""
	p.Extra.Pic.SuccessPicURL = ""
	p.Extra.LinkSha = true
	p.Config.MsgType = "markdown"
	p.Extra.Db.DbLog = true
	p.Extra.Db.DbType = "mysql"
	p.Extra.Db.DbName = "notelog"
	p.Extra.Db.DbHost = "localhost"
	p.Extra.Db.DbPort = 3306
	p.Extra.Db.DbUsername = "dyb"
	p.Extra.Db.DbPassword = "rootdyb"

	p.Drone.Build = struct {
		Status     string
		Link       string
		RepoName   string
		Image      string
		Tag        string
		Port       string
		StartAt    int64
		FinishedAt int64
		Stage      string
		Event      string
	}{Status: "success", Link: "http://172.16.100.99:9053/1/2/3//33//44//55//66//77//88", RepoName: "", Image: "", Tag: "", Port: "", StartAt: 0, FinishedAt: 0, Stage: "build", Event: "push"}

	p.Drone.Commit.Branch = "daiyb"
	p.Drone.Commit.Message = "测试提交"
	p.Drone.Commit.Link = "https://www.baidu.com"
	p.Drone.Commit.Authors = struct {
		Avatar string
		Email  string
		Name   string
	}{Avatar: "", Email: "", Name: "daiyb"}

	err = p.Exec()
	if nil != err {
		t.Errorf("%v", err)
		return
	}
	t.Log("plugin testing finished")
}
