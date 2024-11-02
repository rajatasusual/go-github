# Build the docker image

```
$ docker build -t go-github .
```

# Tag the image

```
$ docker tag go-github your_username/main:1.0.0
```

# Login to docker with your docker Id

```
$ docker login
Login with your Docker ID to push and pull images from Docker Hub. If you do not have a Docker ID, head over to https://hub.docker.com to create one.
Username (your_username): your_username
Password:
Login Succeeded
```

# Push the image to docker hub

```
$ docker push your_username/main:1.0.0
````