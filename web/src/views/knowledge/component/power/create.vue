<template>
  <div class="add-permission-content">
    <div class="content-wrapper" :class="{ 'transfer-mode': transferMode }">
      <!-- 左侧选择面板 -->
      <div class="left-panel">
        <div class="search-section">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索用户名/组织/群组名称"
            class="search-input"
            @input="handleSearch"
          >
            <el-button slot="append" icon="el-icon-search" @click="handleSearch"></el-button>
          </el-input>
        </div>
        
        <div class="selection-tree">
          <el-tree
            :data="treeData"
            :props="treeProps"
            show-checkbox
            node-key="id"
            :default-expand-all="true"
            :check-strictly="false"
            @check="handleTreeCheck"
            class="permission-tree"
          >
            <span class="custom-tree-node" slot-scope="{ node, data }">
              <span class="node-label">{{ node.label }}</span>
              <!-- <span v-if="data.type === 'user' && data.selected" class="remove-icon" @click="removeUser(data)">
                <i class="el-icon-close"></i>
              </span> -->
            </span>
          </el-tree>
        </div>
      </div>
      
      <!-- 右侧已选择面板 - 转让模式下不显示 -->
      <div class="right-panel" v-if="!transferMode">
        <div class="permission-section">
          <div class="permission-label">权限:</div>
          <el-select v-model="selectedPermission" placeholder="请选择权限" class="permission-select">
            <el-option
              v-for="item in permissionOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            ></el-option>
          </el-select>
        </div>
        
        <div class="selected-users-section">
          <div class="selected-users-list">
            <div
              v-for="user in selectedUsers"
              :key="user.id"
              class="selected-user-item"
            >
              <span class="user-info">{{ user.name }} {{ user.organization }}</span>
              <i class="el-icon-close remove-icon" @click="removeSelectedUser(user)"></i>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'AddPermission',
  props: {
    transferMode: {
      type: Boolean,
      default: false
    },
    transferData: {
      type: Object,
      default: null
    }
  },
  computed: {
    // 转让模式下默认选择管理员权限
    defaultPermission() {
      return this.transferMode ? '管理员' : '可读'
    }
  },
  data() {
    return {
      searchKeyword: '',
      selectedPermission: '可读',
      permissionOptions: [
        { label: '管理员', value: '管理员' },
        { label: '编辑者', value: '编辑者' },
        { label: '查看者', value: '查看者' }
      ],
      treeProps: {
        children: 'children',
        label: 'name'
      },
      treeData: [
        {
          id: 'org1',
          name: '组织1',
          type: 'organization',
          children: [
            { id: 'user1_org1', name: '用户1', type: 'user', organization: '组织1', selected: false },
            { id: 'user2_org1', name: '用户2', type: 'user', organization: '组织1', selected: false },
            { id: 'user3_org1', name: '用户3', type: 'user', organization: '组织1', selected: false },
            { id: 'user4_org1', name: '用户4', type: 'user', organization: '组织1', selected: false },
            { id: 'user5_org1', name: '用户5', type: 'user', organization: '组织1', selected: false }
          ]
        },
        {
          id: 'group1',
          name: '群组1',
          type: 'group',
          children: [
            { id: 'user1_group1', name: '用户1', type: 'user', organization: '群组1', selected: true },
            { id: 'user2_group1', name: '用户2', type: 'user', organization: '群组1', selected: true }
          ]
        }
      ],
      selectedUsers: [
        { id: 'user1_group1', name: '用户1', organization: '组织1' },
        { id: 'user2_group1', name: '用户2', organization: '组织1' }
      ]
    }
  },
  watch: {
    transferMode: {
      handler(newVal) {
        if (newVal) {
          this.selectedPermission = '管理员'
        }
      },
      immediate: true
    }
  },
  methods: {
    handleSearch() {
      // 搜索逻辑
      console.log('搜索:', this.searchKeyword)
    },
    handleTreeCheck(data, checkedInfo) {
      // 处理树形选择
      console.log('选择变化:', data, checkedInfo)
    },
    removeUser(user) {
      // 移除用户
      user.selected = false
      this.selectedUsers = this.selectedUsers.filter(u => u.id !== user.id)
    },
    removeSelectedUser(user) {
      // 移除已选择的用户
      this.selectedUsers = this.selectedUsers.filter(u => u.id !== user.id)
      // 同时更新树形数据中的选中状态
      this.updateTreeSelection(user.id, false)
    },
    updateTreeSelection(userId, selected) {
      // 更新树形数据中的选中状态
      const updateNode = (nodes) => {
        nodes.forEach(node => {
          if (node.id === userId) {
            node.selected = selected
          }
          if (node.children) {
            updateNode(node.children)
          }
        })
      }
      updateNode(this.treeData)
    },
    createNewGroup() {
      // 创建新群组
      this.$message.info('创建新群组功能')
    }
  }
}
</script>

