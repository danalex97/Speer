'use strict';

const e = React.createElement;

class Main extends React.Component {
  render() {
    return (
      <h1>Hello World</h1>
    );
  }
}

const domContainer = document.querySelector('#main');
ReactDOM.render(e(Main), domContainer);
