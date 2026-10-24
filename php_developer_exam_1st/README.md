# php_developer_exam_1st

> 本目录原始 README 见文末，本节为快照补充信息。

原始路径：`php_developer_exam_1st`　|　文件 **221**　|　主语言：PHP 177, Shell 27, JavaScript 5

## 目录构成

```
laravel/resources/  — 77
laravel/app/  — 58
laravel/tests/  — 16
docker/rootfs/  — 14
docker/prebuildfs/  — 14
laravel/config/  — 13
laravel/database/  — 11
laravel/  — 8
laravel/routes/  — 3
laravel/bootstrap/  — 2
laravel/public/  — 2
laravel/vendor/  — 1
```

## 关键文件

- 前端入口：`laravel/resources/js/app.js`
- CI/构建：`docker/Dockerfile`

---

## 原始说明

# PHP exam1

## System requirements
* PHP8.4
* Laravel 12
* MySQL8.0+
* docker

## Installation Steps

* clone project and cd
* docker compose up -d
* check cache folders permission
* composer install
* npm install && npm run build
* setup .env
* php artisan key:generate
* php artisan config:clear
* php artisan migrate

## DB migrations
* user table username unique, add balance
* add account_transaction for deposit / withdraw
    * transaction_no for our system, TN+date(YmdHis)+last_insert_id(補0, 4位)
    * type 0 => deposit , type 1 -> withdraw
    * status 0 => pending , 1 ->succeed , 2-> failed
* add game_records table to store user bet records


## crontab
同步 vendor 資料 ，每分鐘跑 (撈兩分內的資料，避免遊戲商延遲問題，可加長時間, maybe 10 mins)
```
# 注意事項: 起始結束時間不得大於八小時, 每分鐘只能跑一次，相同參數不要重複跑，會被擋
* * * * * /usr/bin/docker exec -it bitnami bash -c "php artisan game-records:fetch --vendor=target --start_time=\"$(date -d '2 minute ago' +\%Y-\%m-\%d \%H:\%I:\%S)\" --end_time=$(date +\%Y-\%m-\%d \%H:\%I:\%S)\" "
```

