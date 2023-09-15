package api

import (
	"chaincode/model"
	"chaincode/pkg/utils"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// CreateWeatherPredict 新建预测
func CreateWeatherPredict(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 验证参数
	if len(args) != 2 {
		return shim.Error("参数个数不满足")
	}
	proprietor := args[0]
	predictData := args[1]
	if predictData == "" || proprietor == "" {
		return shim.Error("参数存在空值")
	}
	//判断AI ID是否存在
	resultsProprietor, err := utils.GetStateByPartialCompositeKeys(stub, model.AccountKey, []string{proprietor})
	if err != nil || len(resultsProprietor) != 1 {
		return shim.Error(fmt.Sprintf("proprietor信息验证失败%s", err))
	}
	createTime, _ := stub.GetTxTimestamp()
	unixTime := createTime.GetSeconds()
	unixNano := createTime.GetNanos()
	unixTimestamp := int64(unixTime) + int64(unixNano)/1e9
	weatherPredict := &model.WeatherPredict{
		WeatherPredictID: stub.GetTxID()[:16],
		Proprietor:       proprietor,
		ValiateStatus:    model.ValiateStatusConstant()["waiting"],
		CreateTime:       unixTimestamp,
		PredictData:      predictData,
	}
	// 写入账本
	if err := utils.WriteLedger(weatherPredict, stub, model.WeatherPredictKey, []string{weatherPredict.Proprietor, weatherPredict.WeatherPredictID}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	//将成功创建的信息返回
	weatherPredictByte, err := json.Marshal(weatherPredict)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	// 成功返回
	return shim.Success(weatherPredictByte)
}

// QueryRealEstateList 查询所有AI当日预测
func QueryWeatherPredictList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var weatherPredictList []model.WeatherPredict
	results, err := utils.GetStateByPartialCompositeKeys2(stub, model.WeatherPredictKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var weatherPredict model.WeatherPredict
			err := json.Unmarshal(v, &weatherPredict)
			if err != nil {
				return shim.Error(fmt.Sprintf("QueryWeatherPredictList-反序列化出错: %s", err))
			}
			//unixTimestamp := int64(weatherPredict.CreateTime)
			// 将 Unix 时间戳转换为 time.Time 类型
			//t := time.Unix(unixTimestamp, 0)
			// 获取当前日期
			//currentDate := time.Now().UTC().Truncate(24 * time.Hour)
			// 判断是否是今天
			//isToday := t.After(currentDate) && t.Before(currentDate.Add(24*time.Hour))
			//if isToday {
			weatherPredictList = append(weatherPredictList, weatherPredict)
			//}
		}
	}
	weatherPredictListByte, err := json.Marshal(weatherPredictList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryWeatherPredictList-序列化出错: %s", err))
	}
	return shim.Success(weatherPredictListByte)
}

func IsUnixTimestampInTargetDate(unixTimestamp int64, targetDateString string) (bool, error) {
	// 解析目标日期字符串为 time.Time 类型
	targetDate, err := time.Parse("2006-01-02", targetDateString)
	if err != nil {
		return false, fmt.Errorf("无法解析目标日期字符串: %s %s", err, targetDateString)
	}
	// 将 Unix 时间戳转换为 time.Time 类型
	t := time.Unix(unixTimestamp, 0)
	// 将时间戳转换为目标日期的 UTC 时间
	utcTime := t.UTC()
	// 判断是否在目标日期内
	isTargetDate := utcTime.Year() == targetDate.Year() && utcTime.YearDay() == targetDate.YearDay()
	return isTargetDate, nil
}

