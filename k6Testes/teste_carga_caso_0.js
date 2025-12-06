import http from 'k6/http';
import { check, sleep } from 'k6';

const intervalo = '20s'
const carga = 200
const tempo_de_carga = '5m'

const threashold = '95' //Threashold em %
const threashold_ms = '500' //tempo de resposta

const tolerancia = '0.01'

export const options = {
    stages: [
        { duration: intervalo, target: carga }, 
        { duration: tempo_de_carga, target: carga }, 
        { duration: intervalo, target: 0 },  
    ],

    thresholds: {
    },
};

export default function () {
    const res1 = http.get('http://127.0.0.1:33653/');
    const res2 = http.get('http://127.0.0.1:33653/produtos/1')

    check(res1, {
        'status is 200': (r) => r.status === 200,
    });

    check(res2, {
        'status is 200': (r) => r.status === 200,
    });
    sleep(1); 
}