import moment from 'moment';
export function dateFormat(dateValue) {
    return moment(dateValue).format('YYYY-MM-DD HH:mm:ss');
}
export function todDateFormat(dateString) {
    return moment(dateString);
}
