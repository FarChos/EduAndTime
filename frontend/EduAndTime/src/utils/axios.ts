import axios from 'axios';

const authInstance = axios.create({
    baseURL: 'http://localhost:8080', 
    timeout: 10000,
});
const libreriaInstancia= axios.create({
    baseURL: 'http://localhost:8081', 
    timeout: 10000,
});

authInstance.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('authTokenEAT'); // Recuperar token dinámicamente
        if (token) {
            config.headers.Authorization = token;
        }
        return config;
    },
    (error) => {
        return Promise.reject(
            error instanceof Error ? error : new Error('Request Interceptor Error')
        );
    }
);

libreriaInstancia.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('authTokenEAT'); // Recuperar token dinámicamente
        if (token) {
            config.headers.Authorization = token;
        }
        return config;
    },
    (error) => {
        return Promise.reject(
            error instanceof Error ? error : new Error('Request Interceptor Error')
        );
    }
);

export {authInstance, libreriaInstancia};
