import requestAPI from '@/utils/ajax-request.js';
let request = ""
// 短链接
// 获取短链接列表
if (process.env.NODE_ENV == 'development') {
    request = "/request"
}
export function getShortUrlList(params) {
    return requestAPI(`${request}/url/list`, params);
}
// 添加短链接
export function addShortUrl(params) {
    return requestAPI(`${request}/url`, params, 'POST');
}
// 修改短链接
export function changeShortUrl(params, id) {
    return requestAPI(`${request}/url/${id}`, params, 'PUT');
}
// 批量删除
export function batchDeleteArr(params) {
    return requestAPI(`${request}/url`, params, 'DELETE');
}
// 删除数据
export function deleteUrl(id) {
    return requestAPI(`${request}/url/${id}`, '', 'DELETE');
}
// 批量冻结/解冻
export function batchFreezeArr(params) {
    return requestAPI(`${request}/url/frozen`, params, 'PUT');
}
// 单个冻结url
export function freezeUrl(id) {
    return requestAPI(`${request}/url/frozen/${id}`, '', 'PUT');
}

// 短链接详情
export function urlInfo(id) {
    return requestAPI(`${request}/url/${id}`, '', 'GET');
}


//  黑名单
// 获取黑名单列表
export function getBlackListArr(params) {
    return requestAPI(`${request}/blacklist/list`, params);
}
// 添加黑名单ip
export function addBlacklist(params) {
    return requestAPI(`${request}/blacklist`, params, 'POST');
}
// 修改黑名单ip
export function changeBlacklist(params, id) {
    return requestAPI(`${request}/blacklist/${id}`, params, 'PUT');
}
// 删除黑名单v
export function deleteBlacklist(id) {
    return requestAPI(`${request}/blacklist/${id}`, '', 'DELETE');
}



