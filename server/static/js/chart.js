class BarChart extends React.Component {
	constructor(props) {
		super(props);

		this.createBarChart = this.createBarChart.bind(this);
	}

	componentDidMount() {
    	this.createBarChart();
    }

	componentDidUpdate() {
		this.updateBarChart();
	}

	updateBarChart() {
		const node  = this.node;

		const chart = this.state.chart;
		const xScale = this.state.xScale;
		const yScale = this.state.yScale;

		const height = this.state.height;
		const width  = this.state.width;

		const data = this.props.data;

		const rectWidth = 10;

		console.log(data);

		chart.selectAll('rect')
            .data(data)
            .enter()
            .append('rect');

		chart.selectAll('rect')
            .data(data)
            .exit()
            .remove();

		chart.selectAll('rect')
            .data(data)
            .style('fill', 'black')
            .attr('x', (d, i) => i * rectWidth)
            .attr('y', d => yScale(d))
            .attr('height', d => height - yScale(d))
            .attr('width', rectWidth);
	}

	createBarChart() {
		const margin   = this.props.margin;
		const height   = this.props.size[1] - 2 * margin;
		const width    = this.props.size[0] - 2 * margin;
		const dataSize = this.props.dataSize;

		// build the chart
		const svg   = d3.select('svg');
		const chart = svg.append('g')
    		.attr('transform', `translate(
				${margin},
				${margin}
			)`);

		const dataMax = Math.max(
			d3.max(this.props.data),
			100, // some minimum metric value
		);

		// get the scales
		const yScale = d3.scaleLinear()
			.range([height, 0])
			.domain([0, dataMax]);
		const xScale = d3.scaleLinear()
			.range([0, width])
			.domain([0, dataSize]);

		// build the axis
		chart.append('g')
	    	.call(d3.axisLeft(yScale));
		chart.append('g')
		    .attr('transform', `translate(0, ${height})`)
			.call(d3.axisBottom(xScale));

		this.setState({
			chart  : chart,
			xScale : xScale,
			yScale : yScale,
			width  : width,
			height : height,
		});
	}

	render() {
		if (this.state != null) {
			this.updateBarChart();
		}

		return (<svg
			margin={this.props.margin}
			width={this.props.size[0]}
			height={this.props.size[1]}
		/>);
	}
}
