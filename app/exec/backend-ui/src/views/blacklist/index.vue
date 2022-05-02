<template>
  <div class="app-container">
    <div ref="content" class="content">
      <el-card>
        <el-form
          :inline="true"
          :model="searchTermsValue"
          class="demo-form-inline"
        >
          <el-row :gutter="20">
            <el-col :span="6">
              <el-form-item label="ip:">
                <el-input
                  v-model="searchTermsValue.ip"
                  clearable
                  placeholder="请输入"
                ></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="创建时间:">
                <el-date-picker
                  v-model="searchTermsValue.temporalInterval"
                  :picker-options="pickerOptions"
                  align="right"
                  end-placeholder="结束日期"
                  range-separator="至"
                  start-placeholder="开始日期"
                  style="width: 100%"
                  type="datetimerange"
                ></el-date-picker>
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-button type="primary" @click="filterData">查询</el-button>
            </el-col>
          </el-row>
        </el-form>
      </el-card>

      <el-card class="grid-content">
        <div class="bg-purple">
          <el-button size="medium" type="primary" @click="operation('新增')"
            >新增</el-button
          >
          <div class="table-data">
            <el-table :data="tableData" border style="width: 100%">
              <el-table-column
                label="id"
                min-width="80"
                prop="id"
              ></el-table-column>
              <el-table-column
                label="ip"
                min-width="70"
                prop="ip"
              ></el-table-column>
              <el-table-column
                label="创建时间"
                min-width="150"
                prop="createTime"
              ></el-table-column>

              <el-table-column
                label="修改时间"
                min-width="150"
                prop="updateTime"
              ></el-table-column>
              <el-table-column fixed="right" label="操作" width="120">
                <template slot-scope="scope">
                  <el-button
                    size="small"
                    type="text"
                    @click="operation('编辑', scope.row)"
                    >编辑</el-button
                  >
                  <el-button
                    size="small"
                    type="text"
                    @click="deleteValue(scope.row.id)"
                    >删除</el-button
                  >
                </template>
              </el-table-column>
              <template slot="empty">
                <div class="no-data">
                  <span>暂无数据</span>
                </div>
              </template>
            </el-table>
            <div class="paging">
              <el-pagination
                :current-page.sync="pageNum"
                :page-size.sync="pageSize"
                :page-sizes="[10, 20, 30, 40]"
                :total.sync="total"
                background
                class="table-pagination"
                layout="total, sizes, prev, pager, next, jumper"
                @current-change="onPageNumChange"
                @size-change="onPageSizeChange"
              ></el-pagination>
            </div>
          </div>
        </div>
      </el-card>

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
          <el-form-item class="demo-input" label="ip" prop="ip">
            <el-input v-model="formValue.ip"></el-input>
          </el-form-item>
        </el-form>
        <span slot="footer" class="dialog-footer">
          <el-button @click="clearData">取消</el-button>
          <el-button class="save-btn" @click="sendParams('formValue')"
            >保存</el-button
          >
        </span>
      </el-dialog>
    </div>
  </div>
