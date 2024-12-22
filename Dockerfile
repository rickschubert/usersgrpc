FROM golang:1.23.1
# Install aws-cli
RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
RUN apt-get update
RUN apt install unzip
RUN unzip awscliv2.zip
RUN ./aws/install
# Build and launch app
WORKDIR /app
COPY go.mod go.sum ./
COPY config config
COPY server server
COPY users users
COPY db db
COPY cmd cmd
RUN go build -o serverexecutable ./cmd
CMD "./serverexecutable"
