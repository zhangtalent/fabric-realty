<template>
  <div class="container">
    <el-alert
      type="success"
    >
      <p>账户ID: {{ accountId }}</p>
      <p>用户名: {{ userName }}</p>
      <p>余额: ￥{{ balance }} 元</p>
    </el-alert>
    <div v-if="weatherPredictList.length==0" style="text-align: center;">
      <el-alert
        title="查询不到数据"
        type="warning"
      />
    </div>
    <el-row v-loading="loading" :gutter="20">
      <el-col v-for="(val,index) in weatherPredictList" :key="index" :span="6" :offset="1">
        <el-card class="all-card" style="text-align: unset;">
          <div slot="header" class="clearfix">
            <span>{{ val.sellingStatus }}</span>
            <el-button v-if="roles[0] !== 'admin'&&(val.seller===accountId||val.buyer===accountId)&&val.sellingStatus!=='完成'&&val.sellingStatus!=='已过期'&&val.sellingStatus!=='已取消'" style="float: right; padding: 3px 0" type="text" @click="updateWeatherPredict(val,'cancelled')">取消</el-button>
            <el-button v-if="roles[0] !== 'admin'&&val.seller===accountId&&val.sellingStatus==='交付中'" style="float: right; padding: 3px 8px" type="text" @click="updateWeatherPredict(val,'done')">确认收款</el-button>
            <el-button v-if="roles[0] !== 'admin'&&val.sellingStatus==='销售中'&&val.seller!==accountId" style="float: right; padding: 3px 0" type="text" @click="createWeatherPredictByBuy(val)">购买</el-button>
          </div>
          <div class="item">
            <el-tag>预测ID: </el-tag>
            <span>{{ val.weatherPredictId }}</span>
          </div>
          <div class="item">
            <el-tag type="success">AI ID: </el-tag>
            <span>{{ val.proprietor }}</span>
          </div>
          <div class="item">
            <el-tag type="danger">预测: </el-tag>
            <span>{{ val.predictData }}</span>
          </div>
          <div class="item">
            <el-tag type="warning">状态: </el-tag>
            <span>{{ val.valiateStatus }}</span>
          </div>
          <div class="item">
            <el-tag type="info">创建时间: </el-tag>
            <span>{{ val.createTime }}</span>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { queryWeatherPredictList, createWeatherPredict, updateWeatherPredict } from '@/api/weather'

export default {
  name: 'AllWeatherPredict',
  data() {
    return {
      loading: true,
      weatherPredictList: []
    }
  },
  computed: {
    ...mapGetters([
      'accountId',
      'roles',
      'userName',
      'balance'
    ])
  },
  created() {
    queryWeatherPredictList().then(response => {
      if (response !== null) {
        this.weatherPredictList = response
      }
      this.loading = false
    }).catch(_ => {
      this.loading = false
    })
  },
  methods: {
    createWeatherPredictByBuy(item) {
      this.$confirm('是否立即购买?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'success'
      }).then(() => {
        this.loading = true
        createWeatherPredictByBuy({
          buyer: this.accountId,
          objectOfSale: item.objectOfSale,
          seller: item.seller
        }).then(response => {
          this.loading = false
          if (response !== null) {
            this.$message({
              type: 'success',
              message: '购买成功!'
            })
          } else {
            this.$message({
              type: 'error',
              message: '购买失败!'
            })
          }
          setTimeout(() => {
            window.location.reload()
          }, 1000)
        }).catch(_ => {
          this.loading = false
        })
      }).catch(() => {
        this.loading = false
        this.$message({
          type: 'info',
          message: '已取消购买'
        })
      })
    },
    updateWeatherPredict(item, type) {
      let tip = ''
      if (type === 'done') {
        tip = '确认收款'
      } else {
        tip = '取消操作'
      }
      this.$confirm('是否要' + tip + '?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'success'
      }).then(() => {
        this.loading = true
        updateWeatherPredict({
          buyer: item.buyer,
          objectOfSale: item.objectOfSale,
          seller: item.seller,
          status: type
        }).then(response => {
          this.loading = false
          if (response !== null) {
            this.$message({
              type: 'success',
              message: tip + '操作成功!'
            })
          } else {
            this.$message({
              type: 'error',
              message: tip + '操作失败!'
            })
          }
          setTimeout(() => {
            window.location.reload()
          }, 1000)
        }).catch(_ => {
          this.loading = false
        })
      }).catch(() => {
        this.loading = false
        this.$message({
          type: 'info',
          message: '已取消' + tip
        })
      })
    }
  }
}

</script>

<style>
  .container{
    width: 100%;
    text-align: center;
    min-height: 100%;
    overflow: hidden;
  }
  .tag {
    float: left;
  }

  .item {
    font-size: 14px;
    margin-bottom: 18px;
    color: #999;
  }

  .clearfix:before,
  .clearfix:after {
    display: table;
  }
  .clearfix:after {
    clear: both
  }

  .all-card {
    width: 280px;
    height: 380px;
    margin: 18px;
  }
</style>
