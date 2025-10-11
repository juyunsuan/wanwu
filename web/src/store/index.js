import Vue from 'vue'
import Vuex from 'vuex'
import VuexPersistence from 'vuex-persist'
import { login } from './module/login'
import { user } from './module/user'
import { app } from './module/app'
import { workflow } from './module/workflow'
import knowledge from './module/knowledge'


Vue.use(Vuex)
// 用户信息持久化
const vuexLocal = new VuexPersistence({
    key:'access_cert',
    storage: window.localStorage,
    modules: ['user']
})

export const store = new Vuex.Store({
    modules: {
        login,
        user,
        app,
        workflow,
        knowledge
    },
    plugins: [vuexLocal.plugin]
})
