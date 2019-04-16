import React, { Component } from 'react';
import Textbox from '../components/Textbox';

const SERVER = "http://localhost:8080";

export default class LogUploader extends Component {
  constructor(props) {
    super(props);
    this.state = {name: ''};

    this.handleChange = this.handleChange.bind(this);
    this.handleClick = this.handleClick.bind(this);
  }

  handleChange(event) {
    this.setState({name : event.target.value});
  }

  handleClick(event) {
    this.props.handleLogUpdate(this.state.name);
  }

  render() {
    return <Textbox
      placeholder={this.props.placeholder}
      name={this.state.name}
      handleChange={this.handleChange}
      handleClick={this.handleClick}
      submitText="Load Log"
    />;
  }
}
