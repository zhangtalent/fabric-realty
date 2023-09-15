package v1

import (
	bc "application/blockchain"
	"application/model"
	"application/pkg/app"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WeatherPredictRequestBody struct {
	Proprietor  string `json:"proprietor"`  //所有者(AI ID)
	PredictData string `json:"predictData"` //预测数据
}

type WeatherPredictUpdateRequestBody struct {
	RealityData string `json:"realityData"`
	Date        string `json:"date"`
}
type WeatherPredictQueryRequestBody struct {
	Proprietor string `json:"proprietor"`
}

func CreateWeatherPredict(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(WeatherPredictRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.Proprietor))
	bodyBytes = append(bodyBytes, []byte(body.PredictData))
	//调用智能合约
	resp, err := bc.ChannelExecute("createWeatherPredict", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	var data map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}

func QueryWeatherPredictList(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(WeatherPredictQueryRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	var bodyBytes [][]byte
	if body.Proprietor != "" {
		bodyBytes = append(bodyBytes, []byte(body.Proprietor))
	}
	//调用智能合约
	resp, err := bc.ChannelQuery("queryWeatherPredictList", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	// 反序列化json
	var data []map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}

func UpdateWeatherPredict(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(WeatherPredictUpdateRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	resp, err := bc.ChannelQuery("queryWeatherPredictList", [][]byte{}) //调用智能合约
	if err != nil {
		fmt.Printf("-quer失败%s", err.Error())
		return
	}
	// 反序列化json
	var dataWeather []model.WeatherPredict
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &dataWeather); err != nil {
		fmt.Printf("定时任务-反序列化json失败%s", err.Error())
		return
	}
	for _, v := range dataWeather {
		var bodyBytes [][]byte
		bodyBytes = append(bodyBytes, []byte(body.RealityData))
		bodyBytes = append(bodyBytes, []byte(body.Date))
		bodyBytes = append(bodyBytes, []byte(v.Proprietor))
		bodyBytes = append(bodyBytes, []byte(v.WeatherPredictID))
		//调用智能合约
		resp, err := bc.ChannelExecute("updateWeather", bodyBytes)
		if err != nil {
			appG.Response(http.StatusInternalServerError, "失败", err.Error())
			return
		}
		fmt.Println(resp)
	}
	var data map[string]interface{}
	appG.Response(http.StatusOK, "成功", data)
}
