<template>
    <div class="page-login child-view">
        <div class="login-body">
            <div class="login-box">
                <h2 class="login-title">
                    <span>{{ loginMode == "qrcode" ? $t("authentik-go") : $t("authentik-go") }}</span>
                </h2>
                <p class="login-subtitle">
                    {{ $t("输入您的凭证以访问您的帐户") }}
                </p>
                <transition name="login-mode">
                    <n-form :rules="rules" label-placement="left">
                        <div v-if="loginMode == 'access'" class="login-access">
                            <n-form-item label="" path="email">
                                <n-input v-model:value="formData.email" @blur="onBlur" :placeholder="$t('输入您的电子邮箱')" clearable
                                    size="large">
                                    <template #prefix>
                                        <n-icon :component="MailOutline" />
                                    </template>
                                </n-input>
                            </n-form-item>
                            <n-form-item label="" path="password">
                                <n-input v-model:value="formData.password" @blur="onBlur" :placeholder="$t('输入您的密码')" clearable
                                    size="large">
                                    <template #prefix>
                                        <n-icon :component="LockClosedOutline" />
                                    </template>
                                </n-input>
                            </n-form-item>
                            <n-form-item label="" path="code" v-if="codeNeed">
                                <n-input class="code-load-input" v-model:value="code" :placeholder="$t('输入图形验证码')" clearable size="large">
                                    <template #prefix>
                                        <n-icon :component="CheckmarkCircleOutline" />
                                    </template>
                                    <template #suffix>
                                        <div class="login-code-end" @click="refreshCode">
                                            <div v-if="codeLoad > 0" class="code-load">
                                                <Loading />
                                            </div>
                                            <span v-else-if="codeUrl === 'error'" class="code-error">{{ $t("加载失败") }}</span>
                                            <img v-else :src="codeUrl" />
                                        </div>
                                    </template>
                                </n-input>
                            </n-form-item>
                            <n-form-item label="" path="confirmPassword" v-if="loginType == 'reg'">
                                <n-input type="password" v-model:value="formData.confirmPassword"
                                    :placeholder="$t('输入确认密码')" clearable size="large">
                                    <template #prefix>
                                        <n-icon :component="LockClosedOutline" />
                                    </template>
                                </n-input>
                            </n-form-item>
                            <n-button v-if="loginType == 'login'" :loading="loadIng" @click="handleLogin" type="primary"
                                size="large">{{ $t("登录") }}</n-button>
                            <n-button v-else type="primary" :loading="loadIng" @click="handleReg">{{ $t("注册") }}</n-button>
                            <div class="login-switch">
                                <template v-if="loginType == 'login'">
                                    {{ $t("还没有帐号？") }}
                                    <a href="javascript:void(0)" @click="changeLoginType"> {{ $t("注册帐号") }}</a>
                                </template>
                                <template v-else>
                                    {{ $t("已经有帐号？") }}
                                    <a href="javascript:void(0)" @click="changeLoginType"> {{ $t("登录帐号") }}</a>
                                </template>
                            </div>
                        </div>
                    </n-form>
                </transition>
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ref } from "vue"
import { userLogin, userReg } from "@/api/modules/user"
import { useMessage } from "naive-ui"
import utils from "@/utils/utils"
import { UserStore } from "@/store/user"
import { useRouter } from "vue-router"
import { MailOutline, LockClosedOutline, CheckmarkCircleOutline } from "@vicons/ionicons5"

const router = useRouter()
const message = useMessage()
const loadIng = ref<boolean>(false)
const code = ref("")
const codeUrl = ref("")
const codeLoad = ref(0)
const userState = UserStore()
const loginMode = ref("access") //qrcode
const codeNeed = ref(false)
const codeId = ref("")
const loginType = ref<String>("login")
const formData = ref({
    email: "",
    password: "",
    confirmPassword: "",
    invite: "",
})

const rules = ref({
    email: { required: true, message: $t('输入您的电子邮箱'), trigger: 'blur' },
    password: { required: true, message: $t('输入您的密码'), trigger: 'blur' },
    confirmPassword: { required: true, message: $t('输入您的密码'), trigger: 'blur' },
})

// 登录
const handleLogin = () => {
    if (formData.value.email == "") return message.info($t("请填写邮箱"))
    if (!utils.isEmail(formData.value.email)) return message.info($t("请填写正确邮箱"))
    if (formData.value.password == "") return message.info($t("请填写密码"))
    loadIng.value = true
    userLogin({
        email: formData.value.email,
        password: formData.value.password,
        code_id: codeId.value,
        code: code.value,
    }).then(({ data, msg }) => {
        userState.info = data
        router.replace("/")
    })
    .catch( res => {
        if (res.data.code == "need") {
            onBlur()
        }
    })
    .finally(() => {
        loadIng.value = false
    })
}

