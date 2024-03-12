# Drone CI DingTalk Message Plugin

### Drone CI Plugin Config
`1.0.x`
```yaml
kind: pipeline
name: default

steps:

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

- docker build
```shell
docker build .
```