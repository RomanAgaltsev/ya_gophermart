FROM scratch
WORKDIR /app
COPY cmd/accrual/accrual_linux_amd64 /accrual
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/accrual"]