<template>
    <div class="app-container">
        <div>
            <el-card>
                <el-form
                    ref="form"
                    :inline="true"
                    :model="searchTermsValue"
                    class="demo-form-inline"
                >
                    <el-row :gutter="20">
                        <el-col :span="6">
                            <el-form-item label="短链接">
                                <el-input
                                    v-model="searchTermsValue.shortKey"
                                    clearable
                                    placeholder="请输入"
                                ></el-input>
                            </el-form-item>
                        </el-col>
                        <el-col :span="6">
                            <el-form-item label="原始链接:">
                                <el-input
                                    v-model="searchTermsValue.fullUrl"
                                    clearable
                                    placeholder="请输入"
                                ></el-input>
                            </el-form-item>
                        </el-col>
                        <el-col :span="6">
                            <el-form-item label="冻结:">
                                <el-select
                                    v-model="searchTermsValue.isFrozen"
                                    clearable
                                    placeholder="请选择"
                                    style="width: 100%"
                                >
                                    <el-option
                                        label="请选择"
                                        value=""
                                    ></el-option>
                                    <el-option
                                        :value="1"
                                        label="是"
                                    ></el-option>
                                    <el-option
                                        :value="-1"
                                        label="否"
                                    ></el-option>
                                </el-select>
                            </el-form-item>
                        </el-col>

                        <el-col :span="24">
                            <el-form-item label="创建时间:">
                                <el-date-picker
                                    v-model="searchTermsValue.temporalInterval"
                                    :picker-options="pickerOptions"
                                    align="right"
                                    end-placeholder="结束日期"
                                    range-separator="至"
                                    start-placeholder="开始日期"
                                    type="datetimerange"
                                ></el-date-picker>
                                <el-button
                                    style="
                                        display: inline-block;
                                        margin-left: 20px;
                                    "
                                    type="primary"
                                    @click="filterData"
                                    >查询</el-button
                                >
                            </el-form-item>
                        </el-col>
                    </el-row>
                </el-form>
            </el-card>
            <el-card style="margin-top: 20px">
                <div class="header">
                    <el-button type="primary" @click="operation()">
                        新增
                    </el-button>
                    <div class="header-r">
                        <el-button type="primary" @click="bulkOperation">
                            批量删除
                        </el-button>
                        <el-radio-group style="margin-left: 10px">
                            <el-button-group>
                                <el-button
                                    type="primary"
                                    @click="bulkOperation(1)"
                                    >批量冻结</el-button
                                >
                                <el-button
                                    type="primary"
                                    @click="bulkOperation(0)"
                                    >批量解冻</el-button
                                >
                            </el-button-group>
                        </el-radio-group>
                    </div>
                </div>
                <div class="table-data">
                    <el-table
                        ref="multipleTable"
                        :data="tableData"
                        border
                        style="width: 100%"
                        @selection-change="handleSelectionChange"
                    >
                        <el-table-column
                            type="selection"
                            width="55"
                        ></el-table-column>
                        <el-table-column
                            label="id"
                            min-width="80"
                            prop="id"
                        ></el-table-column>
                        <el-table-column
                            label="短链接"
                            min-width="70"
                            prop="shortKey"
                        ></el-table-column>
                        <el-table-column
                            label="原始链接"
                            min-width="300"
                            prop="fullUrl"
                        ></el-table-column>
                        <el-table-column
                            label="是否冻结"
                            min-width="80"
                            prop="isFrozen"
                        ></el-table-column>
                        <el-table-column
                            label="过期时间"
                            min-width="150"
                            prop="expirationTime"
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
                                    @click="operation(scope.row)"
                                    >编辑</el-button
                                >
                                <el-button
                                    size="small"
                                    type="text"
                                    @click="
                                        deleteOrFreezeValue(
                                            scope.row.id,
                                            '冻结'
                                        )
                                    "
                                    >冻结</el-button
                                >
                                <el-button
                                    size="small"
                                    type="text"
                                    @click="
                                        deleteOrFreezeValue(
                                            scope.row.id,
                                            '删除'
                                        )
                                    "
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
            </el-card>

            <urlModal ref="urlModal" @getTableData="getTableData"></urlModal>
        </div>
    </div>
