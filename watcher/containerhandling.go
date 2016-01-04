package watcher

import (
	"github.com/superordinate/kDaemon/models"
	"github.com/superordinate/kDaemon/logging"
	"github.com/superordinate/kDaemon/database"
	docker "github.com/fsouza/go-dockerclient"
	"encoding/json"
	"strings"
)

//Commands
func AddContainer(job *Job){
	
	job.InUse = true

	newcontainer := models.Container{}
	newcontainer.Status = "Initialized"
	decoder := json.NewDecoder(strings.NewReader(job.Body))
	err := decoder.Decode(&newcontainer)
	if err != nil {
		logging.Log(err)
		job.Complete = true			//bad information, don't try to launch again
		return
	}

	/* Determine node to launch on */
	node, err := DetermineBestNodeForLaunch()
	if err != nil {
		logging.Log(err)
		job.Complete = false 		//bad node, so try to launch in the future
		job.InUse = false
		return 
	}

	/* Get the application information */
	app, err := database.GetApplication(newcontainer.ApplicationID)
	if err != nil {
		logging.Log(err)
		job.Complete = true			//Application doesn't exist, don't try to launch in the future
		return
	}

	//Launch the container on the given node

	err = LaunchAppOnNode(app, node, &newcontainer)
	node.ContainerCount = node.ContainerCount + 1

	if err != nil {
		logging.Log(err)
		job.Complete = true		//Something went wrong, container name conflicts happen here
		return
	}

	_, err = database.UpdateContainer(&newcontainer)

	if err != nil {
		database.CreateContainer(&newcontainer)
	}

	database.UpdateNode(node)

	//save container information to database.
	job.Complete = true
	return

}

func LaunchAppOnNode(app *models.Application, node *models.Node, cont *models.Container) (error) {

	client,err := docker.NewClient(node.DIPAddr + ":" + node.DPort)

	if err != nil {
		logging.Log(err)
	}

	ports := app.GetPorts()

	port := ports[0] +"/tcp"
	exposedPort := map[docker.Port]struct{}{
        docker.Port(port) : {}}
	portbindings:= map[docker.Port][]docker.PortBinding{
        docker.Port(port): {}}

     //try to create container
	containeropts := docker.CreateContainerOptions {
		Name: cont.Name,
		Config: &docker.Config {
				ExposedPorts: exposedPort,
				Image: app.DockerImage,

			},
		HostConfig: &docker.HostConfig {
			PublishAllPorts: true,
			PortBindings: portbindings,
			Privileged: false,			

		},
	}

	cont.Status = "CREATE"

	dock_cont, err := client.CreateContainer(containeropts)
	
	if err != nil {
		logging.Log(err)

		err = client.StartContainer(cont.Name, nil)
	    if err != nil {
	        logging.Log(err)
	        return err
	    }

     	cont.ContainerID = dock_cont.ID
	    cont.NodeID = node.Id
	    cont.ApplicationID = app.Id
	    cont.IsEnabled = true
	    cont.Status = "LAUNCHED"
	    
	    return nil		
	}
		//pull if image not found
		//try to create again

	//start container
	err = client.StartContainer(dock_cont.ID, nil)
    if err != nil {
        logging.Log(err)
        return err
    }

    cont.ContainerID = dock_cont.ID
    cont.NodeID = node.Id
    cont.ApplicationID = app.Id
    cont.IsEnabled = true
    cont.Status = "LAUNCHED"

	
	return nil
}