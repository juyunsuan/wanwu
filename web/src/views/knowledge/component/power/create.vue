<template>
  <div class="add-permission-content">
    <div class="content-wrapper" :class="{ 'transfer-mode': transferMode }">
      <div class="left-panel">
        <div class="search-section">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索用户名/组织/群组名称"
            class="search-input"
          >
          </el-input>
        </div>
        
        <div class="selection-tree">
          <el-tree
            :data="treeData"
            :props="treeProps"
            :show-checkbox="!transferMode"
            node-key="id"
            :default-expand-all="true"
            :check-strictly="false"
            @check="handleTreeCheck"
            @node-click="handleNodeClick"
            :filter-node-method="filterNode"
            class="permission-tree"
            ref="tree"
          >
            <span class="custom-tree-node" slot-scope="{ node, data }" :class="{ 'selected-node': transferMode && data.type === 'user' && isNodeSelected(data.id) }">
              <span class="node-label">{{ node.label }}</span>
              <span v-if="transferMode && data.type === 'user' && isNodeSelected(data.id)" class="selected-icon">
                <i class="el-icon-check"></i>
              </span>
            </span>
          </el-tree>
        </div>
      </div>
      
      <div class="right-panel" v-if="!transferMode">
        <div class="permission-section">
          <div class="permission-label">权限:</div>
          <el-select v-model="selectedPermission" placeholder="请选择权限" class="permission-select">
            <el-option label="可读" :value="0"></el-option>
            <el-option label="可编辑" :value="10"></el-option>
            <el-option label="管理员" :value="20"></el-option>
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
    knowledgeId: {
      type: String,
      default: ''
    }
  },
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
    defaultPermission() {
      return this.transferMode ? '管理员' : '可读'
    }
  },
  data() {
    return {
      searchKeyword: '',
      selectedPermission: '可读',
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
      selectedUsers: []
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
    },
    searchKeyword(val){
      this.$refs.tree.filter(val);
    }
  },
  methods: {
    isNodeSelected(nodeId) {
      return this.selectedUsers.some(user => user.id === nodeId)
    },
    filterNode(value,data){
      if (!value) return true;
      return data.name.indexOf(value) !== -1;
    },
    handleTreeCheck(data, checkedInfo) {
      console.log('选择变化:', data, checkedInfo)
      
      const checkedNodes = checkedInfo.checkedNodes || []
      const halfCheckedNodes = checkedInfo.halfCheckedNodes || []
      
      const selectedUserNodes = checkedNodes.filter(node => node.type === 'user')
      
      this.selectedUsers = selectedUserNodes.map(node => ({
        id: node.id,
        name: node.name,
        organization: node.organization
      }))
      
      this.updateTreeSelectionState(checkedNodes)
    },
    handleNodeClick(data, node) {
      if (this.transferMode && data.type === 'user') {
        this.selectedUsers = [{
          id: data.id,
          name: data.name,
          organization: data.organization
        }]
        
        this.updateTreeSelectionState([data])
        
        this.updateSelectedNodeBackground()
      }
    },
    removeUser(user) {
      user.selected = false
      this.selectedUsers = this.selectedUsers.filter(u => u.id !== user.id)
    },
    removeSelectedUser(user) {
      this.selectedUsers = this.selectedUsers.filter(u => u.id !== user.id)
      
      this.updateTreeSelection(user.id, false)
      
      this.$nextTick(() => {
        if (this.$refs.tree) {
          if (this.transferMode) {
            this.$refs.tree.setCheckedKeys([])
          } else {
            const checkedKeys = this.$refs.tree.getCheckedKeys()
            const newCheckedKeys = checkedKeys.filter(key => key !== user.id)
            this.$refs.tree.setCheckedKeys(newCheckedKeys)
          }
        }
      })
    },
    updateTreeSelection(userId, selected) {
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
    updateTreeSelectionState(checkedNodes) {

      const updateNode = (nodes) => {
        nodes.forEach(node => {
          if (node.type === 'user') {
            const isChecked = checkedNodes.some(checkedNode => checkedNode.id === node.id)
            node.selected = isChecked
          }
          if (node.children) {
            updateNode(node.children)
          }
        })
      }
      updateNode(this.treeData)
    },
    createNewGroup() {
      this.$message.info('创建新群组功能')
    },
    updateSelectedNodeBackground() {
      this.$nextTick(() => {
        const allContents = document.querySelectorAll('.permission-tree .el-tree-node__content')
        allContents.forEach(content => {
          content.classList.remove('selected-content')
        })
        
        this.selectedUsers.forEach(user => {
          const nodeContent = document.querySelector(`[data-key="${user.id}"] .el-tree-node__content`)
          if (nodeContent) {
            nodeContent.classList.add('selected-content')
          }
        })
      })
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
            
            &.selected-content {
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
            
            .selected-icon {
              color: #384BF7;
              font-size: 16px;
              margin-right: 8px;
            }
            
            &.selected-node {
              .node-label {
                color: #384BF7;
              }
            }
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
