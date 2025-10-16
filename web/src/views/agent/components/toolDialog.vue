<template>
    <div>
        <el-dialog
            title="新增工具"
            :visible.sync="dialogVisible"
            width="40%"
            :before-close="handleClose">
            <div class="tool-typ">
                <div class="toolbtn">
                    <div v-for="(item,index) in toolList" :key="index" @click="clickTool(item,index)" :class="[{'active':activeValue === item.value}]">
                        {{item.name}}
                    </div>
                </div>
                <el-input v-model="toolName" placeholder="搜索工具" class="tool-input" suffix-icon="el-icon-search" @keyup.enter.native="searchTool" clearable></el-input>
            </div>
            <div class="toolContent">
                <div @click="goCreate" class="createTool"><span class="el-icon-plus add"></span>{{createText()}}</div>
                <template v-for="(items, type) in contentMap">
                    <div 
                    v-if="activeValue === type"
                    v-for="(item,i) in items"
                    :key="item[type + 'Id'] || item.id"
                    class="toolContent_item"
                    >
                    <span>{{ item.apiName || item.name }}</span>
                    <el-button type="text" @click="openTool($event,item,type)" v-if="!item.checked">添加</el-button>
                    <el-button type="text" v-else style="color:#ccc;">已添加</el-button>
                    </div>
                </template>
            </div>
        </el-dialog>
    </div>
