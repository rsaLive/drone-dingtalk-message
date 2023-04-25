package main

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli"
	"log"
	"os"
)

var Version = "0.1.1"

func main() {
	app := cli.NewApp()
	app.Name = "Drone Dingtalk Message Plugin"
	app.Usage = "Sending message to Dingtalk group by robot using webhook"
	app.Action = run
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "config.debug",
			Usage:  "debug mode",
			EnvVar: "PLUGIN_DEBUG",
		},
		cli.StringFlag{
			Name:   "config.token,access_token,token",
			Usage:  "dingtalk webhook access token",
			EnvVar: "PLUGIN_ACCESS_TOKEN,PLUGIN_TOKEN",
		},
		cli.StringFlag{
			Name:   "config.secret,access_secret,secret",
			Usage:  "dingtalk webhook secret",
			EnvVar: "PLUGIN_ACCESS_SECRET,PLUGIN_SECRET",
		},
		cli.StringFlag{
			Name:   "config.lang",
			Value:  "zh_CN",
			Usage:  "the lang display (zh_CN or en_US, zh_CN is default)",
			EnvVar: "PLUGIN_LANG",
		},
		cli.StringFlag{
			Name:   "config.message.type,message_type",
			Usage:  "dingtalk message type, like text, markdown, action card, link and feed card...",
			EnvVar: "PLUGIN_MSG_TYPE,PLUGIN_TYPE,PLUGIN_MESSAGE_TYPE",
		},
		cli.StringFlag{
			Name:   "config.message.at.all",
			Usage:  "at all in a message(only text and markdown type message can at)",
			EnvVar: "PLUGIN_MSG_AT_ALL",
		},
		cli.StringFlag{
			Name:   "config.message.at.mobiles",
			Usage:  "at someone in a dingtalk group need this guy bind's mobile",
			EnvVar: "PLUGIN_MSG_AT_MOBILES",
		},
		cli.StringFlag{
			Name:   "commit.author.avatar",
			Usage:  "providers the author avatar url for the current commit",
			EnvVar: "DRONE_COMMIT_AUTHOR_AVATAR",
		},
		cli.StringFlag{
			Name:   "commit.author.email",
			Usage:  "providers the author email for the current commit",
			EnvVar: "DRONE_COMMIT_AUTHOR_EMAIL",
		},
		cli.StringFlag{
			Name:   "commit.author.name",
			Usage:  "providers the author name for the current commit",
			EnvVar: "DRONE_COMMIT_AUTHOR",
		},
		cli.StringFlag{
			Name:   "commit.branch",
			Usage:  "providers the branch for the current build",
			EnvVar: "DRONE_COMMIT_BRANCH",
			Value:  "master",
		},
		cli.StringFlag{
			Name:   "commit.link",
			Usage:  "providers the http link to the current commit in the remote source code management system(e.g.GitHub)",
			EnvVar: "DRONE_COMMIT_LINK",
		},
		cli.StringFlag{
			Name:   "commit.message",
			Usage:  "providers the commit message for the current build",
			EnvVar: "DRONE_COMMIT_MESSAGE",
		},
		cli.StringFlag{
			Name:   "commit.sha",
			Usage:  "providers the commit sha for the current build",
			EnvVar: "DRONE_COMMIT_SHA",
		},
		cli.StringFlag{
			Name:   "repo.fullname",
			Usage:  "providers the full name of the repository",
			EnvVar: "DRONE_REPO",
		},
		// DRONE_TAG
		cli.StringFlag{
			Name:   "repo.tag",
			Usage:  "providers the tag of the repository",
			EnvVar: "DRONE_TAG",
		},
		cli.StringFlag{
			Name:   "build.status",
			Usage:  "build status",
			Value:  "success",
			EnvVar: "DRONE_BUILD_STATUS",
		},
		cli.StringFlag{
			Name:   "build.link",
			Usage:  "build link",
			EnvVar: "DRONE_BUILD_LINK",
		},
		cli.StringFlag{
			Name:   "build.port",
			Usage:  "drone build port",
			EnvVar: "PLUGIN_DRONE_PORT",
		},
		cli.StringFlag{
			Name:   "plugin.build.reponamespace",
			Usage:  "build step docker repository",
			EnvVar: "PLUGIN_REPO_NAMESPACE",
		},
		cli.StringFlag{
			Name:   "plugin.build.imagename",
			Usage:  "build step docker repository",
			EnvVar: "PLUGIN_IMAGE_NAME",
		},
		cli.StringFlag{
			Name:   "config.success.pic.url",
			Usage:  "config success picture url",
			EnvVar: "SUCCESS_PICTURE_URL,PLUGIN_SUCCESS_PIC",
		},
		cli.StringFlag{
			Name:   "config.failure.pic.url",
			Usage:  "config failure picture url",
			EnvVar: "FAILURE_PICTURE_URL,PLUGIN_FAILURE_PIC",
		},
		cli.StringFlag{
			Name:   "config.success.color",
			Usage:  "config success color for title in markdown",
			EnvVar: "SUCCESS_COLOR,PLUGIN_SUCCESS_COLOR",
		},
		cli.StringFlag{
			Name:   "config.failure.color",
			Usage:  "config failure color for title in markdown",
			EnvVar: "FAILURE_COLOR,PLUGIN_FAILURE_COLOR",
		},
		cli.BoolFlag{
			Name:   "config.message.color",
			Usage:  "configure the message with color or not",
			EnvVar: "PLUGIN_COLOR,PLUGIN_MESSAGE_COLOR",
		},
		cli.BoolFlag{
			Name:   "config.message.pic",
			Usage:  "configure the message with picture or not",
			EnvVar: "PLUGIN_PIC,PLUGIN_MESSAGE_PIC",
		},
		cli.BoolFlag{
			Name:   "config.message.sha.link",
			Usage:  "link sha source page or not",
			EnvVar: "PLUGIN_SHA_LINK,PLUGIN_MESSAGE_SHA_LINK",
		},
		cli.BoolFlag{
			Name:   "config.db.log",
			Usage:  "write db log",
			EnvVar: "PLUGIN_DB_LOG",
		},
		cli.StringFlag{
			Name:   "config.db.type",
			Usage:  " db type",
			EnvVar: "PLUGIN_DB_TYPE",
		},
		cli.StringFlag{
			Name:   "config.db.name",
			Usage:  " db name",
			EnvVar: "PLUGIN_DB_NAME",
		},
		cli.StringFlag{
			Name:   "config.db.host",
			Usage:  " db host",
			EnvVar: "PLUGIN_DB_HOST",
		}, cli.Int64Flag{
			Name:   "config.db.port",
			Usage:  " db port",
			EnvVar: "PLUGIN_DB_PORT",
		},
		cli.StringFlag{
			Name:   "config.db.username",
			Usage:  " db name",
			EnvVar: "PLUGIN_DB_USERNAME",
		},
		cli.StringFlag{
			Name:   "config.db.password",
			Usage:  " db name",
			EnvVar: "PLUGIN_DB_PASSWORD",
		},
		cli.StringFlag{
			Name:   "config.db.table",
			Usage:  " db table",
			EnvVar: "PLUGIN_DB_TABLE",
		},
		cli.StringFlag{
			Name:   "module.name",
			Usage:  "git update package name",
			EnvVar: "PLUGIN_MODNAME",
		},
		cli.Int64Flag{
			Name:   "build.start",
			Usage:  "drone build start at",
			EnvVar: "DRONE_STAGE_STARTED",
		},
		cli.Int64Flag{
			Name:   "build.end",
			Usage:  "drone build finish at",
			EnvVar: "DRONE_STAGE_FINISHED",
		},
		cli.StringFlag{
			Name:   "build.stage",
			Usage:  "drone build stage name",
			EnvVar: "DRONE_STAGE_NAME",
		},
		cli.StringFlag{
			Name:   "build.event",
			Usage:  "drone build event",
			EnvVar: "DRONE_BUILD_EVENT",
		},
	}

	if err := app.Run(os.Args); nil != err {
		log.Println(err)
	}
}

