import axios from 'axios';

const serverUrl = 'http://localhost:7861';

const createUrl = (url) => `${serverUrl}${url}`;

export default {
    get: (url, params) => axios.get(createUrl(url), { params }),
    post: (url, data) => axios.post(createUrl(url), data),
    put: (url, data) => axios.put(createUrl(url), data),
    delete: (url, data) => axios.delete(createUrl(url), data),
}