import React, { Component } from 'react';

export default class Textbox extends Component {
  constructor(props) {
    super(props);
  }

  render() {
    return <div className="form-group">
      <input
        type="text"
        placeholder={this.props.placeholder}
        className="formControl"
        value={this.props.name}
        onChange={this.props.handleChange} />
        <button
          type="submit"
          className="btn btn-outline-dark"
          onClick={this.props.handleClick}>
            {this.props.submitText}
        </button>
    </div>;
  }
}
