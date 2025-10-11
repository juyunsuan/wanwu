<template>
  <el-dialog
    title="批量编辑元数据值"
    :visible.sync="dialogVisible"
    width="60%"
    :before-close="handleClose"
    class="batch-meta-dialog"
  >
    <div class="dialog-content">
      <!-- 创建元数据值按钮 -->
      <div class="create-section">
        <el-button type="primary" @click="addMetaData" class="create-btn">
          <i class="el-icon-plus"></i>
          创建元数据值
        </el-button>
      </div>

      <!-- 元数据值列表 -->
      <div class="meta-list">
        <div 
          v-for="(item, index) in metaDataList" 
          :key="index" 
          class="meta-item"
        >
          <div class="meta-row">
            <!-- Key选择 -->
            <div class="field-group">
              <label class="field-label">Key:</label>
              <el-select 
                v-model="item.key" 
                placeholder="请选择"
                class="field-select"
                @change="handleKeyChange(item, index)"
              >
                <el-option
                  v-for="meta in keyOptions"
                  :key="meta.metaKey"
                  :label="meta.metaKey + ' | ' + '[ '+meta.metaValueType+' ]'"
                  :value="meta.metaKey"
                />
              </el-select>
            </div>

            <!-- 类型显示 -->
            <div class="field-group type-group">
              <span class="type-label">类型:</span>
              <span class="type-value">[{{ item.type }}]</span>
            </div>

            <el-divider direction="vertical" class="field-divider" />

            <!-- Value输入 -->
            <div class="field-group">
              <label class="field-label">Value:</label>
              <el-input 
                v-model="item.value" 
                placeholder="请输入"
                class="field-input"
              />
            </div>

            <el-divider direction="vertical" class="field-divider" />

            <!-- 删除按钮 -->
            <div class="field-group delete-group">
              <el-button 
                type="text" 
                @click="removeMetaData(index)"
                class="delete-btn"
                icon="el-icon-delete"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- 应用选项 -->
      <div class="apply-section">
        <el-checkbox v-model="applyToAll" class="apply-checkbox">
          应用于所有选定文档
        </el-checkbox>
        <el-tooltip 
          content="若勾选,则自动为所有选定文档创建或编辑元数据值;否则仅对已具有对应元数据值的文档进行编辑。"
          placement="right"
        >
          <i class="el-icon-question question-icon"></i>
        </el-tooltip>
      </div>
    </div>

    <!-- 底部按钮 -->
    <span slot="footer" class="dialog-footer">
      <el-button @click="handleClose" class="cancel-btn">取 消</el-button>
      <el-button type="primary" @click="handleConfirm" class="confirm-btn" :loading="loading">
        确 认
      </el-button>
    </span>
  </el-dialog>
</template>

<script>
import {metaSelect} from "@/api/knowledge"
export default {
  name: 'BatchMetaData',
  data() {
    return {
      dialogVisible: false,
      loading: false,
      applyToAll: false,
      metaDataList: [],
      keyOptions: [],
    }
  },
  created(){
    this.getList();
  },
  methods: {
    getList(){
        const knowledgeId = this.$route.params.id;
        metaSelect({knowledgeId}).then(res =>{
            if(res.code === 0){
                this.keyOptions = res.data.knowledgeMetaList || []
            }
        }).catch(() =>{})
    },
    showDialog() {
      this.dialogVisible = true;
      this.initData();
    },
    
    handleClose() {
      this.dialogVisible = false;
    },
    
    initData() {
      this.metaDataList = [
        {
          key: '',
          type: 'string',
          value: '',
          stringValue: ''
        }
      ];
      this.applyToAll = false;
    },
    
    addMetaData() {
      this.metaDataList.push({
        key: '',
        type: 'string',
        value: '',
        stringValue: ''
      });
    },
    
    removeMetaData(index) {
      this.metaDataList.splice(index, 1);
    },
    
    handleKeyChange(item, index) {
      // 根据选择的key设置对应的类型
      const keyTypeMap = {
        'auto_fetch': 'string',
        'doc_type': 'string',
        'create_time': 'date',
        'author': 'string',
        'version': 'string'
      };
      item.type = keyTypeMap[item.key] || 'string';
    },
    
    handleConfirm() {
      if (this.metaDataList.length === 0) {
        this.$message.warning('请至少添加一个元数据值');
        return;
      }
      
      // 验证必填字段
      for (let i = 0; i < this.metaDataList.length; i++) {
        const item = this.metaDataList[i];
        if (!item.key) {
          this.$message.warning(`第${i + 1}行的Key不能为空`);
          return;
        }
        if (!item.value) {
          this.$message.warning(`第${i + 1}行的Value不能为空`);
          return;
        }
      }
      
      this.loading = true;
      
      // 模拟API调用
      setTimeout(() => {
        this.loading = false;
        this.$message.success('批量编辑成功');
        this.handleClose();
        this.$emit('success', {
          metaDataList: this.metaDataList,
          applyToAll: this.applyToAll
        });
      }, 1000);
    }
  }
}
</script>

