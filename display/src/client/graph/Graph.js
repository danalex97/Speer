import * as d3 from "d3";

import React, { Component } from 'react';
import ReactDOM from 'react-dom';

import Node from './Node';
import Edge from './Edge';

export default class Graph extends Component {
  constructor(props) {
    super(props);

    const nodes   = this.props.nodes.map(x => Object.assign(x, {
      x : 0,
      y : 0
    }));
    const links   = this.props.links.map(x => Object.assign(x, {
      source : x.src,
      target : x.dst,
    }));
    const mapping = new Map(nodes.map(x => [x.node, x]));

    this.state = {
      nodes   : nodes,
      links   : links,
      mapping : mapping,
    };

    this.update = this.update.bind(this);
  }

  componentDidMount() {
    setInterval(this.update, 30);
  }

  update() {
    const forceLink = d3
      .forceLink(this.state.links)
      .id(d => d.node)
      .distance(d => 100);

    const force = d3.forceSimulation()
      .nodes(this.state.nodes)
      .force("charge", d3.forceManyBody())
      .force("link", forceLink)
      .force('center', d3.forceCenter(0, 0))
      .tick(1);

    this.state.nodes.map(d => {
      const ctx = this.props;
      const fact = 0.4;

      if (d.x < -ctx.width * fact) { d.x = -ctx.width * fact; }
      if (d.x >  ctx.width * fact) { d.x =  ctx.width * fact; }
      if (d.y < -ctx.height * fact) { d.y = -ctx.height * fact; }
      if (d.y >  ctx.height * fact) { d.y =  ctx.height * fact; }

      return Object.assign(d, {
        x : d.x,
        y : d.y,
        vx : 0,
        vy : 0
      });
    });

    const mapping = new Map(this.state.nodes.map(x => [x.node, x]));

    this.setState({
      nodes : this.state.nodes,
      links : this.state.links,
      mapping : mapping
    });
  }

  render() {
    const nodes = this.state.nodes.map(node => (<Node
      data={node}
      key={`node#${node.node}`}
      x={node.x}
      y={node.y}
    />));

    const mapping = this.state.mapping;
    const edges = this.state.links.map(edge => {
      const a = mapping.get(edge.src);
      const b = mapping.get(edge.dst);

      return (<Edge
        key={`edge#${edge.src}#${edge.tgt}#${edge.time}`}
        x1={a.x} y1={a.y}
        x2={b.x} y2={b.y}
      />);
    });

    return (<g>
      {nodes}
      {edges}
    </g>);
  }
}