</template>
<script>
import { getList } from '@/api/workflow.js';
import { addWorkFlowInfo, addMcp,customList,addCustom } from "@/api/agent";
import { getExplorationFlowList} from "@/api/workflow";
export default {
    props:['assistantId'],
    data(){
        return {
            toolName:'',
            dialogVisible:false,
            toolIndex:0,
            activeValue:'auto',
            workFlowInfos:[],
            mcpInfos:[],
            customInfos:[],
            mcpList:[],
            workFlowList:[],
            customList:[],
            toolList:[
                {
                    value:'auto',
                    name:'自定义'
                },
                {
                    value:'mcp',
                    name:'MCP'
                },
                {
                    value:'workflow',
                    name:'工作流'
                }
            ]
        }
    },
    computed:{
         contentMap() {
            return {
            auto: this.customInfos,
            mcp: this.mcpInfos,
            workflow: this.workFlowInfos
            }
        }
    },
    created(){
        this.getMcpSelect('');
        this.getWorkflowList('');
        this.getCustomList('')
    },
    methods:{
        getCustomList(name){
            customList({name}).then(res =>{
                if(res.code === 0){
                    this.customInfos = (res.data.list || []).map(m => ({
                        ...m,
                        checked: this.customList.some(item => item.customId === m.customToolId)
                    }));
                }
            }).catch(() =>{

            })
        },
        goCreate(){
            if(this.activeValue === 'auto'){
                this.$router.push({path:'/tool?type=tool&mcp=custom'})
            }else if(this.activeValue === 'mcp'){
                this.$router.push({path:'/tool?type=mcp&mcp=mcp'})
            }else{
                this.$router.push({path:'/appSpace/workflow'})
            }
        },
        createText(){
            if(this.activeValue === 'auto'){
                return '创建自定义工具'
            }else if(this.activeValue === 'mcp'){
                return '导入MCP'
            }else{
                return '创建工作流'
            }
        },
        openTool(e,item,type){
            if(!e) return;
            item.checked = !item.checked
            if(type === 'workflow'){
                this.addWorkFlow(item)
            }else if(type === 'mcp'){
                this.addMcpItem(item)
            }else{
                this.addCustomItem(item)
            }
        },
        addCustomItem(n){
            addCustom({assistantId:this.assistantId,customToolId:n.customToolId}).then(res =>{
                if(res.code === 0){
                    n.checked = true;
                    this.$message.success('工具添加成功');
                    this.$emit('updateDetail');
                }
            }).catch(() =>{

            })
        },
        addMcpItem(n){
            addMcp({assistantId:this.assistantId,mcpId:n.mcpId}).then(res =>{
                if(res.code === 0){
                    n.checked = true;
                    this.$message.success('工具添加成功');
                    this.$emit('updateDetail');
                }
            }).catch(() =>{

            })
        },
        addWorkFlow(n){
            this.doCreateWorkFlow(n,n.appId)
        },
        async doCreateWorkFlow(n,workFlowId, schema){
            let params = {
                assistantId: this.assistantId,
                workFlowId
            };
            let res = await addWorkFlowInfo(params);
            if (res.code === 0) {
                n.checked = true;
                this.$message.success(this.$t('agent.addPluginTips'));
                this.$emit('updateDetail');
            }
        },
        searchTool(){
            if(this.activeValue === 'auto'){
                this.getCustomList(this.toolName)
            }else if(this.activeValue === 'mcp'){
                this.getMcpSelect(this.toolName)
            }else{
                this.getWorkflowList(this.toolName)
            }
        },
        getMcpSelect(name){
            getList({name}).then(res => {
                if(res.code === 0){
                    this.mcpInfos = (res.data.list || []).map(m => ({
                        ...m,
                        checked: this.mcpList.some(item => item.mcpId === m.mcpId)
                    }));
                }
               
            }).catch(err => {

            })
            },
        getWorkflowList(name) {
                getExplorationFlowList({name,appType:'workflow',searchType:'all'}).then(res =>{
                    if (res.code === 0) {
                        this.workFlowInfos = (res.data.list || []).map(m => ({
                            ...m,
                            checked: this.workFlowList.some(item => item.workFlowId === m.appId)
                        }));
                    }
                })
        },
        showDialog(row){
            this.dialogVisible = true;
            this.setMcp(row.mcpInfos);
            this.setWorkflow(row.workFlowInfos);
            this.setCustom(row.customInfos)
            this.mcpList = row.mcpInfos || [];
            this.workFlowList = row.workFlowInfos || [];
            this.customList  = row.customInfos || [];
        },
        setMcp(data){
           this.mcpInfos = this.mcpInfos.map(m => ({
            ...m,
            checked: data.some(item => item.mcpId === m.mcpId)
            }));
        },
        setWorkflow(data){
            this.workFlowInfos = this.workFlowInfos.map(m => ({
            ...m,
            checked: data.some(item => item.workFlowId === m.appId)
            }));
        },
        setCustom(data){
            this.customInfos = this.customInfos.map(m => ({
                ...m,
                checked: data.some(item => item.customId === m.customToolId)
            }));
        },
        handleClose(){
            this.toolIndex = -1;
            this.activeValue = 'auto'
            this.dialogVisible = false;
        },
        clickTool(item,i){
            this.toolIndex = i;
            this.activeValue = item.value;
        }
    }
}
</script>
<style lang="scss" scoped>
/deep/{
    .el-dialog__body{
        padding:10px 20px;
    }
}
.createTool{
    padding:10px;
    cursor: pointer;
    .add{
        padding-right:5px;
    }
}
.createTool:hover{
    color: #384BF7;
}
.tool-typ{
    display:flex;
    justify-content:space-between;
    padding:10px 0;
    border-bottom: 1px solid #dbdbdb;
    .toolbtn{
        display:flex;
        justify-content:flex-start;
        gap:20px;
        div{
            text-align: center;
            padding:5px 20px;
            border-radius:6px;
            border:1px solid #ddd;
            cursor: pointer;
        }
    }
    .tool-input{
        width:200px;
    }
}
.toolContent{
    padding:10px 0;
    max-height:300px;
    overflow-y:auto;
    .toolContent_item{
        padding:5px 20px;
        border:1px solid #dbdbdb;
        border-radius:6px;
        margin-bottom:10px;
        cursor: pointer;
        display: flex;
        align-items:center;
        justify-content:space-between;
    }
    .toolContent_item:hover{
        background:#f4f5ff;
    }
}
.active{
    border:1px solid #384BF7 !important;
    color: #fff;
    background:#384BF7;
}
</style>