// UpdateWeatherPredict 新建预测
func UpdateWeather(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 验证参数
	if len(args) != 4 {
		return shim.Error("参数个数不满足")
	}
	realityData := args[0]
	valiateTime := args[1]
	proprietor := args[2]
	weatherPredictID := args[3]
	if realityData == "" || valiateTime == "" ||
		proprietor == "" || weatherPredictID == "" {
		return shim.Error("参数存在空值")
	}
	//判断AI ID是否存在
	resultsProprietor, err := utils.GetStateByPartialCompositeKeys(stub, model.AccountKey, []string{proprietor})
	if err != nil || len(resultsProprietor) != 1 {
		return shim.Error(fmt.Sprintf("proprietor信息验证失败%s", err))
	}
	resultsWeatherPredict, err := utils.GetStateByPartialCompositeKeys2(stub, model.WeatherPredictKey, []string{proprietor, weatherPredictID})
	if err != nil || len(resultsWeatherPredict) != 1 {
		return shim.Error(fmt.Sprintf("根据%s和%s获取失败: %s", proprietor, weatherPredictID, err))
	}
	var weatherPredict model.WeatherPredict
	if err = json.Unmarshal(resultsWeatherPredict[0], &weatherPredict); err != nil {
		return shim.Error(fmt.Sprintf("UpdateWeatherPredict-反序列化出错: %s", err))
	}
	if weatherPredict.PredictData == realityData {
		weatherPredict.ValiateStatus = model.ValiateStatusConstant()["correct"]
	} else {
		weatherPredict.ValiateStatus = model.ValiateStatusConstant()["error"]
	}
	if err := utils.WriteLedger(weatherPredict, stub, model.WeatherPredictKey, []string{weatherPredict.Proprietor, weatherPredict.WeatherPredictID}); err != nil {
		fmt.Println(err)
		return shim.Error(fmt.Sprintf("%s", err))
	}
	//将成功创建的信息返回
	weatherPredictByte, err := json.Marshal(weatherPredict)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	// 成功返回
	return shim.Success(weatherPredictByte)
}

