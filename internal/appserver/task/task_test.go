package task

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"githup.com/htl/tcclienttest/internal/pkg/options"
	"githup.com/htl/tcclienttest/internal/pkg/util"
	"githup.com/htl/tcclienttest/pkg/log"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func Test_AutoRegiestAndUpgradeSystem(t *testing.T) {

	//往终端发送dass_server地址，让终端自动注册
	util.RunLinuxCmd("sshpass", "-p", "tc123456", "scp", "-o", "StrictHostKeychecking=no", "-o", "ConnectTimeout=3", "./daas-serv.addr", "root@172.20.13.78:/mnt/disk")

	//远程更新脚本
	util.RunRemoteShellcmd("ip", "sh /mnt/disk/test.sh")

}

type terminalLIst struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    data   `json:"data"`
}

type Data struct {
	Transid int64     `json:"trans_id"`
	Content []content `json:"content"`
}

type Content struct {
	Id        int64  `json:"id"`
	IpAddress string `json:"ipAddress"`
	UseStatus string `json:"useStatus"`
}

type TaskTotal struct {
	deviceid int64
	//execTimes    int64
	totalType string
	times     int64
}

func Test_installTask(t *testing.T) {
	//初始化mysql数据库
	ops := options.NewMySQLOptions()
	gormdb, err := ops.NewClient()
	if err != nil {
		fmt.Println(err)
	}

	//从mysql数据库获取终端ip地址
	rows, err := gormdb.Table("userstatus").Select("ip").Where("client_status='power_on' and userId=0").Rows()
	if err != nil {
		fmt.Println(err)
	}

	terminalIPs := make([]string, 0)
	//把数据插入到Ip列表
	for rows.Next() {
		ip := ""
		err = rows.Scan(&ip)
		terminalIPs = append(terminalIPs, ip)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(ip)
	}

	defer rows.Close()
	for _, ip := range terminalIPs {
		//发送daas server ip地址让终端自动注册到daas server
		_, err := util.RunLinuxCmd("sshpass", "-p", "tc123456", "scp", "-o", "StrictHostKeychecking=no", "-o", "ConnectTimeout=3", "./daas-serv.addr", fmt.Sprintf("root@%s:/mnt/disk", ip))
		if err != nil {
			log.Errorf("RunLinuxCmd failed :", err)
			return
		}

		if err != nil {
			log.Errorf("RunLinuxCmd failed :", err)
			return
		}

		//发送linux命令重启终端
		output, err := util.RunLinuxCmd("/opt/tci/bin/tc-client-control", "--reboot ", "--ip", ip, "--ssl")
		if err != nil {
			log.Errorf("RunLinuxCmd failed :", err)
			return
		}

		if strings.Contains(output, "reboot client with IP") {
			fmt.Printf("reboot terminal ip %s sucessfully", ip)
		}

		//check terminal power on

		//var numCount int
		//for numCount == 0 {
		//	sqlrow := gormdb.Table("userstatus").Select("1 as num").Where(fmt.Sprintf("ip='%s' and client_status='power_on' ", ip)).Row()
		//	sqlrow.Scan(&numCount)
		//}

	}

	ipMapDeviceid := make(map[string]int64)
	//获取终端列表，拿到终端注册到daas server中的注册码，deviceid
	getTerminalListFunc := func() {
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
	}

	getTerminalListFunc()

	deleteTerminalImageFunc := func() error {
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

	deleteTerminalImageFunc()

	//终端下载镜像
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

}
