# mml

## Docs
* [4C](https://drive.google.com/open?id=1a3MfEczAGnz4AfXzFKyq_aNJR3IqElML5fQp8-930Nk)

## Schemas
1. [Context](https://drive.google.com/open?id=0B6MswmSTZunJVXByMTN4Zm0tRk0)
2. [Containers](https://drive.google.com/open?id=0B6MswmSTZunJZk5wVldNNl96X3M)

## Run app
```bash
./run (arm|amd64) mail_box_domains
```
Example
```bash
./run amd64 mymail.local
```
## Testing
```bash
# getting agouti
go get github.com/sclevine/agouti
go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega

#installing chromedriver
wget http://chromedriver.storage.googleapis.com/2.21/chromedriver_linux64.zip
unzip chromedriver_linux64.zip
sudo mv -f chromedriver /usr/local/share/chromedriver
sudo ln -s /usr/local/share/chromedriver /usr/local/bin/chromedriver
sudo ln -s /usr/local/share/chromedriver /usr/bin/chromedriver

#running tests
go test mbm
go test
```