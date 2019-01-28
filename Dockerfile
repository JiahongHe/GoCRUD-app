FROM golang:onbuild
EXPOSE 8080
CMD ["go", "run", "main.go", "handlers.go", "types.go", "utils.go"]


