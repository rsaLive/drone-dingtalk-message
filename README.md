# Drone CI DingTalk Message Plugin

### Drone CI Plugin Config
`1.0.x`
```yaml
kind: pipeline
name: default

steps:
...
- name: dingtalk
  image: meetdocker/drone-dingtalk-message
  pull: if-not-exists
  settings:
    token:
      from_secret: dingtalk_token
    secret:
      from_secret: dingtalk_secret
    debug: true
    type: markdown
    message_color: true
    message_pic: true
    sha_link: true
    drone_port: 30000
    db_log: true
    db_type: mysql
    db_name: cicd
    db_host: "localhost"
    db_port: 3306
    db_username: dyb
    db_password: rootdyb
  when:
    status: [failure, success]
```

### Plugin Parameter Reference
`token`(required)

String. Access token for group bot. (you can get the access token when you add a bot in a group)

`secret`(required)

String. Secret for group bot. (you can get the secret when you add a bot in a group)

`drone_port`

Int. Drone run port.

`type`(required)

String. Message type, plan support text, markdown, link and action card, but due to time issue, it's only support `markdown` and `text` now, and you can get the best experience by use markdown.

`message_color`(when `type=markdown`)

Boolean value. This option can change the title and commit message color if turn on.

`success_color`(when `message_color=true`)

String. You can customize the color for the `build success` message by this option, you should input a hex color, example: `008000`.

`failure_color`(when `message_color=true`)

String. You can customize the color for the `build success` message by this option, you should input a hex color, example: `FF0000`.

`sha_link`(when `type=markdown`)

Boolean value. This option can link the sha to your source page when it turn on.

`message_pic`(when `type=markdown`)

Boolean value. If this option turn on,  it will embed a image into the message.

`success_pic`(when `message_pic=true`)

String. You can customize the picture for the `build success` message by this option.

`failure_pic`(when `message_pic=true`)

String. You can customize the picture for the `build failure` message by this option.

### Development
- get this repo
```shell
go get github.com/happy-python/drone-dingtalk-message
```
- build
```shell
go build .
```
- run
```shell
./drone-dingtalk-message -h
```
- docker build
```shell
docker build .
```