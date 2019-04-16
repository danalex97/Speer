import React, { Component } from 'react';

const SERVER = "http://localhost:8080";

export default class LogFetcher {
  constructor(logName) {
    this.logName = logName;
  }

  fetchLog() {
    return new Promise((resolve, reject) => {
      fetch(`${SERVER}/api/getLog/${this.logName}`)
        .then(res => res.json())
        .then(res => {
          resolve(res);
        });
    });
  }
}
