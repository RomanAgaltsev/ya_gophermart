FROM alpine
WORKDIR /
COPY ../cmd/accrual/accrual_linux_amd64 /accrual_linux_amd64
EXPOSE 8080
ENTRYPOINT ["./accrual_linux_amd64"]