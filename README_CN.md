<div align="center">
  <img src="https://github.com/user-attachments/assets/6ceb4269-a861-4545-84db-bad322592156" style="width:45%; height:auto;" />
<p>
  <a href="#-核心功能模块">核心功能模块</a> •
  <a href="#-典型应用场景">典型应用场景</a> •
  <a href="#-快速开始">快速开始</a> •
  <a href="#-使用万悟">使用万悟</a> •
  <a href="#-q--a">Q & A</a> •
  <a href="#-联系我们">联系我们</a> 
</p>
<p>
  <img alt="License" src="https://img.shields.io/badge/license-apache2.0-blue.svg">
  <img alt="Go Version" src="https://img.shields.io/badge/go-%3E%3D%201.24.0-blue">
  </a>
  <a href="https://github.com/UnicomAI/wanwu/releases">
    <img alt="Release Notes" src="https://img.shields.io/github/v/release/UnicomAI/wanwu?label=Release&logo=github&color=green">
  </a>
</p>
<p align="center">
    <a href="https://github.com/UnicomAI/wanwu/blob/main/README.md">English</a> |
    简体中文 |
    <a href="https://github.com/UnicomAI/wanwu/blob/main/README_繁體.md">繁體中文</a>
</p>
</div>


&emsp;&emsp;**元景万悟智能体平台**是一款面向**企业级**场景的**一站式**、**商用license友好**的**智能体开发平台**，致力于为企业提供安全、高效、合规的一站式AI解决方案。我们以"技术开放、生态共建"为核心理念，通过整合大语言模型、业务流程自动化等前沿技术，构建了覆盖模型全生命周期管理、MCP、联网检索、智能体快速开发、企业知识库建设、复杂工作流编排等完整功能体系的AI工程化平台。平台采用模块化架构设计，支持灵活的功能扩展和二次开发，在确保企业数据安全和隐私保护的同时，大幅降低了AI技术的应用门槛。无论是中小型企业快速构建智能化应用，还是大型企业实现复杂业务场景的智能化改造，元景万悟智能体平台都能提供强有力的技术支撑，助力企业加速数字化转型进程，实现降本增效和业务创新。

------

<div>
  <p align="center">
    <a href="https://www.bilibili.com/video/BV1HxpazNEAM"><img width="400" src="https://github.com/user-attachments/assets/54efe5d3-c28d-48fb-9a6e-d6ac536a1f95" /></a>
    <a href="https://www.bilibili.com/video/BV1HxpazNEAM"><img width="394" src="https://github.com/user-attachments/assets/d19831e6-10a3-4ee0-8caf-6c0ebe2af4a5" /></a>
  </p>
</div>

------

### &#x1F525; 采用宽松友好的 Apache 2.0 License，支持开发者自由扩展与二次开发

✔ **企业级工程化**：提供从模型纳管到应用落地的完整工具链，解决LLM技术落地"最后一公里"问题  

✔ **开放开源生态**：采用宽松友好的 **Apache 2.0 License**，支持开发者自由扩展与二次开发  

✔ **全栈技术支持**：配备专业团队为生态伙伴提供 **架构咨询、性能优化** 全周期赋能  

✔ **多租户架构**：提供多租户账号体系，满足用户成本控制、数据安全隔离、业务弹性扩展、行业定制化、快速上线及生态协同等核心需求

✔ **信创适配**：已适配国产信创数据库TiDB和OceanBase

------

### 🚩 核心功能模块

**1. 模型纳管（Model Hub）**

▸ 支持 **数百种专有/开源大模型**（包括GPT、Claude、Llama等系列）的统一接入与生命周期管理

▸ 深度适配 **OpenAI API 标准** 及 **联通元景** 生态模型，实现异构模型的无缝切换

▸ 提供 **多推理后端支持**（vLLM、TGI等）与 **自托管解决方案**，满足不同规模企业的算力需求

#### **2. MCP**

▸ **标准化接口**：使 AI 模型能够无缝连接各种外部工具（如 GitHub、Slack、数据库等），而无需为每个数据源单独开发适配器

▸ **内置丰富精选推荐**：整合100+行业MCP接口，让用户方便快捷，轻松调用

#### **3. 联网检索**（Web Search）

▸ **实时信息获取**：具备强大的联网检索能力，能够实时从互联网获取最新的信息。在问答场景中，当用户的问题需要最新的新闻、数据等信息时，平台可以快速检索并返回准确的结果，提升回答的时效性和准确性

▸ **多源数据整合**：整合了多种互联网数据源，包括新闻网站、学术数据库、行业报告等。通过对多源数据的整合和分析，为用户提供更全面、更深入的信息。例如，在市场调研场景中，可以同时从多个数据源获取相关数据，进行综合分析和评估

