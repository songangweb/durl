<template>
    <div class="short-connection">
        <div class="content" ref="content">
            <div class="operate">
                <div class="grid-content bg-purple">
                    <el-button type="primary" size="medium" @click="operation('新增')">新增</el-button>
                </div>
            </div>
            <div class="search">
                <el-form :inline="true" :model="searchTermsValue" class="demo-form-inline">
                    <el-form-item label="ip:">
                        <el-input v-model="searchTermsValue.ip" placeholder="请输入" size="mini" clearable></el-input>
                    </el-form-item>
                    <el-form-item label="时间区间:">
                        <el-date-picker size="mini" v-model="searchTermsValue.temporalInterval" type="datetimerange" :picker-options="pickerOptions" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" align="right"></el-date-picker>
                    </el-form-item>
                    <el-form-item>
                        <el-button type="primary" @click="filterData" size="small" style="margin-left: 100px">查询</el-button>
                    </el-form-item>
                </el-form>
            </div>
            <div class="table-data">
                <el-table :data="tableData" border style="width: 100%">
                    <el-table-column prop="id" label="id" min-width="80"></el-table-column>
                    <el-table-column prop="ip" label="ip" min-width="70"></el-table-column>
                    <el-table-column prop="createTime" label="创建时间" min-width="150"></el-table-column>

                    <el-table-column prop="updateTime" label="修改时间" min-width="150"></el-table-column>
                    <el-table-column fixed="right" label="操作" width="120">
                        <template slot-scope="scope">
                            <el-button type="text" size="small" @click="operation('编辑', scope.row)">编辑</el-button>
                            <el-button @click="deleteValue(scope.row.id)" type="text" size="small">删除</el-button>
                        </template>
                    </el-table-column>
                    <template slot="empty">
                        <div class="no-data">
                            <span>暂无数据</span>
                        </div>
                    </template>
                </el-table>
                <div class="paging">
                    <el-pagination :current-page.sync="pageNum" :page-sizes="[10, 20, 30, 40]" :page-size.sync="pageSize" :total.sync="total" layout="total, sizes, prev, pager, next, jumper" background @current-change="onPageNumChange" @size-change="onPageSizeChange" class="table-pagination"></el-pagination>
                </div>
            </div>
            <el-dialog :title="amendOrAdd" width="520px" class="dialog" :visible.sync="theVisible" :before-close="clearData">
                <el-form :model="formValue" :rules="formRules" label-width="100px" class="demo-addunit">
                    <el-form-item label="ip" class="demo-input" prop="ip">
                        <el-input v-model="formValue.ip"></el-input>
                    </el-form-item>
                </el-form>
                <span slot="footer" class="dialog-footer">
                    <el-button @click="clearData">取消</el-button>
                    <el-button class="save-btn" @click="sendParams">保存</el-button>
                </span>
            </el-dialog>
        </div>
    </div>
