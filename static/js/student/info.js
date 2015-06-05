$( function () {
    echarts.init( document.getElementById( 'radar') ).setOption( optionRadar );
    $( window).resize( function() {
        echarts.init( document.getElementById( 'radar') ).setOption( optionRadar );
    });
});