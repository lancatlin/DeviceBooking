FROM golang:latest
WORKDIR /app
COPY . /app 
RUN echo "Asia/Taipei" > /etc/timezone && dpkg-reconfigure -f noninteractive tzdata
RUN go build -o app 
CMD ./app