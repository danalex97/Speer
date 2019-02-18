import { PacketEntry } from './../logs/Entries';

const Metrics = {
	'load' : (events) => events.filter((x) => x instanceof PacketEntry).length,
};

class Metric {
  // Receives a metric and a computation method that can be applied to a
  // set of event
	constructor(windowSize, computation) {
		this.windowSize = windowSize;
		this.computation = computation;

    // bind functions
    this.data = this.data.bind(this);
	}

	// this recomputes the metric each time
	data(events) {
    let windowEnd = this.windowSize;

    let current   = new Array(); // list of events in current window
    let metrics   = new Array(); // list of metrics

    // each window is exclusive: [windowStart, windowEnd)
    // with size this.windowSize
    for (let event of events) {
      if (event.time >= windowEnd) {
        // update the metrics vector
        let metric = this.computation(current);
        metrics.push(metric);

        // start a new window
        current    = new Array();
        windowEnd += this.windowSize;
      }

      // add new event to window
      current.push(event);
    }
    // handle last window
    if (current.length > 0) {
      let metric = this.computation(current);
      metrics.push(metric);
    }

    return metrics;
	}

	// TODO: from raw data compute distribution
	distribution(events) {
	}
}

export {
  Metrics,
  Metric
};
