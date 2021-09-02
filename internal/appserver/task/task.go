package task

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"githup.com/htl/tcclienttest/internal/pkg/options"
	"githup.com/htl/tcclienttest/internal/pkg/util"
	"githup.com/htl/tcclienttest/pkg/log"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type terminalRebootDownloadTask struct{}

type TerminalLIst struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    data   `json:"data"`
}

type data struct {
	Transid int64     `json:"trans_id"`
	Content []content `json:"content"`
}

type content struct {
	Id        int64  `json:"id"`
	IpAddress string `json:"ipAddress"`
	UseStatus string `json:"useStatus"`
}

type taskTotal struct {
	deviceid int64
	//execTimes    int64
	totalType string
	times     int64
}

var gormdb *gorm.DB = func() *gorm.DB {
	gormdb, _ := options.NewMySQLOptions().NewClient()
	return gormdb
}()
var terminalIPs = make([]string, 0)
var ipMapDeviceid = make(map[string]int64)

func (t *terminalRebootDownloadTask) PreExecTask() error {
	//从mysql数据库获取终端ip地址
	rows, err := gormdb.Table("userstatus").Select("ip").Where("client_status='power_on' and userId=0").Rows()
	if err != nil {
		return err
	}

	for rows.Next() {
		ip := ""
		err = rows.Scan(&ip)
		if err != nil {
			return err
		}
		terminalIPs = append(terminalIPs, ip)
		if err != nil {
			return err
		}
		//fmt.Println(ip)
	}

	//把数据插入到Ip列表
	defer rows.Close()
	return nil
}

func (t *terminalRebootDownloadTask) ExecTask() error {

	for _, ip := range terminalIPs {

		rebootTerminal(ip)
	}
	getTerminalList()
	deleteTerminalImage()

	return nil
}

func (t *terminalRebootDownloadTask) AfterExecTask() error {
	return totalDownloadTimes()
}

func (t *terminalRebootDownloadTask) DeleteTask() error {
	return nil
}

func (t *terminalRebootDownloadTask) InstallTask() error {
	t.PreExecTask()
	t.ExecTask()
	t.AfterExecTask()

	return nil
}

func rebootTerminal(ip string) error {
	//发送daas server ip地址让终端自动注册到daas server
	_, err := util.RunLinuxCmd("sshpass", "-p", "tc123456", "scp", "-o", "StrictHostKeychecking=no", "-o", "ConnectTimeout=3", "./daas-serv.addr", fmt.Sprintf("root@%s:/mnt/disk", ip))
	if err != nil {
		log.Errorf("RunLinuxCmd failed :", err)
		return err
	}

	if err != nil {
		log.Errorf("RunLinuxCmd failed :", err)
		return err
	}

	//发送linux命令重启终端
	output, err := util.RunLinuxCmd("/opt/tci/bin/tc-client-control", "--reboot ", "--ip", ip, "--ssl")
	if err != nil {
		log.Errorf("RunLinuxCmd failed :", err)
		return err
	}

	if strings.Contains(output, "reboot client with IP") {
		fmt.Printf("reboot terminal ip %s sucessfully", ip)
	}
	return nil
}

func getTerminalList() error {
	url := `http://10.23.17.208:30001/seewo-eddes-manager/v1/schoolCode/0/deviceGroup/0/device/list`
	//contextType := "application/json"
	client := &http.Client{Timeout: 1 * time.Hour}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	jsmeta := strings.ReplaceAll(string(result), `\\`, "")
	fmt.Println(jsmeta)
	t := &TerminalLIst{}
	err = json.Unmarshal(result, t)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(t)

	for _, ipaddr := range t.Data.Content {
		if ipaddr.IpAddress == "" {
			continue
		}
		ipMapDeviceid[ipaddr.IpAddress] = ipaddr.Id
	}

	return nil
}

func deleteTerminalImage() error {
	url := fmt.Sprintf(`http://10.23.17.208:30001/seewo-eddes-manager/v1/device/%d/localCachedImage`, 654065338974998528)
	deleteParam := map[string]interface{}{
		"deviceImageNames": []string{"my_image1.img"},
	}
	jsonStr, _ := json.Marshal(deleteParam)
	req, _ := http.NewRequest("DELETE", url, bytes.NewReader(jsonStr))
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)

	t := &TerminalLIst{}
	err = json.Unmarshal(result, t)

	if t.Code == "999999" || t.Code == "" {
		return errors.New(t.Message)
	}

	return nil
}

func terminalDownloadImage() error {
	for ipStr, deviceid := range ipMapDeviceid {
		for _, ip := range terminalIPs {
			if strings.Contains(ipStr, ip) {
				func(devid int64) {
					//func Post(url string, data interface{}, contentType string) string {
					url := fmt.Sprintf(`http://10.23.17.208:30001/seewo-eddes-manager/v1/device/%d/resource/647577125674549248/image/24`, deviceid)
					contextType := "application/json"

					postParam := map[string]interface{}{
						"oriimg": "my_image1.img", "curimg": "my_image1.img",
					}

					// 超时时间：5秒
					client := &http.Client{Timeout: 1 * time.Hour}
					jsonStr, _ := json.Marshal(postParam)
					resp, err := client.Post(url, contextType, bytes.NewBuffer(jsonStr))
					if err != nil {
						panic(err)
					}
					defer resp.Body.Close()
					result, _ := ioutil.ReadAll(resp.Body)

					t := &TerminalLIst{}
					err = json.Unmarshal(result, t)

					if t.Code == "999999" || t.Code == "" {
						return
					}

					fmt.Println(string(result))
				}(deviceid)
			}
		}
	}
	return nil

}

func totalDownloadTimes() error {

	var state string
	var deviceid int64
	var num int64
	rowNums, err := gormdb.Table("tb_imageprogress").Select("state,device_id,count(*) num").Where(fmt.Sprintf("state in ('finished','failed') group by state,device_id")).Rows()
	if err != nil {
		log.Errorf("select for table tb_imagressgress failed: ", err)
		return err
	}

	total := make([]*taskTotal, 0)

	for rowNums.Next() {
		t := &taskTotal{}
		err := rowNums.Scan(&state, &deviceid, &num)
		if err != nil {
			log.Errorf("scan sql failed: ", err)
			return err
		}
		t.totalType = state
		t.deviceid = deviceid
		t.times = num

		total = append(total, t)

	}
	return nil
}