▸ **智能检索策略**：采用智能检索算法，根据用户的问题自动优化检索策略，提高检索效率和准确性。支持关键词检索、语义检索等多种检索方式，满足不同用户的需求。同时，对检索结果进行智能排序和筛选，优先展示最相关、最有价值的信息

#### **4. 可视化工作流（Workflow Studio）**

▸ 通过 **低代码拖拽画布** 快速构建复杂AI业务流程

▸ 内置 **条件分支、API、大模型、知识库、代码、MCP** 等多种节点，支持端到端流程调试与性能分析

#### **5. 企业级知识库、RAG Pipeline**

▸ 提供**知识库创建**→ **文档解析→向量化→检索→精排** 的全流程知识管理能力，支持pdf/docx/txt/xlsx/csv/pptx等 **多种格式** 文档，还支持网页资源的抓取和接入

▸ 集成 **多模态检索** 、**级联切分** 与 **自适应切分**，显著提升问答准确率

#### **6. 智能体开发框架（Agent Framework）**

▸ 可基于 **函数调用（Function Calling）** 的Agent构建范式，支持工具扩展、私域知识库关联与多轮对话

▸ 支持**在线调试**

#### **7. 后端即服务（BaaS）**

▸ 提供 **RESTful API** ，支持与企业现有系统（OA/CRM/ERP等）深度集成

▸ 提供 **细粒度权限控制**，保障生产环境稳定运行

------

### &#x1F4E2; 功能比较

|    功能     |  元景万悟智能体平台  |       Dify.AI        |       Fastgpt        |      Ragflow       |      Coze开源版      |
| :---------: | :----------------: | :------------------: | :------------------: | :----------------: | :----------------: |
| 模型导入    |         ✅          |          ✅           |        ❌(内置模型)        |         ✅          | ❌(内置模型) |
|   RAG引擎   |         ✅          |          ✅           |          ✅           |         ✅          | ✅ |
|    MCP     |         ✅          |          ✅           |          ✅           |  ✅（需安装工具使用） | ❌ |
| 直接导入OCR |         ✅          |          ❌           |          ❌           |         ❌          | ❌ |
|  搜索增强   |         ✅          |  ✅（需安装工具使用）   |          ✅           |  ✅（需安装工具使用） | ✅ |
|    Agent    |         ✅          |          ✅           |          ✅           |         ✅          | ✅ |
|   工作流    |         ✅          |          ✅           |          ✅           |         ✅          | ✅ |
|  本地部署   |         ✅          |          ✅           |          ✅           |         ✅          | ✅ |
| license友好 |         ✅          |   ❌（商用有限制）    |   ❌（商用有限制）    |     未完全开源     | ✅ |
|   多租户    |         ✅          |   ❌（商用有限制）    |   ❌（商用有限制）    |         ✅          | ✅（但用户间不互通） |

> 截止2025年8月1日对比。

------

### &#x1F3AF; 典型应用场景

- **智能客服**：基于RAG+Agent实现高准确率的业务咨询与工单处理  
- **知识管理**：构建企业专属知识库，支持语义搜索与智能摘要生成  
- **流程自动化**：通过工作流引擎实现合同审核、报销审批等业务的AI辅助决策  

平台已成功应用于 **金融、工业、政务** 等多个行业，助力企业将LLM技术的理论价值转化为实际业务收益。我们诚邀开发者加入开源社区，共同推动AI技术的民主化进程。  

------

### 🚀 快速开始

