FROM ubuntu:kinetic
WORKDIR /app
COPY /bin/resizer /app
EXPOSE 3300
RUN apt update
RUN apt install -y ca-certificates
CMD ["/app/resizer", "-limit","10"]


