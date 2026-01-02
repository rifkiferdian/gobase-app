/*
Template Name: Skote - Admin & Dashboard Template
Author: Themesbrand
Website: https://themesbrand.com/
Contact: themesbrand@gmail.com
File: Dashboard Init Js File
*/

// get colors array from the string
function getChartColorsArray(chartId) {
    if (document.getElementById(chartId) !== null) {
        var colors = document.getElementById(chartId).getAttribute("data-colors");
        if (colors) {
            colors = JSON.parse(colors);
            return colors.map(function (value) {
                var newValue = value.replace(" ", "");
                if (newValue.indexOf(",") === -1) {
                    var color = getComputedStyle(document.documentElement).getPropertyValue(newValue);
                    if (color) return color;
                    else return newValue;;
                } else {
                    var val = value.split(',');
                    if (val.length == 2) {
                        var rgbaColor = getComputedStyle(document.documentElement).getPropertyValue(val[0]);
                        rgbaColor = "rgba(" + rgbaColor + "," + val[1] + ")";
                        return rgbaColor;
                    } else {
                        return newValue;
                    }
                }
            });
        }
    }
}

//  subscribe modal

setTimeout(function () {
    $('#subscribeModal').modal('show');
}, 2000);

// column chart
var columnChartColors = getChartColorsArray("column_chart");
if (columnChartColors) {
    var options = {
        chart: {
            height: 350,
            type: 'bar',
            toolbar: {
                show: false,
            }
        },
        plotOptions: {
            bar: {
                horizontal: false,
                columnWidth: '45%',
                endingShape: 'rounded'
            },
        },
        dataLabels: {
            enabled: false
        },
        stroke: {
            show: true,
            width: 2,
            colors: ['transparent']
        },
        series: [{
            name: 'Net Profit',
            data: [46, 57, 59, 54, 62, 58, 64, 60, 66]
        }, {
            name: 'Revenue',
            data: [74, 83, 102, 97, 86, 106, 93, 114, 94]
        }, {
            name: 'Free Cash Flow',
            data: [37, 42, 38, 26, 47, 50, 54, 55, 43]
        }],
        colors: columnChartColors,
        xaxis: {
            categories: ['Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct'],
        },
        yaxis: {
            title: {
                text: '$ (thousands)',
                style: {
                    fontWeight: '500',
                },
            }
        },
        grid: {
            borderColor: '#f1f1f1',
        },
        fill: {
            opacity: 1

        },
        tooltip: {
            y: {
                formatter: function (val) {
                    return "$ " + val + " thousands"
                }
            }
        }
    }

    var chart = new ApexCharts(
        document.querySelector("#column_chart"),
        options
    );

    chart.render();
}


// stacked column chart
var linechartBasicColors = getChartColorsArray("stacked-column-chart");
if (linechartBasicColors) {
    var options = {
        chart: {
            height: 360,
            type: 'bar',
            stacked: true,
            toolbar: {
                show: false
            },
            zoom: {
                enabled: true
            }
        },

        plotOptions: {
            bar: {
                horizontal: false,
                columnWidth: '15%',
                endingShape: 'rounded',
                borderRadius: 5,
            },
        },

        dataLabels: {
            enabled: false
        },
        series: [{
            name: 'Series A',
            data: [44, 55, 41, 67, 22, 43, 36, 52, 24, 18, 36, 48]
        }, {
            name: 'Series B',
            data: [13, 23, 20, 8, 13, 27, 18, 22, 10, 16, 24, 22]
        }, {
            name: 'Series C',
            data: [11, 17, 15, 15, 21, 14, 11, 18, 17, 12, 20, 18]
        }],
        xaxis: {
            categories: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
        },
        colors: linechartBasicColors,
        legend: {
            position: 'bottom',
        },
        fill: {
            opacity: 1
        },
    }

    var chart = new ApexCharts(
        document.querySelector("#stacked-column-chart"),
        options
    );

    chart.render();
}

// Radial chart
var radialbarColors = getChartColorsArray("radialBar-chart");
if (radialbarColors) {
    var options = {
        chart: {
            height: 200,
            type: 'radialBar',
            offsetY: -10
        },
        plotOptions: {
            radialBar: {
                startAngle: -135,
                endAngle: 135,
                dataLabels: {
                    name: {
                        fontSize: '13px',
                        color: undefined,
                        offsetY: 60
                    },
                    value: {
                        offsetY: 22,
                        fontSize: '16px',
                        color: undefined,
                        formatter: function (val) {
                            return val + "%";
                        }
                    }
                }
            }
        },
        colors: radialbarColors,
        fill: {
            type: 'gradient',
            gradient: {
                shade: 'dark',
                shadeIntensity: 0.15,
                inverseColors: false,
                opacityFrom: 1,
                opacityTo: 1,
                stops: [0, 50, 65, 91]
            },
        },
        stroke: {
            dashArray: 4,
        },
        series: [67],
        labels: ['Series A'],

    }

    var chart = new ApexCharts(
        document.querySelector("#radialBar-chart"),
        options
    );

    chart.render();
}
