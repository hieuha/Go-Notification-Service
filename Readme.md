### Biên dịch
1. `go build notification.go`  

### Hướng dẫn
1. Chép file  `notification.init.d.centos` hoặc `notification.init.d.ubuntu` vào thư mục `/etc/init.d/notification`. Sau đó `chmod +x /etc/init.d/notification`
2. Tạo folder chứa chương trình `mkdir -p /opt/notification`
3. Cấu trúc thư mục  
```
.
├── config.json # Lưu thông tin cấu hình
├── users.sqlite # Lưu thông tin của user
└── notification
```
4. Chép file `login-notify.sh` vào trong `/etc/profile.d/login-notify.sh` của server nào cần monitor.
5. Khởi động dịch vụ `notification`

### Thông tin cấu hình
```
"SMTP_EMAIL": "",
"SMTP_PASSWORD": "",
"SMTP_SERVER": "",
"SMTP_PORT": 25,
"SMS_API": "",
"SMS_USER": "",
"SMS_PASSWORD": "",
"LISTEN": ":8080",
"DEFAULT_EMAIL": "",
"DEFAULT_FULLNAME": "",
"DEFAULT_PHONE": "",
"TIME_START_WORK": 8,
"TIME_STOP_WORK": 18,
"DB_PATH": "./users.sqlite"
```
