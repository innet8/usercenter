import { $t } from "@/lang/index"
import Layout from '@/layout/index.vue'
import Login from '@/views/login.vue'
import Register from '@/views/register.vue'
import Success from '@/views/success.vue'

export const routes = [
    {
        name: "login",
        path: "/login",
        meta: { title: $t('登陆'), login: false },
        component: Login
    },
    {
        name: "register",
        path: "/register",
        meta: { title: $t('注册'), login: false },
        component: Register
    },
    {
        name: "success",
        path: "/success",
        meta: { title: $t('成功'), login: false },
        component: Success
    },
    {
        name: "/",
        path: "/",
        redirect: "/home",
        meta: { title: $t('authentik-go') },
        component: Layout,
        children: [
            {
                name: "home",
                meta: { title: $t('authentik-go') },
                path: 'home',
                component: () => import('@/views/index.vue'),
            },
        ],
    }
]
