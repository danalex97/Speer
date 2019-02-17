import React, { Component } from 'react';
import './app.css';
import BarChart from './charts/BarChart';

export default class App extends Component {
  state = { username: null };

  componentDidMount() {
    fetch('/api/getUsername')
      .then(res => res.json())
      .then(user => this.setState({ username: user.username }));
  }

  render() {
    const { username } = this.state;
    return (
      <div>
        {username ? <h1>{`Hello ${username}`}</h1> : <h1>Loading.. please wait!</h1>}

        <BarChart
            id="chart3"
            data={[1,2,5,2,3,4,5,2]}
            size={[500,500]}
            margin={60}
            dataSize={40}
        />
      </div>
    );
  }
}
