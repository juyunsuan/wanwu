const state = {
  // 表格跨页选中的数据
  selectedRows: [],
  // 当前选中的文档ID列表
  selectedDocIds: [],
  // 表格数据
  tableData: [],
  // 分页信息
  pagination: {
    currentPage: 1,
    pageSize: 10,
    total: 0
  },
  // 搜索条件
  searchParams: {},
  // 批量操作状态
  batchOperationStatus: false
}

const mutations = {
  // 设置选中的行数据
  SET_SELECTED_ROWS(state, rows) {
    console.log(rows,'rows')
    state.selectedRows=rows
  },
  
  // 添加选中的行
  ADD_SELECTED_ROWS(state, rows) {
    // 去重添加
    const existingIds = state.selectedRows.map(row => row.docId)
    const newRows = rows.filter(row => !existingIds.includes(row.docId))
    state.selectedRows = [...state.selectedRows, ...newRows]
  },
  
  // 移除选中的行
  REMOVE_SELECTED_ROWS(state, rowIds) {
    state.selectedRows = state.selectedRows.filter(row => !rowIds.includes(row.docId))
  },
  
  // 清空选中的行
  CLEAR_SELECTED_ROWS(state) {
    state.selectedRows = []
  },
  
  // 设置选中的文档ID列表
  SET_SELECTED_DOC_IDS(state, docIds) {
    state.selectedDocIds = docIds
  },
  
  // 添加选中的文档ID
  ADD_SELECTED_DOC_IDS(state, docIds) {
    // 去重添加
    const newIds = docIds.filter(id => !state.selectedDocIds.includes(id))
    state.selectedDocIds = [...state.selectedDocIds, ...newIds]
  },
  
  // 移除选中的文档ID
  REMOVE_SELECTED_DOC_IDS(state, docIds) {
    state.selectedDocIds = state.selectedDocIds.filter(id => !docIds.includes(id))
  },
  
  // 清空选中的文档ID
  CLEAR_SELECTED_DOC_IDS(state) {
    state.selectedDocIds = []
  },
  
  // 设置表格数据
  SET_TABLE_DATA(state, data) {
    state.tableData = data
  },
  
  // 更新表格数据中的某一行
  UPDATE_TABLE_ROW(state, { index, rowData }) {
    if (index >= 0 && index < state.tableData.length) {
      state.tableData.splice(index, 1, rowData)
    }
  },
  
  // 删除表格数据中的某一行
  DELETE_TABLE_ROW(state, index) {
    if (index >= 0 && index < state.tableData.length) {
      state.tableData.splice(index, 1)
    }
  },
  
  // 设置分页信息
  SET_PAGINATION(state, pagination) {
    state.pagination = { ...state.pagination, ...pagination }
  },
  
  // 设置搜索参数
  SET_SEARCH_PARAMS(state, params) {
    state.searchParams = { ...state.searchParams, ...params }
  },
  
  // 清空搜索参数
  CLEAR_SEARCH_PARAMS(state) {
    state.searchParams = {}
  },
  
  // 设置批量操作状态
  SET_BATCH_OPERATION_STATUS(state, status) {
    state.batchOperationStatus = status
  }
}

const actions = {
  // 选择行
  selectRows({ commit }, rows) {
    commit('ADD_SELECTED_ROWS', rows)
    const docIds = rows.map(row => row.docId)
    commit('ADD_SELECTED_DOC_IDS', docIds)
  },
  
  // 取消选择行
  unselectRows({ commit }, rows) {
    const rowIds = rows.map(row => row.docId)
    commit('REMOVE_SELECTED_ROWS', rowIds)
    commit('REMOVE_SELECTED_DOC_IDS', rowIds)
  },
  
  // 切换行选择状态
  toggleRowSelection({ commit, state }, row) {
    const isSelected = state.selectedRows.some(selectedRow => selectedRow.docId === row.docId)
    if (isSelected) {
      commit('REMOVE_SELECTED_ROWS', [row.docId])
      commit('REMOVE_SELECTED_DOC_IDS', [row.docId])
    } else {
      commit('ADD_SELECTED_ROWS', [row])
      commit('ADD_SELECTED_DOC_IDS', [row.docId])
    }
  },
  
  // 全选当前页
  selectAllCurrentPage({ commit, state }, currentPageRows) {
    commit('ADD_SELECTED_ROWS', currentPageRows)
    const docIds = currentPageRows.map(row => row.docId)
    commit('ADD_SELECTED_DOC_IDS', docIds)
  },
  
  // 取消全选当前页
  unselectAllCurrentPage({ commit, state }, currentPageRows) {
    const rowIds = currentPageRows.map(row => row.docId)
    commit('REMOVE_SELECTED_ROWS', rowIds)
    commit('REMOVE_SELECTED_DOC_IDS', rowIds)
  },
  
  // 清空所有选择
  clearAllSelection({ commit }) {
    commit('CLEAR_SELECTED_ROWS')
    commit('CLEAR_SELECTED_DOC_IDS')
  },
  
  // 更新表格数据
  updateTableData({ commit }, data) {
    commit('SET_TABLE_DATA', data)
  },
  
  // 更新分页信息
  updatePagination({ commit }, pagination) {
    commit('SET_PAGINATION', pagination)
  },
  
  // 更新搜索参数
  updateSearchParams({ commit }, params) {
    commit('SET_SEARCH_PARAMS', params)
  },
  
  // 执行批量操作
  async executeBatchOperation({ commit, state }, operation) {
    commit('SET_BATCH_OPERATION_STATUS', true)
    try {
      const result = await operation(state.selectedDocIds, state.selectedRows)
      return result
    } finally {
      commit('SET_BATCH_OPERATION_STATUS', false)
    }
  }
}

const getters = {
  // 获取选中的行数据
  selectedRows: state => state.selectedRows,
  
  // 获取选中的文档ID列表
  selectedDocIds: state => state.selectedDocIds,
  
  // 获取选中的行数量
  selectedCount: state => state.selectedRows.length,
  
  // 获取表格数据
  tableData: state => state.tableData,
  
  // 获取分页信息
  pagination: state => state.pagination,
  
  // 获取搜索参数
  searchParams: state => state.searchParams,
  
  // 获取批量操作状态
  batchOperationStatus: state => state.batchOperationStatus,
  
  // 检查某行是否被选中
  isRowSelected: state => row => {
    return state.selectedRows.some(selectedRow => selectedRow.docId === row.docId)
  },
  
  // 获取当前页选中的行
  currentPageSelectedRows: state => currentPageRows => {
    return state.selectedRows.filter(selectedRow => 
      currentPageRows.some(currentRow => currentRow.docId === selectedRow.docId)
    )
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions,
  getters
}

