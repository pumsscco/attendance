FROM scratch
MAINTAINER Pluto_Chan
ADD docker.yml /
ADD attendance-docker /
EXPOSE 5578
CMD ["/attendance-docker"]
