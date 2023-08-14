FROM ubuntu:22.04

COPY ./start_server /home
COPY ./userManage/user_files /home/userManage/user_files

EXPOSE 8888 3301 3302 3303 3304 8080

