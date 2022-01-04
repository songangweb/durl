import requestAPI from '@/utils/ajax-request.js';
// 获取短链接列表
export function getShortChainArr(params) {
    return requestAPI(`/request/url/list`, params);
}
// 添加短链接
export function addShortChainValue(params) {
    return requestAPI(`/request/url`, params, 'POST');
}
// 修改短链接
export function changeShortChainValue(params, id) {
    return requestAPI(`/request/url/${id}`, params, 'PUT');
}
// 批量删除
export function batchDeleteArr(params) {
    return requestAPI(`/request/url`, params, 'DELETE');
}
// 删除数据
export function deleteValue(id) {
    return requestAPI(`/request/url/${id}`, 'DELETE');
}
// 批量冻结/解冻
export function batchFreezeArr(params) {
    return requestAPI(`/request/url/frozen`, params, 'PUT');
}
