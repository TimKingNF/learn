$( function () {
    echarts.init( document.getElementById( 'radar') ).setOption( optionRadar );
    echarts.init( document.getElementById( 'pie') ).setOption( optionPie );
    echarts.init( document.getElementById( 'line') ).setOption( optionLine );
    $( window).resize( function() {
        echarts.init( document.getElementById( 'radar') ).setOption( optionRadar );
        echarts.init( document.getElementById( 'pie') ).setOption( optionPie );
        echarts.init( document.getElementById( 'line') ).setOption( optionLine );
    });
});