// 注册
const handleReg = () => {
    if (formData.value.email == "") return message.info($t("请填写邮箱"))
    if (!utils.isEmail(formData.value.email)) return message.info($t("请填写正确邮箱"))
    if (formData.value.password == "") return message.info($t("请填写密码"))
    if (formData.value.confirmPassword == "") return message.info($t("请再次确认密码"))
    if (formData.value.confirmPassword != formData.value.password) return message.info($t("两次填写的密码不符"))
    loadIng.value = true
    userReg({
        email: formData.value.email,
        password: formData.value.password,
    }).then(({ data,msg }) => {
        userState.info = data
        router.replace("/")
    }).finally(() => {
        loadIng.value = false
    })
}

// 变更登录类型
const changeLoginType = () => {
    loginType.value == "login" ? (loginType.value = "reg") : (loginType.value = "login")
    if (loginType.value == "reg") {
        codeNeed.value = false
    } else {
        onBlur()
    }
}

// 判断要不要验证码
const onBlur = () => {
    // const upData = {
    //     email: formData.value.email,
    // }
    // needCode(upData)
    // .then(({ data }) => {
    //     codeNeed.value = data
    //     if (codeNeed.value) {
    //         refreshCode()
    //     }
    // })
}

// 刷新验证码
const refreshCode = () => {
    // codeImg()
    //     .then(({ data }) => {
    //         codeUrl.value = data.image_path
    //         codeId.value = data.captcha_id
    //     })
    //     .catch(() => {
    //         codeUrl.value = "error"
    //     })
}
</script>

<style lang="less">
.page-login {
    @apply bg-bg-login flex items-center;

    .login-body {
        @apply flex items-center flex-col max-h-full overflow-hidden py-32 w-full;

        .login-logo {
            @apply block w-84 h-84 bg-logo mb-36;
        }

        .login-box {
            @apply bg-bg-login-box rounded-2xl w-400 max-w-90p shadow-login-box-Shadow relative;

            .login-mode-switch {
                @apply absolute top-1 right-1 z-10 rounded-lg overflow-hidden;

                .login-mode-switch-box {
                    @apply w-80 h-80 cursor-pointer overflow-hidden bg-primary-color-80;
                    transition: background-color 0.3s;
                    transform: translate(40px, -40px) rotate(45deg);

                    &:hover {
                        @apply bg-primary-color;
                    }

                    .login-mode-switch-icon {
                        @apply absolute text-32 w-50 h-50 bottom-negative-20 left-4 flex items-start justify-start text-white;
                        transform: rotate(-45deg);

                        svg {
                            @apply w-30 h-30 ml-26 mt-8;
                        }
                    }
                }
            }

            .login-title {
                @apply text-24 font-semibold text-center mt-46;
            }

            .login-subtitle {
                @apply text-14 text-text-tips text-center mt-12 px-12;
            }

            .login-qrcode {
                @apply flex items-center justify-center m-auto my-50;
            }

            .login-access {
                @apply mt-30 mx-40 mb-32;

                .n-input {
                    @apply mt-6;
                    transition: all 0s;
                }

                .code-load-input {
                    .n-input-wrapper {
                        @apply pr-0;
                    }

                    .login-code-end {
                        @apply h-38 overflow-hidden cursor-pointer ml-1;

                        .code-load,
                        .code-error {
                            @apply h-full flex items-center justify-center w-5 mx-20;
                        }

                        .code-error {
                            @apply w-auto text-14 opacity-80;
                        }

                        img {
                            @apply h-full min-w-16;
                        }
                    }
                }

                .n-button {
                    @apply mt-24 w-full;
                }

                .login-switch {
                    @apply mt-24 text-text-tips;

                    a {
                        @apply text-primary-color;
                        text-decoration: none;
                    }
                }
            }
        }

        .login-bottom {
            @apply flex items-center justify-between mt-24 w-388;

            .login-setting {
                @apply flex items-center cursor-pointer;
            }

            .login-forgot {
                @apply text-text-tips;

                a {
                    @apply text-primary-color;
                    text-decoration: none;
                }
            }
        }
    }
}

input:-webkit-autofill {
    -webkit-box-shadow: 0 0 0px 1000px white inset;
}

.dark input:-webkit-autofill {
    -webkit-box-shadow: 0 0 0px 1000px #2b2b2b inset;
    -webkit-text-fill-color: #ffffff;
}
</style>
