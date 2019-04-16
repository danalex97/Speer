import * as d3 from "d3";

import React, { Component } from 'react';
import ReactDOM from 'react-dom';

export default class Node extends Component {
  constructor(props) {
    super(props);

    this.state = {
      data : this.props.data,
    };
  }

  render() {
    const data = this.state.data;
    const x = this.props.x;
    const y = this.props.y;

    return (<g className='node'>
      <circle
        cx={x}
        cy={y}
        r='20'/>
      <text>{data.name}</text>
    </g>);
  }
}
