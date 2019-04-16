import * as d3 from "d3";

import React, { Component } from 'react';
import ReactDOM from 'react-dom';

export default class Edge extends Component {
  constructor(props) {
    super(props);
  }

  render() {
    const x1 = this.props.x1;
    const y1 = this.props.y1;
    const x2 = this.props.x2;
    const y2 = this.props.y2;

    return (<g className='link'>
      <line
        stroke="black"
        stroke-width="2px"
        stroke-opacity="0.8"
        x1={x1}
        y1={y1}
        x2={x2}
        y2={y2}
      />
    </g>);
  }
}
