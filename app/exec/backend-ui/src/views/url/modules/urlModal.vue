<template>
  <el-dialog
    :before-close="clearData"
    :title="amendOrAdd"
    :visible.sync="theVisible"
    class="dialog"
    width="520px"
  >
    <el-form
      ref="formValue"
      :model="formValue"
      :rules="formRules"
      class="demo-addunit"
      label-width="100px"
    >
      <el-form-item v-if="model.id" :label="'短链接:'" class="demo-input">
        <el-input
          v-model="formValue.shortKey"
          :disabled="true"
          placeholder="请输入"
        ></el-input>
      </el-form-item>
      <el-form-item :label="'原始链接'" class="demo-input" prop="fullUrl">
        <el-input v-model="formValue.fullUrl"></el-input>
      </el-form-item>
      <el-form-item class="demo-input" label="冻结:" prop="isFrozen">
        <el-radio-group v-model="formValue.isFrozen">
          <el-radio :label="1">是</el-radio>
          <el-radio :label="0">否</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item class="demo-input" label="是否永久:" prop="perpetual">
        <el-radio-group v-model="formValue.perpetual">
          <el-radio :label="1">是</el-radio>
          <el-radio :label="0">否</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item
        v-if="formValue.perpetual === 0 ? true : false"
        class="demo-input"
        label="过期时间:"
        prop="expirationTime"
      >
        <el-date-picker
          v-model="formValue.expirationTime"
          placeholder="选择日期时间"
          type="datetime"
        ></el-date-picker>
      </el-form-item>
    </el-form>
    <span slot="footer" class="dialog-footer">
      <el-button @click="clearData">取消</el-button>
      <el-button class="save-btn" @click="sendParams('formValue')"
        >保存</el-button
      >
    </span>
  </el-dialog>
</template>
<script>
import { urlInfo, changeShortUrl, addShortUrl } from "@/api/request-data";
import { dateFormat } from "@/utils/date-format";
export default {
  data() {
    return {
      amendOrAdd: "",
      theVisible: false,
      model: {},
      formRules: {
        fullUrl: [
          { required: true, message: "请输入链接地址", trigger: "blur" },
        ],
        isFrozen: [
          { required: true, message: "请选择是否冻结", trigger: "blur" },
        ],
        perpetual: [
          { required: true, message: "请选择是否永久", trigger: "blur" },
        ],
        expirationTime: [
          { required: true, message: "请选择过期时间", trigger: "blur" },
        ],
      },
      formValue: {
        fullUrl: "",
        shortKey: "",
        isFrozen: 0,
        perpetual: 1,
        id: 0,
        expirationTime: "",
      },
    };
  },
  methods: {
    show(data) {
      this.model = Object.assign(this.model, data);
      if (data) {
        this.amendOrAdd = "编辑短链接";
      } else {
        this.amendOrAdd = "新增短链接";
      }
      if (this.model.id) {
        urlInfo(this.model.id).then((res) => {
          let { expirationTime, fullUrl, isFrozen, shortKey } = res.data;
          console.log(expirationTime, dateFormat(expirationTime));
          this.formValue = {
            fullUrl: fullUrl,
            isFrozen: isFrozen,
            shortKey: shortKey,
            perpetual: expirationTime == 0 ? 1 : 0,
            expirationTime:
              expirationTime == 0 ? "" : dateFormat(expirationTime),
          };
          this.theVisible = true;
        });
      } else {
        this.theVisible = true;
      }
    },
    sendParams(formName) {
      this.$refs[formName].validate(async (valid) => {
        if (valid) {
          const params = {};
          params.fullUrl = this.formValue.fullUrl;
          params.isFrozen = this.formValue.isFrozen;
          if (this.formValue.perpetual == 1) {
            params.expirationTime = 0;
          } else {
            params.expirationTime =
              this.formValue.expirationTime.getTime() / 1000;
          }

          if (!this.model.id) {
            const requestData = await addShortUrl(params);
            if (requestData.code === 200) {
              this.pageNum = 1;
              this.$message({
                type: "success",
                message: `新增成功!`,
              });
              this.theVisible = false;
              this.$emit("getTableData", 1);
            } else {
              console.log("新增失败");
            }
          } else {
            const requestData = await changeShortUrl(params, this.model.id);
            if (requestData.code === 200) {
              this.$message({
                type: "success",
                message: `修改成功!`,
              });
              this.$emit("getTableData");
              this.theVisible = false;
            } else {
              console.log("修改失败");
            }
          }
        } else {
          console.log("error submit!!");
          return false;
        }
      });
    },
    clearData() {
      this.formValue = {
        fullUrl: "",
        shortKey: "",
        isFrozen: 0,
        perpetual: 1,
        id: 0,
        expirationTime: "",
      };
      this.model = {};
      this.theVisible = false;
    },
  },
};
</script>
