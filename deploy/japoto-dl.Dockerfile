FROM golang:1.17 AS base

RUN apt-get update \
		&& apt-get install \
			-y --no-install-recommends \
			software-properties-common ffmpeg \
		&& rm -rf /var/lib/apt/lists/*

RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go build -o japoto-dl

RUN apk add --no-cache tzdata
ENV TZ=Japan

CMD [ "/app/japoto-dl" ]
