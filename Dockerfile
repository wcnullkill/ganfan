FROM golang:1.15

RUN git config --global url."https://708ab7be0ab66972ff3af60668f83a1673db2ae2@github.com".insteadOf "https://github.com" \
    && go env -w GOPROXY="https://goproxy.cn,direct" \
    && cd /go/src/ \
    && git clone https://github.com/wcnullkill/ganfan.git \
    && cd ganfan \
    && go mod tidy \