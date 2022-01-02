<template>
    <div id="app">
        <div class="content">
            <div class="top-logo">
                <span>Durl</span>
            </div>
            <div class="sidebar">
                <div class="inner">
                    <el-row class="tac">
                        <el-col :span="24">
                            <el-menu default-active="1" class="el-menu-vertical-demo">
                                <el-menu-item index="1" @click="$router.push('/short-connection')">
                                    <i class="el-icon-share"></i>
                                    <span slot="title">短链接</span>
                                </el-menu-item>
                                <el-menu-item index="2" @click="$router.push('/black-list')">
                                    <i class="el-icon-document"></i>
                                    <span slot="title">黑名单</span>
                                </el-menu-item>
                            </el-menu>
                        </el-col>
                    </el-row>
                </div>
            </div>
            <div class="main-content">
                <div class="inner">
                    <router-view></router-view>
                </div>
            </div>
        </div>
    </div>
</template>
<style lang="scss" scoped>
::v-deep .content {
    height: 100vh;
    min-width: 1220px;
    .top-logo {
        height: 70px;
        border-bottom: 1px solid rgb(209, 206, 206);
        > span {
            font-size: 45px;
            line-height: 70px;
            font-weight: 700;
            padding-left: 30px;
        }
    }
    .sidebar {
        display: inline-block;
        height: calc(100vh - 71px);
        width: 190px;
        overflow: hidden;
        background-color: #ffffff;
        > .inner {
            width: 200px;
            overflow: hidden;
        }
    }
    .main-content {
        display: inline-block;
        // border: 1px solid red;

        margin: 10px;
    }
}
</style>
