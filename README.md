# Klouds Cluster Daemon
# Version 0.0 
[![wercker status](https://app.wercker.com/status/7a1a06d652cb003d898554754a8c3c3d/s/master "wercker status")](https://app.wercker.com/project/bykey/7a1a06d652cb003d898554754a8c3c3d)
[![Floobits Status](https://floobits.com/ozzadar/Pauls_Dojo.svg)](https://floobits.com/ozzadar/Pauls_Dojo/redirect)


<img src="http://www.ozzadar.com/klouds.png" align="center"/>


A RESTful API framework for managing containers over a set of Docker Nodes running Prometheus as a monitoring solution.
Will also run a watcher process that will monitor, fix and/or notify administration on the state of the cluster.

###Features
	* Make daemon aware of new docker/prometheus endpoints (adding nodes to the cluster)
	* Launch containers on the cluster (use prometheus to loadbalance and schedule jobs)
	* Forward container access information to Consul for loadbalancing endpoints through a reverse proxy (such as HAProxy or Nginx)


Full stack design can be read here 
	 [Development Plan][Development Plan] and here
	 [Daemon Design] [Daemon Design]


##THE STACK

```
	docker-prometheus --\    /---> consul ---> haproxy-consul  <--\
				   		 \	/									   \
	docker-prometheus --->kDaemon -><-- klouds-frontend <----------User
				   		 /					 
	docker-prometheus __/	 

```


##HOW TO USE:

###Configuration -- config/app.conf

```
	[default]
	bind_ip = 0.0.0.0   			# IP to bind API to
	bind_port = 1337				# Port to bind API to
	mysql_host = localhost 			# Address to mysql server
	mysql_port = 3306				# port for mysql server
	mysql_user = root				# mysql username
	mysql_password = thesecretsauce	# mysql password
	mysql_dbname = kDaemon			# database name

```

###Dependencies

	* Your nodes need to be running the docker API
	* You must have a mysql host

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


### ENDPOINTS [X] = Not implemented
```
 * POST /%API_VERSION%/node/create  			-- Creates a node
 * PUT /%API_VERSION%/nodes/update 			 	-- Edits a node
 * DELETE /%API_VERSION%/nodes/delete  			-- Deletes a node
 * GET /%API_VERSION%/nodes/:id  				-- Gets node information
 * GET /%API_VERSION%/nodes  					-- Gets all nodes
 
 * POST /%API_VERSION%/application/create  		-- Creates an application in the database
 * PUT /%API_VERSION%/applications/update  		-- Edits an application
 * DELETE /%API_VERSION%/applications/delete  	-- Deletes an application
 * GET /%API_VERSION%/applications/:id  		-- Gets application information
 * GET /%API_VERSION%/applications  			-- Gets all applications

 * POST /%API_VERSION%/container/create  		-- Creates a container on the cluster
 * PUT /%API_VERSION%/containers/update  		-- Edits a container [X]
 * DELETE /%API_VERSION%/containers/delete  	-- Deletes a container [X]
 * GET /%API_VERSION%/containers/:id  			-- Gets container information
 * GET /%API_VERSION%/containers  				-- Gets all containers

 * POST /%API_VERSION%/user/create  			-- Creates a user [X]
 * PUT /%API_VERSION%/users/update  			-- Edits a users [X]
 * DELETE /%API_VERSION%/users/delete  			-- Deletes a user [X]
 * GET /%API_VERSION%/users/:id  				-- Gets user information [X]
 ```


#### Main Contributor [Ozzadar](https://github.com/Ozzadar)
[Development Plan]: https://docs.google.com/document/d/1A4-0g1E52wdW9L-hoeAZzay5Uotv1GcBPtXLU1msw2w/edit?usp=sharing
[Daemon Design]: https://docs.google.com/document/d/1EkI7uQzdt1xMwb1etcweYQFCLthK_l9aHZvHOunshzs/edit?usp=sharing