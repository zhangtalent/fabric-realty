<template>
  <div class="app-container">
    <el-form ref="ruleForm" v-loading="loading" :model="ruleForm" label-width="100px">

      <el-form-item label="AI ID" prop="proprietor">
        <el-input v-model="this.accountId" :disabled="true"/>
      </el-form-item>
      <el-form-item label="预测情况" prop="predictData">
        <el-input v-model="ruleForm.predictData" :precision="2" :step="0.1" :min="0" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="submitForm('ruleForm')">立即创建</el-button>
        <el-button @click="resetForm('ruleForm')">重置</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { queryAccountList } from '@/api/account'
import { createWeatherPredict } from '@/api/weather'

export default {
  name: 'createWeatherPredict',
  computed: {
    ...mapGetters([
      'accountId'
    ])
  },
  data() {
    return {
      ruleForm: {
        proprietor: this.accountId,
        predictData: "晴",
      },
      accountList: [],
      loading: false
    }
  },
  created() {
    console.log(this.accountId, 'sss')
  },
  methods: {
    submitForm(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.$confirm('是否立即创建?', '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'success'
          }).then(() => {
            this.loading = true
            createWeatherPredict({
              proprietor: this.accountId,
              predictData: this.ruleForm.predictData,
            }).then(response => {
              this.loading = false
              if (response !== null) {
                this.$message({
                  type: 'success',
                  message: '创建成功!'
                })
              } else {
                this.$message({
                  type: 'error',
                  message: '创建失败!'
                })
              }
            }).catch(_ => {
              this.loading = false
            })
          }).catch(() => {
            this.loading = false
            this.$message({
              type: 'info',
              message: '已取消创建'
            })
          })
        } else {
          return false
        }
      })
    },
    resetForm(formName) {
      this.$refs[formName].resetFields()
    },
  }
}
</script>

<style scoped>
</style>
