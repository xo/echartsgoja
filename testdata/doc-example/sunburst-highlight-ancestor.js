"use strict";
option = {
    silent: true,
    series: {
        radius: ['15%', '95%'],
        center: ['50%', '60%'],
        type: 'sunburst',
        sort: null,
        emphasis: {
            focus: 'ancestor'
        },
        data: [{
                value: 10,
                children: [{
                        value: 4,
                        children: [{
                                value: 2,
                                children: [{
                                        name: 'target',
                                        value: 1
                                    }]
                            }, {
                                value: 1
                            }, {
                                value: 0.5
                            }]
                    }, {
                        value: 2
                    }]
            }, {
                value: 4,
                children: [{
                        children: [{
                                value: 2
                            }]
                    }]
            }],
        label: {
            normal: {
                rotate: 'none',
                color: '#fff'
            }
        },
        levels: [],
        itemStyle: {
            color: 'yellow',
            borderWidth: 2
        },
        emphasis: {
            itemStyle: {
                color: 'red'
            }
        },
        highlight: {
            itemStyle: {
                color: 'orange'
            }
        },
        downplay: {
            itemStyle: {
                color: '#ccc'
            }
        }
    }
};
setTimeout(function () {
    myChart.dispatchAction({
        type: 'sunburstHighlight',
        targetNodeId: 'target'
    });
});
