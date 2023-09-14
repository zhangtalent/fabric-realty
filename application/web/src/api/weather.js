import request from '@/utils/request'

// 新建预测
export function createWeatherPredict(data) {
  return request({
    url: '/createWeatherPredict',
    method: 'post',
    data
  })
}

// 获取房地产信息(空json{}可以查询所有，指定proprietor可以查询指定业主名下房产)
export function queryWeatherPredictList(data) {
  return request({
    url: '/queryWeatherPredictList',
    method: 'post',
    data
  })
}
