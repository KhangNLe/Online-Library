# syntax=docker/dockerfile:1

FROM golang:1.23

# Set destination for COPY
WORKDIR /server/go

#Sqlite3
RUN apt-get update && apt-get install -y sqlite3
# Go modules
COPY /server/go/go.mod /server/go/go.sum ./
RUN go mod download


COPY server/library.db /app/library.db
COPY clientSide/ /app/clientSide/
COPY server/go/ ./

# Build
RUN CGO_ENABLED=1 GOOS=linux go build -o /docker-gs-ping


EXPOSE 6969

# Run
CMD ["/docker-gs-ping"]
