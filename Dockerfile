FROM golang:1.15

RUN git config --global url."https://77a0f93cae0ccf8537faca7a81667a6d648a737e@github.com".insteadOf "https://github.com" \
    && git config --global http.proxy http://10.20.16.124:7890 \
    && go env -w GOPROXY="https://goproxy.cn,direct" \
    && cd /go/src/ \
    && git clone -b dev https://github.com/wcnullkill/ganfan.git \
    && cd ganfan \
    && go mod tidy \