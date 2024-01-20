Date.prototype.setYear = Date.prototype.setFullYear;

function parse(spec) {
  if (typeof spec == "object") {
    return spec;
  }
  return JSON.parse(spec);
}

function version() {
  return echarts.version;
}

function render_options(width, height, opts) {
  var chart = echarts.init(null, null, {
    renderer: "svg",
    ssr: true,
    width: width,
    height: height,
  });
  chart.setOption(parse(opts));
  const s = chart.renderToSVGString();
  chart.dispose();
  chart = null;
  return s;
}
