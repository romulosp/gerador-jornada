import React, { PureComponent } from 'react';

export default class DateInfo extends PureComponent {

  constructor(props) {
    super(props);
    this.state = {
      date: this.getDate()
    };
  }

  componentDidMount = () => {
    this.interval = setInterval(() => this.tick(), 1000);
  }

  componentWillUnmount = () => {
    clearInterval(this.interval);
  }

  getDate() {
    const date = new Date();
    return {
      seconds: this.pad(date.getSeconds(), 2),
      minutes: this.pad(date.getMinutes(), 2),
      hour: this.pad(date.getHours(), 2),
      month: this.pad(date.getMonth() + 1, 2),
      day: this.pad(date.getDate(), 2),
      year: date.getFullYear()
    };
  }

  tick() {
    this.setState({ date: this.getDate() });
  }

  pad(num, size) {
    var s = num + "";
    while (s.length < size) s = "0" + s;
    return s;
  }

  render() {
    return (<React.Fragment>
      <p>BSB - {this.state.date.day}/{this.state.date.month}/{this.state.date.year}</p>
      <p>{this.state.date.hour}:{this.state.date.minutes}:{this.state.date.seconds}</p>
    </React.Fragment>);
  }
}