</template>
<script>
import {
  getBlackListArr,
  addBlacklist,
  changeBlacklist,
  deleteBlacklist,
} from "@/api/request-data.js";
import { dateFormat } from "@/utils/date-format.js";
export default {
  data() {
    return {
      // 筛选参数
      searchTermsValue: {
        ip: "",
        temporalInterval: null,
      },
      filterModel: {},
      // 时间区间快捷选择框
      pickerOptions: {
        shortcuts: [
          {
            text: "最近一周",
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 7);
              picker.$emit("pick", [start, end]);
            },
          },
          {
            text: "最近一个月",
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 30);
              picker.$emit("pick", [start, end]);
            },
          },
          {
            text: "最近三个月",
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 90);
              picker.$emit("pick", [start, end]);
            },
          },
        ],
      },
      // 黑名单列表
      tableData: [],
      // 新增/修改表单
      formValue: {
        ip: "",
        id: 0,
      },
      // 时间范围规则时都开启
      expirationTimeRules: false,
      // 新增/修改
      formRules: {
        ip: [{ required: true, message: "请输入ip地址", trigger: "blur" }],
      },
      // 弹框显示
      theVisible: false,
      amendOrAdd: "新增",
      //列表总数
      total: 0,
      // 每页显示条数
      pageSize: 10,
      // 当前页
      pageNum: 1,
    };
  },
  mounted() {
    // 获取短链列表数据
    this.getTableData();
  },
  methods: {
    // 获取筛选条件
    filterData() {
      this.pageNum = 1;
      this.filterModel = {};
      if (this.searchTermsValue.ip !== "") {
        this.filterModel.ip = this.searchTermsValue.ip;
      }
      if (this.searchTermsValue.temporalInterval !== null) {
        this.filterModel.createTimeL = parseInt(
          this.searchTermsValue.temporalInterval[0].getTime() / 1000
        );
        this.filterModel.createTimeR = parseInt(
          this.searchTermsValue.temporalInterval[1].getTime() / 1000
        );
      }
      this.getTableData();
    },
    // 请求数据
    async getTableData() {
      const params = { ...this.filterModel };
      params.size = this.pageSize;
      params.page = this.pageNum;
      this.loading = true;
      console.log(params);
      const { data } = await getBlackListArr(params);
      this.loading = false;

      this.total = data.len;

      this.tableData = data.list.map((v) => {
        return {
          id: v.id,
          ip: v.ip,
          createTime: dateFormat(v.createTime),
          updateTime: dateFormat(v.updateTime),
        };
      });
      this.$refs.content.scrollTo({ top: 0, behavior: "smooth" });
    },
    //监听页数改变发送请求刷新数据
    onPageNumChange(v) {
      this.pageNum = v;
      this.getTableData();
    },
    // 监听每页显示的条数改变发送请求刷新数据
    onPageSizeChange(v) {
      this.pageSize = v;
      this.getTableData();
    },
    // 弹窗事件
    operation(value, msg) {
      this.amendOrAdd = value + "ip";
      // console.log(msg);
      if (msg) {
        this.formValue = {
          ip: msg.ip,
          id: msg.id,
        };
      }
      this.theVisible = true;
    },
    clearData() {
      this.formValue = {
        ip: "",
      };
      this.theVisible = false;
    },
    // 发送新增/修改参数加入新短链/修改短链接
    async sendParams(formName) {
      this.$refs[formName].validate(async (valid) => {
        if (valid) {
          const params = { ip: this.formValue.ip };
          if (this.formValue.id === 0) {
            console.log(params);
            const requestData = await addBlacklist(params);
            if (requestData.code === 200) {
              this.pageNum = 1;
              this.$message({
                type: "success",
                message: `新增成功!`,
              });
              this.getTableData();
            } else {
              console.log("新增失败");
            }
          } else {
            console.log(params, this.formValue.id);
            const requestData = await changeBlacklist(
              params,
              this.formValue.id
            );
            if (requestData.code === 200) {
              this.$message({
                type: "success",
                message: `修改成功!`,
              });
              this.getTableData();
            } else {
              console.log("修改失败");
            }
          }
          this.clearData();
        }
      });
    },
    async deleteValue(id) {
      this.$confirm(`此操作将永久删除该ip, 是否继续?`, "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      })
        .then(async () => {
          let requestData;
          try {
            requestData = await deleteBlacklist(id);
            console.log(requestData.code);
            if (requestData.code === 200) {
              console.log(11);
              this.$message({
                type: "success",
                message: `删除成功!`,
              });
              this.getTableData();
            }
          } catch (error) {
            console.log("请求失败");
            console.log(error);
          }
        })
        .catch(() => {
          this.$message({
            type: "info",
            message: `已取消删除`,
          });
        });
    },
  },
};
</script>
<style scoped>
.grid-content {
  margin-top: 20px;
}
.table-data {
  margin-top: 20px;
}
.paging {
  margin-top: 20px;
}
</style>
