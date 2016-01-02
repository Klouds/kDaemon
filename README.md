# Klouds Cluster Daemon
# Version 0.0 

<img src="http://www.ozzadar.com/klouds.png" align="center"/>


A RESTful API framework for managing containers over a set of Docker Nodes running Prometheus as a monitoring solution.

###Features
	* Make daemon aware of new docker/prometheus endpoints (adding nodes to the cluster)
	* Launch containers on the cluster (use prometheus to loadbalance and schedule jobs)
	* Forward container access information to Consul for loadbalancing endpoints through a reverse proxy (such as HAProxy or Nginx)

Full stack design can be read here:
	* https://docs.google.com/document/d/1A4-0g1E52wdW9L-hoeAZzay5Uotv1GcBPtXLU1msw2w/edit?usp=sharing
	* https://docs.google.com/document/d/1EkI7uQzdt1xMwb1etcweYQFCLthK_l9aHZvHOunshzs/edit?usp=sharing

##THE STACK

```

docker-prometheus --\   /---> consul ---> haproxy-consul  <---\
			   		  \	/									   \
docker-prometheus ---->kDaemon -><-- klouds-frontend <-------User
			   		  /					 
docker-prometheus __/	 

```


##HOW TO USE:

###Environment Variables

MYSQL_HOST= 127.0.0.1:3306	 			<-- Points to your database
MYSQL_USER= root						<-- User for your mysql database
MYSQL_PASSWORD= iamapassword			<-- password for you mysql user


###To build (linux):


```
go get github.com/superordinate/kDaemon
cd $GOPATH/src/github.com/superordinate/kDaemon
go build .

```
### To Run

``` 

	MYSQL_HOST=127.0.0.1:3306 MYSQL_USER=root MYSQL_PASSWORD=iamapassword ./klouds

```

OR

```
	export MYSQL_HOST= 127.0.0.1:3306
	export MYSQL_USER= root	
	export MYSQL_PASSWORD= iamapassword

	./kDaemon

```

### ENDPOINTS
 * POST /%API_VERSION%/node/create  -- Creates a node
 * PUT /%API_VERSION%/nodes/update  -- Edits a node
 * DELETE /%API_VERSION%/nodes/delete  -- Deletes a node
 * GET /%API_VERSION%/nodes/:id  -- Gets node information
 
 * POST /%API_VERSION%/application/create  -- Creates an application in the database
 * PUT /%API_VERSION%/applications/update  -- Edits an application
 * DELETE /%API_VERSION%/applications/delete  -- Deletes an application
 * GET /%API_VERSION%/applications/:id  -- Gets application information

 * POST /%API_VERSION%/container/create  -- Creates a container on the cluster
 * PUT /%API_VERSION%/containers/update  -- Edits a container
 * DELETE /%API_VERSION%/containers/delete  -- Deletes a container
 * GET /%API_VERSION%/containers/:id  -- Gets container information

 * POST /%API_VERSION%/user/create  -- Creates a user
 * PUT /%API_VERSION%/users/update  -- Edits a user
 * DELETE /%API_VERSION%/users/delete  -- Deletes a user
 * GET /%API_VERSION%/users/:id  -- Gets user information

