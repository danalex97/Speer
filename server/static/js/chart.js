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

	setStyle(chart) {
		// d3.select(this.node)
		// 	.selectAll('rect')
		// 	.data(this.props.data)
		// 	.style('fill', '#fe9922')
		// 	.attr('x', (d, i) => i * height)
		// 	.attr('y', d => this.props.size[1] - yScale(d))
		// 	.attr('height', d => yScale(d))
		// 	.attr('width', height);
	}

	updateBarChart() {
		// d3.select(node)
		// 	.selectAll('rect')
		// 	.data(this.props.data)
		// 	.enter()
		// 	.append('rect');
		// d3.select(node)
		// 	.selectAll('rect')
		// 	.data(this.props.data)
		// 	.exit()
		// 	.remove();

	}

	createBarChart() {
		// const node = this.node;
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

		const dataMax  = d3.max(this.props.data);

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

		this.setStyle(chart);
	}

	render() {
		return (<svg
			ref={node => this.node = node}
			margin={this.props.margin}
			width={this.props.size[0]}
			height={this.props.size[1]}
		/>);
	}
}
