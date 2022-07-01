package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/xuri/excelize/v2"
)

func main() {
	records := readCsvFile("跑绩效.csv")
	for _, row := range records {
		fmt.Println(row[0])
		run(row[0], "amazon.account")
	}
	// run("AL4VR1T7180RT", "amazon.message")
}

func readCsvFile(filename string) [][]string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf(err.Error()))
	}
	fmt.Println(dirname)

	f, err := os.Open(path.Join(dirname, filename))
	if err != nil {
		log.Fatal("Unable to read input file "+filename, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filename, err)
	}

	return records
}

func readFromExcel(filename string) (result []string) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf(err.Error()))
	}
	fmt.Println(dirname)

	f, err := excelize.OpenFile(path.Join(dirname, filename))
	if err != nil {
		fmt.Println(err)
		return []string{}
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return []string{}
	}

	for i, row := range rows {
		if row == nil || i == 0 {
			continue
		}
		fmt.Println(row[0])
		result = append(result, row[0])
	}

	return result
}

func run(shopId string, taskId string) {
	ctx, cancel := context.WithTimeout(context.TODO(), 20*time.Second)
	defer cancel()
	err := removeShopTaskCache(shopId, ctx)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = assignTask(shopId, taskId, ctx)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func removeShopTaskCache(shopId string, ctx context.Context) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "r-wz9hboqezpxr2d2vh0.redis.rds.aliyuncs.com:6379",
		Password: "Tj^IV%kOVN73yJ3Y",
		DB:       2,
	})

	key := "TASK:AMAZON:MESSAGE:" + shopId
	fmt.Println(key)
	val, err := rdb.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	fmt.Println(val)
	return nil
}

func assignTask(shopId string, taskId string, ctx context.Context) error {
	url := "https://chat-svc-rest.ziniao.com/api/task/task"
	data := make(map[string]string)
	data["shopId"] = shopId
	data["taskId"] = taskId
	b, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("http.NewRequest Error: %s", err.Error())
	}

	req = req.WithContext(ctx)
	req.Header.Set("x-api-key", "UOtmiQs03CKZB9V5QGYblpxknIuyVvBnsLZX5xYHFxk")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do Error: %s", err.Error())
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	return nil
}