</template>
<script>
import { getBlackListArr, addBlackValue, changeBlackValue, deleteBlackValue } from '@/api/request-data.js';
import { dateFormat } from '@/utils/date-format.js';
export default {
    data() {
        return {
            // 筛选参数
            searchTermsValue: {
                ip: '',
                temporalInterval: null,
            },
            filterModel: {},
            // 时间区间快捷选择框
            pickerOptions: {
                shortcuts: [
                    {
                        text: '最近一周',
                        onClick(picker) {
                            const end = new Date();
                            const start = new Date();
                            start.setTime(start.getTime() - 3600 * 1000 * 24 * 7);
                            picker.$emit('pick', [start, end]);
                        },
                    },
                    {
                        text: '最近一个月',
                        onClick(picker) {
                            const end = new Date();
                            const start = new Date();
                            start.setTime(start.getTime() - 3600 * 1000 * 24 * 30);
                            picker.$emit('pick', [start, end]);
                        },
                    },
                    {
                        text: '最近三个月',
                        onClick(picker) {
                            const end = new Date();
                            const start = new Date();
                            start.setTime(start.getTime() - 3600 * 1000 * 24 * 90);
                            picker.$emit('pick', [start, end]);
                        },
                    },
                ],
            },
            // 黑名单列表
            tableData: [],
            // 新增/修改表单
            formValue: {
                ip: '',
                id: 0,
            },
            // 时间范围规则时都开启
            expirationTimeRules: false,
            // 新增/修改
            formRules: {
                ip: [{ required: true, message: '请输入ip地址', trigger: 'blur' }],
            },
            // 弹框显示
            theVisible: false,
            amendOrAdd: '新增',
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
            if (this.searchTermsValue.ip !== '') {
                this.filterModel.ip = this.searchTermsValue.ip;
            }
            if (this.searchTermsValue.temporalInterval !== null) {
                this.filterModel.createTimeL = parseInt(this.searchTermsValue.temporalInterval[0].getTime() / 1000);
                this.filterModel.createTimeR = parseInt(this.searchTermsValue.temporalInterval[1].getTime() / 1000);
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
            if (data.len) {
                this.total = data.len;
            }
            this.tableData = data.list.map(v => {
                return {
                    id: v.id,
                    ip: v.ip,
                    createTime: dateFormat(v.createTime),
                    updateTime: dateFormat(v.updateTime),
                };
            });
            this.$refs.content.scrollTo({ top: 0, behavior: 'smooth' });
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
            this.amendOrAdd = value + 'ip';
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
                ip: '',
            };
            this.theVisible = false;
        },
        // 发送新增/修改参数加入新短链/修改短链接
        async sendParams() {
            const params = { ip: this.formValue.ip };
            if (this.formValue.id === 0) {
                console.log(params);
                const requestData = await addBlackValue(params);
                if (requestData.code === 200) {
                    this.pageNum = 1;
                    this.$message({
                        type: 'success',
                        message: `新增成功!`,
                    });
                    this.getTableData();
                } else {
                    console.log('新增失败');
                }
            } else {
                console.log(params, this.formValue.id);
                const requestData = await changeBlackValue(params, this.formValue.id);
                if (requestData.code === 200) {
                    this.$message({
                        type: 'success',
                        message: `修改成功!`,
                    });
                    this.getTableData();
                } else {
                    console.log('修改失败');
                }
            }
            this.clearData();
        },
        async deleteValue(id) {
            this.$confirm(`此操作将永久删除该ip, 是否继续?`, '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning',
            })
                .then(async () => {
                    let requestData;
                    try {
                        requestData = await deleteBlackValue(id);
                        console.log(requestData.code);
                        if (requestData.code === 200) {
                            console.log(11);
                            this.$message({
                                type: 'success',
                                message: `删除成功!`,
                            });
                            this.getTableData();
                        }
                    } catch (error) {
                        console.log('请求失败');
                        console.log(error);
                    }
                })
                .catch(() => {
                    this.$message({
                        type: 'info',
                        message: `已取消删除`,
                    });
                });
        },
    },
};
</script>
<style lang="scss" scoped>
::v-deep.short-connection {
    > .content {
        width: calc(100vw - 230px);
        vertical-align: top;
        height: calc(100vh - 95px);
        box-sizing: border-box;
        padding: 20px;
        overflow-y: scroll;
        > .operate {
            margin-bottom: 20px;
        }
        > .search {
            margin-bottom: 20px;
        }
        .no-data {
            position: relative;
            background-image: url('../assets/no-data.webp');
            background-repeat: no-repeat;
            background-size: 100% 100%;
            width: 260px;
            height: 180px;
            margin: 50px auto;
            & > span {
                font-size: 14px;
                font-weight: 400;
                font-weight: bold;
                color: #333333;
                position: absolute;
                bottom: -55px;
                left: 55%;
                transform: translateX(-50%);
            }
        }
        .paging {
            display: flex;
            align-items: center;
            .el-pagination {
                display: block;
                margin: 30px auto;
                margin-bottom: 20px;
            }
        }
        .dialog {
            .el-dialog__header {
                text-align: center;
                padding: 17.5px 20px !important;
                border-bottom: 1px solid #dfe8f0;
                font-weight: 600;
            }
            .el-input__inner {
                width: 255px;
                height: 40px;
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
            }
            .demo-input {
                width: 380px;
                margin-left: 60px;
            }
            .el-dialog__footer {
                text-align: center;
                .el-button {
                    width: 163px;
                }
                .save-btn {
                    width: 163px;
                    background: #36c9a4;
                    border-color: #36c9a4;
                    color: #ffffff;
                }
            }
        }
    }
}
</style>
