import * as d3 from "d3";

import React, { Component } from 'react';
import ReactDOM from 'react-dom';

import Node from './Node';

export default class Graph extends Component {
  constructor(props) {
    super(props);

    this.state = {
      nodes : this.props.nodes,
    };
  }

  // componentDidMount() {
  //   this.setState(Object.assign(this.state, {
  //     node : d3.select(ReactDOM.findDOMNode(this))
  //   }));
  // }

  render() {
    const nodes = this.state.nodes.map(node => (<Node
      data={node}
      key={`node#${node.node}`}
      x={node.x}
      y={node.y}
    />));

    return (<g>
      {nodes}
    </g>);
  }
}
