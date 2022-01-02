<template>
    <div class="short-connection">
        <div class="content" ref="content">
            <div class="operate">
                <el-row :gutter="24">
                    <el-col :span="15">
                        <div class="grid-content bg-purple">
                            <el-button type="primary" size="medium" @click="operation('新增')">新增</el-button>
                        </div>
                    </el-col>
                    <el-col :span="5">
                        <!-- <div class="grid-content bg-purple">
                            <el-button type="primary" size="medium">批量冻结</el-button>
                        </div> -->
                        <el-radio-group>
                            <el-button-group>
                                <el-button type="primary" size="medium" @click="bulkOperation(1)">批量冻结</el-button>
                                <el-button type="primary" size="medium" @click="bulkOperation(0)">批量解冻</el-button>
                            </el-button-group>
                        </el-radio-group>
                    </el-col>
                    <el-col :span="2">
                        <div class="grid-content bg-purple">
                            <el-button type="primary" size="medium" @click="bulkOperation">批量删除</el-button>
                        </div>
                    </el-col>
                </el-row>
            </div>
            <div class="search">
                <el-form :inline="true" :model="searchTermsValue" class="demo-form-inline">
                    <el-form-item label="短链接:">
                        <el-input v-model="searchTermsValue.shortKey" placeholder="请输入" size="mini" clearable></el-input>
                    </el-form-item>
                    <el-form-item label="原始链接:">
                        <el-input v-model="searchTermsValue.fullUrl" placeholder="请输入" size="mini" clearable></el-input>
                    </el-form-item>
                    <el-form-item label="冻结:">
                        <el-select v-model="searchTermsValue.isFrozen" placeholder="请选择" size="mini" clearable>
                            <el-option label="是" :value="1"></el-option>
                            <el-option label="否" :value="0"></el-option>
                        </el-select>
                    </el-form-item>
                    <br />
                    <el-form-item label="时间区间:">
                        <el-date-picker size="mini" v-model="searchTermsValue.temporalInterval" type="datetimerange" :picker-options="pickerOptions" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" align="right"></el-date-picker>
                    </el-form-item>
                    <el-form-item>
                        <el-button type="primary" @click="filterData" size="small" style="margin-left: 100px">查询</el-button>
                    </el-form-item>
                </el-form>
            </div>
            <div class="table-data">
                <el-table :data="tableData" border style="width: 100%" @selection-change="handleSelectionChange" ref="multipleTable">
                    <el-table-column type="selection" width="55"></el-table-column>
                    <el-table-column prop="id" label="id" min-width="80"></el-table-column>
                    <el-table-column prop="shortKey" label="短链接" min-width="70"></el-table-column>
                    <el-table-column prop="fullUrl" label="原始链接" min-width="300"></el-table-column>
                    <el-table-column prop="isFrozen" label="是否冻结" min-width="80"></el-table-column>
                    <el-table-column prop="expirationTime" label="过期时间" min-width="150"></el-table-column>
                    <el-table-column prop="createTime" label="创建时间" min-width="150"></el-table-column>

                    <el-table-column prop="updateTime" label="修改时间" min-width="150"></el-table-column>
                    <el-table-column fixed="right" label="操作" width="120">
                        <template slot-scope="scope">
                            <el-button type="text" size="small" @click="operation('编辑', scope.row)">编辑</el-button>
                            <el-button type="text" size="small">冻结</el-button>
                            <el-button @click="theDeleteValue(scope.row.id)" type="text" size="small">删除</el-button>
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
                    <el-form-item :label="'短链接:'" class="demo-input" v-if="amendOrAdd === '编辑短链接'">
                        <el-input v-model="formValue.shortKey" placeholder="请输入" :disabled="true"></el-input>
                    </el-form-item>
                    <el-form-item :label="'原始链接'" class="demo-input" prop="fullUrl">
                        <el-input v-model="formValue.fullUrl"></el-input>
                    </el-form-item>
                    <el-form-item label="冻结:" class="demo-input" prop="isFrozen">
                        <el-radio-group v-model="formValue.isFrozen">
                            <el-radio :label="1">是</el-radio>
                            <el-radio :label="0">否</el-radio>
                        </el-radio-group>
                    </el-form-item>
                    <el-form-item label="是否永久:" class="demo-input" prop="perpetual">
                        <el-radio-group v-model="formValue.perpetual">
                            <el-radio :label="1">是</el-radio>
                            <el-radio :label="0">否</el-radio>
                        </el-radio-group>
                    </el-form-item>
                    <el-form-item label="过期时间:" class="demo-input" v-if="formValue.perpetual === 0 ? true : false" prop="expirationTime">
                        <el-date-picker v-model="formValue.expirationTime" type="datetime" placeholder="选择日期时间"></el-date-picker>
                    </el-form-item>
                </el-form>
                <span slot="footer" class="dialog-footer">
                    <el-button>取消</el-button>
                    <el-button class="save-btn" @click="sendParams">保存</el-button>
                </span>
            </el-dialog>
        </div>
    </div>
