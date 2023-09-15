import request from '@/utils/request'

// 新建预测
export function createWeatherPredict(data) {
  return request({
    url: '/createWeatherPredict',
    method: 'post',
    data
  })
}

export function queryWeatherPredictList(data) {
  return request({
    url: '/queryWeatherPredictList',
    method: 'post',
    data
  })
}

export function updateWeatherPredict(data) {
  return request({
    url: '/updateWeatherPredict',
    method: 'post',
    data
  })
}
