<template>
  <div class="power-management">
    <!-- 权限管理弹框 -->
    <el-dialog
      :visible.sync="dialogVisible"
      width="50%"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      class="power-management-dialog"
      @close="handleDialogClose"
    >
      <div slot="title" class="custom-dialog-title">
        <div class="title-content">
          <i class="el-icon-s-custom title-icon"></i>
          <span class="title-text">{{ dialogTitle }}</span>
        </div>
        <div class="title-subtitle" v-if="knowledgeName">
          <span class="knowledge-name">[ {{ knowledgeName }} ]</span>
        </div>
      </div>
        <div class="list-header" v-if="currentView === 'list'">
          <el-button
            type="primary"
            size="small"
            icon="el-icon-plus"
            @click="showCreate"
          >新增</el-button>
        </div>
      <!-- 权限列表 -->
        <PowerList ref="powerList" v-if="currentView === 'list'" @transfer="showTransfer"/>
      <!-- 添加权限 -->
        <PowerCreate ref="powerCreate" v-if="currentView === 'create'"/>
      <!-- 转让权限 -->
        <PowerCreate ref="powerTransfer" v-if="currentView === 'transfer'" :transfer-mode="true" :transfer-data="transferData"/>
      <div
        slot="footer"
        class="dialog-footer"
      >
        <el-button
          v-if="currentView === 'create' || currentView === 'transfer'"
          @click="showList"
        >返回</el-button>
        <el-button
          v-if="currentView === 'create'"
          type="primary"
          @click="handleConfirm"
        >确定</el-button>
        <el-button
          v-if="currentView === 'transfer'"
          type="primary"
          @click="handleTransferConfirm"
        >确定转让</el-button>
        <el-button
          v-if="currentView === 'list'"
          @click="handleCancel"
        >关闭</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import PowerList from "./list.vue";
import PowerCreate from "./create.vue";

export default {
  name: "PowerManagement",
  components: {
    PowerList,
    PowerCreate,
  },
  data() {
    return {
      currentView: "list", // 当前视图：list、create 或 transfer
      dialogVisible: false,
      knowledgeName:'',
      transferData: null, // 转让数据
    };
  },
  computed: {
    dialogTitle() {
      if (this.currentView === "list") {
        return "权限管理";
      } else if (this.currentView === "create") {
        return "添加权限";
      } else if (this.currentView === "transfer") {
        return "转让权限";
      }
      return "权限管理";
    },
  },
  methods: {
    // 显示权限管理弹框
    showDialog() {
      this.currentView = "list";
      this.dialogVisible = true;
    },

    // 显示创建权限页面
    showCreate() {
      this.currentView = "create";
    },

    // 显示转让权限页面
    showTransfer(transferData) {
      this.transferData = transferData;
      this.currentView = "transfer";
    },

    // 显示权限列表页面
    showList() {
      this.currentView = "list";
    },

    // 关闭弹框
    handleCancel() {
      this.dialogVisible = false;
    },

    // 确认添加权限
    handleConfirm() {
      // 这里可以获取 create 组件的数据
      const createData = this.$refs.powerCreate;
      if (createData) {
        console.log("添加权限数据:", {
          selectedPermission: createData.selectedPermission,
          selectedUsers: createData.selectedUsers,
        });

        // 这里可以调用 API 保存数据
        this.$message.success("权限添加成功");

        // 返回列表页面
        this.showList();

        // 刷新列表
        this.refreshList();
      }
    },

    // 对话框关闭事件
    handleDialogClose() {
      // 重置为列表视图
      this.currentView = "list";

      // 可以在这里重置 create 组件的数据
      if (this.$refs.powerCreate) {
        // 重置表单数据等
        console.log("对话框关闭，重置数据");
      }
    },

    // 确认转让权限
    handleTransferConfirm() {
      // 这里可以获取 transfer 组件的数据
      const transferData = this.$refs.powerTransfer;
      if (transferData) {
        console.log("转让权限数据:", {
          selectedUsers: transferData.selectedUsers,
          transferData: this.transferData,
        });

        // 这里可以调用 API 保存转让数据
        this.$message.success("权限转让成功");

        // 转让成功后切换到添加权限页面
        this.showCreate();

        // 刷新列表
        this.refreshList();
      }
    },

    // 刷新列表
    refreshList() {
      if (this.$refs.powerList) {
        // 这里可以调用列表组件的刷新方法
        console.log("刷新权限列表");
      }
    },
  },
};
</script>

<style lang="scss" scoped>
.power-management {
  .list-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    .header-left {
      .page-title {
        font-size: 18px;
        font-weight: 600;
        color: #303133;
      }
    }

    .header-right {
      .el-button {
        border-radius: 4px;
      }
    }
  }
}

.power-management-dialog {
  /deep/ .el-dialog {
    border-radius: 8px;
  }

  /deep/ .el-dialog__header {
    padding: 20px 20px 10px 20px;
    border-bottom: 1px solid #e4e7ed;
  }

  /deep/ .el-dialog__body {
    padding: 20px;
    max-height: 70vh;
    overflow-y: auto;
  }

  /deep/ .el-dialog__footer {
    padding: 10px 20px 20px 20px;
    text-align: right;
    border-top: 1px solid #e4e7ed;
  }
}

.custom-dialog-title {
  display: flex;
  align-items: center;
  
  .title-content {
    display: flex;
    align-items: center;
    
    .title-icon {
      font-size: 20px;
      color: #409eff;
      margin-right: 8px;
    }
    
    .title-text {
      font-size: 18px;
      font-weight: 600;
      color: #303133;
    }
  }
  
  .title-subtitle {
    margin-left: 5px;
    
    .knowledge-name {
      font-size: 14px;
      color: #606266;
      padding: 4px 8px;
    }
  }
}

.dialog-footer {
  .el-button {
    padding: 8px 20px;
    border-radius: 4px;

    &.el-button--primary {
      background-color: #f56c6c;
      border-color: #f56c6c;

      &:hover {
        background-color: #f78989;
        border-color: #f78989;
      }
    }
  }
}
</style>