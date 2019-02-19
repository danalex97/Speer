import React, { Component } from 'react';
import * as d3 from "d3";

import Graph from './Graph';

export default class Canvas extends Component {
	constructor(props) {
		super(props);
  }

	render() {
    const height = this.props.size[1];
    const width  = this.props.size[0];

    const data  = this.props.data;
    const nodes = data.filter(x => x.node != null);
    const links = data.filter(x => x.src != null);

    const graph = <Graph
      nodes={nodes}
      links={links}
      width={width}
      height={height}
    />;

    const translate = `translate(${width / 2}, ${height / 2})`;

		return (<svg
        width={width}
        height={height}>
      <g transform={translate}>
        {graph}
  		</g>
    </svg>);
	}
}
