FROM taekion/taekion-fs-client:latest

RUN mkdir /home/user/tfs_mount && chown user:user /home/user/tfs_mount
USER root
RUN echo user_allow_other > /etc/fuse.conf
USER user