//  run with args
func run(c *cli.Context) {
	plugin := Plugin{
		Drone: Drone{
			//  repo info
			Repo: Repo{
				FullName: c.String("repo.fullname"),
				ModName:  c.String("module.name"),
			},
			//  build info
			Build: Build{
				Status:     c.String("build.status"),
				Link:       c.String("build.link"),
				RepoName:   c.String("plugin.build.reponamespace"),
				Image:      c.String("plugin.build.imagename"),
				Tag:        c.String("repo.tag"),
				Port:       c.String("build.port"),
				StartAt:    c.Int64("build.start"),
				FinishedAt: c.Int64("build.end"),
				Stage:      c.String("build.stage"),
				Event:      c.String("build.event"),
			},
			Commit: Commit{
				Sha:     c.String("commit.sha"),
				Branch:  c.String("commit.branch"),
				Message: c.String("commit.message"),
				Link:    c.String("commit.link"),
				Authors: struct {
					Avatar string
					Email  string
					Name   string
				}{
					Avatar: c.String("commit.author.avatar"),
					Email:  c.String("commit.author.email"),
					Name:   c.String("commit.author.name"),
				},
			},
		},
		//  custom config
		Config: Config{
			AccessToken: c.String("config.token"),
			Secret:      c.String("config.secret"),
			IsAtALL:     c.Bool("config.message.at.all"),
			MsgType:     c.String("config.message.type"),
			Mobiles:     c.String("config.message.at.mobiles"),
			Debug:       c.Bool("config.debug"),
		},
		Extra: Extra{
			Pic: ExtraPic{
				WithPic:       c.Bool("config.message.pic"),
				SuccessPicURL: c.String("config.success.pic.url"),
				FailurePicURL: c.String("config.failure.pic.url"),
			},
			Color: ExtraColor{
				SuccessColor: c.String("config.success.color"),
				FailureColor: c.String("config.failure.color"),
				WithColor:    c.Bool("config.message.color"),
			},
			Db: ExtraDb{
				DbLog:      c.Bool("config.db.log"),
				DbType:     c.String("config.db.type"),
				DbName:     c.String("config.db.name"),
				DbHost:     c.String("config.db.host"),
				DbPort:     c.Int64("config.db.port"),
				DbUsername: c.String("config.db.username"),
				DbPassword: c.String("config.db.password"),
				DbTable:    c.String("config.db.table"),
			},
			LinkSha: c.Bool("config.message.sha.link"),
		},
	}

	if err := plugin.Exec(); nil != err {
		fmt.Println(err)
	}
}
