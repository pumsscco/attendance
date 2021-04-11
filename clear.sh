#!/bin/bash
docker-compose rm -sf
docker rmi -f $(docker images attendance_attendance -aq)
