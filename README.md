
Docker commands:

- docker build --pull --rm -f "DockerFile" -t gopostgretask:latest "."
- docker run -p 3000:3001 -tid gopostgretask:latest