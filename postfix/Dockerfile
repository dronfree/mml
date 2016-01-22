FROM ubuntu:14.04.3
RUN apt-get update && apt-get install -y postfix
ADD main.cf /etc/postfix/main.cf
RUN ["/bin/bash", "-c", "mkdir -p /var/mail/vhosts"]
RUN ["/bin/bash", "-c", "chmod 0777 /var/mail/vhosts"]
ADD vmailbox /etc/postfix/vmailbox
RUN ["/bin/bash", "-c", "postmap /etc/postfix/vmailbox"]
EXPOSE 25:25