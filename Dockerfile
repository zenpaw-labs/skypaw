# just sandbox

FROM fedora:latest
WORKDIR /skypaw
COPY . .
RUN dnf install -y golang && go build ./cmd/skypaw && go clean -cache -modcache -testcache -fuzzcache && dnf remove -y golang && dnf clean all
CMD ["./skypaw"]
