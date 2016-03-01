package watcher

import (
	"errors"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/klouds/kDaemon/logging"
	"github.com/klouds/kDaemon/models"
)

//This object will handle making calls to docker.
type dockerHandler struct {
	client *docker.Client
}

func NewDockerHandler(adress string, port string) (*dockerHandler, error) {
	//Initialize connection to docker
	newclient, err := docker.NewClient(adress + ":" + port)

	if err != nil {
		logging.Log("Couldnt establish docker connection")
		return nil, errors.New("Couldn't establish docker connection")
	}

	//Create the new Handler object
	newdockerhandler := &dockerHandler{
		client: newclient,
	}

	return newdockerhandler, nil
}

//Sweet, so we're creating handlers for our docker client but can't
//do anything with it yet. Let's add some functionality. First, we'll give
//it the ability to detect whether an image is present on the host machine.

//simple message that will poll the host and check whether a container exists.
func (dh *dockerHandler) IsImagePresent(imagename string) bool {
	img, _ := dh.client.InspectImage(imagename)

	if img == nil {
		return false
	}
	return true
}

//Assuming the image doesn't exist, let's give it a function to pull an image
//onto the server for use
// Now this function will return true or false now but will be expanded
// to be its own channel so it can block
func (dh *dockerHandler) PullImage(imagename string) bool {
	imageopts := docker.PullImageOptions{
		Repository: imagename,
	}

	err := dh.client.PullImage(imageopts, docker.AuthConfiguration{})

	if err != nil {
		return false
	}

	return true
}

//After we pull the image, we're going to have to create the container
//this will do just that--- WAIT! We can't do that yet. First we have to check
//whether or not the container already exists in the system!

//Let's check for an existing container then
func (dh *dockerHandler) DoesContainerExist(containerid string) bool {
	cont, _ := dh.client.InspectContainer(containerid)

	if cont == nil {
		return false
	}

	return true
}

//Phew, well that was hard. Let's assume the container didn't exist,
//let's create it now
// How are we going to get the container options?
// what do i need?
// let's go check.
//
//Okay I think I can do this almost as it was before, but changing the returns
func (dh *dockerHandler) CreateContainer(containerid string, app *models.Application) bool {
	logging.Log("Whats up")
	ports := app.GetPorts()

	port := ports[0] + "/tcp"
	exposedPort := map[docker.Port]struct{}{
		docker.Port(port): {}}
	portbindings := map[docker.Port][]docker.PortBinding{
		docker.Port(port): {}}

	//try to create container
	containeropts := docker.CreateContainerOptions{
		Name: containerid,
		Config: &docker.Config{
			ExposedPorts: exposedPort,
			Image:        app.DockerImage,
		},
		HostConfig: &docker.HostConfig{
			PublishAllPorts: true,
			PortBindings:    portbindings,
			Privileged:      false,
		},
	}

	_, err := dh.client.CreateContainer(containeropts)

	if err != nil {
		return false
	}

	return true
}

//Now that the image has been pulled, and the container has been created,
//let's start the bitch up!
func (dh *dockerHandler) StartContainer(containerid string) bool {
	err := dh.client.StartContainer(containerid, nil)

	if err != nil {
		return false
	}

	return true
}

//Everything compiles! Alright! We should in theory be able
//to launch a container now.
//
//Let's wire up a container launch command, as we're going
//to separate the container creation and launching/stopping.
//
//Lets add a stop container command to the handler
//
func (dh *dockerHandler) StopContainer(containerid string) bool {
	err := dh.client.StopContainer(containerid, 1)

	if err != nil {
		return false
	}

	return true
}

//Hey look, I did it!
//wait, now im done.
//lets go back and call this beauty :P

//We need a function to get an exposed port of a container id
func (dh *dockerHandler) InspectContainer(containerid string) *docker.Container {

	cont, err := dh.client.InspectContainer(containerid)

	if err != nil {
		logging.Log("Unable to inspect container. Doesn't exist?")
		return nil
	}

	return cont
}
