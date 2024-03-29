"use strict";
option = {
    baseOption: {
        timeline: {
            axisType: 'category',
            // realtime: false,
            // loop: false,
            autoPlay: false,
            // currentIndex: 2,
            playInterval: 1000,
            controlStyle: {
                stopIcon: 'image://data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAAAXNSR0IArs4c6QAAAGpJREFUWAntlkEKwDAIBKP+/znt82xrD6FIpSSHXDpCQBeXDXOyNQoCBQEz21T1eL7QivUuj/qkO1MTwUm6R3cvPbEw6tO3kJUaH4AABCAAAQhAAAIQgEBJQET2fJxe2udZPuvLWcz/IXACuUMWMJDzUAgAAAAASUVORK5CYII=',
                playIcon: 'path://M41.365908,29.4271388 L41.3664843,29.4265626 L26.3794329,19.1497136 L26.3747509,19.1541315 C26.0642269,18.8592621 25.6429678,18.677793 25.1786824,18.677793 C24.2236284,18.677793 23.4494433,19.4443188 23.4494433,20.3905371 C23.4494433,20.910214 23.4270417,21.9276946 23.4494433,21.9056292 L23.4494433,30.6673861 L23.4494433,39.8901629 C23.4494433,39.8977982 23.4494433,40.4825908 23.4494433,40.9444991 C23.4494433,41.8901412 24.2236284,42.656691 25.1786824,42.656691 C25.6447205,42.656691 26.0677564,42.4740454 26.3782564,42.1764869 L26.3794329,42.1770872 L41.3664843,31.9005503 L41.3659081,31.8996379 C41.6917266,31.5882735 41.894997,31.1514078 41.894997,30.6670739 C41.894997,30.6658974 41.894997,30.6650091 41.894997,30.6635444 C41.894997,30.6623679 41.894997,30.6609273 41.894997,30.6600389 C41.894997,30.175657 41.6917265,29.7384792 41.365908,29.4271388 Z'
            },
            replaceMerge: 'series',
            data: [
                '2 series', '3 series', '1 series'
            ],
        },
        tooltip: {
            trigger: 'axis',
            axisPointer: {
                type: 'shadow'
            }
        },
        legend: {},
        calculable: true,
        grid: {
            top: 80, bottom: 100
        },
        toolbox: {
            left: 'center',
            top: 30,
            feature: {
                dataZoom: {}
            }
        },
        xAxis: {
            type: 'category',
            data: ['CityB', 'CityT', 'CityH', 'CityS'],
            splitLine: { show: false }
        },
        yAxis: [
            {
                type: 'value',
                name: 'GDP'
            }
        ],
        series: []
    },
    options: [
        {
            series: [
                { name: 'a', type: 'bar', data: [12, 33, 44, 55] },
                { name: 'b', type: 'bar', data: [55, 66, 77, 88] },
            ]
        },
        {
            series: [
                { name: 'a', type: 'bar', data: [22, 33, 44, 55] },
                { name: 'b', type: 'bar', data: [55, 66, 77, 88] },
                { name: 'c', type: 'bar', data: [55, 66, 77, 88] }
            ]
        },
        {
            series: [
                { name: 'b', type: 'bar', data: [55, 66, 77, 88] }
            ]
        }
    ]
};
