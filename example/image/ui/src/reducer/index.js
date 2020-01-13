const initialState = {
    isRunning: false,
    refreshRate: '0',
    updateCounter: 0
};

function reducer(state = initialState, action) {
    switch (action.type) {
        case 'FETCH_STATS_SUCCEEDED':
            return Object.assign({}, state, {
                isRunning: action.payload.isRunning
            });
        case 'TOOLBAR_START':
            return Object.assign({}, state, {
                isRunning: true
            });
        case 'TOOLBAR_SELECT':
            return Object.assign({}, state, {
                refreshRate: action.payload.refreshRate
            });
        case 'TOOLBAR_UPDATE':
            return Object.assign({}, state, {
                updateCounter: state.updateCounter + 1
            });
        default:
            return state;
    }
}

export default reducer;