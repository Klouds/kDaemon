# Klouds Cluster Daemon
# Version 0.0 
[![wercker status](https://app.wercker.com/status/7a1a06d652cb003d898554754a8c3c3d/s/master "wercker status")](https://app.wercker.com/project/bykey/7a1a06d652cb003d898554754a8c3c3d)
[![Floobits Status](https://floobits.com/ozzadar/kDaemon.svg)](https://floobits.com/ozzadar/kDaemon/redirect)


<img src="http://www.ozzadar.com/klouds.png" align="center"/>
____

*A RESTful API framework for managing containers over a set of Docker Nodes.*

*Will also run a watcher process that will monitor, fix and rebalance containers on the cluster.*

## Quickstart Guide
[![ScreenShot](https://j.gifs.com/yPprOE.gif)](https://www.youtube.com/watch?v=YlxeknZDP9Y)
Click here for quickstart guide



# Features
* Easily mount new docker nodes through either REST calls or the Web UI
* Create a persistant application library of Docker applications to launch at will
* Balances containers across the cluster in a containers-per-node manner (at present)
* Maintains container uptime by routinely health-checking the cluster

## WEB UI

Web UI for interacting with kDaemon:
[kDaemon_UI][kDaemon_UI]





Full stack design can be read here 
     [Development Plan][Development Plan] and here
     [Daemon Design] [Daemon Design]


## THE STACK
*The stack consists of multiple docker endpoints exposing their APIs to their internal network (or VPN) and managing the state of the cluster from a central location.*

```
    docker --\    
                
    docker -->kDaemon 
                                     
    docker __/   

```
## Dependencies

+ Your nodes need to be running the docker API 
```
sudo docker -H 0.0.0.0:2375 -H unix:///var/run/docker.sock -d 
```

+ Rethinkdb with database name of rethinkdb_dbname and tables
    - containers
    - nodes
    - applications
    
+ You should have golang 1.5+ available on your machine.
+ kDaemon is only tested on Linux. Any machine with an exposed Docker API should be able to become a node however.


# Running the Application
    
## Configuration -- config/app.conf

```
    [default]
    bind_ip = 0.0.0.0               # IP to bind API to
    api_port = 1337                 # Port to bind API to
    rethinkdb_host = 127.0.0.1      # RethinkDB Host
    rethinkdb_port = 28015          # RethinkDB Port
    rethinkdb_dbname = kdaemon      # RethinkDB database name
    api_version = 0.0               # Current API version. (You probably shouldnt change this)

```




## To build (linux):


```
    go get github.com/superordinate/kDaemon
    cd $GOPATH/src/github.com/superordinate/kDaemon
    go build .
```

## To Run

``` 
./kDaemon
```
## How to use
*You can interact with kDaemon in several ways. The easiest way to get started is to use our web-ui at:*
[http://github.com/klouds/kdaemon_ui.][kDaemon_UI] 

If you're interested in working directly with the REST APIs, you here's some documentation for you:
   
#### Adding a node

To add a node, you can POST to **bind_ip:api_port/%API_VERSION%/nodes/create**
```
    {
        "name":"Host1",
        "d_ipaddr":"127.0.0.1",
        "d_port": "2375",
    }
```

#### Adding an application

To add an application, you can POST to **bind_ip:api_port/%API_VERSION%/applications/create**
```
    {
        "name":"ghost-blog",
        "exposed_ports":"2368",
        "docker_image": "ghost"
    }
```

#### Create a container

To create a container on your most available node, you can POST to **bind_ip:api_port/%API_VERSION%/containers/create**

*note: The application_id is generated when adding a new application.*
```
    {
        "name":"ghost-blog-ozzadar",
        "application_id":75186c59-ec80-49d6-beb5-1bfac76e8525

    }
```

#### Launch a container

**POST** to **bind_ip:api_port/%API_VERSION%/container/launch/%CONTAINER_ID%**

## API Reference
UPDATED (Feb 25 2016) -- 
    1) adding a start container command to the list. Separating container creation/launch from each other because what was I thinking?
    2) moved endpoints to /action/id to get rid of annoying pathing errors.


```
 * POST /%API_VERSION%/nodes/create             -- Creates a node
 * PATCH /%API_VERSION%/nodes/update/id         -- Edits a node
 * DELETE /%API_VERSION%/nodes/delete/id        -- Deletes a node
 * GET /%API_VERSION%/nodes/id                  -- Gets node information
 * GET /%API_VERSION%/nodes                     -- Gets all nodes
 
 * POST /%API_VERSION%/applications/create      -- Creates an application in the database
 * PATCH /%API_VERSION%/applications/update/id  -- Edits an application
 * DELETE /%API_VERSION%/applications/delete/id -- Deletes an application
 * GET /%API_VERSION%/applications/id           -- Gets application information
 * GET /%API_VERSION%/applications              -- Gets all applications

 * POST /%API_VERSION%/containers/create        -- Creates a container on the cluster
 * POST /%API_VERSION%/containers/launch/id     -- Launches a container on the cluster
 * POST /%API_VERSION%/containers/launch/id     -- Stops a container on the cluster
 * PATCH /%API_VERSION%/containers/update/id    -- Edits a container 
 * DELETE /%API_VERSION%/containers/delete/id   -- Deletes a container 
 * GET /%API_VERSION%/containers/id             -- Gets container information
 * GET /%API_VERSION%/containers                -- Gets all containers

 ```

#### Questions on how to use? Contact [Ozzadar](https://github.com/Ozzadar) for more info =)

[Development Plan]: https://docs.google.com/document/d/1A4-0g1E52wdW9L-hoeAZzay5Uotv1GcBPtXLU1msw2w/edit?usp=sharing
[Daemon Design]: https://docs.google.com/document/d/1EkI7uQzdt1xMwb1etcweYQFCLthK_l9aHZvHOunshzs/edit?usp=sharing
[Weave]: http://www.weave.works/
[kDaemon_UI]:http://github.com/klouds/kDaemon_ui