<style lang="scss" scoped>
.batch-meta-dialog {
  /deep/ .el-dialog__header {
    padding: 20px 20px 10px;
    border-bottom: 1px solid #f0f0f0;
  }
  
  /deep/ .el-dialog__body {
    padding: 20px;
  }
  
  /deep/ .el-dialog__footer {
    padding: 10px 20px 20px;
    border-top: 1px solid #f0f0f0;
  }
}

.dialog-content {
  .create-section {
    margin-bottom: 20px;
    
    .create-btn {
      background: #384bf7;
      border-color: #384bf7;
      border-radius: 6px;
      padding: 8px 16px;
      
      &:hover {
        background: #2a3cc7;
        border-color: #2a3cc7;
      }
      
      i {
        margin-right: 4px;
      }
    }
  }
  
  .meta-list {
      .meta-item {
        background: #f7f8fa;
        border-radius: 8px;
        padding: 10px;
        margin-bottom: 12px;
      
        .meta-row {
          display: flex;
          align-items: center;
          gap: 16px;
          width: 100%;
          
          .field-divider {
            height: 20px;
            margin: 0 8px;
          }
        
        .field-group {
          display: flex;
          align-items: center;
          flex: 1;
          
          &.type-group {
            flex: 0 0 8%;
            justify-content: center;
          }
          
          &.delete-group {
            flex: 0 0 3%;
            justify-content: center;
          }
          
          .field-label {
            color: #606266;
            font-size: 14px;
            margin-right: 8px;
            min-width: 40px;
          }
          
          .field-select {
            flex: 1;
            min-width: 120px;
            
            /deep/ .el-input__inner {
              border: 1px solid #dcdfe6;
              border-radius: 4px;
              height: 32px;
              line-height: 32px;
            }
          }
          
          .field-input {
            flex: 1;
            min-width: 120px;
            
            /deep/ .el-input__inner {
              border: 1px solid #dcdfe6;
              border-radius: 4px;
              height: 32px;
              line-height: 32px;
            }
          }
          
          .type-label {
            color: #606266;
            font-size: 14px;
            margin-right: 8px;
          }
          
          .type-value {
            color: #384bf7;
            font-size: 14px;
            font-weight: 500;
          }
          
          .string-input {
            flex: 1;
            min-width: 120px;
            
            /deep/ .el-input__inner {
              border: 1px solid #dcdfe6;
              border-radius: 4px;
              height: 32px;
              line-height: 32px;
            }
          }
          
          .delete-btn {
            color: #f56c6c;
            font-size: 16px;
            padding: 4px;
            
            &:hover {
              color: #f78989;
            }
          }
        }
      }
    }
  }
  
  .apply-section {
    display: flex;
    align-items: center;
    margin-top: 20px;
    padding: 0 10px 10px 0;
    .apply-checkbox {
      /deep/ .el-checkbox__label {
        color: #606266;
        font-size: 14px;
      }
      
      /deep/ .el-checkbox__input.is-checked .el-checkbox__inner {
        background-color: #384bf7;
        border-color: #384bf7;
      }
    }
    
    .question-icon {
      color: #909399;
      font-size: 16px;
      margin-left: 8px;
      cursor: pointer;
      
      &:hover {
        color: #606266;
      }
    }
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  
  .cancel-btn {
    background: #fff;
    border: 1px solid #dcdfe6;
    color: #606266;
    border-radius: 4px;
    padding: 8px 16px;
    
    &:hover {
      color: #384bf7;
      border-color: #c6e2ff;
      background-color: #ecf5ff;
    }
  }
  
  .confirm-btn {
    background: #f56c6c;
    border-color: #f56c6c;
    border-radius: 4px;
    padding: 8px 16px;
    
    &:hover {
      background: #f78989;
      border-color: #f78989;
    }
    
    &:focus {
      background: #f56c6c;
      border-color: #f56c6c;
    }
  }
}
</style>
