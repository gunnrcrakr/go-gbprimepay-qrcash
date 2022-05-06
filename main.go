package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	e := echo.New()
	e.POST("/backgroundUrl", background)
	e.Logger.Fatal(e.Start(":5000"))

	call()
}

// Create QRCode
func call() error {

	// สร้าง instance Client
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// เตรียมข้อมูล Form Data
	postData := url.Values{}
	postData.Set("token", os.Getenv("token"))
	postData.Set("amount", "1.12")
	postData.Set("referenceNo", "testkuy003")
	postData.Set("backgroundUrl", os.Getenv("backgroundUrl"))
	postData.Set("detail", "test detail")
	postData.Set("customerName", "test customerName")
	postData.Set("customerEmail", "test customerEmail")
	postData.Set("merchantDefined1", "test merchantDefined1")
	// encode Form Data
	encodedPost := postData.Encode()

	// เตรียม POST
	// ถ้าต้องการให้ Response ออกมาเป็นรูปเลย ให้ตัด /text ข้างหลังออก
	req, err := http.NewRequest("POST", os.Getenv("api_url"), strings.NewReader(encodedPost))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// POST
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	// อย่าลืมปิด Response Body
	defer resp.Body.Close()

	// แปลง Response Body ให้อ่านง่าย
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}

	// แสดง Response Body พร้อมแปลงให้เป็น string
	fmt.Println(string(body))

	return nil
}

// Received from GBPrimePay after finish pay from QRCode
func background(c echo.Context) error {
	json_map := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&json_map)
	if err != nil {
		return err
	} else {
		// เอาตัวแปรไปใช้ต่อไป

		println(json_map)
		return c.JSON(http.StatusOK, json_map)
	}
}
