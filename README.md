* appを動かすとき
```sh
$ docker compose build
$ docker compose up
```

* containerでtestするとき
```sh
docker compose -f docker-compose-test.yml build
docker compose -f docker-compose-test.yml up
```

* localのcircleciでtestするとき
```sh
circleci local execute --job blog_test
```

* 認証
```sh
curl -s -X POST -d 'username=[user name]' -d 'password=[password]' 'localhost:1323/authenticate'
```

* entry一覧取得
```sh
curl localhost:1323/restricted/entries_list -H "Authorization: Bearer [token]"
```

* entryをID指定で取得
```sh
curl localhost:1323/restricted/entries/[entry id] -H "Authorization: Bearer [token]"
```
