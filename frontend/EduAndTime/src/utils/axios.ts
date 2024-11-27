import axios from 'axios';

const authInstance = axios.create({
    baseURL: 'http://localhost:8080', 
    timeout: 10000,
});
const libreriaInstancia= axios.create({
    baseURL: 'http://localhost:8081', 
    timeout: 10000,
});


export {authInstance, libreriaInstancia};
