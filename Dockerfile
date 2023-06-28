FROM golang:1.20-alpine3.18

# Make a directory for the files
RUN mkdir -p /app

# Copies your code file from your action repository to the filesystem path `/` of the container
COPY entrypoint.sh /app/entrypoint.sh

# Make the entrypoint.sh executable
RUN chmod +x /app/entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]
