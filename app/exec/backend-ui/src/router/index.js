import Vue from 'vue';
import VueRouter from 'vue-router';
import shortConnection from '../views/shortConnection.vue';

Vue.use(VueRouter);

// 防止路由重复点击报错

const originalPush = VueRouter.prototype.push;

VueRouter.prototype.push = function push(location) {
    return originalPush.call(this, location).catch(err => err);
};
const routes = [
    {
        path: '/short-connection',
        name: 'shortConnection',
        component: shortConnection,
    },
    {
        path: '/black-list',
        name: 'BlackList',
        component: () => import(/* webpackChunkName: "black-list" */ '../views/blackList.vue'),
    },
    {
        path: '/',
        redirect: '/short-connection',
    },
];

const router = new VueRouter({
    routes,
});

export default router;
