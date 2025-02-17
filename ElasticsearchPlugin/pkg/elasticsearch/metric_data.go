/*
* Copyright 2018 Nokia
* Licensed under BSD 3-Clause Clear License,
* see LICENSE file for details.
 */

package elasticsearch

import (
	"elasticsearchplugin/pkg/config"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"path"
	"strconv"
	"strings"
	"time"
)

type pmResponse []struct {
	PMData       map[string]interface{} `json:"pm_data"`
	PMDataSource map[string]interface{} `json:"pm_data_source"`
	Timestamp    time.Time              `json:"timestamp"`
}

type fmResponse []struct {
	FMData       map[string]interface{} `json:"fm_data"`
	FMDataSource map[string]interface{} `json:"fm_data_source"`
	Timestamp    time.Time              `json:"timestamp"`
}

func pushFMData(filePath string, esConf config.ElasticsearchConf) {
	elkURL := esConf.URL + elkBulkAPI
	log.Infof("Pushing data from %s to elasticsearch", filePath)
	data, err := readFile(filePath)
	if err != nil {
		return
	}

	fileName := path.Base(filePath)
	metricType := strings.Split(fileName, "_")[1]
	metricType = strings.ToLower(metricType)

	var resp fmResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Errorf("Unable to unmarshal json data %s", filePath)
		return
	}

	var postData string
	index := strings.Join([]string{metricType, "fm"}, "-")
	for i, d := range resp {
		if d.FMData["event_time"] == nil {
			continue
		}
		dn := d.FMDataSource["dn"].(string)
		eventTime := d.FMData["event_time"].(string)
		hwID := d.FMDataSource["hw_id"].(string)
		alarmID := d.FMData["alarm_identifier"].(string)
		var id string
		if metricType == "radio" {
			specificProb := d.FMData["specific_problem"].(string)
			id = strings.Join([]string{hwID, dn, alarmID, specificProb, eventTime}, "_")
		} else {
			faultID := d.FMData["fault_id"].(string)
			keys := []string{metricType, hwID, dn, alarmID}
			if faultID != "" {
				keys = append(keys, faultID)
			}
			keys = append(keys, eventTime)
			id = strings.Join(keys, "_")
		}
		d.Timestamp = time.Now()
		source, _ := json.Marshal(d)
		postData += `{"index": {"_index": "` + index + `", "_id": "` + id + `"}}` + "\n"
		postData += string(source) + "\n"
		if i != 0 && i%elkNoOfRecordsPerAPI == 0 || i == len(resp)-1 {
			pushData(elkURL, esConf.User, esConf.Password, postData, filePath)
			postData = ""
		}
	}
}

func pushPMData(filePath string, esConf config.ElasticsearchConf) {
	elkURL := esConf.URL + elkBulkAPI
	log.Infof("Pushing data from %s to elasticsearch", filePath)
	data, err := readFile(filePath)
	if err != nil {
		return
	}

	var resp pmResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Errorf("Unable to unmarshal json data %s", filePath)
		return
	}

	fileName := path.Base(filePath)
	metricType := strings.Split(fileName, "_")[1]
	metricType = strings.ToLower(metricType)

	var postData string
	currTime := time.Now().UTC()
	for i, d := range resp {
		hwID := d.PMDataSource["hw_id"].(string)
		dn := d.PMDataSource["dn"].(string)
		eventTime := d.PMDataSource["timestamp"].(string)
		technology := d.PMDataSource["technology"].(string)
		var objectType string
		for k := range d.PMData {
			objectType = k
			break
		}

		var index, id string
		if metricType == "radio" {
			objectType = strings.ToLower(objectType[:strings.LastIndex(objectType, "_")])
			index = strings.Join([]string{strings.ToLower(technology), "pm", objectType, fmt.Sprintf("%02d", currTime.Month()), strconv.Itoa(currTime.Year())}, "-")
			id = strings.Join([]string{hwID, eventTime, dn}, "_")
		} else {
			index = strings.Join([]string{metricType, "pm"}, "-")
			id = strings.Join([]string{objectType, dn, eventTime}, "_")
		}

		d.Timestamp = time.Now()
		source, _ := json.Marshal(d)
		postData += `{"index": {"_index": "` + index + `", "_id": "` + id + `"}}` + "\n"
		postData += string(source) + "\n"
		if i != 0 && i%elkNoOfRecordsPerAPI == 0 || i == len(resp)-1 {
			pushData(elkURL, esConf.User, esConf.Password, postData, filePath)
			postData = ""
		}
	}
}
