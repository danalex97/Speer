import React, { Component } from 'react';
import * as d3 from "d3";

export default class BarChart extends Component {
	constructor(props) {
		super(props);
	}

	componentDidMount() {
		const margin = this.props.margin;
		const height = this.props.size[1] - 2 * margin;
		const width  = this.props.size[0] - 2 * margin;

		// build the chart
		const svg   = d3.select('#' + this.props.id);
		const chart = svg.append('g')
				.attr('transform', `translate(
				${margin},
				${margin}
			)`);

		const data       = this.props.data;
		const windowSize = this.props.windowSize;

		const dataMax    = d3.max(data);
		const dataSize   = data.length * windowSize;

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

		// set values
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
            .attr('x', (d, i) => xScale(i * windowSize))
            .attr('y', d => yScale(d))
            .attr('height', d => height - yScale(d))
            .attr('width', xScale(windowSize));
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
