import Vue from 'vue'
import VueRouter from 'vue-router'
import shortConnection from '../views/shortConnection.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/short-connection',
    name: 'shortConnection',
    component: shortConnection
  },
  {
    path: '/black-list',
    name: 'BlackList',
    component: () => import(/* webpackChunkName: "black-list" */ '../views/blackList.vue')
  },
  {
    path: '/',
    redirect:'/short-connection'
  },
]

const router = new VueRouter({
  routes
})

export default router
