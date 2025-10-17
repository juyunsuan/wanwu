import service from "@/utils/request"
const BASE_URL = '/user/api/v1'
export const createAgent = (data)=>{
    return service({
        url: `${BASE_URL}/assistant`,
        method: 'post',
        data
    })
}

export const updateAgent = (data)=>{
    return service({
        url: `${BASE_URL}/assistant`,
        method: 'put',
        data
    })
}
export const delAgent = (data)=>{
    return service({
        url: `${BASE_URL}/assistant`,
        method: 'delete',
        data
    })
}
export const getAgentInfo = (params)=>{
    return service({
        url: `${BASE_URL}/assistant`,
        method: 'get',
        params
    })
}
export const putAgentInfo = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/config`,
        method: 'put',
        data
    })
}
export const createConversation = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/conversation`,
        method: 'post',
        data
    })
}
export const delConversation = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/conversation`,
        method: 'delete',
        data
    })
}
export const getConversationHistory = (params)=>{
    return service({
        url: `${BASE_URL}/assistant/conversation/detail`,
        method: 'get',
        params
    })
}
export const getConversationlist = (params)=>{
    return service({
        url: `${BASE_URL}/assistant/conversation/list`,
        method: 'get',
        params
    })
}
export const getActionInfo = (params)=>{
    return service({
        url: `${BASE_URL}/assistant/action`,
        method: 'get',
        params
    })
}
export const editActionInfo = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/action`,
        method: 'put',
        data
    })
}
export const addActionInfo = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/action`,
        method: 'post',
        data
    })
}
export const delActionInfo = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/action`,
        method: 'delete',
        data
    })
}
export const enableAction = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/action/enable`,
        method: 'put',
        data
    })
}
export const addWorkFlowInfo = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/tool/workflow`,
        method: 'post',
        data
    })
}
export const delWorkFlowInfo = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/tool/workflow`,
        method: 'delete',
        data
    })
}
export const enableWorkFlow = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/tool/workflow/switch`,
        method: 'put',
        data
    })
}
export const agentStream = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/stream`,
        method: 'post',
        data
    })
}
export const agentTestStream = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/test/stream`,
        method: 'post',
        data
    })
}
export const getAgentList = (params)=>{
    return service({
        url: `${BASE_URL}/assistant/list`,
        method: 'get',
        params
    })
}

//删除mcp
export const deleteMcp = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/tool/mcp`,
        method: 'delete',
        data
    })
}
//添加mcp
export const addMcp = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/tool/mcp`,
        method: 'post',
        data
    })
}
//启停mcp
export const enableMcp = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/tool/mcp/switch`,
        method: 'put',
        data
    })
}
//删除自定义工具
export const deleteCustom = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/tool/custom`,
        method: 'delete',
        data
    })
}
//添加自定义工具
export const addCustom = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/tool/custom`,
        method: 'post',
        data
    })
}
// 自定义工具列表
export const customList = (data)=>{
    return service({
        url: `${BASE_URL}/tool/custom/select`,
        method: 'get',
        params:data
    })
}
//启停自定义工具
export const enableCustom = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/tool/custom/switch`,
        method: 'put',
        data
    })
}
//编辑url
export const editOpenurl = (data)=>{
    return service({
        url: `${BASE_URL}/appspace/app/openurl`,
        method: 'put',
        data
    })
}
//创建url
export const createOpenurl = (data)=>{
    return service({
        url: `${BASE_URL}/appspace/app/openurl`,
        method: 'post',
        data
    })
}
//删除应用url
export const delOpenurl = (data)=>{
    return service({
        url: `${BASE_URL}/appspace/app/openurl`,
        method: 'delete',
        data
    })
}
//获取应用url列表
export const getOpenurl = (data)=>{
    return service({
        url: `${BASE_URL}/appspace/app/openurl/list`,
        method: 'get',
        params:data
    })
}
//启停应用url状态
export const switchOpenurl = (data)=>{
    return service({
        url: `${BASE_URL}/appspace/app/openurl/status`,
        method: 'put',
        data
    })
}


//获取智能体openurl信息
export const getOpenurlInfo = (suffix,config={})=>{
    return service({
        url: `/service/url/openurl/v1/agent/${suffix}`,
        method: 'get',
        ...config,
        isOpenUrl:true
    })
}
//智能体openurl创建智能体对话
export const openurlConversation = (data,suffix,config={})=>{
    return service({
        url: `/service/url/openurl/v1/agent/${suffix}/conversation`,
        method: 'post',
        data,
        ...config,
        isOpenUrl:true
    })
}
//删除智能体openurl创建智能体对话
export const delOpenurlConversation = (data,suffix,config={})=>{
    return service({
        url: `/service/url/openurl/v1/agent/${suffix}/conversation`,
        method: 'delete',
        data,
        ...config,
        isOpenUrl:true
    })
}
//智能体openurl详情历史列表
export const OpenurlConverHistory = (data,suffix,config={})=>{
    return service({
        url: `/service/url/openurl/v1/agent/${suffix}/conversation/detail`,
        method: 'get',
        params:data,
        ...config,
        isOpenUrl:true
    })
}
//智能体openurl对话列表
export const OpenurlConverList = (suffix,config={})=>{
    return service({
        url: `/service/url/openurl/v1/agent/${suffix}/conversation/list`,
        method: 'get',
        ...config,
        isOpenUrl:true
    })
}
//智能体openurl流式对话
export const OpenurlStream = (data,suffix,config={})=>{
    return service({
        url: `/service/url/openurl/v1/agent/${suffix}/stream`,
        method: 'post',
        data,
        ...config,
        isOpenUrl:true
    })
}