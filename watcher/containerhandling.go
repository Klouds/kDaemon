package watcher

import (
	"encoding/json"
	"errors"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/superordinate/kDaemon/database"
	"github.com/superordinate/kDaemon/logging"
	"github.com/superordinate/kDaemon/models"
	"strings"
)

//Commands
func AddContainer(job *Job) {

	job.InUse = true

	newcontainer := models.Container{}
	newcontainer.Status = "Initialized"

	decoder := json.NewDecoder(strings.NewReader(job.Body))
	err := decoder.Decode(&newcontainer)
	if err != nil {
		logging.Log(err)
		job.Complete = true //bad information, don't try to launch again
		return
	}

	/* Determine node to launch on */
	node, err := DetermineBestNodeForLaunch()
	if err != nil {
		logging.Log(err)
		job.Complete = false //bad node, so try to launch in the future
		job.InUse = false
		return
	}

	/* Get the application information */
	app, err := database.GetApplication(newcontainer.ApplicationID)
	if err != nil {
		logging.Log(err)
		job.Complete = true //Application doesn't exist, don't try to launch in the future
		return
	}

	//Launch the container on the given node

	err = LaunchAppOnNode(app, node, &newcontainer)
	if err != nil {
		logging.Log(err)
		job.Complete = true //Something went wrong, container name conflicts happen here
		job.InUse = false
		return
	}

	//save container information to database.
	job.Complete = true
	return

}

func LaunchAppOnNode(app *models.Application, node *models.Node, cont *models.Container) error {

	client, err := docker.NewClient(node.DIPAddr + ":" + node.DPort)

	if err != nil {
		logging.Log(err)
	}

	//Does image exist
	if imageExists(app.DockerImage, client) == nil {
		//does container exist
		if containerExists(cont.Name, client) == nil {
			//delete container
			if removeContainer(cont.Name, client) != nil {
				return errors.New("Container exists but cannot be removed.")
			}
		} else {
			logging.Log("Container doesn't exist")
		}

	} else { //image doesnt exist
		//pull image
		if pullImage(app.DockerImage, client) == nil {

			return errors.New("Pulling image")
		} else {

		}

	}

	//try to create container
	existingContainer := database.GetContainerByName(cont.Name)

	var useContainer *models.Container

	if existingContainer != nil {
		useContainer = existingContainer
	} else {
		useContainer, _, err = database.CreateContainer(cont)
		logging.Log(useContainer.Status)
		if err != nil {
			logging.Log("LC > Could not create the container on the database.")
		}
	}

	if useContainer.Status != "LAUNCHED" &&
		createContainer(cont.Name, app, client) == nil {

		logging.Log("Created Container")

		//start container
		if startContainer(cont.Name, client) == nil {
			logging.Log("starting Container")
			node.ContainerCount = node.ContainerCount + 1
			useContainer.ContainerID = cont.Name
			useContainer.NodeID = node.Id
			useContainer.ApplicationID = app.Id
			useContainer.IsEnabled = true
			useContainer.Status = "LAUNCHED"

			_, err := database.UpdateContainer(useContainer)

			if err != nil {
				logging.Log("LC > Could not save the container to the database.")
			}

			_, err = database.UpdateNode(node)

			if err != nil {
				logging.Log("LC > Could not save the node to the database.")
			}

		} else {
			return errors.New("Tried to start container2, but it failed.")
		}
	} else {
		return errors.New("Tried to create container, but it failed.")
	}

	return nil
}

func RemoveContainer(job *Job) {
	job.InUse = true

	//Get container info
	newcontainer := models.Container{}
	decoder := json.NewDecoder(strings.NewReader(job.Body))
	err := decoder.Decode(&newcontainer)
	if err != nil {
		logging.Log(err)
		job.Complete = true //bad information, don't try to launch again
		return
	}

	node, err := database.GetNode(newcontainer.NodeID)
	logging.Log("WHAT IS NODE ID", newcontainer.NodeID)

	if err != nil {
		logging.Log("RC > Node doesn't exist in Database")
		database.DeleteContainer(newcontainer.Id)
		job.Complete = true
		job.InUse = false
		return
	}

	if RemoveContainerFromNode(node, &newcontainer) == nil {
		//if successful

		node.ContainerCount = node.ContainerCount - 1
		database.UpdateNode(node)

		logging.Log("RC > Container Removed from node successfully")
		job.Complete = true
	} else {
		logging.Log("RC > Container removal unsuccessful. Must not be running on Node.")
		job.Complete = false
		job.InUse = false
	}
}

func RemoveContainerFromNode(node *models.Node, cont *models.Container) error {

	client, err := docker.NewClient(node.DIPAddr + ":" + node.DPort)
	logging.Log("REMOVE:", node)
	if err != nil {
		logging.Log(node)
		logging.Log("THE THING IS BROKEN")
	}

	if containerExists(cont.Name, client) == nil {
		//delete container
		if removeContainer(cont.Name, client) != nil {
			return errors.New("RC > Container exists but cannot be removed.")
		}
	} else {
		logging.Log("RC > Container doesn't exist")
	}

	err = database.DeleteContainer(cont.Id)

	return err
}

func createContainer(name string, app *models.Application, client *docker.Client) error {
	ports := app.GetPorts()

	port := ports[0] + "/tcp"
	exposedPort := map[docker.Port]struct{}{
		docker.Port(port): {}}
	portbindings := map[docker.Port][]docker.PortBinding{
		docker.Port(port): {}}

	//try to create container
	containeropts := docker.CreateContainerOptions{
		Name: name,
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

	_, err := client.CreateContainer(containeropts)

	return err
}

func startContainer(name string, client *docker.Client) error {
	return client.StartContainer(name, nil)
}

func imageExists(image string, client *docker.Client) error {
	img, err := client.InspectImage(image)

	if img == nil {
		err = errors.New("Image doesn't exist")
	} else {
		err = nil
	}

	return err
}

func containerExists(name string, client *docker.Client) error {
	cont, err := client.InspectContainer(name)

	if cont != nil {
		err = nil
	} else {
		err = errors.New("Container exists")
	}

	return err
}

func removeContainer(name string, client *docker.Client) error {

	logging.Log("Removing container")
	rmcontopts := docker.RemoveContainerOptions{
		ID:            name,
		RemoveVolumes: true,
		Force:         true,
	}

	return client.RemoveContainer(rmcontopts)
}

func pullImage(image string, client *docker.Client) error {
	logging.Log("pulling image")

	imageopts := docker.PullImageOptions{
		Repository: image,
	}

	err := client.PullImage(imageopts, docker.AuthConfiguration{})

	return err

}
