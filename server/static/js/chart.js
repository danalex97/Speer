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
		const height = this.state.height;
		const width  = this.state.width;
		const rectWidth = 10;

		const data = this.props.data;
		const dataSize = this.props.dataSize;

		// update axis
		const rhs = Math.max(dataSize, data.length);
		const lhs = rhs - dataSize;

		const xScale = d3.scaleLinear()
			.range([0, width])
			.domain([lhs, rhs]);
		// TODO: change numbers on bottom axis
		chart.selectAll("g.x.axis")
			.call(d3.axisBottom(xScale));

		// keep yScale
		const yScale = this.state.yScale;

		// update values
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
		const svg   = d3.select('#' + this.props.id);
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
			.attr("class", "x axis")
			.attr('transform', `translate(0, ${height})`)
			.call(d3.axisBottom(xScale));

		this.setState({
			chart  : chart,
			yScale : yScale,
			width  : width,
			height : height,
		});
	}

	render() {
		return (<svg
			id={this.props.id}
			margin={this.props.margin}
			width={this.props.size[0]}
			height={this.props.size[1]}
		/>);
	}
}
