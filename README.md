#### Runnel
Service to stream live output of command line interface over http.


##### RUN SERVER (docker - recommended)

- Clone project and run server inside docker container

        git clone https://github.com/VeritasOS/runnel.git
        cd runnel

        # build docker image
        docker build -t runnel:latest .

        # start runnel server
        docker run --name runnel -d -p 127.0.0.1:9090:9090 runnel:latest


##### USAGE

- Trigger below commands to get live stream

        # Fire command with curl
        curl -X POST -H 'Content-Type: application/json' -d '{"cmd":"ping -c 3 google.com"}' http://localhost:9090/command

        # Get live output stream
        # replace uuid with one which you get from above command
        while true; do curl -sS -H 'Content-Type: application/json' 'http://localhost:9090/stream/fd4b1a38-94f4-4eba-80e7-50578ac4baae' | jq '.response'; done

- Options

        # timeout (seconds) - wait time for your command to start logging
        http://localhost:9090/stream/fd4b1a38-94f4-4eba-80e7-50578ac4baae?timeout=30


##### RUN SERVER (on host - not recommended)

- Clone project and run server directly on host

        git clone https://github.com/VeritasOS/runnel.git
        cd runnel

        ./bin/linux_64/runnel_server -p localhost:9090


##### LIBRARY

- If you are working with golang you can use the api directly
- Install redis server

```
        # Get lib
	    go get github.com/veritasos/runnel/runnel

        # Fire command
        client := runnel.NewClient()
        key, err := client.RunCommand("ping", "-c 2 google.com")

        # Get output stream
        client := runnel.NewClient()
        output, err := client.Stream(key, 10)
```


##### NOTES

- Run runnel service locally and wrap it with your rest api service or cli
- If you are running runnel server on host, make sure to create a new user with strict permissions.
- It's recommended to run runnel server with Dockerfile


##### WARNING

- Do not run the runnel server with static ip address, its recommended to use `127.0.0.1` or `localhost`
- Do not run the runnel server with `root` user


##### TODO
- Integrate websocket