// UpdateWeatherPredict 新建预测
func UpdateWeatherPredict(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 验证参数
	if len(args) != 4 {
		return shim.Error("参数个数不满足")
	}

	proprietor := args[2]
	if proprietor == "" {
		return shim.Error("参数存在空值")
	}

	fmt.Println(proprietor)

	//判断AI ID是否存在
	resultsProprietor, err := utils.GetStateByPartialCompositeKeys(stub, model.AccountKey, []string{proprietor})
	if err != nil || len(resultsProprietor) != 1 {
		return shim.Error(fmt.Sprintf("proprietor信息验证失败%s", err))
	}

	createTime, _ := stub.GetTxTimestamp()
	unixTime := createTime.GetSeconds()
	unixNano := createTime.GetNanos()
	unixTimestamp := int64(unixTime) + int64(unixNano)/1e9
	weatherPredict := &model.WeatherPredict{
		WeatherPredictID: stub.GetTxID()[:16],
		Proprietor:       proprietor,
		ValiateStatus:    model.ValiateStatusConstant()["waiting"],
		CreateTime:       unixTimestamp,
		PredictData:      "GOOD",
	}
	// 写入账本
	if err := utils.WriteLedger(weatherPredict, stub, model.WeatherPredictKey, []string{weatherPredict.Proprietor, weatherPredict.WeatherPredictID}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	//将成功创建的信息返回
	weatherPredictByte, err := json.Marshal(weatherPredict)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	// 成功返回
	return shim.Success(weatherPredictByte)
}

// UpdateWeatherPredict
func dddUpdateWeatherPredict(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 验证参数
	if len(args) != 4 {
		return shim.Error("参数个数不满足")
	}
	//realityData := args[0]
	//valiateTime := args[1]
	proprietor := args[2]
	//weatherPredictID := args[3]
	//判断AI ID是否存在
	resultsProprietor, err := utils.GetStateByPartialCompositeKeys(stub, model.AccountKey, []string{proprietor})
	if err != nil || len(resultsProprietor) != 1 {
		return shim.Error(fmt.Sprintf("proprietor信息验证失败%s", err))
	}
	createTime, _ := stub.GetTxTimestamp()
	unixTime := createTime.GetSeconds()
	unixNano := createTime.GetNanos()
	unixTimestamp := int64(unixTime) + int64(unixNano)/1e9
	weatherPredict := &model.WeatherPredict{
		WeatherPredictID: stub.GetTxID()[:16],
		Proprietor:       proprietor,
		ValiateStatus:    model.ValiateStatusConstant()["waiting"],
		CreateTime:       unixTimestamp,
		PredictData:      "sss",
	}
	// 写入账本
	if err := utils.WriteLedger(weatherPredict, stub, model.WeatherPredictKey, []string{weatherPredict.Proprietor, weatherPredict.WeatherPredictID}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	/*resultsWeatherPredict, err := utils.GetStateByPartialCompositeKeys2(stub, model.WeatherPredictKey, []string{proprietor, weatherPredictID})
	if err != nil || len(resultsWeatherPredict) != 1 {
		return shim.Error(fmt.Sprintf("根据%s和%s获取失败: %s", proprietor, weatherPredictID, err))
	}
	var weatherPredict model.WeatherPredict
	if err = json.Unmarshal(resultsWeatherPredict[0], &weatherPredict); err != nil {
		return shim.Error(fmt.Sprintf("UpdateWeatherPredict-反序列化出错: %s", err))
	}
	if weatherPredict.ValiateStatus == model.ValiateStatusConstant()["waiting"] {
		unixTimestamp := int64(weatherPredict.CreateTime)
		// 假设目标日期是一个字符串
		targetDateString := valiateTime
		// 调用函数进行判断
		isInTargetDate, err := IsUnixTimestampInTargetDate(unixTimestamp, targetDateString)
		fmt.Println(unixTimestamp)
		fmt.Println(targetDateString)
		fmt.Println(isInTargetDate)
		if err != nil {
			return shim.Error(fmt.Sprintf("UpdateWeatherPredict-日期转换出错: %s", err))
		}
		if isInTargetDate {
			fmt.Println(1)
		}
		if weatherPredict.PredictData == realityData {
			weatherPredict.ValiateStatus = model.ValiateStatusConstant()["correct"]
		} else {
			weatherPredict.ValiateStatus = model.ValiateStatusConstant()["error"]
		}
		fmt.Println(weatherPredict)
		fmt.Println(weatherPredict.PredictData == realityData)
		if err := utils.WriteLedger(weatherPredict, stub, model.WeatherPredictKey, []string{weatherPredict.Proprietor, weatherPredict.WeatherPredictID}); err != nil {
			fmt.Println(err)
			return shim.Error(fmt.Sprintf("%s", err))
		}
	}*/
	weatherPredictByte, err := json.Marshal(weatherPredict)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	/*for i, _ := range results {
		var weatherPredict model.WeatherPredict
		if err = json.Unmarshal(results[i], &weatherPredict); err != nil {
			return shim.Error(fmt.Sprintf("CreateSellingBuy-反序列化出错: %s", err))
		}
		weatherPredict.ValiateStatus = "HHHH"
		if err := utils.WriteLedger(weatherPredict, stub, model.WeatherPredictKey, []string{weatherPredict.Proprietor, weatherPredict.WeatherPredictID}); err != nil {
			return shim.Error(fmt.Sprintf("将status更新 失败%s", err))
		}
		if v != nil {
			var weatherPredict model.WeatherPredict
			err := json.Unmarshal(v, &weatherPredict)
			if weatherPredict.ValiateStatus != model.ValiateStatusConstant()["waiting"] {
				continue
			}
			if err != nil {
				return shim.Error(fmt.Sprintf("UpdateWeatherPredict-反序列化出错: %s", err))
			}
			unixTimestamp := int64(weatherPredict.CreateTime)
			// 假设目标日期是一个字符串
			targetDateString := valiateTime
			// 调用函数进行判断
			isInTargetDate, err := IsUnixTimestampInTargetDate(unixTimestamp, targetDateString)
			fmt.Println(unixTimestamp)
			fmt.Println(targetDateString)
			fmt.Println(isInTargetDate)
			if err != nil {
				return shim.Error(fmt.Sprintf("UpdateWeatherPredict-日期转换出错: %s", err))
			}
			if isInTargetDate {
				fmt.Println(1)
			}
			if weatherPredict.PredictData == realityData {
				weatherPredict.ValiateStatus = model.ValiateStatusConstant()["correct"]
			} else {
				weatherPredict.ValiateStatus = model.ValiateStatusConstant()["error"]
			}
			fmt.Println(weatherPredict)
			fmt.Println(weatherPredict.PredictData == realityData)
			//清除原来的信息
			if err := utils.DelLedger(stub, model.WeatherPredictKey, []string{weatherPredict.Proprietor, weatherPredict.WeatherPredictID}); err != nil {
				return shim.Error(fmt.Sprintf("%s", err))
			}
			if err := utils.WriteLedger(weatherPredict, stub, model.WeatherPredictKey, []string{weatherPredict.Proprietor, weatherPredict.WeatherPredictID}); err != nil {
				fmt.Println(err)
				return shim.Error(fmt.Sprintf("%s", err))
			}
		}
	}*/

	return shim.Success(weatherPredictByte)
}
