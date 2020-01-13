import axios from 'axios';

export function fetchStatsSucceeded(stats) {
    return {
        type: 'FETCH_STATS_SUCCEEDED',
        payload: {
            isRunning: stats.isRunning
        }
    }
}

export function fetchStats() {
    return dispatch => {
        axios.get('http://localhost:3001/stats')
            .then(resp => {
                dispatch(fetchStatsSucceeded(resp.data))
            })
    }
}

export function start() {
    return {
        type: 'TOOLBAR_START'
    };
}

export function sendStart() {
    return dispatch => {
        axios.get('http://localhost:3001/start')
            .then(resp => {
                dispatch(start())
            })
    }
}

export function update() {
    return {
        type: 'TOOLBAR_UPDATE'
    }
}

export function select(rate) {
    return {
        type: 'TOOLBAR_SELECT',
        payload: { refreshRate: rate }
    }
}