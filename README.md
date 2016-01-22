# mml
Mail for an hour
```bash
cd postfix && docker build -t postfix .
```
Run container
```bash
docker run -d -p 8025:25 -t postfix
```