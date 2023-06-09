//go:build !windows
// +build !windows

package main

import (
	"fmt"
	"time"
	"os/exec"
	"strconv"
	"strings"
)

func (server *NetFlowServer) flowCollect() (in, out int64, err error) {

	//因为Linux的流量是累加值，所以要通过历史的计数器相减再除采集间隔获取秒级的出入口流量

	server.mux.Lock()
	defer server.mux.Unlock()

	for _, collectInfo := range server.portsFlowCounters {

		if collectInfo.port <= 0 {
			continue
		}
        
		InOlderCounter := collectInfo.inFlow
		OutOlderCounter := collectInfo.outFlow

		InCurrentCounter, errIn := getPortInFlowByIptables(collectInfo.port)
		OutCurrentCounter, errOut := getPortOutFlowByIptables(collectInfo.port)

        var tempIn, tempOut int64

		if errIn == nil {
			collectInfo.inFlow = InCurrentCounter
			tempIn = (InCurrentCounter - InOlderCounter) / int64(server.collectIntervalSec)
			if tempIn < 0 {
				tempIn = 0
			}
			in += tempIn
		} else {
			LOG_ERROR(errIn)
		}

		if errOut == nil {
			collectInfo.outFlow = OutCurrentCounter
			tempOut = (OutCurrentCounter - OutOlderCounter) / int64(server.collectIntervalSec)
			if tempOut < 0 {
				tempOut = 0
			}
			out += tempOut
		} else {
			LOG_ERROR(errOut)
		}

		if tempIn>0 || tempOut>0 {
		    LOG_INFO(fmt.Sprintf("netflow port: %d, in_bytes: %d, out_bytes:%d, date:%s, Timestamp:%d", collectInfo.port, tempIn, tempOut, time.Now().Format("2006-01-02 15:04"), time.Now().Unix()))
		}

	}
	return
}

//通过iptables获取入站流量
func getPortInFlowByIptables(port int) (int64, error) {
	portStr := strconv.Itoa(port)
	cmd := []*exec.Cmd{
		exec.Command("iptables", "-L", "-v", "-n", "-x"),
		exec.Command("grep", "tcp dpt:"+portStr),
		exec.Command("awk", "{print $2}"),
		exec.Command("head", "-n", "1"),
	}
	data, err := ExecPipeLine(cmd...)
	if err != nil {
		return 0, err
	}

	data = strings.ReplaceAll(data, "\n", "")

	count, err := strconv.Atoi(data)
	if err != nil {
		return 0, err
	}
	return int64(count), nil
}

//通过iptables获取出站流量
func getPortOutFlowByIptables(port int) (int64, error) {
	portStr := strconv.Itoa(port)
	cmd := []*exec.Cmd{
		exec.Command("iptables", "-L", "-v", "-n", "-x"),
		exec.Command("grep", "tcp spt:"+portStr),
		exec.Command("awk", "{print $2}"),
		exec.Command("head", "-n", "1"),
	}
	data, err := ExecPipeLine(cmd...)
	if err != nil {
		return 0, err
	}

	data = strings.ReplaceAll(data, "\n", "")

	count, err := strconv.Atoi(data)
	if err != nil {
		return 0, err
	}
	return int64(count), nil
}

func (server *NetFlowServer) cleanRecords() {

	LOG_DEBUG(">>>>>>>>>>>> clean records")

	portsList := server.portsList
	for _, port := range portsList {
		LOG_INFO_F("init iptables with port : %d", port)

		dcmd := []*exec.Cmd{
			exec.Command("iptables", "-D", "INPUT", "-p", "tcp", "--dport", fmt.Sprintf("%d", port)),
		}
		ExecPipeLine(dcmd...)

		scmd := []*exec.Cmd{
			exec.Command("iptables", "-D", "INPUT", "-p", "tcp", "--sport", fmt.Sprintf("%d", port)),
		}
		ExecPipeLine(scmd...)

		dcmd1 := []*exec.Cmd{
			exec.Command("iptables", "-D", "OUTPUT", "-p", "tcp", "--dport", fmt.Sprintf("%d", port)),
		}
		ExecPipeLine(dcmd1...)

		scmd1 := []*exec.Cmd{
			exec.Command("iptables", "-D", "OUTPUT", "-p", "tcp", "--sport", fmt.Sprintf("%d", port)),
		}
		ExecPipeLine(scmd1...)
	}
}

func (server *NetFlowServer) setupRecords() {

	LOG_DEBUG(">>>>>>>>>>>> setup records")

	portsList := server.portsList
	for _, port := range portsList {
		LOG_INFO_F("init iptables with port : %d", port)

		cmd1 := []*exec.Cmd{
			exec.Command("iptables", "-A", "INPUT", "-p", "tcp", "--dport", fmt.Sprintf("%d", port)),
		}
		ExecPipeLine(cmd1...)

		cmd2 := []*exec.Cmd{
			exec.Command("iptables", "-A", "OUTPUT", "-p", "tcp", "--sport", fmt.Sprintf("%d", port)),
		}
		ExecPipeLine(cmd2...)
	}
}