</template>
<script>
import { getShortChainArr, addShortChainValue, changeShortChainValue, batchFreezeArr, batchDeleteArr, deleteValue } from '@/api/request-data.js';
import { dateFormat, todDateFormat } from '@/utils/date-format.js';
export default {
    data() {
        return {
            // 筛选参数
            searchTermsValue: {
                shortKey: '',
                fullUrl: '',
                isFrozen: '',
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
            // 短链接列表
            tableData: [],
            // 新增/修改表单
            formValue: {
                fullUrl: '',
                shortKey: '',
                isFrozen: 0,
                perpetual: 1,
                expirationTime: {},
                id: 0,
            },
            expirationTimeRules: false,
            formRules: {
                fullUrl: [{ required: this.amendOrAdd === '新增短链接', message: '请输入链接地址', trigger: 'blur' }],
                isFrozen: [{ required: true, message: '请选择是否冻结', trigger: 'blur' }],
                perpetual: [{ required: true, message: '请选择是否永久', trigger: 'blur' }],
                expirationTime: [{ required: true, message: '请选择过期时间', trigger: 'blur' }],
            },
            // 弹框显示
            theVisible: false,
            multipleSelection: [],
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
    watch: {
        perpetual(newV) {
            this.expirationTimeRules = newV === 0 ? true : false;
        },
    },
    methods: {
        // 获取筛选条件
        filterData() {
            this.pageNum = 1;
            this.filterModel = {};
            // if (this.searchTermsValue.shortKey !== '') {
            //     this.filterModel.shortKey = this.searchTermsValue.shortKey;
            // }
            if (this.searchTermsValue.fullUrl !== '') {
                this.filterModel.url = this.searchTermsValue.fullUrl;
            }
            // if (this.searchTermsValue.isFrozen !== '') {
            //     this.filterModel.isFrozen = this.searchTermsValue.isFrozen;
            // }
            if (this.searchTermsValue.temporalInterval !== null) {
                this.filterModel.createTimeL = parseInt(this.searchTermsValue.temporalInterval[0].getTime() / 1000);
                this.filterModel.createTimeR = parseInt(this.searchTermsValue.temporalInterval[1].getTime() / 1000);
            }
            this.getTableData();
        },
        // 请求数据
        async getTableData() {
            this.multipleSelection = [];
            const params = { ...this.filterModel };
            params.size = this.pageSize;
            params.page = this.pageNum;
            this.loading = true;
            console.log(params);
            const { data } = await getShortChainArr(params);
            this.loading = false;
            if (data.len) {
                this.total = data.len;
            }
            this.tableData = data.list.map(v => {
                return {
                    createTime: dateFormat(v.createTime),
                    expirationTime: v.expirationTime === 0 ? '永久' : dateFormat(v.expirationTime),
                    fullUrl: v.fullUrl,
                    id: v.id,
                    isFrozen: v.isFrozen ? '是' : '否',
                    shortKey: v.shortKey,
                    updateTime: dateFormat(v.updateTime),
                };
            });
            this.$refs.content.scrollTo({ top: 0, behavior: 'smooth' });
            // console.log(this.tableData);
        },
        //监听页数
        onPageNumChange(v) {
            this.pageNum = v;
            this.getTableData();
        },
        // 监听每页显示的条数
        onPageSizeChange(v) {
            this.pageSize = v;
            this.getTableData();
            console.log(...this.multipleSelection, 111);
        },
        // 获取备选id
        handleSelectionChange(val) {
            this.multipleSelection = val.map(v => {
                return v.id;
            });
        },
        // 弹窗事件
        operation(value, msg) {
            this.amendOrAdd = value + '短链接';
            console.log(msg);
            if (msg) {
                this.formValue = {
                    shortKey: msg.shortKey,
                    fullUrl: msg.fullUrl,
                    isFrozen: msg.isFrozen === '是' ? 1 : 0,
                    perpetual: msg.expirationTime === '永久' ? 1 : 0,
                    id: msg.id,
                };
                if (msg.expirationTime !== '永久') {
                    this.formValue.perpetual = todDateFormat(msg.expirationTime)._d.getTime();
                }
            }
            this.theVisible = true;
        },
        // 批量冻结/解冻/删除 事件
        async bulkOperation(v) {
            if (this.multipleSelection.length > 0) {
                if (v === 0 || v === 1) {
                    const params = { ids: this.multipleSelection, isFrozen: v };
                    console.log(params);
                    try {
                        const requestData = await batchFreezeArr(params);
                        if (requestData.code === 200) {
                            this.getTableData();
                        }
                    } catch (error) {
                        console.log('请求失败');
                    }
                } else {
                    const params = { ids: this.multipleSelection };
                    console.log(params, 111);
                    try {
                        const requestData = await batchDeleteArr(params);
                        if (requestData.code === 200) {
                            this.getTableData();
                        }
                    } catch (error) {
                        console.log('请求失败');
                    }
                }
            }
        },
        clearData() {
            this.formValue = {
                fullUrl: '',
                shortKey: '',
                isFrozen: 0,
                perpetual: 1,
                id: 0,
                expirationTime: {},
            };
            this.theVisible = false;
        },
        // 发送新增/修改参数加入新短链/修改短链接
        async sendParams() {
            const params = {};
            params.fullUrl = this.formValue.fullUrl;
            params.isFrozen = this.formValue.isFrozen;
            try {
                params.expirationTime = this.formValue.expirationTime.getTime() / 1000;
            } catch (error) {
                params.expirationTime = 0;
            }
            if (this.formValue.id === 0) {
                console.log(params);
                const requestData = await addShortChainValue(params);
                if (requestData.code === 200) {
                    this.pageNum = 1;
                    this.getTableData();
                } else {
                    console.log('新增失败');
                }
            } else {
                console.log(params, this.formValue.id);
                const requestData = await changeShortChainValue(params, this.formValue.id);
                if (requestData.code === 200) {
                    this.getTableData();
                } else {
                    console.log('修改失败');
                }
            }
            this.formValue = {
                fullUrl: '',
                shortKey: '',
                isFrozen: 0,
                perpetual: 1,
                expirationTime: {},
            };
            this.theVisible = false;
        },
        async theDeleteValue(id) {
            const requestData = await deleteValue(id);
            if (requestData.code === 200) {
                console.log(11);
                this.getTableData();
            }
            // try {

            // } catch (error) {
            //     console.log('请求失败...');
            // }
        },
    },
};
</script>
<style lang="scss" scoped>
::v-deep.short-connection {
    > .content {
        width: calc(100vw - 230px);
        vertical-align: top;
        height: calc(100vh - 90px);
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
