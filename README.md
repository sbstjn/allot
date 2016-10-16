# allot [![CircleCI](https://img.shields.io/circleci/project/sbstjn/allot.svg)](https://circleci.com/gh/sbstjn/allot) [![Go Report Card](https://goreportcard.com/badge/github.com/sbstjn/allot)](https://goreportcard.com/report/github.com/sbstjn/allot) [![Coverage Status](https://coveralls.io/repos/github/sbstjn/allot/badge.svg?branch=master)](https://coveralls.io/github/sbstjn/allot?branch=master)

**allot** is a small `Golang` library to match and parse commands with pre-defined strings. For example use **allot** to define a list of commands your CLI application or Slackbot supports and check if incoming requests are matching your commands.

The **allot** library supports placeholders and regular expressions for parameter matching and parsing.

## Usage

```go
cmd := allot.NewCommand("deploy <project:string> to <env:string>")
req := allot.NewRequest("deploy allot to live")

if cmd.Matches(req) {
    fmt.Printf(
      "Deploy \"%s\" to \"%s\"",
      cmd.ParamString(req, "project"),
      cmd.ParamString(req, "env"),
    )
}
```


## Examples

See the [hanu Slackbot](https://github.com/sbstjn/hanu) framework for a usecase for **allot**:

* [Host a Golang Slack bot on Heroku](https://sbstjn.com/host-golang-slackbot-on-heroku-with-hanu.html)

## Credits
 * [Go coverage script from Mathias Lafeldt](https://mlafeldt.github.io/blog/test-coverage-in-go/)
