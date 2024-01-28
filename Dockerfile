FROM ubuntu:20.04

RUN apt-get update && \
apt-get install -y --no-install-recommends openssh-server vim less sudo

RUN mkdir /var/run/sshd

RUN echo 'root:root' | chpasswd

RUN useradd -m -d /home/casone -s /bin/bash casone

RUN echo 'casone:casone' | chpasswd

RUN sed -i 's/#PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config

RUN sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin no/' /etc/ssh/sshd_config

RUN sed -i 's/#Port 22/Port 20022/' /etc/ssh/sshd_config

COPY ./.ssh/go-scp.pub /home/casone/.ssh/authorized_keys

EXPOSE 20022
CMD ["/usr/sbin/sshd", "-D"]
