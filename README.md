# Klouds Cluster Daemon
# Version 0.0 
[![wercker status](https://app.wercker.com/status/7a1a06d652cb003d898554754a8c3c3d/s/master "wercker status")](https://app.wercker.com/project/bykey/7a1a06d652cb003d898554754a8c3c3d)
[![Floobits Status](https://floobits.com/ozzadar/kDaemon.svg)](https://floobits.com/ozzadar/kDaemon/redirect)


<img src="http://www.ozzadar.com/klouds.png" align="center"/>


A RESTful API framework for managing containers over a set of Docker Nodes running Prometheus as a monitoring solution.
Will also run a watcher process that will monitor, fix and/or notify administration on the state of the cluster.


##WEB UI -- MOVED to new REPO

The user interface currently included has no functionality and is rather ugly. New UI is being built at

[kDaemon_UI][kDaemon_UI]

with the websocket being served at /ws

###Features
	* Make daemon aware of new docker/prometheus endpoints (adding nodes to the cluster)
	* Launch containers on the cluster (use prometheus to loadbalance and schedule jobs)
	* Forward container access information to Consul for loadbalancing endpoints through a reverse proxy (such as HAProxy or Nginx)


Full stack design can be read here 
	 [Development Plan][Development Plan] and here
	 [Daemon Design] [Daemon Design]


###THE STACK

```
	docker-prometheus --\    /---> consul ---> haproxy-consul  <--\
				   		 \	/									   \
	docker-prometheus --->kDaemon -><-- klouds-frontend <----------User
				   		 /					 
	docker-prometheus __/	 

```

###Configuration -- config/app.conf


```
	[default]
	bind_ip = 0.0.0.0   			# IP to bind API to
	api_port = 1337				    # Port to bind API to
	rethinkdb_host = 127.0.0.1		# RethinkDB Host
	rethinkdb_port = 28015			# RethinkDB Port
	rethinkdb_dbname = kdaemon      # RethinkDB database name

```

###Dependencies

+ Your nodes need to be running the docker API
+ Rethinkdb with database name of rethinkdb_dbname and tables
	- containers
	- nodes
	- applications

	(will be creating tables if they don't exist in the future, not implemented at the moment though)


###To build (linux):


```
	go get github.com/superordinate/kDaemon
	cd $GOPATH/src/github.com/superordinate/kDaemon
	go build .

```
### To Run

``` 
	./kDaemon

```

##HOW TO USE:

For development, it is okay to publicly expose your docker endpoints with
```
	docker -H 0.0.0.0:2375 -H unix:///var/run/docker.sock -d &
```

If attempting to deploy in production, you will want to point your docker hosts into a VPN and have kDaemon linked through this VPN. We suggest using [Weave][Weave]. Support for Docker authentication is in the pipeline though you should use VPNs regardless.

####Adding a node

To add a node, you can POST to **bind_ip:api_port/%API_VERSION%/nodes/create**
```
    {
        "hostname":"Host1",
        "d_ipaddr":"127.0.0.1",
        "d_port": "2375",
    }
```

####Adding an application

To add an application, you can POST to **bind_ip:api_port/%API_VERSION%/application/create**
```
    {
        "name":"ghost-blog",
        "exposed_ports":"2368",
        "docker_image": "ghost",
        "dependencies":"",
        "isenabled":true
    }
```

####Create a container

To create a container on your most available node, you can POST to **bind_ip:api_port/%API_VERSION%/container/create**
```
    {
        "name":"ghost-blog-ozzadar",
        "application_id":1
    }
```

####Launch a container

**POST** to **bind_ip:api_port/%API_VERSION%/container/launch/%CONTAINER_ID%**

## API Reference
UPDATED (Feb 25 2016) -- 
    1) adding a start container command to the list. Separating container creation/launch from each other because what was I thinking?
    2) moved endpoints to /action/id to get rid of annoying pathing errors.


```
 * POST /%API_VERSION%/nodes/create  			-- Creates a node
 * PATCH /%API_VERSION%/nodes/update/id 		-- Edits a node
 * DELETE /%API_VERSION%/nodes/delete/id  		-- Deletes a node
 * GET /%API_VERSION%/nodes/id  				-- Gets node information
 * GET /%API_VERSION%/nodes  					-- Gets all nodes
 
 * POST /%API_VERSION%/applications/create  	-- Creates an application in the database
 * PATCH /%API_VERSION%/applications/update/id  -- Edits an application
 * DELETE /%API_VERSION%/applications/delete/id -- Deletes an application
 * GET /%API_VERSION%/applications/id  		-- Gets application information
 * GET /%API_VERSION%/applications  			-- Gets all applications

 * POST /%API_VERSION%/containers/create        -- Creates a container on the cluster
 * POST /%API_VERSION%/containers/launch/id  	-- Launches a container on the cluster
 * PATCH /%API_VERSION%/containers/update/id  	-- Edits a container 
 * DELETE /%API_VERSION%/containers/delete/id  	-- Deletes a container 
 * GET /%API_VERSION%/containers/id  			-- Gets container information
 * GET /%API_VERSION%/containers  				-- Gets all containers

 ```

#### Main Contributor [Ozzadar](https://github.com/Ozzadar)
[Development Plan]: https://docs.google.com/document/d/1A4-0g1E52wdW9L-hoeAZzay5Uotv1GcBPtXLU1msw2w/edit?usp=sharing
[Daemon Design]: https://docs.google.com/document/d/1EkI7uQzdt1xMwb1etcweYQFCLthK_l9aHZvHOunshzs/edit?usp=sharing
[Weave]: http://www.weave.works/
[kDaemon_UI]:http://github.com/klouds/kDaemon_ui
