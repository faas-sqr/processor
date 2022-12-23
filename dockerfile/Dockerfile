FROM ubuntu:22.04
ADD  ./python-ubuntu.tar.gz /home/function

WORKDIR /home/function
RUN /bin/sh -c 'cp -r /home/function/python-ubuntu  /home/function/functionA && mv /home/function/python-ubuntu /home/function/functionB'

COPY ./bin  /usr/local/bin/
CMD ["processor"]
#CMD ["sleep", "100d"]