- 元景万悟智能体平台的工作流模块使用的是以下项目，可到其仓库查看详情。

  - v0.1.8及以前：wanwu-agentscope 项目

  - v0.2.0开始：[wanwu-workflow](https://github.com/UnicomAI/wanwu-workflow/tree/dev/wanwu-backend) 项目

- **Docker安装（推荐）**

1. 首次运行前

    1.1 拷贝环境变量文件

    ```bash
    cp .env.bak .env
    ```

    1.2 根据系统修改.env文件中的`WANWU_ARCH`、`WANWU_EXTERNAL_IP`变量

    ```
    # amd64 / arm64
    WANWU_ARCH=amd64
    
    # external ip port（注意如果浏览器访问非localhost部署的万悟，则需要修改localhost为对外ip，例如192.168.xx.xx）
    WANWU_EXTERNAL_IP=localhost
    ```

    1.3 配置.env文件中的`WANWU_BFF_JWT_SIGNING_KEY`变量，一串自定义复杂随机字符串，用于生成jwt token

    ```
    # bff
    WANWU_BFF_JWT_SIGNING_KEY=
    ```

    1.4 创建docker运行网络
    ```
    docker network create wanwu-net
    ```

2. 启动服务（首次运行会自动从Docker Hub拉取镜像）
    ```bash
    # amd64系统执行:
    docker compose --env-file .env --env-file .env.image.amd64 up -d
    # arm64系统执行:
    docker compose --env-file .env --env-file .env.image.arm64 up -d
    ```

3. 登录系统：http://localhost:8081
    ```
    默认用户：admin
    默认密码：Wanwu123456
    ```

4. 关闭服务
    ```bash
    # amd64系统执行:
    docker compose --env-file .env --env-file .env.image.amd64 down
    # arm64系统执行:
    docker compose --env-file .env --env-file .env.image.arm64 down
    ```

- **源码启动（开发）**

1. 基于上述Docker安装步骤，将系统服务完整启动

2. 以后端bff-service服务为例

    2.1 停止bff-service
    ```
    make -f Makefile.develop stop-bff
    ```

    2.2 编译bff-service可执行文件
    ```
    # amd64系统执行:
    make build-bff-amd64
    # arm64系统执行:
    make build-bff-arm64
    ```

    2.3 启动bff-service
    ```
    make -f Makefile.develop run-bff
    ```

------

### ⬆️ 版本升级

1. 基于上述Docker安装步骤，将系统服务完整停止

2. 更新至最新版本代码

    2.1 wanwu仓库目录内，更新代码
    ```bash
    # 切换到main分支
    git checkout main
    # 拉取最新代码
    git pull
    ```

    2.2 重新拷贝环境变量文件（如果有环境变量修改，请自行重新修改）
    ```bash
    # 备份当前.env文件
    cp .env .env.old
    # 拷贝.env文件
    cp .env.bak .env
    ```

3. 基于上述Docker安装步骤，将系统服务完整启动

------

### &#x1F4D1; 使用万悟

为了帮助您快速上手本项目，我们强烈推荐先查看[ 文档操作手册](https://github.com/UnicomAI/wanwu/tree/main/configs/microservice/bff-service/static/manual)。我们为用户提供了交互式、结构化的操作指南，您可以直接在其中查看操作说明、接口文档等，极大地降低了学习和使用的门槛。详细功能清单如下：

|                             功能                             |                           详细描述                           |
| :----------------------------------------------------------: | :----------------------------------------------------------: |
| [模型管理](https://github.com/UnicomAI/wanwu/blob/main/configs/microservice/bff-service/static/manual/1.%E6%A8%A1%E5%9E%8B%E7%AE%A1%E7%90%86.md) | 支持用户导入包括联通元景、OpenAI-API-compatible、Ollama、通义千问、火山引擎等模型供应商的LLM、Embedding、Rerank模型。[ 模型导入方式-详细版](https://github.com/UnicomAI/wanwu/blob/main/configs/microservice/bff-service/static/manual/%E6%A8%A1%E5%9E%8B%E5%AF%BC%E5%85%A5%E6%96%B9%E5%BC%8F-%E8%AF%A6%E7%BB%86%E7%89%88.md) |
| [知识库](https://github.com/UnicomAI/wanwu/tree/main/configs/microservice/bff-service/static/manual/2.%E7%9F%A5%E8%AF%86%E5%BA%93) | 在文档解析能力方面:支持12种文件类型的上传，支持ur解析;文档解析方式支持OCR和[ 高精度模型解析(标题/表格/公式)](https://github.com/UnicomAI/DocParserServer/tree/main)，文档分段设置支持通用分段和父子分段。在调优能力方面:支持元数据管理及元数据过滤查询，支持分段内容增删改，支持对分段设置关键词标签提升召回效果，支持分段启停操作，支持命中测试等功能。在检索能力方面:支持向量检索、全文检索、混合检索多种检索模式;在问答能力方面:支持自动引用出处，支持图文并茂的生成答案 |
| [资源库](https://github.com/UnicomAI/wanwu/blob/main/configs/microservice/bff-service/static/manual/3.%E5%B7%A5%E5%85%B7%E5%B9%BF%E5%9C%BA.md) | 同时支持导入自己的MCP服务或自定义工具，并在工作流和智能体中使用 |
| [安全护栏](https://github.com/UnicomAI/wanwu/blob/main/configs/microservice/bff-service/static/manual/4.%E5%AE%89%E5%85%A8%E6%8A%A4%E6%A0%8F.md) |        用户可以创建敏感词表，控制模型反馈结果的安全性        |
| [文本问答](https://github.com/UnicomAI/wanwu/blob/main/configs/microservice/bff-service/static/manual/5.%E6%96%87%E6%9C%AC%E9%97%AE%E7%AD%94.md) | 基于私人知识库的专属知识顾问，支持知识库管理、知识问答、知识总结、个性参数配置、安全护栏、检索配置等功能，提高知识管理与学习的效率。支持公开或私密发布文本问答应用，支持发布为API |
| [工作流](https://github.com/UnicomAI/wanwu/tree/main/configs/microservice/bff-service/static/manual/6.%E5%B7%A5%E4%BD%9C%E6%B5%81) | 可以扩展智能体能力边界，由节点组成，提供可视化工作流编辑能力，用户可以编排多个不同的工作流节点，实现复杂且稳定的业务流程。支持公开或私密发布工作流应用，支持发布为API，支持导入导出 |
| [智能体](https://github.com/UnicomAI/wanwu/blob/main/configs/microservice/bff-service/static/manual/7.%E6%99%BA%E8%83%BD%E4%BD%93.md) | 基于用户使用场景和业务需求创建智能体，支持选模型、设置提示词、联网检索、知识库选择、MCP、工作流、自定义工具等。支持公开或私密发布智能体应用，支持发布为API和Web Url |
| [应用广场](https://github.com/UnicomAI/wanwu/blob/main/configs/microservice/bff-service/static/manual/8.%E5%BA%94%E7%94%A8%E5%B9%BF%E5%9C%BA.md) |    支持用户体验已发布的应用，包括文本问答、工作流和智能体    |
| [MCP广场](https://github.com/UnicomAI/wanwu/blob/main/configs/microservice/bff-service/static/manual/9.MCP%E5%B9%BF%E5%9C%BA.md) |             内置100+优选行业MCP server，即选即用             |
| [设置](https://github.com/UnicomAI/wanwu/blob/main/configs/microservice/bff-service/static/manual/9.%E8%AE%BE%E7%BD%AE.md) | 平台支持多租户，允许用户进行组织、角色、用户管理、平台基础配置 |

------

### &#x1F4F0; TODO LIST

- [ ] 多模态模型接入
- [ ] 支持自定义MCP Server，即可以把工作流、智能体、或者符合OpenAPI规范的API作为tools添加到MCP Server里进行发布
- [ ] 知识库共享
- [ ] 智能体和模型测评
- [ ] 智能体监控统计
- [ ] 模型体验
- [ ] 提示词工程

------

### &#128172; Q & A

- **【Q】Linux系统Elastic(elastic-wanwu)启动报错：Memory limited without swap.**

    【A】关闭服务，执行 `sudo sysctl -w vm.max_map_count=262144` 后，重启服务

- **【Q】系统服务正常启动后，mysql-wanwu-setup和elastic-wanwu-setup容器退出：状态码为Exited (0)**

    【A】正常，这两个容器用于完成一些初始化任务，执行完成后会自动退出

- **【Q】模型导入相关**

    【A】以导入联通元景LLM为例（导入OpenAI-API-compatible或导入Embedding、Rerank类型类似）：
    ```
    1. 联通元景MaaS云LLM的Open API接口例如：https://maas.ai-yuanjing.com/openapi/compatible-mode/v1/chat/completions
    
    2. 用户在联通元景MaaS云上申请到的API Key形如：sk-abc********************xyz
    
    3. 确认API与Key可正确请求LLM，以请求yuanjing-70b-chat为例：
        curl --location 'https://maas.ai-yuanjing.com/openapi/compatible-mode/v1/chat/completions' \
        --header 'Content-Type: application/json' \
        --header 'Accept: application/json' \
        --header 'Authorization: Bearer sk-abc********************xyz' \
        --data '{
                "model": "yuanjing-70b-chat",
                "messages": [{
                        "role": "user",
                        "content": "你好"
                }]
        }'
    
    4. 导入模型：
    4.1【模型名称】必须为上述curl中可以正确请求的model；例如 yuanjing-70b-chat
    4.2【API Key】必须为上述curl中可以正确请求的key；例如 sk-abc********************xyz（注意不填Bearer前缀）
    4.3【推理URL】必须为上述curl中可以正确请求的url；例如 https://maas.ai-yuanjing.com/openapi/compatible-mode/v1（注意不带 /chat/completions 后缀）
    
    5. 导入Embedding模型同上述导入LLM，注意推理URL不带 /embeddings 后缀
    
    6. 导入Rerank模型同上述导入LLM，注意推理URL不带 /rerank 后缀
    ```

------

### &#x1F517; 致谢

- [Coze](https://github.com/coze-dev)
- [LangChain](https://github.com/langchain-ai/langchain)
- [Qwen-Agent](https://github.com/QwenLM/Qwen-Agent)

------

### ⚖️ 许可证
元景万悟智能体平台根据Apache License 2.0发布。

------

### &#x1F4E9; 联系我们
| QQ 群1(已满):490071123                                    | QQ 群2:1026898615                                         |
| ------------------------------------------------------------ | ------------------------------------------------------------ |
| <img width="183" height="258" alt="image" src="https://github.com/user-attachments/assets/010f1d68-78e9-446d-baf1-0a7339efb48e" /> | <img width="183" height="258" alt="image" src="https://github.com/user-attachments/assets/10796f69-5c18-4f21-adbb-b22b6ef88df2" /> |