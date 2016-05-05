# mml

## Docs
* [4C](https://drive.google.com/open?id=1a3MfEczAGnz4AfXzFKyq_aNJR3IqElML5fQp8-930Nk)

## Schemas
1. [Context](https://drive.google.com/open?id=0B6MswmSTZunJVXByMTN4Zm0tRk0)
2. [Containers](https://drive.google.com/open?id=0B6MswmSTZunJZk5wVldNNl96X3M)

## Run app
```bash
./run (arm|amd64) [mailbox_domain] [nginx_port] [postfix_port] [mbm_port]
```
Example
```bash
./run amd64 mymail.local 80 25 8080
```
## Testing
```bash
# getting agouti
go get github.com/sclevine/agouti
go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega

# installing chromedriver
wget http://chromedriver.storage.googleapis.com/2.21/chromedriver_linux64.zip
unzip chromedriver_linux64.zip
sudo mv -f chromedriver /usr/local/share/chromedriver
sudo ln -s /usr/local/share/chromedriver /usr/local/bin/chromedriver
sudo ln -s /usr/local/share/chromedriver /usr/bin/chromedriver

# running tests
go test mbm
go test mailbox
# running only one test (in general "-run" flag equals regexp)
go test -run=TestReadMultiPartMail04 mailbox
```
## Developing front
```bash
docker pull nginx
docker run --name nginx-for-front -d -p 8090:80 -v /home/you/workspace/mml/front/public:/usr/share/nginx/html:ro nginx
cp dev/nginx-for-front.conf.dist dev/nginx-for-front.conf
# put your <mbm-container-IP-172.17.0.2> in dev/nginx-for-front.conf
docker cp dev/nginx-for-front.conf nginx-for-front:/etc/nginx/conf.d/default.conf
docker exec -t nginx-for-front service nginx reload
```