<style lang="scss" scoped>
.add-permission-content {
  background: #fff;
  border-radius: 8px;
  
  .content-wrapper {
    display: flex;
    gap: 15px;
    height: 400px;
    
    &.transfer-mode {
      .left-panel {
        flex: 1;
        width: 100%;
      }
    }
  
    .left-panel {
      flex: 1;
      display: flex;
      flex-direction: column;
      border: 1px solid #e4e7ed;
      border-radius: 4px;
      padding: 15px;
    
      .search-section {
        margin-bottom: 15px;
        
        .search-input {
          /deep/ .el-input-group__append {
            background-color: #384BF7;
            border-color: #384BF7;
            color: white;
            
            &:hover {
              background-color: #384BF7 !important;
              border-color: #384BF7 !important;
              color: white !important;
            }
            
            &:active {
              background-color: #384BF7 !important;
              border-color: #384BF7 !important;
              color: white !important;
            }
            
            .el-button {
              background-color: transparent !important;
              border-color: transparent !important;
              color: white !important;
              width: 30px !important;
              height: 30px !important;
              padding: 0 !important;
              margin: 0 !important;
              min-width: 30px !important;
              min-height: 30px !important;
              
              &:hover {
                background-color: transparent !important;
                border-color: transparent !important;
                color: white !important;
                width: 30px !important;
                height: 30px !important;
                padding: 0 !important;
                margin: 0 !important;
                min-width: 30px !important;
                min-height: 30px !important;
              }
              
              &:active {
                background-color: transparent !important;
                border-color: transparent !important;
                color: white !important;
                width: 30px !important;
                height: 30px !important;
                padding: 0 !important;
                margin: 0 !important;
                min-width: 30px !important;
                min-height: 30px !important;
              }
            }
          }
        }
      }
      
      .selection-tree {
        flex: 1;
        overflow-y: auto;
        .permission-tree {
          /deep/ .el-tree-node__content {
            height: 32px;
            line-height: 32px;
            
            &:hover {
              background-color: #f5f7fa;
            }
          }
          
          /deep/ .el-checkbox {
            margin-right: 8px;
          }
          
          .custom-tree-node {
            display: flex;
            align-items: center;
            justify-content: space-between;
            width: 100%;
            
            .node-label {
              flex: 1;
              font-size: 14px;
              color: #606266;
            }
          }
        }
      }
      
      .new-group-link {
        margin-top: 10px;
        text-align: center;
        
        .new-group-btn {
          color: #384BF7;
          font-size: 14px;
          
          &:hover {
            color: #5a6cff;
          }
        }
      }
    }
    
    .right-panel {
      flex: 1;
      display: flex;
      flex-direction: column;
      border: 1px solid #e4e7ed;
      border-radius: 4px;
      padding: 15px;
      
      .permission-section {
        margin-bottom: 20px;
        display: flex;
        align-items: center;
        
        .permission-label {
          font-size: 14px;
          color: #606266;
          margin-right: 10px;
          white-space: nowrap;
        }
        
        .permission-select {
          flex: 1;
        }
      }
      
      .selected-users-section {
        flex: 1;
        
        .selected-users-list {
          max-height: 300px;
          overflow-y: auto;
          
          .selected-user-item {
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding: 8px 10px;
            cursor: pointer;
            border-radius: 4px;
            background-color: #f5f7fa;
            border: 1px solid transparent;
            transition: all 0.3s ease;
            margin-bottom: 8px;
            
            &:last-child {
              margin-bottom: 0;
            }
            
            &:hover {
              background-color: #f5f7fa;
              border-color: #384BF7;
            }
            
            .user-info {
              font-size: 14px;
              color: #606266;
            }
            
            .remove-icon {
              color: #384BF7;
              cursor: pointer;
              font-size: 12px;
              padding: 2px;
              border-radius: 2px;
              opacity: 0;
              transition: opacity 0.3s ease;
            }
          }
          
          .selected-user-item:hover .remove-icon {
            opacity: 1;
          }
        }
      }
    }
  }
}
</style>
