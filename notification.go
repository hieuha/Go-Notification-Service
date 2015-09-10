/**
Author: HieuHT (HieuHT@vnoss.org)
**/
package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/mail"
	"net/smtp"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

/** Data Structure **/
type RepsonseData struct {
	Response string
	Message  string
}

type Config struct {
	SMTP_EMAIL    string
	SMTP_PASSWORD string
	SMTP_SERVER   string
	SMTP_PORT     int
	SMS_API       string
	SMS_USER      string
	SMS_PASSWORD  string
	LISTEN        string

	DEFAULT_EMAIL    string
	DEFAULT_FULLNAME string
	DEFAULT_PHONE    string

	TIME_START_WORK int
	TIME_STOP_WORK  int

	DB_PATH string
}

/** Variables **/
var (
	rs                         = RepsonseData{Response: "", Message: ""}
	db                         = &sql.DB{}
	config                     = &Config{}
	message_sms, message_email string
	user_fullname              string
	user_email                 string
	user_phone                 string
	response                   string
	user_host                  []byte
)

/** Utilitis **/
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func loadConfig() *Config {
	conf := Config{}
	content, e := ioutil.ReadFile("./config.json")
	checkErr(e)
	err := json.Unmarshal(content, &conf)
	checkErr(err)
	return &conf
}

func wirteJson(data RepsonseData, w http.ResponseWriter) {
	js, err := json.Marshal(data)
	checkErr(err)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func doSendEmail(subject string, body string, to mail.Address) error {
	auth := smtp.PlainAuth("", config.SMTP_EMAIL, config.SMTP_PASSWORD, config.SMTP_SERVER)
	from := mail.Address{"Notification System", config.SMTP_EMAIL}
	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", config.SMTP_SERVER, config.SMTP_PORT),
		auth,
		from.Address,
		[]string{to.Address},
		[]byte(message),
	)
	return err
}

func doSendSMS(phones string, message string) string {
	res := ""
	//Code here
	return res
}

func checkOfficeWorkingTime() bool {
	server_time := time.Now()
	server_time_hour := server_time.Hour() //[0-23]
	if server_time_hour < config.TIME_START_WORK || server_time_hour > config.TIME_STOP_WORK {
		return true
	}
	return false
}

func checkExistInArray(str string, array []string) bool {
	for _, value := range array {
		if value == strings.TrimSpace(str) {
			return true
		}
	}
	return false
}

func sendNotification(email_subject string, email_to mail.Address, message_email string, user_phone string, message_sms string) string {
	err := doSendEmail(email_subject, message_email, email_to)
	checkErr(err)
	res := doSendSMS(user_phone, message_sms)
	return res
}

/** Functions **/
func notifyLogin(w http.ResponseWriter, r *http.Request) {
	/** Get Parameters **/
	username := strings.TrimSpace(r.URL.Query().Get("user"))
	remoteip := strings.TrimSpace(r.URL.Query().Get("remoteip"))
	servername := strings.TrimSpace(r.URL.Query().Get("servername"))
	extra := strings.TrimSpace(r.URL.Query().Get("extra"))
	/** The Query to get detail User **/
	row, db_err := db.Query("SELECT fullname, username, email, phone, host FROM `users` WHERE `username`=?", username)
	checkErr(db_err)

	/** Get User Structure **/
	user_time := time.Now().Format(time.RFC822Z)
	if row.Next() {
		db_err = row.Scan(&user_fullname, &username, &user_email, &user_phone, &user_host)
		checkErr(db_err)
	} else {
		user_email = config.DEFAULT_EMAIL
		user_fullname = config.DEFAULT_FULLNAME
		user_phone = config.DEFAULT_PHONE
		user_host = []byte("")
	}
	trusted_hosts := strings.Split(string(user_host), ",")

	/** Email **/
	message_email = fmt.Sprintf("%s - user %s from %s login %s, %s", user_time, username, remoteip, servername, extra)
	email_subject := fmt.Sprintf("Alert %s@%s login to %s", username, remoteip, servername)
	email_to := mail.Address{user_fullname, user_email}

	/** SMS **/
	message_sms = fmt.Sprintf("%s - %s@%s login to %s, %s", user_time, username, remoteip, servername, extra)
	response = "NO SENDING EMAIL or SMS"
	/** Send Notification **/
	if checkOfficeWorkingTime() {
		if !checkExistInArray(remoteip, trusted_hosts) && "" != string(user_host) {
			response = sendNotification(email_subject, email_to, message_email, user_phone, message_sms)
		} else if "" == string(user_host) {
			response = sendNotification(email_subject, email_to, message_email, user_phone, message_sms)
		}
	}

	/** JSON Response **/
	rs = RepsonseData{Response: response, Message: message_email}
	fmt.Println(message_email) //Write logs to Standard output
	wirteJson(rs, w)
}

/** Main **/
func main() {

	/** Load Json Configuration **/
	config = loadConfig()

	/** Route **/
	http.HandleFunc("/notify/login", notifyLogin)

	/** Database **/
	db, _ = sql.Open("sqlite3", config.DB_PATH)

	/** Web Server **/
	err := http.ListenAndServe(config.LISTEN, nil)
	checkErr(err)
}
