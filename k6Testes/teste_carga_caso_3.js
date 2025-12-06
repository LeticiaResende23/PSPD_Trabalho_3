import http from 'k6/http';
import { check, sleep } from 'k6';

const url = 'http://127.0.0.1:38491/'

const intervalo = '20s'
const carga = 400
const tempo_de_carga = '5m'

const threashold = '95' //Threashold em %
const threashold_ms = '20000' //tempo de resposta

const tolerancia = '0.1'

export const options = {
    stages: [
        { duration: intervalo, target: carga }, 
        { duration: tempo_de_carga, target: carga }, 
        { duration: intervalo, target: 0 },  
    ],

    thresholds: {
        'http_req_duration': ['p('+ threashold +')<'+ threashold_ms], 
        'http_req_failed': ['rate<' + tolerancia], 
    },
    failOnThreshold: true,
};

export default function () {
    const res1 = http.get(url);
    const res2 = http.get(url + 'produtos/1')

    check(res1, {
        'status is 200': (r) => r.status === 200,
    });

    check(res2, {
        'status is 200': (r) => r.status === 200,
    });
    sleep(1); 
}