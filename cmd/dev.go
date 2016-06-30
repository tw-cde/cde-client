package cmd

import (
	"fmt"
	"github.com/sjkyspa/stacks/client/pkg"
	"github.com/sjkyspa/stacks/client/config"
	"github.com/sjkyspa/stacks/controller/api/api"
	"github.com/sjkyspa/stacks/controller/api/net"
	"net/url"
	"io/ioutil"
	"os/exec"
	//"bytes"
	//"strings"
	"os"
	"bytes"
	"strings"
)

func DevUp() error {

	if (!git.IsGitDirectory()) {
		return fmt.Errorf("Execute inside the app dir")
	}

	configRepository := config.NewConfigRepository(func(error) {})
	appRepository := api.NewAppRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))
	uri, err := url.Parse(configRepository.Endpoint())
	appId, err := git.DetectAppName(uri.Host)

	if err != nil || appId == "" {
		return fmt.Errorf("Please use the -remote to specfiy the app")
	}

	app, err := appRepository.GetApp(appId)
	if err != nil {
		return err
	}

	stack, err := app.GetStack()
	if err != nil {
		return err
	}

	f, err := toCompose(stack)
	if err != nil {
		return err
	}

	dockerComposeUp := exec.Command("docker-compose", "-f", f, "up", "-d")

	var out bytes.Buffer
	var errout bytes.Buffer
	dockerComposeUp.Stdin = strings.NewReader("test")
	dockerComposeUp.Stdout = &out
	dockerComposeUp.Stderr = &errout
	err = dockerComposeUp.Run()
	if err != nil {
		return err
	}
	fmt.Println(out.String())
	fmt.Println(errout.String())

	dockerExec := exec.Command("docker", "exec", "-it", "template_runtime_1", "sh")
	dockerExec.Stdin = os.Stdin
	dockerExec.Stderr = os.Stderr
	dockerExec.Stdout = os.Stdout
	err = dockerExec.Run()
	if err != nil {
		return err
	}

	return nil
}

func toCompose(stack api.Stack) (string, error) {
	aa :=
	`version: '2'
services:
  runtime:
    image: hub.deepi.cn/jersey-mysql-build
    entrypoint: /bin/sh
    command: -c 'tail -f /dev/null'
    volumes:
      - /Mac/workspace/tmp/cde-stacks/jersey-mysql/template:/codee
      - /var/run/docker.sock:/var/run/docker.sock
    links:
      - mysql
  mysql:
    image: tutum/mysql
    ports:
     - 5000:5000`

	err := ioutil.WriteFile("dockercompose.yml", []byte(aa), 0600)
	if err != nil {
		return "", err
	}

	return "dockercompose.yml", nil
}

func DevDown() error {
	fmt.Println("dev down")
	return nil
}

func DevDestroy() error {
	fmt.Println("dev destroy")
	return nil
}
