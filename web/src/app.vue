<template>
    <!-- :theme="theme" -->
    <n-config-provider  :locale="locale" :date-locale="dateLocale"
        :theme-overrides="theme === null ? lightThemeOverrides : darkThemeOverrides">
        <n-loading-bar-provider>
            <n-message-provider>
                <n-notification-provider>
                    <n-dialog-provider>
                        <Result v-if="resultCode > 0" />
                        <router-view v-else />
                    </n-dialog-provider>
                </n-notification-provider>
            </n-message-provider>
        </n-loading-bar-provider>
    </n-config-provider>
</template>

<script lang="ts" setup>
import { GlobalStore } from "@/store"
import { zhCN, dateZhCN, enUS, dateEnUS } from "naive-ui"
import result from "@/utils/result"
import { computed } from "vue"
import { NLocale } from "naive-ui/es/locales/common/enUS"
import { NDateLocale } from "naive-ui/es/locales"
import themeOverrides from "@/utils/naive.config"

const lightThemeOverrides = themeOverrides.lightThemeOverrides
const darkThemeOverrides = themeOverrides.darkThemeOverrides
const resultCode = result.code()
const globalStore = GlobalStore()
const { themeName, theme } = globalStore.appSetup()

const locale = computed((): NLocale => {
    if (globalStore.language == "en") return enUS
    return zhCN
})

const dateLocale = computed((): NDateLocale => {
    if (globalStore.language == "en") return dateEnUS
    return dateZhCN
})

watch(themeName,(newVal,oldVal)=>{
    document.querySelector("body").classList.remove('y-theme-'+oldVal)
    document.querySelector("body").classList.add('y-theme-'+newVal)
},{immediate:true})

</script>

<style lang="less">
.child-view {
    @apply absolute w-full bottom-0 top-0 right-0 left-0 min-h-full;
}
</style>

