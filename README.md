# Eribo - (エリボ)
F-list bot.

Launching a testbot:

```
go install -race && eribo \
-account=kusubooru \
-password='<password>' \
-character=testbot2 \
-datasource='kusubooru:kusubooru@()/eribo?parseTime=true' \
-join='["lab"]'
```
