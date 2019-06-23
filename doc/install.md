# 安裝流程
依賴工具：
- Docker
- docker-compose
- Git

先下載 repo
```
$ git clone https://github.com/lancatlin/DeviceBooking.git
```

編輯 .env，設定環境變數  

```
DB_NAME=
DB_USER=
DB_PASSWORD=
PORT=
```

將 key.pem, cert.pem 檔，放到目錄中，如果沒有程式會自動生成

```
```

啟動程式，如果使用者不在 Docker group 中，需要管理員權限  

```
$ docker-compose up -d --build
```

MariaDB 的檔案會存在 docker volume devicebooking_mariadb-data 中