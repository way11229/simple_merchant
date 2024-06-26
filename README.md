# Simple Merchant Service

該服務提供使用email註冊、驗證email、登入、登出、取得推薦商品列表等功能。

## API Document

[API 文件](https://gitlab.com/way11229/simple_merchant/-/blob/main/doc/swagger/service.swagger.json)

亦可參考./doc/swagger/service.swagger.json

Error List

2: unknown

3: invalid parameters

5: not found

6: already exists

7: permission denied

10: aborted

13: internal server error

|  code   | error message  |
|  ----  | ----  |
| 2  | unknown error |
| 3  | missing required parameters |
| 3  | invalid user name |
| 3  | invalid email |
| 3  | invalid user password |
| 3  | email has verified |
| 3  | invalid product name |
| 5  | record not found |
| 5  | email not found |
| 6  | user email duplicated |
| 7  | permission deny |
| 10  | login aborted |
| 10  | invalid verification code |
| 10  | verification code expired |
| 10  | send verification code in short period |
| 13  | internal server error |

## System Design Documentation

1. 該服務使用go 1.22.2
2. 程式架構參考 [go-clean-arch](https://github.com/bxcodec/go-clean-arch)
3. 實現grpc接口，監控9000 port; 並使用[grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)實現http接口及Restful API，監控8080 port
4. 使用[go-migrate](https://github.com/golang-migrate/migrate)控制database schema版本
5. 使用[sqlc](https://github.com/sqlc-dev/sqlc)實現database操作
6. 使用[go-redis](https://github.com/redis/go-redis)實現redis操作
7. 使用alpine3.19作為執行容器

為提昇推薦商品列表API效能，使用redis sorted sets依照設定商品權重儲存商品ID，並將推薦商品資訊以json格式存於redis strings中，作為快取。

### Structured Project Layout

![image](https://github.com/way11229/simple_merchant/blob/main/simple_merchant_struct_project_layout.png)

## Instructions

本專案根目錄有提供docker-compose.yml，可於根目錄直接啟動服務，並依照需求，修改docker-compose.yml中的參數即可。

因本專案未實際實現驗證email寄送，故將驗證碼直接使用log顯示。

參數說明

|  參數   | 說明  |
|  ----  | ----  |
|  MYSQL_SQL_DRIVER_NAME  |  sql mysql driver name  |
|  MYSQL_SQL_DATA_SOURCE_NAME  |  sql mysql 連線資訊  |
|  MYSQL_MIGRATION_SOURCE_URL  |  go migrate mysql schema 檔案路徑  |
|  MYSQL_MIGRATION_DATABASE_URL  |  go migrate mysql 連線資訊 |
|  LOGIN_TOKEN_EXPIRE_SECONDS  |  登入token過期時間（秒）  |
|  LOGIN_TOKEN_CACHE_EXPIRE_SECONDS  |  登入token快取過期時間（秒）  |
|  USER_EMAIL_VERIFICATION_CODE_LEN  |  使用者email驗證碼長度  |
|  USER_EMAIL_VERIFICATION_CODE_MAX_TRY  |  使用者email驗證最多嘗試次數  |
|  USER_EMAIL_VERIFICATION_CODE_EXPIRED_SECONDS  |  使用者email驗證碼過期時間（秒）  |
|  USER_EMAIL_VERIFICATION_CODE_ISSUE_LIMIT_SECONDS  |  發送使用者email驗證碼間隔時間（秒）  |
|  VERIFICATION_EMAIL_SUBJECT  |  使用者email驗證信標題  |
|  VERIFICATION_EMAIL_CONTENT  |  使用者email驗證信內容  |
|  SYMMETRIC_KEY  |  access token 金鑰  |
|  REDIS_ADDR  |  redis 連線資訊  |
|  REDIS_PWD  |  redis 密碼  |
|  RECOMMENDED_PRODUCT_CACHE_EXPIRED_SECONDS  |  推薦商品快取過期時間  |

## Tests

本專案於tests/acceptance_tests實做驗收測試，可依照需求，修改tests/acceptance_tests/test.env中的參數即可。

因本專案有權限設計，故部份驗收測試需提供access token來執行，請將環境中測試帳號的access token填入test.env中的ACCESS_TOKEN。

### Stress test

因本專案取得推薦商品API預計每分鐘會有300次取用，故使用jmeter驗證其效能。

CPU 2.90GHz × 16

RAM 48G

每分鐘300次請求
![image](https://github.com/way11229/simple_merchant/blob/main/stress_test_300_60.png)

每分鐘1000次請求
![image](https://github.com/way11229/simple_merchant/blob/main/stress_test_1000_60.png)
