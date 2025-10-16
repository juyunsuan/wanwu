<template>
  <div class="page-wrapper mcp-management">
    <div class="common_bg">
      <!-- tabs -->
      <div class="mcp-tabs">
        <div :class="['mcp-tab',{ 'active': tabActive === 0 }]" @click="tabClick(0)">内置工具</div>
        <div :class="['mcp-tab',{ 'active': tabActive === 1 }]" @click="tabClick(1)">自定义</div>
      </div>

      <builtInTools ref="builtInTools" v-if="tabActive === 0"/>
      <autoTools ref="autoTools" v-if="tabActive === 1"/>
    </div>
  </div>
</template>
<script>
import builtInTools from './builtInTools'
import autoTools from './autoTools'
export default {
  data() {
    return {
      tabActive: 0
    };
  },
  watch: {
    $route: {
      handler() {
        if (this.$route.query.mcp === "builtIn") this.tabActive = 0
        if (this.$route.query.mcp === "custom") this.tabActive = 1
      },
      // 深度观察监听
      deep: true
    }
  },
  mounted() {
    if (this.$route.query.mcp === "builtIn") this.tabActive = 0
    if (this.$route.query.mcp === "custom") this.tabActive = 1
  },
  methods: {
    tabClick(status){
      this.tabActive = status
    },
  },
  components: {
    builtInTools,
    autoTools
  },
};
</script>
<style lang="scss" scoped>
.mcp-tabs{
  .mcp-tab{
    width: 100px !important;
    height: 30px !important;
    line-height: 30px !important;
    font-size: 10px !important;
    color: #666666 !important;
    border-bottom: 1px solid #CCCCCC !important;
  }
  .active{
    background: #CCCCCC !important;
    color: #333 !important;
  }
}
</style>