FROM golang:1.18.1 AS builder

#Установим необходимые переменные среды, необходимые для нашего изображения
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

#Создаем и переходим в рабочую директорию /build
WORKDIR /build

#Копируем и загружаем зависимости go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

#Копируем полностью весь код в контейнер
COPY . .

#Создаем бинарный файл и переносим его в контейнер
RUN go build -o main .

#Создаем новую рабочую директорию /dist
WORKDIR /dist

#Копируем бинарный файл в новую рабочую директорию /dist
RUN cp /build/main .

#Уменьшаем размер образа
FROM scratch

COPY --from=builder /dist/main /

#Запускаем бинарный файл
#ENTRYPOINT ["/main"]
CMD ["/main"]