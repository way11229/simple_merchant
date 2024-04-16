# Simple Merchant Service

該服務提供使用email註冊、驗證email、登入、登出、取得推薦商品列表等功能。

## API Document

[API 文件](https://gitlab.com/way11229/simple_merchant/-/blob/main/doc/service.swagger.json?ref_type=heads)

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

### Structured Project Layout

![image](https://github.com/way11229/simple_merchant/blob/main/simple_merchant_struct_project_layout.png)

## Instructions

本專案根目錄有提供docker-compose.yml，可於根目錄直接啟動服務，並依照需求，修改docker-compose.yml中的參數即可。

## Tests

本專案於tests/acceptance_tests有實做驗收測試，可依照需求，修改tests/acceptance_tests/test.env中的參數即可。

因本專案有權限設計，故部份驗收測試需提供access token來執行，請將環境中測試帳號的access token填入test.env中的ACCESS_TOKEN。
