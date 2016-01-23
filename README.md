# mml
Build docker image
```bash
cd postfix && docker build -t postfix .
```
Run container
```bash
docker run -d -p 25:25 -t postfix
```