# HR_System

---

## Requirements

> Write Go GIN RESTful API Service include below Condition (Push to Github)
> ・HR System Backend API
> ・MySQL as  Database
> ・Redis as Cache
> ・GORM Migration
> ・GORM MySQL SEED data
> ・Unit Test
> ・Makefile for build and deploy
> ・Run all service using docker-compose

## System Architecture

* Application HR API
* PostgreSQL as database
* Redis as cache server
* OpenTelemetry Collector as telemetry data collector
* Telemetry services : Loki, Tempo, Mimir
* Dashboard : Grafana
* Testing tool : Grafana k6

## API Endpoints

** POST /api/v1/auth/login
用途: 用於用戶登入系統，取得身份驗證 Token。

** GET /api/v1/auth/me
用途: 驗證當前登入用戶的資訊，並返回其角色和基本資料。

** GET /api/v1/employees
用途: 獲取所有員工的列表（HR 可看到詳細資料，其他角色可能僅能看到部分資料）。

** GET /api/v1/employees/:id
用途: 根據員工 ID 獲取該員工的詳細資料。（僅 HR 與該員工本人可查看）

** POST /api/v1/employees
用途: 新增一名員工（僅 HR 可操作）。

** PUT /api/v1/employees/:id
用途: 更新指定員工的資料。

** DELETE /api/v1/employees/:id
用途: 軟刪除指定員工資料（僅 HR 可操作）。

** GET /api/v1/attendance/:employ_id
用途: 獲取指定員工的所有出勤記錄。（僅 HR 與該員工本人以及同部門主管可查看）

** POST /api/v1/attendance/:employ_id
用途: 員工打卡（上班或下班）。同一天內的第一次打卡視為上班，第二次打卡視為下班。

** GET /api/v1/attendance/:employ_id/:date
用途: 獲取指定員工在某一天的出勤記錄。

** GET /api/v1/leave/:id
用途: 獲取指定請假單的詳細資料。（僅 HR 與該員工本人以及同部門主管可查看）

** GET /api/v1/leave/:employ_id
用途: 獲取指定員工的所有請假記錄。

** POST /api/v1/leave/:id
用途: 提交新的請假申請。

** GET /api/v1/leave/:employ_id/:date
用途: 獲取指定員工在某一天的請假記錄。

** PUT /api/v1/leave/:id/approve
用途: 同部門的主管批准請假申請。

** GET /api/v1/reports/attendance
用途: 獲取所有員工的出勤報告（按部門或日期篩選）。