</template>
<script>
import {
    getShortUrlList,
    batchFreezeArr,
    batchDeleteArr,
    deleteUrl,
    freezeUrl
} from '@/api/request-data.js';
import { dateFormat, todDateFormat } from '@/utils/date-format.js';
import urlModal from './modules/urlModal.vue';
export default {
    components: { urlModal },
    data() {
        return {
            // 筛选参数
            searchTermsValue: {
                shortKey: '',
                fullUrl: '',
                isFrozen: '',
                temporalInterval: null
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
                            start.setTime(
                                start.getTime() - 3600 * 1000 * 24 * 7
                            );
                            picker.$emit('pick', [start, end]);
                        }
                    },
                    {
                        text: '最近一个月',
                        onClick(picker) {
                            const end = new Date();
                            const start = new Date();
                            start.setTime(
                                start.getTime() - 3600 * 1000 * 24 * 30
                            );
                            picker.$emit('pick', [start, end]);
                        }
                    },
                    {
                        text: '最近三个月',
                        onClick(picker) {
                            const end = new Date();
                            const start = new Date();
                            start.setTime(
                                start.getTime() - 3600 * 1000 * 24 * 90
                            );
                            picker.$emit('pick', [start, end]);
                        }
                    }
                ]
            },
            // 短链接列表
            tableData: [],
            // 新增/修改表单

            // 时间范围规则时都开启
            expirationTimeRules: false,

            // 多选数组
            multipleSelection: [],
            amendOrAdd: '新增',
            //列表总数
            total: 0,
            // 每页显示条数
            pageSize: 10,
            // 当前页
            pageNum: 1
        };
    },
    mounted() {
        // 获取短链列表数据
        this.getTableData();
    },
    watch: {
        // 新增多选框内时间范围的表单验证是否启用
        perpetual(newV) {
            this.expirationTimeRules = newV === 0 ? true : false;
        }
    },
    methods: {
        // 获取筛选条件
        filterData() {
            this.pageNum = 1;
            this.filterModel = {};
            if (this.searchTermsValue.shortKey !== '') {
                this.filterModel.shortKey = this.searchTermsValue.shortKey;
            }
            if (this.searchTermsValue.fullUrl !== '') {
                this.filterModel.fullUrl = this.searchTermsValue.fullUrl;
            }
            if (this.searchTermsValue.temporalInterval !== null) {
                this.filterModel.createTimeL = parseInt(
                    this.searchTermsValue.temporalInterval[0].getTime() / 1000
                );
                this.filterModel.createTimeR = parseInt(
                    this.searchTermsValue.temporalInterval[1].getTime() / 1000
                );
            }
            this.filterModel.isFrozen = this.searchTermsValue.isFrozen;
            this.getTableData();
        },
        // 请求数据
        async getTableData(page) {
            this.multipleSelection = [];
            const params = { ...this.filterModel };
            params.size = this.pageSize;
            params.page = page ? page : this.pageNum;
            this.pageNum = page ? page : this.pageNum;
            this.loading = true;
            console.log(params);
            const { data } = await getShortUrlList(params);
            this.loading = false;

            this.total = data.len;

            this.tableData = data.list.map((v) => {
                return {
                    createTime: dateFormat(v.createTime),
                    expirationTime:
                        v.expirationTime === 0
                            ? '永久'
                            : dateFormat(v.expirationTime),
                    fullUrl: v.fullUrl,
                    id: v.id,
                    isFrozen: v.isFrozen ? '是' : '否',
                    shortKey: v.shortKey,
                    updateTime: dateFormat(v.updateTime)
                };
            });
            //   this.$refs.content.scrollTo({ top: 0, behavior: "smooth" });
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
            console.log(...this.multipleSelection, 111);
        },
        // 获取备选id
        handleSelectionChange(val) {
            this.multipleSelection = val.map((v) => {
                return v.id;
            });
        },
        // 弹窗事件
        operation(data) {
            this.$refs.urlModal.show(data);
        },
        // 批量冻结/解冻/删除 事件
        async bulkOperation(v) {
            if (this.multipleSelection.length > 0) {
                let operate;
                switch (v) {
                    case 1:
                        operate = '冻结';
                        break;
                    case 0:
                        operate = '解冻';
                        break;
                    default:
                        operate = '删除';
                        break;
                }
                this.$confirm(
                    `此操作将批量${operate}该链接, 是否继续?`,
                    '提示',
                    {
                        confirmButtonText: '确定',
                        cancelButtonText: '取消',
                        type: 'warning'
                    }
                )
                    .then(async () => {
                        let requestData;
                        try {
                            if (v === 0 || v === 1) {
                                const params = {
                                    ids: this.multipleSelection,
                                    isFrozen: v
                                };
                                console.log(params);
                                requestData = await batchFreezeArr(params);
                            } else {
                                const params = { ids: this.multipleSelection };
                                requestData = await batchDeleteArr(params);
                            }
                            console.log(requestData.code);
                            if (requestData.code === 200) {
                                console.log(11);
                                this.$message({
                                    type: 'success',
                                    message: `批量${operate}成功!`
                                });
                                this.multipleSelection = [];
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
                            message: `已取消${operate}`
                        });
                    });
            } else {
                this.$message({
                    type: 'info ',
                    message: `无勾选内容!`
                });
            }
        },

        async deleteOrFreezeValue(id, operate) {
            this.$confirm(`此操作将${operate}该链接, 是否继续?`, '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            })
                .then(async () => {
                    let requestData;
                    try {
                        if (operate === '删除') {
                            requestData = await deleteUrl(id);
                        } else if (operate === '冻结') {
                            requestData = await freezeUrl(id);
                        }
                        console.log(requestData.code);
                        if (requestData.code === 200) {
                            console.log(11);
                            this.$message({
                                type: 'success',
                                message: `${operate}成功!`
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
                        message: `已取消${operate}`
                    });
                });
        }
    }
};
</script>

<style scoped>
.table-data {
    margin-top: 20px;
}
.header {
    width: 100%;
    overflow: hidden;
}
.header-r {
    float: right;
}
.paging {
    margin-top: 20px;
}
</style>
